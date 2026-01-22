# 设计文档

## 概述

创建 Go 后端项目的基础框架，为后续模块开发奠定基础。本阶段聚焦于项目结构、配置管理、数据库连接和基础模型定义。

## 架构

### 项目结构

```
go-backend/
├── cmd/
│   └── server/
│       └── main.go              # Server 入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── database/
│   │   └── database.go          # 数据库连接
│   ├── model/                   # 数据模型（扁平结构，Go 风格）
│   │   ├── organization.go
│   │   ├── target.go
│   │   ├── scan.go
│   │   ├── subdomain.go
│   │   ├── website.go
│   │   ├── endpoint.go
│   │   ├── directory.go
│   │   ├── vulnerability.go
│   │   ├── host_port_mapping.go
│   │   ├── worker_node.go
│   │   ├── scan_engine.go
│   │   └── user.go
│   ├── handler/                 # HTTP 处理器
│   │   └── health.go
│   ├── middleware/
│   │   ├── logger.go
│   │   └── recovery.go
│   └── pkg/                     # 内部工具包
│       ├── logger.go            # 日志工具
│       └── response.go          # 响应工具
├── go.mod
├── go.sum
├── Makefile
└── .env.example
```

### Go 风格规范

| 规范 | 说明 |
|------|------|
| 包名 | 小写单词，不用下划线（`model` 不是 `models`） |
| 文件名 | 小写 + 下划线（`worker_node.go`） |
| 导出标识符 | PascalCase（`Target`, `GetByID`） |
| 私有标识符 | camelCase（`targetID`, `getByID`） |
| 接口命名 | 动词 + er（`Reader`, `Scanner`） |
| 错误变量 | `Err` 前缀（`ErrNotFound`） |
| 常量 | PascalCase 或全大写（`MaxRetries`, `DEFAULT_PORT`） |

### 技术选型

| 组件 | 选择 | 版本 |
|------|------|------|
| Web 框架 | Gin | v1.9+ |
| ORM | GORM | v1.25+ |
| 配置管理 | Viper | v1.18+ |
| 日志 | Zap | v1.26+ |
| PostgreSQL 驱动 | pgx | v5+ |

## 组件和接口

### 配置结构

```go
// internal/config/config.go

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Log      LogConfig
}

type ServerConfig struct {
    Port int    `mapstructure:"SERVER_PORT" default:"8888"`
    Mode string `mapstructure:"GIN_MODE" default:"release"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"DB_HOST" default:"localhost"`
    Port     int    `mapstructure:"DB_PORT" default:"5432"`
    User     string `mapstructure:"DB_USER" default:"postgres"`
    Password string `mapstructure:"DB_PASSWORD"`
    Name     string `mapstructure:"DB_NAME" default:"xingrin"`
    SSLMode  string `mapstructure:"DB_SSLMODE" default:"disable"`
    
    MaxOpenConns    int `mapstructure:"DB_MAX_OPEN_CONNS" default:"25"`
    MaxIdleConns    int `mapstructure:"DB_MAX_IDLE_CONNS" default:"5"`
    ConnMaxLifetime int `mapstructure:"DB_CONN_MAX_LIFETIME" default:"300"`
}

type RedisConfig struct {
    Host     string `mapstructure:"REDIS_HOST" default:"localhost"`
    Port     int    `mapstructure:"REDIS_PORT" default:"6379"`
    Password string `mapstructure:"REDIS_PASSWORD"`
    DB       int    `mapstructure:"REDIS_DB" default:"0"`
}

type LogConfig struct {
    Level  string `mapstructure:"LOG_LEVEL" default:"info"`
    Format string `mapstructure:"LOG_FORMAT" default:"json"`
}
```

### 数据库连接

```go
// internal/database/database.go

func NewDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,  // 使用单数表名
        },
    })
    if err != nil {
        return nil, err
    }
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
    
    return db, nil
}
```

## 数据模型

### 基础模型定义

```go
// internal/model/target.go

type Target struct {
    ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
    Name          string     `gorm:"column:name;size:300" json:"name"`
    Type          string     `gorm:"column:type;size:20;default:'domain'" json:"type"`
    CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    LastScannedAt *time.Time `gorm:"column:last_scanned_at" json:"lastScannedAt"`
    DeletedAt     *time.Time `gorm:"column:deleted_at;index" json:"-"`
}

func (Target) TableName() string {
    return "target"
}

// internal/model/organization.go

type Organization struct {
    ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string     `gorm:"column:name;size:300" json:"name"`
    Description string     `gorm:"column:description;size:1000" json:"description"`
    CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"-"`
}

