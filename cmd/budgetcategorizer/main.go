package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/parser"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	parser := parser.NewParser()
	for _, record := range s3Event.Records {
		s3event := record.S3
		objectKey := strings.ReplaceAll(s3event.Object.Key, "input/", "")
		execute(s3event.Bucket.Name, objectKey, downloader, uploader, parser)
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

func execute(bucketName string, objectKey string, downloader s3manageriface.DownloaderAPI, uploader s3manageriface.UploaderAPI, p parser.Parser) {
	//download file
	content, _ := downloadFile(objectKey, bucketName, downloader)
	//read transactions from file
	transactions, _ := p.ParseTransactions(bytes.NewReader(content))
	//write transactions to temp folder
	resultFileName := getResultFileName(objectKey)
	// writeResult(transactions, "/tmp/"+resultFileName)
	output, _ := convertToCSV(transactions)
	fmt.Printf("Output has size %v \n", len(output))
	//	upload to s3
	uploadResult(output, resultFileName, bucketName, uploader)
}

func downloadFile(objectKey string, bucketName string, downloader s3manageriface.DownloaderAPI) ([]byte, error) {
	fmt.Printf("Downloading file '%v' from bucket '%v' \n", objectKey, bucketName)
	buff := &aws.WriteAtBuffer{}
	n, err := downloader.Download(buff, &s3.GetObjectInput{
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

func getResultFileName(fileName string) string {
	resultFileName := []byte(fileName)
	re := regexp.MustCompile(`(\.CSV)`)
	resultFileName = re.ReplaceAll(resultFileName, []byte("-result.txt"))
	return string(resultFileName)
}

func convertToCSV(transactions []*budget.Transaction) ([]byte, error) {
	var b bytes.Buffer
	writer := csv.NewWriter(&b)

	for _, value := range transactions {
		tunasse := strconv.FormatFloat(value.Value, 'f', -1, 64)

		line := []string{value.Date, value.Description, value.Comment, value.Category, string(tunasse)}
		err := writer.Write(line)
		if err != nil {
			fmt.Printf("%v\n", err)
			return nil, err
		}
	}
	writer.Flush()
	return b.Bytes(), nil
}

func uploadResult(result []byte, fileName string, bucketName string, uploader s3manageriface.UploaderAPI) (string, error) {
	r := bytes.NewReader(result)
	objectKey := "output/" + fileName

	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   r,
	}
	o, err := uploader.Upload(upParams)
	if err != nil {
		return "", err
	}
	fmt.Printf("Success uploading file to location %v \n", o.Location)

	return fileName, err
}
