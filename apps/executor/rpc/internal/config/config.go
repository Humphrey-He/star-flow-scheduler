package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	SchedulerRpc zrpc.RpcClientConf
	Executor     ExecutorConf
	Runtime      RuntimeConf
	Registry     RegistryConf
	Reporter     ReporterConf
}

type ExecutorConf struct {
	ExecutorCode string
	Host         string
	IP           string
	GrpcAddr     string
	HttpAddr     string
	Tags         []string
	Capacity     int32
	Version      string
}

type RuntimeConf struct {
	WorkerCount        int
	QueueSize          int
	DefaultTimeoutMs   int64
	ShutdownTimeoutSec int64
}

type RegistryConf struct {
	HeartbeatIntervalSec int64
	RegisterRetryTimes   int
}

type ReporterConf struct {
	QueueSize       int
	RetryTimes      int
	RetryIntervalMs int64
}
