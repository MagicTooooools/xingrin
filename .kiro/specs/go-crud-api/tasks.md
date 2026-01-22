# 实现计划: Go CRUD API

## 概述

实现基础 CRUD API，采用三层架构。

## 任务

- [x] 1. 创建基础设施
  - [x] 1.1 创建 dto 包（请求/响应/分页）
  - [x] 1.2 创建 repository 基础接口

- [x] 2. 实现用户管理
  - [x] 2.1 创建 user_repository.go
  - [x] 2.2 创建 user_service.go
  - [x] 2.3 创建 user_handler.go
  - [x] 2.4 注册路由
  - _需求: 1_

- [x] 3. 实现组织管理
  - [x] 3.1 创建 organization_repository.go
  - [x] 3.2 创建 organization_service.go
  - [x] 3.3 创建 organization_handler.go
  - [x] 3.4 注册路由
  - _需求: 2_

- [x] 4. 实现目标管理
  - [x] 4.1 创建 target_repository.go
  - [x] 4.2 创建 target_service.go（含类型检测）
  - [x] 4.3 创建 target_handler.go
  - [x] 4.4 注册路由
  - _需求: 3_

- [x] 5. 实现引擎管理
  - [x] 5.1 创建 engine_repository.go
  - [x] 5.2 创建 engine_service.go
  - [x] 5.3 创建 engine_handler.go
  - [x] 5.4 注册路由
  - _需求: 4_

- [x] 6. 检查点
  - 所有测试通过
  - API 可正常调用
