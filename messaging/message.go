package messaging

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	budget "github.com/jbleduigou/budgetcategorizer"
	"go.uber.org/zap"
)

// Broker provides an interface for sending messaging to a queue
type Broker interface {
	Send(t budget.Transaction) error
}

// NewBroker will provide an instance of a Broker, implementation is not exposed
func NewBroker(queue string, svc sqsiface.SQSAPI) Broker {
	return &sqsbroker{queueURL: queue, svc: svc}
}

type sqsbroker struct {
	queueURL string
	svc      sqsiface.SQSAPI
}

func (b *sqsbroker) Send(t budget.Transaction) error {
	amount := strconv.FormatFloat(t.Value, 'f', -1, 64)

	payload, _ := json.Marshal(t)
	//Decided to ignore error, it is only returned is case of:
	// - UnsupportedTypeError: Channel, complex, and function values
	// - UnsupportedValueError: cyclic data structures

	result, err := b.svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Date": {
				DataType:    aws.String("String"),
				StringValue: aws.String(t.Date),
			},
			"Description": {
				DataType:    aws.String("String"),
				StringValue: aws.String(t.Description),
			},
			"Category": {
				DataType:    aws.String("String"),
				StringValue: aws.String(t.Category),
			},
			"Value": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(amount),
			},
		},
		MessageBody: aws.String(string(payload)),
		QueueUrl:    &b.queueURL,
	})

	if err != nil {
		zap.S().Errorf("Error while sending message", err)
		return err
	}

	zap.S().Infof("Message send successfully with id '%v'", *result.MessageId)
	return nil
}
