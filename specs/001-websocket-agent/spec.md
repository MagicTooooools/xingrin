# Feature Specification: WebSocket Agent System

**Feature Branch**: `001-websocket-agent`
**Created**: 2026-01-21
**Status**: Draft
**Input**: 实现基于 WebSocket 的轻量级 Agent 系统，用于替代 SSH 任务分发方式。Agent 作为常驻服务运行在远程 VPS 上，主动连接 Server 建立长连接，通过 Pull 模式拉取任务并启动临时 Worker 容器执行扫描。

## User Scenarios & Testing

### User Story 1 - Agent 部署和连接 (Priority: P1)

作为运维人员，我希望能够快速在远程 VPS 上部署 Agent 并连接到 Server，以便开始接收和执行扫描任务。

**Why this priority**: 这是整个系统的基础，没有 Agent 连接就无法执行任何任务。这是 MVP 的核心功能。

**Independent Test**: 可以通过在远程机器上运行安装命令，验证 Agent 成功连接到 Server 并显示在线状态来独立测试。

**Acceptance Scenarios**:

1. **Given** 我在 Web 界面创建了一个新 Agent，**When** 我在远程 VPS 上执行提供的安装命令，**Then** Agent 应该成功启动并在 Web 界面显示为在线状态
2. **Given** Agent 已连接到 Server，**When** 网络临时断开，**Then** Agent 应该自动重连并恢复在线状态
3. **Given** 我提供了错误的 API Key，**When** Agent 尝试连接，**Then** 连接应该被拒绝并显示认证失败错误

---

### User Story 2 - 任务执行和状态跟踪 (Priority: P1)

作为系统管理员，我希望 Agent 能够自动拉取任务并执行扫描，同时实时更新任务状态，以便我能够监控扫描进度。

**Why this priority**: 这是系统的核心价值所在，Agent 必须能够执行任务才能替代 SSH 方式。

**Independent Test**: 可以通过创建一个扫描任务，验证 Agent 自动拉取、执行并更新状态来独立测试。

**Acceptance Scenarios**:

1. **Given** Agent 处于空闲状态且系统负载正常，**When** Server 有新的扫描任务，**Then** Agent 应该自动拉取任务并启动 Worker 容器执行
2. **Given** Worker 容器正在执行扫描，**When** 扫描完成（退出码为 0），**Then** 任务状态应该更新为 completed
3. **Given** Worker 容器执行失败（退出码非 0），**When** 容器退出，**Then** 任务状态应该更新为 failed 并包含错误日志
4. **Given** 用户在 Web 界面取消任务，**When** 取消指令发送到 Agent，**Then** Agent 应该停止对应的 Worker 容器并更新状态为 cancelled

---

### User Story 3 - 负载监控和智能调度 (Priority: P2)

作为系统管理员，我希望 Agent 能够监控自身负载并智能决定是否接受新任务，以便避免系统过载。

**Why this priority**: 这确保了系统的稳定性和可靠性，防止单个 Agent 过载导致任务失败。

**Independent Test**: 可以通过模拟高负载场景，验证 Agent 拒绝新任务来独立测试。

**Acceptance Scenarios**:

1. **Given** Agent 的 CPU 使用率超过阈值（默认 85%），**When** 检查是否可以接受新任务，**Then** Agent 应该等待直到负载降低
2. **Given** Agent 已达到最大并发任务数（默认 5），**When** 尝试拉取新任务，**Then** Agent 应该等待直到有任务完成
3. **Given** Agent 定期上报心跳，**When** Server 收到心跳数据，**Then** Server 应该记录 Agent 的 CPU、内存、磁盘使用率和当前任务数

---

### User Story 4 - 配置动态更新 (Priority: P2)

作为系统管理员，我希望能够在 Web 界面动态调整 Agent 的配置参数，以便根据实际情况优化性能。

**Why this priority**: 这提供了灵活性，允许在不重启 Agent 的情况下调整配置。

**Independent Test**: 可以通过在 Web 界面修改配置，验证 Agent 立即应用新配置来独立测试。

**Acceptance Scenarios**:

1. **Given** Agent 已连接到 Server，**When** 管理员在 Web 界面修改最大任务数，**Then** Agent 应该立即接收并应用新配置
2. **Given** Agent 使用旧的负载阈值，**When** Server 推送新的阈值配置，**Then** Agent 应该使用新阈值进行负载检查

