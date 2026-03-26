package redisx

const (
	delayQueueKeyPrefix      = "sched:delay:zset"
	readyQueueKeyPrefix      = "sched:ready:list"
	lockKeyPrefix            = "sched:lock:"
	executorHeartbeatKeyPref = "sched:executor:hb:"
	idempotentKeyPrefix      = "sched:idem:"
)

func DelayQueueKey() string {
	return delayQueueKeyPrefix
}

func ReadyQueueKey() string {
	return readyQueueKeyPrefix
}

func LockKey(name string) string {
	return lockKeyPrefix + name
}

func ExecutorHeartbeatKey(code string) string {
	return executorHeartbeatKeyPref + code
}

func IdempotentKey(key string) string {
	return idempotentKeyPrefix + key
}
