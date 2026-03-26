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
	workflows         *pkgrepo.WorkflowRepository
	nodes             *pkgrepo.WorkflowNodeRepository
	nodeInstances     *pkgrepo.WorkflowNodeInstanceRepository
	workflowInstances *pkgrepo.WorkflowInstanceRepository
	jobs              *pkgrepo.JobRepository
	dispatcher        *dispatch.Service
}

func NewRuntimeService(
	entClient *ent.Client,
	workflows *pkgrepo.WorkflowRepository,
	nodes *pkgrepo.WorkflowNodeRepository,
	nodeInstances *pkgrepo.WorkflowNodeInstanceRepository,
	workflowInstances *pkgrepo.WorkflowInstanceRepository,
	jobs *pkgrepo.JobRepository,
	dispatcher *dispatch.Service,
) *RuntimeService {
	return &RuntimeService{
		ent:               entClient,
		workflows:         workflows,
		nodes:             nodes,
		nodeInstances:     nodeInstances,
		workflowInstances: workflowInstances,
		jobs:              jobs,
		dispatcher:        dispatcher,
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

func timePtr(t time.Time) *time.Time {
	return &t
}
