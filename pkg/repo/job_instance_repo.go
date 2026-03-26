package repo

import (
	"context"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobdefinition"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobinstance"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
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

type JobInstanceCreate struct {
	InstanceNo    string
	JobID         int64
	WorkflowID    *int64
	TriggerType   string
	TriggerTime   time.Time
	ScheduledTime time.Time
	Status        string
	Payload       *string
	ShardTotal    int
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
		query = query.Where(jobinstance.JobIDEQ(int64(job.ID)))
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

func (r *JobInstanceRepository) Create(ctx context.Context, req JobInstanceCreate) (*ent.JobInstance, error) {
	create := r.client.JobInstance.Create().
		SetInstanceNo(req.InstanceNo).
		SetJobID(req.JobID).
		SetTriggerType(req.TriggerType).
		SetTriggerTime(req.TriggerTime).
		SetScheduledTime(req.ScheduledTime).
		SetStatus(req.Status).
		SetShardTotal(req.ShardTotal)

	if req.WorkflowID != nil {
		create.SetWorkflowID(*req.WorkflowID)
	}
	if req.Payload != nil {
		create.SetPayload(*req.Payload)
	}

	return create.Save(ctx)
}

func (r *JobInstanceRepository) UpdateStatusIf(ctx context.Context, instanceNo string, fromStatus types.InstanceStatus, toStatus types.InstanceStatus) (int, error) {
	return r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(string(fromStatus))).
		SetStatus(string(toStatus)).
		Save(ctx)
}

func (r *JobInstanceRepository) UpdateResultIfStatus(ctx context.Context, instanceNo string, fromStatus types.InstanceStatus, toStatus types.InstanceStatus, startTime *time.Time, finishTime *time.Time, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	upd := r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(string(fromStatus))).
		SetStatus(string(toStatus))

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

func (r *JobInstanceRepository) MarkDispatchedIfPending(ctx context.Context, instanceNo string, executorID int64, dispatchTime time.Time) (int, error) {
	return r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(string(types.InstanceStatusPending))).
		SetStatus(string(types.InstanceStatusDispatched)).
		SetExecutorID(executorID).
		SetDispatchTime(dispatchTime).
		Save(ctx)
}

func (r *JobInstanceRepository) MarkRunningIfDispatched(ctx context.Context, instanceNo string, startTime time.Time) (int, error) {
	return r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(string(types.InstanceStatusDispatched))).
		SetStatus(string(types.InstanceStatusRunning)).
		SetStartTime(startTime).
		Save(ctx)
}

func (r *JobInstanceRepository) MarkSuccessIfRunning(ctx context.Context, instanceNo string, finishTime time.Time, resultSummary *string) (int, error) {
	upd := r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(string(types.InstanceStatusRunning))).
		SetStatus(string(types.InstanceStatusSuccess)).
		SetFinishTime(finishTime)

	if resultSummary != nil {
		upd.SetResultSummary(*resultSummary)
	}

	return upd.Save(ctx)
}

func (r *JobInstanceRepository) MarkFailedIfRunning(ctx context.Context, instanceNo string, finishTime time.Time, errorCode *string, errorMessage *string) (int, error) {
	upd := r.client.JobInstance.Update().
		Where(jobinstance.InstanceNoEQ(instanceNo), jobinstance.StatusEQ(string(types.InstanceStatusRunning))).
		SetStatus(string(types.InstanceStatusFailed)).
		SetFinishTime(finishTime)

	if errorCode != nil {
		upd.SetErrorCode(*errorCode)
	}
	if errorMessage != nil {
		upd.SetErrorMessage(*errorMessage)
	}

	return upd.Save(ctx)
}
