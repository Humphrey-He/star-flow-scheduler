package dispatchservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/svc"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchJobLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchJobLogic {
	return &DispatchJobLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DispatchJobLogic) DispatchJob(in *schedulerv1_schedulev1.DispatchJobRequest) (*schedulerv1_schedulev1.DispatchJobResponse, error) {
	l.Logger.Infof("dispatch job instance=%s shard=%s handler=%s", in.InstanceNo, in.ShardNo, in.HandlerName)
	return &schedulerv1_schedulev1.DispatchJobResponse{Accepted: true, Message: "accepted"}, nil
}
