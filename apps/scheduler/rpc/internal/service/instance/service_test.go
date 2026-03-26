package instance

import (
	"context"
	"testing"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
)

type fakeRepo struct {
	status  string
	updated int
	err     error
}

func (f *fakeRepo) GetStatusByInstanceNo(ctx context.Context, instanceNo string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.status, nil
}

func (f *fakeRepo) UpdateStatusIf(ctx context.Context, instanceNo string, fromStatus string, toStatus string) (int, error) {
	if f.status != fromStatus {
		return 0, nil
	}
	f.status = toStatus
	f.updated++
	return 1, nil
}

func (f *fakeRepo) UpdateResultIfStatus(ctx context.Context, instanceNo string, fromStatus string, toStatus string, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	if f.status != fromStatus {
		return 0, nil
	}
	f.status = toStatus
	f.updated++
	return 1, nil
}

func TestReportResultIdempotent(t *testing.T) {
	repo := &fakeRepo{status: string(state.StatusSuccess)}
	svc := NewService(repo)

	rows, err := svc.ReportResult(context.Background(), "inst-1", state.StatusSuccess, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if rows != 0 {
		t.Fatalf("expected 0 rows, got %d", rows)
	}
}

func TestReportResultInvalidTransition(t *testing.T) {
	repo := &fakeRepo{status: string(state.StatusSuccess)}
	svc := NewService(repo)

	_, err := svc.ReportResult(context.Background(), "inst-1", state.StatusFailed, nil, nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestReportResultConflict(t *testing.T) {
	repo := &fakeRepo{status: string(state.StatusRunning)}
	svc := NewService(repo)

	repo.status = string(state.StatusDispatched)
	_, err := svc.ReportResult(context.Background(), "inst-1", state.StatusSuccess, nil, nil, nil)
	if err == nil {
		t.Fatalf("expected conflict error")
	}
}

func TestTransitionConditionalUpdate(t *testing.T) {
	repo := &fakeRepo{status: string(state.StatusPending)}
	svc := NewService(repo)

	updated, err := svc.Transition(context.Background(), "inst-1", state.StatusPending, state.StatusDispatched)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !updated {
		t.Fatalf("expected update true")
	}

	updated, err = svc.Transition(context.Background(), "inst-1", state.StatusPending, state.StatusDispatched)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if updated {
		t.Fatalf("expected update false for conditional update")
	}
}
