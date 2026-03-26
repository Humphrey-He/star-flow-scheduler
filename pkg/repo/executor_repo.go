package repo

import (
	"context"
	"time"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/executor"
)

type ExecutorUpsert struct {
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
	Metadata      map[string]interface{}
}

type ExecutorRepository struct {
	client *ent.Client
}

func NewExecutorRepository(client *ent.Client) *ExecutorRepository {
	return &ExecutorRepository{client: client}
}

func (r *ExecutorRepository) Upsert(ctx context.Context, req ExecutorUpsert) (*ent.Executor, error) {
	existing, err := r.client.Executor.Query().Where(executor.ExecutorCodeEQ(req.ExecutorCode)).Only(ctx)
	if err == nil {
		update := r.client.Executor.UpdateOne(existing).
			SetHost(req.Host).
			SetIP(req.IP).
			SetGrpcAddr(req.GrpcAddr).
			SetCapacity(req.Capacity).
			SetCurrentLoad(req.CurrentLoad).
			SetStatus(req.Status).
			SetLastHeartbeatAt(req.LastHeartbeat)

		if req.HttpAddr != nil {
			update.SetHTTPAddr(*req.HttpAddr)
		}
		if req.Tags != nil {
			update.SetTags(*req.Tags)
		}
		if req.Version != nil {
			update.SetVersion(*req.Version)
		}
		if req.Metadata != nil {
			update.SetMetadata(req.Metadata)
		}

		return update.Save(ctx)
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}

	create := r.client.Executor.Create().
		SetExecutorCode(req.ExecutorCode).
		SetHost(req.Host).
		SetIP(req.IP).
		SetGrpcAddr(req.GrpcAddr).
		SetCapacity(req.Capacity).
		SetCurrentLoad(req.CurrentLoad).
		SetStatus(req.Status).
		SetLastHeartbeatAt(req.LastHeartbeat)

	if req.HttpAddr != nil {
		create.SetHTTPAddr(*req.HttpAddr)
	}
	if req.Tags != nil {
		create.SetTags(*req.Tags)
	}
	if req.Version != nil {
		create.SetVersion(*req.Version)
	}
	if req.Metadata != nil {
		create.SetMetadata(req.Metadata)
	}

	return create.Save(ctx)
}

func (r *ExecutorRepository) UpdateHeartbeat(ctx context.Context, executorCode string, currentLoad int) error {
	return r.client.Executor.Update().
		Where(executor.ExecutorCodeEQ(executorCode)).
		SetCurrentLoad(currentLoad).
		SetLastHeartbeatAt(time.Now()).
		SetStatus("online").
		Exec(ctx)
}
