# Data Model: WebSocket Agent System

**Feature**: WebSocket Agent System
**Branch**: 001-websocket-agent
**Date**: 2026-01-21

## Overview

本文档定义 WebSocket Agent 系统的数据模型，包括数据库表结构、实体关系和状态机。

## Database Schema

### 1. agent 表

Agent 元数据和配置信息。

```sql
CREATE TABLE agent (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,              -- 自动生成（如 "agent-hostname"）
    api_key         VARCHAR(8) NOT NULL UNIQUE,         -- 认证密钥（8字符hex，4字节随机）
    status          VARCHAR(20) DEFAULT 'online',       -- online/offline（注册时直接 online）
    hostname        VARCHAR(255),                       -- 主机名
    ip_address      VARCHAR(45),                        -- 最后连接的 IP
    version         VARCHAR(20),                        -- Agent 版本号

    -- 调度配置（可通过 API 动态修改）
    max_tasks       INT DEFAULT 5,                      -- 最大并发任务数
    cpu_threshold   INT DEFAULT 85,                     -- CPU 负载阈值 (%)
    mem_threshold   INT DEFAULT 85,                     -- 内存负载阈值 (%)
    disk_threshold  INT DEFAULT 90,                     -- 磁盘空间阈值 (%)

    -- 自注册相关
    registration_token  VARCHAR(8),                     -- 注册时使用的 token（用于审计）

    -- 时间戳
    connected_at    TIMESTAMP,                          -- 最后连接时间
    last_heartbeat  TIMESTAMP,                          -- 最后心跳时间
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_agent_status ON agent(status);
CREATE INDEX idx_agent_api_key ON agent(api_key);
```

**字段说明**：
- `name`: 自动生成（格式：`agent-{hostname}`）
- `api_key`: 8 字符的随机字符串，用于 WebSocket 和 HTTP API 认证
- `status`:
  - `online`: Agent 已连接，心跳正常
  - `offline`: Agent 断开连接或心跳超时（>120 秒）
- `hostname`: Agent 上报的主机名
- `ip_address`: 从 WebSocket 连接中提取的客户端 IP
- `version`: Agent 上报的版本号，用于自动更新判断
- `registration_token`: 注册时使用的 token（用于审计追溯）

### 2. registration_token 表

注册令牌表，用于控制 Agent 自注册的准入权限。

```sql
CREATE TABLE registration_token (
    id              SERIAL PRIMARY KEY,
    token           VARCHAR(8) NOT NULL UNIQUE,         -- 注册令牌（8字符hex）
    expires_at      TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '1 hour'),  -- 过期时间（固定1小时）
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_registration_token_token ON registration_token(token);
CREATE INDEX idx_registration_token_expires ON registration_token(expires_at);
```

**字段说明**：
- `token`: 8 字符的随机字符串（hex），用于 Agent 注册时的身份验证
- `expires_at`: 过期时间，固定为创建后 1 小时
- `created_at`: 创建时间

**使用场景**：
- 批量部署：创建 token 后 1 小时内完成所有 Agent 部署
- Token 在 1 小时内可以无限次使用
- 过期后需要重新生成新的 token

### 3. scan_task 表

任务队列，支持优先级调度。

```sql
CREATE TABLE scan_task (
    id              SERIAL PRIMARY KEY,
    scan_id         INT NOT NULL REFERENCES scan(id) ON DELETE CASCADE,
    stage           INT NOT NULL DEFAULT 0,             -- 当前阶段（单 workflow 时固定为 0）
    workflow_name   VARCHAR(100) NOT NULL,              -- e.g. 'subdomain_discovery'
    status          VARCHAR(20) DEFAULT 'pending',      -- pending/running/completed/failed/cancelled
    version         VARCHAR(20),                        -- Worker 版本号（从 VERSION 文件读取）

    -- 分配信息
    agent_id        INT REFERENCES agent(id),           -- 分配给哪个 Agent
    config          TEXT,                               -- YAML 配置
    error_message   VARCHAR(4096),                      -- 错误信息（Agent 截断，对齐 K8s termination message）

    -- 重试控制
    retry_count     INT DEFAULT 0,                      -- 已重试次数

    -- 时间戳
    created_at      TIMESTAMP DEFAULT NOW(),
    started_at      TIMESTAMP,                          -- Worker 启动时间
    completed_at    TIMESTAMP                           -- 完成时间
);

CREATE INDEX idx_scan_task_pending_order ON scan_task(status, stage DESC, created_at ASC);
CREATE INDEX idx_scan_task_agent_id ON scan_task(agent_id);
CREATE INDEX idx_scan_task_scan_id ON scan_task(scan_id);
```

**字段说明**：
- `stage`: 任务阶段，用于多 workflow 串联（当前固定为 0）
- `workflow_name`: 工作流名称，对应 Worker 的 templates.yaml 中的定义
- `status`: 任务状态（pending/running/completed/failed/cancelled）
- `version`: Worker 版本号（Server 从 VERSION 文件读取，Agent 拼接镜像名称：yyhuni/orbit-worker:v{VERSION}）
- `agent_id`: 分配给哪个 Agent（拉取时写入）
- `config`: YAML 格式的任务配置

