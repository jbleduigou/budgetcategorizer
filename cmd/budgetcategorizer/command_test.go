package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestDownloadFile(t *testing.T) {
	os.Setenv("SQS_QUEUE_URL", "unit-test")
	m := mock.NewDownloader("test")
	m.On("GetObject",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader("test"))}, nil)
	c := &command{downloader: m}

	content, err := c.downloadFile(context.Background(), "input/CA20191220_1142.CSV", "mybucket")
	assert.Equal(t, []byte("test"), content)
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestDownloadFileWithError(t *testing.T) {
	os.Setenv("SQS_QUEUE_URL", "unit-test")
	m := mock.NewDownloader("")
	m.On("GetObject",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(nil, fmt.Errorf("error for unit test"))
	c := &command{downloader: m}

	content, err := c.downloadFile(context.Background(), "input/CA20191220_1142.CSV", "mybucket")
	assert.Equal(t, []byte(nil), content)
	assert.Equal(t, "error for unit test", err.Error())
	m.AssertExpectations(t)
}

func TestExecute(t *testing.T) {
	os.Setenv("SQS_QUEUE_URL", "unit-test")
	d := mock.NewDownloader("")
	d.On("GetObject",
		mock.Anything,
		&s3.GetObjectInput{Bucket: aws.String("mybucket"), Key: aws.String("input/CA20191220_1142.CSV")},
		mock.Anything).Return(&s3.GetObjectOutput{Body: io.NopCloser(strings.NewReader("test"))}, nil)
	p := mock.NewParser()
	p.On("ParseTransactions", mock.Anything).Return([]budget.Transaction{budget.NewTransaction("19/12/2019", "Paiement Par Carte Express Proxi Saint Thonan 17/12", "", "", 13.37)}, nil)
	keywords := make(map[string]string)
	keywords["Express Proxi Saint Thonan"] = "Courses Alimentation"
	cat := categorizer.NewCategorizer(keywords)
	b := mock.NewBroker()
	b.On("Send", []budget.Transaction{budget.NewTransaction("19/12/2019", "Paiement Par Carte Express Proxi Saint Thonan 17/12", "", "Courses Alimentation", 13.37)}).Return(nil)
	c := &command{downloader: d, parser: p, bucketName: "mybucket", objectKey: "input/CA20191220_1142.CSV", categorizer: cat, broker: b}

	c.execute(context.Background())

	d.AssertExpectations(t)
	p.AssertExpectations(t)
	b.AssertExpectations(t)
}
func TestExecuteMissingSqsEnvVariable(t *testing.T) {
	os.Unsetenv("SQS_QUEUE_URL")
	d := mock.NewDownloader("")
	p := mock.NewParser()
	keywords := make(map[string]string)
	keywords["Express Proxi Saint Thonan"] = "Courses Alimentation"
	cat := categorizer.NewCategorizer(keywords)
	b := mock.NewBroker()
	c := &command{downloader: d, parser: p, bucketName: "mybucket", objectKey: "input/CA20191220_1142.CSV", categorizer: cat, broker: b}

	c.execute(context.Background())

	d.AssertExpectations(t)
	p.AssertExpectations(t)
	b.AssertExpectations(t)
}

func TestMissingSqsEnvVariable(t *testing.T) {
	os.Unsetenv("SQS_QUEUE_URL")
	c := &command{}

	err := c.verifyEnvVariables()

	assert.Equal(t, "no value defined for variable SQS_QUEUE_URL", err.Error())
}
