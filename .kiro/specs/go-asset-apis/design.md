# Design Document: Go Asset APIs

## Overview

本设计文档描述了 Go 后端资产 API 的架构设计，包括 Subdomain、Endpoint 和 Directory 三类资产的 CRUD 操作。设计遵循现有 Go 后端的分层架构模式（Handler → Service → Repository），确保与前端 API 的兼容性。

## Architecture

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer (Gin)                        │
├─────────────────────────────────────────────────────────────┤
│  SubdomainHandler  │  EndpointHandler  │  DirectoryHandler  │
├─────────────────────────────────────────────────────────────┤
│  SubdomainService  │  EndpointService  │  DirectoryService  │
├─────────────────────────────────────────────────────────────┤
│  SubdomainRepo     │  EndpointRepo     │  DirectoryRepo     │
├─────────────────────────────────────────────────────────────┤
│                      PostgreSQL (GORM)                       │
└─────────────────────────────────────────────────────────────┘
```

### API 路由设计

```
# Subdomain APIs
GET    /api/targets/:id/subdomains           # 列表查询
POST   /api/targets/:id/subdomains/bulk-create  # 批量创建
GET    /api/targets/:id/subdomains/export    # 导出
POST   /api/assets/subdomains/bulk-delete    # 批量删除

# Endpoint APIs
GET    /api/targets/:id/endpoints            # 列表查询
POST   /api/targets/:id/endpoints/bulk-create   # 批量创建
GET    /api/targets/:id/endpoints/export     # 导出
GET    /api/endpoints/:id                    # 详情查询
DELETE /api/endpoints/:id                    # 单个删除
POST   /api/assets/endpoints/bulk-delete     # 批量删除

# Directory APIs
GET    /api/targets/:id/directories          # 列表查询
POST   /api/targets/:id/directories/bulk-create # 批量创建
GET    /api/targets/:id/directories/export   # 导出
POST   /api/assets/directories/bulk-delete   # 批量删除
```

## Components and Interfaces

### Handler Layer

```go
// SubdomainHandler handles subdomain HTTP endpoints
type SubdomainHandler struct {
    svc *service.SubdomainService
}

func (h *SubdomainHandler) List(c *gin.Context)       // GET /targets/:id/subdomains
func (h *SubdomainHandler) BulkCreate(c *gin.Context) // POST /targets/:id/subdomains/bulk-create
func (h *SubdomainHandler) Export(c *gin.Context)     // GET /targets/:id/subdomains/export
func (h *SubdomainHandler) BulkDelete(c *gin.Context) // POST /assets/subdomains/bulk-delete

// EndpointHandler handles endpoint HTTP endpoints
type EndpointHandler struct {
    svc *service.EndpointService
}

func (h *EndpointHandler) List(c *gin.Context)       // GET /targets/:id/endpoints
func (h *EndpointHandler) GetByID(c *gin.Context)    // GET /endpoints/:id
func (h *EndpointHandler) BulkCreate(c *gin.Context) // POST /targets/:id/endpoints/bulk-create
func (h *EndpointHandler) Delete(c *gin.Context)     // DELETE /endpoints/:id
func (h *EndpointHandler) Export(c *gin.Context)     // GET /targets/:id/endpoints/export
func (h *EndpointHandler) BulkDelete(c *gin.Context) // POST /assets/endpoints/bulk-delete

// DirectoryHandler handles directory HTTP endpoints
type DirectoryHandler struct {
    svc *service.DirectoryService
}

func (h *DirectoryHandler) List(c *gin.Context)       // GET /targets/:id/directories
func (h *DirectoryHandler) BulkCreate(c *gin.Context) // POST /targets/:id/directories/bulk-create
func (h *DirectoryHandler) Export(c *gin.Context)     // GET /targets/:id/directories/export
func (h *DirectoryHandler) BulkDelete(c *gin.Context) // POST /assets/directories/bulk-delete
```

### Service Layer

```go
// SubdomainService handles subdomain business logic
type SubdomainService struct {
    repo       *repository.SubdomainRepository
    targetRepo *repository.TargetRepository
}

