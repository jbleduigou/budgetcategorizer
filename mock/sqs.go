package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jbleduigou/budgetcategorizer/iface"
	"github.com/stretchr/testify/mock"
)

// NewSQSClient provides a mock instance of a SQS client
func NewSQSClient() *SQS {
	return &SQS{}
}

// SQS is an implementation of the SQSAPI interface with a mock, use for testing not for production
type SQS struct {
	iface.SQSSendMessageAPI
	mock.Mock
}

// SendMessageBatch API operation for Amazon Simple Queue Service.
// Delivers a bunch of messages to the specified queue.
func (_m *SQS) SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error) {
	ret := _m.Called(ctx, params, optFns)
	if ret.Get(0) == nil && ret.Get(1) == nil {
		return nil, nil
	}
	if ret.Get(0) == nil {
		return nil, ret.Get(1).(error)
	}
	if ret.Get(1) == nil {
		return ret.Get(0).(*sqs.SendMessageBatchOutput), nil
	}
	return ret.Get(0).(*sqs.SendMessageBatchOutput), ret.Get(1).(error)
}
