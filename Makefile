# StarFlow Scheduler
# goctl file name style: use underscore instead of hyphen

GOCTL_STYLE ?= go_zero
API_FILE := api/scheduler.api
SWAGGER_DIR := docs/swagger

.PHONY: api rpc-scheduler rpc-executor swagger ent-generate ent-migrate fmt

api:
	goctl api go --api $(API_FILE) --dir apps/scheduler/api --style $(GOCTL_STYLE)

rpc-scheduler:
	goctl rpc protoc proto/executor.proto --proto_path=proto --go_out=proto/pb --go-grpc_out=proto/pb --zrpc_out=apps/scheduler/rpc -m --style $(GOCTL_STYLE)

rpc-executor:
	goctl rpc protoc proto/executor.proto --proto_path=proto --go_out=proto/pb --go-grpc_out=proto/pb --zrpc_out=apps/executor/rpc -m --style $(GOCTL_STYLE)

swagger:
	goctl api plugin --plugin goctl-swagger --api $(API_FILE) --dir $(SWAGGER_DIR) --style $(GOCTL_STYLE)

ent-generate:
	go generate ./pkg/ent

ent-migrate:
	bash pkg/ent/migrate/run_migrations.sh

fmt:
	gofmt -w apps pkg cmd
