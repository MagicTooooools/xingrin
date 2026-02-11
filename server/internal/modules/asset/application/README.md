# asset/application

asset 模块 application 命名规范：

- `ports.go`：应用层端口接口（仓储、目标查询等依赖抽象）。
- `facade_*.go`：按聚合能力暴露对外入口（如 website/subdomain/endpoint 等）。
- `*_query.go` / `*_command.go`：按读写职责拆分应用逻辑。
- `errors.go`：应用层错误定义。

说明：

- 默认实现优先放在 `infrastructure`，并按能力命名（如 `clock.go`、`token_generator.go`、`codec.go`）。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不再使用 `defaults.go`。
- 新代码不在 `application` 层新增 `default_impls.go`。
- `application` 层不直接依赖 `dto`；DTO 映射放在 `handler`/`wiring` 边界层。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
