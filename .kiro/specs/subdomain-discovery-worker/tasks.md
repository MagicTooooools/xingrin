# 任务列表

## 任务 1: 初始化 Worker 项目结构

**需求:** 需求 1

**依赖:** 无

### 子任务

- [x] 1.1 创建 `worker/` 目录和 `go.mod`
- [x] 1.2 创建项目目录结构（cmd/、internal/）
- [x] 1.3 创建 `cmd/worker/main.go` 入口文件
- [x] 1.4 创建 `Makefile` 和 `.env.example`
- [x] 1.5 添加基础依赖（gin、zap、viper）

---

## 任务 2: 实现配置管理

**需求:** 需求 2

**依赖:** 任务 1

### 子任务

- [x] 2.1 创建 `internal/config/config.go`
- [x] 2.2 实现从环境变量和 .env 文件加载配置
- [x] 2.3 实现必需配置项验证
- [ ] 2.4 编写配置加载单元测试

---

## 任务 3: 实现日志工具

**需求:** 需求 10

**依赖:** 任务 2

### 子任务

- [x] 3.1 创建 `internal/pkg/logger.go`
- [x] 3.2 实现 JSON 格式结构化日志
- [x] 3.3 实现开发环境人类可读日志
- [x] 3.4 支持日志级别配置

---

## 任务 4: 实现 Server API 客户端

**需求:** 需求 3

**依赖:** 任务 2, 任务 3

### 子任务

- [x] 4.1 创建 `internal/client/server_client.go`
- [x] 4.2 实现 UpdateScanStatus 方法
- [x] 4.3 实现 SaveSubdomains 方法
- [x] 4.4 实现 WriteScanLog 方法
- [x] 4.5 实现指数退避重试逻辑
- [ ] 4.6 编写 ServerClient 单元测试

---

## 任务 5: 实现命令模板系统

**需求:** 需求 4

**依赖:** 任务 1

### 子任务

- [x] 5.1 创建 `internal/tool/templates.go`
- [x] 5.2 定义子域名发现工具的命令模板
- [x] 5.3 创建 `internal/tool/command_builder.go`
- [x] 5.4 实现占位符替换逻辑
- [x] 5.5 实现可选参数追加逻辑
- [ ] 5.6 编写命令构建单元测试

---

## 任务 6: 实现工具执行器

**需求:** 需求 5

**依赖:** 任务 3, 任务 5

### 子任务

- [x] 6.1 创建 `internal/tool/runner.go`
- [x] 6.2 实现单工具执行（带超时）
- [x] 6.3 实现 stdout/stderr 捕获
- [x] 6.4 实现日志文件写入
- [x] 6.5 实现并行执行多个工具
- [ ] 6.6 编写工具执行器单元测试

---

## 任务 7: 实现结果解析器

**需求:** 需求 7

**依赖:** 任务 1

### 子任务

- [x] 7.1 创建 `internal/parser/subdomain.go`
- [x] 7.2 实现从文件解析子域名（每行一个）
- [x] 7.3 实现子域名去重逻辑
- [ ] 7.4 编写解析器单元测试

---

## 任务 8: 实现 Flow 框架

**需求:** 需求 6, 需求 8

**依赖:** 任务 4, 任务 6, 任务 7

### 子任务

- [x] 8.1 创建 `internal/flow/hooks.go`（回调钩子定义）
- [x] 8.2 创建 `internal/flow/types.go`（Flow 类型定义）
- [x] 8.3 创建 `internal/flow/subdomain_discovery.go`
- [x] 8.4 实现 Stage 1: 被动收集（并行执行）
- [x] 8.5 实现 Stage 2: 字典爆破（可选）
- [x] 8.6 实现 Stage 3: 变异生成 + 验证（可选）
- [x] 8.7 实现 Stage 4: DNS 存活验证（可选）
- [x] 8.8 实现阶段跳过逻辑
- [x] 8.9 实现结果合并去重
- [ ] 8.10 编写 Flow 单元测试

---

## 任务 9: 实现 HTTP API

**需求:** 需求 9

**依赖:** 任务 8

### 子任务

- [x] 9.1 创建 `internal/handler/health.go`
- [x] 9.2 创建 `internal/handler/scan.go`
- [x] 9.3 实现请求验证
- [x] 9.4 实现异步扫描执行
- [ ] 9.5 编写 HTTP API 集成测试

---

## 任务 10: 集成和端到端测试

**需求:** 所有需求

**依赖:** 任务 9

### 子任务

- [ ] 10.1 编写 Flow 集成测试（使用 mock 工具）
- [ ] 10.2 编写完整扫描流程测试
- [ ] 10.3 更新 README 文档
- [ ] 10.4 创建 Dockerfile
