//go:build entgen
// +build entgen

package repo

import (
    "context"

    "github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
    "github.com/Humphrey-He/star-flow-scheduler/pkg/ent/workflownode"
)

type WorkflowNodeCreate struct {
    WorkflowID       int64
    NodeCode         string
    NodeName         string
    JobCode          string
    UpstreamCodes    *string
    TriggerCondition string
    FailStrategy     string
    TimeoutMs        int
    SortOrder        int
}

type WorkflowNodeRepository struct {
    client *ent.Client
}

func NewWorkflowNodeRepository(client *ent.Client) *WorkflowNodeRepository {
    return &WorkflowNodeRepository{client: client}
}

func (r *WorkflowNodeRepository) CreateBatch(ctx context.Context, nodes []WorkflowNodeCreate) ([]*ent.WorkflowNode, error) {
    if len(nodes) == 0 {
        return []*ent.WorkflowNode{}, nil
    }
    bulk := make([]*ent.WorkflowNodeCreate, 0, len(nodes))
    for _, n := range nodes {
        create := r.client.WorkflowNode.Create().
            SetWorkflowID(n.WorkflowID).
            SetNodeCode(n.NodeCode).
            SetNodeName(n.NodeName).
            SetJobCode(n.JobCode).
            SetTriggerCondition(n.TriggerCondition).
            SetFailStrategy(n.FailStrategy).
            SetTimeoutMs(n.TimeoutMs).
            SetSortOrder(n.SortOrder)
        if n.UpstreamCodes != nil {
            create.SetUpstreamCodes(*n.UpstreamCodes)
        }
        bulk = append(bulk, create)
    }
    return r.client.WorkflowNode.CreateBulk(bulk...).Save(ctx)
}

func (r *WorkflowNodeRepository) ListByWorkflowID(ctx context.Context, workflowID int64) ([]*ent.WorkflowNode, error) {
    return r.client.WorkflowNode.Query().
        Where(workflownode.WorkflowIDEQ(workflowID)).
        Order(ent.Asc(workflownode.FieldSortOrder)).
        All(ctx)
}

func (r *WorkflowNodeRepository) DeleteByWorkflowID(ctx context.Context, workflowID int64) (int, error) {
    return r.client.WorkflowNode.Delete().
        Where(workflownode.WorkflowIDEQ(workflowID)).
        Exec(ctx)
}

func (r *WorkflowNodeRepository) GetByNodeCode(ctx context.Context, workflowID int64, nodeCode string) (*ent.WorkflowNode, error) {
    return r.client.WorkflowNode.Query().
        Where(workflownode.WorkflowIDEQ(workflowID), workflownode.NodeCodeEQ(nodeCode)).
        Only(ctx)
}
