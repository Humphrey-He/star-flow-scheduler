package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/jobshard"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
)

type ShardCreate struct {
	InstanceID int64
	ShardNo    string
	ShardIndex int
	ShardTotal int
	RouteKey   *string
	ExecutorID *int64
	Status     string
	Payload    *string
}

type ShardRepository struct {
	client *ent.Client
}

func NewShardRepository(client *ent.Client) *ShardRepository {
	return &ShardRepository{client: client}
}

func (r *ShardRepository) BatchCreate(ctx context.Context, shards []ShardCreate) ([]*ent.JobShard, error) {
	if len(shards) == 0 {
		return []*ent.JobShard{}, nil
	}

	builders := make([]*ent.JobShardCreate, 0, len(shards))
	for _, s := range shards {
		create := r.client.JobShard.Create().
			SetInstanceID(s.InstanceID).
			SetShardNo(s.ShardNo).
			SetShardIndex(s.ShardIndex).
			SetShardTotal(s.ShardTotal).
			SetStatus(s.Status)

		if s.RouteKey != nil {
			create.SetRouteKey(*s.RouteKey)
		}
		if s.ExecutorID != nil {
			create.SetExecutorID(*s.ExecutorID)
		}
		if s.Payload != nil {
			create.SetPayload(*s.Payload)
		}
		builders = append(builders, create)
	}

	return r.client.JobShard.CreateBulk(builders...).Save(ctx)
}

func (r *ShardRepository) ListByInstanceNo(ctx context.Context, instanceID int64) ([]*ent.JobShard, error) {
	return r.client.JobShard.Query().
		Where(jobshard.InstanceIDEQ(instanceID)).
		Order(ent.Asc(jobshard.FieldShardIndex)).
		All(ctx)
}

func (r *ShardRepository) UpdateStatusIf(ctx context.Context, shardNo string, fromStatus string, toStatus string) (int, error) {
	return r.client.JobShard.Update().
		Where(jobshard.ShardNoEQ(shardNo), jobshard.StatusEQ(fromStatus)).
		SetStatus(toStatus).
		Save(ctx)
}

func (r *ShardRepository) CountByInstanceNoAndStatus(ctx context.Context, instanceID int64, status string) (int, error) {
	return r.client.JobShard.Query().
		Where(jobshard.InstanceIDEQ(instanceID), jobshard.StatusEQ(status)).
		Count(ctx)
}

func (r *ShardRepository) MarkSuccessIfRunning(ctx context.Context, shardNo string) (int, error) {
	return r.UpdateStatusIf(ctx, shardNo, string(types.ShardStatusRunning), string(types.ShardStatusSuccess))
}

func (r *ShardRepository) MarkFailedIfRunning(ctx context.Context, shardNo string) (int, error) {
	return r.UpdateStatusIf(ctx, shardNo, string(types.ShardStatusRunning), string(types.ShardStatusFailed))
}
