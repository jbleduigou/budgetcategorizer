package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetResultFileName(t *testing.T) {
	c := &command{}

	output := c.getResultFileName("CA20191220_1142.CSV")
	assert.Equal(t, "CA20191220_1142-result.txt", output)
}

func TestDownloadFile(t *testing.T) {
	m := mock.NewDownloader()
	m.On("Download",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(int64(1337), nil)
	c := &command{downloader: m}

	content, err := c.downloadFile("CA20191220_1142.CSV", "mybucket")
	assert.Equal(t, []byte("test"), content)
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestUploadFile(t *testing.T) {
	m := mock.NewUploader()
	m.On("Upload", mock.Anything, mock.Anything).Return(&s3manager.UploadOutput{Location: ""}, nil)
	c := &command{uploader: m}

	err := c.uploadResult([]byte("test"), "CA20191220_1142-result.txt", "mybucket")
	assert.Nil(t, err)
	m.AssertExpectations(t)
}
