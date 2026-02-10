# scan/router

## 结构说明
- 公开入口：`scan_module_routes.go`（`RegisterScanRoutes`）
- worker 入口：`worker_scan_routes.go`（`RegisterWorkerScanRoutes`）
- 其余文件为包内私有子注册函数

## 文件
- `scan_routes.go`：scan 资源路由
- `scan_log_routes.go`：scan-log 资源路由

## 约束
- 普通业务路由统一通过 `RegisterScanRoutes` 组合注册
- worker 路由单独入口，避免与业务鉴权组混用