func (s *SubdomainService) ListByTarget(targetID int, query *dto.SubdomainListQuery) ([]model.Subdomain, int64, error)
func (s *SubdomainService) BulkCreate(targetID int, names []string) (int, error)
func (s *SubdomainService) BulkDelete(ids []int) (int64, error)
func (s *SubdomainService) StreamByTarget(targetID int) (*sql.Rows, error)
func (s *SubdomainService) CountByTarget(targetID int) (int64, error)

// EndpointService handles endpoint business logic
type EndpointService struct {
    repo       *repository.EndpointRepository
    targetRepo *repository.TargetRepository
}

func (s *EndpointService) ListByTarget(targetID int, query *dto.EndpointListQuery) ([]model.Endpoint, int64, error)
func (s *EndpointService) GetByID(id int) (*model.Endpoint, error)
func (s *EndpointService) BulkCreate(targetID int, urls []string) (int, error)
func (s *EndpointService) Delete(id int) error
func (s *EndpointService) BulkDelete(ids []int) (int64, error)
func (s *EndpointService) StreamByTarget(targetID int) (*sql.Rows, error)
func (s *EndpointService) CountByTarget(targetID int) (int64, error)

// DirectoryService handles directory business logic
type DirectoryService struct {
    repo       *repository.DirectoryRepository
    targetRepo *repository.TargetRepository
}

func (s *DirectoryService) ListByTarget(targetID int, query *dto.DirectoryListQuery) ([]model.Directory, int64, error)
func (s *DirectoryService) BulkCreate(targetID int, urls []string) (int, error)
func (s *DirectoryService) BulkDelete(ids []int) (int64, error)
func (s *DirectoryService) StreamByTarget(targetID int) (*sql.Rows, error)
func (s *DirectoryService) CountByTarget(targetID int) (int64, error)
```

### Repository Layer

```go
// SubdomainRepository handles subdomain database operations
type SubdomainRepository struct {
    db *gorm.DB
}

func (r *SubdomainRepository) FindByTargetID(targetID, page, pageSize int, filter string) ([]model.Subdomain, int64, error)
func (r *SubdomainRepository) BulkCreate(subdomains []model.Subdomain) (int, error)
func (r *SubdomainRepository) BulkDelete(ids []int) (int64, error)
func (r *SubdomainRepository) StreamByTargetID(targetID int) (*sql.Rows, error)
func (r *SubdomainRepository) CountByTargetID(targetID int) (int64, error)

// EndpointRepository handles endpoint database operations
type EndpointRepository struct {
    db *gorm.DB
}

func (r *EndpointRepository) FindByTargetID(targetID, page, pageSize int, filter string) ([]model.Endpoint, int64, error)
func (r *EndpointRepository) FindByID(id int) (*model.Endpoint, error)
func (r *EndpointRepository) BulkCreate(endpoints []model.Endpoint) (int, error)
func (r *EndpointRepository) Delete(id int) error
func (r *EndpointRepository) BulkDelete(ids []int) (int64, error)
func (r *EndpointRepository) StreamByTargetID(targetID int) (*sql.Rows, error)
func (r *EndpointRepository) CountByTargetID(targetID int) (int64, error)

// DirectoryRepository handles directory database operations
type DirectoryRepository struct {
    db *gorm.DB
}

