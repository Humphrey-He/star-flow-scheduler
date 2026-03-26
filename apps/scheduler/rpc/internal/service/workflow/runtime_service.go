package workflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/service/dispatch"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/db"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
)

type RuntimeService struct {
	ent               *ent.Client
	workflows         workflowRepository
	nodes             workflowNodeRepository
	nodeInstances     workflowNodeInstanceRepository
	workflowInstances workflowInstanceRepository
	jobs              jobRepository
	dispatcher        instanceCreator
	resolver          DependencyResolver
}

type workflowRepository interface {
	GetByCode(ctx context.Context, workflowCode string) (*ent.WorkflowDefinition, error)
}

type workflowNodeRepository interface {
	ListByWorkflowID(ctx context.Context, workflowID int64) ([]*ent.WorkflowNode, error)
}

type workflowNodeInstanceRepository interface {
	BatchCreate(ctx context.Context, items []pkgrepo.WorkflowNodeInstanceCreate) ([]*ent.WorkflowNodeInstance, error)
	ListByWorkflowInstanceID(ctx context.Context, workflowInstanceID int64) ([]*ent.WorkflowNodeInstance, error)
	GetByWorkflowInstanceIDAndNodeCode(ctx context.Context, workflowInstanceID int64, nodeCode string) (*ent.WorkflowNodeInstance, error)
	GetByJobInstanceID(ctx context.Context, jobInstanceID int64) (*ent.WorkflowNodeInstance, error)
	UpdateStatusIf(ctx context.Context, workflowInstanceID int64, nodeCode string, fromStatus string, toStatus string, startTime *time.Time) (int, error)
	UpdateJobInstanceID(ctx context.Context, workflowInstanceID int64, nodeCode string, jobInstanceID int64) (int, error)
}

type workflowInstanceRepository interface {
	Create(ctx context.Context, req pkgrepo.WorkflowInstanceCreate) (*ent.WorkflowInstance, error)
}

type jobRepository interface {
	GetByCode(ctx context.Context, jobCode string) (*ent.JobDefinition, error)
}

type instanceCreator interface {
	CreateInstance(ctx context.Context, jobCode string, triggerType string, payload *string) (*ent.JobInstance, error)
}

func NewRuntimeService(
	entClient *ent.Client,
	workflows workflowRepository,
	nodes workflowNodeRepository,
	nodeInstances workflowNodeInstanceRepository,
	workflowInstances workflowInstanceRepository,
	jobs jobRepository,
	dispatcher instanceCreator,
) *RuntimeService {
	return &RuntimeService{
		ent:               entClient,
		workflows:         workflows,
		nodes:             nodes,
		nodeInstances:     nodeInstances,
		workflowInstances: workflowInstances,
		jobs:              jobs,
		dispatcher:        dispatcher,
		resolver:          NewResolver(),
	}
}

