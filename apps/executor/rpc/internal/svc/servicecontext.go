package svc

import "github.com/Humphrey-He/star-flow-scheduler/apps/executor/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
