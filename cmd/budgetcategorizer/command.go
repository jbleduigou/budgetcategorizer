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
)

type command struct {
	bucketName  string
	objectKey   string
	downloader  s3manageriface.DownloaderAPI
	parser      parser.Parser
	categorizer categorizer.Categorizer
	broker      messaging.Broker
	requestID   string
}

func (c *command) execute() {
	//verify that required variable are defined
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
	for _, t := range categorized {
		c.broker.Send(t)
	}
}

func (c *command) verifyEnvVariables() error {
	_, ok := os.LookupEnv("CONFIGURATION_FILE_BUCKET")
	if !ok {
		fmt.Printf("[WARN] %v No value defined for variable CONFIGURATION_FILE_BUCKET\n", c.requestID)
	}
	_, ok = os.LookupEnv("CONFIGURATION_FILE_OBJECT_KEY")
	if !ok {
		fmt.Printf("[WARN] %v No value defined for variable CONFIGURATION_FILE_OBJECT_KEY\n", c.requestID)
	}
	_, ok = os.LookupEnv("SQS_QUEUE_URL")
	if !ok {
		fmt.Printf("[ERROR] %v No value defined for variable SQS_QUEUE_URL\n", c.requestID)
		return fmt.Errorf("No value defined for variable SQS_QUEUE_URL")
	}
	return nil
}

func (c *command) downloadFile(objectKey string, bucketName string) ([]byte, error) {
	fmt.Printf("[INFO] %v Downloading file '%v' from bucket '%v' \n", c.requestID, objectKey, bucketName)
	buff := &aws.WriteAtBuffer{}
	n, err := c.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fmt.Printf("[ERROR] %v Failed to download file %v, %v\n", c.requestID, objectKey, err)
		return nil, err
	}
	fmt.Printf("[INFO] %v File %v downloaded, read %d bytes\n", c.requestID, objectKey, n)
	return buff.Bytes(), nil
}

func mapTransactions(input []budget.Transaction, f func(budget.Transaction) budget.Transaction) []budget.Transaction {
	output := make([]budget.Transaction, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
