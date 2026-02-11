# agent/application

agent 模块 application 命名规范：

- `ports.go`：应用层端口接口（仓储、消息、时钟、令牌生成等依赖抽象）。
- `agent_facade.go`：对外聚合入口，组合 query/command/registration 能力。
- `agent_*_service.go`：按职责拆分的应用服务实现。
- `errors.go`：应用层错误定义。

说明：

- 端口默认实现下沉到 `server/internal/modules/agent/infrastructure/clock.go` 与 `server/internal/modules/agent/infrastructure/token_generator.go`，由 wiring 显式注入。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不在 `application` 层新增默认实现类型。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
