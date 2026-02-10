# Application Contract Standard

本文档定义 `internal/modules/*/application` 的契约（contract）组织规范，用于统一接口放置与命名方式。

## 目标

- 保持应用层依赖倒置：应用层声明能力，基础设施层提供实现。
- 降低接口分散/重复导致的维护成本。
- 在统一风格的同时，避免“一刀切”重构。

## 核心策略（混合策略）

- **单用例接口就近定义**：仅被单个服务/文件使用的接口，放在对应 `*_query.go`、`*_command.go` 或用例文件内。
- **跨用例复用接口集中定义**：被多个应用层文件复用的接口，放在 `contracts.go`。

## 命名规范

- 查询能力接口：`XxxQueryStore`
- 写入能力接口：`XxxCommandStore`
- 组合接口：`XxxStore`（通过 interface embedding 组合 Query/Command）

示例（agent 模块）：

- `AgentQueryStore`
- `AgentCommandStore`
- `AgentStore`（组合前两者）

## `contracts.go` 边界

`contracts.go` 允许：

- 接口定义
- 接口组合
- 错误别名（`var ErrXxx = ...`）
- 类型别名（`type Xxx = ...`）

`contracts.go` 禁止：

- Service 结构体定义
- 业务流程函数实现
- 基础设施实现细节（如 GORM 查询/更新逻辑）
- 与应用层契约无关的默认实现代码

## 渐进迁移准则

- 先改“当前活跃模块”，避免全仓一次性迁移。
- 每次迁移仅做接口组织与依赖收敛，不改变业务行为。
- 迁移后至少通过：
  - 模块 application 单测
  - `check-naming-conventions`
  - `check-layer-dependencies`

## 首批试点

- 模块：`agent`
- 目标：形成可复用样板（接口拆分、构造函数依赖收窄、命名检查守护）
- 后续批次建议顺序：`identity` → `snapshot` → `catalog` → `asset` → `scan` → `security`
