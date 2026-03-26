package repo

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobinstance"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type JobInstanceRepository struct {
	client *ent.Client
	inner  *repo.JobInstanceRepository
}

func NewJobInstanceRepository(client *ent.Client) *JobInstanceRepository {
	return &JobInstanceRepository{
		client: client,
		inner:  repo.NewJobInstanceRepository(client),
	}
}

func (r *JobInstanceRepository) GetByInstanceNo(ctx context.Context, instanceNo string) (*ent.JobInstance, error) {
	return r.client.JobInstance.Query().Where(jobinstance.InstanceNoEQ(instanceNo)).Only(ctx)
}

func (r *JobInstanceRepository) Create(ctx context.Context, req repo.JobInstanceCreate) (*ent.JobInstance, error) {
	return r.inner.Create(ctx, req)
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

func (r *JobInstanceRepository) UpdateResultIfStatus(ctx context.Context, instanceNo string, fromStatus string, toStatus string, startTime *time.Time, finishTime *time.Time, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	upd := r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(fromStatus)).
		SetStatus(toStatus)

	if startTime != nil {
		upd.SetStartTime(*startTime)
	}
	if finishTime != nil {
		upd.SetFinishTime(*finishTime)
	}
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

func (r *JobInstanceRepository) UpdateDispatchInfoIfStatus(ctx context.Context, instanceNo string, fromStatus string, toStatus string, executorID int64, dispatchTime time.Time) (int, error) {
	return r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(fromStatus)).
		SetStatus(toStatus).
		SetExecutorID(executorID).
		SetDispatchTime(dispatchTime).
		Save(ctx)
}

func (r *JobInstanceRepository) ListDueInstances(ctx context.Context, now time.Time, limit int) ([]string, error) {
	if limit <= 0 {
		return []string{}, nil
	}
	rows, err := r.client.JobInstance.Query().
		Where(
			jobinstance.StatusEQ("pending"),
			jobinstance.ScheduledTimeLTE(now),
		).
		Order(jobinstance.ByScheduledTime()).
		Limit(limit).
		Select(jobinstance.FieldInstanceNo).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.InstanceNo)
	}
	return out, nil
}
