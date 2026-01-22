# 需求文档: Go JWT 认证

## 简介

为 Go 后端实现 JWT 认证系统，包括登录、Token 刷新、认证中间件。

## 术语表

- **JWT**: JSON Web Token，无状态认证令牌
- **Access Token**: 短期令牌，用于 API 认证（15分钟）
- **Refresh Token**: 长期令牌，用于刷新 Access Token（7天）

## 需求

### 需求 1: 登录接口

**用户故事:** 作为用户，我希望通过用户名密码登录，获取 JWT Token。

#### 验收标准

1. THE API SHALL 提供 `POST /api/auth/login` 接口
2. THE API SHALL 验证用户名密码（兼容 Django pbkdf2_sha256 密码格式）
3. THE API SHALL 返回 access_token 和 refresh_token
4. THE API SHALL 在登录失败时返回 401 错误

### 需求 2: Token 刷新接口

**用户故事:** 作为用户，我希望在 Access Token 过期前刷新它。

#### 验收标准

1. THE API SHALL 提供 `POST /api/auth/refresh` 接口
2. THE API SHALL 验证 Refresh Token 有效性
3. THE API SHALL 返回新的 access_token
4. THE API SHALL 在 Refresh Token 无效时返回 401 错误

### 需求 3: 认证中间件

**用户故事:** 作为开发者，我希望有一个中间件自动验证 JWT Token。

#### 验收标准

1. THE Middleware SHALL 从 Authorization header 提取 Bearer Token
2. THE Middleware SHALL 验证 Token 签名和过期时间
3. THE Middleware SHALL 将用户信息注入到 Gin Context
4. THE Middleware SHALL 在 Token 无效时返回 401 错误

### 需求 4: 密码验证

**用户故事:** 作为系统，我需要验证用户密码与 Django 存储的密码哈希匹配。

#### 验收标准

1. THE System SHALL 支持 Django pbkdf2_sha256 密码格式
2. THE System SHALL 正确解析 `pbkdf2_sha256$iterations$salt$hash` 格式
3. THE System SHALL 使用相同算法验证密码

### 需求 5: 配置管理

**用户故事:** 作为运维，我希望通过环境变量配置 JWT 参数。

#### 验收标准

1. THE Config SHALL 支持 `JWT_SECRET` 环境变量
2. THE Config SHALL 支持 `JWT_ACCESS_EXPIRE` 环境变量（默认 15 分钟）
3. THE Config SHALL 支持 `JWT_REFRESH_EXPIRE` 环境变量（默认 7 天）

## 非功能需求

- Token 签名算法使用 HS256
- 密码验证使用 PBKDF2-SHA256（兼容 Django）
- 所有敏感信息不记录到日志
