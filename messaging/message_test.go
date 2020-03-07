package messaging

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
)

func TestSendSuccess(t *testing.T) {
	m := mock.NewSQSClient()
	messageID := "67043d7c-49db-43bb-af88-d0c30d62234f"
	request := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Date": {
				DataType:    aws.String("String"),
				StringValue: aws.String("<date>"),
			},
			"Description": {
				DataType:    aws.String("String"),
				StringValue: aws.String("<description>"),
			},
			"Category": {
				DataType:    aws.String("String"),
				StringValue: aws.String("<category>"),
			},
			"Value": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("13.37"),
			},
		},
		MessageBody: aws.String("<description>"),
		QueueUrl:    aws.String("https://sqs.eu-west-3.amazonaws.com/959789434/testing"),
	}

	m.On("SendMessage", request).Return(&sqs.SendMessageOutput{MessageId: &messageID}, nil)

	b := NewBroker("https://sqs.eu-west-3.amazonaws.com/959789434/testing", m)

	err := b.Send(budget.NewTransaction("<date>", "<description>", "<comment>", "<category>", 13.37))

	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestSendError(t *testing.T) {
	m := mock.NewSQSClient()
	request := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Date": {
				DataType:    aws.String("String"),
				StringValue: aws.String("<date>"),
			},
			"Description": {
				DataType:    aws.String("String"),
				StringValue: aws.String("<description>"),
			},
			"Category": {
				DataType:    aws.String("String"),
				StringValue: aws.String("<category>"),
			},
			"Value": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("13.37"),
			},
		},
		MessageBody: aws.String("<description>"),
		QueueUrl:    aws.String("https://sqs.eu-west-3.amazonaws.com/959789434/testing"),
	}

	m.On("SendMessage", request).Return(nil, fmt.Errorf("Error for unit tests"))

	b := NewBroker("https://sqs.eu-west-3.amazonaws.com/959789434/testing", m)

	err := b.Send(budget.NewTransaction("<date>", "<description>", "<comment>", "<category>", 13.37))

	assert.Equal(t, err.Error(), "Error for unit tests")
	m.AssertExpectations(t)
}