---

### User Story 5 - 自动更新 (Priority: P3)

作为运维人员，我希望 Agent 能够自动更新到最新版本，以便无需手动干预即可获得新功能和修复。

**Why this priority**: 这是便利性功能，可以简化运维工作，但不是核心功能。

**Independent Test**: 可以通过发布新版本，验证 Agent 自动拉取镜像并重启来独立测试。

**Acceptance Scenarios**:

1. **Given** Agent 运行旧版本，**When** Server 检测到版本不匹配，**Then** Server 应该发送更新指令
2. **Given** Agent 收到更新指令，**When** 开始更新流程，**Then** Agent 应该拉取新镜像、启动新容器并退出旧容器
3. **Given** 镜像拉取失败，**When** 更新失败，**Then** Agent 应该记录错误并继续运行当前版本

---

### Edge Cases

- 当 Agent 正在执行任务时收到更新指令，应该等待任务完成后再更新
- 当网络不稳定导致频繁断连时，应该使用指数退避策略避免过度重连
- 当 Docker 守护进程不可用时，Agent 应该记录错误并等待 Docker 恢复
- 当磁盘空间不足时，Agent 应该拒绝新任务并上报警告
- 当 Worker 容器长时间未响应时，应该有超时机制强制停止
- 当多个 Agent 同时拉取任务时，应该使用数据库锁避免重复分配

## Requirements

### Functional Requirements

- **FR-001**: Agent 必须作为独立的 Go 二进制文件编译，支持 Linux amd64 和 arm64 架构
- **FR-002**: Agent 必须支持通过命令行参数配置 Server 地址和 API Key
- **FR-003**: Agent 必须主动连接 Server 的 WebSocket 端点并进行认证
- **FR-004**: Agent 必须在连接失败时使用指数退避策略自动重连（1s, 2s, 4s, 8s, 最大 60s）
- **FR-005**: Agent 必须每 5 秒发送一次心跳消息，包含 CPU、内存、磁盘使用率、任务数、版本号和主机名
- **FR-006**: Server 必须在 120 秒未收到心跳时将 Agent 标记为离线
- **FR-007**: Agent 必须通过 HTTP API 主动拉取任务（Pull 模式），API 路径为 `/api/agent/tasks/*`（操作 scan_task，非 scan）
- **FR-026**: Agent 拉取策略：收到 WS task_available 通知时立即拉取；拉取返回 204 后退避等待（5s/10s/30s，最大 60s）；收到新通知时重置退避；拉取间隔根据当前负载动态调整（负载 <50% 时 1 秒，50-80% 时 3 秒，>80% 时 10 秒），实现自动负载均衡
- **FR-008**: Agent 必须在满足以下条件时拉取任务：当前任务数 < max_tasks 且 CPU/内存/磁盘使用率低于阈值
- **FR-009**: Server 必须使用 PostgreSQL 行级锁（FOR UPDATE SKIP LOCKED）确保任务不被重复分配
- **FR-010**: Agent 必须使用 Docker SDK 启动 Worker 容器，并传递任务参数作为环境变量
- **FR-011**: Agent 必须监控 Worker 容器的退出码，并根据退出码更新任务状态（0=completed, 非0=failed）
- **FR-012**: Agent 必须在 Worker 失败时读取容器最后 100 行日志作为错误信息（超过 4KB 时截断，数据库字段 error_message 为 VARCHAR(4096)）
- **FR-013**: Agent 必须响应 Server 的任务取消指令，停止对应的 Worker 容器
- **FR-014**: Agent 必须在 Worker 容器退出后先读取日志再删除容器（不使用 --rm，手动清理）
- **FR-025**: Agent 必须对每个 Worker 容器设置最大运行时长（默认 7 天），超时强制停止并标记任务为 failed
- **FR-015**: Server 必须提供一键安装脚本，包含 Agent 部署命令和 API Key。脚本应包含：拉取 Agent Docker 镜像、使用提供的 API Key 和 Server 地址启动 Agent 容器（挂载 Docker socket 和 /opt/orbit 目录）
  - **Server 地址生成规则**（最少交互）：
    1) 若配置了 `PUBLIC_URL`（完整 URL，含协议/域名/端口），直接使用；
    2) 否则从用户访问 `GET /api/agents/install.sh` 的请求 URL 推断（基于 Host/Proto/Port 头）。
  - 安装脚本内写入 `SERVER_URL=<PUBLIC_URL or inferred URL>`，Agent 启动后使用该值访问 HTTP API，并将 `https→wss`、`http→ws` 自动转换用于 WebSocket 连接。
