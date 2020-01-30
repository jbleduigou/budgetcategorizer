package mock

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/mock"
)

const (
	// Anything is used in Diff and Assert when the argument being tested
	// shouldn't be taken into consideration.
	Anything = "mock.Anything"
)

// NewDownloader provides a mock instance of a Downloader
func NewDownloader(content string) *Downloader {
	return &Downloader{content: content}
}

// Downloader is an implementation of the Downloader interface with a mock, use for testing not for production
type Downloader struct {
	mock.Mock
	content string
}

// Download downloads an object using bucket name and object key provided
func (_m *Downloader) Download(_a0 io.WriterAt, _a1 *s3.GetObjectInput, _a2 ...func(*s3manager.Downloader)) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2)
	_a0.WriteAt([]byte(_m.content), 0)
	if ret.Get(1) != nil {
		return ret.Get(0).(int64), ret.Get(1).(error)
	}
	return ret.Get(0).(int64), nil
}

// DownloadWithContext is the same as Download with the additional support for Context input parameters.
func (_m *Downloader) DownloadWithContext(_a0 aws.Context, _a1 io.WriterAt, _a2 *s3.GetObjectInput, _a3 ...func(*s3manager.Downloader)) (int64, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)
	_a1.WriteAt([]byte(_m.content), 0)
	return ret.Get(0).(int64), ret.Get(1).(error)
}