func (r *DirectoryRepository) FindByTargetID(targetID, page, pageSize int, filter string) ([]model.Directory, int64, error)
func (r *DirectoryRepository) BulkCreate(directories []model.Directory) (int, error)
func (r *DirectoryRepository) BulkDelete(ids []int) (int64, error)
func (r *DirectoryRepository) StreamByTargetID(targetID int) (*sql.Rows, error)
func (r *DirectoryRepository) CountByTargetID(targetID int) (int64, error)
```

## Data Models

### 现有模型（已定义）

模型已在 `go-backend/internal/model/` 中定义，无需修改：

- `Subdomain`: id, target_id, name, created_at
- `Endpoint`: id, target_id, url, host, location, title, status_code, content_length, content_type, tech, webserver, response_body, response_headers, vhost, matched_gf_patterns, created_at
- `Directory`: id, target_id, url, status, content_length, words, lines, content_type, duration, created_at

## Asset-Target Matching Validation

资产创建时必须验证资产是否属于目标，确保数据一致性。

### 验证规则

| 资产类型 | Target 类型限制 | 匹配规则 |
|---------|----------------|---------|
| Subdomain | 仅 `domain` | subdomain == target 或 subdomain 以 `.target` 结尾 |
| Website | 无限制 | URL hostname 匹配 target |
| Endpoint | 无限制 | URL hostname 匹配 target |
| Directory | 无限制 | URL hostname 匹配 target |

### URL Hostname 匹配规则

根据 target 类型，URL hostname 的匹配方式不同：

| Target 类型 | 匹配规则 | 示例 |
|------------|---------|------|
| `domain` | hostname == target 或 hostname 以 `.target` 结尾 | target=`example.com` → `example.com` ✓, `api.example.com` ✓, `other.com` ✗ |
| `ip` | hostname == target | target=`192.168.1.1` → `192.168.1.1` ✓, `192.168.1.2` ✗ |
| `cidr` | hostname 是 IP 且在 CIDR 范围内 | target=`10.0.0.0/8` → `10.1.2.3` ✓, `192.168.1.1` ✗ |

### 验证函数设计

新增文件：`go-backend/internal/pkg/validator/target.go`

```go
package validator

import (
    "net"
    "net/url"
    "strings"
)

// IsURLMatchTarget checks if URL hostname matches target
// Returns true if the URL's hostname belongs to the target
func IsURLMatchTarget(urlStr, targetName, targetType string) bool {
    parsed, err := url.Parse(urlStr)
    if err != nil {
        return false
    }
    
    hostname := strings.ToLower(parsed.Hostname())
    if hostname == "" {
        return false
    }
    
    targetName = strings.ToLower(targetName)
    
    switch targetType {
    case "domain":
        // hostname equals target or ends with .target
        return hostname == targetName || strings.HasSuffix(hostname, "."+targetName)
    
    case "ip":
        // hostname must exactly equal target
        return hostname == targetName
    
    case "cidr":
        // hostname must be an IP within the CIDR range
        ip := net.ParseIP(hostname)
        if ip == nil {
            return false
        }
        _, network, err := net.ParseCIDR(targetName)
        if err != nil {
            return false
        }
        return network.Contains(ip)
    
    default:
        return false
    }
}

// IsSubdomainMatchTarget checks if subdomain belongs to target domain
// Returns true if subdomain equals target or ends with .target
func IsSubdomainMatchTarget(subdomain, targetDomain string) bool {
    subdomain = strings.ToLower(strings.TrimSpace(subdomain))
    targetDomain = strings.ToLower(strings.TrimSpace(targetDomain))
    
    if subdomain == "" || targetDomain == "" {
        return false
    }
    
    return subdomain == targetDomain || strings.HasSuffix(subdomain, "."+targetDomain)
}
```

### Service 层验证流程

批量创建时的验证流程：

```
输入数据 → 过滤空白 → 验证匹配 → 去重 → 批量插入
                         ↓
                    不匹配的静默跳过（计入 skipped）
