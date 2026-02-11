# scan/application

scan 模块 application 命名规范：

- `ports.go`：应用层端口接口（仓储、运行时、取消器、通知器等依赖抽象）。
- `scan_*_inputs.go`：应用层输入模型（创建请求、列表查询）。
- `scan_*_outputs.go`：应用层输出模型（统计、目标信息等）。
- `facade_*.go`：按用例聚合对外能力（create/query/lifecycle/pull/status）。
- `*_service.go` / `*_runtime.go` / `*_log.go`：按职责拆分应用逻辑。
- `errors.go`：可选；仅在该模块定义应用错误时创建。

说明：

- 默认实现优先放在 `infrastructure`，并按能力命名（如 `clock.go`、`token_generator.go`、`codec.go`）。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不再使用 `defaults.go`。
- 新代码不在 `application` 层新增 `default_impls.go`。
- `application` 层不直接依赖 `dto`；DTO 映射放在 `handler`/`wiring` 边界层。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
