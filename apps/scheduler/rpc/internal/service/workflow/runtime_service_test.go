package workflow

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
)

type fakeWorkflowRepo struct {
	def *ent.WorkflowDefinition
}

func (f *fakeWorkflowRepo) GetByCode(_ context.Context, _ string) (*ent.WorkflowDefinition, error) {
	return f.def, nil
}

type fakeWorkflowNodeRepo struct {
	nodes []*ent.WorkflowNode
}

func (f *fakeWorkflowNodeRepo) ListByWorkflowID(_ context.Context, _ int64) ([]*ent.WorkflowNode, error) {
	return f.nodes, nil
}

type fakeWorkflowInstanceRepo struct {
	nextID int64
}

func (f *fakeWorkflowInstanceRepo) Create(_ context.Context, req pkgrepo.WorkflowInstanceCreate) (*ent.WorkflowInstance, error) {
	f.nextID++
	return &ent.WorkflowInstance{
		ID:                 f.nextID,
		WorkflowInstanceNo: req.WorkflowInstanceNo,
		WorkflowID:         req.WorkflowID,
		WorkflowCode:       req.WorkflowCode,
		Status:             req.Status,
	}, nil
}

type fakeJobRepo struct {
	jobs map[string]*ent.JobDefinition
}

func (f *fakeJobRepo) GetByCode(_ context.Context, jobCode string) (*ent.JobDefinition, error) {
	return f.jobs[jobCode], nil
}

type fakeNodeInstanceRepo struct {
	mu    sync.Mutex
	items map[string]*ent.WorkflowNodeInstance
}

func newFakeNodeInstanceRepo() *fakeNodeInstanceRepo {
	return &fakeNodeInstanceRepo{items: make(map[string]*ent.WorkflowNodeInstance)}
}

func (r *fakeNodeInstanceRepo) key(workflowInstanceID int64, nodeCode string) string {
	return stringKey(workflowInstanceID, nodeCode)
}

func (r *fakeNodeInstanceRepo) BatchCreate(_ context.Context, items []pkgrepo.WorkflowNodeInstanceCreate) ([]*ent.WorkflowNodeInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*ent.WorkflowNodeInstance, 0, len(items))
	for _, item := range items {
		inst := &ent.WorkflowNodeInstance{
			WorkflowInstanceID: item.WorkflowInstanceID,
			WorkflowID:         item.WorkflowID,
			NodeCode:           item.NodeCode,
			JobID:              item.JobID,
			Status:             item.Status,
		}
		if item.JobInstanceID != nil {
			inst.JobInstanceID = item.JobInstanceID
		}
		r.items[r.key(item.WorkflowInstanceID, item.NodeCode)] = inst
		out = append(out, inst)
	}
	return out, nil
}

func (r *fakeNodeInstanceRepo) ListByWorkflowInstanceID(_ context.Context, workflowInstanceID int64) ([]*ent.WorkflowNodeInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*ent.WorkflowNodeInstance, 0)
	for _, item := range r.items {
		if item.WorkflowInstanceID == workflowInstanceID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *fakeNodeInstanceRepo) GetByWorkflowInstanceIDAndNodeCode(_ context.Context, workflowInstanceID int64, nodeCode string) (*ent.WorkflowNodeInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items[r.key(workflowInstanceID, nodeCode)], nil
}

func (r *fakeNodeInstanceRepo) GetByJobInstanceID(_ context.Context, jobInstanceID int64) (*ent.WorkflowNodeInstance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range r.items {
		if item.JobInstanceID != nil && *item.JobInstanceID == jobInstanceID {
			return item, nil
		}
	}
	return nil, ent.ErrNotFound
}

func (r *fakeNodeInstanceRepo) UpdateStatusIf(_ context.Context, workflowInstanceID int64, nodeCode string, fromStatus string, toStatus string, startTime *time.Time) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item := r.items[r.key(workflowInstanceID, nodeCode)]
	if item == nil || item.Status != fromStatus {
		return 0, nil
	}
	item.Status = toStatus
	if startTime != nil {
		item.StartTime = startTime
	}
	return 1, nil
}

func (r *fakeNodeInstanceRepo) UpdateJobInstanceID(_ context.Context, workflowInstanceID int64, nodeCode string, jobInstanceID int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item := r.items[r.key(workflowInstanceID, nodeCode)]
	if item == nil {
		return 0, nil
	}
	item.JobInstanceID = &jobInstanceID
	return 1, nil
}

type fakeDispatcher struct {
	mu    sync.Mutex
	calls int
	next  int64
}

func (f *fakeDispatcher) CreateInstance(_ context.Context, _ string, _ string, _ *string) (*ent.JobInstance, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	f.next++
	return &ent.JobInstance{ID: int(f.next)}, nil
}

