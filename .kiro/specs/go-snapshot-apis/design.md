# 设计文档

## 概述

本文档描述 Go 后端网站快照（WebsiteSnapshot）API 的技术设计。该 API 作为扫描结果的写入入口，同时提供快照数据的查询和导出功能。

核心设计原则：
- **单一写入入口**：Worker 通过一个接口提交扫描结果，内部自动同步到快照表和资产表
- **Service 层解耦**：WebsiteSnapshotService 独立于 WebsiteService，通过组合调用实现同步
- **复用现有代码**：复用已有的 WebsiteService.BulkUpsert 方法写入资产表

## 架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Layer                               │
├─────────────────────────────────────────────────────────────────┤
│  WebsiteSnapshotHandler                                          │
│  - BulkUpsert()  POST /scans/{id}/websites/bulk-upsert          │
│  - List()        GET  /scans/{id}/websites/                      │
│  - Export()      GET  /scans/{id}/websites/export/               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                        Service Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  WebsiteSnapshotService                                          │
│  - SaveAndSync()      写入快照 + 同步资产                         │
│  - ListByScan()       查询快照                                    │
│  - StreamByScan()     流式导出                                    │
│                              │                                   │
│                              ↓ 调用                              │
│  WebsiteService（已有）                                           │
│  - BulkUpsert()       写入资产表                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       Repository Layer                           │
├─────────────────────────────────────────────────────────────────┤
│  WebsiteSnapshotRepository（新）    WebsiteRepository（已有）     │
│  - BulkCreate()                    - BulkUpsert()               │
│  - FindByScanID()                                                │
│  - StreamByScanID()                                              │
│  - CountByScanID()                                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                        Database Layer                            │
├─────────────────────────────────────────────────────────────────┤
│  website_snapshot 表              website 表                     │
│  - scan_id (FK)                   - target_id (FK)              │
│  - url (unique with scan_id)      - url (unique with target_id) │
│  - host, title, status_code...    - host, title, status_code... │
└─────────────────────────────────────────────────────────────────┘
```

## 组件和接口

### Handler: WebsiteSnapshotHandler

```go
// go-backend/internal/handler/website_snapshot.go

type WebsiteSnapshotHandler struct {
    svc *service.WebsiteSnapshotService
}

// BulkUpsert 批量写入网站快照（扫描结果导入）
// POST /api/scans/:scan_id/websites/bulk-upsert
func (h *WebsiteSnapshotHandler) BulkUpsert(c *gin.Context)

// List 查询网站快照列表
// GET /api/scans/:scan_id/websites/
func (h *WebsiteSnapshotHandler) List(c *gin.Context)

// Export 导出网站快照为 CSV
// GET /api/scans/:scan_id/websites/export/
func (h *WebsiteSnapshotHandler) Export(c *gin.Context)
```

### Service: WebsiteSnapshotService

```go
// go-backend/internal/service/website_snapshot.go

type WebsiteSnapshotService struct {
    snapshotRepo   *repository.WebsiteSnapshotRepository
    scanRepo       *repository.ScanRepository
    websiteService *WebsiteService  // 复用已有的资产 Service
}

// SaveAndSync 保存快照并同步到资产表
// 1. 验证 Scan 存在且未被软删除
// 2. 写入 website_snapshot 表
// 3. 调用 WebsiteService.BulkUpsert 写入 website 表
func (s *WebsiteSnapshotService) SaveAndSync(scanID int, items []dto.WebsiteSnapshotItem) (int64, error)

// ListByScan 查询指定扫描的快照列表
func (s *WebsiteSnapshotService) ListByScan(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error)

// StreamByScan 流式获取快照数据（用于 CSV 导出）
func (s *WebsiteSnapshotService) StreamByScan(scanID int) (*sql.Rows, error)

// CountByScan 获取快照数量
func (s *WebsiteSnapshotService) CountByScan(scanID int) (int64, error)
```

### Repository: WebsiteSnapshotRepository

```go
// go-backend/internal/repository/website_snapshot.go

type WebsiteSnapshotRepository struct {
    db *gorm.DB
}

// BulkCreate 批量创建快照（ON CONFLICT DO NOTHING）
func (r *WebsiteSnapshotRepository) BulkCreate(snapshots []model.WebsiteSnapshot) (int64, error)

// FindByScanID 查询指定扫描的快照（支持分页、过滤、排序）
func (r *WebsiteSnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter, ordering string) ([]model.WebsiteSnapshot, int64, error)

// StreamByScanID 流式获取快照数据
func (r *WebsiteSnapshotRepository) StreamByScanID(scanID int) (*sql.Rows, error)

// CountByScanID 获取快照数量
func (r *WebsiteSnapshotRepository) CountByScanID(scanID int) (int64, error)

// ScanRow 扫描单行数据
func (r *WebsiteSnapshotRepository) ScanRow(rows *sql.Rows) (*model.WebsiteSnapshot, error)
```

## 数据模型

### 请求 DTO

```go
// go-backend/internal/dto/website_snapshot.go

// WebsiteSnapshotItem 单个网站快照数据
type WebsiteSnapshotItem struct {
    URL             string   `json:"url" binding:"required,url"`
    Host            string   `json:"host"`
    Title           string   `json:"title"`
    StatusCode      *int     `json:"statusCode"`
    ContentLength   *int64   `json:"contentLength"`
    Location        string   `json:"location"`
    Webserver       string   `json:"webserver"`
    ContentType     string   `json:"contentType"`
    Tech            []string `json:"tech"`
    ResponseBody    string   `json:"responseBody"`
    Vhost           *bool    `json:"vhost"`
    ResponseHeaders string   `json:"responseHeaders"`
}

