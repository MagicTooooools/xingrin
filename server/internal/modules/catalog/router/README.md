# catalog/router

## 结构说明
- 公开入口：`routes.go`（`RegisterCatalogRoutes`）
- 其余文件为包内私有子注册函数（`register*Routes`）

## 文件
- `targets.go`：target 路由
- `engines.go`：engine 路由
- `presets.go`：preset 路由
- `wordlists.go`：wordlist 路由
- `worker.go`：worker 路由

## 约束
- 外部只调用 `RegisterCatalogRoutes`
- 子路由按资源拆分，函数保持私有
