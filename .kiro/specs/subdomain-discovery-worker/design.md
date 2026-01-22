# 设计文档

## 概述

本文档描述 Go 版本子域名发现 Worker 的技术设计。Worker 作为独立服务运行，通过 HTTP API 与 Server 通信，负责执行扫描工具并将结果回传。

## 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              整体架构                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────┐         HTTP API          ┌─────────────────┐
│     Server      │ ────────────────────────▶ │     Worker      │
│   (server/)     │                           │   (worker/)     │
│                 │ ◀──────────────────────── │                 │
│ - 发起扫描       │      状态/结果回传         │ - 接收任务       │
│ - 存储结果       │                           │ - 执行工具       │
│ - 提供 API      │                           │ - 解析结果       │
└─────────────────┘                           └─────────────────┘
```

## 组件和接口

### 项目结构

```
worker/
├── cmd/
│   └── worker/
│       └── main.go              # 入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── client/
│   │   └── server_client.go     # Server API 客户端
│   ├── flow/
│   │   ├── hooks.go             # Flow 回调钩子
│   │   ├── runner.go            # Flow 执行器
│   │   └── subdomain_discovery.go  # 子域名发现 Flow
│   ├── tool/
│   │   ├── command_builder.go   # 命令构建器
│   │   ├── runner.go            # 工具执行器
│   │   └── templates.go         # 命令模板
│   ├── parser/
│   │   └── subdomain.go         # 结果解析器
│   ├── handler/
│   │   ├── scan.go              # 扫描 API 处理器
│   │   └── health.go            # 健康检查
│   └── pkg/
│       └── logger.go            # 日志工具
├── go.mod
├── go.sum
├── Makefile
└── .env.example
```

### 核心组件

#### 1. Config（配置管理）

```go
type Config struct {
    ServerURL        string // Server API 地址
    ServerToken      string // 认证 token
    ScanToolsBasePath string // 扫描工具路径
    ResultsBasePath  string // 结果存储路径
    LogLevel         string // 日志级别
    Port             int    // HTTP 服务端口
}
```

#### 2. ServerClient（Server API 客户端）

```go
type ServerClient interface {
    // 更新扫描状态
    UpdateScanStatus(scanID int, status string, errorMsg string) error
    
    // 批量保存子域名
    SaveSubdomains(scanID int, targetID int, subdomains []string) error
    
    // 写入扫描日志
    WriteScanLog(scanID int, flowName string, message string, level string) error
}
```

#### 3. ToolRunner（工具执行器）

```go
type ToolRunner interface {
    // 执行单个工具
    Run(ctx context.Context, cmd string, timeout time.Duration) (*ToolResult, error)
    
    // 并行执行多个工具
    RunParallel(ctx context.Context, cmds []ToolCommand) []*ToolResult
}

type ToolResult struct {
    Tool       string
    OutputFile string
    ExitCode   int
    Duration   time.Duration
    Error      error
}
```

#### 4. CommandBuilder（命令构建器）

```go
type CommandBuilder interface {
    // 构建扫描命令
    Build(toolName string, scanType string, params map[string]string, config map[string]interface{}) (string, error)
}
```

#### 5. Flow（扫描流程）

```go
type Flow interface {
    Name() string
    Execute(ctx context.Context, input *FlowInput) (*FlowOutput, error)
}

type FlowInput struct {
    ScanID       int
    TargetID     int
    TargetName   string
    WorkspaceDir string
    Config       map[string]interface{}
}

type FlowOutput struct {
    Success        bool
    ProcessedCount int
    FailedTools    []string
    SuccessfulTools []string
}
```

#### 6. FlowHooks（回调钩子）

```go
type FlowHooks struct {
    OnStart    func(scanID int, flowName string)
    OnComplete func(scanID int, flowName string, output *FlowOutput)
    OnFailure  func(scanID int, flowName string, err error)
}
```

### HTTP API

#### POST /api/scans/execute

接收扫描任务。

请求体：
```json
{
    "scanId": 123,
    "targetId": 456,
    "targetName": "example.com",
    "workspaceDir": "/opt/xingrin/results/scan_20260115_123456",
    "config": "subdomain_discovery:\n  passive_tools:\n    subfinder:\n      enabled: true\n..."
}
```

响应：
- 202 Accepted: 任务已接收，异步执行
- 400 Bad Request: 请求参数无效

#### GET /health

健康检查。

响应：
```json
{
    "status": "ok",
    "version": "1.0.0"
}
```

## 数据模型

### 扫描配置（YAML）

沿用 Python 版本的配置格式：

```yaml
subdomain_discovery:
  passive_tools:
    subfinder:
      enabled: true
      timeout: 3600
    sublist3r:
      enabled: true
      timeout: 3600
    assetfinder:
      enabled: true
      timeout: 3600
  bruteforce:
    enabled: false
    subdomain_bruteforce:
      wordlist_name: subdomains-top1million-110000.txt
  permutation:
    enabled: true
    subdomain_permutation_resolve:
      timeout: 7200
  resolve:
    enabled: true
    subdomain_resolve:
      timeout: auto
