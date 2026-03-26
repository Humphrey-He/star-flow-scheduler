package httpapi

import (
    "context"
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "time"

    "starflow-scheduler/internal/models"
    "starflow-scheduler/internal/repo"
)

type Server struct {
    db       *sql.DB
    jobs     *repo.JobRepository
    instances *repo.JobInstanceRepository
}

func NewServer(db *sql.DB) *Server {
    return &Server{
        db:        db,
        jobs:      repo.NewJobRepository(db),
        instances: repo.NewJobInstanceRepository(db),
    }
}

func (s *Server) Routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/api/v1/jobs", s.handleJobs)
    mux.HandleFunc("/api/v1/jobs/", s.handleJobByCode)
    mux.HandleFunc("/api/v1/job-instances", s.handleJobInstances)
    mux.HandleFunc("/api/v1/job-instances/", s.handleJobInstanceByNo)

    return mux
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        s.createJobDefinition(w, r)
    default:
        writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
    }
}

func (s *Server) handleJobByCode(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
        return
    }
    jobCode := strings.TrimPrefix(r.URL.Path, "/api/v1/jobs/")
    if jobCode == "" {
        writeError(w, http.StatusBadRequest, "invalid_request", "job_code is required")
        return
    }

    ctx := r.Context()
    job, err := s.jobs.GetByCode(ctx, jobCode)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            writeError(w, http.StatusNotFound, "not_found", "job not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "db_error", err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "job": job,
    })
}

func (s *Server) handleJobInstances(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
        return
    }

    q := r.URL.Query()
    filter := repo.JobInstanceFilter{
        JobCode:  q.Get("job_code"),
        Status:   q.Get("status"),
        Page:     parseIntDefault(q.Get("page"), 1),
        PageSize: parseIntDefault(q.Get("page_size"), 20),
    }

    if start := q.Get("start_time"); start != "" {
        if t, err := parseTime(start); err == nil {
            filter.StartAt = &t
        }
    }
    if end := q.Get("end_time"); end != "" {
        if t, err := parseTime(end); err == nil {
            filter.EndAt = &t
        }
    }

    ctx := r.Context()
    items, total, err := s.instances.List(ctx, filter)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "db_error", err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "items": items,
        "total": total,
        "page": filter.Page,
        "page_size": filter.PageSize,
    })
}

func (s *Server) handleJobInstanceByNo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
        return
    }
    instanceNo := strings.TrimPrefix(r.URL.Path, "/api/v1/job-instances/")
    if instanceNo == "" {
        writeError(w, http.StatusBadRequest, "invalid_request", "instance_no is required")
        return
    }

    ctx := r.Context()
    item, err := s.instances.GetByInstanceNo(ctx, instanceNo)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            writeError(w, http.StatusNotFound, "not_found", "instance not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "db_error", err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "instance": item,
    })
}

func (s *Server) createJobDefinition(w http.ResponseWriter, r *http.Request) {
    var req struct {
        JobCode        string          `json:"job_code"`
        JobName        string          `json:"job_name"`
        JobType        string          `json:"job_type"`
        ScheduleExpr   *string         `json:"schedule_expr"`
        DelayMs        *int64          `json:"delay_ms"`
        ExecuteMode    string          `json:"execute_mode"`
        HandlerName    string          `json:"handler_name"`
        Payload        json.RawMessage `json:"payload"`
        TimeoutMs      int             `json:"timeout_ms"`
        RetryLimit     int             `json:"retry_limit"`
        RetryBackoff   string          `json:"retry_backoff"`
        Priority       int             `json:"priority"`
        ShardTotal     int             `json:"shard_total"`
        RouteStrategy  string          `json:"route_strategy"`
        ExecutorTag    *string         `json:"executor_tag"`
        IdempotentExpr *string         `json:"idempotent_key_expr"`
        Status         string          `json:"status"`
        CreatedBy      *string         `json:"created_by"`
        UpdatedBy      *string         `json:"updated_by"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
        return
    }

    job := &models.JobDefinition{
        JobCode:           req.JobCode,
        JobName:           req.JobName,
        JobType:           req.JobType,
        ScheduleExpr:      req.ScheduleExpr,
        DelayMs:           req.DelayMs,
        ExecuteMode:       req.ExecuteMode,
        HandlerName:       req.HandlerName,
        TimeoutMs:         req.TimeoutMs,
        RetryLimit:        req.RetryLimit,
        RetryBackoff:      req.RetryBackoff,
        Priority:          req.Priority,
        ShardTotal:        req.ShardTotal,
        RouteStrategy:     req.RouteStrategy,
        ExecutorTag:       req.ExecutorTag,
        IdempotentKeyExpr: req.IdempotentExpr,
        Status:            req.Status,
        CreatedBy:         req.CreatedBy,
        UpdatedBy:         req.UpdatedBy,
    }

    if len(req.Payload) > 0 && string(req.Payload) != "null" {
        payload := string(req.Payload)
        job.HandlerPayload = &payload
    }

    repo.BuildJobDefinitionDefaults(job)
    if err := repo.ValidateJobDefinition(job); err != nil {
        writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    exists, err := s.jobs.ExistsByCode(ctx, job.JobCode)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "db_error", err.Error())
        return
    }
    if exists {
        writeError(w, http.StatusConflict, "already_exists", "job_code already exists")
        return
    }

    id, err := s.jobs.Create(ctx, job)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "db_error", err.Error())
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "job_id":   id,
        "job_code": job.JobCode,
    })
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
    writeJSON(w, status, map[string]any{
        "error": map[string]any{
            "code":    code,
            "message": message,
        },
    })
}

func parseIntDefault(raw string, def int) int {
    if raw == "" {
        return def
    }
    var n int
    if _, err := fmt.Sscanf(raw, "%d", &n); err != nil {
        return def
    }
    return n
}

func parseTime(raw string) (time.Time, error) {
    if t, err := time.Parse(time.RFC3339, raw); err == nil {
        return t, nil
    }
    return time.Parse("2006-01-02 15:04:05", raw)
}
