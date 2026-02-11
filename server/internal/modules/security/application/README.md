# security/application

security 模块 application 命名规范：

- `ports.go`：应用层端口接口（漏洞存储、目标查询、原始输出编解码等依赖抽象）。
- `codec.go`：应用层内的编解码职责实现（如 `vulnerabilityJSONRawOutputCodec`）。
- `facade_*.go`：按领域能力聚合对外入口（如 vulnerability）。
- `*_service.go`：可选；仅在 facade 无法承载复杂编排时拆分。
- `errors.go`：可选；仅在该模块定义应用错误时创建。

说明：

- 若实现属于外部能力默认实现，优先放在 `infrastructure` 并按能力命名。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不再使用 `defaults.go`。
- 新代码不在 `application` 层新增 `default_impls.go`。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