```

### 命令模板

沿用 Python 版本的模板格式：

```go
var SubdomainDiscoveryCommands = map[string]CommandTemplate{
    "subfinder": {
        Base: "subfinder -d {domain} -all -o '{output_file}' -v",
        Optional: map[string]string{
            "threads":         "-t {threads}",
            "provider_config": "-pc '{provider_config}'",
        },
    },
    "sublist3r": {
        Base: "python3 '/usr/local/share/Sublist3r/sublist3r.py' -d {domain} -o '{output_file}'",
        Optional: map[string]string{
            "threads": "-t {threads}",
        },
    },
    // ... 其他工具
}
```

## 正确性属性

*正确性属性是系统在所有有效执行中都应该保持的特性。每个属性都是一个可以通过属性测试验证的形式化规范。*

### Property 1: 配置验证完整性

*对于任意* 配置输入，如果缺少必需的配置项（SERVER_URL、SERVER_TOKEN），配置加载应当返回错误；如果所有必需项都存在，配置加载应当成功。

**验证: 需求 2.2, 2.3**

### Property 2: 命令构建正确性

*对于任意* 有效的工具名、参数映射和配置，命令构建器应当：
1. 用实际值替换所有占位符
2. 当配置中存在可选参数时，将其追加到命令中
3. 生成的命令应当与 Python 版本兼容

**验证: 需求 4.2, 4.3**

### Property 3: 工具执行错误处理

*对于任意* 工具执行，当工具超时时应当返回超时错误，当工具失败时应当返回退出码和错误信息，正常完成时应当返回输出文件路径。

**验证: 需求 5.4, 5.5**

### Property 4: 并行执行等待

*对于任意* 并行执行的工具集合，执行器应当等待所有工具完成后才返回，返回的结果数量应当等于输入的工具数量。

**验证: 需求 5.6, 6.3**

### Property 5: Flow 阶段跳过

*对于任意* 扫描配置，当某个阶段的 `enabled` 为 false 时，该阶段应当被跳过，不执行任何工具。

**验证: 需求 6.2**

### Property 6: 结果去重

*对于任意* 子域名结果集合，解析和合并后的结果应当不包含重复项，且包含所有唯一的子域名。

**验证: 需求 6.4, 7.2**

### Property 7: API 请求认证

*对于任意* ServerClient 发出的 HTTP 请求，请求头中应当包含 Authorization token。

**验证: 需求 3.4**

### Property 8: 请求验证

*对于任意* 扫描执行请求，如果缺少必需字段（scanId、targetId、targetName、config），应当返回 400 错误。

**验证: 需求 9.2, 9.3**

## 错误处理

### 错误类型

```go
var (
    ErrConfigMissing    = errors.New("missing required configuration")
    ErrToolTimeout      = errors.New("tool execution timeout")
    ErrToolFailed       = errors.New("tool execution failed")
    ErrAPICallFailed    = errors.New("server API call failed")
    ErrInvalidRequest   = errors.New("invalid request")
)
```

### 重试策略

ServerClient API 调用失败时：
1. 最多重试 3 次
2. 使用指数退避：1s, 2s, 4s
3. 所有重试失败后记录错误并继续执行

### 错误传播

```
工具执行失败
    │
    ▼
记录到 FailedTools 列表
    │
    ▼
继续执行其他工具
    │
    ▼
Flow 完成时汇总失败信息
    │
    ▼
回调 Server 更新状态
```

## 测试策略

### 单元测试

- 配置加载和验证
- 命令构建逻辑
- 结果解析和去重
- 请求验证

### 属性测试

使用 `gopter` 库进行属性测试：

- Property 1: 配置验证完整性
- Property 2: 命令构建正确性
- Property 6: 结果去重

### 集成测试

- ServerClient 与 mock server 的交互
- Flow 完整执行流程（使用 mock 工具）
- HTTP API 端点测试

### 测试配置

- 属性测试最少运行 100 次迭代
- 使用 `testify` 进行断言
- 使用 `httptest` 进行 HTTP 测试
