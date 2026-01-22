# 任务列表

## 阶段 1：Agent 端

## 任务 1: Agent 项目初始化

**需求**: 需求 1

**说明**: 创建 Agent 项目的基础结构，包括 Go 模块、目录结构和构建配置。

### 子任务

- [ ] 1.1 创建 `agent/` 目录和 `go.mod`（模块名 `github.com/orbit/agent`）
- [ ] 1.2 创建项目目录结构（cmd/agent、internal/config、internal/connection 等）
- [ ] 1.3 创建 Makefile，支持 `build`、`build-linux-amd64`、`build-linux-arm64` 目标
- [ ] 1.4 创建 `cmd/agent/main.go` 入口文件（基础框架）

---

## 任务 2: 配置管理

**需求**: 需求 2

**说明**: 实现命令行参数解析和环境变量读取。

### 子任务

- [ ] 2.1 创建 `internal/config/config.go`，定义 Config 结构体
- [ ] 2.2 实现命令行参数解析（--server、--key、--name、--cpu-threshold、--mem-threshold）
- [ ] 2.3 实现环境变量读取（AGENT_SERVER_HOST、AGENT_API_KEY）
- [ ] 2.4 实现配置验证（必需参数检查）
- [ ] 2.5 在 main.go 中集成配置加载
- [ ] 2.6 基于 `--server` 派生 WebSocket URL（wss://<server>/api/agents/ws）与 Worker HTTP URL（https://<server>，强制 https/wss），TLS 校验默认关闭

---

## 任务 3: 日志和系统信息

**需求**: 需求 10, 需求 4

**说明**: 实现结构化日志和系统信息采集。使用 `github.com/shirou/gopsutil/v3` 采集系统负载，正确处理容器 cgroup 限制。

### 子任务

- [ ] 3.1 创建 `internal/pkg/logger.go`，封装 zap 日志
- [ ] 3.2 支持日志级别配置（debug/info/warn/error）
- [ ] 3.3 创建 `internal/pkg/system.go`，使用 gopsutil 实现 CPU/内存使用率采集
- [ ] 3.4 编写系统信息采集的单元测试

---

## 任务 4: WebSocket 连接管理

**需求**: 需求 3

**说明**: 实现 WebSocket 客户端，包括连接、重连和消息收发。

### 子任务

- [ ] 4.1 创建 `internal/connection/message.go`，定义消息结构
- [ ] 4.2 创建 `internal/connection/client.go`，实现 WebSocketClient 接口
- [ ] 4.3 实现连接建立（携带 X-Agent-Key Header）
- [ ] 4.4 实现消息发送和接收
- [ ] 4.5 创建 `internal/connection/reconnect.go`，实现指数退避重连
- [ ] 4.6 实现 ping/pong 响应
- [ ] 4.7 编写连接管理的单元测试
- [ ] 4.8 在连接建立时记录 last_seen_ip（RemoteAddr / X-Forwarded-For）

---

## 任务 5: 心跳上报

**需求**: 需求 4

**说明**: 实现定时心跳上报，包含系统负载信息。

### 子任务

- [ ] 5.1 创建 `internal/heartbeat/reporter.go`，实现 HeartbeatReporter
- [ ] 5.2 实现每 5 秒发送心跳消息
- [ ] 5.3 心跳包含 CPU、内存、任务数、版本号
- [ ] 5.4 在 main.go 中启动心跳协程
- [ ] 5.5 心跳携带 hostname，并在 Server 侧保存

---

## 任务 6: 任务执行器

**需求**: 需求 5, 需求 6

**说明**: 实现基于 channel 的任务执行器，包含队列、调度和执行逻辑。使用 gopsutil 进行负载检查。

### 子任务

- [ ] 6.1 创建 `internal/task/executor.go`，实现 TaskExecutor
- [ ] 6.2 使用 channel 作为任务队列
- [ ] 6.3 实现 for 循环消费任务
- [ ] 6.4 实现负载检查（使用 gopsutil 获取 CPU/内存，与阈值比较）
- [ ] 6.5 实现任务取消（标记 + 停止容器）
- [ ] 6.6 编写任务执行器的单元测试

---

## 任务 7: Docker 容器管理

**需求**: 需求 7

**说明**: 实现 Worker 容器的启动、监控和清理。

### 子任务

- [ ] 7.1 创建 `internal/docker/runner.go`，实现 DockerRunner 接口
- [ ] 7.2 实现容器启动（Docker SDK，自动清理）
- [ ] 7.3 实现环境变量传递（SERVER_URL、SCAN_ID 等）
- [ ] 7.4 实现目录挂载（/opt/orbit）
- [ ] 7.5 实现容器停止（用于任务取消）
- [ ] 7.6 实现等待容器退出并获取退出码
- [ ] 7.7 编写 Docker SDK 操作的集成测试

---

## 任务 8: 扫描状态管理

**需求**: 需求 8

**说明**: 实现 Agent 调用 HTTP API 更新扫描状态。

### 子任务

- [ ] 8.1 创建 `internal/server/client.go`，实现 HTTP 客户端
- [ ] 8.2 实现 `UpdateScanStatus(scanID, status, errorMessage)` 方法
- [ ] 8.3 在任务执行器中集成状态更新：
  - 收到任务 → `scheduled`
  - 启动 Worker → `running`
  - Worker 退出码 0 → `completed`
  - Worker 退出码非 0 → `failed`
  - 任务取消 → `cancelled`
- [ ] 8.4 处理 HTTP 请求失败（重试 + 日志），并使用 X-Agent-Key 认证

