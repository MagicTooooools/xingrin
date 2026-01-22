---
inclusion: fileMatch
fileMatchPattern: "go-backend/**/*.go"
---

# Go Backend Code Conventions

## Language Requirements

All code in `go-backend/` must be written in English:
- Comments, error messages, log messages, API responses
- Variable names and function names

## Project Structure

```
go-backend/
├── cmd/server/main.go      # Application entrypoint
├── internal/
│   ├── config/             # Configuration loading
│   ├── database/           # Database connection
│   ├── model/              # Domain models (GORM)
│   ├── dto/                # Request/Response DTOs
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic layer
│   ├── handler/            # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── auth/               # Authentication (JWT)
│   └── pkg/                # Internal utilities
├── go.mod
└── Makefile
```

## File Naming

Use concise names without redundant suffixes:

```
✅ handler/target.go
✅ service/target.go
✅ repository/target.go

❌ handler/target_handler.go
❌ service/target_service.go
```

## Architecture Layers

```
HTTP Request → Handler → Service → Repository → Database
                 ↓          ↓           ↓
               DTO      Business     GORM Model
                        Logic
```

- **Handler**: Parse request, validate input, call service, return response
- **Service**: Business logic, orchestrate repositories, no HTTP awareness
- **Repository**: Database operations only, return models

## Error Handling

```go
// ✅ Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create target: %w", err)
}

// ✅ Define domain errors in service layer
var (
    ErrTargetNotFound = errors.New("target not found")
    ErrTargetExists   = errors.New("target already exists")
)

// ✅ Check specific errors
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrTargetNotFound
}
```

## Interface Design

```go
// ✅ Define interfaces where they are used (service layer)
type TargetRepository interface {
    Create(target *model.Target) error
    FindByID(id int) (*model.Target, error)
}

// ✅ Accept interfaces, return structs
func NewTargetService(repo TargetRepository) *TargetService {
    return &TargetService{repo: repo}
}
```

## Context Usage

```go
// ✅ Pass context as first parameter
func (r *TargetRepository) FindByID(ctx context.Context, id int) (*model.Target, error) {
    var target model.Target
    err := r.db.WithContext(ctx).First(&target, id).Error
    return &target, err
}

// ✅ Use context for cancellation and timeouts
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

## Dependency Injection

```go
// ✅ Constructor injection
func NewTargetService(repo *TargetRepository, orgRepo *OrganizationRepository) *TargetService {
    return &TargetService{
        repo:    repo,
        orgRepo: orgRepo,
    }
}

