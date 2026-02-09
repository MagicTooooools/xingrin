# asset/dto

asset 模块 DTO 规范：

- `*_models.go`：按资源拆分放置 asset 业务 DTO（subdomain/website/endpoint/directory/host-port/screenshot）。
- `common_http.go`：仅作为薄适配层，重导出 `server/internal/modules/httpdto` 的共享 HTTP DTO 能力（绑定、分页、统一错误响应）。

约束：

- 不得在 DTO 模型文件（`models.go` 或 `*_models.go`）中复用 `server/internal/dto` 的业务 DTO 别名。
- 不得 import 其他模块的 `dto`（如 `security/snapshot`）。
- 跨模块数据转换放在 service/application 层，不放在 DTO 层。
