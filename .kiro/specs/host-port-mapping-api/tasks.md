# Implementation Plan: Host Port Mapping API

## Overview

为 Go 后端实现 host_port_mapping 资产的 CRUD API，包括列表查询（按 IP 聚合）、CSV 导出、批量 upsert 和批量删除功能。

## Tasks

- [x] 1. 创建 DTO 定义
  - [x] 1.1 创建 `go-backend/internal/dto/host_port_mapping.go`
    - 定义 `HostPortMappingListQuery`（分页 + 过滤）
    - 定义 `HostPortMappingResponse`（聚合格式：ip, hosts[], ports[], createdAt）
    - 定义 `HostPortMappingListResponse`（分页响应）
    - 定义 `HostPortMappingItem`（单条映射）
    - 定义 `BulkUpsertHostPortMappingsRequest/Response`
    - 定义 `BulkDeleteHostPortMappingsRequest/Response`（按 IP 列表删除）
    - _Requirements: 1.1, 1.4, 3.1, 4.1_

- [x] 2. 创建 Repository 层
  - [x] 2.1 创建 `go-backend/internal/repository/host_port_mapping.go`
    - 实现 `GetIPAggregation()` - 按 IP 分组查询
    - 实现 `GetHostsAndPortsByIP()` - 获取指定 IP 的 hosts 和 ports
    - 实现 `StreamByTargetID()` - 流式导出原始数据
    - 实现 `BulkUpsert()` - 批量插入（ON CONFLICT DO NOTHING）
    - 实现 `DeleteByIPs()` - 按 IP 列表删除
    - 实现 `ScanRow()` - 扫描单行数据
    - _Requirements: 1.1, 1.2, 2.1, 3.2, 4.1_

- [x] 3. 创建 Service 层
  - [x] 3.1 创建 `go-backend/internal/service/host_port_mapping.go`
    - 实现 `ListByTarget()` - 返回按 IP 聚合的分页数据
    - 实现 `StreamByTarget()` - 流式导出
    - 实现 `BulkUpsert()` - 批量 upsert（分批处理）
    - 实现 `BulkDeleteByIPs()` - 按 IP 删除
    - _Requirements: 1.1, 1.2, 2.2, 3.1, 4.1_

- [x] 4. 创建 Handler 层
  - [x] 4.1 创建 `go-backend/internal/handler/host_port_mapping.go`
    - 实现 `List()` - GET /targets/:id/host-port-mappings
    - 实现 `Export()` - GET /targets/:id/host-port-mappings/export
    - 实现 `BulkUpsert()` - POST /targets/:id/host-port-mappings/bulk-upsert
    - 实现 `BulkDelete()` - POST /host-port-mappings/bulk-delete
    - _Requirements: 1.1, 1.3, 2.1, 2.3, 3.3, 3.4, 4.2_

- [x] 5. 注册路由
  - [x] 5.1 更新 `go-backend/cmd/server/main.go`
    - 注册 `/targets/:id/host-port-mappings` 路由组
    - 注册 `/host-port-mappings/bulk-delete` 独立路由
    - _Requirements: 1.1, 2.1, 3.1, 4.1_

- [x] 6. Checkpoint - 编译测试
  - 确保代码编译通过
  - 手动测试 API 端点
  - 如有问题请询问用户

## Notes

- 遵循项目现有的资产 API 设计模式（参考 website、subdomain）
- List 接口返回按 IP 聚合的数据，与前端 `IPAddress` 类型匹配
- Export 接口返回原始格式 CSV（每行一个 host+ip+port 组合）
- Bulk Delete 接收 IP 字符串列表（不是 ID），删除这些 IP 的所有映射记录
- 批量操作使用 100 条/批次，避免 PostgreSQL 参数限制
