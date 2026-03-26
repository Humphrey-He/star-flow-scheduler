package config

import (
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PostgresDSN string `json:",optional"`
	Redis       redisx.Config `json:",optional"`
	Scanner     ScannerConf   `json:",optional"`
	Dispatcher  DispatcherConf `json:",optional"`
	Registry    RegistryConf  `json:",optional"`
}

type ScannerConf struct {
	TickIntervalMs int64 `json:",optional"`
	BatchSize      int64 `json:",optional"`
	LockTTLms      int64 `json:",optional"`
	RequeueDelayMs int64 `json:",optional"`
}

type DispatcherConf struct {
	PopTimeoutMs int64 `json:",optional"`
	IdleSleepMs  int64 `json:",optional"`
	RequeueMs    int64 `json:",optional"`
}

type RegistryConf struct {
	HeartbeatCacheTtlMs int64 `json:",optional"`
}
