# LunaFox Repository Naming Standard v2

## 1. 目标

统一 `server/agent/worker/frontend` 命名风格，降低认知成本，配合自动守卫防止回退。

## 2. Go 命名规则（server/agent/worker）

### 2.1 application 层

- 允许角色后缀：`query`、`command`、`facade`、`mapper`、`contracts`、`errors`、`aliases`、`adapter`、`codec`、`service`
- 禁止泛名文件：`commands.go`、`ports.go`、`types_alias.go`
- 禁止资源无前缀泛名：`service.go`（必须资源化，如 `scan_command_service.go`、`agent_service.go`）

### 2.2 facade 约束

- `facade_*.go` 只允许：
  - 用例编排
  - 错误语义映射
  - 调用 Query/Command Service
- `facade_*.go` 禁止 DTO 映射函数（`*FromDTO` / `*ToDTO`）
- DTO 映射统一放到 `dto_mappers.go`

### 2.3 术语一致性

- 文件命名统一使用 `host_port`
- 禁止在文件名中出现 `hostport`

### 2.4 agent handler 命名规则

- `internal/modules/agent/handler` 仅允许 `*_handler.go`、`*_mapper.go`
- 禁止 `types.go`、`helpers.go`、`ws_types.go` 等泛名文件
- 运行时接口统一为 `api/agent/*`，管理接口统一为 `api/admin/agents/*`

## 3. Frontend 命名规则

- 目录范围：`frontend/app`、`frontend/components`、`frontend/hooks`、`frontend/lib`、`frontend/services`、`frontend/types`、`frontend/messages`
- `ts/tsx/css` 文件名统一 kebab-case（禁止大写）
- `frontend/services` 下统一 `*.service.ts`
- 禁止 `*.api.ts`

## 4. 守卫脚本

- `server/scripts/check-naming-conventions.sh`
- `scripts/check-frontend-naming.mjs`
- `scripts/check-repo-naming.sh`

## 5. CI 要求

- 在 `.github/workflows/test.yml` 中执行 `bash ./scripts/check-repo-naming.sh`
- 命名守卫失败视为 CI 失败
