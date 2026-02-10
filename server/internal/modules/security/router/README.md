# security/router

## 结构说明
- 公开入口：`routes.go`（`RegisterSecurityRoutes`）
- 资源子路由：`vulnerabilities.go`（私有 `registerVulnerabilityRoutes`）

## 约束
- 外部仅通过 `RegisterSecurityRoutes` 注册安全模块路由
- 资源级注册函数保持私有
