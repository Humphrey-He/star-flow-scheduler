package validator

import (
	"strings"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
)

func ValidateCreateJob(req *types.CreateJobRequest) error {
	if strings.TrimSpace(req.JobCode) == "" {
		return errx.InvalidParam("job_code is required")
	}
	if strings.TrimSpace(req.JobName) == "" {
		return errx.InvalidParam("job_name is required")
	}
	if strings.TrimSpace(req.JobType) == "" {
		return errx.InvalidParam("job_type is required")
	}
	if req.JobType == "cron" {
		if req.ScheduleExpr == nil || strings.TrimSpace(*req.ScheduleExpr) == "" {
			return errx.InvalidParam("schedule_expr is required for cron jobs")
		}
	}
	if req.JobType == "delay" {
		if req.DelayMs == nil || *req.DelayMs <= 0 {
			return errx.InvalidParam("delay_ms must be > 0 for delay jobs")
		}
	}
	if req.JobType == "once" {
		if req.ScheduleExpr != nil && strings.TrimSpace(*req.ScheduleExpr) != "" {
			return errx.InvalidParam("schedule_expr must be empty for once jobs")
		}
	}
	if req.ExecuteMode == "sharding" && req.ShardTotal < 2 {
		return errx.InvalidParam("shard_total must be >= 2 for sharding mode")
	}
	if req.JobType == "dag" && strings.TrimSpace(req.HandlerName) != "" {
		return errx.InvalidParam("handler_name should be empty for dag jobs")
	}
	if req.JobType != "dag" && strings.TrimSpace(req.HandlerName) == "" {
		return errx.InvalidParam("handler_name is required")
	}

	return nil
}
