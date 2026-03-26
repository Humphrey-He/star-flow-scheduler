package dispatch

import (
	"context"
	"fmt"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/receiver"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"
)

type Service struct {
	receiver *receiver.Receiver
}

func NewService(receiver *receiver.Receiver) *Service {
	return &Service{receiver: receiver}
}

func (s *Service) DispatchJob(ctx context.Context, in *schedulerv1_schedulev1.DispatchJobRequest) error {
	if in == nil {
		return fmt.Errorf("request is nil")
	}
	if in.InstanceNo == "" || in.HandlerName == "" {
		return fmt.Errorf("instance_no or handler_name is empty")
	}

	payload := []byte(nil)
	if in.Payload != nil {
		payload = in.Payload.Raw
	}

	task := &model.Task{
		InstanceNo:    in.InstanceNo,
		ShardNo:       in.ShardNo,
		JobCode:       in.JobCode,
		HandlerName:   in.HandlerName,
		Payload:       payload,
		TimeoutMs:     in.TimeoutMs,
		TraceID:       in.TraceId,
		IdempotentKey: in.IdempotentKey,
		ShardIndex:    in.ShardIndex,
		ShardTotal:    in.ShardTotal,
	}

	return s.receiver.Accept(ctx, task)
}
