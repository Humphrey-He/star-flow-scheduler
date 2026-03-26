package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobexecutionlog"
)

type ExecutionLogCreate struct {
	InstanceNo string
	ShardNo    *string
	ExecutorID *int64
	LogLevel   string
	Phase      string
	Message    string
	TraceID    *string
}

type ExecutionLogRepository struct {
	client *ent.Client
}

func NewExecutionLogRepository(client *ent.Client) *ExecutionLogRepository {
	return &ExecutionLogRepository{client: client}
}

func (r *ExecutionLogRepository) Create(ctx context.Context, req ExecutionLogCreate) (*ent.JobExecutionLog, error) {
	create := r.client.JobExecutionLog.Create().
		SetInstanceNo(req.InstanceNo).
		SetLogLevel(req.LogLevel).
		SetPhase(req.Phase).
		SetMessage(req.Message)

	if req.ShardNo != nil {
		create.SetShardNo(*req.ShardNo)
	}
	if req.ExecutorID != nil {
		create.SetExecutorID(*req.ExecutorID)
	}
	if req.TraceID != nil {
		create.SetTraceID(*req.TraceID)
	}

	return create.Save(ctx)
}

func (r *ExecutionLogRepository) ListByInstanceNo(ctx context.Context, instanceNo string, limit int) ([]*ent.JobExecutionLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	return r.client.JobExecutionLog.Query().
		Where(jobexecutionlog.InstanceNoEQ(instanceNo)).
		Order(ent.Desc(jobexecutionlog.FieldID)).
		Limit(limit).
		All(ctx)
}

func (r *ExecutionLogRepository) ListByShardNo(ctx context.Context, shardNo string, limit int) ([]*ent.JobExecutionLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	return r.client.JobExecutionLog.Query().
		Where(jobexecutionlog.ShardNoEQ(shardNo)).
		Order(ent.Desc(jobexecutionlog.FieldID)).
		Limit(limit).
		All(ctx)
}
