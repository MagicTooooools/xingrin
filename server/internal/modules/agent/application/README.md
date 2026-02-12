# agent/application

通用规则请遵循：`docs/server-application-naming-template-v1.md`

agent 模块补充规则：

- **入口聚合**：对外入口单点收敛在 `agent_facade.go`。
- **服务编排**：主流程按 query/command/runtime/task/registration 拆分为 `agent_*_service.go`。
- **端口拆分**：在 `query_ports.go` / `command_ports.go` 之外，依赖能力单独拆分（如 `agent_store_ports.go`、`agent_message_ports.go`、`clock_ports.go`、`token_generator_ports.go`、`heartbeat_cache_ports.go`）。
- **模型命名**：新增跨边界模型优先资源化命名（如 `*_item_models.go`、`*_query_inputs.go`）；默认实现保持在 `infrastructure` 并由 wiring 注入。
- **历史迁移**：`errors.go` 为历史保留，后续新增错误文件优先资源化 `*_errors.go`。