```

#### Subdomain BulkCreate 验证

```go
func (s *SubdomainService) BulkCreate(targetID int, names []string) (int, error) {
    // 1. 获取 target
    target, err := s.targetRepo.FindByID(targetID)
    if err != nil {
        return 0, ErrTargetNotFound
    }
    
    // 2. 验证 target 类型必须是 domain
    if target.Type != "domain" {
        return 0, ErrInvalidTargetType
    }
    
    // 3. 过滤并验证
    var validSubdomains []model.Subdomain
    for _, name := range names {
        name = strings.TrimSpace(name)
        if name == "" {
            continue // 跳过空白
        }
        if !validator.IsSubdomainMatchTarget(name, target.Name) {
            continue // 跳过不匹配的
        }
        validSubdomains = append(validSubdomains, model.Subdomain{
            TargetID: targetID,
            Name:     name,
        })
    }
    
    // 4. 批量插入（去重由数据库 ON CONFLICT 处理）
    return s.repo.BulkCreate(validSubdomains)
}
```

#### Website/Endpoint/Directory BulkCreate 验证

```go
func (s *EndpointService) BulkCreate(targetID int, urls []string) (int, error) {
    // 1. 获取 target
    target, err := s.targetRepo.FindByID(targetID)
    if err != nil {
        return 0, ErrTargetNotFound
    }
    
    // 2. 过滤并验证
    var validEndpoints []model.Endpoint
    for _, u := range urls {
        u = strings.TrimSpace(u)
        if u == "" {
            continue // 跳过空白
        }
        if !validator.IsURLMatchTarget(u, target.Name, target.Type) {
            continue // 跳过不匹配的
        }
        validEndpoints = append(validEndpoints, model.Endpoint{
            TargetID: targetID,
            URL:      u,
            Host:     extractHostFromURL(u),
        })
    }
    
    // 3. 批量插入
    return s.repo.BulkCreate(validEndpoints)
}
```

### 新增错误类型

```go
var (
    ErrInvalidTargetType = errors.New("invalid target type: subdomain can only be created for domain-type targets")
)
```

### DTO 定义

```go
// SubdomainListQuery represents subdomain list query parameters
type SubdomainListQuery struct {
    PaginationQuery
    Filter string `form:"filter"`
}

// SubdomainResponse represents subdomain response
type SubdomainResponse struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}

// BulkCreateSubdomainsRequest represents bulk create subdomains request
type BulkCreateSubdomainsRequest struct {
    Subdomains []string `json:"subdomains" binding:"required,min=1,max=10000"`
}

// BulkCreateSubdomainsResponse represents bulk create subdomains response
type BulkCreateSubdomainsResponse struct {
    CreatedCount int `json:"createdCount"`
}

// EndpointListQuery represents endpoint list query parameters
type EndpointListQuery struct {
    PaginationQuery
    Filter string `form:"filter"`
}

// EndpointResponse represents endpoint response
type EndpointResponse struct {
    ID                int       `json:"id"`
    URL               string    `json:"url"`
    Host              string    `json:"host"`
    Location          string    `json:"location"`
    Title             string    `json:"title"`
    Webserver         string    `json:"webserver"`
    ContentType       string    `json:"contentType"`
    StatusCode        *int      `json:"statusCode"`
    ContentLength     *int      `json:"contentLength"`
    ResponseBody      string    `json:"responseBody"`
    Tech              []string  `json:"tech"`
    Vhost             *bool     `json:"vhost"`
    MatchedGFPatterns []string  `json:"matchedGfPatterns"`
    ResponseHeaders   string    `json:"responseHeaders"`
    CreatedAt         time.Time `json:"createdAt"`
}

// BulkCreateEndpointsRequest represents bulk create endpoints request
type BulkCreateEndpointsRequest struct {
    URLs []string `json:"urls" binding:"required,min=1,max=10000"`
}

// BulkCreateEndpointsResponse represents bulk create endpoints response
type BulkCreateEndpointsResponse struct {
    CreatedCount int `json:"createdCount"`
}

// DirectoryListQuery represents directory list query parameters
type DirectoryListQuery struct {
    PaginationQuery
    Filter string `form:"filter"`
}

// DirectoryResponse represents directory response
type DirectoryResponse struct {
    ID            int       `json:"id"`
    URL           string    `json:"url"`
    Status        *int      `json:"status"`
    ContentLength *int64    `json:"contentLength"`
    Words         *int      `json:"words"`
    Lines         *int      `json:"lines"`
    ContentType   string    `json:"contentType"`
    Duration      *int64    `json:"duration"`
    CreatedAt     time.Time `json:"createdAt"`
}

// BulkCreateDirectoriesRequest represents bulk create directories request
type BulkCreateDirectoriesRequest struct {
    URLs []string `json:"urls" binding:"required,min=1,max=10000"`
}

// BulkCreateDirectoriesResponse represents bulk create directories response
type BulkCreateDirectoriesResponse struct {
    CreatedCount int `json:"createdCount"`
}

// BulkDeleteRequest represents bulk delete request (shared)
type BulkDeleteRequest struct {
    IDs []int `json:"ids" binding:"required,min=1,max=10000"`
}

