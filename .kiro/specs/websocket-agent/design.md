# 设计文档

## 概述

本文档描述 WebSocket Agent 的技术设计。Agent 作为轻量级常驻服务运行在远程 VPS 上，通过 WebSocket 主动连接 Server，接收任务指令并启动临时 Worker 容器执行扫描。

## 架构

```

**地址派生规则（强制 https/wss）**：
- WebSocket URL：`wss://<server>/api/agents/ws`
- Worker HTTP URL：`https://<server>`（用于注入到 Worker 的 `SERVER_URL`）
- TLS 校验默认关闭（跳过证书校验）
┌─────────────────────────────────────────────────────────────────────────────┐
│                              整体架构                                        │
└─────────────────────────────────────────────────────────────────────────────┘

                                    WebSocket
┌─────────────────┐            (Agent 主动连接)           ┌─────────────────┐
│     Server      │ ◀─────────────────────────────────── │     Agent       │
│   (server/)     │                                       │   (agent/)      │
│                 │ ─────────────────────────────────────▶│                 │
│ - WebSocket Hub │           控制通知/心跳                │ - 连接管理       │
│ - 任务调度       │                                       │ - 任务接收       │
│ - Agent 管理    │                                       │ - Worker 启动   │
└─────────────────┘                                       └─────────────────┘
                                                                  │
                                                                  │ Docker SDK
                                                                  ▼
                                                          ┌─────────────────┐
                                                          │     Worker      │
                                                          │   (临时容器)     │
                                                          │                 │
                                                          │ - 执行扫描       │
                                                          │ - 完成后退出     │
                                                          └─────────────────┘
```

## 任务调度架构

### 调度模式：Agent Pull（业界标准）

采用 **Pull 模式**，参考 Temporal、Celery 等成熟系统：
- Agent 主动拉取任务，而非 Server 推送
- Agent 只在空闲时拉取，自动负载均衡
- 避免任务饥饿问题

```
                    PostgreSQL                              Agent
+-------------------------------------+          +-------------------------+
|  scan_task 表（优先级队列）           |          |                         |
|  +-----------------------------+    |          |  1. 心跳上报负载         |
|  | A2 (stage=1) priority=160  |<---+----------+  2. 检查是否有空闲槽位    |
|  | B2 (stage=1) priority=130  |    |  拉取    |  3. 调用 API 拉取任务     |
|  | C1 (stage=0) priority=20   |    |          |  4. 执行并上报结果        |
|  | D1 (stage=0) priority=10   |    |          |                         |
|  +-----------------------------+    |          +-------------------------+
+-------------------------------------+
```

### 当前实现：单 Workflow（subdomain_discovery）

**当前阶段**：一个 Scan 对应一个 Task，执行单个 Workflow（subdomain_discovery）。

Workflow 内部的多阶段（recon → bruteforce → permutation → resolve）由 Worker 内部编排，定义在 `worker/internal/workflow/subdomain_discovery/templates.yaml`。

```
Scan
+-- scan_task (workflow_name=subdomain_discovery, stage=0)
        │
        └── Worker 内部执行 4 个阶段（recon/bruteforce/permutation/resolve）
```

### 未来扩展：多 Workflow 编排（预留）

**设计预留**：未来支持多个 Workflow 串联（如 subdomain → port_scan → vuln_scan）。

```
# 未来架构（暂不实现）
Scan A
+-- Task A1: subdomain_discovery (stage=0)
+-- Task A2: port_scan (stage=1, depends_on=[A1])
+-- Task A3: vuln_scan (stage=2, depends_on=[A2])
```

**扩展方式**：
1. Worker 端 `templates.yaml` 新增 `depends_on` 和 `output_type` 字段
2. Server 端解析所有 workflow metadata，自动构建 DAG
3. 无需单独的 Server 端 pipeline 配置文件

**依赖触发**：A1 完成 → A2 变为 ready → 加入队列

### 优先级计算

```
priority = stage * 100 + wait_seconds

