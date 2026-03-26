package scheduler

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/internal/models"
	"github.com/Humphrey-He/star-flow-scheduler/internal/repo"

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

	job := &models.JobDefinition{
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

	repo.BuildJobDefinitionDefaults(job)
	if err := repo.ValidateJobDefinition(job); err != nil {
		return nil, err
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

	id, err := l.svcCtx.Jobs.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	return &types.CreateJobResponse{
		JobId:   id,
		JobCode: job.JobCode,
	}, nil
}
