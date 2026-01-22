# 实现计划: Go 后端基础框架

## 概述

创建 Go 后端项目的基础框架，包括项目结构、配置管理、数据库连接、基础数据模型、日志系统和健康检查端点。

## 任务

- [x] 1. 初始化 Go 项目结构
  - [x] 1.1 创建 `go-backend/` 目录和标准布局
    - 创建 `cmd/server/`, `internal/config/`, `internal/database/`, `internal/model/`, `internal/handler/`, `internal/middleware/`, `internal/pkg/` 目录
    - _需求: 1.1, 1.4_
  - [x] 1.2 初始化 Go 模块和依赖
    - 创建 `go.mod` 文件，添加 Gin, GORM, Viper, Zap 依赖
    - _需求: 1.2_
  - [x] 1.3 创建 Makefile
    - 包含 `build`, `run`, `test`, `lint` 命令
    - _需求: 1.3_
  - [x] 1.4 创建 `.env.example` 配置示例文件
    - _需求: 2.1, 2.2, 2.3_

- [x] 2. 实现配置管理
  - [x] 2.1 实现配置结构体和加载逻辑
    - 创建 `internal/config/config.go`
    - 定义 ServerConfig, DatabaseConfig, RedisConfig, LogConfig 结构体
    - 使用 Viper 从环境变量加载配置
    - _需求: 2.1, 2.2, 2.3, 2.4, 2.5_
  - [x] 2.2 编写配置默认值属性测试
    - **Property 4: 配置默认值正确性**
    - **验证: 需求 2.4**

- [x] 3. 实现日志系统
  - [x] 3.1 实现 Zap 日志封装
    - 创建 `internal/pkg/logger.go`
    - 支持 JSON 格式和可配置日志级别
    - _需求: 5.1, 5.2, 5.3, 5.4_

- [x] 4. 实现数据库连接
  - [x] 4.1 实现数据库连接和连接池
    - 创建 `internal/database/database.go`
    - 使用 GORM 连接 PostgreSQL
    - 配置连接池参数
    - _需求: 3.1, 3.3, 3.4, 3.5_

- [x] 5. 实现基础数据模型
  - [x] 5.1 实现 Organization 模型
    - 创建 `internal/model/organization.go`
    - _需求: 4.1, 4.2, 4.3, 4.6_
  - [x] 5.2 实现 Target 模型
    - 创建 `internal/model/target.go`
    - _需求: 4.1, 4.2, 4.3, 4.6_
  - [x] 5.3 实现 Scan 模型
    - 创建 `internal/model/scan.go`
    - 处理 PostgreSQL 数组类型和 JSONB 类型
    - _需求: 4.1, 4.2, 4.4, 4.5, 4.6_
  - [x] 5.4 实现资产模型 (Subdomain, WebSite)
    - 创建 `internal/model/subdomain.go`, `internal/model/website.go`
    - _需求: 4.1, 4.2, 4.6_
  - [x] 5.5 实现引擎模型 (WorkerNode, ScanEngine)
    - 创建 `internal/model/worker_node.go`, `internal/model/scan_engine.go`
    - _需求: 4.1, 4.2, 4.6_
  - [x] 5.6 编写模型表名映射属性测试
    - **Property 1: 数据库表名映射正确性**
    - **验证: 需求 4.1**
  - [x] 5.7 编写 JSON 字段名属性测试
    - **Property 2: JSON 字段名转换正确性**
    - **验证: 需求 4.6**
  - [x] 5.8 编写数据库字段映射属性测试
    - **Property 3: 数据库字段映射正确性**
    - **验证: 需求 4.2**

- [x] 6. 实现中间件
  - [x] 6.1 实现日志中间件
    - 创建 `internal/middleware/logger.go`
    - 记录请求信息和请求 ID
    - _需求: 5.3_
  - [x] 6.2 实现 Recovery 中间件
    - 创建 `internal/middleware/recovery.go`
    - 捕获 panic 并记录错误堆栈
    - _需求: 5.4_

- [x] 7. 实现健康检查端点
  - [x] 7.1 实现健康检查 Handler
    - 创建 `internal/handler/health.go`
    - 检查数据库和 Redis 连接状态
    - _需求: 6.1, 6.2, 6.3, 6.4_
  - [x] 7.2 实现响应工具
    - 创建 `internal/pkg/response.go`
    - 统一 API 响应格式
    - _需求: 6.4_

- [x] 8. 实现服务入口
  - [x] 8.1 实现 main.go
    - 创建 `cmd/server/main.go`
    - 初始化配置、日志、数据库
    - 注册路由和中间件
    - 启动 HTTP 服务
    - _需求: 3.4, 3.5_

- [x] 9. 检查点 - 验证基础框架
  - 确保所有测试通过
  - 验证服务能正常启动并连接数据库
  - 验证 `/health` 端点返回正确状态
  - 如有问题请询问用户

## 备注

- 每个任务都引用了具体的需求以便追踪
- 检查点用于确保增量验证
- 属性测试验证通用正确性属性
