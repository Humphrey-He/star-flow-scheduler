package dispatchservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

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

func (l *DispatchJobLogic) DispatchJob(in *schedulev1.DispatchJobRequest) (*schedulev1.DispatchJobResponse, error) {
	_ = l.svcCtx.DispatchSvc
	return &schedulev1.DispatchJobResponse{
		Accepted: false,
		Message:  "dispatch not supported on scheduler rpc",
	}, nil
}