### 4. Agent 自注册流程

系统采用 **Token 控制的自注册模式**，Agent 通过注册令牌自动注册到系统。

**完整流程**：
```
1. Admin 生成注册 Token
   → POST /api/registration-tokens
   → Server 返回 token 和安装命令（token 1小时后过期）

2. Admin 批量部署
   → for vps in vps-{1..100}; do
       ssh $vps "curl install.sh | bash -s <token>"
     done

3. Agent 使用 Token 注册
   → POST /api/agents/register
     { "token": "<token>", "hostname": "vps-1", "version": "1.0.0" }
   → Server 验证 token（检查是否过期）
   → Server 生成专属 api_key，创建 Agent 记录（status=online）
   → 返回 api_key 给 Agent

4. Agent 保存 api_key 并发送心跳
   → Agent 将 api_key 保存到本地配置文件
   → 使用 api_key 建立 WebSocket 连接
   → 发送心跳数据
```

**Token 生命周期**：
```
创建 Token (默认1小时后过期)
    ↓
VPS-1 注册 → 返回 key_001
VPS-2 注册 → 返回 key_002
...
VPS-N 注册 → 返回 key_N (在过期前可无限次使用)
    ↓
Token 过期 → 新注册请求被拒绝
    ↓
已注册的 Agent 继续使用各自的 api_key 正常工作
```

**Token 验证逻辑**：
```sql
-- 验证 token 是否有效
SELECT * FROM registration_token
WHERE token = $1
  AND expires_at > NOW();
```

**适用场景**：
- 批量部署（一个 token 部署多台服务器）
- 快速上手（一键安装命令）
- 公网环境（需要准入控制）
- 临时授权（设置过期时间）

### 5. Redis 缓存结构

#### 心跳数据

```
Key:    agent:{agent_id}:heartbeat
Value:  {
    "cpu": 45.2,              // CPU 使用率 (0-100)
    "mem": 62.1,              // 内存使用率 (0-100)
    "disk": 78.5,             // 磁盘使用率 (0-100)
    "tasks": 2,               // 运行中任务数
    "version": "1.0.0",       // Agent 版本
    "hostname": "vps-1",      // 主机名
    "uptime": 86400,          // 运行时长（秒）
    "updated_at": "2026-01-21T10:30:00Z"
}
TTL:    15 秒
```

## Entity Relationships

```
┌──────────────────────┐
│ RegistrationToken    │
│  (注册令牌)           │
└──────────┬───────────┘
           │ 1
           │ (用于注册)
           │ N
           ▼
┌─────────────────┐
│     Agent       │
│  (常驻服务)      │
└────────┬────────┘
         │ 1
         │
         │ N
         ▼
┌─────────────────┐         ┌─────────────────┐
│   ScanTask      │ N     1 │      Scan       │
│  (任务队列)      │◀────────│   (扫描记录)     │
└─────────────────┘         └─────────────────┘
         │
         │ depends_on (预留)
         │
         ▼
┌─────────────────┐
│   ScanTask      │
│  (下游任务)      │
└─────────────────┘
```

**关系说明**：
- 一个 RegistrationToken 可以用于注册多个 Agent（在过期前无限次使用）
- 一个 Agent 可以执行多个 ScanTask
- 一个 Scan 对应一个或多个 ScanTask（当前为 1:1，未来支持 1:N）
- ScanTask 之间可以有依赖关系（预留，当前未使用）

## State Machines

### Agent 状态机

```
online ──────────────→ offline
       (心跳超时/断开)
          │
          └──────────────────→ online
                (重连成功)
```

**状态转换规则**：
- Agent 注册时直接创建为 `online` 状态
- `online → offline`: 120 秒未收到心跳或连接断开
- `offline → online`: Agent 重连成功

### ScanTask 状态机

**统一状态**（scan.status 和 scan_task.status 使用相同的 5 个状态）：
```
pending ────────────────→ running ───────────────→ completed
        (Agent 拉取并启动)         (退出码=0)
                                          │
                                          ├──→ failed
                                          │    (退出码≠0)
                                          │
                                          └──→ cancelled
                                               (用户取消)
```

**状态转换职责**：
- **Server**:
  - 创建 scan_task: `status=pending`
  - 分配任务（Agent pull 成功）: `pending → running`（分配和启动合并为一步）
  - 任务回收（Agent 离线）: `running → pending`（重试）或 `→ failed`（超过重试次数）
- **Agent**:
  - Worker 退出: `running → completed/failed/cancelled`

**未来多 workflow 扩展**（预留）：
- 依赖通过查询条件过滤，不新增状态

### 任务回收机制（Agent 离线时）

当 Agent 离线（心跳超时 >120 秒）时，其名下所有 `running` 状态的任务需要回收：

