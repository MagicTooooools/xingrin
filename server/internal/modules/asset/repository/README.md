# asset/repository

asset 模块 repository 规范：

- 每个资源统一拆分为三类文件：
  - `<resource>.go`：仓储结构体、构造函数、筛选映射/公共类型。
  - `<resource>_query.go`：查询方法（`Find/Get/List/Count/Stream/Scan`）。
  - `<resource>_command.go`：写操作方法（`Create/Update/Delete/Bulk*`）。
- 当前资源包括：`website / subdomain / endpoint / directory / host_port / screenshot`。

约束：

- 禁止使用 `*_mutation.go` 命名。
- 禁止使用泛名 `types.go`。
- `*_query.go` 不得出现写操作方法；`*_command.go` 不得出现查询方法。
