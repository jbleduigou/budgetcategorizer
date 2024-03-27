package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/mock"
)

const (
	// Anything is used in Diff and Assert when the argument being tested
	// shouldn't be taken into consideration.
	Anything = "mock.Anything"
)

// NewDownloader provides a mock instance of a Downloader
func NewDownloader(content string) *Downloader {
	return &Downloader{}
}

// Downloader is an implementation of the Downloader interface with a mock, use for testing not for production
type Downloader struct {
	mock.Mock
}

func (_m *Downloader) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	ret := _m.Called(ctx, input, optFns)
	if ret.Get(1) != nil {
		return nil, ret.Get(1).(error)
	}
	return ret.Get(0).(*s3.GetObjectOutput), nil
}