当前（单 workflow）：所有任务 stage=0，按 wait_seconds 排序
未来（多 workflow）：下游任务 stage 更高，优先执行
```

### 任务队列存储：PostgreSQL

**选择 PostgreSQL 而非 Redis 的原因**：
- 任务持久化，不怕丢失
- 支持复杂查询（统计、报表）
- 事务保证（ACID）
- 扫描任务是分钟级，不需要 Redis 的毫秒级性能

**Redis 用途保留**：
- 心跳缓存（已有）
- 实时通知（Pub/Sub，可选）

## 消息协议

### 消息格式

所有 WebSocket 消息使用 JSON 格式：

```json
{
    "type": "message_type",
    "payload": { ... },
    "timestamp": "2026-01-15T10:30:00Z"
}
```

### 消息类型（WebSocket 控制通道）

#### Agent -> Server

| 类型 | 说明 | Payload |
|------|------|---------|
| `heartbeat` | 心跳 | `{cpu, mem, disk, tasks, version, hostname, uptime}` |

#### Server -> Agent

| 类型 | 说明 | Payload |
|------|------|---------|
| `task_available` | 有新任务可拉取（可选优化） | `{}` |
| `task_cancel` | 取消任务 | `{taskId}` |
| `config_update` | 配置更新 | `{maxTasks, cpuThreshold, memThreshold, diskThreshold}` |
| `update_required` | 需要更新 | `{version, image}` |
| `ping` | 心跳检测 | `{}` |

### 状态管理职责

**核心原则**：Worker 只执行扫描，Agent 负责状态管理。

| 组件 | 职责 |
|------|------|
| Worker | 执行扫描、保存结果到 Server、退出时返回 exit code |
| Agent | 监控 Worker 退出码、调 HTTP API 更新扫描状态 |

**状态流转**：
- `dispatched` - Server 将任务分配给 Agent（拉取成功）
- `running` - Agent 启动 Worker 成功
- `completed` - Worker 退出码为 0
- `failed` - Worker 退出码非 0
- `cancelled` - 任务被取消（用户或系统触发）

**错误信息获取**：
- Worker 失败时，Agent 使用 Docker SDK 读取容器最后 100 行日志
- 日志超过 4KB 时截断，并添加 `[truncated]\n` 前缀
- 错误信息通过 HTTP API 的 `errorMessage` 字段传递给 Server 并存入数据库

这种设计参考了 Kubernetes、Temporal、Celery 等成熟系统的做法：执行者（Worker）只负责执行，调度者（Agent）负责状态管理。

## 组件和接口

### Agent 项目结构

```
agent/
├── cmd/
│   └── agent/
│       └── main.go              # 入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── connection/
│   │   ├── client.go            # WebSocket 客户端
│   │   └── reconnect.go         # 重连策略
│   ├── heartbeat/
│   │   └── reporter.go          # 心跳上报
│   ├── task/
│   │   ├── manager.go           # 任务管理
│   │   └── executor.go          # 任务执行
│   ├── docker/
│   │   └── runner.go            # Docker 容器管理
│   └── pkg/
│       ├── logger.go            # 日志工具
│       └── system.go            # 系统信息采集
├── go.mod
├── go.sum
├── Makefile
└── scripts/
    └── install.sh               # 安装脚本模板
```

### Server 端新增结构

```
server/internal/
├── websocket/
│   ├── hub.go                   # 连接管理中心
│   ├── client.go                # 单个 Agent 连接
│   └── message.go               # 消息定义
├── handler/
│   └── agent_ws.go              # WebSocket 端点处理
└── service/
    └── task_dispatcher.go       # 任务分发服务
```

### 核心组件

#### 1. Config（Agent 配置）

```go
type Config struct {
    // 启动参数（CLI 配置）
    ServerURL      string  // Server 基础地址（IP/域名/端口，不包含协议）
    APIKey         string  // 认证密钥
    AgentName      string  // Agent 名称（默认主机名）
    WorkerImage    string  // Worker 镜像名称
    LogLevel       string  // 日志级别
    
    // 调度参数（Server 动态下发，存储在 DB）
    MaxTasks       int     // 最大并发任务数（默认 5）
    CPUThreshold   int     // CPU 负载阈值（默认 85%）
    MemThreshold   int     // 内存负载阈值（默认 85%）
    DiskThreshold  int     // 磁盘空间阈值（默认 90%）
}
```

#### 2. WebSocketClient（连接管理）

```go
type WebSocketClient interface {
    // 连接 Server
    Connect(ctx context.Context) error
    
    // 发送消息
    Send(msg *Message) error
    
    // 接收消息（阻塞）
    Receive() (*Message, error)
    
    // 关闭连接
    Close() error
    
    // 连接状态
    IsConnected() bool
}
```

#### 3. HeartbeatReporter（心跳上报）

```go
type HeartbeatReporter interface {
    // 启动心跳（每 5 秒）
    Start(ctx context.Context, client WebSocketClient)
    
    // 停止心跳
    Stop()
}

type HeartbeatPayload struct {
    CPU         float64 `json:"cpu"`         // CPU 使用率 (0-100)
    Mem         float64 `json:"mem"`         // 内存使用率 (0-100)
    Disk        float64 `json:"disk"`        // 磁盘使用率 (0-100)
    Tasks       int     `json:"tasks"`       // 运行中任务数
    QueuedTasks int     `json:"queuedTasks"` // 队列中任务数
    Version     string  `json:"version"`     // Agent 版本
    Hostname    string  `json:"hostname"`    // 主机名
    Uptime      int64   `json:"uptime"`      // 运行时长（秒）
}
```

#### 4. TaskExecutor（任务执行器 - Pull 模式）

Agent 主动调用 Server API 拉取任务：

```go
import (
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/shirou/gopsutil/v3/disk"
)

