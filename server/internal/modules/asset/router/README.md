# asset/router

## 结构说明
- 公开入口：`routes.go`（`RegisterAssetRoutes`）
- 健康检查：`health.go`（`RegisterHealthRoutes`）
- 其余文件为包内私有子注册函数（`register*Routes`）

## 文件
- `assets.go`：website/subdomain/directory 路由
- `endpoints.go`：endpoint 路由
- `host_ports.go`：host-port 路由
- `screenshots.go`：screenshot 路由
- `public.go`：无需鉴权的公共路由

## 约束
- 仅保留一个模块公开入口函数（health 这类跨模块启动入口除外）
- 子路由注册函数保持小写私有，避免对外暴露实现细节
- `bulk-upsert` 路由统一委托给 `snapshot` 模块 handler：
  - `/targets/:id/websites/bulk-upsert` -> `WebsiteSnapshotHandler`
  - `/targets/:id/endpoints/bulk-upsert` -> `EndpointSnapshotHandler`
  - `/targets/:id/directories/bulk-upsert` -> `DirectorySnapshotHandler`
  - `/targets/:id/host-ports/bulk-upsert` -> `HostPortSnapshotHandler`
  - `/targets/:id/screenshots/bulk-upsert` -> `ScreenshotSnapshotHandler`
