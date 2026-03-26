package dispatchservicelogic

import (
	"context"
	"testing"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"
)

type fakeInstanceSvc struct {
	calls int
}

func (f *fakeInstanceSvc) ReportResult(_ context.Context, _ string, _ state.InstanceStatus, _ *time.Time, _ *time.Time, _ *string, _ *string, _ *string) (int, error) {
	f.calls++
	return 1, nil
}

type fakeInstanceRepo struct {
	instance *ent.JobInstance
}

func (f *fakeInstanceRepo) GetByInstanceNo(_ context.Context, _ string) (*ent.JobInstance, error) {
	return f.instance, nil
}

type fakeWorkflowRuntime struct {
	calls int
	last  int64
}

func (f *fakeWorkflowRuntime) OnJobInstanceFinished(_ context.Context, jobInstanceID int64, _ types.WorkflowNodeStatus) error {
	f.calls++
	f.last = jobInstanceID
	return nil
}

func TestReportResultTriggersWorkflow(t *testing.T) {
	workflowID := int64(10)
	instance := &ent.JobInstance{ID: 12, WorkflowID: &workflowID}
	instanceRepo := &fakeInstanceRepo{instance: instance}
	instanceSvc := &fakeInstanceSvc{}
	workflowRuntime := &fakeWorkflowRuntime{}

	svcCtx := &svc.ServiceContext{
		InstanceRepo:    instanceRepo,
		InstanceSvc:     instanceSvc,
		WorkflowRuntime: workflowRuntime,
	}

	logic := NewReportResultLogic(context.Background(), svcCtx)
	_, err := logic.ReportResult(&schedulev1.ReportResultRequest{
		InstanceNo: "inst-1",
		Status:     schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS,
	})
	if err != nil {
		t.Fatalf("report result err: %v", err)
	}
	if workflowRuntime.calls != 1 {
		t.Fatalf("expected workflow runtime called once, got %d", workflowRuntime.calls)
	}
	if workflowRuntime.last != int64(instance.ID) {
		t.Fatalf("expected workflow runtime job id %d got %d", instance.ID, workflowRuntime.last)
	}
}

func TestReportResultNoWorkflow(t *testing.T) {
	instance := &ent.JobInstance{ID: 15}
	instanceRepo := &fakeInstanceRepo{instance: instance}
	instanceSvc := &fakeInstanceSvc{}
	workflowRuntime := &fakeWorkflowRuntime{}

	svcCtx := &svc.ServiceContext{
		InstanceRepo:    instanceRepo,
		InstanceSvc:     instanceSvc,
		WorkflowRuntime: workflowRuntime,
	}

	logic := NewReportResultLogic(context.Background(), svcCtx)
	_, err := logic.ReportResult(&schedulev1.ReportResultRequest{
		InstanceNo: "inst-2",
		Status:     schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS,
	})
	if err != nil {
		t.Fatalf("report result err: %v", err)
	}
	if workflowRuntime.calls != 0 {
		t.Fatalf("expected workflow runtime not called, got %d", workflowRuntime.calls)
	}
}
