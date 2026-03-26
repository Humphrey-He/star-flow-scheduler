package instance

import (
	"context"
	"errors"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/rpc/internal/state"
)

type Service struct {
	instances instanceRepo
}

type instanceRepo interface {
	GetStatusByInstanceNo(ctx context.Context, instanceNo string) (string, error)
	UpdateStatusIf(ctx context.Context, instanceNo string, fromStatus string, toStatus string) (int, error)
	UpdateResultIfStatus(ctx context.Context, instanceNo string, fromStatus string, toStatus string, resultSummary *string, errorCode *string, errorMessage *string) (int, error)
}

func NewService(instances instanceRepo) *Service {
	return &Service{instances: instances}
}

func (s *Service) Transition(ctx context.Context, instanceNo string, from state.InstanceStatus, to state.InstanceStatus) (bool, error) {
	if err := state.ValidateTransition(from, to); err != nil {
		return false, err
	}

	rows, err := s.instances.UpdateStatusIf(ctx, instanceNo, string(from), string(to))
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (s *Service) ReportResult(ctx context.Context, instanceNo string, status state.InstanceStatus, resultSummary *string, errorCode *string, errorMessage *string) (int, error) {
	currentStatus, err := s.instances.GetStatusByInstanceNo(ctx, instanceNo)
	if err != nil {
		return 0, err
	}

	from := state.InstanceStatus(currentStatus)
	to := status

	if from == to {
		return 0, nil
	}

	if err := state.ValidateTransition(from, to); err != nil {
		return 0, err
	}

	rows, err := s.instances.UpdateResultIfStatus(ctx, instanceNo, string(from), string(to), resultSummary, errorCode, errorMessage)
	if err != nil {
		return 0, err
	}

	if rows == 0 {
		latestStatus, err := s.instances.GetStatusByInstanceNo(ctx, instanceNo)
		if err == nil && state.InstanceStatus(latestStatus) == to {
			return 0, nil
		}
		return 0, errors.New("status conflict")
	}

	return rows, nil
}
