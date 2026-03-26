package runtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/handler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/model"
	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"
)

type testHandler struct {
	err error
}

func (h *testHandler) Name() string {
	return "test"
}

func (h *testHandler) Execute(ctx context.Context, payload []byte) error {
	return h.err
}

type testReporter struct {
	ch chan *model.TaskResult
}

func (r *testReporter) Report(result *model.TaskResult) {
	r.ch <- result
}

func TestRuntimeExecutesTask(t *testing.T) {
	reg := handler.NewRegistry()
	_ = reg.Register(&testHandler{})

	rep := &testReporter{ch: make(chan *model.TaskResult, 1)}
	rt := NewRuntime(Config{
		WorkerCount:      1,
		QueueSize:        1,
		DefaultTimeoutMs: 1000,
	}, reg, rep)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rt.Start(ctx)

	err := rt.Enqueue(&model.Task{
		InstanceNo:  "inst-1",
		HandlerName: "test",
	})
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	select {
	case res := <-rep.ch:
		if res.Status != schedulev1.InstanceStatus_INSTANCE_STATUS_SUCCESS {
			t.Fatalf("expected success got %v", res.Status)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for result")
	}
}

func TestRuntimeHandlesHandlerError(t *testing.T) {
	reg := handler.NewRegistry()
	_ = reg.Register(&testHandler{err: errors.New("boom")})

	rep := &testReporter{ch: make(chan *model.TaskResult, 1)}
	rt := NewRuntime(Config{
		WorkerCount:      1,
		QueueSize:        1,
		DefaultTimeoutMs: 1000,
	}, reg, rep)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rt.Start(ctx)

	err := rt.Enqueue(&model.Task{
		InstanceNo:  "inst-2",
		HandlerName: "test",
	})
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	select {
	case res := <-rep.ch:
		if res.Status != schedulev1.InstanceStatus_INSTANCE_STATUS_FAILED {
			t.Fatalf("expected failed got %v", res.Status)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for result")
	}
}

func TestQueueFull(t *testing.T) {
	reg := handler.NewRegistry()
	rep := &testReporter{ch: make(chan *model.TaskResult, 1)}
	rt := NewRuntime(Config{
		WorkerCount:      1,
		QueueSize:        1,
		DefaultTimeoutMs: 1000,
	}, reg, rep)

	if err := rt.Enqueue(&model.Task{InstanceNo: "1"}); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}
	if err := rt.Enqueue(&model.Task{InstanceNo: "2"}); !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected queue full error")
	}
}

func TestRuntimeDrainRejectsEnqueue(t *testing.T) {
	reg := handler.NewRegistry()
	rep := &testReporter{ch: make(chan *model.TaskResult, 1)}
	rt := NewRuntime(Config{
		WorkerCount:        1,
		QueueSize:          1,
		DefaultTimeoutMs:   1000,
		ShutdownTimeoutSec: 1,
	}, reg, rep)

	rt.Drain()
	if err := rt.Enqueue(&model.Task{InstanceNo: "3"}); err == nil {
		t.Fatalf("expected draining error")
	}
}
