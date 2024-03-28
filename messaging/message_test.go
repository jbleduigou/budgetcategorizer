package messaging

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/jbleduigou/budgetcategorizer/mock"
	"github.com/stretchr/testify/assert"
	testify "github.com/stretchr/testify/mock"
)

func TestSendSuccess(t *testing.T) {
	m := mock.NewSQSClient()

	m.On("SendMessageBatch", mock.Anything, testify.MatchedBy(func(req *sqs.SendMessageBatchInput) bool { return len(req.Entries) == 10 }), mock.Anything).Return(getSendMessageBatchOutput(), nil)
	m.On("SendMessageBatch", mock.Anything, testify.MatchedBy(func(req *sqs.SendMessageBatchInput) bool { return len(req.Entries) == 1 }), mock.Anything).Return(getSendMessageBatchOutput(), nil)

	b := NewBroker("https://sqs.eu-west-3.amazonaws.com/959789434/testing", m)

	list := make([]budget.Transaction, 11)

	err := b.Send(list)

	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestSendError(t *testing.T) {
	m := mock.NewSQSClient()

	m.On("SendMessageBatch", mock.Anything, testify.MatchedBy(func(req *sqs.SendMessageBatchInput) bool { return len(req.Entries) == 5 }), mock.Anything).Return(getSendMessageBatchOutput(), fmt.Errorf("Error for unit tests"))

	b := NewBroker("https://sqs.eu-west-3.amazonaws.com/959789434/testing", m)

	list := make([]budget.Transaction, 5)

	err := b.Send(list)

	assert.Equal(t, err.Error(), "Error for unit tests")
	m.AssertExpectations(t)
}

func getSendMessageBatchOutput() *sqs.SendMessageBatchOutput {
	messageID := "5aab4335-527a-41d5-bba6-7e5cfdb8228d"
	return &sqs.SendMessageBatchOutput{
		Successful: []types.SendMessageBatchResultEntry{{MessageId: &messageID}},
		Failed:     []types.BatchResultErrorEntry{{Id: &messageID, Message: &messageID}},
	}
}
