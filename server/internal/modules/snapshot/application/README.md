# snapshot/application

snapshot 模块 application 命名规范：

- `ports.go`：应用层端口接口（query/command store、asset sync、scan lookup、codec 等依赖抽象）。
- `query_inputs.go`：应用层查询入参模型（分页、过滤、排序）。
- `facade_*.go`：按业务视角聚合对外能力（web/discovery/port-capture/vulnerability）。
- `*_snapshot.go`：按资产类型拆分快照 query/command 逻辑。

说明：

- 默认实现优先放在 `infrastructure`，并按能力命名（如 `clock.go`、`token_generator.go`、`codec.go`）。

约束：

- 新代码不再使用 `contracts.go`，统一使用 `ports.go`。
- 新代码不再使用 `defaults.go`。
- 新代码不在 `application` 层新增 `default_impls.go`。
- `application` 层不直接依赖 `dto`；DTO 映射放在 `handler`/`wiring` 边界层。
- 避免使用弱语义泛名文件（如 `types.go`、`common.go`）。
