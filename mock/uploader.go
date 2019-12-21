package mock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/mock"
)

func NewUploader() *MockUploader {
	return &MockUploader{}
}

type MockUploader struct {
	mock.Mock
}

func (_m *MockUploader) Upload(_a0 *s3manager.UploadInput, _a1 ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	ret := _m.Called(_a0, _a1)
	if ret.Get(1) != nil {
		return ret.Get(0).(*s3manager.UploadOutput), ret.Get(1).(error)
	}
	return ret.Get(0).(*s3manager.UploadOutput), nil
}

func (_m *MockUploader) UploadWithContext(_a0 aws.Context, _a1 *s3manager.UploadInput, _a2 ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	_m.Called(_a0, _a1, _a2)
	return nil, nil
	// return ret.Get(0).(*s3manager.UploadOutput), nil
}
