package repo

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/internal/models"
)

type JobInstanceRepository struct {
	db *sql.DB
}

func NewJobInstanceRepository(db *sql.DB) *JobInstanceRepository {
	return &JobInstanceRepository{db: db}
}

type JobInstanceFilter struct {
	JobCode  string
	Status   string
	StartAt  *time.Time
	EndAt    *time.Time
	Page     int
	PageSize int
}

func (r *JobInstanceRepository) List(ctx context.Context, filter JobInstanceFilter) ([]models.JobInstance, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 200 {
		filter.PageSize = 20
	}

	var where []string
	var args []any

	base := `FROM job_instances ji`
	if filter.JobCode != "" {
		base += " INNER JOIN job_definitions jd ON jd.id = ji.job_id"
		where = append(where, "jd.job_code = ?")
		args = append(args, filter.JobCode)
	}
	if filter.Status != "" {
		where = append(where, "ji.status = ?")
		args = append(args, filter.Status)
	}
	if filter.StartAt != nil {
		where = append(where, "ji.trigger_time >= ?")
		args = append(args, *filter.StartAt)
	}
	if filter.EndAt != nil {
		where = append(where, "ji.trigger_time <= ?")
		args = append(args, *filter.EndAt)
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = " WHERE " + strings.Join(where, " AND ")
	}

	var total int
	countQuery := "SELECT COUNT(1) " + base + whereSQL
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	listQuery := `SELECT ji.id, ji.instance_no, ji.job_id, ji.workflow_id, ji.trigger_type,
        ji.trigger_time, ji.scheduled_time, ji.dispatch_time, ji.start_time, ji.finish_time,
        ji.status, ji.retry_count, ji.current_backoff_ms, ji.executor_id, ji.shard_total,
        ji.success_shards, ji.failed_shards, ji.payload, ji.result_summary, ji.error_code,
        ji.error_message, ji.created_at, ji.updated_at ` + base + whereSQL +
		" ORDER BY ji.id DESC LIMIT ? OFFSET ?"

	argsWithPage := append(append([]any{}, args...), filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, argsWithPage...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []models.JobInstance
	for rows.Next() {
		var item models.JobInstance
		var workflowID sql.NullInt64
		var dispatchTime sql.NullTime
		var startTime sql.NullTime
		var finishTime sql.NullTime
		var executorID sql.NullInt64
		var payload sql.NullString
		var resultSummary sql.NullString
		var errorCode sql.NullString
		var errorMessage sql.NullString

		if err := rows.Scan(
			&item.ID, &item.InstanceNo, &item.JobID, &workflowID, &item.TriggerType,
			&item.TriggerTime, &item.ScheduledTime, &dispatchTime, &startTime, &finishTime,
			&item.Status, &item.RetryCount, &item.CurrentBackoff, &executorID, &item.ShardTotal,
			&item.SuccessShards, &item.FailedShards, &payload, &resultSummary, &errorCode,
			&errorMessage, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if workflowID.Valid {
			v := workflowID.Int64
			item.WorkflowID = &v
		}
		if dispatchTime.Valid {
			v := dispatchTime.Time
			item.DispatchTime = &v
		}
		if startTime.Valid {
			v := startTime.Time
			item.StartTime = &v
		}
		if finishTime.Valid {
			v := finishTime.Time
			item.FinishTime = &v
		}
		if executorID.Valid {
			v := executorID.Int64
			item.ExecutorID = &v
		}
		if payload.Valid {
			item.Payload = &payload.String
		}
		if resultSummary.Valid {
			item.ResultSummary = &resultSummary.String
		}
		if errorCode.Valid {
			item.ErrorCode = &errorCode.String
		}
		if errorMessage.Valid {
			item.ErrorMessage = &errorMessage.String
		}

		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *JobInstanceRepository) GetByInstanceNo(ctx context.Context, instanceNo string) (*models.JobInstance, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT id, instance_no, job_id, workflow_id, trigger_type, trigger_time, scheduled_time,
               dispatch_time, start_time, finish_time, status, retry_count, current_backoff_ms,
               executor_id, shard_total, success_shards, failed_shards, payload, result_summary,
               error_code, error_message, created_at, updated_at
        FROM job_instances WHERE instance_no = ?`, instanceNo)

	var item models.JobInstance
	var workflowID sql.NullInt64
	var dispatchTime sql.NullTime
	var startTime sql.NullTime
	var finishTime sql.NullTime
	var executorID sql.NullInt64
	var payload sql.NullString
	var resultSummary sql.NullString
	var errorCode sql.NullString
	var errorMessage sql.NullString

	if err := row.Scan(
		&item.ID, &item.InstanceNo, &item.JobID, &workflowID, &item.TriggerType, &item.TriggerTime,
		&item.ScheduledTime, &dispatchTime, &startTime, &finishTime, &item.Status,
		&item.RetryCount, &item.CurrentBackoff, &executorID, &item.ShardTotal,
		&item.SuccessShards, &item.FailedShards, &payload, &resultSummary, &errorCode,
		&errorMessage, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if workflowID.Valid {
		v := workflowID.Int64
		item.WorkflowID = &v
	}
	if dispatchTime.Valid {
		v := dispatchTime.Time
		item.DispatchTime = &v
	}
	if startTime.Valid {
		v := startTime.Time
		item.StartTime = &v
	}
	if finishTime.Valid {
		v := finishTime.Time
		item.FinishTime = &v
	}
	if executorID.Valid {
		v := executorID.Int64
		item.ExecutorID = &v
	}
	if payload.Valid {
		item.Payload = &payload.String
	}
	if resultSummary.Valid {
		item.ResultSummary = &resultSummary.String
	}
	if errorCode.Valid {
		item.ErrorCode = &errorCode.String
	}
	if errorMessage.Valid {
		item.ErrorMessage = &errorMessage.String
	}

	return &item, nil
}
