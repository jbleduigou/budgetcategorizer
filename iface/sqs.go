package iface

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSSendMessageAPI interface {
	SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error)
}
