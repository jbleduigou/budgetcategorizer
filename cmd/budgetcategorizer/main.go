package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/config"
	"github.com/jbleduigou/budgetcategorizer/exporter"
	"github.com/jbleduigou/budgetcategorizer/messaging"
	"github.com/jbleduigou/budgetcategorizer/parser"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	// Create all collaborators for command
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	parser := parser.NewParser()
	config := config.GetConfiguration(downloader)
	categorizer := categorizer.NewCategorizer(config.Keywords)
	exporter := exporter.NewExporter()
	for _, record := range s3Event.Records {
		// Retrieve data from S3 event
		s3event := record.S3
		objectKey := strings.ReplaceAll(s3event.Object.Key, "input/", "")
		// Instantiate a command
		c := &command{s3event.Bucket.Name, objectKey, downloader, uploader, parser, categorizer, exporter, messaging.NewBroker(os.Getenv("SQS_QUEUE_URL"))}
		// Execute the command
		c.execute()
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
