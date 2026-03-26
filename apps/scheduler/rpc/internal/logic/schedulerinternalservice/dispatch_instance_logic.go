package schedulerinternalservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchInstanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchInstanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchInstanceLogic {
	return &DispatchInstanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DispatchInstanceLogic) DispatchInstance(in *schedulerv1_schedulev1.DispatchInstanceRequest) (*schedulerv1_schedulev1.DispatchInstanceResponse, error) {
	if in == nil {
		return nil, nil
	}

	exec, err := l.svcCtx.DispatchSvc.DispatchInstance(l.ctx, in.InstanceNo)
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.DispatchInstanceResponse{
		Dispatched:   true,
		ExecutorCode: exec.ExecutorCode,
	}, nil
}