---

## 任务 9: Agent 主循环

**需求**: 需求 3, 需求 5

**说明**: 实现 Agent 主循环，协调各组件。

### 子任务

- [ ] 9.1 在 main.go 中实现主循环
- [ ] 9.2 启动 WebSocket 连接（带重连）
- [ ] 9.3 启动心跳上报协程
- [ ] 9.4 启动任务执行器协程
- [ ] 9.5 实现消息分发（根据消息类型调用对应处理器）
- [ ] 9.6 实现优雅关闭（SIGINT/SIGTERM）
- [ ] 9.7 收到 task_assign 后发送 task_ack（入队即确认）
- [ ] 9.8 Server 端使用 ack_timeout=5s 进行投递重发

---

## 任务 10: Docker 镜像构建

**需求**: 需求 10

**说明**: 创建 Agent 的 Dockerfile 和构建配置。

### 子任务

- [ ] 10.1 创建 `agent/Dockerfile`（基于 alpine，不包含 docker-cli）
- [ ] 10.2 配置多架构构建（amd64/arm64）
- [ ] 10.3 在 Makefile 中添加 `docker-build` 目标
- [ ] 10.4 创建安装脚本 `install-agent.sh`
- [ ] 10.5 编写部署文档（README.md）

---

## 任务 11: 自动更新

**需求**: 需求 12

**说明**: 实现 Agent 自动更新功能。

### 子任务

- [ ] 11.1 创建 `internal/updater/updater.go`，实现 Updater 组件
- [ ] 11.2 实现镜像拉取（docker pull）
- [ ] 11.3 实现启动新版本容器（使用相同配置）
- [ ] 11.4 实现退出当前进程（让新容器接管）
- [ ] 11.5 在消息处理中集成 `update_required` 处理
- [ ] 11.6 编写更新流程的集成测试

---

## 阶段 2：Server 端（后续实现）

## 任务 12: Server 端 WebSocket Hub

**需求**: 需求 9

**说明**: 在 Server 端实现 WebSocket 连接管理中心。

### 子任务

- [ ] 12.1 创建 `server/internal/websocket/message.go`，定义消息结构
- [ ] 12.2 创建 `server/internal/websocket/client.go`，封装单个 Agent 连接
- [ ] 12.3 创建 `server/internal/websocket/hub.go`，实现连接管理
- [ ] 12.4 实现 Agent 注册/注销
- [ ] 12.5 实现向指定 Agent 发送消息
- [ ] 12.6 实现心跳超时检测（15 秒）

---

## 任务 13: Server 端 WebSocket 端点

**需求**: 需求 9

**说明**: 实现 WebSocket API 端点和认证。

### 子任务

- [ ] 13.1 创建 `server/internal/handler/agent_ws.go`
- [ ] 13.2 实现 `/api/agents/ws` WebSocket 端点
- [ ] 13.3 实现 API Key 认证（从 Header 读取 X-Agent-Key，查数据库验证）
- [ ] 13.4 认证成功后更新 agent 记录（status=online, connected_at, hostname 等）
- [ ] 13.5 实现消息处理（心跳、任务状态等）
- [ ] 13.6 在路由中注册 WebSocket 端点

---

## 任务 14: Agent 数据模型和 API

**需求**: 需求 9

**说明**: 实现 Agent 数据模型和管理 API。

### 子任务

- [ ] 14.1 创建数据库迁移文件（agent 表）
- [ ] 14.2 创建 `server/internal/model/agent.go` 模型
- [ ] 14.3 实现 Agent CRUD API（创建、列表、删除）
- [ ] 14.4 创建 Agent 时生成 API Key
- [ ] 14.5 返回部署命令（包含 Key）
- [ ] 14.6 创建 `GET /api/agents/install.sh` 安装脚本接口
- [ ] 14.7 创建 Agent 时返回 installCommand 与 installScriptUrl
- [ ] 14.8 API 响应包含 last_seen_ip 与 hostname（供前端展示）

---

## 任务 15: Server 端任务分发

**需求**: 需求 9

**说明**: 实现基于 WebSocket 的任务分发服务。

### 子任务

- [ ] 15.1 创建 `server/internal/service/task_dispatcher.go`
- [ ] 15.2 实现选择最优 Agent（基于心跳负载数据）
- [ ] 15.3 实现任务推送（通过 WebSocket）
- [ ] 15.4 实现任务状态回调处理
- [ ] 15.5 更新扫描 API，使用新的任务分发服务
- [ ] 15.6 task_assign payload 包含 targetType、workflowName、workerImage；若未提供 workerImage，Agent 使用默认镜像

---

## 任务 16: Server 端版本检测

**需求**: 需求 12

**说明**: 在 Server 端实现 Agent 版本检测和更新触发。

### 子任务

- [ ] 16.1 在心跳处理中比较 Agent 版本和 Server 版本
- [ ] 16.2 版本不匹配时发送 `update_required` 消息
- [ ] 16.3 使用 Redis 锁防止重复触发（60 秒内只触发一次）
- [ ] 16.4 记录更新日志

---

## 任务 17: 集成测试

**需求**: 全部

**说明**: 编写端到端集成测试。

### 子任务

- [ ] 17.1 编写 Agent 连接和认证测试
- [ ] 17.2 编写心跳上报测试
- [ ] 17.3 编写任务分配和执行测试
- [ ] 17.4 编写断线重连测试
- [ ] 17.5 编写多 Agent 负载均衡测试
- [ ] 17.6 编写自动更新测试
