# Implementation Plan: Go Asset APIs

## Overview

实现 Go 后端的 Subdomain、Endpoint、Directory 三类资产 API，包括列表查询、批量创建、批量删除、导出功能，以及 Asset-Target 匹配验证。

## Tasks

- [x] 1. 创建 Target 匹配验证函数
  - [x] 1.1 创建 `go-backend/internal/pkg/validator/target.go`
    - 实现 `IsURLMatchTarget(urlStr, targetName, targetType string) bool`
    - 实现 `IsSubdomainMatchTarget(subdomain, targetDomain string) bool`
    - 实现 `DetectTargetType(name string) string`
    - 使用 Go 标准库 `net/url`、`net`、`strings`
    - _Requirements: 2.6, 2.7, 7.6, 12.6_
  - [x] 1.2 编写验证函数单元测试
    - 测试 domain 类型匹配（精确匹配、后缀匹配）
    - 测试 IP 类型匹配
    - 测试 CIDR 类型匹配
    - 测试边界情况（空字符串、无效 URL）
    - _Requirements: 2.6, 2.7, 7.6, 12.6_

- [x] 2. 实现 Subdomain API
  - [x] 2.1 创建 `go-backend/internal/dto/subdomain.go`
  - [x] 2.2 创建 `go-backend/internal/repository/subdomain.go`
  - [x] 2.3 创建 `go-backend/internal/service/subdomain.go`
  - [x] 2.4 创建 `go-backend/internal/handler/subdomain.go`
  - [x] 2.5 注册 Subdomain 路由到 main.go

- [x] 3. Checkpoint - 验证 Subdomain API ✓

- [x] 4. 实现 Endpoint API
  - [x] 4.1 创建 `go-backend/internal/dto/endpoint.go`
  - [x] 4.2 创建 `go-backend/internal/repository/endpoint.go`
  - [x] 4.3 创建 `go-backend/internal/service/endpoint.go`
  - [x] 4.4 创建 `go-backend/internal/handler/endpoint.go`
  - [x] 4.5 注册 Endpoint 路由到 main.go

- [x] 5. Checkpoint - 验证 Endpoint API ✓

- [x] 6. 实现 Directory API
  - [x] 6.1 创建 `go-backend/internal/dto/directory.go`
  - [x] 6.2 创建 `go-backend/internal/repository/directory.go`
  - [x] 6.3 创建 `go-backend/internal/service/directory.go`
  - [x] 6.4 创建 `go-backend/internal/handler/directory.go`
  - [x] 6.5 注册 Directory 路由到 main.go

- [x] 7. Checkpoint - 验证 Directory API ✓

- [x] 8. 更新 Website Service 添加验证
  - [x] 8.1 修改 `go-backend/internal/service/website.go`
    - 在 BulkCreate 中添加 URL 匹配验证
    - 使用 validator.IsURLMatchTarget 函数

- [x] 9. 更新种子数据生成
  - [x] 9.1 修改 `go-backend/cmd/seed/main.go`
    - 添加 createSubdomains 函数（仅为 domain 类型 target 生成）
    - 添加 createEndpoints 函数
    - 添加 createDirectories 函数
    - 每个资产类型生成 20 条数据
    - 更新 clearData 函数添加新表

- [x] 10. Final Checkpoint ✓
  - 所有代码编译通过
  - 所有 API 路由已注册

## Notes

- 验证函数使用 Go 标准库，不需要第三方库
- CSV 导出使用流式处理，避免内存问题
- 所有 filter 字段都已有索引覆盖
- 数组字段（tech）使用 GIN 索引，scope 包已支持

## Created Files

- `go-backend/internal/dto/subdomain.go`
- `go-backend/internal/dto/endpoint.go`
- `go-backend/internal/dto/directory.go`
- `go-backend/internal/repository/subdomain.go`
- `go-backend/internal/repository/endpoint.go`
- `go-backend/internal/repository/directory.go`
- `go-backend/internal/service/subdomain.go`
- `go-backend/internal/service/endpoint.go`
- `go-backend/internal/service/directory.go`
- `go-backend/internal/handler/subdomain.go`
- `go-backend/internal/handler/endpoint.go`
- `go-backend/internal/handler/directory.go`

## Modified Files

- `go-backend/internal/pkg/validator/target.go` - 添加 DetectTargetType 函数
- `go-backend/internal/service/website.go` - 添加 URL 匹配验证
- `go-backend/cmd/server/main.go` - 注册新路由
- `go-backend/cmd/seed/main.go` - 添加新资产类型种子数据生成
