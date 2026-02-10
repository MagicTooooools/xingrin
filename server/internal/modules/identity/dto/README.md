# identity/dto

identity 模块 DTO 规范：

- `*_dto.go`：按资源拆分放置 identity 业务 DTO（organization/user/target）。
- `common_http.go`：仅作为薄适配层，重导出 `server/internal/modules/httpdto` 的共享 HTTP DTO 能力（绑定、分页、统一错误响应）。

约束：

- 不得在 DTO 文件（`dto.go` 或 `*_dto.go`）中复用 `server/internal/dto` 的业务 DTO 别名。
- 不得 import 其他业务模块的 `dto`。
