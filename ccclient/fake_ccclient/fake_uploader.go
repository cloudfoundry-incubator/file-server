// This file was generated by counterfeiter
package fake_ccclient

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/cloudfoundry-incubator/file-server/ccclient"
)

type FakeUploader struct {
	UploadStub        func(uploadURL *url.URL, filename string, r *http.Request) (*http.Response, *url.URL, error)
	uploadMutex       sync.RWMutex
	uploadArgsForCall []struct {
		uploadURL *url.URL
		filename  string
		r         *http.Request
	}
	uploadReturns struct {
		result1 *http.Response
		result2 *url.URL
		result3 error
	}
}

func (fake *FakeUploader) Upload(uploadURL *url.URL, filename string, r *http.Request) (*http.Response, *url.URL, error) {
	fake.uploadMutex.Lock()
	fake.uploadArgsForCall = append(fake.uploadArgsForCall, struct {
		uploadURL *url.URL
		filename  string
		r         *http.Request
	}{uploadURL, filename, r})
	fake.uploadMutex.Unlock()
	if fake.UploadStub != nil {
		return fake.UploadStub(uploadURL, filename, r)
	} else {
		return fake.uploadReturns.result1, fake.uploadReturns.result2, fake.uploadReturns.result3
	}
}

func (fake *FakeUploader) UploadCallCount() int {
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	return len(fake.uploadArgsForCall)
}

func (fake *FakeUploader) UploadArgsForCall(i int) (*url.URL, string, *http.Request) {
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	return fake.uploadArgsForCall[i].uploadURL, fake.uploadArgsForCall[i].filename, fake.uploadArgsForCall[i].r
}

func (fake *FakeUploader) UploadReturns(result1 *http.Response, result2 *url.URL, result3 error) {
	fake.UploadStub = nil
	fake.uploadReturns = struct {
		result1 *http.Response
		result2 *url.URL
		result3 error
	}{result1, result2, result3}
}

var _ ccclient.Uploader = new(FakeUploader)