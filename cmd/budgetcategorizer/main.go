package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	budget "github.com/jbleduigou/budgetcategorizer"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	for _, record := range s3Event.Records {
		s3event := record.S3
		objectKey := strings.ReplaceAll(s3event.Object.Key, "input/", "")
		execute(s3event.Bucket.Name, objectKey, downloader, uploader)
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

func execute(bucketName string, objectKey string, downloader *s3manager.Downloader, uploader *s3manager.Uploader) {
	//download file
	content, _ := downloadFile(objectKey, bucketName, downloader)
	//read transactions from file
	transactions := readTransactions(bytes.NewReader(content))
	//write transactions to temp folder
	resultFileName := getResultFileName(objectKey)
	writeResult(transactions, "/tmp/"+resultFileName)
	//	upload to s3
	uploadResult(resultFileName, bucketName, uploader)
}

func downloadFile(objectKey string, bucketName string, downloader *s3manager.Downloader) ([]byte, error) {
	{
		buff := &aws.WriteAtBuffer{}
		n, err := downloader.Download(buff, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String("input/" + objectKey),
		})
		if err != nil {
			fmt.Printf("failed to download file\n, %v", err)
			return nil, err
		}
		fmt.Printf("file downloaded, %d bytes\n", n)
		return buff.Bytes(), nil
	}
}

func readTransactions(r io.Reader) (transactions []*budget.Transaction) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.FieldsPerRecord = 4

	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, each := range rawCSVdata {
		if len(each) == 4 {
			date := each[0]
			libelle := sanitizeDescription(each[1])
			debit, err := parseAmount(each[2])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			t := budget.NewTransaction(date, libelle, "", "Courses Alimentation", debit)
			transactions = append(transactions, t)
		}
		if len(each) == 5 {
			date := each[0]
			if "Date" != date {
				libelle := sanitizeDescription(each[1])
				credit, err := parseAmount(each[3])
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				t := budget.NewTransaction(date, libelle, "", "", -credit)
				transactions = append(transactions, t)
			}
		}
	}
	fmt.Printf("Found %v transactions\n", len(transactions))
	return transactions
}

func getResultFileName(fileName string) string {
	resultFileName := []byte(fileName)
	re := regexp.MustCompile(`(\.CSV)`)
	resultFileName = re.ReplaceAll(resultFileName, []byte("-result.txt"))
	return string(resultFileName)
}

func sanitizeDescription(d string) string {
	libelle := []byte(d)
	{
		re := regexp.MustCompile(`\n`)
		libelle = re.ReplaceAll(libelle, []byte(" "))
	}
	{
		re := regexp.MustCompile(`[\s]+`)
		libelle = re.ReplaceAll(libelle, []byte(" "))
	}
	return string(libelle)
}

func parseAmount(a string) (float64, error) {
	creditStr := []byte(a)
	{
		re := regexp.MustCompile(`,`)
		creditStr = re.ReplaceAll(creditStr, []byte("."))
	}
	credit, err := strconv.ParseFloat(string(creditStr), 64)
	return credit, err
}

func writeResult(transactions []*budget.Transaction, fileName string) {
	fmt.Printf("Writing result to file %v \n", fileName)
	file, err := os.Create(fileName)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range transactions {
		tunasse := strconv.FormatFloat(value.Value, 'f', -1, 64)

		line := []string{value.Date, value.Description, value.Comment, value.Category, string(tunasse)}
		err := writer.Write(line)
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func uploadResult(fileName string, bucketName string, uploader *s3manager.Uploader) (string, error) {
	file, err := os.Open("/tmp/" + fileName)
	checkError("Cannot open file", err)
	defer file.Close()

	objectKey := "output/" + fileName

	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   file,
	}
	o, err := uploader.Upload(upParams)
	if err != nil {
		return "", err
	}
	fmt.Printf("Success uploading file to location %v \n", o.Location)

	return fileName, err
}
