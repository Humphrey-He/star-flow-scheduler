package service

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
)

type InstanceService struct {
	instances *repo.JobInstanceRepository
}

func NewInstanceService(instances *repo.JobInstanceRepository) *InstanceService {
	return &InstanceService{instances: instances}
}

func (s *InstanceService) List(ctx context.Context, filter repo.JobInstanceFilter) ([]*ent.JobInstance, int, error) {
	return s.instances.List(ctx, filter)
}

func (s *InstanceService) GetByInstanceNo(ctx context.Context, instanceNo string) (*ent.JobInstance, error) {
	return s.instances.GetByInstanceNo(ctx, instanceNo)
}