func (s *RuntimeService) CreateWorkflowInstance(ctx context.Context, workflowCode string) (*ent.WorkflowInstance, []*ent.WorkflowNode, error) {
	def, err := s.workflows.GetByCode(ctx, workflowCode)
	if err != nil {
		return nil, nil, err
	}
	nodes, err := s.nodes.ListByWorkflowID(ctx, def.ID)
	if err != nil {
		return nil, nil, err
	}
	if len(nodes) == 0 {
		return nil, nil, fmt.Errorf("workflow nodes empty")
	}

	rootNodes := findRootNodes(nodes)
	if len(rootNodes) == 0 {
		return nil, nil, fmt.Errorf("workflow has no root nodes")
	}

	workflowInstanceNo := newWorkflowInstanceNo()
	var created *ent.WorkflowInstance
	err = db.Transact(ctx, s.ent, func(ctx context.Context, tx *ent.Tx) error {
		txNodes := pkgrepo.NewWorkflowNodeRepository(tx.Client())
		txNodeInstances := pkgrepo.NewWorkflowNodeInstanceRepository(tx.Client())
		txWorkflowInstances := pkgrepo.NewWorkflowInstanceRepository(tx.Client())

		instance, err := txWorkflowInstances.Create(ctx, pkgrepo.WorkflowInstanceCreate{
			WorkflowInstanceNo: workflowInstanceNo,
			WorkflowID:         def.ID,
			WorkflowCode:       def.WorkflowCode,
			Status:             string(types.WorkflowStatusRunning),
		})
		if err != nil {
			return err
		}

		nodeDefs, err := txNodes.ListByWorkflowID(ctx, def.ID)
		if err != nil {
			return err
		}
		txJobs := pkgrepo.NewJobRepository(tx.Client())
		nodeItems := make([]pkgrepo.WorkflowNodeInstanceCreate, 0, len(nodeDefs))
		for _, node := range nodeDefs {
			job, err := txJobs.GetByCode(ctx, node.JobCode)
			if err != nil {
				return err
			}
			status := string(types.WorkflowNodeStatusPending)
			if isRootNode(node) {
				status = string(types.WorkflowNodeStatusReady)
			}
			nodeItems = append(nodeItems, pkgrepo.WorkflowNodeInstanceCreate{
				WorkflowInstanceID: instance.ID,
				WorkflowID:         def.ID,
				NodeCode:           node.NodeCode,
				JobID:              int64(job.ID),
				Status:             status,
			})
		}
		if _, err := txNodeInstances.BatchCreate(ctx, nodeItems); err != nil {
			return err
		}

		created = instance
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return created, rootNodes, nil
}

func (s *RuntimeService) TriggerWorkflow(ctx context.Context, workflowCode string) (*ent.WorkflowInstance, error) {
	instance, rootNodes, err := s.CreateWorkflowInstance(ctx, workflowCode)
	if err != nil {
		return nil, err
	}
	for _, node := range rootNodes {
		jobInstance, err := s.dispatcher.CreateInstance(ctx, node.JobCode, "workflow", nil)
		if err != nil {
			continue
		}
		_, _ = s.nodeInstances.UpdateJobInstanceID(ctx, instance.ID, node.NodeCode, int64(jobInstance.ID))
		_, _ = s.nodeInstances.UpdateStatusIf(ctx, instance.ID, node.NodeCode, string(types.WorkflowNodeStatusReady), string(types.WorkflowNodeStatusRunning), timePtr(time.Now()))
	}
	return instance, nil
}

func (s *RuntimeService) OnNodeJobFinished(ctx context.Context, workflowInstanceID int64, nodeCode string, status types.WorkflowNodeStatus) error {
	nodeInst, err := s.nodeInstances.GetByWorkflowInstanceIDAndNodeCode(ctx, workflowInstanceID, nodeCode)
	if err != nil {
		return err
	}
	if nodeInst.Status != string(types.WorkflowNodeStatusRunning) {
		return nil
	}

	_, err = s.nodeInstances.UpdateStatusIf(ctx, workflowInstanceID, nodeCode, string(types.WorkflowNodeStatusRunning), string(status), timePtr(time.Now()))
	if err != nil {
		return err
	}

	nodes, err := s.nodes.ListByWorkflowID(ctx, nodeInst.WorkflowID)
	if err != nil {
		return err
	}
	nodeMap := buildNodeMap(nodes)
	current, ok := nodeMap[nodeCode]
	if !ok {
		return nil
	}
	downstreams := downstreamNodes(nodeMap, current.NodeCode)
	if len(downstreams) == 0 {
		return nil
	}

	nodeInstances, err := s.nodeInstances.ListByWorkflowInstanceID(ctx, workflowInstanceID)
	if err != nil {
		return err
	}
	nodeInstanceMap := buildNodeInstanceMap(nodeInstances)

	s.triggerDownstreamNodes(ctx, workflowInstanceID, downstreams, nodeMap, nodeInstanceMap)

	return nil
}

func (s *RuntimeService) OnJobInstanceFinished(ctx context.Context, jobInstanceID int64, status types.WorkflowNodeStatus) error {
	nodeInst, err := s.nodeInstances.GetByJobInstanceID(ctx, jobInstanceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return err
	}
	return s.OnNodeJobFinished(ctx, nodeInst.WorkflowInstanceID, nodeInst.NodeCode, status)
}

func newWorkflowInstanceNo() string {
	return fmt.Sprintf("WF-%d", time.Now().UnixNano())
}

func findRootNodes(nodes []*ent.WorkflowNode) []*ent.WorkflowNode {
	out := make([]*ent.WorkflowNode, 0)
	for _, node := range nodes {
		if isRootNode(node) {
			out = append(out, node)
		}
	}
	return out
}

func isRootNode(node *ent.WorkflowNode) bool {
	if node.UpstreamCodes == nil || strings.TrimSpace(*node.UpstreamCodes) == "" {
		return true
	}
	return false
}

func buildNodeMap(nodes []*ent.WorkflowNode) map[string]*ent.WorkflowNode {
	out := make(map[string]*ent.WorkflowNode, len(nodes))
	for _, node := range nodes {
		out[node.NodeCode] = node
	}
	return out
}

func buildNodeInstanceMap(nodes []*ent.WorkflowNodeInstance) map[string]*ent.WorkflowNodeInstance {
	out := make(map[string]*ent.WorkflowNodeInstance, len(nodes))
	for _, node := range nodes {
		out[node.NodeCode] = node
	}
	return out
}

func downstreamNodes(nodes map[string]*ent.WorkflowNode, nodeCode string) []*ent.WorkflowNode {
	out := make([]*ent.WorkflowNode, 0)
	for _, node := range nodes {
		if isUpstreamOf(node, nodeCode) {
			out = append(out, node)
		}
	}
	return out
}

func isUpstreamOf(node *ent.WorkflowNode, upstream string) bool {
	if node.UpstreamCodes == nil {
		return false
	}
	parts := strings.Split(*node.UpstreamCodes, ",")
	for _, part := range parts {
		if strings.TrimSpace(part) == upstream {
			return true
		}
	}
	return false
}

func collectUpstreamStatuses(nodes map[string]*ent.WorkflowNode, nodeInstances map[string]*ent.WorkflowNodeInstance, node *ent.WorkflowNode) []types.WorkflowNodeStatus {
	if node.UpstreamCodes == nil || strings.TrimSpace(*node.UpstreamCodes) == "" {
		return []types.WorkflowNodeStatus{}
	}
	parts := strings.Split(*node.UpstreamCodes, ",")
	out := make([]types.WorkflowNodeStatus, 0, len(parts))
	for _, part := range parts {
		code := strings.TrimSpace(part)
		if code == "" {
			continue
		}
		inst := nodeInstances[code]
		if inst == nil {
			return []types.WorkflowNodeStatus{}
		}
		out = append(out, types.WorkflowNodeStatus(inst.Status))
	}
	return out
}

func triggerConditionOf(node *ent.WorkflowNode) string {
	if node.TriggerCondition == "" {
		return "all_success"
	}
	return node.TriggerCondition
}

func (s *RuntimeService) triggerDownstreamNodes(ctx context.Context, workflowInstanceID int64, downstreams []*ent.WorkflowNode, nodeMap map[string]*ent.WorkflowNode, nodeInstanceMap map[string]*ent.WorkflowNodeInstance) {
	for _, downstream := range downstreams {
		downstreamInst := nodeInstanceMap[downstream.NodeCode]
		if downstreamInst == nil {
			continue
		}
		if downstreamInst.JobInstanceID != nil || downstreamInst.Status != string(types.WorkflowNodeStatusPending) {
			continue
		}
		upstreamStatuses := collectUpstreamStatuses(nodeMap, nodeInstanceMap, downstream)
		if len(upstreamStatuses) == 0 {
			continue
		}
		ready, err := s.resolver.CanTrigger(triggerConditionOf(downstream), upstreamStatuses)
		if err != nil || !ready {
			continue
		}

		updated, err := s.nodeInstances.UpdateStatusIf(ctx, workflowInstanceID, downstream.NodeCode, string(types.WorkflowNodeStatusPending), string(types.WorkflowNodeStatusReady), nil)
		if err != nil || updated == 0 {
			continue
		}

		jobInstance, err := s.dispatcher.CreateInstance(ctx, downstream.JobCode, "workflow", nil)
		if err != nil {
			continue
		}
		_, _ = s.nodeInstances.UpdateJobInstanceID(ctx, workflowInstanceID, downstream.NodeCode, int64(jobInstance.ID))
		_, _ = s.nodeInstances.UpdateStatusIf(ctx, workflowInstanceID, downstream.NodeCode, string(types.WorkflowNodeStatusReady), string(types.WorkflowNodeStatusRunning), timePtr(time.Now()))
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
