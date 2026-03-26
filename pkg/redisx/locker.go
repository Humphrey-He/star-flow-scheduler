package redisx

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Locker interface {
	TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
	Renew(ctx context.Context, key string, ttl time.Duration) error
}

type locker struct {
	client *redis.Client
	mu     sync.Mutex
	tokens map[string]string
}

func NewLocker(client *redis.Client) Locker {
	return &locker{
		client: client,
		tokens: make(map[string]string),
	}
}

func (l *locker) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if l.client == nil {
		return false, ErrNotFound
	}
	token := uuid.NewString()
	ok, err := l.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil || !ok {
		return ok, err
	}
	l.mu.Lock()
	l.tokens[key] = token
	l.mu.Unlock()
	return true, nil
}

func (l *locker) Unlock(ctx context.Context, key string) error {
	if l.client == nil {
		return ErrNotFound
	}
	token, ok := l.getToken(key)
	if !ok {
		return ErrLockNotHeld
	}
	_, err := lockReleaseScript.Run(ctx, l.client, []string{key}, token).Result()
	if err != nil {
		return err
	}
	l.deleteToken(key)
	return nil
}

func (l *locker) Renew(ctx context.Context, key string, ttl time.Duration) error {
	if l.client == nil {
		return ErrNotFound
	}
	token, ok := l.getToken(key)
	if !ok {
		return ErrLockNotHeld
	}
	_, err := lockRenewScript.Run(ctx, l.client, []string{key}, token, int64(ttl/time.Millisecond)).Result()
	return err
}

func (l *locker) getToken(key string) (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	token, ok := l.tokens[key]
	return token, ok
}

func (l *locker) deleteToken(key string) {
	l.mu.Lock()
	delete(l.tokens, key)
	l.mu.Unlock()
}

var lockReleaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`)

var lockRenewScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
	return 0
end
`)
