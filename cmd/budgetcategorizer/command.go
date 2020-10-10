package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/messaging"
	"github.com/jbleduigou/budgetcategorizer/parser"
	"go.uber.org/zap"
)

type command struct {
	bucketName  string
	objectKey   string
	downloader  s3manageriface.DownloaderAPI
	parser      parser.Parser
	categorizer categorizer.Categorizer
	broker      messaging.Broker
}

func (c *command) execute() {
	//verify that required variables are defined
	err := c.verifyEnvVariables()
	if err != nil {
		return
	}
	//download file
	content, _ := c.downloadFile(c.objectKey, c.bucketName)
	//read transactions from file
	transactions, _ := c.parser.ParseTransactions(bytes.NewReader(content))
	//categorize transactions
	categorized := mapTransactions(transactions, c.categorizer.Categorize)

	c.broker.Send(categorized)
}

func (c *command) verifyEnvVariables() error {
	_, ok := os.LookupEnv("CONFIGURATION_FILE_BUCKET")
	if !ok {
		zap.S().Warnf("No value defined for variable CONFIGURATION_FILE_BUCKET")
	}
	_, ok = os.LookupEnv("CONFIGURATION_FILE_OBJECT_KEY")
	if !ok {
		zap.S().Warnf("No value defined for variable CONFIGURATION_FILE_OBJECT_KEY")
	}
	_, ok = os.LookupEnv("SQS_QUEUE_URL")
	if !ok {
		zap.S().Errorf("No value defined for variable SQS_QUEUE_URL")
		return fmt.Errorf("No value defined for variable SQS_QUEUE_URL")
	}
	return nil
}

func (c *command) downloadFile(objectKey string, bucketName string) ([]byte, error) {
	zap.S().Infof("Downloading file '%v' from bucket '%v'", objectKey, bucketName)
	buff := &aws.WriteAtBuffer{}
	n, err := c.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		zap.S().Errorf("Failed to download file %v", objectKey, err)
		return nil, err
	}
	zap.S().Infof("File %v downloaded, read %d bytes", objectKey, n)
	return buff.Bytes(), nil
}

func mapTransactions(input []budget.Transaction, f func(budget.Transaction) budget.Transaction) []budget.Transaction {
	output := make([]budget.Transaction, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
