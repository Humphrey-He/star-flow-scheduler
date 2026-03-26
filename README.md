# StarFlow Scheduler（分布式任务调度系统）

## 项目简介
StarFlow Scheduler 是一个基于 Go 的分布式任务调度系统，采用“调度中心 + 执行器集群”的架构形态，提供北向 HTTP/JSON 接口与南向 gRPC 内部通信，支持定时/延时/一次性/DAG 任务、分片、幂等与重试等核心能力。

本仓库使用 go-zero 作为服务框架，PostgreSQL 作为主数据存储，ent 作为 ORM，确保类型安全与可演进 schema。

## 技术栈
- Go 1.20+
- go-zero（API + RPC）
- gRPC
- PostgreSQL
- ent ORM

## 目录结构
- api/                     — goctl API 定义
- apps/scheduler/api/       — 调度中心 HTTP API 服务（go-zero）
- apps/scheduler/rpc/       — 调度中心 RPC 服务（go-zero）
- apps/executor/rpc/        — 执行器 RPC 服务（go-zero）
- pkg/ent/                  — ent 生成代码与 schema
- pkg/repo/                 — Repository 层（业务查询与封装）
- pkg/db/                   — PGX 连接、事务封装
- proto/                    — proto 定义与 pb 代码
- scripts/                  — 生成/迁移脚本

## 开发规范（摘要）
- 所有业务代码通过 pkg/repo 访问数据库，禁止在 handler/logic 里直接使用 ent.Client。
- schema 只在 pkg/ent/schema 下维护。
- 迁移脚本纳入版本控制，变更须包含 migration。

## 快速开始

### 1. 生成代码
```
make api
make rpc-scheduler
make rpc-executor
```

### 2. ent 生成与迁移
```
make ent-generate
DATABASE_URL="postgres://postgres:password@127.0.0.1:5432/starflow?sslmode=disable" make ent-migrate
```

### 3. 启动服务
```
# 调度中心 HTTP API
go run ./apps/scheduler/api -f apps/scheduler/api/etc/scheduler-api.yaml

# 调度中心 RPC
go run ./apps/scheduler/rpc -f apps/scheduler/rpc/etc/executor.yaml

# 执行器 RPC
go run ./apps/executor/rpc -f apps/executor/rpc/etc/executor.yaml
```

## goctl 文件命名规则
本项目统一使用下划线风格（a_b），通过 Makefile 传入 `--style go` 并提供 `goctl.yaml` 配置。

## Swagger
如果已安装 goctl-swagger 插件，可执行：
```
make swagger
```

## 注意事项
- `go generate ./pkg/ent` 可能需要下载 ent 依赖。若本地网络受限，请先配置代理或手动准备依赖。
- CI 需要执行 ent 生成与迁移检查，避免 schema 与生成代码不一致。
