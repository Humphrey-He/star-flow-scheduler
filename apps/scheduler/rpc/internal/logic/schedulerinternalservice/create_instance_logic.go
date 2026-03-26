package schedulerinternalservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateInstanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInstanceLogic {
	return &CreateInstanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateInstanceLogic) CreateInstance(in *schedulerv1_schedulev1.CreateInstanceRequest) (*schedulerv1_schedulev1.CreateInstanceResponse, error) {
	if in == nil {
		return nil, nil
	}

	payload := payloadString(in.Payload)
	trigger := in.TriggerType
	if trigger == "" {
		trigger = "manual"
	}

	instance, err := l.svcCtx.DispatchSvc.CreateInstance(l.ctx, in.JobCode, trigger, payload)
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.CreateInstanceResponse{
		InstanceNo: instance.InstanceNo,
		Status:     mapStatus(state.StatusPending),
	}, nil
}

func payloadString(payload *schedulerv1_schedulev1.JobPayload) *string {
	if payload == nil || len(payload.Raw) == 0 {
		return nil
	}
	s := string(payload.Raw)
	return &s
}

func mapStatus(status state.InstanceStatus) schedulerv1_schedulev1.InstanceStatus {
	switch status {
	case state.StatusPending:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_PENDING
	case state.StatusDispatched:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_DISPATCHED
	case state.StatusRunning:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_RUNNING
	case state.StatusSuccess:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS
	case state.StatusFailed:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED
	case state.StatusRetryWait:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_RETRY_WAIT
	case state.StatusDead:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_DEAD
	case state.StatusCanceled:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_CANCELED
	default:
		return schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_UNSPECIFIED
	}
}
