# StarFlow Scheduler（分布式任务调度系统）

## 项目简介
StarFlow Scheduler 是一个基于 Go 的分布式任务调度系统，采用“调度中心 + 执行器集群”的架构形态，提供北向 HTTP/JSON 接口与南向 gRPC 内部通信，支持定时/延时/一次性/DAG 任务、分片、幂等与重试等核心能力。

本仓库使用 go-zero 作为服务框架，PostgreSQL 作为主数据存储，ent 作为 ORM，Redis 作为高性能调度辅助层，确保类型安全、可演进 schema 与调度高吞吐。

## 技术栈
- Go 1.20+
- go-zero（API + RPC）
- gRPC
- PostgreSQL
- ent ORM
- Redis

## 目录结构
- api/                     — goctl API 定义
- apps/scheduler/api/       — 调度中心 HTTP API 服务（go-zero）
- apps/scheduler/rpc/       — 调度中心 RPC 服务（go-zero）
- apps/executor/rpc/        — 执行器 RPC 服务（go-zero）
- pkg/ent/                  — ent 生成代码与 schema
- pkg/repo/                 — Repository 层（业务查询与封装）
- pkg/db/                   — PGX 连接、事务封装
- pkg/redisx/               — Redis client、队列、锁、心跳缓存
- pkg/metricsx/             — 统一指标埋点门面（expvar）
- proto/                    — proto 定义与 pb 代码
- scripts/                  — 生成/迁移脚本

## 架构与核心流程

### 架构概览
- **Scheduler** 负责实例创建、调度决策、派发、回调推进、失败闭环（重试/死信）、Workflow 编排推进。
- **Executor** 负责接收派发任务、执行 handler、上报结果、心跳/负载上报与优雅下线。
- **PostgreSQL** 是最终真相源（实例状态、死信、执行日志、Workflow 运行态）。
- **Redis** 作为加速层（Delay Queue / Ready Queue / 分布式锁 / 心跳缓存）。

### 任务调度主链路
1. **CreateInstance**：调度中心创建任务实例（PG 落库），必要时写入 Redis Delay Queue。
2. **Delay Scanner**：到期任务从 Delay Queue 拉取，进入 Ready Queue。
3. **Ready Dispatcher**：从 Ready Queue 取实例，路由选择执行器并 gRPC 派发。
4. **Executor Runtime**：入队 -> Worker 执行 -> Reporter 异步上报结果。
5. **ReportResult**：Scheduler 回调推进实例状态；若关联 Workflow，则驱动下游节点推进。

### Workflow 编排链路（DAG）
1. Workflow 定义校验（唯一性、依赖合法性、环检测）。
2. 创建 Workflow Instance 并初始化 Node Runtime。
3. Root 节点触发 JobInstance。
4. Job 完成回调驱动 Resolver 判断下游是否触发。
5. fail_strategy 决定失败后是否继续推进（stop / continue）。

## 可观测性
### 指标（统一前缀）
- Scheduler：`scheduler_*`
- Executor：`executor_*`
- Workflow：`workflow_*`（预留/扩展）

目前已覆盖：
- Scheduler Dispatch/Scanner/ReportResult 核心计数与耗时
- Executor 执行、上报核心计数与耗时
- Retry / Dead Letter 计数
- 在线执行器数量统计

### 日志字段
核心链路逐步收敛结构化字段：
- `instance_no` / `executor_code` / `job_code` / `handler_name`
- `workflow_instance_no`
- `trace_id` / `request_id`（逐步覆盖）

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
go run ./apps/scheduler/api -f apps/scheduler/api/etc/scheduler_api.yaml

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
