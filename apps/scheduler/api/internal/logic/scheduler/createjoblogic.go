package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateJobLogic {
	return &CreateJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateJobLogic) CreateJob(req *types.CreateJobRequest) (resp *types.CreateJobResponse, err error) {
	payload, err := marshalPayload(req.Payload)
	if err != nil {
		return nil, err
	}

	job := repo.JobDefinitionCreate{
		JobCode:           req.JobCode,
		JobName:           req.JobName,
		JobType:           req.JobType,
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
	}

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
	if job.JobType == "cron" && (job.ScheduleExpr == nil || *job.ScheduleExpr == "") {
		return nil, fmt.Errorf("schedule_expr is required for cron jobs")
	}
	if job.JobType == "delay" && (job.DelayMs == nil || *job.DelayMs <= 0) {
		return nil, fmt.Errorf("delay_ms must be > 0 for delay jobs")
	}
	if job.JobType == "once" && job.ScheduleExpr != nil && *job.ScheduleExpr != "" {
		return nil, fmt.Errorf("schedule_expr must be empty for once jobs")
	}
	if job.JobType == "dag" && job.HandlerName != "" {
		return nil, fmt.Errorf("handler_name should be empty for dag jobs")
	}
	if job.JobType != "dag" && job.HandlerName == "" {
		return nil, fmt.Errorf("handler_name is required")
	}

	ctx, cancel := context.WithTimeout(l.ctx, 5*time.Second)
	defer cancel()

	exists, err := l.svcCtx.Jobs.ExistsByCode(ctx, job.JobCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, repo.ErrAlreadyExists
	}

	created, err := l.svcCtx.Jobs.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	return &types.CreateJobResponse{
		JobId:   created.ID,
		JobCode: created.JobCode,
	}, nil
}
