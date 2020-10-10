package messaging

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
)

func TestSendSuccess(t *testing.T) {
	m := mock.NewSQSClient()

	m.On("SendMessageBatch", testify.MatchedBy(func(req *sqs.SendMessageBatchInput) bool { return len(req.Entries) == 10 })).Return(getSendMessageBatchOutput(), nil)
	m.On("SendMessageBatch", testify.MatchedBy(func(req *sqs.SendMessageBatchInput) bool { return len(req.Entries) == 1 })).Return(getSendMessageBatchOutput(), nil)

	b := NewBroker("https://sqs.eu-west-3.amazonaws.com/959789434/testing", m)

	list := make([]budget.Transaction, 11)

	err := b.Send(list)

	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestSendError(t *testing.T) {
	m := mock.NewSQSClient()

	m.On("SendMessageBatch", testify.MatchedBy(func(req *sqs.SendMessageBatchInput) bool { return len(req.Entries) == 5 })).Return(getSendMessageBatchOutput(), fmt.Errorf("Error for unit tests"))

	b := NewBroker("https://sqs.eu-west-3.amazonaws.com/959789434/testing", m)

	list := make([]budget.Transaction, 5)

	err := b.Send(list)

	assert.Equal(t, err.Error(), "Error for unit tests")
	m.AssertExpectations(t)
}

func getSendMessageBatchOutput() *sqs.SendMessageBatchOutput {
	messageID := "5aab4335-527a-41d5-bba6-7e5cfdb8228d"
	return &sqs.SendMessageBatchOutput{
		Successful: []*sqs.SendMessageBatchResultEntry{&sqs.SendMessageBatchResultEntry{MessageId: &messageID}},
		Failed:     []*sqs.BatchResultErrorEntry{&sqs.BatchResultErrorEntry{Id: &messageID, Message: &messageID}},
	}
}