func (Organization) TableName() string {
    return "organization"
}

// internal/model/scan.go

type Scan struct {
    ID                int             `gorm:"primaryKey;autoIncrement" json:"id"`
    TargetID          int             `gorm:"column:target_id;not null" json:"targetId"`
    EngineIDs         pq.Int64Array   `gorm:"column:engine_ids;type:integer[]" json:"engineIds"`
    EngineNames       datatypes.JSON  `gorm:"column:engine_names;type:jsonb" json:"engineNames"`
    YamlConfiguration string          `gorm:"column:yaml_configuration;type:text" json:"yamlConfiguration"`
    ScanMode          string          `gorm:"column:scan_mode;size:10;default:'full'" json:"scanMode"`
    Status            string          `gorm:"column:status;size:20;default:'initiated'" json:"status"`
    ResultsDir        string          `gorm:"column:results_dir;size:100" json:"resultsDir"`
    ContainerIDs      pq.StringArray  `gorm:"column:container_ids;type:varchar(100)[]" json:"containerIds"`
    WorkerID          *int            `gorm:"column:worker_id" json:"workerId"`
    ErrorMessage      string          `gorm:"column:error_message;size:2000" json:"errorMessage"`
    Progress          int             `gorm:"column:progress;default:0" json:"progress"`
    CurrentStage      string          `gorm:"column:current_stage;size:50" json:"currentStage"`
    StageProgress     datatypes.JSON  `gorm:"column:stage_progress;type:jsonb" json:"stageProgress"`
    CreatedAt         time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    StoppedAt         *time.Time      `gorm:"column:stopped_at" json:"stoppedAt"`
    DeletedAt         *time.Time      `gorm:"column:deleted_at;index" json:"-"`
    
    // 缓存统计
    CachedSubdomainsCount int `gorm:"column:cached_subdomains_count" json:"cachedSubdomainsCount"`
    CachedWebsitesCount   int `gorm:"column:cached_websites_count" json:"cachedWebsitesCount"`
    CachedEndpointsCount  int `gorm:"column:cached_endpoints_count" json:"cachedEndpointsCount"`
    CachedIPsCount        int `gorm:"column:cached_ips_count" json:"cachedIpsCount"`
    CachedVulnsTotal      int `gorm:"column:cached_vulns_total" json:"cachedVulnsTotal"`
}

func (Scan) TableName() string {
    return "scan"
}

// internal/model/asset.go

