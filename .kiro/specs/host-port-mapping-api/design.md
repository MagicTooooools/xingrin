# Design Document: Host Port Mapping API

## Overview

为 Go 后端实现 host_port_mapping 资产的 CRUD API。该 API 遵循项目现有的资产 API 设计模式（参考 website、subdomain），提供列表查询、CSV 导出、批量创建、批量 upsert 和批量删除功能。

## Architecture

```
HTTP Request → Handler → Service → Repository → Database
                 ↓          ↓           ↓
               DTO      Business     GORM Model
                        Logic
```

### API Routes

| Method | Route | Handler | Description |
|--------|-------|---------|-------------|
| GET | `/targets/:id/host-port-mappings` | List | 分页列表（按 IP 聚合） |
| GET | `/targets/:id/host-port-mappings/export` | Export | CSV 导出（原始格式） |
| POST | `/targets/:id/host-port-mappings/bulk-upsert` | BulkUpsert | 批量 upsert（扫描器写入） |
| POST | `/host-port-mappings/bulk-delete` | BulkDelete | 批量删除（按 IP 删除） |

### 响应格式说明

**List 接口**：返回按 IP 聚合的数据，与前端 `IPAddress` 类型匹配：
```typescript
// 前端期望的格式
interface IPAddress {
  ip: string        // IP 地址（唯一标识）
  hosts: string[]   // 关联的主机名列表
  ports: number[]   // 关联的端口列表
  createdAt: string // 最早创建时间
}
```

**Export 接口**：返回原始格式 CSV（每行一个 host+ip+port 组合）：
```csv
ip,host,port,created_at
192.168.1.1,example.com,80,2024-01-01 12:00:00
192.168.1.1,example.com,443,2024-01-01 12:00:00
```

**Bulk Delete 接口**：接收 IP 字符串列表（不是 ID），删除这些 IP 的所有映射记录

## Components and Interfaces

### Handler Layer

```go
// HostPortMappingHandler handles HTTP requests
type HostPortMappingHandler struct {
    svc *HostPortMappingService
}

func (h *HostPortMappingHandler) List(c *gin.Context)
func (h *HostPortMappingHandler) Export(c *gin.Context)
func (h *HostPortMappingHandler) BulkUpsert(c *gin.Context)
func (h *HostPortMappingHandler) BulkDelete(c *gin.Context)
```

### Service Layer

```go
// HostPortMappingService handles business logic
type HostPortMappingService struct {
    repo       *HostPortMappingRepository
    targetRepo *TargetRepository
}

// ListByTarget 返回按 IP 聚合的数据（与 Python 后端一致）
// 聚合逻辑：
// 1. 按 IP 分组，获取每个 IP 的最早 created_at
// 2. 对每个 IP，收集其所有 hosts 和 ports（去重排序）
func (s *HostPortMappingService) ListByTarget(targetID int, query *dto.HostPortMappingListQuery) ([]dto.HostPortMappingResponse, int64, error)

// StreamByTarget 流式返回原始数据用于 CSV 导出
func (s *HostPortMappingService) StreamByTarget(targetID int) (*sql.Rows, error)

// BulkUpsert 批量创建（忽略冲突）
// 使用 ON CONFLICT DO NOTHING，因为所有字段都在唯一约束中
func (s *HostPortMappingService) BulkUpsert(targetID int, items []dto.HostPortMappingItem) (int64, error)

// BulkDeleteByIPs 按 IP 列表删除所有相关映射
// 与 Python 后端一致：传入 IP 列表，删除这些 IP 的所有记录
func (s *HostPortMappingService) BulkDeleteByIPs(ips []string) (int64, error)
```

### Repository Layer

```go
// HostPortMappingRepository handles database operations
type HostPortMappingRepository struct {
    db *gorm.DB
}

// GetIPAggregation 获取按 IP 聚合的数据
// SQL: SELECT ip, MIN(created_at) FROM host_port_mapping WHERE target_id = ? GROUP BY ip ORDER BY MIN(created_at) DESC
func (r *HostPortMappingRepository) GetIPAggregation(targetID int, filter string) ([]IPAggregationRow, error)

// GetHostsAndPortsByIP 获取指定 IP 的所有 hosts 和 ports
func (r *HostPortMappingRepository) GetHostsAndPortsByIP(targetID int, ip string, filter string) (hosts []string, ports []int, error)

// StreamByTargetID 流式返回原始数据用于 CSV 导出
func (r *HostPortMappingRepository) StreamByTargetID(targetID int) (*sql.Rows, error)

// BulkUpsert 批量插入（忽略冲突）
// 使用 ON CONFLICT (target_id, host, ip, port) DO NOTHING
func (r *HostPortMappingRepository) BulkUpsert(mappings []model.HostPortMapping) (int64, error)

// DeleteByIPs 按 IP 列表删除
// SQL: DELETE FROM host_port_mapping WHERE ip IN (?)
func (r *HostPortMappingRepository) DeleteByIPs(ips []string) (int64, error)

// ScanRow 扫描单行数据
func (r *HostPortMappingRepository) ScanRow(rows *sql.Rows) (*model.HostPortMapping, error)
```

