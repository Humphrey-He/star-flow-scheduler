//go:build entgen
// +build entgen

package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/workflowdefinition"
)

type WorkflowCreate struct {
	WorkflowCode string
	WorkflowName string
	Description  *string
	Status       string
}

type WorkflowRepository struct {
	client *ent.Client
}

func NewWorkflowRepository(client *ent.Client) *WorkflowRepository {
	return &WorkflowRepository{client: client}
}

func (r *WorkflowRepository) Create(ctx context.Context, req WorkflowCreate) (*ent.WorkflowDefinition, error) {
	create := r.client.WorkflowDefinition.Create().
		SetWorkflowCode(req.WorkflowCode).
		SetWorkflowName(req.WorkflowName).
		SetStatus(req.Status)

	if req.Description != nil {
		create.SetDescription(*req.Description)
	}

	return create.Save(ctx)
}

func (r *WorkflowRepository) GetByCode(ctx context.Context, workflowCode string) (*ent.WorkflowDefinition, error) {
	return r.client.WorkflowDefinition.Query().Where(workflowdefinition.WorkflowCodeEQ(workflowCode)).Only(ctx)
}

func (r *WorkflowRepository) List(ctx context.Context, status string, limit int) ([]*ent.WorkflowDefinition, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	query := r.client.WorkflowDefinition.Query()
	if status != "" {
		query = query.Where(workflowdefinition.StatusEQ(status))
	}
	return query.Order(ent.Desc(workflowdefinition.FieldID)).Limit(limit).All(ctx)
}
