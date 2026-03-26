package dispatch

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	rpcrepo "github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	pkgrepo "github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	schedulerv1_schedulev1 "github.com/Humphrey-He/star-flow-scheduler/proto/pb/github.com/Humphrey-He/star-flow-scheduler/proto/schedulerv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	jobs      *rpcrepo.JobRepository
	instances *rpcrepo.JobInstanceRepository
	executors *rpcrepo.ExecutorRepository
}

func NewService(jobs *rpcrepo.JobRepository, instances *rpcrepo.JobInstanceRepository, executors *rpcrepo.ExecutorRepository) *Service {
	return &Service{
		jobs:      jobs,
		instances: instances,
		executors: executors,
	}
}

func (s *Service) CreateInstance(ctx context.Context, jobCode string, triggerType string, payload *string) (*ent.JobInstance, error) {
	job, err := s.jobs.GetByCode(ctx, jobCode)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	instanceNo := newInstanceNo()
	req := pkgrepo.JobInstanceCreate{
		InstanceNo:    instanceNo,
		JobID:         int64(job.ID),
		TriggerType:   triggerType,
		TriggerTime:   now,
		ScheduledTime: now,
		Status:        string(state.StatusPending),
		Payload:       payload,
		ShardTotal:    job.ShardTotal,
	}

	return s.instances.Create(ctx, req)
}

func (s *Service) DispatchInstance(ctx context.Context, instanceNo string) (*ent.Executor, error) {
	instance, err := s.instances.GetByInstanceNo(ctx, instanceNo)
	if err != nil {
		return nil, err
	}

	job, err := s.jobs.GetByID(ctx, instance.JobID)
	if err != nil {
		return nil, err
	}

	execs, err := s.executors.ListOnline(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(execs) == 0 {
		return nil, fmt.Errorf("no online executor")
	}

	exec := execs[0]
	if err := s.dispatchToExecutor(ctx, exec, instance, job); err != nil {
		return nil, err
	}

	_, err = s.instances.UpdateDispatchInfoIfStatus(ctx, instanceNo, string(state.StatusPending), string(state.StatusDispatched), int64(exec.ID), time.Now())
	if err != nil {
		return nil, err
	}

	return exec, nil
}

func (s *Service) dispatchToExecutor(ctx context.Context, exec *ent.Executor, instance *ent.JobInstance, job *ent.JobDefinition) error {
	conn, err := grpc.DialContext(ctx, exec.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := schedulerv1_schedulev1.NewDispatchServiceClient(conn)

	payload := &schedulerv1_schedulev1.JobPayload{Raw: []byte(ptrString(instance.Payload))}
	_, err = client.DispatchJob(ctx, &schedulerv1_schedulev1.DispatchJobRequest{
		InstanceNo:  instance.InstanceNo,
		JobCode:     job.JobCode,
		HandlerName: job.HandlerName,
		ShardNo:     "",
		TimeoutMs:   int32(job.TimeoutMs),
		Payload:     payload,
		TraceId:     "",
	})
	return err
}

func ptrString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func newInstanceNo() string {
	return fmt.Sprintf("JI-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}
