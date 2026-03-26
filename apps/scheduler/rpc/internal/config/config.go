package config

import (
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PostgresDSN string `json:",optional"`
	Redis       redisx.Config `json:",optional"`
}
