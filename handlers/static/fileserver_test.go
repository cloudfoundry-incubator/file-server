package static_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"code.cloudfoundry.org/fileserver/handlers/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileServer", func() {
	var (
		servedDirectory string
		fileServer      *httptest.Server
	)

	BeforeEach(func() {
		var err error
		servedDirectory, err = ioutil.TempDir("", "fileserver-test")
		Expect(err).NotTo(HaveOccurred())
		os.Mkdir(filepath.Join(servedDirectory, "testdir"), os.ModePerm)

		tenHoursAgo := time.Now().Add(-10 * time.Hour)

		ioutil.WriteFile(filepath.Join(servedDirectory, "test"), []byte("hello"), os.ModePerm)
		ioutil.WriteFile(filepath.Join(servedDirectory, "test.sha1"), []byte("some-hash\n"), os.ModePerm)
		os.Chtimes(filepath.Join(servedDirectory, "test"), tenHoursAgo, tenHoursAgo)

		ioutil.WriteFile(filepath.Join(servedDirectory, "test2.."), []byte("world"), os.ModePerm)
		ioutil.WriteFile(filepath.Join(servedDirectory, "test2...sha1"), []byte("some-hash-2\n"), os.ModePerm)

		ioutil.WriteFile(filepath.Join(servedDirectory, "no-sha"), []byte("hello"), os.ModePerm)

		fileServer = httptest.NewServer(static.NewFileServer(servedDirectory))
	})

	AfterEach(func() {
		fileServer.Close()
		os.RemoveAll(servedDirectory)
	})

	Context("when the file and its sha1 file exists", func() {
		It("returns a 200 OK and the file and its ETag as the sha1sum", func() {
			resp, err := http.Get(fmt.Sprintf("%s/test", fileServer.URL))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("ETag")).To(Equal(fmt.Sprintf(`"%s"`, "some-hash")))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("hello"))
		})

		Context("when the file name contains dot dot", func() {
			It("returns a 200 OK and the file and its ETag as the sha1sum", func() {
				resp, err := http.Get(fmt.Sprintf("%s/test2..", fileServer.URL))
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("ETag")).To(Equal(fmt.Sprintf(`"%s"`, "some-hash-2")))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("world"))
			})
		})

		Context("when the request provides an 'If-None-Match' header that matches the sha1sum", func() {
			It("returns a 304 Not Modfied and does not return the file", func() {
				req, err := http.NewRequest("GET", fmt.Sprintf("%s/test", fileServer.URL), nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("If-None-Match", `"some-hash"`)

				resp, err := http.DefaultClient.Do(req)
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotModified))
				Expect(resp.Header.Get("ETag")).To(Equal(fmt.Sprintf(`"%s"`, "some-hash")))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal(""))

			})
		})

		Context("when the request provides an 'If-None-Match' header that does not match the sha1sum", func() {
			It("returns a 200 OK and the file even if the file had an older mtime than the time in If-Modified-Since", func() {
				req, err := http.NewRequest("GET", fmt.Sprintf("%s/test", fileServer.URL), nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("If-None-Match", `"different-hash"`)
				req.Header.Set("If-Modified-Since", time.Now().Add(-2*time.Hour).Format(http.TimeFormat))

				resp, err := http.DefaultClient.Do(req)
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Header.Get("ETag")).To(Equal(fmt.Sprintf(`"%s"`, "some-hash")))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("hello"))
			})
		})
	})

	It("returns 400 on filepaths with dot dot", func() {
		resp, err := http.Get(fmt.Sprintf("%s/../protected-file", fileServer.URL))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("returns 404 on files that don't exist", func() {
		resp, err := http.Get(fmt.Sprintf("%s/does-not-exist", fileServer.URL))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("returns 404 when the file exists but the sha file is missing", func() {
		resp, err := http.Get(fmt.Sprintf("%s/no-sha", fileServer.URL))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("returns 401 when accessing a directory", func() {
		resp, err := http.Get(fmt.Sprintf("%s/testdir", fileServer.URL))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

})
