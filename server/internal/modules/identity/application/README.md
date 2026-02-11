# identity/application

identity 模块 application 命名规范：

- `ports.go`：应用层端口接口（仓储、密码哈希、token 能力等依赖抽象）。
- `facade_*.go`：按领域能力暴露对外入口（auth/user/organization）。
- `*_query.go` / `*_command.go`：按读写职责拆分应用逻辑。
- `errors.go`：可选；仅在该模块定义应用错误时创建。

说明：

- 默认实现优先放在 `infrastructure`，并按能力命名（如 `clock.go`、`token_generator.go`、`codec.go`）。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不再使用 `defaults.go`。
- 新代码不在 `application` 层新增 `default_impls.go`。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
