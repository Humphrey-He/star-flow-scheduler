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

func (r *JobInstanceRepository) GetByInstanceNo(ctx context.Context, instanceNo string) (*ent.JobInstance, error) {
	return r.client.JobInstance.Query().Where(jobinstance.InstanceNoEQ(instanceNo)).Only(ctx)
}

func (r *JobInstanceRepository) GetStatusByInstanceNo(ctx context.Context, instanceNo string) (string, error) {
	row, err := r.client.JobInstance.Query().
		Where(jobinstance.InstanceNoEQ(instanceNo)).
		Select(jobinstance.FieldStatus).
		Only(ctx)
	if err != nil {
		return "", err
	}
	return row.Status, nil
}

func (r *JobInstanceRepository) UpdateStatusIf(ctx context.Context, instanceNo string, fromStatus string, toStatus string) (int, error) {
	return r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(fromStatus)).
		SetStatus(toStatus).
		Save(ctx)
}

func (r *JobInstanceRepository) UpdateResultIfStatus(ctx context.Context, instanceNo string, fromStatus string, toStatus string, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	upd := r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(fromStatus)).
		SetStatus(toStatus)

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
