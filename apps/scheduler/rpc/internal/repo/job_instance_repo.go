package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobinstance"
)

type JobInstanceRepository struct {
	client *ent.Client
}

func NewJobInstanceRepository(client *ent.Client) *JobInstanceRepository {
	return &JobInstanceRepository{client: client}
}

func (r *JobInstanceRepository) UpdateStatusIf(ctx context.Context, instanceNo string, fromStatus string, toStatus string) (int, error) {
	return r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(fromStatus)).
		SetStatus(toStatus).
		Save(ctx)
}

func (r *JobInstanceRepository) UpdateResult(ctx context.Context, instanceNo string, status string, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	upd := r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo)).
		SetStatus(status)

	if resultSummary != nil {
		upd.SetResultSummary(*resultSummary)
	}
	if errorCode != nil {
		upd.SetErrorCode(*errorCode)
	}
	if errorMessage != nil {
		upd.SetErrorMessage(*errorMessage)
	}

	return upd.Save(ctx)
}