// BulkUpsertWebsiteSnapshotsRequest 批量写入请求
type BulkUpsertWebsiteSnapshotsRequest struct {
    TargetID int                   `json:"targetId" binding:"required"`
    Websites []WebsiteSnapshotItem `json:"websites" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertWebsiteSnapshotsResponse 批量写入响应
type BulkUpsertWebsiteSnapshotsResponse struct {
    SnapshotCount int `json:"snapshotCount"`
    AssetCount    int `json:"assetCount"`
}

// WebsiteSnapshotListQuery 列表查询参数
type WebsiteSnapshotListQuery struct {
    PaginationQuery
    Filter   string `form:"filter"`
    Ordering string `form:"ordering"`
}

// WebsiteSnapshotResponse 快照响应
type WebsiteSnapshotResponse struct {
    ID              int       `json:"id"`
    ScanID          int       `json:"scanId"`
    URL             string    `json:"url"`
    Host            string    `json:"host"`
    Title           string    `json:"title"`
    StatusCode      *int      `json:"statusCode"`
    ContentLength   *int64    `json:"contentLength"`
    Location        string    `json:"location"`
    Webserver       string    `json:"webserver"`
    ContentType     string    `json:"contentType"`
    Tech            []string  `json:"tech"`
    ResponseBody    string    `json:"responseBody"`
    Vhost           *bool     `json:"vhost"`
    ResponseHeaders string    `json:"responseHeaders"`
    CreatedAt       time.Time `json:"createdAt"`
}
```

### 过滤字段映射

```go
var WebsiteSnapshotFilterMapping = scope.FilterMapping{
    "url":       {Column: "url"},
    "host":      {Column: "host"},
    "title":     {Column: "title"},
    "status":    {Column: "status_code", IsNumeric: true},
    "webserver": {Column: "webserver"},
    "tech":      {Column: "tech", IsArray: true},
}
```

## 正确性属性

*正确性属性是系统在所有有效执行中应保持为真的特征或行为——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 快照和资产同步写入

*For any* 有效的网站快照数据，通过 bulk-upsert 接口写入后，数据应同时存在于 website_snapshot 表和 website 表中，且字段值一致（除了 scan_id/target_id 的差异）。

**Validates: Requirements 1.1, 1.2**

### Property 2: 快照去重

*For any* 包含重复 URL 的请求，写入后 website_snapshot 表中同一 scan_id 下的 URL 应唯一，重复项被忽略。

**Validates: Requirements 1.4**

### Property 3: 资产 Upsert 保留 created_at

*For any* 已存在于 website 表的记录，通过快照同步更新后，其 created_at 字段应保持不变。

**Validates: Requirements 1.5**

### Property 4: 分页正确性

*For any* 分页查询，返回的记录数应不超过 pageSize，且 total 应等于该 scan 下的所有快照数量。

**Validates: Requirements 2.3, 2.6**

### Property 5: 过滤正确性

*For any* 带 filter 参数的查询，返回的所有记录应满足过滤条件（文本字段模糊匹配，数字字段精确匹配）。

**Validates: Requirements 3.1, 3.2, 3.4**

### Property 6: 排序正确性

*For any* 带 ordering 参数的查询，返回的记录应按指定字段和方向排序。

**Validates: Requirements 4.1, 4.2, 4.3**

### Property 7: CSV 导出完整性

*For any* CSV 导出请求，导出的记录数应等于该 scan 下的所有快照数量，且每条记录包含所有必需字段。

**Validates: Requirements 5.1, 5.3, 5.4**

### Property 8: Scan 存在性验证

*For any* 快照请求（读或写），如果 scan_id 不存在或已被软删除，应返回 404 错误。

**Validates: Requirements 7.1, 7.2, 7.3, 7.4**

## 错误处理

| 场景 | HTTP 状态码 | 错误码 | 错误信息 |
|------|------------|--------|----------|
| scan_id 无效（非数字） | 400 | BAD_REQUEST | Invalid scan ID |
| scan 不存在 | 404 | NOT_FOUND | Scan not found |
| scan 已软删除 | 404 | NOT_FOUND | Scan not found |
| 请求体格式错误 | 400 | VALIDATION_ERROR | Invalid request body |
| 请求体为空 | 400 | VALIDATION_ERROR | websites is required |
| 超过最大数量限制 | 400 | VALIDATION_ERROR | websites must have at most 5000 items |
| 数据库错误 | 500 | SERVER_ERROR | Failed to save snapshots |

## 测试策略

### 单元测试

- Handler 层：测试请求解析、参数验证、错误响应
- Service 层：测试业务逻辑、Scan 验证、同步调用
- Repository 层：测试 SQL 查询、批量操作、去重逻辑

### 属性测试

使用 `github.com/leanovate/gopter` 进行属性测试：

1. **Property 1 测试**：生成随机网站数据，调用 SaveAndSync，验证两张表数据一致
2. **Property 2 测试**：生成包含重复 URL 的数据，验证去重行为
3. **Property 3 测试**：先插入资产，再通过快照同步更新，验证 created_at 不变
4. **Property 4 测试**：生成随机分页参数，验证返回结果符合分页规则
5. **Property 5 测试**：生成随机过滤条件，验证返回结果满足条件
6. **Property 6 测试**：生成随机排序参数，验证返回结果有序
7. **Property 7 测试**：生成随机快照数据，导出 CSV，验证完整性
8. **Property 8 测试**：使用不存在/已删除的 scan_id，验证返回 404

### 测试配置

- 每个属性测试运行 100 次迭代
- 使用测试数据库，每次测试前清理数据
- 测试标签格式：`Feature: go-snapshot-apis, Property N: {property_text}`