type TaskExecutor struct {
    httpClient     *http.Client         // HTTP 客户端，调用 Server API
    cancelledTasks sync.Map             // 已取消的任务 ID（taskId）
    runningTasks   sync.Map             // 正在执行的任务 {taskId: containerID}
    maxTasks       int                  // 最大并发任务数
    cpuThreshold   float64
    memThreshold   float64
    diskThreshold  float64
    serverURL      string
    apiKey         string
}

// RunningCount 返回当前运行的任务数
func (e *TaskExecutor) RunningCount() int {
    count := 0
    e.runningTasks.Range(func(_, _ any) bool {
        count++
        return true
    })
    return count
}

// Run 核心循环（Pull 模式）
func (e *TaskExecutor) Run(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // 检查是否有空闲槽位
            if !e.canAcceptTask(ctx) {
                continue
            }
            
            // 从 Server 拉取任务
            task, err := e.pullTask(ctx)
            if err != nil {
                log.Error("拉取任务失败", "error", err)
                continue
            }
            if task == nil {
                continue // 无任务可拉取
            }
            
            // 启动 Worker 执行任务
            go e.executeTask(ctx, task)
        }
    }
}

// canAcceptTask 检查是否可以接受新任务
func (e *TaskExecutor) canAcceptTask(ctx context.Context) bool {
    // 检查任务数限制
    if e.RunningCount() >= e.maxTasks {
        return false
    }
    
    // 检查系统负载
    cpuPercent, memPercent, diskPercent, err := getSystemLoad()
    if err != nil {
        return true // 获取失败时允许拉取
    }
    
    return cpuPercent < e.cpuThreshold && 
           memPercent < e.memThreshold && 
           diskPercent < e.diskThreshold
}

// pullTask 调用 Server API 拉取任务
func (e *TaskExecutor) pullTask(ctx context.Context) (*Task, error) {
    req, _ := http.NewRequestWithContext(ctx, "POST", 
        e.serverURL+"/api/agent/tasks/pull", nil)
    req.Header.Set("X-Agent-Key", e.apiKey)
    
    resp, err := e.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 204 {
        return nil, nil // 无任务
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    var task Task
    if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
        return nil, err
    }
    return &task, nil
}

