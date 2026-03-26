package assembler

import (
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
)

const timeLayout = time.RFC3339

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func formatTimePtr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.Format(timeLayout)
	return &s
}

func MapJobDefinition(job *ent.JobDefinition) types.JobDefinition {
	return types.JobDefinition{
		Id:                job.ID,
		JobCode:           job.JobCode,
		JobName:           job.JobName,
		JobType:           job.JobType,
		ScheduleExpr:      job.ScheduleExpr,
		DelayMs:           job.DelayMs,
		ExecuteMode:       job.ExecuteMode,
		HandlerName:       job.HandlerName,
		HandlerPayload:    job.HandlerPayload,
		TimeoutMs:         job.TimeoutMs,
		RetryLimit:        job.RetryLimit,
		RetryBackoff:      job.RetryBackoff,
		Priority:          job.Priority,
		ShardTotal:        job.ShardTotal,
		RouteStrategy:     job.RouteStrategy,
		ExecutorTag:       job.ExecutorTag,
		IdempotentKeyExpr: job.IdempotentKeyExpr,
		Status:            job.Status,
		CreatedBy:         job.CreatedBy,
		UpdatedBy:         job.UpdatedBy,
		CreatedAt:         formatTime(job.CreatedAt),
		UpdatedAt:         formatTime(job.UpdatedAt),
	}
}

func MapJobInstance(instance *ent.JobInstance) types.JobInstance {
	return types.JobInstance{
		Id:             instance.ID,
		InstanceNo:     instance.InstanceNo,
		JobId:          instance.JobID,
		WorkflowId:     instance.WorkflowID,
		TriggerType:    instance.TriggerType,
		TriggerTime:    formatTime(instance.TriggerTime),
		ScheduledTime:  formatTime(instance.ScheduledTime),
		DispatchTime:   formatTimePtr(instance.DispatchTime),
		StartTime:      formatTimePtr(instance.StartTime),
		FinishTime:     formatTimePtr(instance.FinishTime),
		Status:         instance.Status,
		RetryCount:     instance.RetryCount,
		CurrentBackoff: instance.CurrentBackoffMs,
		ExecutorId:     instance.ExecutorID,
		ShardTotal:     instance.ShardTotal,
		SuccessShards:  instance.SuccessShards,
		FailedShards:   instance.FailedShards,
		Payload:        instance.Payload,
		ResultSummary:  instance.ResultSummary,
		ErrorCode:      instance.ErrorCode,
		ErrorMessage:   instance.ErrorMessage,
		CreatedAt:      formatTime(instance.CreatedAt),
		UpdatedAt:      formatTime(instance.UpdatedAt),
	}
}

func ParseTime(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t, nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", raw)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
