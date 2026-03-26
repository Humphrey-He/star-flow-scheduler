package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReadyQueue interface {
	Push(ctx context.Context, instanceNo string) error
	Pop(ctx context.Context, timeout time.Duration) (string, error)
	Len(ctx context.Context) (int64, error)
}

type readyQueue struct {
	client *redis.Client
	key    string
}

func NewReadyQueue(client *redis.Client) ReadyQueue {
	return &readyQueue{
		client: client,
		key:    ReadyQueueKey(),
	}
}

func (q *readyQueue) Push(ctx context.Context, instanceNo string) error {
	if q.client == nil {
		return ErrNotFound
	}
	return q.client.LPush(ctx, q.key, instanceNo).Err()
}

func (q *readyQueue) Pop(ctx context.Context, timeout time.Duration) (string, error) {
	if q.client == nil {
		return "", ErrNotFound
	}
	res, err := q.client.BRPop(ctx, timeout, q.key).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	if len(res) != 2 {
		return "", ErrNotFound
	}
	return res[1], nil
}

func (q *readyQueue) Len(ctx context.Context) (int64, error) {
	if q.client == nil {
		return 0, ErrNotFound
	}
	return q.client.LLen(ctx, q.key).Result()
}
