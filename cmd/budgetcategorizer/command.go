package main

import (
	"bytes"
	"fmt"
	"time"

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
	//download file
	content, _ := c.downloadFile(c.objectKey, c.bucketName)
	//read transactions from file
	transactions, _ := c.parser.ParseTransactions(bytes.NewReader(content))
	//categorize transactions
	categorized := mapTransactions(transactions, c.categorizer.Categorize)
	for _, t := range categorized {
		c.broker.Send(t)
		time.Sleep(500 * time.Millisecond)
	}
}

func (c *command) downloadFile(objectKey string, bucketName string) ([]byte, error) {
	fmt.Printf("Downloading file '%v' from bucket '%v' \n", objectKey, bucketName)
	buff := &aws.WriteAtBuffer{}
	n, err := c.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("input/" + objectKey),
	})
	if err != nil {
		fmt.Printf("Failed to download file %v\n, %v", objectKey, err)
		return nil, err
	}
	fmt.Printf("File %v downloaded, read %d bytes\n", objectKey, n)
	return buff.Bytes(), nil
}

func mapTransactions(input []budget.Transaction, f func(budget.Transaction) budget.Transaction) []budget.Transaction {
	output := make([]budget.Transaction, len(input))
	for i, v := range input {
		output[i] = f(v)
	}
	return output
}
