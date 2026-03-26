package redisx

import "errors"

var (
	ErrNotFound    = errors.New("redis: not found")
	ErrLockNotHeld = errors.New("redis: lock not held")
)
