# 实现计划: Go JWT 认证

## 概述

实现 JWT 认证系统，包括登录、Token 刷新、认证中间件。

## 任务

- [x] 1. 添加依赖和配置
  - [x] 1.1 添加 JWT 和 PBKDF2 依赖
    - `go get github.com/golang-jwt/jwt/v5`
    - `go get golang.org/x/crypto`
    - _需求: 5_
  - [x] 1.2 更新 config.go 添加 JWT 配置
    - 添加 JWTSecret, JWTAccessExpire, JWTRefreshExpire
    - _需求: 5_

- [x] 2. 实现核心认证逻辑
  - [x] 2.1 创建 internal/auth/jwt.go
    - 实现 GenerateAccessToken()
    - 实现 GenerateRefreshToken()
    - 实现 ValidateToken()
    - 定义 Claims 结构
    - _需求: 1, 2_
  - [x] 2.2 创建 internal/auth/password.go
    - 实现 VerifyDjangoPassword() - 兼容 Django pbkdf2_sha256
    - 实现 HashPassword() - 用于创建新用户
    - _需求: 4_

- [x] 3. 实现认证中间件
  - [x] 3.1 创建 internal/middleware/auth.go
    - 实现 AuthMiddleware()
    - 从 Authorization header 提取 Token
    - 验证 Token 并注入用户信息
    - _需求: 3_

- [x] 4. 实现认证接口
  - [x] 4.1 创建 internal/handler/auth_handler.go
    - 实现 Login() - POST /api/auth/login
    - 实现 RefreshToken() - POST /api/auth/refresh
    - 实现 GetCurrentUser() - GET /api/auth/me
    - _需求: 1, 2_
  - [x] 4.2 注册路由
    - 在 main.go 中注册认证路由
    - _需求: 1, 2_

- [x] 5. 编写测试
  - [x] 5.1 创建 internal/auth/auth_test.go
    - 测试 JWT 生成和验证
    - 测试 Django 密码验证
    - _需求: 1, 2, 4_

- [x] 6. 检查点 - 验证认证功能
  - 所有测试通过 ✅
