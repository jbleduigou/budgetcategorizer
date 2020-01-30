package mock

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/mock"
)

// NewSQSClient provides a mock instance of a SQS client
func NewSQSClient() *SQS {
	return &SQS{}
}

// SQS is an implementation of the SQSAPI interface with a mock, use for testing not for production
type SQS struct {
	sqsiface.SQSAPI
	mock.Mock
}

// SendMessage API operation for Amazon Simple Queue Service.
// Delivers a message to the specified queue.
func (_m *SQS) SendMessage(_a0 *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	ret := _m.Called(_a0)
	if ret.Get(0) == nil && ret.Get(1) == nil {
		return nil, nil
	}
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(error)
	}
	if ret.Get(1) == nil {
		return ret.Get(0).(*sqs.SendMessageOutput), nil
	}
	return ret.Get(0).(*sqs.SendMessageOutput), ret.Get(1).(error)
}