// BulkDeleteResponse represents bulk delete response (shared)
type BulkDeleteResponse struct {
    DeletedCount int64 `json:"deletedCount"`
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: List returns only items belonging to target

*For any* target with associated assets (subdomains/endpoints/directories), when listing assets for that target, all returned items SHALL have target_id equal to the requested target ID.

**Validates: Requirements 1.1, 5.1, 11.1**

### Property 2: Pagination respects page and pageSize parameters

*For any* valid page and pageSize parameters, the returned results SHALL contain at most pageSize items, and the page number SHALL match the requested page.

**Validates: Requirements 1.2, 5.2, 11.2**

### Property 3: Filter returns only matching items

*For any* filter string and list of assets, all returned items SHALL contain the filter text in the appropriate field(s) (name for subdomains, url/host/title for endpoints, url for directories).

**Validates: Requirements 1.3, 5.3, 11.3**

### Property 4: Results are sorted by createdAt DESC

*For any* list of returned assets, the createdAt timestamps SHALL be in descending order (newest first).

**Validates: Requirements 1.5, 5.5, 11.5**

### Property 5: Bulk create is idempotent

*For any* list of asset names/URLs, calling bulk create twice with the same data SHALL result in the same final state (no duplicates created).

**Validates: Requirements 2.2, 7.2, 12.2**

### Property 6: Bulk create count matches actual created records

*For any* bulk create operation, the returned createdCount SHALL equal the number of new records actually inserted into the database.

**Validates: Requirements 2.3, 7.3, 12.3**

### Property 7: Bulk delete removes specified items

*For any* list of valid asset IDs, after bulk delete, none of those IDs SHALL exist in the database.

**Validates: Requirements 3.1, 9.1, 13.1**

### Property 8: Bulk delete count matches actual deleted records

*For any* bulk delete operation, the returned deletedCount SHALL equal the number of records actually removed from the database.

**Validates: Requirements 3.2, 9.2, 13.2**

### Property 9: Export content matches database records

*For any* target with assets, the exported file content SHALL contain exactly the same assets as querying the database directly.

**Validates: Requirements 4.1, 10.1, 14.1**

### Property 10: Pagination response format is correct

*For any* paginated response, totalPages SHALL equal ceil(total / pageSize), and results SHALL be an empty array (not null) when no items exist.

**Validates: Requirements 15.1, 15.3, 15.4**

## Error Handling

### 错误类型定义

```go
var (
    ErrSubdomainNotFound = errors.New("subdomain not found")
    ErrEndpointNotFound  = errors.New("endpoint not found")
    ErrDirectoryNotFound = errors.New("directory not found")
    ErrTargetNotFound    = errors.New("target not found")  // 已存在
)
```

### HTTP 错误响应

| 场景 | HTTP Status | 响应格式 |
|------|-------------|----------|
| Target 不存在 | 404 | `{"error": "Target not found"}` |
| Asset 不存在 | 404 | `{"error": "Subdomain/Endpoint/Directory not found"}` |
| 请求参数无效 | 400 | `{"error": "Invalid request", "details": [...]}` |
| 服务器内部错误 | 500 | `{"error": "Internal server error"}` |

### 输入验证

- 批量创建：最大 10000 条记录
- 批量删除：最大 10000 个 ID
- 分页：pageSize 最大 1000
- 空字符串和纯空白字符串在批量创建时静默跳过

## Seed Data Generation

种子数据生成文件 `go-backend/cmd/seed/main.go` 需要更新以支持新的资产类型。

### 新增生成函数

```go
// createSubdomains creates subdomains for domain-type targets
func createSubdomains(db *gorm.DB, targetIDs []int, subdomainsPerTarget int) error

// createEndpoints creates endpoints for targets
func createEndpoints(db *gorm.DB, targetIDs []int, endpointsPerTarget int) error

// createDirectories creates directories for targets
func createDirectories(db *gorm.DB, targetIDs []int, directoriesPerTarget int) error
```

### 数据生成策略

| 资产类型 | 每个 Target 数量 | 数据特征 |
|---------|-----------------|---------|
| Subdomain | 20 | 仅为 domain 类型 target 生成，格式：`{prefix}.{target_domain}` |
| Endpoint | 20 | 包含 URL、状态码、技术栈等 HTTP 元数据 |
| Directory | 20 | 包含 URL、状态码、内容长度等目录扫描结果 |

### clearData 更新

```go
func clearData(db *gorm.DB) error {
    tables := []string{
        "directory",      // 新增
        "endpoint",       // 新增
        "subdomain",      // 新增
        "website",
        "organization_target",
        "target",
        "organization",
    }
    // ...
}
```

### 命令行参数

```bash
# 默认生成
go run cmd/seed/main.go

# 清除后重新生成
go run cmd/seed/main.go -clear

# 自定义数量
go run cmd/seed/main.go -orgs=50
```

## Python 后端差异分析

对比 Python 后端实现，以下是需要注意的差异：

### 1. Subdomain 批量创建响应

**Python 后端**返回详细统计：
```json
{
  "createdCount": 10,
  "skippedCount": 2,
  "invalidCount": 1,
  "mismatchedCount": 1,
  "totalReceived": 14
}
```

**Go 设计**简化为：
```json
{
  "createdCount": 10
}
```

**决策**：Go 版本简化响应，因为前端主要只使用 `createdCount`。如需详细统计，可后续扩展。

### 3. 导出格式

**Python 后端**：CSV 格式，导出所有数据库字段

**Go 设计**：CSV 格式，导出模型的所有字段（与数据库表结构一致）

**导出原则**：模型有多少字段就导出多少，不做字段筛选。

### 4. 批量删除路由

**Python 后端**：`POST /api/assets/subdomains/bulk-delete/`

**Go 设计**：`POST /api/subdomains/bulk-delete/`（遵循 go-backend-conventions.md，不加 assets 前缀）

**决策**：Go 版本使用更简洁的路由，前端需要相应调整。

### 5. 验证逻辑（已对齐）

**Python 后端**：
- Subdomain: 验证 target.type == "domain" + 域名后缀匹配
- Website/Endpoint/Directory: 验证 URL hostname 匹配 target

**Go 设计**：完全对齐 Python 行为，详见 "Asset-Target Matching Validation" 章节。

## Testing Strategy

### 索引覆盖分析

所有 filter 查询字段都已有索引覆盖：

| 资产类型 | Filter 字段 | 索引名称 | 索引类型 |
|---------|------------|---------|---------|
| Subdomain | `name` | `idx_subdomain_name` | B-tree |
| Endpoint | `url` | `idx_endpoint_url` | B-tree |
| Endpoint | `host` | `idx_endpoint_host` | B-tree |
| Endpoint | `title` | `idx_endpoint_title` | B-tree |
| Endpoint | `status_code` | `idx_endpoint_status_code` | B-tree |
| Endpoint | `tech` | `idx_endpoint_tech_gin` | GIN |
| Directory | `url` | `idx_directory_url` | B-tree |
| Directory | `status` | `idx_directory_status` | B-tree |

其他查询字段索引：
- `target_id`: 各模型都有 `idx_xxx_target` 索引
- `created_at`: 各模型都有 `idx_xxx_created_at` 索引（用于排序）
- `matched_gf_patterns`: `idx_endpoint_matched_gf_patterns_gin` (GIN)

**结论：不需要添加新索引。**

### 测试框架

- 单元测试：Go 标准库 `testing`
- HTTP 测试：`net/http/httptest`
- 属性测试：`github.com/leanovate/gopter`
- Mock：`github.com/stretchr/testify/mock`

### 测试层次

1. **Repository 层测试**：使用真实数据库（测试容器）验证 SQL 查询
2. **Service 层测试**：Mock Repository，验证业务逻辑
3. **Handler 层测试**：使用 httptest，验证 HTTP 接口

### 属性测试配置

- 每个属性测试运行 100 次迭代
- 使用 gopter 生成随机测试数据
- 测试标签格式：`Feature: go-asset-apis, Property N: {property_text}`

### 单元测试覆盖

- 正常流程测试
- 边界条件测试（空列表、最大分页等）
- 错误处理测试（404、400 等）
