# 实现计划: Go 网站快照 API

## 概述

实现 Go 后端的网站快照 API，包括批量写入（同步到资产表）、列表查询、CSV 导出三个接口。

## 任务

- [x] 1. 创建 DTO 定义
  - 创建 `go-backend/internal/dto/website_snapshot.go`
  - 定义请求和响应结构体
  - _Requirements: 1.6, 2.7, 5.3, 6.1, 6.2_

- [x] 2. 创建 Repository 层
  - [x] 2.1 创建 WebsiteSnapshotRepository
    - 创建 `go-backend/internal/repository/website_snapshot.go`
    - 实现 BulkCreate（ON CONFLICT DO NOTHING）
    - 实现 FindByScanID（分页、过滤、排序）
    - 实现 StreamByScanID 和 CountByScanID
    - _Requirements: 1.4, 2.1, 3.1, 4.1_

  - [x] 2.2 编写 Repository 单元测试
    - 测试 BulkCreate 去重逻辑
    - 测试分页、过滤、排序
    - _Requirements: 1.4, 2.3, 3.1, 4.1_

- [x] 3. 创建 Service 层
  - [x] 3.1 创建 WebsiteSnapshotService
    - 创建 `go-backend/internal/service/website_snapshot.go`
    - 实现 SaveAndSync（写快照 + 调用 WebsiteService.BulkUpsert）
    - 实现 ListByScan、StreamByScan、CountByScan
    - 添加 Scan 存在性验证
    - _Requirements: 1.1, 1.2, 1.3, 7.1, 7.3_

  - [x] 3.2 编写 Property 测试：快照和资产同步写入
    - **Property 1: 快照和资产同步写入**
    - **Validates: Requirements 1.1, 1.2**

  - [x] 3.3 编写 Property 测试：Scan 存在性验证
    - **Property 8: Scan 存在性验证**
    - **Validates: Requirements 7.1, 7.2, 7.3, 7.4**

- [x] 4. 创建 Handler 层
  - [x] 4.1 创建 WebsiteSnapshotHandler
    - 创建 `go-backend/internal/handler/website_snapshot.go`
    - 实现 BulkUpsert 接口
    - 实现 List 接口
    - 实现 Export 接口
    - _Requirements: 1.1, 2.1, 5.1_

  - [x] 4.2 编写 Property 测试：分页正确性 (合并到集成测试)
    - **Property 4: 分页正确性**
    - **Validates: Requirements 2.3, 2.6**

  - [x] 4.3 编写 Property 测试：过滤正确性 (合并到集成测试)
    - **Property 5: 过滤正确性**
    - **Validates: Requirements 3.1, 3.2, 3.4**

- [x] 5. 注册路由和依赖注入
  - [x] 5.1 更新路由配置
    - 在 `go-backend/cmd/server/main.go` 或路由文件中注册新路由
    - 配置依赖注入
    - _Requirements: 1.1, 2.1, 5.1_

- [x] 6. Checkpoint - 确保所有测试通过
  - 运行所有测试，确保通过
  - 如有问题，询问用户

- [x] 7. 编写集成测试
  - [x] 7.1 编写 API 集成测试
    - 测试完整的请求-响应流程
    - 测试错误场景
    - _Requirements: 1.3, 2.2, 5.2, 6.3, 6.4_
    - 注: Handler 测试已覆盖主要场景，包括分页、过滤、错误处理

- [x] 8. 最终 Checkpoint
  - 确保所有测试通过
  - 如有问题，询问用户

## 备注

- 每个任务引用具体的需求以便追溯
- Property 测试验证核心正确性属性
- 单元测试验证具体示例和边界情况
