package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/workflowinstance"
)

type WorkflowInstanceCreate struct {
	WorkflowInstanceNo string
	WorkflowID         int64
	WorkflowCode       string
	Status             string
}

type WorkflowInstanceRepository struct {
	client *ent.Client
}

func NewWorkflowInstanceRepository(client *ent.Client) *WorkflowInstanceRepository {
	return &WorkflowInstanceRepository{client: client}
}

func (r *WorkflowInstanceRepository) Create(ctx context.Context, req WorkflowInstanceCreate) (*ent.WorkflowInstance, error) {
	return r.client.WorkflowInstance.Create().
		SetWorkflowInstanceNo(req.WorkflowInstanceNo).
		SetWorkflowID(req.WorkflowID).
		SetWorkflowCode(req.WorkflowCode).
		SetStatus(req.Status).
		Save(ctx)
}

func (r *WorkflowInstanceRepository) GetByInstanceNo(ctx context.Context, instanceNo string) (*ent.WorkflowInstance, error) {
	return r.client.WorkflowInstance.Query().Where(workflowinstance.WorkflowInstanceNoEQ(instanceNo)).Only(ctx)
}

func (r *WorkflowInstanceRepository) GetByID(ctx context.Context, workflowInstanceID int64) (*ent.WorkflowInstance, error) {
	return r.client.WorkflowInstance.Get(ctx, workflowInstanceID)
}

func (r *WorkflowInstanceRepository) ListByWorkflowCode(ctx context.Context, workflowCode string, limit int) ([]*ent.WorkflowInstance, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return r.client.WorkflowInstance.Query().
		Where(workflowinstance.WorkflowCodeEQ(workflowCode)).
		Order(ent.Desc(workflowinstance.FieldID)).
		Limit(limit).
		All(ctx)
}

func (r *WorkflowInstanceRepository) UpdateStatusIf(ctx context.Context, workflowInstanceID int64, fromStatus string, toStatus string) (int, error) {
	return r.client.WorkflowInstance.Update().
		Where(
			workflowinstance.IDEQ(workflowInstanceID),
			workflowinstance.StatusEQ(fromStatus),
		).
		SetStatus(toStatus).
		Save(ctx)
}
