package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
)

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

func TestExecute(t *testing.T) {
	d := mock.NewDownloader("")
	d.On("Download",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(int64(1337), nil)
	p := mock.NewParser()
	p.On("ParseTransactions", mock.Anything).Return([]budget.Transaction{}, nil)
	keywords := make(map[string]string)
	cat := categorizer.NewCategorizer(keywords)
	c := &command{downloader: d, parser: p, bucketName: "mybucket", objectKey: "CA20191220_1142.CSV", categorizer: cat}

	c.execute()

	d.AssertExpectations(t)
	p.AssertExpectations(t)
}
