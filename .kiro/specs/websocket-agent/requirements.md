# 需求文档

## 简介

本功能实现基于 WebSocket 的轻量级 Agent，用于替代当前的 SSH 任务分发方式。Agent 作为常驻服务运行在远程 VPS 上，主动连接 Server 建立长连接，接收任务指令并启动临时 Worker 容器执行扫描。

## 核心优势

- **无需公网 IP**：Agent 主动连接 Server，支持 NAT 穿透
- **无需 SSH**：不暴露 SSH 端口，更安全
- **一键安装**：单个二进制文件 + systemd 服务
- **轻量级**：Go 编译，~10MB 内存占用

## 术语表

- **Agent**: 常驻在远程 VPS 上的轻量服务，负责接收任务、启动 Worker、上报状态
- **Server**: Go 后端 API 服务（`server/`），负责任务调度和 WebSocket 连接管理
- **Worker**: 临时容器，执行具体的扫描任务，完成后退出
- **Task**: 一个扫描任务，包含 task_id、scan_id、target、config 等信息
- **Heartbeat**: Agent 定期发送的心跳消息，包含系统负载信息
- **Agent API Key**: 绑定到单个 Agent 记录的认证密钥，用于 WebSocket 连接和 /api/agent/** HTTP 调用

## 需求

### 需求 1: Agent 项目初始化

**用户故事:** 作为开发者，我希望创建一个独立的 Go 模块用于 Agent，以便它可以编译为单个二进制文件分发。

#### 验收标准

1. Agent 应当作为独立的 Go 模块创建在 `agent/` 目录下
2. Agent 应当有自己的 `go.mod`，模块名为 `github.com/orbit/agent`
3. Agent 应当编译为单个静态链接的二进制文件（~10MB）
4. Agent 应当支持 Linux amd64 和 arm64 架构

### 需求 2: 配置管理

**用户故事:** 作为运维人员，我希望通过命令行参数配置 Agent，以便快速部署。

#### 验收标准

1. Agent 应当支持以下命令行参数：
   - `--server`: Server 基础地址（必需，仅 IP/域名/端口，不包含协议）
   - `--key`: API Key 用于认证（必需）
   - `--name`: Agent 名称（可选，默认为主机名）
2. 当缺少必需参数时，Agent 应当显示帮助信息并退出
3. 调度参数（maxTasks、cpuThreshold、memThreshold）由 Server 动态下发，不通过 CLI 配置
3. Agent 应当支持从环境变量读取配置（`AGENT_SERVER_HOST`、`AGENT_API_KEY`），其中 `AGENT_SERVER_HOST` 为基础地址（不含协议）
4. Agent 根据 `--server` 自动派生（强制 https/wss）：
   - WebSocket URL：`wss://<server>/api/agents/ws`
   - Worker HTTP URL：`https://<server>`（用于注入到 Worker 的 `SERVER_URL`）
5. TLS 校验默认关闭（跳过证书校验）

### 需求 3: WebSocket 连接管理

**用户故事:** 作为 Agent，我希望与 Server 建立稳定的 WebSocket 长连接，以便接收实时控制通知（配置更新/取消/更新）。

#### 验收标准

1. Agent 启动时应当主动连接 Server 的 WebSocket 端点
2. 连接时应当在 Header 中携带 API Key 进行认证
3. 当连接失败时，Agent 应当使用指数退避策略重试（1s, 2s, 4s, 8s, 最大 60s）
4. 当连接断开时，Agent 应当自动重连
5. Agent 应当响应 Server 的 ping 消息以保持连接活跃

### 需求 4: 心跳与负载上报

**用户故事:** 作为 Server，我希望实时了解 Agent 的负载情况，以便进行任务调度。

#### 验收标准

1. Agent 应当每 5 秒发送一次心跳消息
2. 心跳消息应当包含：
   - `cpu`: CPU 使用率（%）
   - `mem`: 内存使用率（%）
   - `disk`: 磁盘使用率（%）
   - `tasks`: 当前运行的任务数
   - `version`: Agent 版本号
   - `hostname`: 主机名
   - `uptime`: Agent 运行时长（秒）
3. 当 Server 超过 15 秒未收到心跳时，应当将 Agent 标记为离线

**说明（IP/Host 记录）**：
- Server 记录 `last_seen_ip`（来自 WebSocket 连接的 RemoteAddr 或代理头）
- Agent 在心跳中上报 `hostname`
- IP 仅用于诊断与展示，不作为身份校验

### 需求 5: 任务拉取与执行

**用户故事:** 作为 Agent，我希望通过 HTTP 主动拉取任务并执行，而不是被动接收推送。

#### 验收标准

1. Agent 应当通过 `POST /api/agent/tasks/pull` 拉取任务（不使用 WebSocket 推送）
2. 任务拉取响应应当包含：
   - task_id: 任务 ID（scan_task.id）
   - scan_id: 扫描 ID（任务所属 scan）
   - target_id: 目标 ID
   - target_name: 目标名称
   - target_type: 目标类型（domain/ip/cidr/url）
   - workflow_name: 工作流名称（如 subdomain_discovery）
   - workspace_dir: 工作目录
   - config: YAML 格式的扫描配置
   - worker_image: Worker 镜像名称和版本（可选，如 `yyhuni/orbit-worker:v1.0.19`）
3. Agent 不维护本地等待队列，仅跟踪运行中任务
4. Agent 应当跟踪每个任务的状态（running/completed/failed）

**说明**：`taskId` 为 scan_task.id；`scanId` 表示所属 scan（单 workflow 时为 1:1，但仍区分）。

**说明（Worker 版本策略）**：
- Server 负责决定 Worker 镜像版本（通常与 Server 版本同步）
- 任务中携带 `worker_image` 用于显式指定版本
- 若任务未指定 `worker_image`，Agent 使用默认 Worker 镜像

### 需求 6: 任务调度（Agent Pull + PostgreSQL 队列）

**用户故事:** 作为 Agent，我希望在空闲时向 Server 主动拉取最高优先级的任务，避免饥饿并实现负载均衡。

#### 当前实现

**单 Workflow 模式**：一个 Scan 对应一个 scan_task，执行 `subdomain_discovery` workflow。

- 创建 Scan 时，同时创建一个 `status='ready'` 的 scan_task
- workflow 内部的多阶段（recon/bruteforce/permutation/resolve）由 Worker 内部编排
- `depends_on` 字段预留为空数组，未来支持多 workflow 串联

#### 验收标准

1. Agent 应当持续监控本机 CPU、内存和磁盘使用率
2. 当同时满足以下条件时，调用 Server API 拉取任务（pull）：
   - 当前运行任务数 < max_tasks（默认 5）
   - CPU 使用率 < cpu_threshold（默认 85%）
   - 内存使用率 < mem_threshold（默认 85%）
   - 磁盘使用率 < disk_threshold（默认 90%）
3. Server 从 PostgreSQL 中按优先级选择一个 `ready` 任务返回，并进行行级锁定，防止并发重复分配
4. 优先级计算：`priority = stage * 100 + wait_seconds`（当前所有任务 stage=0）
5. 当任一条件不满足时，Agent 等待后再检查
6. 调度参数由 Server 动态下发（通过 `config_update` 消息）
7. 磁盘使用率通过 gopsutil 的 `disk.Usage("/")` 获取根分区使用率

#### Server 端 SQL（行级锁）

```sql
-- 原子地取出一个最高优先级任务并锁住
WITH c AS (
  SELECT id
  FROM scan_task
  WHERE status = 'ready'
  ORDER BY priority DESC, id ASC
  LIMIT 1
  FOR UPDATE SKIP LOCKED
)
UPDATE scan_task t
SET status = 'dispatched', dispatched_at = NOW(), agent_id = $1
FROM c
WHERE t.id = c.id
RETURNING t.*;
```

#### API
- `POST /api/agent/tasks/pull`（Header: `X-Agent-Key`）
  - Request: `{}`
  - Response: `{taskId, scanId, stage, workerImage, config, ...}` 或 `204 No Content`（无任务）

#### 未来扩展（预留）
- 支持多 workflow 串联（subdomain → port_scan → vuln_scan）
- Worker 端 `templates.yaml` 新增 `depends_on` 和 `output_type` 字段
- Server 端解析 workflow metadata，自动构建 DAG
- 任务完成后触发下游任务 `pending` → `ready`

### 需求 7: Worker 容器管理

**用户故事:** 作为 Agent，我希望管理 Worker 容器的生命周期，以便正确执行和清理任务。

#### 验收标准

1. Agent 应当使用 Docker API（Go SDK）启动 Worker 容器
2. Worker 容器应当使用自动清理策略（等价于 `--rm`），完成后删除
3. Agent 应当使用任务中指定的 Worker 镜像版本
4. Agent 应当使用 “缺失时拉取” 策略（本地有镜像则直接使用，没有才拉取）
5. Agent 应当将任务参数通过环境变量传递给 Worker
6. Agent 应当挂载必要的目录（/opt/orbit）
7. 当任务被取消时，Agent 应当停止对应的 Worker 容器
8. Agent 应当捕获 Worker 的退出码判断任务成功/失败

### 需求 8: 任务状态管理

**用户故事:** 作为 Server，我希望 Agent 负责管理任务状态（scan_task），以便统一状态更新入口。

#### 验收标准

1. Server 在拉取分配任务时将 `scan_task.status` 置为 `dispatched`，并同步将 `scan.status` 置为 `scheduled`
2. Agent 启动 Worker 成功后，应当调用 HTTP API 将 `scan_task.status` 更新为 `running`，并同步将 `scan.status` 置为 `running`
3. Agent 监控 Worker 退出码：
   - 退出码为 0：调用 HTTP API 更新为 `completed`（同步 `scan.status=completed`，并写入 `scan.stopped_at`）
   - 退出码非 0：调用 HTTP API 更新为 `failed`（同步 `scan.status=failed`，写入 `scan.stopped_at`，并附带错误信息）
4. 任务被取消时，Agent 应当调用 HTTP API 更新为 `cancelled`（同步 `scan.status=cancelled`，并写入 `scan.stopped_at`）
5. Worker 不再负责更新状态，只负责执行扫描和保存结果
6. 错误信息获取：Agent 使用 Docker SDK 读取 Worker 容器的最后 100 行日志作为 `errorMessage`
7. 日志截断：如果日志超过 4KB，截断并添加 `[truncated]\n` 前缀

**HTTP API 端点**：
```
PATCH /api/agent/tasks/:id/status
Header: X-Agent-Key: <agent_api_key>
Body: {"status": "running|completed|failed|cancelled", "errorMessage": "..."}
```

**说明**：
- /api/agent/** 使用 X-Agent-Key
- /api/worker/** 仍使用 X-Worker-Token（Worker 容器内调用）

### 需求 9: Server 端 WebSocket 支持

**用户故事:** 作为 Server，我希望管理多个 Agent 的 WebSocket 连接，以便接收心跳并下发实时控制通知（取消/配置更新/更新）。

#### 验收标准

1. Server 应当暴露 `/api/agents/ws` WebSocket 端点
2. Server 应当验证连接时的 API Key
3. Server 应当维护在线 Agent 列表
4. Server 应当能向指定 Agent 推送控制消息（task_cancel/config_update/update_required，以及可选的 task_available）
5. Server 应当处理 Agent 断开连接的情况
6. Agent 连接成功后，Server 应推送 `config_update` 消息（包含 maxTasks、cpuThreshold、memThreshold）
7. 管理员通过 API 修改配置时，Server 应推送 `config_update` 给在线 Agent

### 需求 10: Docker 镜像构建

**用户故事:** 作为运维人员，我希望通过一条 docker run 命令部署 Agent，以便快速部署到多台机器。

#### 验收标准

1. Agent 应当提供 Dockerfile，构建为 Docker 镜像
2. 镜像应当支持 linux/amd64 和 linux/arm64 架构
3. 镜像不需要包含 Docker CLI（使用 Docker SDK）
4. 部署命令应当挂载 `/var/run/docker.sock` 和 `/opt/orbit`
5. 使用 `--restart=always` 实现自动重启
6. Agent 容器使用 `--oom-score-adj=-500` 降低被 OOM 杀死的优先级
7. Worker 容器使用 `oom-score-adj=500` 提高被 OOM 杀死的优先级（保护 Agent）

### 需求 10.1: 安装流程与脚本

**用户故事:** 作为管理员，我希望 Server 能生成安装命令和脚本，方便一键部署 Agent。

#### 验收标准

1. Server 创建 Agent 时生成 API Key，并绑定到该 Agent
2. 创建 Agent 的响应应返回：
   - agentId
   - apiKey
   - installCommand（已拼好 --server 和 --key）
   - installScriptUrl（例如 `/api/agents/install.sh?key=...`）
3. Server 应提供安装脚本下载接口（`GET /api/agents/install.sh`）
4. 安装脚本应：
   - 检查 Docker 是否安装
   - 创建 `/opt/orbit` 目录
   - 使用 `docker run --restart=always ... --server <base> --key <apiKey>` 启动 Agent
   - 若已有旧容器，先 stop/remove 再启动
5. Agent API Key 支持重置/禁用（可选）
6. 仅管理员可创建/删除 Agent（权限控制）

### 需求 11: 错误处理与日志

**用户故事:** 作为运维人员，我希望有清晰的日志输出，以便排查问题。

#### 验收标准

1. Agent 应当记录所有重要事件（连接、断开、任务开始/完成）
2. Agent 应当使用结构化日志（JSON 格式）
3. Agent 应当支持日志级别配置（debug/info/warn/error）
4. 当发生错误时，Agent 应当记录详细的错误信息和上下文

### 需求 12: 自动更新

**用户故事:** 作为运维人员，我希望 Agent 能自动更新到最新版本，无需手动干预。

#### 验收标准

1. Agent 心跳时应当上报当前版本号
2. Server 检测到版本不匹配时，应当发送 `update_required` 消息
3. Agent 收到更新指令后，应当：
   - 拉取新版本镜像（docker pull）
   - 启动新版本容器（使用相同配置）
   - 退出当前容器（让新容器接管）
4. 更新过程应当记录详细日志
5. 更新失败时，Agent 应当继续运行并上报错误
