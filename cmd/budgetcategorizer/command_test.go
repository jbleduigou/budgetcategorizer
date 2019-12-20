package main

import (
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
)

type mockDownloader struct {
	downloadCalled            bool
	downloadWithContextCalled bool
}

func (_m *mockDownloader) Download(_a0 io.WriterAt, _a1 *s3.GetObjectInput, _a2 ...func(*s3manager.Downloader)) (int64, error) {
	_m.downloadCalled = true
	_a0.WriteAt([]byte("test"), 0)
	return 1337, nil
}

func (_m *mockDownloader) DownloadWithContext(_a0 aws.Context, _a1 io.WriterAt, _a2 *s3.GetObjectInput, _a3 ...func(*s3manager.Downloader)) (int64, error) {
	_m.downloadWithContextCalled = true
	_a1.WriteAt([]byte("test"), 0)
	return 1337, nil
}

func TestGetResultFileName(t *testing.T) {
	c := &command{downloader: &mockDownloader{}}

	output := c.getResultFileName("CA20191220_1142.CSV")
	assert.Equal(t, "CA20191220_1142-result.txt", output)
}

func TestDownloadFile(t *testing.T) {
	m := &mockDownloader{}
	c := &command{downloader: m}

	content, err := c.downloadFile("CA20191220_1142.CSV", "mybucket")
	assert.Equal(t, []byte("test"), content)
	assert.Nil(t, err)
	assert.Equal(t, true, m.downloadCalled)
	assert.Equal(t, false, m.downloadWithContextCalled)
}