- **FR-016**: Agent 必须支持接收 Server 推送的配置更新（maxTasks、cpuThreshold、memThreshold、diskThreshold）
- **FR-017**: Agent 必须在收到更新指令时拉取新版本镜像、启动新容器并退出当前进程
- **FR-018**: Worker 容器必须使用 oom-score-adj=500 提高被 OOM 杀死的优先级，保护 Agent（oom-score-adj=-500）
- **FR-019**: Agent 必须使用 gopsutil 库正确处理容器 cgroup 限制来采集系统负载（在容器环境中从 /sys/fs/cgroup 读取指标）
- **FR-020**: Server 必须在分配任务时将 scan_task.status 更新为 running，并同步更新 scan.status 为 running（分配和启动合并为一步）
- **FR-027**: Server 必须在创建 Scan 时同时创建对应的 scan_task 记录（status=pending, workflow_name 取 YamlConfiguration 的第一个顶层 key，即扫描配置 YAML 中定义的工作流名称，如 "subdomain_discovery"，version 从 VERSION 文件读取，Agent 收到后自行拼接镜像名称为 yyhuni/orbit-worker:v{VERSION}）
- **FR-028**: Server 必须运行后台 Job（每分钟执行），负责：1) 标记心跳超时的 Agent 为 offline；2) 回收离线 Agent 的任务
- **FR-021**: Server 必须在 Agent 离线时回收其名下的 running 任务，重试次数 <3 时重置为 pending，否则标记为 failed
- **FR-022**: Server 必须校验状态更新请求的 Agent 所有权（agent_id 匹配），不匹配返回 403
- **FR-023**: Server 必须保证状态更新幂等，重复上报相同状态返回 200
- **FR-024**: Server 必须拒绝非法状态转换（仅允许 pending→running、running→completed/failed/cancelled）

### Key Entities

- **Agent**: 常驻服务，包含属性：ID、名称、API Key、状态（pending/online/offline）、主机名、IP 地址、版本号、调度配置（max_tasks、cpu_threshold、mem_threshold、disk_threshold）、连接时间、最后心跳时间
- **ScanTask**: 扫描任务，包含属性：ID、scan_id、stage、workflow_name、状态（pending/running/completed/failed/cancelled）、agent_id、worker_image、配置（YAML）、错误信息、retry_count、时间戳
- **Heartbeat**: 心跳数据，包含属性：CPU 使用率、内存使用率、磁盘使用率、运行中任务数、版本号、主机名、运行时长
- **Worker**: 临时容器，执行具体扫描任务，完成后自动删除

## Success Criteria

### Measurable Outcomes

- **SC-001**: Agent 安装和连接过程在 2 分钟内完成（从执行安装命令到显示在线状态）
- **SC-002**: Agent 在网络恢复后 120 秒内自动重连成功
- **SC-003**: 任务从创建到被 Agent 拉取的延迟不超过 5 秒（在 Agent 空闲且负载正常的情况下）
- **SC-004**: 任务状态更新的延迟不超过 2 秒（从 Worker 退出到状态更新完成）
- **SC-005**: Agent 在 CPU 使用率超过 85% 时不接受新任务，确保系统稳定性
- **SC-006**: 单个 Agent 支持同时运行至少 5 个并发任务
- **SC-007**: 心跳数据每 5 秒更新一次，Server 在 120 秒未收到心跳时准确标记 Agent 离线
- **SC-008**: 配置更新在推送后 5 秒内被 Agent 应用
- **SC-009**: Agent 自动更新过程在 5 分钟内完成（包括镜像拉取和容器重启）
- **SC-010**: 多个 Agent 同时拉取任务时，不会出现任务重复分配（通过数据库锁保证）
- **SC-011**: Agent 内存占用不超过 50MB（空闲状态）
- **SC-012**: Worker 容器在任务完成后 100% 被清理，不留下僵尸容器

## Assumptions

