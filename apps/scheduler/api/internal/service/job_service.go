package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/validator"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type JobService struct {
	jobs *repo.JobRepository
}

func NewJobService(jobs *repo.JobRepository) *JobService {
	return &JobService{jobs: jobs}
}

func (s *JobService) Create(ctx context.Context, req *types.CreateJobRequest) (*ent.JobDefinition, error) {
	if err := validator.ValidateCreateJob(req); err != nil {
		return nil, err
	}

	payload, err := marshalPayload(req.Payload)
	if err != nil {
		return nil, errx.InvalidParam("invalid payload")
	}

	job := repo.JobDefinitionCreate{
		JobCode:           strings.TrimSpace(req.JobCode),
		JobName:           strings.TrimSpace(req.JobName),
		JobType:           strings.TrimSpace(req.JobType),
		ScheduleExpr:      req.ScheduleExpr,
		DelayMs:           req.DelayMs,
		ExecuteMode:       req.ExecuteMode,
		HandlerName:       req.HandlerName,
		HandlerPayload:    payload,
		TimeoutMs:         req.TimeoutMs,
		RetryLimit:        req.RetryLimit,
		RetryBackoff:      req.RetryBackoff,
		Priority:          req.Priority,
		ShardTotal:        req.ShardTotal,
		RouteStrategy:     req.RouteStrategy,
		ExecutorTag:       req.ExecutorTag,
		IdempotentKeyExpr: req.IdempotentKeyExpr,
		Status:            req.Status,
		CreatedBy:         req.CreatedBy,
		UpdatedBy:         req.UpdatedBy,
		CreatedAt:         ptrTime(time.Now()),
	}

	applyDefaults(&job)

	exists, err := s.jobs.ExistsByCode(ctx, job.JobCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errx.Conflict("job_code already exists")
	}

	return s.jobs.Create(ctx, job)
}

func (s *JobService) GetByCode(ctx context.Context, jobCode string) (*ent.JobDefinition, error) {
	return s.jobs.GetByCode(ctx, strings.TrimSpace(jobCode))
}

func applyDefaults(job *repo.JobDefinitionCreate) {
	if job.ExecuteMode == "" {
		job.ExecuteMode = "standalone"
	}
	if job.TimeoutMs == 0 {
		job.TimeoutMs = 60000
	}
	if job.RetryLimit == 0 {
		job.RetryLimit = 3
	}
	if job.RetryBackoff == "" {
		job.RetryBackoff = "1s,3s,5s"
	}
	if job.Priority == 0 {
		job.Priority = 5
	}
	if job.ShardTotal == 0 {
		job.ShardTotal = 1
	}
	if job.RouteStrategy == "" {
		job.RouteStrategy = "least_load"
	}
	if job.Status == "" {
		job.Status = "enabled"
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

func ptrTime(t time.Time) *time.Time {
	return &t
}
