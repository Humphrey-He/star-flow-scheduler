package repo

import (
	"context"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/deadletter"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/types"
)

type DeadLetterFilter struct {
	Status     string
	DeadReason string
	Page       int
	PageSize   int
}

type DeadLetterCreate struct {
	InstanceNo     string
	JobID          int64
	DeadReason     string
	DeadMessage    *string
	RetrySnapshot  *string
	LastExecutorID *int64
	Status         string
}

type DeadLetterRepository struct {
	client *ent.Client
}

func NewDeadLetterRepository(client *ent.Client) *DeadLetterRepository {
	return &DeadLetterRepository{client: client}
}

func (r *DeadLetterRepository) Create(ctx context.Context, req DeadLetterCreate) (*ent.DeadLetter, error) {
	create := r.client.DeadLetter.Create().
		SetInstanceNo(req.InstanceNo).
		SetJobID(req.JobID).
		SetDeadReason(req.DeadReason).
		SetStatus(req.Status)

	if req.DeadMessage != nil {
		create.SetDeadMessage(*req.DeadMessage)
	}
	if req.RetrySnapshot != nil {
		create.SetRetrySnapshot(*req.RetrySnapshot)
	}
	if req.LastExecutorID != nil {
		create.SetLastExecutorID(*req.LastExecutorID)
	}

	return create.Save(ctx)
}

func (r *DeadLetterRepository) GetByID(ctx context.Context, id int64) (*ent.DeadLetter, error) {
	return r.client.DeadLetter.Get(ctx, int(id))
}

func (r *DeadLetterRepository) List(ctx context.Context, filter DeadLetterFilter) ([]*ent.DeadLetter, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 200 {
		filter.PageSize = 20
	}

	query := r.client.DeadLetter.Query()

	if filter.Status != "" {
		query = query.Where(deadletter.StatusEQ(filter.Status))
	}
	if filter.DeadReason != "" {
		query = query.Where(deadletter.DeadReasonEQ(filter.DeadReason))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	items, err := query.Order(ent.Desc(deadletter.FieldID)).Limit(filter.PageSize).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *DeadLetterRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.client.DeadLetter.UpdateOneID(int(id)).SetStatus(status).Exec(ctx)
}

func (r *DeadLetterRepository) BatchListByIDs(ctx context.Context, ids []int64) ([]*ent.DeadLetter, error) {
	if len(ids) == 0 {
		return []*ent.DeadLetter{}, nil
	}
	idInts := make([]int, 0, len(ids))
	for _, id := range ids {
		idInts = append(idInts, int(id))
	}
	return r.client.DeadLetter.Query().Where(deadletter.IDIn(idInts...)).All(ctx)
}

func (r *DeadLetterRepository) MarkRetriedIfOpen(ctx context.Context, id int64) (int, error) {
	return r.client.DeadLetter.Update().
		Where(deadletter.IDEQ(int(id)), deadletter.StatusEQ(string(types.DeadLetterStatusOpen))).
		SetStatus(string(types.DeadLetterStatusRetried)).
		Save(ctx)
}

func (r *DeadLetterRepository) MarkClosedIfOpen(ctx context.Context, id int64) (int, error) {
	return r.client.DeadLetter.Update().
		Where(deadletter.IDEQ(int(id)), deadletter.StatusEQ(string(types.DeadLetterStatusOpen))).
		SetStatus(string(types.DeadLetterStatusClosed)).
		Save(ctx)
}