- Docker 已在远程 VPS 上安装并正常运行
- Server 和 Agent 之间的网络连接支持 WebSocket（wss://）
- PostgreSQL 数据库版本支持 FOR UPDATE SKIP LOCKED 语法（9.5+）
- Agent 运行的机器有足够的磁盘空间存储 Worker 镜像和扫描结果
- Server 使用 HTTPS/WSS 协议，TLS 证书校验默认关闭（用于自签名证书场景）
- 单个 Worker 容器的资源需求不会超过 Agent 机器的总资源
- 扫描任务的执行时间通常在分钟到小时级别，不需要毫秒级的任务调度
- Agent 和 Worker 使用相同的 /opt/orbit 目录进行数据交换

## Out of Scope

- Agent 的 Web 管理界面（通过 Server 的 Web 界面管理）
- Agent 之间的直接通信（所有通信通过 Server 中转）
- 任务优先级的手动调整（优先级由系统自动计算）
- 多 Workflow 串联执行（当前仅支持单个 subdomain_discovery workflow，多 workflow 功能预留）
- Agent 的日志聚合和分析（Agent 只负责本地日志记录）
- Worker 容器的资源限制配置（使用 Docker 默认设置）
- Agent 的健康检查端点（通过心跳机制实现健康监控）
- 任务级别的超时控制（由 Worker 内部实现；Agent 仅实现容器级别的 7 天超时保护，见 FR-025）

## Dependencies

- Docker Engine（Agent 使用 Docker SDK 管理容器）
- PostgreSQL 数据库（Server 端任务队列存储）
- Redis（Server 端心跳数据缓存）
- Go 1.21+（Agent 开发语言）
- gopsutil v3（系统负载采集库）
- gorilla/websocket（WebSocket 客户端库）

## Version Management

**统一版本管理**：所有组件（Server、Agent、Worker）的版本号统一由项目根目录的 `VERSION` 文件管理。

**版本文件位置**：`/Users/yangyang/Desktop/orbit/VERSION`

**版本使用方式**：
- **Server**：启动时读取 VERSION 文件，创建 scan_task 时拼接 Worker 镜像名称（格式：`yyhuni/orbit-worker:v{VERSION}`）
- **Agent**：启动时读取 VERSION 文件，在心跳消息中上报版本号给 Server
- **Worker**：Docker 镜像的 tag 使用 VERSION 文件内容（如：`yyhuni/orbit-worker:v1.5.12-dev`）

**版本升级流程**：
1. 更新 `VERSION` 文件内容（如：`v1.5.12-dev` → `v1.5.13`）
2. 重新构建并发布 Worker 镜像（tag 使用新版本号）
3. 重启 Server（读取新版本号）
4. 新创建的 scan_task 自动使用新版本 Worker 镜像
5. Agent 拉取任务时获取新镜像名称，自动拉取并使用新版本 Worker

**优势**：
- 统一管理：只需修改一个文件即可更新所有组件版本
- 版本一致：确保 Server、Agent、Worker 版本同步
- 自动升级：无需手动修改配置，新任务自动使用新版本
- 易于追踪：所有组件版本号一致，便于问题排查和回滚

## Security Considerations

### 认证体系（3 种调用者）

| 调用者 | 认证方式 | Header | 用途 |
|--------|---------|--------|------|
| 用户（前端） | JWT | `Authorization: Bearer <token>` | Web 界面操作 |
| Worker（容器） | 全局静态 Token | `X-Worker-Token` | 保存扫描结果 |
| Agent | 每个 Agent 独立 Key | `X-Agent-Key` | 任务拉取、状态更新 |

- Agent 的 `api_key` 存储在 `agent` 表，每个 Agent 一个独立的 key
- Worker 的 token 是全局配置（Worker 是 Agent 启动的临时容器，不需要独立认证）
- WebSocket 认证：由于部分环境不支持自定义 Header，支持两种方式：
  - Header: `X-Agent-Key: <key>`
  - Query: `wss://server/api/agents/ws?key=<key>`

### 其他安全要求

- API Key 必须安全存储，不应出现在日志或错误信息中
- WebSocket 连接必须使用 wss:// 协议（加密传输）
- Agent 不应信任来自 Worker 的任何输入（Worker 只能通过 HTTP API 上报结果）
- Docker socket 挂载到 Agent 容器时需要注意权限控制
- 错误日志在上报前应该截断，避免泄露敏感信息
- Agent 的 oom-score-adj 设置为 -500，确保在内存不足时优先保护 Agent