// updateTaskStatus 更新任务状态
func (e *TaskExecutor) updateTaskStatus(ctx context.Context, taskID int, status, errorMsg string) error {
    body, _ := json.Marshal(map[string]string{
        "status": status,
        "errorMessage": errorMsg,
    })
    req, _ := http.NewRequestWithContext(ctx, "PATCH",
        fmt.Sprintf("%s/api/agent/tasks/%d/status", e.serverURL, taskID),
        bytes.NewReader(body))
    req.Header.Set("X-Agent-Key", e.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := e.httpClient.Do(req)
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}

// executeTask 执行任务
func (e *TaskExecutor) executeTask(ctx context.Context, task *Task) {
    // 更新状态为 running
    e.updateTaskStatus(ctx, task.ID, "running", "")
    e.runningTasks.Store(task.ID, "")
    defer e.runningTasks.Delete(task.ID)
    
    // 启动 Worker 容器
    containerID, err := e.docker.Run(ctx, task)
    if err != nil {
        e.updateTaskStatus(ctx, task.ID, "failed", err.Error())
        return
    }
    e.runningTasks.Store(task.ID, containerID)
    
    // 等待容器退出
    exitCode, err := e.docker.Wait(containerID)
    
    if exitCode == 0 {
        e.updateTaskStatus(ctx, task.ID, "completed", "")
    } else {
        // 获取错误日志
        logs := e.docker.GetLogs(containerID, 100)
        e.updateTaskStatus(ctx, task.ID, "failed", logs)
    }
}

// getSystemLoad 获取系统负载
func getSystemLoad() (cpuPercent, memPercent, diskPercent float64, err error) {
    cpuPercents, err := cpu.Percent(0, false)
    if err != nil || len(cpuPercents) == 0 {
        return 0, 0, 0, err
    }
    
    memInfo, err := mem.VirtualMemory()
    if err != nil {
        return 0, 0, 0, err
    }
    
    diskInfo, err := disk.Usage("/")
    if err != nil {
        return 0, 0, 0, err
    }
    
    return cpuPercents[0], memInfo.UsedPercent, diskInfo.UsedPercent, nil
}

// Cancel 取消任务
func (e *TaskExecutor) Cancel(taskID int) {
    e.cancelledTasks.Store(taskID, true)
    
    if containerID, ok := e.runningTasks.Load(taskID); ok {
        e.docker.StopContainer(containerID.(string))
    }
}
```

#### 5. DockerRunner（容器管理，使用 Docker SDK）

```go
type DockerRunner interface {
    // 启动 Worker 容器
    // - 使用 OomScoreAdj=500 让 Worker 优先被 OOM 杀死，保护 Agent
    Run(ctx context.Context, task *Task) (containerID string, err error)
    
    // 停止容器
    Stop(containerID string) error
    
    // 等待容器退出
    Wait(containerID string) (exitCode int, err error)
    
    // 拉取镜像
    Pull(image string) error
    
    // 启动新版本 Agent（用于自更新）
    RunAgent(image string, args []string) error
}
```

#### 6. Updater（自动更新）

```go
type Updater struct {
    docker      DockerRunner
    config      *Config
    currentVer  string
}

// 处理更新指令
func (u *Updater) HandleUpdateRequired(newVersion, image string) error {
    log.Info("收到更新指令", "currentVersion", u.currentVer, "newVersion", newVersion)
    
    // 1. 拉取新镜像
    fullImage := fmt.Sprintf("%s:%s", image, newVersion)
    if err := u.docker.Pull(fullImage); err != nil {
        log.Error("拉取镜像失败", "error", err)
        return err
    }
    
    // 2. 启动新版本容器
    args := []string{
        "--server", u.config.ServerURL,
        "--key", u.config.APIKey,
        "--name", u.config.AgentName,
    }
    if err := u.docker.RunAgent(fullImage, args); err != nil {
        log.Error("启动新容器失败", "error", err)
        return err
    }
    
    // 3. 退出当前进程（让新容器接管）
    log.Info("新版本已启动，退出当前进程")
    os.Exit(0)
    return nil
}
```

#### 7. Hub（Server 端连接管理）

```go
type Hub struct {
    // 注册新连接
    Register(client *AgentClient)
    
    // 注销连接
    Unregister(client *AgentClient)
    
    // 向指定 Agent 发送消息
    SendTo(agentID int, msg *Message) error
    
    // 广播消息
    Broadcast(msg *Message)
    
    // 获取在线 Agent 列表
    GetOnlineAgents() []*AgentClient
    
    // 选择最优 Agent（负载最低）
    SelectBestAgent() *AgentClient
}
```

## 数据模型

### Agent 表

```sql
CREATE TABLE agent (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    api_key         VARCHAR(64) NOT NULL UNIQUE,
    status          VARCHAR(20) DEFAULT 'pending',  -- pending/online/offline
    hostname        VARCHAR(255),
    ip_address      VARCHAR(45),
    version         VARCHAR(20),
    -- 调度配置（可通过 API 动态修改）
    max_tasks       INT DEFAULT 5,
    cpu_threshold   INT DEFAULT 85,
    mem_threshold   INT DEFAULT 85,
    disk_threshold  INT DEFAULT 90,
    connected_at    TIMESTAMP,
    last_heartbeat  TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 状态说明

| 状态 | 说明 |
|------|------|
| pending | 已创建，等待 Agent 连接 |
| online | Agent 已连接，心跳正常 |
| offline | Agent 断开连接或心跳超时 |

### 前端展示优先级（Agent 列表）

1. 主标题：用户自定义名称（若为空，回退 hostname）
2. 副标题：hostname + public IP（last_seen_ip）
3. 其他信息：状态、版本、负载

### scan_task 表（任务队列）

```sql
CREATE TABLE scan_task (
    id              SERIAL PRIMARY KEY,
    scan_id         INT NOT NULL REFERENCES scan(id),
    stage           INT NOT NULL DEFAULT 0,           -- 当前阶段（单 workflow 时固定为 0）
    workflow_name   VARCHAR(100) NOT NULL,            -- e.g. 'subdomain_discovery'
    depends_on      INT[],                            -- 依赖的前置 task id 列表（当前为空，预留）
    status          VARCHAR(20) DEFAULT 'ready',      -- ready/dispatched/running/completed/failed/cancelled (pending 预留用于多 workflow)
    priority        INT DEFAULT 0,                    -- 动态计算：stage*100 + wait_seconds
    agent_id        INT REFERENCES agent(id),         -- 分配给哪个 Agent
    worker_image    VARCHAR(255),
    config          TEXT,                             -- YAML 配置
    error_message   TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    dispatched_at   TIMESTAMP,
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP
);

CREATE INDEX idx_scan_task_status_priority ON scan_task(status, priority DESC);
```

**当前简化**：
- 创建 Scan 时，直接创建一个 `status='ready'` 的 scan_task
- 无需 `pending` 状态（没有依赖需要等待）
- `depends_on` 字段预留，当前为空数组

### 任务状态流转

**当前（单 workflow）**：
```
ready ──────────────→ dispatched ───────→ running ───────→ completed/failed/cancelled
      (Agent 拉取)            (Worker 启动)       (Worker 退出)
```

**未来（多 workflow，预留）**：
```
pending ──────────────→ ready ──────────────→ dispatched ───────→ running
        (依赖完成)               (Agent 拉取)            (Worker 启动)
                                                              │
                                      ┌───────────────────────┴───────────────────────┐
                                      │                                               │
                                      ▼                                               ▼
                                  completed                                        failed
```

### 状态变更职责（scan_task）

scan_task 的状态变更由 Server / Agent 分工完成：

- **Server**：
  - 创建 scan_task：`status=ready`
  - 分配任务（Agent pull 成功）：`ready → dispatched`（写入 `agent_id`, `dispatched_at`）
  - 未来多 workflow：依赖满足时 `pending → ready`
- **Agent**：
  - 启动 Worker 成功后上报：`dispatched → running`（写入 `started_at`）
  - Worker 退出后上报：`running → completed/failed`（写入 `completed_at`, `error_message`）
  - 取消任务：上报 `cancelled`
- **Worker**：只负责执行扫描、保存结果，不直接写 scan_task 状态。

### scan.status 与 scan_task.status 同步（兼容现有 UI/API）

当前阶段仍保留 `scan.status`（用于列表/过滤/历史等）。`scan_task` 是调度层的真实来源，`scan.status` 作为聚合状态同步写入。

**当前（单 workflow, 1 scan = 1 scan_task）同步规则**：
- `scan_task.ready` → `scan.pending`
- `scan_task.dispatched` → `scan.scheduled`
- `scan_task.running` → `scan.running`
- `scan_task.completed` → `scan.completed`
- `scan_task.failed` → `scan.failed`
- `scan_task.cancelled` → `scan.cancelled`

**写入时机**：
- 创建 Scan + scan_task：`scan.status=pending`，`scan_task.status=ready`
- PullTask 分配：同一事务内写入 `scan_task.status=dispatched` + `scan.status=scheduled`
- UpdateStatus（Agent 上报）：更新 `scan_task.status`，同时同步更新 `scan.status`；终态时写入 `scan.stopped_at=NOW()`

**未来（多 workflow）**：
- `scan.status` 由同一 `scan_id` 下所有 scan_task 聚合决定（例如：任一 running 则 running；任一 failed 则 failed；全部 completed 则 completed；全部 cancelled 则 cancelled；否则 pending/scheduled）。

## Agent 注册流程

```
┌──────────────────────────────────────────────────────────────────┐
│  1. 用户在 Web 界面点击「添加 Agent」，输入名称                      │
└──────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────┐
│  2. Server 生成 API Key，创建 agent 记录（status=pending）         │
└──────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────┐
│  3. Web 界面显示部署命令                                           │
│                                                                  │
│     curl ... | bash -s -- --server 1.2.3.4 --key abc123...      │
└──────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────┐
│  4. 用户在远程机器执行命令，Agent 启动并连接 Server                  │
└──────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌──────────────────────────────────────────────────────────────────┐
│  5. Server 验证 Key，更新 agent 记录                               │
│     - status = online                                            │
│     - hostname, ip_address, version（Agent 上报）                 │
│     - connected_at = NOW()                                       │
└──────────────────────────────────────────────────────────────────┘
```

## 数据流

### 连接建立流程

```
Agent                                    Server                     Redis
  │                                        │                          │
  │──────── WebSocket Connect ────────────▶│                          │
  │         (Header: X-Agent-Key)          │                          │
  │                                        │                          │
  │◀──────── config_update ─────────────│                          │
  │         {maxTasks, cpuThreshold,       │                          │
  │          memThreshold, diskThreshold}  │                          │
  │                                        │                          │
  │──────── heartbeat (每5秒) ────────────▶│────── SET ──────────────▶│
  │         {cpu, mem, tasks}              │   agent:{id}:heartbeat   │
  │                                        │   TTL=15s                │
  │                                        │                          │
```

### 心跳数据存储

Server 收到心跳后存入 Redis：

```
Key:    agent:{agent_id}:heartbeat
Value:  {
    "cpu": 45.2,
    "mem": 62.1,
    "disk": 78.5,
    "tasks": 2,
    "version": "1.0.0",
    "hostname": "vps-1",
    "uptime": 86400,
    "updated_at": "..."
}
TTL:    15 秒
```

任务调度时从 Redis 读取所有在线 Agent 的负载数据，选择最优节点。这与现有 Python 版本的 `worker_load_service` 机制一致。

### 任务执行流程（Pull 模式）

```
Agent                          Server                      PostgreSQL             Worker
  │                              │                              │                    │
  │  (检查空闲: tasks < max       │                              │                    │
  │   && cpu/mem/disk < threshold) │                              │                    │
  │                              │                              │                    │
  │── POST /api/agent/tasks/pull ▶│                              │                    │
  │                              │─── SELECT ... FOR UPDATE ───▶│                    │
  │                              │       SKIP LOCKED            │                    │
  │                              │◀─── 返回最高优先级任务 ──────│                    │
  │                              │                              │                    │
  │                              │─── UPDATE status=dispatched ▶│                    │
  │                              │                              │                    │
  │◀─── 200 {taskId, config...} ─│                              │                    │
  │    (或 204 无任务)              │                              │                    │
  │                              │                              │                    │
  │── PATCH .../status=running ─▶│─── UPDATE status=running ──▶│                    │
  │                              │                              │                    │
  │───── Docker SDK run ───────────────────────────────────────────────▶│
  │                              │                              │                    │
  │                              │                              │      (扫描执行中)    │
  │                              │                              │                    │
  │                              │◀──────────────────────────────── HTTP: 保存结果 ───│
  │                              │                              │                    │
  │◀─────────────────────────────────────────────────────── exit (exitCode) ─│
  │                              │                              │                    │
  │── PATCH .../status=completed ▶│── UPDATE status=completed ─▶│                    │
  │    (或 failed + errorMsg)     │                              │                    │
  │                              │── 触发下游任务 ready ──────▶│                    │
```

**说明**：
- Agent 主动拉取（Pull），而非 Server 推送任务
- Server 使用 `FOR UPDATE SKIP LOCKED` 保证并发安全
- Agent 负责任务状态更新（running → completed/failed/cancelled）
- Server 在分配时将任务状态置为 `dispatched`
- Worker 只负责执行扫描和保存结果，退出时返回 exit code
- 任务完成后，Server 自动检查并触发依赖它的下游任务

### HTTP API 端点

#### 1. 拉取任务

```
POST /api/agent/tasks/pull
Header: X-Agent-Key: <agent_api_key>
Request: {}
Response (200): {
    "taskId": 123,
    "scanId": 456,
    "stage": 0,
    "workflowName": "subdomain_discovery",
    "targetId": 789,
    "targetName": "example.com",
    "targetType": "domain",
    "workspaceDir": "/opt/orbit/results/scan_456",
    "config": "subdomain_discovery:\n  ...",
    "workerImage": "yyhuni/orbit-worker:v1.0.19"
}
Response (204): 无任务可拉取
```

#### 2. 更新任务状态

```
PATCH /api/agent/tasks/:id/status
Header: X-Agent-Key: <agent_api_key>
Body: {"status": "running|completed|failed|cancelled", "errorMessage": "..."}
```

### Server 端实现骨架

#### Repository

```go
// internal/repository/scan_task_repository.go
type ScanTaskRepository interface {
    // 原子地拉取最高优先级任务并锁住
    PullTask(ctx context.Context, agentID int) (*model.ScanTask, error)
    
    // 更新任务状态
    UpdateStatus(ctx context.Context, taskID int, status string, errorMsg string) error
    
    // 触发下游任务（依赖完成时将 pending 转为 ready）
    TriggerDependents(ctx context.Context, taskID int) error
    
    // 更新优先级（定时任务，根据等待时间调整）
    RefreshPriorities(ctx context.Context) error
}
```

#### Handler

```go
// internal/handler/agent_task_handler.go
func (h *AgentTaskHandler) PullTask(c *gin.Context) {
    agentID := c.GetInt("agentID") // 从中间件获取
    
    task, err := h.taskRepo.PullTask(c.Request.Context(), agentID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    if task == nil {
        c.Status(204) // No Content
        return
    }
    
    c.JSON(200, dto.TaskAssignResponse{
        TaskID:       task.ID,
        ScanID:       task.ScanID,
        Stage:        task.Stage,
        WorkflowName: task.WorkflowName,
        // ... 其他字段
    })
}

func (h *AgentTaskHandler) UpdateStatus(c *gin.Context) {
    taskID := c.Param("id")
    var req dto.UpdateStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.taskRepo.UpdateStatus(c.Request.Context(), taskID, req.Status, req.ErrorMessage); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 如果任务完成，触发下游任务
    if req.Status == "completed" {
        _ = h.taskRepo.TriggerDependents(c.Request.Context(), taskID)
    }
    
    c.Status(200)
}
```

#### SQL 实现

```go
// PullTask 原子拉取
func (r *scanTaskRepository) PullTask(ctx context.Context, agentID int) (*model.ScanTask, error) {
    var task model.ScanTask
    
    err := r.db.WithContext(ctx).Raw(`
        WITH c AS (
            SELECT id
            FROM scan_task
            WHERE status = 'ready'
            ORDER BY priority DESC, id ASC
            LIMIT 1
            FOR UPDATE SKIP LOCKED
        )
        UPDATE scan_task t
        SET status = 'dispatched', dispatched_at = NOW(), agent_id = ?
        FROM c
        WHERE t.id = c.id
        RETURNING t.*
    `, agentID).Scan(&task).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &task, err
}

// TriggerDependents 触发下游任务
func (r *scanTaskRepository) TriggerDependents(ctx context.Context, taskID int) error {
    // 将依赖已完成任务的 pending 任务转为 ready
    return r.db.WithContext(ctx).Exec(`
        UPDATE scan_task
        SET status = 'ready',
            priority = stage * 100 + EXTRACT(EPOCH FROM (NOW() - created_at))::INT
        WHERE status = 'pending'
          AND ? = ALL(SELECT st.id FROM scan_task st 
                      WHERE st.id = ANY(depends_on) AND st.status = 'completed')
    `, taskID).Error
}
```

### 任务取消流程

```
Server                     Agent                      Worker
  │                          │                          │
  │─── task_cancel ─────────▶│                          │
  │    {taskId}              │                          │
  │                          │──── Docker SDK stop ────▶│
  │                          │                          │
  │◀── HTTP: 更新状态 ───────│                          │
  │    (cancelled)           │                          │
```

### 自动更新流程

```
Agent                          Server                         Docker Hub
  │                              │                                │
  │─── heartbeat ───────────────▶│                                │
  │    {version: "v1.0.8"}       │                                │
  │                              │                                │
  │                              │── 比较版本 ──                   │
  │                              │   agent: v1.0.8                │
  │                              │   server: v1.0.19              │
  │                              │                                │
  │◀── update_required ─────────│                                │
  │    {version: "v1.0.19",      │                                │
  │     image: "orbit-agent"}  │                                │
  │                              │                                │
  │── Docker SDK pull ─────────────────────────────────────────▶│
  │   orbit-agent:v1.0.19      │                                │
  │                              │                                │
  │── Docker SDK run (新容器) ──│                                │
  │   orbit-agent-new          │                                │
  │                              │                                │
  │── 退出当前进程 ──             │                                │
  │   (旧容器停止)               │                                │
  │                              │                                │
  │                              │◀── 新 Agent 连接 ──────────────│
  │                              │    {version: "v1.0.19"}        │
```

**更新策略**：

1. **版本检测**：Server 在收到心跳时比较 `agent_version` 和 `settings.IMAGE_TAG`
2. **触发更新**：版本不匹配时发送 `update_required` 消息
3. **Agent 自更新**：
   - 拉取新版本镜像
   - 启动新版本容器（使用临时名称 `orbit-agent-new`）
   - 新容器启动成功后，退出当前进程
4. **容器切换**：
   - 新容器连接 Server 后，Server 检测到版本匹配
   - 新容器重命名为 `orbit-agent`（或保持 `-new` 后缀）
5. **失败处理**：
   - 镜像拉取失败：记录日志，继续运行
   - 新容器启动失败：记录日志，继续运行
   - 不会影响当前 Agent 的正常工作

**优势**：
- 无需 SSH：Agent 自己执行更新，不需要 Server SSH 到远程机器
- 更安全：不需要存储 SSH 密码
- 更简单：Agent 自己拉取镜像并重启

### Worker 版本管理

Worker 是临时容器，每次任务都是新启动的，因此不需要像 Agent 那样自更新。

```
Server                         Agent                         Docker Hub
  │                              │                                │
  │─── HTTP 拉取返回 ────────────▶│                                │
  │    {workerImage:             │                                │
  │     "worker:v1.0.19"}         │                                │
  │                              │                                │
  │                              │── Docker SDK pull(if missing)+run ─▶│
  │                              │   worker:v1.0.19               │
  │                              │                                │
  │                              │   (本地有就用本地，没有才拉取)   │
  │                              │                                │
```

**版本管理策略**：

| 组件 | 生命周期 | 更新方式 | 触发时机 |
|------|---------|---------|---------|
| Agent | 常驻容器 | 自更新（拉取新镜像 + 重启） | Server 检测到版本不匹配 |
| Worker | 临时容器 | 按需拉取（Docker SDK pull-if-missing） | 每次任务启动时 |

**任务拉取响应格式**：

```json
{
  "taskId": 123,
  "scanId": 456,
  "targetId": 789,
  "targetName": "example.com",
  "targetType": "domain",
  "workflowName": "subdomain_discovery",
  "workspaceDir": "/opt/orbit/results/scan_456",
  "config": "subdomain_discovery:\n  ...",
  "workerImage": "yyhuni/orbit-worker:v1.0.19"
}
```

**优势**：
- Worker 不需要额外的更新逻辑
- 版本由 Server 统一控制
- 支持不同任务使用不同版本（灰度发布）
- 本地有镜像时直接使用，无网络延迟

**默认/覆盖策略**：
- Server 推荐在任务拉取响应中明确指定 `workerImage`（通常与 Server 版本一致）
- Agent 可配置一个默认 Worker 镜像（当响应未提供时使用）

## 重连策略

### 指数退避

```
重试次数    等待时间
   1          1s
   2          2s
   3          4s
   4          8s
   5         16s
   6         32s
   7+        60s (最大)
```

### 重连逻辑

```go
func (c *Client) reconnectLoop(ctx context.Context) {
    backoff := 1 * time.Second
    maxBackoff := 60 * time.Second
    
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }
        
        err := c.Connect(ctx)
        if err == nil {
            backoff = 1 * time.Second // 重置
            c.handleMessages(ctx)
        }
        
        time.Sleep(backoff)
        backoff = min(backoff*2, maxBackoff)
    }
}
```

## Docker 部署

### 一键部署命令

```bash
curl -sSL https://your-server/install-agent.sh | bash -s -- \
  --server 1.2.3.4 \
  --key YOUR_API_KEY
```

用户只需要传 Server 的 IP 或域名，Agent 会自动派生 WebSocket/HTTP 地址。

### Server 生成安装命令/脚本

创建 Agent 时，Server 返回安装所需的字段（示例）：

```json
{
  "agentId": 12,
  "apiKey": "agent_xxx",
  "installCommand": "curl -sSL https://your-server/api/agents/install.sh?key=agent_xxx | bash -s -- --server 1.2.3.4 --key agent_xxx",
  "installScriptUrl": "https://your-server/api/agents/install.sh?key=agent_xxx"
}
```

安装脚本行为：
- 检查 Docker 是否安装
- 创建 `/opt/orbit`
- `docker run --restart=always ... --server <base> --key <apiKey>`
- 若已有旧容器，先 stop/remove 再启动

### 安装脚本功能

1. 检查 Docker 是否安装
2. 创建 `/opt/orbit` 目录（如果不存在）
3. 拉取 Agent 镜像
4. 启动 Agent 容器

### 安装脚本示例

```bash
#!/bin/bash
set -e

SERVER=""
KEY=""
IMAGE="orbit-agent:latest"

# 解析参数
while [[ $# -gt 0 ]]; do
  case $1 in
    --server) SERVER="$2"; shift 2 ;;
    --key) KEY="$2"; shift 2 ;;
    --image) IMAGE="$2"; shift 2 ;;
    *) shift ;;
  esac
done

# 检查参数
if [ -z "$SERVER" ] || [ -z "$KEY" ]; then
  echo "Usage: install-agent.sh --server <ip_or_domain> --key <api_key>"
  exit 1
fi

# 检查 Docker
if ! command -v docker &> /dev/null; then
  echo "Error: Docker is not installed"
  exit 1
fi


# 创建目录
mkdir -p /opt/orbit

# 停止旧容器（如果存在）
docker rm -f orbit-agent 2>/dev/null || true

# 启动 Agent
docker run -d --name orbit-agent \
  --hostname "$(hostname)" \
  --restart=always \
  --oom-score-adj=-500 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /opt/orbit:/opt/orbit \
  "$IMAGE" \
  --server "$SERVER" \
  --key "$KEY"

echo "Agent started successfully"
```

### 升级流程

```bash
docker pull orbit-agent:latest
docker restart orbit-agent
```

### 查看日志

```bash
docker logs -f orbit-agent
```

### 卸载

```bash
docker rm -f orbit-agent
```

## 正确性属性

### Property 1: 连接认证

*对于任意* WebSocket 连接请求，如果 API Key 无效或缺失，Server 应当拒绝连接并返回认证失败消息。

**验证: 需求 3.2, 8.2**

### Property 2: 心跳超时检测

*对于任意* 已连接的 Agent，如果 Server 超过 15 秒未收到心跳，应当将该 Agent 标记为离线并关闭连接。

**验证: 需求 4.3**

### Property 3: 任务容量限制

*对于任意* 任务分配请求，当 Agent 当前任务数达到 max_tasks 时，应当拒绝新任务。

**验证: 需求 5.4**

### Property 4: 重连指数退避

*对于任意* 连接失败序列，重试间隔应当按指数增长（1s, 2s, 4s...），且不超过 60 秒。

**验证: 需求 3.3**

### Property 5: 任务状态完整性

*对于任意* 已接收的任务，Agent 必须发送至少一个终态消息（complete/failed/cancelled）。

**验证: 需求 7.1, 7.2, 7.3**

### Property 6: 容器清理

*对于任意* 启动的 Worker 容器，无论任务成功、失败还是取消，容器都应当被清理（--rm 或手动 stop）。

**验证: 需求 6.2, 6.5**

## 错误处理

### 错误类型

```go
var (
    ErrAuthFailed       = errors.New("authentication failed")
    ErrConnectionLost   = errors.New("connection lost")
    ErrTaskCapacityFull = errors.New("task capacity full")
    ErrTaskNotFound     = errors.New("task not found")
    ErrDockerFailed     = errors.New("docker operation failed")
)
```

### 错误恢复

| 错误场景 | 恢复策略 |
|---------|---------|
| 连接断开 | 自动重连（指数退避） |
| 认证失败 | 记录错误，退出程序 |
| Docker 启动失败 | 上报 task_failed，继续接收新任务 |
| 任务超时 | 停止容器，上报 task_failed |

## 测试策略

### 单元测试

- 配置解析和验证
- 消息序列化/反序列化
- 重连退避计算
- 系统信息采集

### 集成测试

- WebSocket 连接和认证
- 心跳发送和接收
- 任务分配和状态回调
- 容器启动和清理

### 端到端测试

- 完整的任务执行流程
- 断线重连场景
- 多 Agent 负载均衡
