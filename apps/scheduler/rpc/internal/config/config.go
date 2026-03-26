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
}

type ScannerConf struct {
	TickIntervalMs int64 `json:",optional"`
	BatchSize      int64 `json:",optional"`
	LockTTLms      int64 `json:",optional"`
	RequeueDelayMs int64 `json:",optional"`
}
