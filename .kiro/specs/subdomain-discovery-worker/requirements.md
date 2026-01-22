# 需求文档

## 简介

本功能实现 Go 版本的子域名发现扫描 Worker，作为独立服务运行，通过 HTTP API 与 Server 通信。Worker 负责执行扫描工具、解析结果、并将结果回传给 Server。

本 spec 聚焦于子域名发现（subdomain_discovery）这一个扫描类型，为后续其他扫描类型的迁移建立基础架构。

## 术语表

- **Worker**: 独立的扫描执行服务，运行在容器中，调用各种扫描工具
- **Server**: Go 后端 API 服务（`server/`），提供数据存储和 API 接口
- **Flow**: 一个扫描类型的完整执行流程（如 subdomain_discovery_flow）
- **Stage**: Flow 内部的执行阶段（如被动收集、字典爆破）
- **Tool**: 具体的扫描工具（如 subfinder、assetfinder）
- **ToolRunner**: 负责构建命令、执行工具、捕获输出的组件
- **ResultParser**: 负责解析工具输出的组件
- **ServerClient**: Worker 调用 Server API 的 HTTP 客户端

## 需求

### 需求 1: Worker 项目初始化

**用户故事:** 作为开发者，我希望创建一个独立的 Go 模块用于 Worker，以便它可以独立于 Server 进行开发和部署。

#### 验收标准

1. Worker 应当作为独立的 Go 模块创建在 `worker/` 目录下
2. Worker 应当有自己的 `go.mod`，模块名为 `github.com/xingrin/worker`
3. Worker 应当遵循与 Server 相同的项目结构（`cmd/`、`internal/`）
4. Worker 应当使用与 Server 相同的技术栈（Gin 用于 HTTP、Zap 用于日志、Viper 用于配置）

### 需求 2: 配置管理

**用户故事:** 作为开发者，我希望 Worker 能从环境变量和配置文件加载配置，以便在不同环境中轻松配置。

#### 验收标准

1. Worker 应当从环境变量和 `.env` 文件加载配置
2. Worker 应当支持以下配置项：
   - `SERVER_URL`: Server API 地址（如 `http://server:8888`）
   - `SERVER_TOKEN`: 访问 Server API 的认证 token
   - `SCAN_TOOLS_BASE_PATH`: 扫描工具的基础路径（默认 `/usr/local/bin`）
   - `RESULTS_BASE_PATH`: 扫描结果的基础路径（默认 `/opt/xingrin/results`）
   - `LOG_LEVEL`: 日志级别（默认 `info`）
3. 当缺少必需的配置项时，Worker 应当启动失败并显示清晰的错误信息

### 需求 3: Server API 客户端

**用户故事:** 作为 Worker，我希望通过 HTTP API 与 Server 通信，以便接收扫描任务并提交结果。

#### 验收标准

1. ServerClient 应当支持以下操作：
   - 更新扫描状态（running/completed/failed）
   - 批量保存子域名结果
   - 写入扫描日志
2. 当 API 调用失败时，ServerClient 应当使用指数退避策略重试最多 3 次
3. 当所有重试都失败时，ServerClient 应当记录错误并继续执行
4. ServerClient 应当在所有请求中包含认证 token

### 需求 4: 命令模板系统

**用户故事:** 作为开发者，我希望为每个扫描工具定义命令模板，以便根据配置动态构建命令。

#### 验收标准

1. Worker 应当为子域名发现工具定义命令模板：
   - subfinder
   - sublist3r
   - assetfinder
   - subdomain_bruteforce（puredns）
   - subdomain_resolve（puredns）
   - subdomain_permutation_resolve（dnsgen + puredns）
2. 当构建命令时，Worker 应当用实际值替换占位符
3. 当配置中提供了可选参数时，Worker 应当将其追加到命令中
4. 命令模板格式应当与现有 Python 实现兼容

### 需求 5: 工具执行器（ToolRunner）

**用户故事:** 作为 Worker，我希望执行扫描工具并捕获其输出，以便收集扫描结果。

#### 验收标准

1. ToolRunner 应当执行带有可配置超时的 shell 命令
2. ToolRunner 应当捕获 stdout 和 stderr
3. ToolRunner 应当将工具输出写入日志文件
4. 当工具超时时，ToolRunner 应当终止进程并返回超时错误
5. 当工具失败时，ToolRunner 应当返回退出码和错误信息
6. ToolRunner 应当支持多个工具的并行执行

### 需求 6: 子域名发现 Flow

**用户故事:** 作为 Worker，我希望执行包含多个阶段的子域名发现流程，以便全面发现子域名。

#### 验收标准

1. subdomain_discovery_flow 应当支持 4 个阶段：
   - Stage 1: 被动收集（并行执行 subfinder、sublist3r、assetfinder）
   - Stage 2: 字典爆破（可选，使用 puredns bruteforce）
   - Stage 3: 变异生成 + 验证（可选，使用 dnsgen + puredns resolve）
   - Stage 4: DNS 存活验证（可选，使用 puredns resolve）
2. 当某个阶段在配置中被禁用时，Flow 应当跳过该阶段
3. 当 Stage 1 的工具并行运行时，Flow 应当等待所有工具完成后再继续
4. 当多个阶段产生结果时，Flow 应当合并并去重
5. Flow 应当在所有阶段完成后调用 Server API 保存结果

### 需求 7: 结果解析与保存

**用户故事:** 作为 Worker，我希望解析工具输出并将结果保存到 Server，以便发现的子域名被持久化。

#### 验收标准

1. Worker 应当从工具输出文件解析子域名结果（每行一个域名）
2. Worker 应当在保存前对子域名去重
3. Worker 应当调用 Server API 批量保存子域名
4. 保存结果时，Worker 应当包含 scan_id 和 target_id

### 需求 8: 状态回调机制

**用户故事:** 作为 Worker，我希望向 Server 报告扫描状态，以便用户可以跟踪扫描进度。

#### 验收标准

1. 当 Flow 开始时，Worker 应当调用 Server API 将状态更新为 "running"
2. 当 Flow 成功完成时，Worker 应当调用 Server API 将状态更新为 "completed"
3. 当 Flow 失败时，Worker 应当调用 Server API 将状态更新为 "failed" 并附带错误信息
4. Worker 应当通过 API 将扫描日志写入 Server 以便用户查看

### 需求 9: HTTP API 接口

**用户故事:** 作为 Server，我希望通过 HTTP API 在 Worker 上触发扫描，以便编排扫描执行。

#### 验收标准

1. Worker 应当暴露 POST `/api/scans/execute` 端点来接收扫描任务
2. 请求体应当包含：
   - scan_id: 扫描 ID
   - target_id: 目标 ID
   - target_name: 目标名称（域名）
   - workspace_dir: 工作目录路径
   - config: 扫描配置（YAML 格式）
3. Worker 应当验证请求，如果无效则返回 400
4. Worker 应当异步执行扫描并立即返回 202 Accepted
5. Worker 应当暴露 GET `/health` 端点用于健康检查

### 需求 10: 错误处理与日志

**用户故事:** 作为开发者，我希望有全面的错误处理和日志记录，以便轻松调试问题。

#### 验收标准

1. Worker 应当记录所有工具执行，包括命令、耗时和退出码
2. Worker 应当记录所有对 Server 的 API 调用，包括请求/响应详情
3. 当发生错误时，Worker 应当记录带有上下文的错误（scan_id、tool_name 等）
4. Worker 应当在生产环境使用 JSON 格式的结构化日志
5. Worker 应当在开发环境支持人类可读的日志