func (f *fakeDispatcher) Calls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func stringKey(instanceID int64, nodeCode string) string {
	return fmtKey(instanceID) + ":" + nodeCode
}

func fmtKey(id int64) string {
	return strconv.FormatInt(id, 10)
}

func TestOnNodeJobFinishedIdempotent(t *testing.T) {
	nodes := []*ent.WorkflowNode{
		{NodeCode: "A", JobCode: "jobA"},
		{NodeCode: "B", JobCode: "jobB", UpstreamCodes: strPtr("A")},
	}
	nodeRepo := &fakeWorkflowNodeRepo{nodes: nodes}
	nodeInstRepo := newFakeNodeInstanceRepo()
	workflowRepo := &fakeWorkflowRepo{def: &ent.WorkflowDefinition{ID: 1, WorkflowCode: "wf"}}
	instanceRepo := &fakeWorkflowInstanceRepo{}
	jobs := &fakeJobRepo{
		jobs: map[string]*ent.JobDefinition{
			"jobA": {ID: 1, JobCode: "jobA"},
			"jobB": {ID: 2, JobCode: "jobB"},
		},
	}
	dispatcher := &fakeDispatcher{}

	svc := NewRuntimeService(nil, workflowRepo, nodeRepo, nodeInstRepo, instanceRepo, jobs, dispatcher)

	instance, roots, err := svc.CreateWorkflowInstance(context.Background(), "wf")
	if err != nil {
		t.Fatalf("create workflow instance err: %v", err)
	}
	if len(roots) != 1 || roots[0].NodeCode != "A" {
		t.Fatalf("expected root A")
	}
	_, _ = nodeInstRepo.UpdateStatusIf(context.Background(), instance.ID, "A", string(types.WorkflowNodeStatusReady), string(types.WorkflowNodeStatusRunning), nil)

	if err := svc.OnNodeJobFinished(context.Background(), instance.ID, "A", types.WorkflowNodeStatusSuccess); err != nil {
		t.Fatalf("on node finished err: %v", err)
	}
	if err := svc.OnNodeJobFinished(context.Background(), instance.ID, "A", types.WorkflowNodeStatusSuccess); err != nil {
		t.Fatalf("on node finished err: %v", err)
	}

	if dispatcher.Calls() != 1 {
		t.Fatalf("expected downstream triggered once got %d", dispatcher.Calls())
	}
}

func TestTriggerConditionAnySuccess(t *testing.T) {
	nodes := []*ent.WorkflowNode{
		{NodeCode: "A", JobCode: "jobA"},
		{NodeCode: "B", JobCode: "jobB"},
		{NodeCode: "C", JobCode: "jobC", UpstreamCodes: strPtr("A,B"), TriggerCondition: "any_success"},
	}
	nodeRepo := &fakeWorkflowNodeRepo{nodes: nodes}
	nodeInstRepo := newFakeNodeInstanceRepo()
	workflowRepo := &fakeWorkflowRepo{def: &ent.WorkflowDefinition{ID: 1, WorkflowCode: "wf"}}
	instanceRepo := &fakeWorkflowInstanceRepo{}
	jobs := &fakeJobRepo{
		jobs: map[string]*ent.JobDefinition{
			"jobA": {ID: 1, JobCode: "jobA"},
			"jobB": {ID: 2, JobCode: "jobB"},
			"jobC": {ID: 3, JobCode: "jobC"},
		},
	}
	dispatcher := &fakeDispatcher{}

	svc := NewRuntimeService(nil, workflowRepo, nodeRepo, nodeInstRepo, instanceRepo, jobs, dispatcher)

	instance, _, err := svc.CreateWorkflowInstance(context.Background(), "wf")
	if err != nil {
		t.Fatalf("create workflow instance err: %v", err)
	}
	_, _ = nodeInstRepo.UpdateStatusIf(context.Background(), instance.ID, "A", string(types.WorkflowNodeStatusReady), string(types.WorkflowNodeStatusRunning), nil)
	_, _ = nodeInstRepo.UpdateStatusIf(context.Background(), instance.ID, "B", string(types.WorkflowNodeStatusReady), string(types.WorkflowNodeStatusRunning), nil)

	if err := svc.OnNodeJobFinished(context.Background(), instance.ID, "A", types.WorkflowNodeStatusSuccess); err != nil {
		t.Fatalf("on node finished err: %v", err)
	}
	if dispatcher.Calls() != 1 {
		t.Fatalf("expected downstream triggered once got %d", dispatcher.Calls())
	}
}

func strPtr(v string) *string {
	return &v
}