type Subdomain struct {
    ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
    TargetID  int       `gorm:"column:target_id;not null" json:"targetId"`
    Name      string    `gorm:"column:name;size:1000" json:"name"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Subdomain) TableName() string {
    return "subdomain"
}

type WebSite struct {
    ID              int            `gorm:"primaryKey;autoIncrement" json:"id"`
    TargetID        int            `gorm:"column:target_id;not null" json:"targetId"`
    URL             string         `gorm:"column:url;type:text" json:"url"`
    Host            string         `gorm:"column:host;size:253" json:"host"`
    Title           string         `gorm:"column:title;type:text" json:"title"`
    StatusCode      *int           `gorm:"column:status_code" json:"statusCode"`
    ContentLength   *int           `gorm:"column:content_length" json:"contentLength"`
    Tech            pq.StringArray `gorm:"column:tech;type:varchar(100)[]" json:"tech"`
    Webserver       string         `gorm:"column:webserver;type:text" json:"webserver"`
    ResponseHeaders string         `gorm:"column:response_headers;type:text" json:"responseHeaders"`
    CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (WebSite) TableName() string {
    return "website"
}

// internal/model/engine.go

type WorkerNode struct {
    ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"column:name;size:100;uniqueIndex" json:"name"`
    IPAddress string    `gorm:"column:ip_address;type:inet" json:"ipAddress"`
    SSHPort   int       `gorm:"column:ssh_port;default:22" json:"sshPort"`
    Username  string    `gorm:"column:username;size:50;default:'root'" json:"username"`
    Password  string    `gorm:"column:password;size:200" json:"-"`
    IsLocal   bool      `gorm:"column:is_local;default:false" json:"isLocal"`
    Status    string    `gorm:"column:status;size:20;default:'pending'" json:"status"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (WorkerNode) TableName() string {
    return "worker_node"
}

type ScanEngine struct {
    ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
    Name          string    `gorm:"column:name;size:200;uniqueIndex" json:"name"`
    Configuration string    `gorm:"column:configuration;size:10000" json:"configuration"`
    CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
    UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (ScanEngine) TableName() string {
    return "scan_engine"
}
```

### 健康检查端点

```go
// internal/handler/health.go

type HealthHandler struct {
    db    *gorm.DB
    redis *redis.Client
}

type HealthResponse struct {
    Status   string            `json:"status"`
    Database string            `json:"database"`
    Redis    string            `json:"redis"`
    Details  map[string]string `json:"details,omitempty"`
}

func (h *HealthHandler) Check(c *gin.Context) {
    resp := HealthResponse{
        Status:   "healthy",
        Database: "connected",
        Redis:    "connected",
    }
    
    // 检查数据库
    sqlDB, _ := h.db.DB()
    if err := sqlDB.Ping(); err != nil {
        resp.Status = "unhealthy"
        resp.Database = "disconnected"
    }
    
    // 检查 Redis
    if err := h.redis.Ping(c).Err(); err != nil {
        resp.Status = "unhealthy"
        resp.Redis = "disconnected"
    }
    
    if resp.Status == "healthy" {
        c.JSON(200, resp)
    } else {
        c.JSON(503, resp)
    }
}
```

## 正确性属性

*正确性属性是系统在所有有效执行中都应保持的特性。*

### Property 1: 数据库表名映射正确性
*对于任意* Go 模型，其 TableName() 方法返回的表名应与 Django 模型的 db_table 一致。
**验证: 需求 4.1**

### Property 2: JSON 字段名转换正确性
*对于任意* Go 模型序列化为 JSON，所有字段名应为 camelCase 格式。
**验证: 需求 4.6**

### Property 3: 数据库字段映射正确性
*对于任意* Go 模型字段，其 gorm column tag 应与数据库实际列名（snake_case）一致。
**验证: 需求 4.2**

### Property 4: 配置默认值正确性
*对于任意* 缺失的环境变量，配置系统应返回预定义的默认值。
**验证: 需求 2.4**

## 错误处理

### 启动错误

```go
func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("配置加载失败", zap.Error(err))
    }
    
    // 连接数据库
    db, err := database.NewDatabase(&cfg.Database)
    if err != nil {
        log.Fatal("数据库连接失败", zap.Error(err))
    }
    
    // 验证数据库连接
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("数据库 Ping 失败", zap.Error(err))
    }
    
    log.Info("服务启动成功", zap.Int("port", cfg.Server.Port))
}
```

## 测试策略

### 单元测试

1. **配置测试**: 验证环境变量读取和默认值
2. **模型测试**: 验证表名和字段映射
3. **数据库测试**: 验证连接和基本查询

### 属性测试

```go
// 测试 JSON 字段名为 camelCase
func TestJSONFieldNames(t *testing.T) {
    target := model.Target{ID: 1, Name: "test.com"}
    jsonBytes, _ := json.Marshal(target)
    jsonStr := string(jsonBytes)
    
    // 不应包含 snake_case
    assert.NotContains(t, jsonStr, "created_at")
    assert.NotContains(t, jsonStr, "last_scanned_at")
    
    // 应包含 camelCase
    assert.Contains(t, jsonStr, "createdAt")
    assert.Contains(t, jsonStr, "lastScannedAt")
}

// 测试表名映射
func TestTableNames(t *testing.T) {
    tests := []struct {
        model     interface{ TableName() string }
        expected  string
    }{
        {model.Target{}, "target"},
        {model.Organization{}, "organization"},
        {model.Scan{}, "scan"},
        {model.Subdomain{}, "subdomain"},
        {model.WebSite{}, "website"},
        {model.WorkerNode{}, "worker_node"},
        {model.ScanEngine{}, "scan_engine"},
    }
    
    for _, tt := range tests {
        assert.Equal(t, tt.expected, tt.model.TableName())
    }
}
```
