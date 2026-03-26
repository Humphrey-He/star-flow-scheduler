package models

import "time"

type JobDefinition struct {
	ID                int64
	JobCode           string
	JobName           string
	JobType           string
	ScheduleExpr      *string
	DelayMs           *int64
	ExecuteMode       string
	HandlerName       string
	HandlerPayload    *string
	TimeoutMs         int
	RetryLimit        int
	RetryBackoff      string
	Priority          int
	ShardTotal        int
	RouteStrategy     string
	ExecutorTag       *string
	IdempotentKeyExpr *string
	Status            string
	CreatedBy         *string
	UpdatedBy         *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type JobInstance struct {
	ID             int64
	InstanceNo     string
	JobID          int64
	WorkflowID     *int64
	TriggerType    string
	TriggerTime    time.Time
	ScheduledTime  time.Time
	DispatchTime   *time.Time
	StartTime      *time.Time
	FinishTime     *time.Time
	Status         string
	RetryCount     int
	CurrentBackoff int64
	ExecutorID     *int64
	ShardTotal     int
	SuccessShards  int
	FailedShards   int
	Payload        *string
	ResultSummary  *string
	ErrorCode      *string
	ErrorMessage   *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Executor struct {
	ID            int64
	ExecutorCode  string
	Host          string
	IP            string
	GrpcAddr      string
	HttpAddr      *string
	Tags          *string
	Capacity      int
	CurrentLoad   int
	Version       *string
	Status        string
	LastHeartbeat time.Time
	Metadata      *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
