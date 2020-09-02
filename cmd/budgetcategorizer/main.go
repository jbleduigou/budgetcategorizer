package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jbleduigou/budgetcategorizer/categorizer"
	"github.com/jbleduigou/budgetcategorizer/config"
	"github.com/jbleduigou/budgetcategorizer/messaging"
	"github.com/jbleduigou/budgetcategorizer/parser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func handleS3Event(ctx context.Context, s3Event events.S3Event) {
	// Create all collaborators for command
	initLogger(ctx)
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	parser := parser.NewParser()
	config := config.GetConfiguration(downloader)
	categorizer := categorizer.NewCategorizer(config.Keywords)
	sqs := sqs.New(sess)
	for _, record := range s3Event.Records {
		// Retrieve data from S3 event
		s3event := record.S3
		// Instantiate a command
		c := &command{s3event.Bucket.Name, s3event.Object.Key, downloader, parser, categorizer, messaging.NewBroker(os.Getenv("SQS_QUEUE_URL"), sqs)}
		// Execute the command
		c.execute()
	}
}

func initLogger(ctx context.Context) {
	// Retrieve AWS Request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := lc.AwsRequestID
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"request-id": requestID},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync() // flushes buffer, if any
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handleS3Event)
}