// ❌ Avoid global state
var globalDB *gorm.DB // Don't do this
```

## Testing

```go
// ✅ Table-driven tests
func TestDetectTargetType(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"domain", "example.com", "domain"},
        {"ip", "192.168.1.1", "ip"},
        {"cidr", "10.0.0.0/8", "cidr"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := detectTargetType(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

## GORM Model Conventions

```go
type Target struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"column:name;uniqueIndex" json:"name"`
    Type      string    `gorm:"column:type" json:"type"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// Use singular table names
func (Target) TableName() string {
    return "target"
}
```

## HTTP Handler Pattern

```go
func (h *TargetHandler) Create(c *gin.Context) {
    // 1. Parse and validate request
    var req dto.CreateTargetRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        dto.BadRequest(c, "invalid request body")
        return
    }
    
    // 2. Call service
    target, err := h.svc.Create(&req)
    if err != nil {
        // 3. Handle domain errors
        if errors.Is(err, service.ErrTargetExists) {
            dto.Conflict(c, err.Error())
            return
        }
        dto.InternalError(c, "failed to create target")
        return
    }
    
    // 4. Return response
    dto.Created(c, target)
}
```

## Concurrency Safety

```go
// ✅ Use channels for communication
results := make(chan Result, len(items))
for _, item := range items {
    go func(item Item) {
        results <- process(item)
    }(item)
}

// ✅ Use sync primitives when needed
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// ✅ Always close resources
defer rows.Close()
defer resp.Body.Close()
```

## Logging

```go
// ✅ Use structured logging (zap)
pkg.Info("target created",
    zap.Int("id", target.ID),
    zap.String("name", target.Name),
)

// ✅ Log errors with context
pkg.Error("failed to create target",
    zap.Error(err),
    zap.String("name", req.Name),
)
```

## Key Principles

1. **Simplicity** - Write clear, straightforward code
2. **Explicit** - No magic, explicit dependencies and error handling
3. **Testable** - Design for easy unit testing with interfaces
4. **Layered** - Clear separation between handler/service/repository
5. **Idiomatic** - Follow Go conventions and standard library patterns


## API Response Format (Industry Standard)

Follow Stripe/GitHub style - return data directly without wrapper.

### Success Response

```go
// Single resource - return directly
{"id": 1, "name": "example.com", "type": "domain"}

// List/Operation result - return directly
{"count": 3, "items": [...]}

// Use dto helpers
dto.OK(c, target)           // 200 - returns data directly
dto.Created(c, target)      // 201 - returns data directly
dto.NoContent(c)            // 204 - no body
```

### Error Response

```go
// Error format with code for i18n
{
    "error": {
        "code": "NOT_FOUND",
        "message": "target not found"  // Debug info, not for display
    }
}

// With field-level details
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid input data",
        "details": [
            {"field": "name", "message": "required"}
        ]
    }
}

// Use dto helpers
dto.BadRequest(c, "message")       // 400
dto.Unauthorized(c, "message")     // 401
dto.NotFound(c, "message")         // 404
dto.Conflict(c, "message")         // 409
dto.InternalError(c, "message")    // 500
dto.ValidationError(c, details)    // 400 with field details
```

## Pagination

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number (1-based) |
| pageSize | int | 20 | Items per page (max 100) |

### Paginated Response Format

```go
// Matches Python format for frontend compatibility
{
    "results": [...],      // Data array
    "total": 100,          // Total count
    "page": 1,             // Current page
    "pageSize": 20,        // Items per page
    "totalPages": 5        // Total pages
}
```

### Usage

```go
// In handler
func (h *TargetHandler) List(c *gin.Context) {
    var query dto.TargetListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        dto.BadRequest(c, "invalid query parameters")
        return
    }
    
    targets, total, err := h.svc.List(&query)
    if err != nil {
        dto.InternalError(c, "failed to list targets")
        return
    }
    
    // Use generic Paginated helper
    dto.Paginated(c, targets, total, query.GetPage(), query.GetPageSize())
}
```

### Empty Array Handling

`dto.Paginated` automatically converts `nil` slices to empty arrays `[]` in JSON output.

```go
// ✅ Handler can use simple var declaration
var resp []dto.TargetResponse
for _, t := range targets {
    resp = append(resp, toResponse(&t))
}
dto.Paginated(c, resp, total, page, pageSize)
// Output: {"results": [], ...}  (not null)

// ❌ For non-paginated responses, initialize explicitly
failedTargets := []dto.FailedTarget{}  // Not: var failedTargets []dto.FailedTarget
dto.Success(c, BatchResponse{FailedTargets: failedTargets})
```

**Rule**: Use `dto.Paginated` for list APIs - it handles nil → `[]` automatically. For other responses with array fields, initialize with `[]T{}` to avoid `null` in JSON.


## Filter (Unified Search Parameter)

All list APIs use a unified `filter` parameter for searching and filtering.

### Query Parameter

| Parameter | Type | Description |
|-----------|------|-------------|
| filter | string | Plain text or smart filter syntax |

### Filter Syntax

| Syntax | Description | Example |
|--------|-------------|---------|
| Plain text | Fuzzy search on default field | `filter=portal` |
| `field="value"` | Fuzzy match (ILIKE) | `filter=name="portal"` |
| `field=="value"` | Exact match | `filter=name=="example.com"` |
| `field!="value"` | Not equal | `filter=type!="ip"` |
| `\|\|` or `or` | OR logic | `filter=name="a" \|\| name="b"` |
| `&&` or `and` | AND logic | `filter=name="a" && type="domain"` |

### Implementation

```go
// 1. Define filter mapping in repository
var TargetFilterMapping = scope.FilterMapping{
    "name": {Column: "target.name", IsArray: false},
    "type": {Column: "target.type", IsArray: false},
}

// 2. Use WithFilterDefault in repository (with default field for plain text)
query = query.Scopes(scope.WithFilterDefault(filter, TargetFilterMapping, "name"))

// 3. DTO uses "filter" parameter
type TargetListQuery struct {
    PaginationQuery
    Filter string `form:"filter"`
    Type   string `form:"type"`
}
```

### Frontend Usage

```typescript
// Plain text search (searches default field)
const response = await api.get('/targets/', { params: { filter: 'portal' } })

// Smart filter syntax
const response = await api.get('/targets/', { params: { filter: 'name="portal"' } })
```


## Validation

Use `github.com/asaskevich/govalidator` for common validations. Don't write custom regex.

```go
import "github.com/asaskevich/govalidator"

// ✅ Use govalidator
govalidator.IsDNSName("example.com")     // Domain
govalidator.IsURL("https://example.com") // URL
govalidator.IsIP("192.168.1.1")          // IP

// ✅ Use standard library for IP/CIDR
net.ParseIP("192.168.1.1")
net.ParseCIDR("10.0.0.0/8")

// ❌ Don't write custom regex for domain/IP/URL validation
regexp.MustCompile(`^([a-zA-Z0-9]...`)
```


## Delete API Convention

### Single Delete (RESTful)

```
DELETE /api/{resource}/{id}/
Response: 204 No Content (no body)
```

```go
// Handler
func (h *TargetHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    
    if err := h.svc.Delete(id); err != nil {
        if errors.Is(err, service.ErrTargetNotFound) {
            dto.NotFound(c, "Target not found")
            return
        }
        dto.InternalError(c, "Failed to delete target")
        return
    }
    
    dto.NoContent(c)  // 204
}
```

### Bulk Delete

```
POST /api/{resource}/bulk-delete/
Request: {"ids": [1, 2, 3]}
Response: 200 {"deletedCount": 3}
```

```go
// Handler
func (h *TargetHandler) BulkDelete(c *gin.Context) {
    var req dto.BulkDeleteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        dto.BadRequest(c, "Invalid request body")
        return
    }
    
    deletedCount, err := h.svc.BulkDelete(req.IDs)
    if err != nil {
        dto.InternalError(c, "Failed to delete targets")
        return
    }
    
    dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}
```

### Frontend Integration

- Single delete: Frontend uses local data for success toast (name already known)
- Bulk delete: Frontend uses `deletedCount` from response for success toast


## Batch Operation Limits

| Operation | Max Items | Notes |
|-----------|-----------|-------|
| Batch Create | 5000 | `binding:"max=5000"` |
| Bulk Delete | No limit | But consider performance |

```go
// Example: Batch create with limit
type BatchCreateTargetRequest struct {
    Targets []TargetItem `json:"targets" binding:"required,min=1,max=5000,dive"`
}
```


## Request Binding (Industry Standard)

Use `dto.BindJSON`, `dto.BindQuery`, `dto.BindURI` helpers for automatic validation error handling.

```go
// ✅ Use dto.BindJSON - one line, auto error response
var req dto.CreateTargetRequest
if !dto.BindJSON(c, &req) {
    return
}

// ❌ Don't use c.ShouldBindJSON directly
if err := c.ShouldBindJSON(&req); err != nil {
    // Manual error handling...
}
```

Benefits:
- Automatic validation error translation (e.g., "targets must have at most 5000 items")
- Consistent error response format across all handlers
- Less boilerplate code


## Standard Library First

Prefer Go standard library over custom implementations or third-party packages.

```go
// ✅ Use encoding/csv for CSV operations
import "encoding/csv"
writer := csv.NewWriter(w)
writer.Write([]string{"a", "b", "c"})

// ✅ Use encoding/json for JSON
import "encoding/json"

// ✅ Use net/http for HTTP clients
import "net/http"

// ✅ Use time for time operations
import "time"

// ❌ Don't write custom CSV escaping
func escapeCSV(s string) string { ... }  // Use encoding/csv instead
```

Benefits:
- Well-tested and maintained
- Handles edge cases (escaping, encoding)
- Familiar to other Go developers


## File Export (CSV/Excel)

### Streaming Export Pattern

For large data exports, use streaming to avoid memory issues:

```go
func (h *Handler) Export(c *gin.Context) {
    // 1. Get streaming cursor
    rows, err := h.svc.Stream(id)
    
    // 2. Use csv.StreamCSV helper (no Content-Length for streaming)
    csv.StreamCSV(c, rows, headers, filename, mapper, 0)
}
```

### IMPORTANT: Don't Set Content-Length for Streaming

```go
// ❌ DON'T set Content-Length for streaming exports
// Estimated size may not match actual size, causing connection to close early
c.Header("Content-Length", strconv.Itoa(estimatedSize))

// ✅ DO use chunked transfer encoding (default when no Content-Length)
// Browser will show "unknown size" but download completes correctly
```

### UTF-8 BOM for Excel

Add UTF-8 BOM at the start for Excel to recognize Chinese characters:

```go
var UTF8BOM = []byte{0xEF, 0xBB, 0xBF}
c.Writer.Write(UTF8BOM)
```


## RESTful Route Design

### Nested Resources (belongs to parent)

Use nested routes when the child resource is always accessed in context of a parent:

```
GET    /targets/:id/websites          # List websites for a target
POST   /targets/:id/websites/bulk-create  # Create websites for a target
GET    /targets/:id/websites/export   # Export websites for a target
```

### Standalone Resources (independent operations)

Use standalone routes when the operation doesn't require parent context:

```
GET    /websites/:id           # Get single website by ID
DELETE /websites/:id           # Delete single website by ID
POST   /websites/bulk-delete   # Bulk delete by IDs (no target needed)
```

### Route Naming Rules

```
✅ /targets/:id/websites       # Nested under parent
✅ /websites/:id               # Standalone by ID
✅ /websites/bulk-delete       # Action on resource

❌ /assets/websites/bulk-delete  # Don't add unnecessary prefixes
❌ /api/v1/assets/targets/...    # Keep URLs short and clean
```


## Bulk Create vs Bulk Upsert

For asset tables (website, subdomain, endpoint, etc.), provide two separate interfaces:

### bulk-create (Frontend manual add)

```
POST /targets/:id/websites/bulk-create
Request: {"urls": ["https://..."]}
Behavior: ON CONFLICT DO NOTHING (only create new, ignore duplicates)
Use case: User manually adds assets from frontend
```

### bulk-upsert (Scanner import)

```
POST /targets/:id/websites/bulk-upsert
Request: {"websites": [{url, title, statusCode, tech, ...}]}
Behavior: ON CONFLICT DO UPDATE with COALESCE + array merge
Use case: Scanner imports assets with full data
```

### Target Ownership Validation

Both interfaces MUST validate that assets belong to the target:

```go
// Service layer - filter items that match target
for _, item := range items {
    if validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
        // Only matching URLs are processed
        websites = append(websites, ...)
    }
    // Non-matching URLs are silently filtered out
}
```

This ensures:
- Website URL must match target domain/IP/CIDR
- Subdomain must be under target domain
- Invalid items are silently skipped (not rejected with error)

### Upsert Update Strategy

| Field Type | Strategy | Example |
|------------|----------|---------|
| String fields | `COALESCE(NULLIF(new, ''), old)` | Only update if new value is non-empty |
| Nullable fields | `COALESCE(new, old)` | Only update if new value is not null |
| Array fields | Merge + deduplicate + sort | `tech` array merges and removes duplicates |
| Primary key | Don't update | `url`, `target_id` |
| Timestamp | Don't update | `created_at` |

### Implementation Example

```go
// Repository layer - use GORM OnConflict with custom assignments
result := r.db.Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "url"}, {Name: "target_id"}},
    DoUpdates: clause.Assignments(map[string]interface{}{
        "title":       gorm.Expr("COALESCE(NULLIF(EXCLUDED.title, ''), website.title)"),
        "status_code": gorm.Expr("COALESCE(EXCLUDED.status_code, website.status_code)"),
        "tech": gorm.Expr(`(
            SELECT ARRAY(SELECT DISTINCT unnest FROM unnest(
                COALESCE(website.tech, ARRAY[]::varchar(100)[]) ||
                COALESCE(EXCLUDED.tech, ARRAY[]::varchar(100)[])
            ) ORDER BY unnest)
        )`),
    }),
}).Create(&websites)
```


## GORM Bulk Operations (IMPORTANT)

Always use GORM's `clause.OnConflict` for bulk insert/upsert operations. Never use raw SQL with loops.

### ✅ Correct: GORM Batch Insert

```go
// BulkCreate - ON CONFLICT DO NOTHING
result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&batch)

// BulkUpsert - ON CONFLICT DO UPDATE
result := r.db.Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "scan_id"}, {Name: "url"}},
    DoUpdates: clause.Assignments(map[string]interface{}{
        "status_code": gorm.Expr("COALESCE(EXCLUDED.status_code, table.status_code)"),
        "image":       gorm.Expr("COALESCE(EXCLUDED.image, table.image)"),
    }),
}).Create(&batch)
```

### ❌ Wrong: Raw SQL Loop (N times slower)

```go
// DON'T DO THIS - executes N SQL statements instead of 1
sql := `INSERT INTO ... ON CONFLICT DO UPDATE SET ...`
for _, item := range items {
    r.db.Exec(sql, item.Field1, item.Field2, ...)  // ❌ One SQL per item
}
```

### Performance Comparison

| Method | 1000 Records | SQL Statements |
|--------|--------------|----------------|
| GORM `Create(&batch)` | ~50ms | 1 |
| Raw SQL loop | ~5000ms | 1000 |

### When Raw SQL is Acceptable

- Simple DELETE/UPDATE with WHERE clause (not bulk insert)
- Complex queries that GORM can't express
- One-time operations (not in hot paths)


## PostgreSQL Parameter Limits

PostgreSQL has a hard limit of **65535 parameters** per SQL statement. Bulk operations MUST use batching.

### Calculation

```
Parameters = Records × Fields
Example: 5000 websites × 14 fields = 70000 > 65535 ❌
```

### Required Pattern

```go
func (r *Repository) BulkUpsert(items []model.Item) (int64, error) {
    var totalAffected int64
    
    // MUST batch to avoid "too many parameters" error
    batchSize := 100  // Safe batch size
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }
        batch := items[i:end]
        
        affected, err := r.upsertBatch(batch)
        if err != nil {
            return totalAffected, err
        }
        totalAffected += affected
    }
    
    return totalAffected, nil
}
```

### Batch Size Guidelines

| Fields per Record | Max Batch Size | Recommended |
|-------------------|----------------|-------------|
| 5-10 | ~6000 | 500 |
| 10-15 | ~4000 | 100 |
| 15-20 | ~3000 | 100 |
| 20+ | ~2000 | 50 |
