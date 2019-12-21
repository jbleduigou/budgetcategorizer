package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/jbleduigou/budgetcategorizer/exporter"
	"github.com/jbleduigou/budgetcategorizer/parser"
)

type command struct {
	bucketName string
	objectKey  string
	downloader s3manageriface.DownloaderAPI
	uploader   s3manageriface.UploaderAPI
	p          parser.Parser
	e          exporter.Exporter
}

func (c *command) execute() {
	//download file
	content, _ := c.downloadFile(c.objectKey, c.bucketName)
	//read transactions from file
	transactions, _ := c.p.ParseTransactions(bytes.NewReader(content))
	//write transactions to temp folder
	resultFileName := c.getResultFileName(c.objectKey)
	//export transactions to proper csv format
	output, _ := c.e.Export(transactions)
	fmt.Printf("Output has size %v \n", len(output))
	//	upload to s3
	c.uploadResult(output, resultFileName, c.bucketName)
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

func (c *command) getResultFileName(fileName string) string {
	resultFileName := []byte(fileName)
	re := regexp.MustCompile(`(\.CSV)`)
	resultFileName = re.ReplaceAll(resultFileName, []byte("-result.txt"))
	return string(resultFileName)
}

func (c *command) uploadResult(result []byte, fileName string, bucketName string) error {
	r := bytes.NewReader(result)
	objectKey := "output/" + fileName

	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   r,
	}
	o, err := c.uploader.Upload(upParams)
	if err != nil {
		return err
	}
	fmt.Printf("Success uploading file to location %v \n", o.Location)
	return nil
}
