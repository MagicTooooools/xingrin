# catalog/dto

catalog 模块 DTO 规范：

- `*_models.go`：按资源拆分放置 catalog 业务 DTO（target/engine/wordlist/worker/preset）。
- `common_http.go`：仅作为薄适配层，重导出 `server/internal/modules/httpdto` 的共享 HTTP DTO 能力（绑定、分页、统一错误响应）。

约束：

- 不得在 DTO 模型文件（`models.go` 或 `*_models.go`）中复用 `server/internal/dto` 的业务 DTO 别名。
- 不得 import 其他业务模块的 `dto`。
- 与其他模块的字段映射通过 service/application 显式转换。
