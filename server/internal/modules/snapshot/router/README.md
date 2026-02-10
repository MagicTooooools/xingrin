# snapshot/router

## 结构说明
- 公开入口：`snapshots.go`（`RegisterScanSnapshotRoutes`）
- 该路由组用于承接扫描快照的写入（bulk-upsert/bulk-create）与查询导出。

## 路由矩阵（挂载后实际前缀为 `/api`）

### 扫描维度快照（`/api/scans/:id/*`）
- `POST /api/scans/:id/websites/bulk-upsert`
- `POST /api/scans/:id/subdomains/bulk-upsert`
- `POST /api/scans/:id/endpoints/bulk-upsert`
- `POST /api/scans/:id/directories/bulk-upsert`
- `POST /api/scans/:id/host-ports/bulk-upsert`
- `POST /api/scans/:id/screenshots/bulk-upsert`
- `POST /api/scans/:id/vulnerabilities/bulk-create`
- 以及对应 `GET` 列表/导出路由

### 全局漏洞快照查询
- `GET /api/vulnerability-snapshots`
- `GET /api/vulnerability-snapshots/:id`

## 与 asset upsert 对应关系

`asset/router` 下的目标维度 upsert 路由统一委托到 `snapshot` 模块 handler：

- `/api/targets/:id/websites/bulk-upsert` -> `WebsiteSnapshotHandler.BulkUpsert`
- `/api/targets/:id/endpoints/bulk-upsert` -> `EndpointSnapshotHandler.BulkUpsert`
- `/api/targets/:id/directories/bulk-upsert` -> `DirectorySnapshotHandler.BulkUpsert`
- `/api/targets/:id/host-ports/bulk-upsert` -> `HostPortSnapshotHandler.BulkUpsert`
- `/api/targets/:id/screenshots/bulk-upsert` -> `ScreenshotSnapshotHandler.BulkUpsert`

## worker 写入入口（同一套快照 handler）
- `POST /api/worker/scans/:id/subdomains/bulk-upsert`
- `POST /api/worker/scans/:id/websites/bulk-upsert`
- `POST /api/worker/scans/:id/endpoints/bulk-upsert`

以上 worker 路由与 scan 快照路由复用同一批 snapshot handler，保证写入语义一致。
