package messaging

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/google/uuid"
	budget "github.com/jbleduigou/budgetcategorizer"
)

// Broker provides an interface for sending messaging to a queue
type Broker interface {
	Send(t []budget.Transaction) error
}

// NewBroker will provide an instance of a Broker, implementation is not exposed
func NewBroker(queue string, svc sqsiface.SQSAPI) Broker {
	return &sqsbroker{queueURL: queue, svc: svc}
}

type sqsbroker struct {
	queueURL string
	svc      sqsiface.SQSAPI
}

func (b *sqsbroker) Send(list []budget.Transaction) error {
	batchSize := 10
	for start := 0; start < len(list); start += batchSize {
		end := start + batchSize
		if end > len(list) {
			end = len(list)
		}
		err := b.sendBatch(list[start:end])
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *sqsbroker) sendBatch(list []budget.Transaction) error {
	request := &sqs.SendMessageBatchInput{QueueUrl: &b.queueURL}
	entries := make([]*sqs.SendMessageBatchRequestEntry, len(list))
	for i, t := range list {

		//Decided to ignore error, it is only returned is case of:
		// - UnsupportedTypeError: Channel, complex, and function values
		// - UnsupportedValueError: cyclic data structures
		payload, _ := json.Marshal(t)

		m := &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(uuid.New().String()),
			MessageBody: aws.String(string(payload))}
		entries[i] = m
	}
	request.SetEntries(entries)
	slog.Info("Sending a batch of messages to SQS",
		slog.Int("batch-size", len(list)))
	result, err := b.svc.SendMessageBatch(request)
	if err != nil {
		slog.Error("Error while sending message", "error", err)
		return err
	}
	for _, s := range result.Successful {
		slog.Info("Message send successfully to SQS",
			slog.String("message-id", aws.StringValue(s.MessageId)))
	}
	for _, f := range result.Failed {
		slog.Warn("Error while sending message to SQS",
			slog.String("error-message", aws.StringValue(f.Message)),
			slog.String("message-id", aws.StringValue(f.Id)))
	}
	return nil
}
