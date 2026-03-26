package dispatchservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportResultLogic {
	return &ReportResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportResultLogic) ReportResult(in *schedulerv1_schedulev1.ReportResultRequest) (*schedulerv1_schedulev1.ReportResultResponse, error) {
	_ = l.svcCtx.InstanceSvc
	l.Logger.Infof("report result instance=%s shard=%s status=%v", in.InstanceNo, in.ShardNo, in.Status)
	return &schedulerv1_schedulev1.ReportResultResponse{Ok: true}, nil
}
