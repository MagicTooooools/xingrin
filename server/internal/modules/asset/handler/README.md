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
- `write.go`：写接口
- `export.go`：导出接口（仅有导出能力的资源）
