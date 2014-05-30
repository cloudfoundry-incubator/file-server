package fileserver_runner

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type Runner struct {
	fileServerBin string
	etcdCluster   []string
	dir           string
	port          int
	Session       *gexec.Session
	ccAddress     string
	ccUsername    string
	ccPassword    string
}

func New(fileServerBin string, port int, etcdCluster []string, ccAddress, ccUsername, ccPassword string) *Runner {
	return &Runner{
		fileServerBin: fileServerBin,
		etcdCluster:   etcdCluster,
		port:          port,
		ccAddress:     ccAddress,
		ccUsername:    ccUsername,
		ccPassword:    ccPassword,
	}
}

func (r *Runner) Start(extras ...string) {
	r.StartWithoutCheck(extras...)

	Eventually(func() int {
		resp, _ := http.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/static/ready", r.port))
		if resp != nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, 4.0).Should(Equal(http.StatusOK))
}

func (r *Runner) StartWithoutCheck(extras ...string) {
	if r.Session != nil {
		panic("starting an already started fileserver runner!!!")
	}

	tempDir, err := ioutil.TempDir("", "inigo-file-server")
	Ω(err).ShouldNot(HaveOccurred())

	r.dir = tempDir

	ioutil.WriteFile(filepath.Join(r.dir, "ready"), []byte("ready"), os.ModePerm)

	args := append(
		extras,
		"-staticDirectory", r.dir,
		"-port", fmt.Sprintf("%d", r.port),
		"-etcdCluster", strings.Join(r.etcdCluster, ","),
		"-ccAddress", r.ccAddress,
		"-ccUsername", r.ccUsername,
		"-ccPassword", r.ccPassword,
		"-skipCertVerify",
	)

	executorSession, err := gexec.Start(exec.Command(r.fileServerBin, args...), ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())

	r.Session = executorSession
}

func (r *Runner) ExitCode() int {
	return r.Session.ExitCode()
}

func (r *Runner) CreateAndServeFile(name string, reader io.Reader) {
	data, err := ioutil.ReadAll(reader)
	Ω(err).ShouldNot(HaveOccurred())

	ioutil.WriteFile(filepath.Join(r.dir, name), data, os.ModePerm)
}

func (r *Runner) ServeFile(name string, path string) {
	data, err := ioutil.ReadFile(path)
	Ω(err).ShouldNot(HaveOccurred())

	ioutil.WriteFile(filepath.Join(r.dir, name), data, os.ModePerm)
}

func (r *Runner) StaticDir() string {
	return r.dir
}

func (r *Runner) Stop() {
	if r.Session != nil {
		r.Session.Interrupt().Wait(5 * time.Second)

		// TODO: kill static dir
		r.Session = nil
	}
}

func (r *Runner) KillWithFire() {
	if r.Session != nil {
		r.Session.Kill().Wait(5 * time.Second)

		// TODO: kill static dir
		r.Session = nil
	}
}
