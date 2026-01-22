# 需求文档

## 简介

创建 Go 后端项目的基础框架，包括项目结构、配置管理、数据库连接、基础数据模型。这是 Go 重构的第一阶段，为后续模块（认证、API、Worker）奠定基础。

## 术语表

- **Server**: Go 后端 API 服务
- **GORM**: Go 的 ORM 库，用于数据库操作
- **Gin**: Go 的 Web 框架
- **Viper**: Go 的配置管理库

## 需求

### 需求 1: Go 项目结构

**用户故事:** 作为开发者，我希望有清晰的项目结构，以便快速理解和扩展代码。

#### 验收标准

1. THE Go_Project SHALL 使用标准 Go 项目布局（cmd/, internal/, pkg/）
2. THE Go_Project SHALL 包含 go.mod 和 go.sum 依赖管理文件
3. THE Go_Project SHALL 包含 Makefile 用于常用构建命令
4. THE Go_Project SHALL 放置在 `go-backend/` 目录下，与现有 `backend/` 并存

### 需求 2: 配置管理

**用户故事:** 作为 DevOps 工程师，我希望使用环境变量配置服务，以便在不同环境部署。

#### 验收标准

1. THE Go_Server SHALL 从环境变量读取数据库连接信息
2. THE Go_Server SHALL 从环境变量读取 Redis 连接信息
3. THE Go_Server SHALL 从环境变量读取服务端口
4. WHEN 环境变量缺失时，THE Go_Server SHALL 使用合理的默认值
5. THE Go_Server SHALL 支持与现有 Django 相同的环境变量名称

### 需求 3: 数据库连接

**用户故事:** 作为开发者，我希望 Go 服务能连接现有 PostgreSQL 数据库，以便复用现有数据。

#### 验收标准

1. THE Go_Server SHALL 使用 GORM 连接 PostgreSQL 数据库
2. THE Go_Server SHALL 复用现有数据库表，不创建新表
3. THE Go_Server SHALL 支持数据库连接池配置
4. WHEN 数据库连接失败时，THE Go_Server SHALL 记录错误并退出
5. THE Go_Server SHALL 在启动时验证数据库连接

### 需求 4: 基础数据模型

**用户故事:** 作为开发者，我希望 Go 模型与 Django 模型兼容，以便读写相同的数据。

#### 验收标准

1. THE Go_Model SHALL 映射到现有数据库表（使用相同表名）
2. THE Go_Model SHALL 使用相同的字段名（snake_case）
3. THE Go_Model SHALL 支持软删除（deleted_at 字段）
4. THE Go_Model SHALL 正确处理 PostgreSQL 数组类型（如 engine_ids）
5. THE Go_Model SHALL 正确处理 JSONB 类型（如 stage_progress）
6. WHEN 序列化为 JSON 时，THE Go_Model SHALL 输出 camelCase 字段名

### 需求 5: 日志系统

**用户故事:** 作为运维人员，我希望有结构化日志，以便排查问题。

#### 验收标准

1. THE Go_Server SHALL 使用结构化 JSON 日志格式
2. THE Go_Server SHALL 支持可配置的日志级别（DEBUG, INFO, WARN, ERROR）
3. THE Go_Server SHALL 在日志中包含请求 ID 用于追踪
4. WHEN 发生错误时，THE Go_Server SHALL 记录完整的错误堆栈

### 需求 6: 健康检查

**用户故事:** 作为运维人员，我希望有健康检查端点，以便监控服务状态。

#### 验收标准

1. THE Go_Server SHALL 提供 `/health` 端点
2. WHEN 数据库连接正常时，THE Go_Server SHALL 返回 200 OK
3. WHEN 数据库连接异常时，THE Go_Server SHALL 返回 503 Service Unavailable
4. THE 健康检查响应 SHALL 包含数据库和 Redis 的连接状态
