package config

import (
	"github.com/Humphrey-He/star-flow-scheduler/pkg/redisx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	PostgresDSN string `json:",optional"`
	Redis       redisx.Config `json:",optional"`
}
