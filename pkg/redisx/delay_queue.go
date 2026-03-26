package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type DelayQueue interface {
	Add(ctx context.Context, instanceNo string, executeAt time.Time) error
	PopDue(ctx context.Context, now time.Time, limit int64) ([]string, error)
	Remove(ctx context.Context, instanceNo string) error
}

type delayQueue struct {
	client *redis.Client
	key    string
}

func NewDelayQueue(client *redis.Client) DelayQueue {
	return &delayQueue{
		client: client,
		key:    DelayQueueKey(),
	}
}

func (q *delayQueue) Add(ctx context.Context, instanceNo string, executeAt time.Time) error {
	if q.client == nil {
		return ErrNotFound
	}
	score := float64(executeAt.UnixMilli())
	return q.client.ZAdd(ctx, q.key, redis.Z{
		Score:  score,
		Member: instanceNo,
	}).Err()
}

func (q *delayQueue) PopDue(ctx context.Context, now time.Time, limit int64) ([]string, error) {
	if q.client == nil {
		return nil, ErrNotFound
	}
	if limit <= 0 {
		return []string{}, nil
	}
	res, err := delayQueuePopScript.Run(ctx, q.client, []string{q.key}, strconv.FormatInt(now.UnixMilli(), 10), strconv.FormatInt(limit, 10)).Result()
	if err == redis.Nil {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	items, ok := res.([]interface{})
	if !ok {
		return []string{}, nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out, nil
}

func (q *delayQueue) Remove(ctx context.Context, instanceNo string) error {
	if q.client == nil {
		return ErrNotFound
	}
	return q.client.ZRem(ctx, q.key, instanceNo).Err()
}

var delayQueuePopScript = redis.NewScript(`
local items = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
if #items > 0 then
	redis.call('ZREM', KEYS[1], unpack(items))
end
return items
`)
