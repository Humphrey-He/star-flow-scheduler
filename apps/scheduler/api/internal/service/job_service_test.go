package service

import (
	"context"
	"os"
	"testing"

	"entgo.io/ent/dialect"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/enttest"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

func TestJobServiceCreate(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping ent integration test")
	}

	client := enttest.Open(t, dialect.Postgres, dsn)
	t.Cleanup(func() { _ = client.Close() })

	svc := NewJobService(repo.NewJobRepository(client))

	expr := "* * * * *"
	req := &types.CreateJobRequest{
		JobCode:      "job_1",
		JobName:      "Job 1",
		JobType:      "cron",
		ScheduleExpr: &expr,
		ExecuteMode:  "standalone",
		HandlerName:  "handler",
	}

	created, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.TimeoutMs != 60000 {
		t.Fatalf("expected default timeout 60000, got %d", created.TimeoutMs)
	}
	if created.Priority != 5 {
		t.Fatalf("expected default priority 5, got %d", created.Priority)
	}

	_, err = svc.Create(context.Background(), req)
	if err == nil {
		t.Fatalf("expected conflict error")
	}

	be := errx.FromError(err)
	if be.Code != errx.CodeStatusConflict {
		t.Fatalf("expected conflict code, got %d", be.Code)
	}
}
