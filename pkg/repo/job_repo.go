package repo

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobdefinition"
)

type JobDefinitionCreate struct {
	JobCode           string
	JobName           string
	JobType           string
	ScheduleExpr      *string
	DelayMs           *int64
	ExecuteMode       string
	HandlerName       string
	HandlerPayload    *string
	TimeoutMs         int
	RetryLimit        int
	RetryBackoff      string
	Priority          int
	ShardTotal        int
	RouteStrategy     string
	ExecutorTag       *string
	IdempotentKeyExpr *string
	Status            string
	CreatedBy         *string
	UpdatedBy         *string
	CreatedAt         *time.Time
}

type JobRepository struct {
	client *ent.Client
}

func NewJobRepository(client *ent.Client) *JobRepository {
	return &JobRepository{client: client}
}

func (r *JobRepository) Create(ctx context.Context, req JobDefinitionCreate) (*ent.JobDefinition, error) {
	create := r.client.JobDefinition.Create().
		SetJobCode(req.JobCode).
		SetJobName(req.JobName).
		SetJobType(req.JobType).
		SetExecuteMode(req.ExecuteMode).
		SetHandlerName(req.HandlerName).
		SetTimeoutMs(req.TimeoutMs).
		SetRetryLimit(req.RetryLimit).
		SetRetryBackoff(req.RetryBackoff).
		SetPriority(req.Priority).
		SetShardTotal(req.ShardTotal).
		SetRouteStrategy(req.RouteStrategy).
		SetStatus(req.Status)

	if req.ScheduleExpr != nil {
		create.SetScheduleExpr(*req.ScheduleExpr)
	}
	if req.DelayMs != nil {
		create.SetDelayMs(*req.DelayMs)
	}
	if req.HandlerPayload != nil {
		create.SetHandlerPayload(*req.HandlerPayload)
	}
	if req.ExecutorTag != nil {
		create.SetExecutorTag(*req.ExecutorTag)
	}
	if req.IdempotentKeyExpr != nil {
		create.SetIdempotentKeyExpr(*req.IdempotentKeyExpr)
	}
	if req.CreatedBy != nil {
		create.SetCreatedBy(*req.CreatedBy)
	}
	if req.UpdatedBy != nil {
		create.SetUpdatedBy(*req.UpdatedBy)
	}
	if req.CreatedAt != nil {
		create.SetCreatedAt(*req.CreatedAt)
	}

	return create.Save(ctx)
}

func (r *JobRepository) GetByCode(ctx context.Context, jobCode string) (*ent.JobDefinition, error) {
	return r.client.JobDefinition.Query().Where(jobdefinition.JobCodeEQ(jobCode)).Only(ctx)
}

func (r *JobRepository) GetByID(ctx context.Context, id int64) (*ent.JobDefinition, error) {
	return r.client.JobDefinition.Get(ctx, int(id))
}

func (r *JobRepository) ExistsByCode(ctx context.Context, jobCode string) (bool, error) {
	return r.client.JobDefinition.Query().Where(jobdefinition.JobCodeEQ(jobCode)).Exist(ctx)
}
