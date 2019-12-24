package mock

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/mock"
)

func NewSQSClient() *MockSQS {
	return &MockSQS{}
}

type MockSQS struct {
	sqsiface.SQSAPI
	mock.Mock
}

func (_m *MockSQS) SendMessage(_a0 *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
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