**回收规则**：
- `retry_count < 3`: 重置为 `pending`，`retry_count += 1`，重新进入队列
- `retry_count >= 3`: 标记为 `failed`，错误信息为 "Agent lost, max retries exceeded"

**后台 Job 逻辑**（每分钟执行一次）：

```sql
-- Step 1: 标记心跳超时的 Agent 为 offline
UPDATE agent 
SET status = 'offline' 
WHERE status = 'online' 
  AND last_heartbeat < NOW() - INTERVAL '120 seconds';

-- Step 2: 回收离线 Agent 的任务
UPDATE scan_task
SET 
  status = CASE WHEN retry_count >= 3 THEN 'failed' ELSE 'pending' END,
  agent_id = NULL,
  retry_count = retry_count + 1,
  error_message = CASE WHEN retry_count >= 3 THEN 'Agent lost, max retries exceeded' ELSE NULL END
WHERE status = 'running'
  AND agent_id IN (SELECT id FROM agent WHERE status = 'offline');

-- Step 3: 同步更新受影响的 scan.status（状态已统一，直接复制）
UPDATE scan s
SET status = t.status
FROM scan_task t
WHERE s.id = t.scan_id 
  AND t.status IN ('pending', 'failed')
  AND s.status NOT IN ('completed', 'failed', 'cancelled');
```

**设计理由**：
- Agent 离线 ≠ 任务本身有问题，换个 Agent 跑大概率能成功
- 重试次数限制防止无限循环
- 扫描任务是幂等的，偶发重复执行不会造成数据错误

## Data Validation Rules

### Agent
- `name`: 1-100 字符，非空，自动生成格式：`agent-{hostname}`
- `api_key`: 8 字符，唯一，自动生成
- `status`: 枚举值 `online|offline`
- `max_tasks`: 1-100，默认 5
- `cpu_threshold`: 1-100，默认 85
- `mem_threshold`: 1-100，默认 85
- `disk_threshold`: 1-100，默认 90
- `registration_token`: 8 字符（可选，仅自注册时填充）

### RegistrationToken
- `token`: 8 字符，唯一，自动生成（hex）
- `expires_at`: 时间戳，必填，默认创建后 1 小时过期
- **业务规则**：
  - 过期的 token 不能用于注册
  - Token 在过期前可以无限次使用

### ScanTask
- `workflow_name`: 非空，必须在 Worker templates.yaml 中定义
- `status`: 枚举值 `pending|running|completed|failed|cancelled`
- `config`: 有效的 YAML 格式
- `error_message`: 最大 4KB（Agent 截断，对齐 K8s termination message 限制）

## Indexes and Performance

### 关键索引
- `idx_scan_task_pending_order`: 支持任务拉取查询（按 stage DESC, created_at ASC 排序）
- `idx_agent_status`: 支持在线 Agent 查询
- `idx_scan_task_agent_id`: 支持按 Agent 查询任务

### 查询优化
- 任务拉取使用 `FOR UPDATE SKIP LOCKED` 避免锁竞争
- 心跳数据存储在 Redis，避免频繁写入 PostgreSQL
- Agent 状态更新使用乐观锁（`updated_at` 检查）

### 任务拉取 SQL 示例
```sql
-- 原子操作：选取一个 pending 任务并设置为 running
WITH selected AS (
  SELECT id FROM scan_task
  WHERE status = 'pending'
  ORDER BY stage DESC, created_at ASC
  FOR UPDATE SKIP LOCKED
  LIMIT 1
)
UPDATE scan_task t
SET status = 'running', agent_id = $1, started_at = NOW()
FROM selected
WHERE t.id = selected.id
RETURNING t.*;
```

## Migration Strategy

### 新增表
- `agent`: 全新表，替代现有的 `worker_node` 表
- `scan_task`: 全新表
- `registration_token`: 全新表，用于自注册模式

### 修改现有表
- `scan`: 删除 `worker_id` 字段（任务分配关系移至 `scan_task.agent_id`）
- `scan.status` 与 `scan_task.status` 同步更新（当前 1:1 映射）

### 弃用表
- `worker_node`: SSH 方式已弃用，由 `agent` 表替代

### 数据迁移
- 新创建的 `scan` 自动创建对应的 `scan_task`（在同一事务中）
- 状态统一迁移：
  ```sql
  -- 将旧状态统一为新的 5 状态模型
  UPDATE scan SET status = 'pending' WHERE status IN ('initiated', 'scheduled');
  ```

## Future Extensions

### 多 Workflow 支持（预留）
- `scan_task.depends_on`: 存储依赖的 task ID 数组
- 依赖通过拉取时的查询条件过滤，不新增状态

### 任务排序规则
- 拉取时按 `ORDER BY stage DESC, created_at ASC` 排序
- 高 stage 优先（让流水线尽快完成）
- 同 stage 内先创建的先执行（防止饿死）

### Agent 分组（可选）
- 新增 `agent_group` 表
- 支持将 Agent 分配到不同的组
- 任务可以指定目标 Agent 组