## Data Models

### Existing Model (host_port_mapping.go)

```go
type HostPortMapping struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    TargetID  int       `gorm:"column:target_id" json:"targetId"`
    Host      string    `gorm:"column:host" json:"host"`
    IP        string    `gorm:"column:ip;type:inet" json:"ip"`
    Port      int       `gorm:"column:port" json:"port"`
    CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}
```

### DTOs

```go
// Query
type HostPortMappingListQuery struct {
    PaginationQuery
    Filter string `form:"filter"`
}

// Response (聚合格式，按 IP 分组)
type HostPortMappingResponse struct {
    IP        string    `json:"ip"`        // IP 地址（唯一标识）
    Hosts     []string  `json:"hosts"`     // 关联的主机名列表
    Ports     []int     `json:"ports"`     // 关联的端口列表
    CreatedAt time.Time `json:"createdAt"` // 最早创建时间
}

// Paginated Response
type HostPortMappingListResponse struct {
    Results    []HostPortMappingResponse `json:"results"`
    Total      int64                     `json:"total"`
    Page       int                       `json:"page"`
    PageSize   int                       `json:"pageSize"`
    TotalPages int                       `json:"totalPages"`
}

// Request Item (for bulk upsert)
type HostPortMappingItem struct {
    Host string `json:"host" binding:"required"`
    IP   string `json:"ip" binding:"required,ip"`
    Port int    `json:"port" binding:"required,min=1,max=65535"`
}

// Bulk Upsert Request (for scanner import)
type BulkUpsertHostPortMappingsRequest struct {
    Mappings []HostPortMappingItem `json:"mappings" binding:"required,min=1,max=5000,dive"`
}

// Bulk Upsert Response
type BulkUpsertHostPortMappingsResponse struct {
    UpsertedCount int `json:"upsertedCount"`
}

// Bulk Delete Request (by IP list)
type BulkDeleteHostPortMappingsRequest struct {
    IPs []string `json:"ips" binding:"required,min=1,dive,ip"`
}

// Bulk Delete Response
type BulkDeleteHostPortMappingsResponse struct {
    DeletedCount int64 `json:"deletedCount"`
}
```

### Filter Mapping

```go
var HostPortMappingFilterMapping = scope.FilterMapping{
    "host": {Column: "host"},
    "ip":   {Column: "ip"},
    "port": {Column: "port"},
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Pagination returns correct subset

*For any* target with N mappings, requesting page P with pageSize S should return at most S items, and the total count should equal N.

**Validates: Requirements 1.1, 1.4**

### Property 2: Filter returns only matching results

*For any* filter query on host/ip/port, all returned mappings should match the filter criteria.

**Validates: Requirements 1.2**

### Property 3: Upsert creates new and preserves existing

*For any* set of mappings, upsert should create new records for non-existing combinations and preserve existing records.

**Validates: Requirements 4.1**

### Property 4: Bulk delete handles non-existent IPs gracefully

*For any* list of IPs (including non-existent ones), bulk delete should remove only existing records and return the count of actually deleted records.

**Validates: Requirements 4.1, 4.3**

## Error Handling

| Error | HTTP Status | Code |
|-------|-------------|------|
| Target not found | 404 | NOT_FOUND |
| Invalid request body | 400 | BAD_REQUEST |
| Invalid target ID | 400 | BAD_REQUEST |
| Internal error | 500 | INTERNAL_ERROR |

## Testing Strategy

### Unit Tests

- Handler: 测试请求解析和响应格式
- Service: 测试业务逻辑和错误处理
- Repository: 测试数据库操作

### Property-Based Tests

使用 `github.com/leanovate/gopter` 进行属性测试：

- Property 1: 分页正确性
- Property 2: 筛选正确性
- Property 3: Upsert 行为
- Property 4: 删除容错性

每个属性测试运行 100 次迭代。
