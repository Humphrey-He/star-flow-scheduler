package executorregistryservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartbeatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatLogic {
	return &HeartbeatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HeartbeatLogic) Heartbeat(in *schedulerv1_schedulev1.HeartbeatRequest) (*schedulerv1_schedulev1.HeartbeatResponse, error) {
	err := l.svcCtx.RegistrySvc.Heartbeat(l.ctx, in.ExecutorCode, int(in.CurrentLoad))
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.HeartbeatResponse{
		Status:        schedulerv1_schedulev1.ExecutorStatus_EXECUTOR_STATUS_ONLINE,
		AcceptNewJobs: true,
	}, nil
}
