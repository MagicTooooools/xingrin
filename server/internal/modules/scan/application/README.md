# scan/application

scan 模块 application 命名规范：

- `ports.go`：应用层端口接口（仓储、运行时、取消器、通知器等依赖抽象）。
- `default_impls.go`：可选；仅在存在端口默认实现时创建。
- `facade_*.go`：按用例聚合对外能力（create/query/lifecycle/pull/status）。
- `*_service.go` / `*_runtime.go` / `*_log.go`：按职责拆分应用逻辑。
- `errors.go`：可选；仅在该模块定义应用错误时创建。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不再使用 `defaults.go`，统一使用 `default_impls.go`。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
