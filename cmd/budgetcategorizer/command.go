package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/messaging"
	"github.com/jbleduigou/budgetcategorizer/parser"
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
		slog.Warn("No value defined for variable CONFIGURATION_FILE_BUCKET")
	}
	_, ok = os.LookupEnv("CONFIGURATION_FILE_OBJECT_KEY")
	if !ok {
		slog.Warn("No value defined for variable CONFIGURATION_FILE_OBJECT_KEY")
	}
	_, ok = os.LookupEnv("SQS_QUEUE_URL")
	if !ok {
		slog.Error("No value defined for variable SQS_QUEUE_URL")
		return fmt.Errorf("no value defined for variable SQS_QUEUE_URL")
	}
	return nil
}

func (c *command) downloadFile(objectKey string, bucketName string) ([]byte, error) {
	slog.Info("Downloading file from bucket", "object-key", objectKey, "bucket-name", bucketName)
	buff := &aws.WriteAtBuffer{}
	n, err := c.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		slog.Error("Failed to download file %v", "object-key", objectKey, "bucket-name", bucketName, "error", err)
		return nil, err
	}
	slog.Info("File downloaded with sucess", "object-key", objectKey, "bucket-name", bucketName, "bytes-read", n)
	return buff.Bytes(), nil
}

func mapTransactions(input []budget.Transaction, f func(budget.Transaction) budget.Transaction) []budget.Transaction {
	output := make([]budget.Transaction, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
