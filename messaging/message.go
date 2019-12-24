package messaging

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	budget "github.com/jbleduigou/budgetcategorizer"
)

type Broker interface {
	Send(t budget.Transaction) error
}

func NewBroker(queue string) Broker {
	return &sqsbroker{queue: queue}
}

type sqsbroker struct {
	queue string
}

func (b *sqsbroker) Send(t budget.Transaction) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	amount := strconv.FormatFloat(t.Value, 'f', -1, 64)

	result, err := svc.SendMessage(&sqs.SendMessageInput{
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
			// "Comment": &sqs.MessageAttributeValue{
			// 	DataType:    aws.String("String"),
			// 	StringValue: aws.String(t.Comment),
			// },
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
		QueueUrl:    &b.queue,
	})

	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	fmt.Println("Success", *result.MessageId)
	return nil
}
