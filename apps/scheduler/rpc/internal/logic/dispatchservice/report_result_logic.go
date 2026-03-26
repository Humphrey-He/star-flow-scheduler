package dispatchservicelogic

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/metricsx"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
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
		l.Logger.Errorw("scheduler report result failed",
			append(reportFields(in, status, l.resolveWorkflowInstanceNo(in.InstanceNo)), logx.Field("error_message", err.Error()))...,
		)
		return nil, err
	}
	if status == state.StatusSuccess {
		metricsx.Inc("scheduler_report_result_success_total")
	} else {
		metricsx.Inc("scheduler_report_result_fail_total")
	}
	l.Logger.Infow("scheduler report result",
		reportFields(in, status, l.resolveWorkflowInstanceNo(in.InstanceNo))...,
	)
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

func reportFields(in *schedulev1.ReportResultRequest, status state.InstanceStatus, workflowInstanceNo string) []logx.LogField {
	fields := []logx.LogField{
		logx.Field("instance_no", in.InstanceNo),
		logx.Field("status", status),
	}
	if in.Meta != nil {
		if in.Meta.TraceId != "" {
			fields = append(fields, logx.Field("trace_id", in.Meta.TraceId))
		}
		if in.Meta.RequestId != "" {
			fields = append(fields, logx.Field("request_id", in.Meta.RequestId))
		}
	}
	if workflowInstanceNo != "" {
		fields = append(fields, logx.Field("workflow_instance_no", workflowInstanceNo))
	}
	if in.ShardNo != "" {
		fields = append(fields, logx.Field("shard_no", in.ShardNo))
	}
	return fields
}

func (l *ReportResultLogic) resolveWorkflowInstanceNo(instanceNo string) string {
	if l.svcCtx == nil || l.svcCtx.Ent == nil || l.svcCtx.InstanceRepo == nil {
		return ""
	}
	instance, err := l.svcCtx.InstanceRepo.GetByInstanceNo(l.ctx, instanceNo)
	if err != nil {
		return ""
	}
	nodeRepo := pkgrepo.NewWorkflowNodeInstanceRepository(l.svcCtx.Ent)
	nodeInst, err := nodeRepo.GetByJobInstanceID(l.ctx, int64(instance.ID))
	if err != nil {
		if ent.IsNotFound(err) {
			return ""
		}
		return ""
	}
	workflowRepo := pkgrepo.NewWorkflowInstanceRepository(l.svcCtx.Ent)
	workflowInst, err := workflowRepo.GetByID(l.ctx, nodeInst.WorkflowInstanceID)
	if err != nil {
		return ""
	}
	return workflowInst.WorkflowInstanceNo
}
