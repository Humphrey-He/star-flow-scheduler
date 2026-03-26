package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type ExecutorRepository struct {
	inner *repo.ExecutorRepository
}

func NewExecutorRepository(inner *repo.ExecutorRepository) *ExecutorRepository {
	return &ExecutorRepository{inner: inner}
}

func (r *ExecutorRepository) Upsert(ctx context.Context, req repo.ExecutorUpsert) (*ent.Executor, error) {
	return r.inner.Upsert(ctx, req)
}

func (r *ExecutorRepository) UpdateHeartbeat(ctx context.Context, executorCode string, currentLoad int) error {
	return r.inner.UpdateHeartbeat(ctx, executorCode, currentLoad)
}
