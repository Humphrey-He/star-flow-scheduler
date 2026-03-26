package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultPoolSize      = 10
	defaultTimeoutMillis = 1000
)

func NewRedis(c Config) (*redis.Client, error) {
	if c.Addr == "" {
		return nil, nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     intOrDefault(c.PoolSize, defaultPoolSize),
		DialTimeout:  msToDuration(c.DialTimeoutMs),
		ReadTimeout:  msToDuration(c.ReadTimeoutMs),
		WriteTimeout: msToDuration(c.WriteTimeoutMs),
	})

	ctx, cancel := context.WithTimeout(context.Background(), msToDuration(c.DialTimeoutMs))
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func MustNewRedis(c Config) *redis.Client {
	client, err := NewRedis(c)
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("redis config is empty")
	}
	return client
}

func msToDuration(ms int64) time.Duration {
	if ms <= 0 {
		ms = defaultTimeoutMillis
	}
	return time.Duration(ms) * time.Millisecond
}

func intOrDefault(val int, def int) int {
	if val <= 0 {
		return def
	}
	return val
}
