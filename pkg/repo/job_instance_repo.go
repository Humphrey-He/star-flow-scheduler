package repo

import (
	"context"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobdefinition"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobinstance"
)

type JobInstanceFilter struct {
	JobCode  string
	Status   string
	StartAt  *time.Time
	EndAt    *time.Time
	Page     int
	PageSize int
}

type JobInstanceRepository struct {
	client *ent.Client
}

func NewJobInstanceRepository(client *ent.Client) *JobInstanceRepository {
	return &JobInstanceRepository{client: client}
}

func (r *JobInstanceRepository) List(ctx context.Context, filter JobInstanceFilter) ([]*ent.JobInstance, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 200 {
		filter.PageSize = 20
	}

	query := r.client.JobInstance.Query()

	if filter.JobCode != "" {
		job, err := r.client.JobDefinition.Query().Where(jobdefinition.JobCodeEQ(filter.JobCode)).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return []*ent.JobInstance{}, 0, nil
			}
			return nil, 0, err
		}
		query = query.Where(jobinstance.JobIDEQ(job.ID))
	}

	if filter.Status != "" {
		query = query.Where(jobinstance.StatusEQ(filter.Status))
	}

	if filter.StartAt != nil {
		query = query.Where(jobinstance.TriggerTimeGTE(*filter.StartAt))
	}

	if filter.EndAt != nil {
		query = query.Where(jobinstance.TriggerTimeLTE(*filter.EndAt))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	items, err := query.Order(ent.Desc(jobinstance.FieldID)).Limit(filter.PageSize).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *JobInstanceRepository) GetByInstanceNo(ctx context.Context, instanceNo string) (*ent.JobInstance, error) {
	instanceNo = strings.TrimSpace(instanceNo)
	return r.client.JobInstance.Query().Where(jobinstance.InstanceNoEQ(instanceNo)).Only(ctx)
}
