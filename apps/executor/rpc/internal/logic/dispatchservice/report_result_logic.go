package dispatchservicelogic

import (
	"context"
	"errors"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/svc"
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
	return nil, errors.New("report_result is handled by scheduler rpc")
}
