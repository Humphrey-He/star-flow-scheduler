package service

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/repo"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
)

type InstanceService struct {
	instances *repo.JobInstanceRepository
}

func NewInstanceService(instances *repo.JobInstanceRepository) *InstanceService {
	return &InstanceService{instances: instances}
}

func (s *InstanceService) Transition(ctx context.Context, instanceNo string, from state.InstanceStatus, to state.InstanceStatus) (bool, error) {
	if err := state.ValidateTransition(from, to); err != nil {
		return false, err
	}

	rows, err := s.instances.UpdateStatusIf(ctx, instanceNo, string(from), string(to))
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (s *InstanceService) ReportResult(ctx context.Context, instanceNo string, status state.InstanceStatus, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	return s.instances.UpdateResult(ctx, instanceNo, string(status), resultSummary, errorCode, errorMessage)
}
