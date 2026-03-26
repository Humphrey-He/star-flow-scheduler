package dispatchservicelogic

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
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
	metricsx.Inc("scheduler_report_result_total")
	status := mapReportStatus(in.Status)
	startAt := unixMilliPtr(in.StartTime)
	finishAt := unixMilliPtr(in.FinishTime)
	if startAt != nil && finishAt != nil && finishAt.After(*startAt) {
		metricsx.ObserveDurationMs("scheduler_report_result_duration_ms", finishAt.Sub(*startAt))
	}
	_, err := l.svcCtx.InstanceSvc.ReportResult(l.ctx, in.InstanceNo, status, startAt, finishAt, strPtr(in.ResultSummary), strPtr(in.ErrorCode), strPtr(in.ErrorMessage))
	if err != nil {
		metricsx.Inc("scheduler_report_result_fail_total")
		return nil, err
	}
	if status == state.StatusSuccess {
		metricsx.Inc("scheduler_report_result_success_total")
	} else {
		metricsx.Inc("scheduler_report_result_fail_total")
	}
	if l.svcCtx.WorkflowRuntime != nil {
		instance, err := l.svcCtx.InstanceRepo.GetByInstanceNo(l.ctx, in.InstanceNo)
		if err == nil && instance.WorkflowID != nil {
			nodeStatus := mapWorkflowNodeStatus(status)
			_ = l.svcCtx.WorkflowRuntime.OnJobInstanceFinished(l.ctx, int64(instance.ID), nodeStatus)
		}
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

func mapWorkflowNodeStatus(status state.InstanceStatus) types.WorkflowNodeStatus {
	switch status {
	case state.StatusSuccess:
		return types.WorkflowNodeStatusSuccess
	case state.StatusFailed:
		return types.WorkflowNodeStatusFailed
	default:
		return types.WorkflowNodeStatusFailed
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
