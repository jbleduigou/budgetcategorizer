package messaging

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	budget "github.com/jbleduigou/budgetcategorizer"
)

type Broker interface {
	Send(t budget.Transaction) error
}

func NewBroker(queue string, svc sqsiface.SQSAPI) Broker {
	return &sqsbroker{queueURL: queue, svc: svc}
}

type sqsbroker struct {
	queueURL string
	svc      sqsiface.SQSAPI
}

func (b *sqsbroker) Send(t budget.Transaction) error {
	amount := strconv.FormatFloat(t.Value, 'f', -1, 64)

	result, err := b.svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Date": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(t.Date),
			},
			"Description": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(t.Description),
			},
			"Category": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(t.Category),
			},
			"Value": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(amount),
			},
		},
		MessageBody: aws.String(t.Description),
		QueueUrl:    &b.queueURL,
	})

	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	fmt.Println("Success", *result.MessageId)
	return nil
}
