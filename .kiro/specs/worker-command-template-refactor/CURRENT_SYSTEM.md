# 当前系统文档：Worker 命令模板系统

> 本文档记录了重构前的系统架构、代码实现和配置方式，供后续 AI 研究和对比。
> 
> 创建时间：2026-01-17
> 
> 状态：待重构

---

## 目录

1. [系统概述](#系统概述)
2. [架构设计](#架构设计)
3. [核心组件](#核心组件)
4. [配置系统](#配置系统)
5. [数据流](#数据流)
6. [代码示例](#代码示例)
7. [存在的问题](#存在的问题)

---

## 系统概述

### 功能

Worker 命令模板系统负责：
1. 定义扫描工具的命令模板
2. 根据用户配置构建实际执行的命令
3. 管理参数的占位符替换
4. 支持可选参数的动态添加

### 技术栈

- **语言**: Go 1.21+
- **配置格式**: YAML
- **模板引擎**: 字符串替换（无第三方库）
- **嵌入资源**: `embed.FS`

---

## 架构设计

### 系统组件图

```
┌─────────────────────────────────────────────────────────────┐
│  Server 配置 (subdomain_discovery.yaml)                     │
│  - 用户定义启用哪些工具                                      │
│  - 用户覆盖参数值（timeout, threads 等）                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 传递配置
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Workflow (subdomain_discovery)                             │
│  - 读取 Server 配置                                         │
│  - 调用各个 stage 执行扫描                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 调用
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Stage (stage_passive.go, stage_bruteforce.go)             │
│  - 遍历启用的工具                                            │
│  - 调用 buildCommand() 构建命令                             │
│  - 调用 Runner 执行命令                                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 调用
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  buildCommand() (helpers.go)                                │
│  - 获取模板                                                  │
│  - 调用 CommandBuilder.Build()                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 调用
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  TemplateLoader (template_loader.go)                        │
│  - 从 embed.FS 加载 templates.yaml                          │
│  - 使用 sync.Once 缓存模板                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 返回模板
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  CommandBuilder (command_builder.go)                        │
│  - 替换必需占位符                                            │
│  - 添加可选参数                                              │
│  - 返回最终命令字符串                                         │
└─────────────────────────────────────────────────────────────┘
```

### 两层配置系统

系统使用两层配置：

1. **Worker 模板层** (`templates.yaml`)
   - 开发者定义
   - 包含命令模板和可选参数映射
   - 嵌入到 Worker 二进制文件中

2. **Server 配置层** (`subdomain_discovery.yaml`)
   - 用户定义
   - 控制启用哪些工具
   - 覆盖参数值（如 timeout, threads）

---

## 核心组件

### 1. CommandTemplate 结构体

**文件**: `worker/internal/activity/command_template.go`

```go
package activity

// CommandTemplate defines a command template for an activity
type CommandTemplate struct {
	Base     string            `yaml:"base"`     // Base command with required placeholders
	Optional map[string]string `yaml:"optional"` // Optional parameters and their flags
}
```

**特点**:
- 扁平结构
- `Base`: 基础命令，包含必需占位符（如 `{domain}`, `{output-file}`）
- `Optional`: 可选参数的标志模板（如 `"-t {threads}"`）
- **没有默认值字段** - 默认值硬编码在命令字符串中

### 2. CommandBuilder

**文件**: `worker/internal/activity/command_builder.go`

```go
type CommandBuilder struct{}

func (b *CommandBuilder) Build(
    tmpl CommandTemplate, 
    params map[string]string, 
    config map[string]any,
) (string, error) {
    // 1. 从 base 开始
    cmd := tmpl.Base
    
    // 2. 替换必需占位符
    for key, value := range params {
        placeholder := "{" + key + "}"
        cmd = strings.ReplaceAll(cmd, placeholder, value)
    }
    
    // 3. 添加可选参数（如果用户配置中存在）
    for configKey, flagTemplate := range tmpl.Optional {
        if value, ok := getConfigValue(config, configKey); ok {
            flag := strings.ReplaceAll(flagTemplate, "{"+configKey+"}", fmt.Sprintf("%v", value))
            cmd = cmd + " " + flag
        }
    }
    
    // 4. 检查未替换的占位符
    if strings.Contains(cmd, "{") && strings.Contains(cmd, "}") {
        return "", fmt.Errorf("command contains unreplaced placeholders: %s", cmd)
    }
    
    return cmd, nil
}
```

**工作流程**:
1. 从 `Base` 命令开始
2. 替换必需占位符（如 `{domain}` → `example.com`）
3. 如果用户配置中有可选参数，添加对应的标志
4. 检查是否有未替换的占位符
5. 返回最终命令字符串

### 3. TemplateLoader

**文件**: `worker/internal/workflow/subdomain_discovery/template_loader.go`

```go
package subdomain_discovery

import (
	"embed"
	"github.com/orbit/worker/internal/activity"
)

//go:embed templates.yaml
var templatesFS embed.FS

// loader is the template loader for subdomain discovery workflow
var loader = activity.NewTemplateLoader(templatesFS, "templates.yaml")

// getTemplate returns the command template for a given tool
func getTemplate(toolName string) (activity.CommandTemplate, error) {
	return loader.Get(toolName)
}
```

**特点**:
- 使用 `embed.FS` 嵌入 YAML 文件
- 使用 `sync.Once` 缓存已加载的模板
- 提供简单的 `getTemplate()` 接口

### 4. buildCommand 辅助函数

**文件**: `worker/internal/workflow/subdomain_discovery/helpers.go`

```go
func buildCommand(toolName string, params map[string]string, config map[string]any) (string, error) {
	tmpl, err := getTemplate(toolName)
	if err != nil {
		return "", err
	}
	builder := activity.NewCommandBuilder()
	return builder.Build(tmpl, params, config)
}
```

**作用**:
- 封装模板获取和命令构建的流程
- 被各个 stage 调用

---

## 配置系统

### Worker 模板配置

**文件**: `worker/internal/workflow/subdomain_discovery/templates.yaml`

```yaml
# 被动收集工具
subfinder:
  base: "subfinder -d {domain} -all -o '{output-file}' -v"
  optional:
    threads: "-t {threads}"
    provider-config: "-pc '{provider-config}'"
    timeout: "-timeout {timeout}"

sublist3r:
  base: "python3 '/usr/local/share/Sublist3r/sublist3r.py' -d {domain} -o '{output-file}'"
  optional:
    threads: "-t {threads}"

assetfinder:
  base: "assetfinder --subs-only {domain} > '{output-file}'"
  optional: {}

# 主动扫描工具
subdomain-bruteforce:
  base: "puredns bruteforce '{wordlist}' {domain} -r '{resolvers}' --write '{output-file}' --quiet"
  optional:
    threads: "-t {threads}"
    rate-limit: "--rate-limit {rate-limit}"
    wildcard-tests: "--wildcard-tests {wildcard-tests}"
    wildcard-batch: "--wildcard-batch {wildcard-batch}"

subdomain-resolve:
  base: "puredns resolve '{input-file}' -r '{resolvers}' --write '{output-file}' --wildcard-tests 50 --wildcard-batch 1000000 --quiet"
  optional:
    threads: "-t {threads}"
    rate-limit: "--rate-limit {rate-limit}"
    wildcard-tests: "--wildcard-tests {wildcard-tests}"
    wildcard-batch: "--wildcard-batch {wildcard-batch}"
```

**特点**:
- 扁平结构：`base` 和 `optional` 分开
- 默认值硬编码在 `base` 中（如 `--wildcard-tests 50`）
- 没有参数类型定义
- 没有参数描述

### Server 配置

**文件**: `server/configs/engines/subdomain_discovery.yaml`

```yaml
# Stage 1: Passive Collection
passive-tools:
  subfinder:
    enabled: true
    timeout: 3600  # 覆盖默认值
    # threads: 10  # 可选，注释掉表示使用默认值

  sublist3r:
    enabled: true
    timeout: 3600

  assetfinder:
    enabled: true
    timeout: 3600

# Stage 2: Dictionary Bruteforce
bruteforce:
  enabled: false
  subdomain-bruteforce:
    timeout: 86400
    wordlist-name: subdomains-top1million-110000.txt

# Stage 3: Permutation + Resolve
permutation:
  enabled: true
  subdomain-permutation-resolve:
    timeout: 86400

# Stage 4: DNS Resolution Validation
resolve:
  enabled: true
  subdomain-resolve:
    timeout: 86400
```

**特点**:
- 扁平结构
- 使用 `enabled` 控制工具是否启用
- 用户可以覆盖参数值
- 没有参数说明（用户不知道有哪些参数可配置）

---

## 数据流

### 完整的命令构建流程

```
1. 用户配置
   ↓
   server/configs/engines/subdomain_discovery.yaml
   {
     "passive-tools": {
       "subfinder": {
         "enabled": true,
         "timeout": 3600,
         "threads": 20
       }
     }
   }

2. Workflow 读取配置
   ↓
   stage_passive.go: runPassiveStage()
   - 遍历 passive-tools
   - 检查 enabled: true
   - 调用 createPassiveCommand()

3. 构建命令参数
   ↓
   stage_passive.go: createPassiveCommand()
   params = {
     "domain": "example.com",
     "output-file": "/path/to/output.txt"
   }
   config = {
     "timeout": 3600,
     "threads": 20
   }

4. 调用 buildCommand()
   ↓
   helpers.go: buildCommand()
   - 调用 getTemplate("subfinder")
   - 调用 CommandBuilder.Build()

5. 获取模板
   ↓
   template_loader.go: getTemplate()
   返回:
   {
     "base": "subfinder -d {domain} -all -o '{output-file}' -v",
     "optional": {
       "threads": "-t {threads}",
       "timeout": "-timeout {timeout}"
     }
   }

6. 构建命令
   ↓
   command_builder.go: Build()
   
   步骤 1: 从 base 开始
   cmd = "subfinder -d {domain} -all -o '{output-file}' -v"
   
   步骤 2: 替换必需占位符
   cmd = "subfinder -d example.com -all -o '/path/to/output.txt' -v"
   
   步骤 3: 添加可选参数
   - config 中有 "timeout": 3600
   - 添加 "-timeout 3600"
   - config 中有 "threads": 20
   - 添加 "-t 20"
   
   cmd = "subfinder -d example.com -all -o '/path/to/output.txt' -v -timeout 3600 -t 20"
   
   步骤 4: 检查未替换占位符
   - 没有 "{...}" 格式的字符串
   - 验证通过

7. 返回最终命令
   ↓
   "subfinder -d example.com -all -o '/path/to/output.txt' -v -timeout 3600 -t 20"

8. 执行命令
   ↓
   Runner.RunParallel()
```

### 参数优先级

当前系统的参数值来源：

1. **硬编码默认值** (最低优先级)
   - 在 `base` 命令中硬编码
   - 例如：`--wildcard-tests 50`

2. **用户配置** (最高优先级)
   - 在 Server 配置文件中定义
   - 例如：`timeout: 3600`

**问题**: 如果默认值在 `base` 中硬编码，用户无法覆盖它们（除非修改 Worker 模板）

---

## 代码示例

### 示例 1: 构建 subfinder 命令

```go
// 输入
toolName := "subfinder"
params := map[string]string{
    "domain":      "example.com",
    "output-file": "/tmp/subfinder_output.txt",
}
config := map[string]any{
    "timeout": 3600,
    "threads": 20,
}

// 调用
cmd, err := buildCommand(toolName, params, config)

// 输出
// cmd = "subfinder -d example.com -all -o '/tmp/subfinder_output.txt' -v -timeout 3600 -t 20"
```

### 示例 2: 构建 puredns resolve 命令

```go
// 输入
toolName := "subdomain-resolve"
params := map[string]string{
    "input-file":  "/tmp/subdomains.txt",
    "output-file": "/tmp/resolved.txt",
    "resolvers":   "/etc/resolvers.txt",
}
config := map[string]any{
    "threads": 100,
}

// 调用
cmd, err := buildCommand(toolName, params, config)

// 输出
// cmd = "puredns resolve '/tmp/subdomains.txt' -r '/etc/resolvers.txt' --write '/tmp/resolved.txt' --wildcard-tests 50 --wildcard-batch 1000000 --quiet -t 100"
// 
// 注意：--wildcard-tests 50 和 --wildcard-batch 1000000 是硬编码的，无法覆盖
```

### 示例 3: Stage 执行流程

```go
// stage_passive.go
func (w *Workflow) runPassiveStage(ctx *workflowContext) stageResult {
    // 1. 获取 stage 配置
    stageConfig, ok := ctx.config[stagePassive].(map[string]any)
    if !ok {
        return stageResult{}
    }

    var commands []activity.Command

    // 2. 遍历所有域名
    for _, domain := range ctx.domains {
        // 3. 遍历所有被动工具
        for _, toolName := range passiveTools {
            // 4. 检查工具是否启用
            if !isToolEnabled(stageConfig, toolName) {
                continue
            }

            // 5. 获取工具配置
            toolConfig, _ := stageConfig[toolName].(map[string]any)
            
            // 6. 创建命令
            cmd := w.createPassiveCommand(ctx, domain, toolName, toolConfig)
            if cmd != nil {
                commands = append(commands, *cmd)
            }
        }
    }

    // 7. 并行执行所有命令
    results := w.runner.RunParallel(ctx.ctx, commands)
    return processResults(results)
}
```

---

## 存在的问题

### 1. 默认值硬编码

**问题**: 默认参数值直接写在 `base` 命令中

```yaml
subdomain-resolve:
  base: "puredns resolve '{input-file}' -r '{resolvers}' --write '{output-file}' --wildcard-tests 50 --wildcard-batch 1000000 --quiet"
```

**影响**:
- 用户无法覆盖这些默认值
- 修改默认值需要修改 Worker 模板并重新编译
- 默认值不可见（用户不知道有哪些默认值）

### 2. 配置重复

**问题**: 多个工具有相同的参数配置，但需要重复定义

```yaml
# 每个工具都要定义 threads, timeout
subfinder:
  optional:
    threads: "-t {threads}"
    timeout: "-timeout {timeout}"

sublist3r:
  optional:
    threads: "-t {threads}"
    timeout: "-timeout {timeout}"  # 重复

subdomain-bruteforce:
  optional:
    threads: "-t {threads}"  # 重复
    rate-limit: "--rate-limit {rate-limit}"
    wildcard-tests: "--wildcard-tests {wildcard-tests}"
    wildcard-batch: "--wildcard-batch {wildcard-batch}"

subdomain-resolve:
  optional:
    threads: "-t {threads}"  # 重复
    rate-limit: "--rate-limit {rate-limit}"  # 重复
    wildcard-tests: "--wildcard-tests {wildcard-tests}"  # 重复
    wildcard-batch: "--wildcard-batch {wildcard-batch}"  # 重复
```

**影响**:
- 修改共享配置需要改多处
- 容易出现不一致
- 维护成本高

### 3. 缺少验证

**问题**: 模板加载时没有验证

```go
// template_loader.go 只是简单加载，没有验证
func (l *TemplateLoader) Load() (map[string]CommandTemplate, error) {
    // 读取 YAML
    data, err := l.fs.ReadFile(l.filename)
    
    // 解析 YAML
    if err := yaml.Unmarshal(data, &l.cache); err != nil {
        return nil, err
    }
    
    // 没有验证！
    return l.cache, nil
}
```

**影响**:
- 模板错误只在运行时才能发现
- 错误信息不清晰
- 难以调试

### 4. 类型不明确

**问题**: 参数类型未定义

```yaml
subfinder:
  optional:
    threads: "-t {threads}"  # threads 是什么类型？int? string?
    timeout: "-timeout {timeout}"  # timeout 是什么类型？
```

**影响**:
- 容易出现类型错误
- 没有类型验证
- 用户不知道应该传什么类型的值

### 5. 错误信息不清晰

**问题**: 构建失败时错误信息简单

```go
if strings.Contains(cmd, "{") && strings.Contains(cmd, "}") {
    return "", fmt.Errorf("command contains unreplaced placeholders: %s", cmd)
}
```

**影响**:
- 不知道哪个占位符未替换
- 不知道缺少哪个参数
- 难以快速定位问题

### 6. 缺少参数文档

**问题**: Server 配置文件没有参数说明

```yaml
passive-tools:
  subfinder:
    enabled: true
    timeout: 3600
    # 用户不知道还有哪些参数可以配置
    # 用户不知道参数的含义和取值范围
```

**影响**:
- 用户需要查看 Worker 模板才知道有哪些参数
- 没有参数说明和默认值
- 配置困难

### 7. 配置结构不一致

**问题**: Worker 模板和 Server 配置结构不同

```yaml
# Worker 模板（扁平）
subfinder:
  base: "..."
  optional:
    threads: "-t {threads}"

# Server 配置（也是扁平，但结构不同）
passive-tools:
  subfinder:
    enabled: true
    threads: 20
```

**影响**:
- 两层配置难以对应
- 用户不知道 Worker 模板中定义了什么
- 维护困难

---

## 总结

### 当前系统的优点

1. ✅ **简单直观**: 扁平结构，易于理解
2. ✅ **性能良好**: 使用 `sync.Once` 缓存，避免重复加载
3. ✅ **灵活**: 支持任意参数的动态添加

### 当前系统的缺点

1. ❌ **默认值硬编码**: 无法覆盖，不易维护
2. ❌ **配置重复**: 大量重复定义，维护成本高
3. ❌ **缺少验证**: 错误只在运行时发现
4. ❌ **类型不明确**: 没有类型定义和验证
5. ❌ **错误信息简单**: 难以快速定位问题
6. ❌ **缺少文档**: 用户不知道有哪些参数可配置
7. ❌ **结构不一致**: Worker 模板和 Server 配置难以对应

### 重构目标

基于以上问题，重构应该实现：

1. ✅ 默认值在模板中明确定义
2. ✅ 使用 YAML 锚点消除重复
3. ✅ 启动时验证所有模板
4. ✅ 明确的参数类型定义
5. ✅ 详细的错误信息
6. ✅ 自动生成配置文档
7. ✅ 统一的配置结构（符合业界标准）

---

## 附录

### 相关文件清单

**Worker 端**:
- `worker/internal/activity/command_template.go` - 模板结构定义
- `worker/internal/activity/command_builder.go` - 命令构建器
- `worker/internal/activity/template_loader.go` - 通用模板加载器
- `worker/internal/workflow/subdomain_discovery/templates.yaml` - 工具模板
- `worker/internal/workflow/subdomain_discovery/template_loader.go` - Workflow 特定加载器
- `worker/internal/workflow/subdomain_discovery/helpers.go` - 辅助函数
- `worker/internal/workflow/subdomain_discovery/stage_passive.go` - 被动收集阶段
- `worker/internal/workflow/subdomain_discovery/stage_bruteforce.go` - 爆破阶段

**Server 端**:
- `server/configs/engines/subdomain_discovery.yaml` - 用户配置

### 关键常量

```go
// worker/internal/workflow/subdomain_discovery/workflow.go
const (
    // Stage names
    stagePassive     = "passive-tools"
    stageBruteforce  = "bruteforce"
    stagePermutation = "permutation"
    stageResolve     = "resolve"
    
    // Tool names
    toolSubfinder              = "subfinder"
    toolSublist3r              = "sublist3r"
    toolAssetfinder            = "assetfinder"
    toolSubdomainBruteforce    = "subdomain-bruteforce"
    toolSubdomainResolve       = "subdomain-resolve"
    toolSubdomainPermutation   = "subdomain-permutation-resolve"
    
    // Default timeout
    defaultTimeout = 86400 // 24 hours
)
```

---

**文档结束**
