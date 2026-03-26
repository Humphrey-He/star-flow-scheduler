package redisx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type ExecutorHeartbeat struct {
	ExecutorCode string `json:"executor_code"`
	CurrentLoad  int64  `json:"current_load"`
	RunningJobs  int64  `json:"running_jobs"`
	UpdatedAtMs  int64  `json:"updated_at_ms"`
}

type HeartbeatCache interface {
	Set(ctx context.Context, executorCode string, value ExecutorHeartbeat, ttl time.Duration) error
	Get(ctx context.Context, executorCode string) (*ExecutorHeartbeat, error)
	Delete(ctx context.Context, executorCode string) error
}

type heartbeatCache struct {
	client *redis.Client
}

func NewHeartbeatCache(client *redis.Client) HeartbeatCache {
	return &heartbeatCache{client: client}
}

func (c *heartbeatCache) Set(ctx context.Context, executorCode string, value ExecutorHeartbeat, ttl time.Duration) error {
	if c.client == nil {
		return ErrNotFound
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, ExecutorHeartbeatKey(executorCode), payload, ttl).Err()
}

func (c *heartbeatCache) Get(ctx context.Context, executorCode string) (*ExecutorHeartbeat, error) {
	if c.client == nil {
		return nil, ErrNotFound
	}
	raw, err := c.client.Get(ctx, ExecutorHeartbeatKey(executorCode)).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var out ExecutorHeartbeat
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *heartbeatCache) Delete(ctx context.Context, executorCode string) error {
	if c.client == nil {
		return ErrNotFound
	}
	return c.client.Del(ctx, ExecutorHeartbeatKey(executorCode)).Err()
}
