package dispatchservicelogic

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
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
	status := mapReportStatus(in.Status)
	_, err := l.svcCtx.InstanceSvc.ReportResult(l.ctx, in.InstanceNo, status, strPtr(in.ResultSummary), strPtr(in.ErrorCode), strPtr(in.ErrorMessage))
	if err != nil {
		return nil, err
	}

	return &schedulerv1_schedulev1.ReportResultResponse{Ok: true}, nil
}

func mapReportStatus(status schedulerv1_schedulev1.InstanceStatus) state.InstanceStatus {
	switch status {
	case schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS:
		return state.StatusSuccess
	case schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED:
		return state.StatusFailed
	case schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_RUNNING:
		return state.StatusRunning
	case schedulerv1_schedulev1.InstanceStatus_INSTANCE_STATUS_DISPATCHED:
		return state.StatusDispatched
	default:
		return state.StatusFailed
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
