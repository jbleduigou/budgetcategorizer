package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetResultFileName(t *testing.T) {
	c := &command{}

	output := c.getResultFileName("CA20191220_1142.CSV")
	assert.Equal(t, "CA20191220_1142-result.txt", output)
}

func TestDownloadFile(t *testing.T) {
	m := mock.NewDownloader("test")
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

func TestDownloadFileWithError(t *testing.T) {
	m := mock.NewDownloader("")
	m.On("Download",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(int64(0), fmt.Errorf("error for unit test"))
	c := &command{downloader: m}

	content, err := c.downloadFile("CA20191220_1142.CSV", "mybucket")
	assert.Equal(t, []byte(nil), content)
	assert.Equal(t, "error for unit test", err.Error())
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

func TestUploadFileWithError(t *testing.T) {
	m := mock.NewUploader()
	m.On("Upload", mock.Anything, mock.Anything).Return(&s3manager.UploadOutput{Location: ""}, fmt.Errorf("error for unit test"))
	c := &command{uploader: m}

	err := c.uploadResult([]byte("test"), "CA20191220_1142-result.txt", "mybucket")
	assert.Equal(t, "error for unit test", err.Error())
	m.AssertExpectations(t)
}

//TODO use better return values and arguments for mocks
func TestExecute(t *testing.T) {
	d := mock.NewDownloader("")
	d.On("Download",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(int64(1337), nil)
	p := mock.NewParser()
	p.On("ParseTransactions", mock.Anything).Return([]*budget.Transaction{}, nil)
	e := mock.NewExporter()
	e.On("Export", mock.Anything, mock.Anything).Return([]byte(""), nil)
	u := mock.NewUploader()
	u.On("Upload", mock.Anything, mock.Anything).Return(&s3manager.UploadOutput{Location: ""}, nil)
	keywords := make(map[string]string)
	cat := categorizer.NewCategorizer(keywords)
	c := &command{downloader: d, parser: p, exporter: e, uploader: u, bucketName: "mybucket", objectKey: "CA20191220_1142.CSV", categorizer: cat}

	c.execute()

	d.AssertExpectations(t)
	p.AssertExpectations(t)
	e.AssertExpectations(t)
	u.AssertExpectations(t)
}
