# 设计文档: Go JWT 认证

## 概述

实现基于 JWT 的认证系统，使用 `golang-jwt/jwt` 库。

## 架构

```
请求 → AuthMiddleware → Handler
         ↓
    验证 JWT Token
         ↓
    注入用户信息到 Context
```

## 目录结构

```
internal/
├── auth/
│   ├── jwt.go              # JWT 生成和验证
│   ├── password.go         # 密码验证（Django 兼容）
│   └── auth_test.go        # 测试
├── handler/
│   └── auth_handler.go     # 登录/刷新接口
├── middleware/
│   └── auth.go             # 认证中间件
└── config/
    └── config.go           # 添加 JWT 配置
```

## 核心组件

### 1. JWT Token 结构

```go
type Claims struct {
    UserID   int    `json:"userId"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

### 2. 登录流程

```
用户输入 → 查询数据库 → 验证密码 → 生成 Token → 返回
```

### 3. 密码验证（Django 兼容）

Django 密码格式: `pbkdf2_sha256$iterations$salt$hash`

```go
func VerifyPassword(password, encoded string) bool {
    // 1. 解析 encoded 字符串
    // 2. 使用相同参数计算 PBKDF2
    // 3. 比较哈希值
}
```

### 4. 中间件

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 提取 Authorization header
        // 2. 验证 Bearer Token
        // 3. 解析 Claims
        // 4. 注入到 Context
        c.Set("user", claims)
        c.Next()
    }
}
```

## API 设计

### POST /api/auth/login

请求:
```json
{
    "username": "admin",
    "password": "admin"
}
```

响应:
```json
{
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
}
```

### POST /api/auth/refresh

请求:
```json
{
    "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

响应:
```json
{
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "expiresIn": 900
}
```

## 依赖

- `github.com/golang-jwt/jwt/v5` - JWT 库
- `golang.org/x/crypto/pbkdf2` - PBKDF2 密码验证

## 配置

```yaml
jwt:
  secret: "your-secret-key"      # JWT 签名密钥
  accessExpire: 15m              # Access Token 过期时间
  refreshExpire: 168h            # Refresh Token 过期时间 (7天)
```

## 安全考虑

1. JWT Secret 必须足够长（至少 32 字符）
2. 生产环境必须通过环境变量配置 Secret
3. 密码验证失败不透露具体原因（统一返回"用户名或密码错误"）
4. Token 不记录到日志
