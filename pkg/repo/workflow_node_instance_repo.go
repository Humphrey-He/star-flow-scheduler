package repo

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/workflownodeinstance"
)

type WorkflowNodeInstanceCreate struct {
	WorkflowInstanceID int64
	WorkflowID         int64
	NodeCode           string
	JobID              int64
	JobInstanceID      *int64
	Status             string
	ErrorMessage       *string
	StartTime          *time.Time
	FinishTime         *time.Time
}

type WorkflowNodeInstanceRepository struct {
	client *ent.Client
}

func NewWorkflowNodeInstanceRepository(client *ent.Client) *WorkflowNodeInstanceRepository {
	return &WorkflowNodeInstanceRepository{client: client}
}

func (r *WorkflowNodeInstanceRepository) BatchCreate(ctx context.Context, items []WorkflowNodeInstanceCreate) ([]*ent.WorkflowNodeInstance, error) {
	if len(items) == 0 {
		return []*ent.WorkflowNodeInstance{}, nil
	}
	bulk := make([]*ent.WorkflowNodeInstanceCreate, 0, len(items))
	for _, item := range items {
		create := r.client.WorkflowNodeInstance.Create().
			SetWorkflowInstanceID(item.WorkflowInstanceID).
			SetWorkflowID(item.WorkflowID).
			SetNodeCode(item.NodeCode).
			SetJobID(item.JobID).
			SetStatus(item.Status)
		if item.JobInstanceID != nil {
			create.SetJobInstanceID(*item.JobInstanceID)
		}
		if item.ErrorMessage != nil {
			create.SetErrorMessage(*item.ErrorMessage)
		}
		if item.StartTime != nil {
			create.SetStartTime(*item.StartTime)
		}
		if item.FinishTime != nil {
			create.SetFinishTime(*item.FinishTime)
		}
		bulk = append(bulk, create)
	}
	return r.client.WorkflowNodeInstance.CreateBulk(bulk...).Save(ctx)
}

func (r *WorkflowNodeInstanceRepository) ListByWorkflowInstanceID(ctx context.Context, workflowInstanceID int64) ([]*ent.WorkflowNodeInstance, error) {
	return r.client.WorkflowNodeInstance.Query().
		Where(workflownodeinstance.WorkflowInstanceIDEQ(workflowInstanceID)).
		Order(ent.Asc(workflownodeinstance.FieldID)).
		All(ctx)
}

func (r *WorkflowNodeInstanceRepository) GetByWorkflowInstanceIDAndNodeCode(ctx context.Context, workflowInstanceID int64, nodeCode string) (*ent.WorkflowNodeInstance, error) {
	return r.client.WorkflowNodeInstance.Query().
		Where(
			workflownodeinstance.WorkflowInstanceIDEQ(workflowInstanceID),
			workflownodeinstance.NodeCodeEQ(nodeCode),
		).
		Only(ctx)
}

func (r *WorkflowNodeInstanceRepository) GetByJobInstanceID(ctx context.Context, jobInstanceID int64) (*ent.WorkflowNodeInstance, error) {
	return r.client.WorkflowNodeInstance.Query().
		Where(workflownodeinstance.JobInstanceIDEQ(jobInstanceID)).
		Only(ctx)
}

func (r *WorkflowNodeInstanceRepository) UpdateStatusIf(ctx context.Context, workflowInstanceID int64, nodeCode string, fromStatus string, toStatus string, startTime *time.Time) (int, error) {
	update := r.client.WorkflowNodeInstance.Update().
		Where(
			workflownodeinstance.WorkflowInstanceIDEQ(workflowInstanceID),
			workflownodeinstance.NodeCodeEQ(nodeCode),
			workflownodeinstance.StatusEQ(fromStatus),
		).
		SetStatus(toStatus)
	if startTime != nil {
		update.SetStartTime(*startTime)
	}
	return update.Save(ctx)
}

func (r *WorkflowNodeInstanceRepository) UpdateJobInstanceID(ctx context.Context, workflowInstanceID int64, nodeCode string, jobInstanceID int64) (int, error) {
	return r.client.WorkflowNodeInstance.Update().
		Where(
			workflownodeinstance.WorkflowInstanceIDEQ(workflowInstanceID),
			workflownodeinstance.NodeCodeEQ(nodeCode),
		).
		SetJobInstanceID(jobInstanceID).
		Save(ctx)
}

func (r *WorkflowNodeInstanceRepository) ListByStatus(ctx context.Context, workflowInstanceID int64, status string) ([]*ent.WorkflowNodeInstance, error) {
	return r.client.WorkflowNodeInstance.Query().
		Where(
			workflownodeinstance.WorkflowInstanceIDEQ(workflowInstanceID),
			workflownodeinstance.StatusEQ(status),
		).
		Order(ent.Asc(workflownodeinstance.FieldID)).
		All(ctx)
}
