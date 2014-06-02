package integration_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/cloudfoundry-incubator/file-server/integration/fileserver_runner"
	Bbs "github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/router"
	steno "github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/gunk/urljoiner"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type ByteEmitter struct {
	written int
	length  int
}

func NewEmitter(length int) *ByteEmitter {
	return &ByteEmitter{
		length:  length,
		written: 0,
	}
}

func (emitter *ByteEmitter) Read(p []byte) (n int, err error) {
	if emitter.written >= emitter.length {
		return 0, io.EOF
	}
	time.Sleep(time.Millisecond)
	p[0] = 0xF1
	emitter.written++
	return 1, nil
}

var _ = Describe("File server", func() {
	var (
		bbs     *Bbs.BBS
		address string
		port    int
		baseUrl string
		err     error
		appGuid = "app-guid"
		runner  *fileserver_runner.Runner
	)

	dropletUploadRequest := func(appGuid string, body io.Reader, contentLength int) *http.Request {
		route, ok := router.NewFileServerRoutes().RouteForHandler(router.FS_UPLOAD_DROPLET)
		Ω(ok).Should(BeTrue())

		path, err := route.PathWithParams(map[string]string{"guid": appGuid})
		Ω(err).ShouldNot(HaveOccurred())

		postRequest, err := http.NewRequest("POST", urljoiner.Join(baseUrl, path), body)
		Ω(err).ShouldNot(HaveOccurred())
		postRequest.ContentLength = int64(contentLength)
		postRequest.Header.Set("Content-Type", "application/octet-stream")

		return postRequest
	}

	BeforeEach(func() {
		logSink := steno.NewTestingSink()
		steno.Init(&steno.Config{
			Sinks: []steno.Sink{logSink},
		})
		logger := steno.NewLogger("the-logger")
		steno.EnterTestMode()

		bbs = Bbs.NewBBS(etcdRunner.Adapter(), timeprovider.NewTimeProvider(), logger)

		address = "localhost"
		port = 8182 + config.GinkgoConfig.ParallelNode
		baseUrl = fmt.Sprintf("http://%s:%d/", address, port)

		runner = fileserver_runner.New(fileServerBinary, port, etcdRunner.NodeURLS(), fakeCC.Address(), fakeCC.Username(), fakeCC.Password())
	})

	AfterEach(func() {
		runner.KillWithFire()
	})

	Context("when file server exits", func() {
		It("should remove its presence", func() {
			runner.Start("-address", address, "-heartbeatInterval", "10s")

			_, err = bbs.GetAvailableFileServer()
			Ω(err).ShouldNot(HaveOccurred())

			runner.Stop()

			_, err = bbs.GetAvailableFileServer()
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when it fails to maintain presence", func() {
		BeforeEach(func() {
			runner.Start("-address", address, "-heartbeatInterval", "1s")
		})

		It("should retry", func() {
			_, err := bbs.GetAvailableFileServer()
			Ω(err).ShouldNot(HaveOccurred())

			etcdRunner.Stop()

			Eventually(func() error {
				_, err := bbs.GetAvailableFileServer()
				return err
			}).Should(HaveOccurred())

			Consistently(runner, 1).ShouldNot(gexec.Exit())

			etcdRunner.Start()

			Eventually(func() error {
				_, err := bbs.GetAvailableFileServer()
				return err
			}, 3).ShouldNot(HaveOccurred())
		})
	})

	Context("when started correctly", func() {
		BeforeEach(func() {
			runner.Start("-address", address)
		})

		It("should maintain presence in ETCD", func() {
			fileServerURL, err := bbs.GetAvailableFileServer()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(fileServerURL).Should(Equal(baseUrl))
		})

		Context("and serving a file", func() {
			BeforeEach(func() {
				runner.CreateAndServeFile("test", bytes.NewBufferString("hello"))
			})

			It("should return that file on GET request", func() {
				resp, err := http.Get(urljoiner.Join(baseUrl, "/v1/static/test"))
				Ω(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()

				Ω(resp.StatusCode).Should(Equal(http.StatusOK))

				body, err := ioutil.ReadAll(resp.Body)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(string(body)).Should(Equal("hello"))
			})
		})
	})

	Context("when an address is not specified", func() {
		BeforeEach(func() {
			runner.Start()
		})

		It("publishes its url properly", func() {
			fileServerURL, err := bbs.GetAvailableFileServer()
			Ω(err).ShouldNot(HaveOccurred())

			serverURL, err := url.Parse(fileServerURL)
			Ω(err).ShouldNot(HaveOccurred())

			host, _, err := net.SplitHostPort(serverURL.Host)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(host).ShouldNot(Equal(""))
		})
	})

	Describe("uploading a file", func() {
		var contentLength = 100

		BeforeEach(func() {
			runner.Start("-address", address)
		})

		It("should upload the file...", func(done Done) {
			emitter := NewEmitter(contentLength)
			postRequest := dropletUploadRequest(appGuid, emitter, contentLength)
			resp, err := http.DefaultClient.Do(postRequest)
			Ω(err).ShouldNot(HaveOccurred())
			defer resp.Body.Close()

			Ω(resp.StatusCode).Should(Equal(http.StatusCreated))
			Ω(len(fakeCC.UploadedDroplets[appGuid])).Should(Equal(contentLength))
			close(done)
		}, 2.0)
	})

	// BUG(tedsuo): we appear to be unable to test fileserver drain
	XDescribe("when the fileserver receives SIGINT", func() {
		var sendStarted chan struct{}

		BeforeEach(func() {
			runtime.GOMAXPROCS(8)
			runner.Start("-address", address)
			sendStarted = make(chan struct{})
			go func() {
				defer GinkgoRecover()
				<-sendStarted
				time.Sleep(1000 * time.Millisecond)
				println("******** INTERRUPT ********")
				runner.Stop()
			}()
		})

		Describe("and file requests are in flight", func() {
			var contentLength = 100000

			It("completes in-flight file requests", func(done Done) {
				close(sendStarted)
				emitter := NewEmitter(contentLength)
				postRequest := dropletUploadRequest(appGuid, emitter, contentLength)

				client := http.Client{
					Transport: &http.Transport{},
				}

				resp, err := client.Do(postRequest)
				Ω(err).ShouldNot(HaveOccurred())
				resp.Body.Close()

				Ω(resp.StatusCode).Should(Equal(http.StatusCreated))
				Ω(len(fakeCC.UploadedDroplets[appGuid])).Should(Equal(contentLength))

				close(done)
			}, 8.0)
		})
	})
})
