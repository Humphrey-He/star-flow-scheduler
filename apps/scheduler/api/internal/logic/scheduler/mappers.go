package scheduler

import (
	"encoding/json"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/internal/models"
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

func mapJobDefinition(job *models.JobDefinition) types.JobDefinition {
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

func mapJobInstance(instance *models.JobInstance) types.JobInstance {
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
		CurrentBackoff: instance.CurrentBackoff,
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

func marshalPayload(payload map[string]interface{}) (*string, error) {
	if payload == nil {
		return nil, nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	s := string(data)
	return &s, nil
}

func parseTime(raw string) (*time.Time, error) {
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
