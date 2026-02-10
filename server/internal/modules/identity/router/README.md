# identity/router

## 结构说明
- 公开入口：`identity_module_routes.go`（`RegisterIdentityRoutes`）
- 其余文件为包内私有子注册函数（`register*Routes`）

## 文件
- `auth_routes.go`：auth 相关路由
- `user_routes.go`：user 相关路由
- `organization_routes.go`：organization 相关路由

## 约束
- 仅暴露模块入口函数
- 子注册函数只在本包内组合使用
