package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/internal/models"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *models.JobDefinition) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
        INSERT INTO job_definitions (
            job_code, job_name, job_type, schedule_expr, delay_ms, execute_mode,
            handler_name, handler_payload, timeout_ms, retry_limit, retry_backoff,
            priority, shard_total, route_strategy, executor_tag, idempotent_key_expr,
            status, created_by, updated_by
        ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		job.JobCode, job.JobName, job.JobType, job.ScheduleExpr, job.DelayMs, job.ExecuteMode,
		job.HandlerName, job.HandlerPayload, job.TimeoutMs, job.RetryLimit, job.RetryBackoff,
		job.Priority, job.ShardTotal, job.RouteStrategy, job.ExecutorTag, job.IdempotentKeyExpr,
		job.Status, job.CreatedBy, job.UpdatedBy,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *JobRepository) GetByCode(ctx context.Context, jobCode string) (*models.JobDefinition, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT id, job_code, job_name, job_type, schedule_expr, delay_ms, execute_mode,
               handler_name, handler_payload, timeout_ms, retry_limit, retry_backoff,
               priority, shard_total, route_strategy, executor_tag, idempotent_key_expr,
               status, created_by, updated_by, created_at, updated_at
        FROM job_definitions WHERE job_code = ?`, jobCode)

	var job models.JobDefinition
	var scheduleExpr sql.NullString
	var delayMs sql.NullInt64
	var handlerPayload sql.NullString
	var executorTag sql.NullString
	var idemKeyExpr sql.NullString
	var createdBy sql.NullString
	var updatedBy sql.NullString

	if err := row.Scan(
		&job.ID, &job.JobCode, &job.JobName, &job.JobType, &scheduleExpr, &delayMs, &job.ExecuteMode,
		&job.HandlerName, &handlerPayload, &job.TimeoutMs, &job.RetryLimit, &job.RetryBackoff,
		&job.Priority, &job.ShardTotal, &job.RouteStrategy, &executorTag, &idemKeyExpr,
		&job.Status, &createdBy, &updatedBy, &job.CreatedAt, &job.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if scheduleExpr.Valid {
		job.ScheduleExpr = &scheduleExpr.String
	}
	if delayMs.Valid {
		v := delayMs.Int64
		job.DelayMs = &v
	}
	if handlerPayload.Valid {
		job.HandlerPayload = &handlerPayload.String
	}
	if executorTag.Valid {
		job.ExecutorTag = &executorTag.String
	}
	if idemKeyExpr.Valid {
		job.IdempotentKeyExpr = &idemKeyExpr.String
	}
	if createdBy.Valid {
		job.CreatedBy = &createdBy.String
	}
	if updatedBy.Valid {
		job.UpdatedBy = &updatedBy.String
	}

	return &job, nil
}

func (r *JobRepository) ExistsByCode(ctx context.Context, jobCode string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM job_definitions WHERE job_code = ?", jobCode).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ValidateJobDefinition(job *models.JobDefinition) error {
	if job.JobCode == "" {
		return fmt.Errorf("job_code is required")
	}
	if job.JobName == "" {
		return fmt.Errorf("job_name is required")
	}
	if job.JobType == "" {
		return fmt.Errorf("job_type is required")
	}
	if job.JobType != "dag" && job.HandlerName == "" {
		return fmt.Errorf("handler_name is required")
	}

	switch job.JobType {
	case "cron":
		if job.ScheduleExpr == nil || *job.ScheduleExpr == "" {
			return fmt.Errorf("schedule_expr is required for cron jobs")
		}
	case "delay":
		if job.DelayMs == nil || *job.DelayMs <= 0 {
			return fmt.Errorf("delay_ms must be > 0 for delay jobs")
		}
	case "once":
		if job.ScheduleExpr != nil && *job.ScheduleExpr != "" {
			return fmt.Errorf("schedule_expr must be empty for once jobs")
		}
	case "dag":
		// DAG jobs should be triggered by workflow nodes instead of handler binding.
		if job.HandlerName != "" {
			return fmt.Errorf("handler_name should be empty for dag jobs")
		}
	default:
		return fmt.Errorf("unsupported job_type: %s", job.JobType)
	}

	return nil
}

func BuildJobDefinitionDefaults(job *models.JobDefinition) {
	if job.ExecuteMode == "" {
		job.ExecuteMode = "standalone"
	}
	if job.TimeoutMs == 0 {
		job.TimeoutMs = 60000
	}
	if job.RetryLimit == 0 {
		job.RetryLimit = 3
	}
	if job.RetryBackoff == "" {
		job.RetryBackoff = "1s,3s,5s"
	}
	if job.Priority == 0 {
		job.Priority = 5
	}
	if job.ShardTotal == 0 {
		job.ShardTotal = 1
	}
	if job.RouteStrategy == "" {
		job.RouteStrategy = "least_load"
	}
	if job.Status == "" {
		job.Status = "enabled"
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
}
