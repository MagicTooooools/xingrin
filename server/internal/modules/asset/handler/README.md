# asset/handler

asset 模块 handler 层统一按“资源子目录”组织：

- `health.go`：健康检查（模块级）
- `website/`
- `endpoint/`
- `subdomain/`
- `directory/`
- `host_port/`
- `screenshot/`

各资源子目录遵循固定职责拆分：

- `handler.go`：Handler 结构体与构造函数
- `read.go`：读接口
- `write.go`：写接口（不含 `bulk-upsert`）
- `export.go`：导出接口（仅有导出能力的资源）

## 约定

- `bulk-upsert` 统一由 `snapshot` 模块 handler 承接，不在 `asset/handler/*/write.go` 内实现。
- `asset` 资源 handler 仅保留查询、导出、创建、删除等资产侧接口。
