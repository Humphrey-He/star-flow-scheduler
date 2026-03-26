package model

import (
	"time"

	schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"
)

type TaskResult struct {
	InstanceNo    string
	ShardNo       string
	Status        schedulev1.InstanceStatus
	StartTime     time.Time
	FinishTime    time.Time
	ErrorCode     string
	ErrorMessage  string
	ResultSummary string
}
