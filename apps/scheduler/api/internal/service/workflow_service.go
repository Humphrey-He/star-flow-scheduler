package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/db"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type WorkflowService struct {
	ent       *ent.Client
	workflows *repo.WorkflowRepository
	nodes     *repo.WorkflowNodeRepository
	jobs      *repo.JobRepository
}

type WorkflowNodeSpec struct {
	NodeCode         string
	NodeName         string
	JobCode          string
	UpstreamCodes    []string
	TriggerCondition string
	FailStrategy     string
	TimeoutMs        int
	SortOrder        int
}

type CreateWorkflowRequest struct {
	WorkflowCode string
	WorkflowName string
	Description  *string
	Nodes        []WorkflowNodeSpec
}

type UpdateWorkflowRequest struct {
	WorkflowCode string
	WorkflowName string
	Description  *string
	Nodes        []WorkflowNodeSpec
}

func NewWorkflowService(entClient *ent.Client, workflows *repo.WorkflowRepository, nodes *repo.WorkflowNodeRepository, jobs *repo.JobRepository) *WorkflowService {
	return &WorkflowService{
		ent:       entClient,
		workflows: workflows,
		nodes:     nodes,
		jobs:      jobs,
	}
}

func (s *WorkflowService) CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (*ent.WorkflowDefinition, error) {
	if err := s.validateDefinition(ctx, req.WorkflowCode, req.Nodes); err != nil {
		return nil, err
	}

	var created *ent.WorkflowDefinition
	err := db.Transact(ctx, s.ent, func(ctx context.Context, tx *ent.Tx) error {
		txWorkflows := repo.NewWorkflowRepository(tx.Client())
		txNodes := repo.NewWorkflowNodeRepository(tx.Client())

		if _, err := txWorkflows.GetByCode(ctx, req.WorkflowCode); err == nil {
			return errx.Conflict("workflow_code already exists")
		} else if !ent.IsNotFound(err) {
			return err
		}

		createdDef, err := txWorkflows.Create(ctx, repo.WorkflowCreate{
			WorkflowCode: req.WorkflowCode,
			WorkflowName: req.WorkflowName,
			Description:  req.Description,
			Status:       "enabled",
		})
		if err != nil {
			return err
		}

		nodes := buildNodeCreates(createdDef.ID, req.Nodes)
		if _, err := txNodes.CreateBatch(ctx, nodes); err != nil {
			return err
		}
		created = createdDef
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *WorkflowService) UpdateWorkflow(ctx context.Context, req UpdateWorkflowRequest) (*ent.WorkflowDefinition, error) {
	if err := s.validateDefinition(ctx, req.WorkflowCode, req.Nodes); err != nil {
		return nil, err
	}

	var updated *ent.WorkflowDefinition
	err := db.Transact(ctx, s.ent, func(ctx context.Context, tx *ent.Tx) error {
		txWorkflows := repo.NewWorkflowRepository(tx.Client())
		txNodes := repo.NewWorkflowNodeRepository(tx.Client())

		existing, err := txWorkflows.GetByCode(ctx, req.WorkflowCode)
		if err != nil {
			if ent.IsNotFound(err) {
				return errx.NotFound("workflow not found")
			}
			return err
		}

		updatedDef, err := txWorkflows.UpdateByID(ctx, existing.ID, req.WorkflowName, req.Description)
		if err != nil {
			return err
		}

		if _, err := txNodes.DeleteByWorkflowID(ctx, existing.ID); err != nil {
			return err
		}
		nodes := buildNodeCreates(existing.ID, req.Nodes)
		if _, err := txNodes.CreateBatch(ctx, nodes); err != nil {
			return err
		}
		updated = updatedDef
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *WorkflowService) GetWorkflow(ctx context.Context, workflowCode string) (*ent.WorkflowDefinition, []*ent.WorkflowNode, error) {
	workflow, err := s.workflows.GetByCode(ctx, workflowCode)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, errx.NotFound("workflow not found")
		}
		return nil, nil, err
	}
	nodes, err := s.nodes.ListByWorkflowID(ctx, workflow.ID)
	if err != nil {
		return nil, nil, err
	}
	return workflow, nodes, nil
}

func (s *WorkflowService) validateDefinition(ctx context.Context, workflowCode string, nodes []WorkflowNodeSpec) error {
	if strings.TrimSpace(workflowCode) == "" {
		return errx.InvalidParam("workflow_code is empty")
	}
	if len(nodes) == 0 {
		return errx.InvalidParam("workflow nodes cannot be empty")
	}

	nodeMap := make(map[string]WorkflowNodeSpec, len(nodes))
	for _, n := range nodes {
		if strings.TrimSpace(n.NodeCode) == "" {
			return errx.InvalidParam("node_code is empty")
		}
		if _, exists := nodeMap[n.NodeCode]; exists {
			return errx.InvalidParam("node_code must be unique")
		}
		nodeMap[n.NodeCode] = n
		if strings.TrimSpace(n.JobCode) == "" {
			return errx.InvalidParam("job_code is empty")
		}
		if _, err := s.jobs.GetByCode(ctx, n.JobCode); err != nil {
			if ent.IsNotFound(err) {
				return errx.InvalidParam(fmt.Sprintf("job_code not found: %s", n.JobCode))
			}
			return err
		}
	}

	rootCount := 0
	for _, n := range nodes {
		if len(n.UpstreamCodes) == 0 {
			rootCount++
			continue
		}
		for _, up := range n.UpstreamCodes {
			if _, ok := nodeMap[up]; !ok {
				return errx.InvalidParam(fmt.Sprintf("upstream node not found: %s", up))
			}
		}
	}
	if rootCount == 0 {
		return errx.InvalidParam("workflow must have at least one root node")
	}

	if hasCycle(nodeMap) {
		return errx.InvalidParam("workflow has cycle")
	}
	return nil
}

func hasCycle(nodes map[string]WorkflowNodeSpec) bool {
	inDegree := make(map[string]int, len(nodes))
	graph := make(map[string][]string, len(nodes))
	for code := range nodes {
		inDegree[code] = 0
	}
	for _, node := range nodes {
		for _, up := range node.UpstreamCodes {
			graph[up] = append(graph[up], node.NodeCode)
			inDegree[node.NodeCode]++
		}
	}
	queue := make([]string, 0, len(nodes))
	for code, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, code)
		}
	}
	visited := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visited++
		for _, next := range graph[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	return visited != len(nodes)
}

func buildNodeCreates(workflowID int64, nodes []WorkflowNodeSpec) []repo.WorkflowNodeCreate {
	out := make([]repo.WorkflowNodeCreate, 0, len(nodes))
	for _, n := range nodes {
		var upstream *string
		if len(n.UpstreamCodes) > 0 {
			joined := strings.Join(n.UpstreamCodes, ",")
			upstream = &joined
		}
		trigger := n.TriggerCondition
		if trigger == "" {
			trigger = "all_success"
		}
		failStrategy := n.FailStrategy
		if failStrategy == "" {
			failStrategy = "stop"
		}
		timeout := n.TimeoutMs
		if timeout <= 0 {
			timeout = 60000
		}
		out = append(out, repo.WorkflowNodeCreate{
			WorkflowID:       workflowID,
			NodeCode:         n.NodeCode,
			NodeName:         n.NodeName,
			JobCode:          n.JobCode,
			UpstreamCodes:    upstream,
			TriggerCondition: trigger,
			FailStrategy:     failStrategy,
			TimeoutMs:        timeout,
			SortOrder:        n.SortOrder,
		})
	}
	return out
}
