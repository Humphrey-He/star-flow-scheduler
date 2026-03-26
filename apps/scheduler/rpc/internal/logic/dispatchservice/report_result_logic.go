package dispatchservicelogic

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

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

func (l *ReportResultLogic) ReportResult(in *schedulev1.ReportResultRequest) (*schedulev1.ReportResultResponse, error) {
	status := mapReportStatus(in.Status)
	startAt := unixMilliPtr(in.StartTime)
	finishAt := unixMilliPtr(in.FinishTime)
	_, err := l.svcCtx.InstanceSvc.ReportResult(l.ctx, in.InstanceNo, status, startAt, finishAt, strPtr(in.ResultSummary), strPtr(in.ErrorCode), strPtr(in.ErrorMessage))
	if err != nil {
		return nil, err
	}

	return &schedulev1.ReportResultResponse{Ok: true}, nil
}

func mapReportStatus(status schedulev1.InstanceStatus) state.InstanceStatus {
	switch status {
	case schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS:
		return state.StatusSuccess
	case schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED:
		return state.StatusFailed
	case schedulev1.InstanceStatus_INSTANCE_STATUS_RUNNING:
		return state.StatusRunning
	case schedulev1.InstanceStatus_INSTANCE_STATUS_DISPATCHED:
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

func unixMilliPtr(ms int64) *time.Time {
	if ms <= 0 {
		return nil
	}
	t := time.UnixMilli(ms)
	return &t
}
