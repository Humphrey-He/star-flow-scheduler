package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type JobRepository struct {
	inner *repo.JobRepository
}

func NewJobRepository(inner *repo.JobRepository) *JobRepository {
	return &JobRepository{inner: inner}
}

func (r *JobRepository) GetByCode(ctx context.Context, jobCode string) (*ent.JobDefinition, error) {
	return r.inner.GetByCode(ctx, jobCode)
}

func (r *JobRepository) GetByID(ctx context.Context, id int64) (*ent.JobDefinition, error) {
	return r.inner.GetByID(ctx, id)
}
