package validator

import (
	"testing"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
)

func TestValidateCreateJob(t *testing.T) {
	makeReq := func() *types.CreateJobRequest {
		return &types.CreateJobRequest{
			JobCode:     "job_1",
			JobName:     "Job 1",
			JobType:     "cron",
			ExecuteMode: "standalone",
			HandlerName: "handler",
		}
	}

	req := makeReq()
	if err := ValidateCreateJob(req); err == nil {
		t.Fatalf("expected error for missing schedule_expr")
	}

	expr := "* * * * *"
	req.ScheduleExpr = &expr
	if err := ValidateCreateJob(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req.JobType = "delay"
	if err := ValidateCreateJob(req); err == nil {
		t.Fatalf("expected error for missing delay_ms")
	}

	delay := int64(1000)
	req.DelayMs = &delay
	if err := ValidateCreateJob(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req.ExecuteMode = "sharding"
	req.ShardTotal = 1
	if err := ValidateCreateJob(req); err == nil {
		t.Fatalf("expected error for shard_total")
	}

	req.ShardTotal = 2
	if err := ValidateCreateJob(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req.JobType = "dag"
	if err := ValidateCreateJob(req); err == nil {
		t.Fatalf("expected error for dag handler_name")
	}
}
