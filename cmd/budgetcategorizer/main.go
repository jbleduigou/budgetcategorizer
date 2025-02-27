package main

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	profile "github.com/jbleduigou/aws-lambda-profile"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/config"
	"github.com/jbleduigou/budgetcategorizer/iface"
	"github.com/jbleduigou/budgetcategorizer/messaging"
	"github.com/jbleduigou/budgetcategorizer/parser"
	slogawslambda "github.com/jbleduigou/slog-aws-lambda"
)

var singleton *config.Configuration
var lock = &sync.Mutex{}

func handleS3Event(ctx context.Context, s3Event events.S3Event) {
	defer profile.Start(profile.S3Bucket(os.Getenv("PROFILING_S3_BUCKET")), profile.AWSRegion(os.Getenv("REGION"))).Stop(ctx)

	// Create all collaborators for command
	initLogger(ctx)
	awscfg, _ := awsconfig.LoadDefaultConfig(ctx)
	parser := parser.NewParser()
	downloader := s3.NewFromConfig(awscfg)

	cfg := getConfig(ctx, downloader)
	categorizer := categorizer.NewCategorizer(cfg.Keywords)
	sqs := sqs.NewFromConfig(awscfg)
	for _, record := range s3Event.Records {
		// Retrieve data from S3 event
		s3event := record.S3
		// Instantiate a command
		c := &command{s3event.Bucket.Name, s3event.Object.Key, downloader, parser, categorizer, messaging.NewBroker(os.Getenv("SQS_QUEUE_URL"), sqs)}
		// Execute the command
		c.execute(ctx)
	}
}

// This function looks for a config in the cache, if not found it will download it from S3
func getConfig(ctx context.Context, downloader iface.S3DownloadAPI) config.Configuration {
	lock.Lock()
	defer lock.Unlock()
	if singleton != nil {
		slog.Info("Using cached configuration")
		return *singleton
	}
	cfg := config.GetConfiguration(ctx, downloader)
	singleton = &cfg
	return cfg
}

func initLogger(ctx context.Context) {
	slog.SetDefault(slog.New(slogawslambda.NewAWSLambdaHandler(ctx, nil, "REGION")))
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handleS3Event)
}
