# Design Document

## Overview

本文档描述 Worker 命令模板系统的重构设计，采用业界标准方案实现配置管理、验证和文档生成。

### 当前问题

1. **默认值硬编码**: 默认参数值直接写在命令字符串中（如 `-timeout 3600`），难以维护
2. **配置重复**: 多个工具的相同参数（如 timeout）需要重复定义
3. **缺少验证**: 模板加载时不验证，运行时才发现错误
4. **错误信息不清晰**: 命令构建失败时难以定位问题
5. **文档不足**: 缺少参数说明和示例
6. **配置验证缺失**: Server 端无法验证用户配置的正确性

### 设计目标

1. **集中管理默认值**: 在 Worker 模板中明确定义所有默认值
2. **消除重复配置**: 使用 YAML 锚点共享配置
3. **启动时验证**: 加载模板时验证所有配置
4. **清晰的错误信息**: 提供详细的错误上下文
5. **简化配置**: 支持 `enabled: true` 即可使用默认值
6. **自动生成 Schema**: 从 Worker 模板生成 JSON Schema
7. **自动生成文档**: 从 Worker 模板生成配置参考文档
8. **Server 端验证**: 使用 JSON Schema 验证用户配置

### 业界参考

本设计参考以下业界标准实践：

- **Helm**: `helm-schema-gen` - 从 values.yaml 生成 values.schema.json
- **Terraform**: `terraform-plugin-docs` - 从 Provider Schema 生成文档
- **GitHub Actions**: action.yml 的 inputs 定义格式
- **Cloudflare**: 从 OpenAPI 生成 Terraform Provider

## Architecture

### 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│  开发阶段（go generate）                                      │
├─────────────────────────────────────────────────────────────┤
│  Worker 模板 (templates.yaml) [单一数据源]                   │
│      ↓                                                       │
│  SchemaGenerator → config.schema.json                       │
│  DocGenerator    → config-reference.md                      │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│  Server 启动                                                 │
├─────────────────────────────────────────────────────────────┤
│  ConfigValidator.Load(config.schema.json)                   │
│  ConfigValidator.Validate(user_config.yaml)                 │
│      ↓ 验证通过                                              │
│  启动成功                                                    │
└─────────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────────┐
│  Worker 启动                                                 │
├─────────────────────────────────────────────────────────────┤
│  TemplateLoader.Load(templates.yaml)                        │
│      ├─ 解析 YAML（支持锚点）                                │
│      ├─ 验证模板结构                                         │
│      └─ 缓存到内存（sync.Once）                              │
└─────────────────────────────────────────────────────────────┘

                        ↓
┌─────────────────────────────────────────────────────────────┐
│  扫描任务执行                                                │
├─────────────────────────────────────────────────────────────┤
│  CommandBuilder.Build(toolName, userConfig)                 │
│      ├─ 获取缓存的模板                                       │
│      ├─ 合并默认值和用户配置                                  │
│      ├─ 验证参数类型                                         │
│      ├─ 替换占位符                                           │
│      └─ 返回最终命令                                         │
└─────────────────────────────────────────────────────────────┘
```

### 数据流

```
Worker 模板 (templates.yaml) [单一数据源]
    ↓ go generate
    ├── config.schema.json (JSON Schema)
    └── config-reference.md (配置文档)
    ↓ 用于验证
Server 配置 (subdomain_discovery.yaml) [用户编写]
    ↓ 传递配置
Worker 执行
```

### 关键设计决策

1. **单一数据源**: Worker 模板是唯一的真实来源，Schema 和文档都从它生成
2. **自动化生成**: 使用 `go generate` 自动生成 Schema 和文档，避免手动维护
3. **早期验证**: Server 启动时验证配置，Worker 启动时验证模板，快速失败
4. **保持简单**: 用户配置保持扁平结构，不需要嵌套

## Components and Interfaces

### 1. CommandTemplate 结构体

采用嵌套结构，每个参数包含所有属性，使用 Go Template 语法：

```go
// CommandTemplate 定义工具的命令模板（使用 Go Template 语法）
type CommandTemplate struct {
    Metadata    ToolMetadata          `yaml:"metadata"`
    BaseCommand string                `yaml:"base_command"` // 使用 {{.Var}} 占位符
    Parameters  map[string]Parameter  `yaml:"parameters"`
}

// Parameter 定义单个参数的所有属性
type Parameter struct {
    Flag               string      `yaml:"flag"`        // 使用 {{.Var}} 占位符
    Default            interface{} `yaml:"default"`
    Type               string      `yaml:"type"`        // "string", "int", "bool"
    Required           bool        `yaml:"required"`
    Description        string      `yaml:"description"`
    DeprecationMessage string      `yaml:"deprecation_message,omitempty"`
}

// ToolMetadata 定义工具的元数据
type ToolMetadata struct {
    DisplayName string `yaml:"display_name"`
    Description string `yaml:"description"`
    Stage       string `yaml:"stage"`              // 所属阶段（必需）
    Warning     string `yaml:"warning,omitempty"`
}

// 说明：
// - Stage: 工具所属的扫描阶段（如 passive, bruteforce）
//   **关键作用**: 工作流代码通过 Stage 字段动态发现可用工具，实现模板驱动
//   **必须**: 引用 WorkflowMetadata.Stages 中已定义的阶段 ID
// - Warning: 警告信息（如主动扫描警告、API Keys 需求等）
// 
// 注意：
// - 阶段依赖关系在 StageMetadata.DependsOn 中定义，不在工具级别定义
// - API Keys 相关信息通过 Description 或 Warning 字段说明
// 
// 动态工具发现示例：
// ```go
// // ❌ 错误 - 硬编码工具名称
// tools := []string{"subfinder", "sublist3r", "assetfinder"}
// 
// // ✅ 正确 - 从模板动态获取
// tools := getToolsByStage("passive")  // 自动发现所有 passive 阶段的工具
// ```
```

**设计理由**:
- 使用 Go 原生 `text/template`，代码量更少（~50 行 vs ~140 行）
- 自动检测缺失字段（`missingkey=error`）
- 支持高级特性（条件、循环、函数）
- 符合业界标准（Helm、Kubernetes）
- 更好的错误信息（Go template 提供详细的错误位置）

**占位符格式**:
- 必需字段：`{{.Domain}}`, `{{.OutputFile}}`（PascalCase）
- 可选字段：`{{.Timeout}}`, `{{.Threads}}`
- 函数调用：`{{quote .Domain}}`, `{{default 3600 .Timeout}}`

### 2. WorkflowMetadata 结构体

定义 Workflow 的整体元数据：

```go
// WorkflowMetadata 定义 Workflow 的元数据
type WorkflowMetadata struct {
    Name         string          `yaml:"name"`
    DisplayName  string          `yaml:"display_name"`
    Description  string          `yaml:"description"`
    Version      string          `yaml:"version"`
    TargetTypes  []string        `yaml:"target_types"`
    Stages       []StageMetadata `yaml:"stages"`
}

// StageMetadata 定义阶段的元数据
type StageMetadata struct {
    ID          string   `yaml:"id"`
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    Order       int      `yaml:"order"`
    Required    bool     `yaml:"required"`
    Parallel    bool     `yaml:"parallel"`
    DependsOn   []string `yaml:"depends_on"`
    Outputs     []string `yaml:"outputs"`
}
```

### 3. TemplateSource 接口

定义模板数据源的抽象接口，支持多种加载方式：

```go
// worker/internal/activity/template_source.go

// TemplateSource 定义模板数据源接口
type TemplateSource interface {
    // Read 读取模板数据
    Read() ([]byte, error)
}

// EmbedSource 从嵌入文件读取（当前实现）
type EmbedSource struct {
    fs       embed.FS
    filename string
}

func (s *EmbedSource) Read() ([]byte, error) {
    return s.fs.ReadFile(s.filename)
}

// MemorySource 从内存读取（用于测试）
type MemorySource struct {
    data []byte
}

func (s *MemorySource) Read() ([]byte, error) {
    return s.data, nil
}
```

**设计说明**:
- 接口只有一个方法 `Read()`，保持简单
- 当前只实现 `EmbedSource`（编译时嵌入）和 `MemorySource`（测试用）
- 未来可以添加 `FileSource`、`URLSource`、`DBSource` 等实现
- 使用接口抽象，易于扩展和测试

### 4. TemplateLoader

负责加载和验证 Worker 模板（使用 TemplateSource 接口）：

```go
// worker/internal/activity/template_loader.go

// TemplateLoader 加载和缓存命令模板
type TemplateLoader struct {
    source   TemplateSource             // 数据源接口
    once     sync.Once                  // 保证只加载一次
    cache    map[string]CommandTemplate // 缓存的模板
    metadata WorkflowMetadata           // 工作流元数据
    err      error                      // 加载错误
}

// NewTemplateLoader 创建模板加载器
func NewTemplateLoader(source TemplateSource) *TemplateLoader {
    return &TemplateLoader{
        source: source,
    }
}

// Load 加载模板（使用 sync.Once 缓存）
func (l *TemplateLoader) Load() (map[string]CommandTemplate, error) {
    l.once.Do(func() {
        data, err := l.source.Read()  // 使用接口方法
        if err != nil {
            l.err = fmt.Errorf("failed to read template source: %w", err)
            pkg.Logger.Error("Failed to read template source", zap.Error(l.err))
            return
        }
        
        // 解析 YAML
        l.cache = make(map[string]CommandTemplate)
        if err := yaml.Unmarshal(data, &l.cache); err != nil {
            l.err = fmt.Errorf("failed to parse templates: %w", err)
            pkg.Logger.Error("Failed to parse templates", zap.Error(l.err))
            return
        }
        
        // 验证模板
        if err := l.validate(); err != nil {
            l.err = err
            pkg.Logger.Error("Failed to validate templates", zap.Error(l.err))
            return
        }
        
        pkg.Logger.Info("Templates loaded successfully",
            zap.Int("count", len(l.cache)))
    })
    
    return l.cache, l.err
}

// Get 获取指定工具的模板
func (l *TemplateLoader) Get(name string) (CommandTemplate, error) {
    templates, err := l.Load()
    if err != nil {
        return CommandTemplate{}, fmt.Errorf("templates not loaded: %w", err)
    }
    
    tmpl, ok := templates[name]
    if !ok {
        return CommandTemplate{}, fmt.Errorf("template not found: %s", name)
    }
    
    return tmpl, nil
}

// GetMetadata 获取 Workflow 元数据
func (l *TemplateLoader) GetMetadata() WorkflowMetadata {
    return l.metadata
}

// validate 验证所有模板的正确性（内部方法）
func (l *TemplateLoader) validate() error {
    for name, tmpl := range l.cache {
        if tmpl.BaseCommand == "" {
            return fmt.Errorf("template %s: base_command is required", name)
        }
        
        // 验证 Go Template 语法
        if _, err := template.New(name).Parse(tmpl.BaseCommand); err != nil {
            return fmt.Errorf("template %s: invalid base_command syntax: %w", name, err)
        }
        
        // 验证 Metadata
        if tmpl.Metadata.Stage == "" {
            return fmt.Errorf("template %s: metadata.stage is required", name)
        }
        
        // 验证 Stage 引用的阶段存在
        stageExists := false
        for _, stage := range l.metadata.Stages {
            if stage.ID == tmpl.Metadata.Stage {
                stageExists = true
                break
            }
        }
        if !stageExists {
            return fmt.Errorf("template %s: stage %s not found in workflow metadata", 
                name, tmpl.Metadata.Stage)
        }
        
        // 验证参数
        for paramName, param := range tmpl.Parameters {
            // 验证类型
            if param.Type != "string" && param.Type != "int" && param.Type != "bool" {
                return fmt.Errorf("template %s, parameter %s: invalid type %s (must be string/int/bool)",
                    name, paramName, param.Type)
            }
            
            // 验证默认值类型
            if param.Default != nil {
                if err := validateType(param.Default, param.Type); err != nil {
                    return fmt.Errorf("template %s, parameter %s: %w", name, paramName, err)
                }
            }
            
            // 验证 flag 的 Go Template 语法
            if param.Flag != "" {
                if _, err := template.New(name + "_" + paramName).Parse(param.Flag); err != nil {
                    return fmt.Errorf("template %s, parameter %s: invalid flag syntax: %w",
                        name, paramName, err)
                }
            }
        }
    }
    
    // 验证每个阶段至少有一个工具
    stageTools := make(map[string]int)
    for _, tmpl := range l.cache {
        stageTools[tmpl.Metadata.Stage]++
    }
    for _, stage := range l.metadata.Stages {
        if stageTools[stage.ID] == 0 {
            return fmt.Errorf("stage %s has no tools defined", stage.ID)
        }
    }
    
    return nil
}

// ValidateStageDependencies 验证阶段依赖关系
func (l *TemplateLoader) ValidateStageDependencies(config map[string]any) error {
    // 获取元数据
    metadata := l.GetMetadata()
    
    // 构建启用的阶段集合
    enabledStages := make(map[string]bool)
    for _, stage := range metadata.Stages {
        // 检查配置中是否启用了该阶段
        if isStageEnabled(config, stage.ID) {
            enabledStages[stage.ID] = true
        }
    }
    
    // 验证依赖关系
    for _, stage := range metadata.Stages {
        if !enabledStages[stage.ID] {
            continue
        }
        
        // 检查依赖的阶段是否都已启用
        for _, depID := range stage.DependsOn {
            if !enabledStages[depID] {
                return fmt.Errorf("stage %s depends on %s, but %s is not enabled",
                    stage.ID, depID, depID)
            }
        }
    }
    
    return nil
}
```

**设计说明**:
- 使用 `TemplateSource` 接口，支持多种数据源
- 当前只需要使用 `EmbedSource`
- 未来扩展只需添加新的 Source 实现，不需要修改 TemplateLoader
- 使用 `sync.Once` 保证只加载一次
- 提供完整的验证逻辑

**使用示例**:
```go
// worker/internal/workflow/subdomain_discovery/stages.go
package subdomain_discovery

import (
    "embed"
    "github.com/orbit/worker/internal/activity"
    "github.com/orbit/worker/internal/workflow"
)

//go:embed templates.yaml
var templatesFS embed.FS

// 使用 EmbedSource 创建加载器
var templateLoader = activity.NewTemplateLoader(&activity.EmbedSource{
    fs:       templatesFS,
    filename: "templates.yaml",
})

// Metadata 定义
var Metadata = workflow.WorkflowMetadata{
    Name: "subdomain_discovery",
    // ...
}

// 在需要的地方使用
func runPassiveStage(ctx context.Context, config map[string]any) error {
    tmpl, err := templateLoader.Get("subfinder")
    if err != nil {
        return err
    }
    // 使用模板构建命令...
}
```


### 4. CommandBuilder

负责使用 Go Template 构建最终命令：

```go
// CommandBuilder 根据模板和配置构建命令（使用 Go Template）
type CommandBuilder struct {
    funcMap template.FuncMap
}

// Build 构建命令字符串
// toolName: 工具名称（如 "subfinder"）
// params: 必需参数（如 Domain, OutputFile）
// config: 用户配置（可选参数）
func (b *CommandBuilder) Build(
    tmpl CommandTemplate,
    params map[string]any,
    config map[string]any,
) (string, error)
```

**构建流程**:
1. 合并必需参数和可选参数到一个 data map
2. 应用覆盖优先级：用户配置 > 默认值
3. 构建完整的命令模板字符串（base_command + 启用的参数 flags）
4. 使用 Go Template 执行模板
5. 返回最终命令

**Go Template 配置**:
- 使用 `missingkey=error` 自动检测缺失字段
- 提供 funcMap：
  - `quote`: 引号转义 `{{quote .Domain}}`
  - `default`: 提供默认值 `{{default 3600 .Timeout}}`
  - `lower`/`upper`: 大小写转换
  - `join`: 数组连接

**示例**:
```go
// 输入
tmpl := CommandTemplate{
    BaseCommand: "subfinder -d {{.Domain}} -o {{quote .OutputFile}}",
    Parameters: map[string]Parameter{
        "Timeout": {
            Flag: "-timeout {{.Timeout}}",
            Default: 3600,
            Type: "int",
        },
    },
}
params := map[string]any{"Domain": "example.com", "OutputFile": "/tmp/out.txt"}
config := map[string]any{"Timeout": 7200}

// 输出
// "subfinder -d example.com -o \"/tmp/out.txt\" -timeout 7200"
```

### 4. SchemaGenerator（新增）

从 Worker 模板生成 JSON Schema：

```go
// SchemaGenerator 从模板生成 JSON Schema
type SchemaGenerator struct{}

// Generate 生成 JSON Schema
// metadata: Workflow 元数据
// templates: 所有工具的模板
// 返回: JSON Schema 字符串
func (g *SchemaGenerator) Generate(metadata WorkflowMetadata, templates map[string]CommandTemplate) (string, error)
```

**生成规则**:
- 每个工具对应一个 Schema 对象
- 参数类型映射：string → string, int → integer, bool → boolean
- Required 参数添加到 required 数组
- Description 映射到 description 字段
- Default 映射到 default 字段
- **使用元数据生成更丰富的 Schema**（工具说明、阶段信息）

**输出示例**:
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Subdomain Discovery Configuration",
  "description": "通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名",
  "type": "object",
  "properties": {
    "subfinder": {
      "type": "object",
      "description": "使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名",
      "properties": {
        "enabled": {"type": "boolean", "default": false},
        "timeout": {
          "type": "integer", 
          "default": 3600, 
          "description": "扫描超时时间（秒）"
        },
        "all": {"type": "boolean", "default": true, "description": "使用所有数据源"}
      }
    }
  }
}
```

### 5. DocGenerator（新增）

从 Worker 模板生成配置文档：

```go
// DocGenerator 从模板生成配置文档
type DocGenerator struct{}

// Generate 生成 Markdown 文档
// metadata: Workflow 元数据
// templates: 所有工具的模板
// 返回: Markdown 字符串
func (g *DocGenerator) Generate(metadata WorkflowMetadata, templates map[string]CommandTemplate) (string, error)
```

**生成规则**:
- 使用 Workflow 元数据生成概述章节
- 使用 Stage 元数据生成扫描流程章节
- 每个工具一个章节，使用工具元数据
- 参数表格：名称、类型、默认值、必需、描述
- 标记废弃的参数
- 包含使用示例
- **显示工具所属阶段和依赖关系**
- **显示警告信息**（如主动扫描警告）

**输出示例**:
```markdown
# 子域名发现配置参考

## 概述

**名称**: 子域名发现 (subdomain_discovery)  
**版本**: 1.0.0  
**描述**: 通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名

## 扫描流程

### 阶段 1: 被动收集 (passive) [必需]

**描述**: 使用多个数据源被动收集子域名，不产生主动扫描流量  
**执行方式**: 并行执行  
**依赖**: 无

## 工具配置

### subfinder

**描述**: 使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名  
**阶段**: passive  
**主页**: https://github.com/projectdiscovery/subfinder  
**需要 API Keys**: 是

### 参数

| 参数 | 类型 | 默认值 | 必需 | 描述 |
|------|------|--------|------|------|
| timeout | int | 3600 | 否 | 扫描超时时间（秒） |
| all | bool | true | 否 | 使用所有数据源 |

### 示例

```yaml
subfinder:
  enabled: true
  timeout: 7200
```
```


### 6. ConfigValidator（新增）

Server 端配置验证器：

```go
// ConfigValidator 验证用户配置
type ConfigValidator struct {
    schema   *jsonschema.Schema
    metadata WorkflowMetadata
}

// LoadSchema 加载 JSON Schema 和元数据
func (v *ConfigValidator) LoadSchema(schemaPath string, metadataPath string) error

// Validate 验证用户配置
// config: 用户配置（YAML 解析后的 map）
// 返回: 验证错误列表
func (v *ConfigValidator) Validate(config map[string]interface{}) []error

// ValidateStageDependencies 验证阶段依赖关系
func (v *ConfigValidator) ValidateStageDependencies(config map[string]interface{}) []error
```

**验证规则**:
- 使用 JSON Schema 验证库（如 `github.com/xeipuuv/gojsonschema`）
- 验证类型匹配
- 验证必需字段
- 验证值范围（如果 Schema 中定义）
- **验证阶段依赖关系**（依赖的阶段必须启用）
- 返回详细的验证错误

### 7. 参数覆盖逻辑

```go
// mergeParameters 合并必需参数、默认值和用户配置
func mergeParameters(
    tmpl CommandTemplate,
    params map[string]any,      // 必需参数（Domain, OutputFile 等）
    config map[string]any,      // 用户配置
) (map[string]any, error) {
    result := make(map[string]any)
    
    // 1. 添加必需参数
    for key, value := range params {
        result[key] = value
    }
    
    // 2. 添加可选参数（默认值 + 用户配置）
    for name, param := range tmpl.Parameters {
        if userValue, exists := config[name]; exists {
            // 用户配置优先
            if err := validateType(userValue, param.Type); err != nil {
                return nil, fmt.Errorf("parameter %s: %w", name, err)
            }
            result[name] = userValue
        } else if param.Default != nil {
            // 使用默认值
            result[name] = param.Default
        } else if param.Required {
            // 必需参数缺失
            return nil, fmt.Errorf("required parameter %s is missing", name)
        }
        // 可选参数且无默认值：不添加到结果
    }
    
    return result, nil
}

// buildCommandTemplate 构建完整的命令模板字符串
func buildCommandTemplate(tmpl CommandTemplate, data map[string]any) string {
    cmd := tmpl.BaseCommand
    
    // 只添加有值的参数的 flag
    for name, param := range tmpl.Parameters {
        if _, exists := data[name]; exists {
            cmd += " " + param.Flag
        }
    }
    
    return cmd
}
```

## Data Models

### Worker 模板格式（templates.yaml）

使用 Go Template 语法、YAML 锚点和元数据：

```yaml
# 文件头注释：说明整体结构和使用方式
# 本文件定义所有扫描工具的命令模板
# 使用 Go Template 语法（{{.Var}}）和 YAML 锚点共享通用参数定义

# Workflow 元数据
metadata:
  name: "subdomain_discovery"
  display_name: "子域名发现"
  description: "通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名"
  version: "1.0.0"
  target_types: ["domain"]
  
  # 阶段定义
  stages:
    - id: "passive"
      name: "被动收集"
      description: "使用多个数据源被动收集子域名，不产生主动扫描流量"
      order: 1
      required: true
      parallel: true
      depends_on: []
      outputs: ["subdomains"]
    
    - id: "bruteforce"
      name: "字典爆破"
      description: "使用字典对域名进行爆破，发现未公开的子域名"
      order: 2
      required: false
      parallel: false
      depends_on: []
      outputs: ["subdomains"]
    
    - id: "permutation"
      name: "排列组合"
      description: "对已发现的子域名进行排列组合，生成新的可能子域名"
      order: 3
      required: false
      parallel: false
      depends_on: ["passive", "bruteforce"]
      outputs: ["subdomains"]

# 共享参数定义（使用 YAML 锚点）
x-common-params: &common-params
  Timeout:
    flag: "-timeout {{.Timeout}}"
    default: 3600
    type: "int"
    required: false
    description: "扫描超时时间（秒）"
  
  RateLimit:
    flag: "-rl {{.RateLimit}}"
    default: 150
    type: "int"
    required: false
    description: "每秒请求数限制"

# 工具模板
tools:
  subfinder:
    metadata:
      display_name: "Subfinder"
      description: "使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名"
      stage: "passive"
      warning: "建议配置 API Keys（Shodan、Censys、VirusTotal、SecurityTrails）以提高收集效率"
    
    base_command: "subfinder -d {{.Domain}} -all -o {{quote .OutputFile}} -v"
    parameters:
      <<: *common-params  # 引用共享参数
      Threads:
        flag: "-t {{.Threads}}"
        default: 10
        type: "int"
        required: false
        description: "并发线程数"
      
      ProviderConfig:
        flag: "-pc {{quote .ProviderConfig}}"
        default: null
        type: "string"
        required: false
        description: "Provider 配置文件路径"

  httpx:
    metadata:
      display_name: "HTTPX"
      description: "HTTP 探测工具"
      stage: "http_probe"
    
    base_command: "httpx -l {{quote .InputFile}} -o {{quote .OutputFile}}"
    parameters:
      <<: *common-params  # 复用相同的 Timeout 和 RateLimit
      Threads:
        flag: "-threads {{.Threads}}"
        default: 50
        type: "int"
        required: false
        description: "并发线程数"
      
      # 本地配置覆盖锚点定义
      Timeout:
        flag: "-timeout {{.Timeout}}"
        default: 1800  # httpx 使用更短的超时时间
        type: "int"
        required: false
        description: "扫描超时时间（秒）"

  # 条件渲染示例
  nuclei:
    metadata:
      display_name: "Nuclei"
      description: "漏洞扫描工具"
      stage: "vulnerability_scan"
      warning: "漏洞扫描可能触发目标的安全防护系统"
    
    base_command: "nuclei -u {{.Target}} {{if .Verbose}}-v{{end}} {{if .Debug}}-debug{{end}}"
    parameters:
      Verbose:
        flag: ""  # 条件在 base_command 中处理
        default: false
        type: "bool"
        required: false
        description: "启用详细输出"
      
      Debug:
        flag: ""
        default: false
        type: "bool"
        required: false
        description: "启用调试模式"
      
      Templates:
        flag: "-t {{.Templates}}"
        default: null
        type: "string"
        required: false
        description: "模板路径"

  subdomain-bruteforce:
    metadata:
      display_name: "Subdomain Bruteforce"
      description: "使用字典对域名进行 DNS 爆破"
      stage: "bruteforce"
      warning: "主动扫描会产生大量 DNS 请求，可能被目标检测"
    
    base_command: "puredns bruteforce {{quote .Wordlist}} {{.Domain}} -r {{quote .Resolvers}} --write {{quote .OutputFile}} --quiet"
    parameters:
      <<: *common-params
      Threads:
        flag: "-t {{.Threads}}"
        default: 100
        type: "int"
        required: false
        description: "并发线程数"
      
      WildcardTests:
        flag: "--wildcard-tests {{.WildcardTests}}"
        default: 50
        type: "int"
        required: false
        description: "泛解析检测测试次数"
```

**命名规范**:
- 占位符使用 **PascalCase**（大驼峰）：`{{.Domain}}`, `{{.OutputFile}}`
- 参数名使用 **PascalCase**：`Timeout`, `RateLimit`, `ProviderConfig`
- 必需字段：`Domain`, `OutputFile`, `InputFile`, `Target`
- 可选字段：`Timeout`, `Threads`, `RateLimit`, `Verbose`

**Go Template 特性**:
- 条件渲染：`{{if .Verbose}}-v{{end}}`
- 函数调用：`{{quote .OutputFile}}`（自动添加引号）
- 默认值：`{{default 3600 .Timeout}}`
- 循环：`{{range .Domains}}-d {{.}}{{end}}`

**元数据字段**:
- `metadata`: Workflow 整体元数据（名称、版本、阶段定义）
- `tools.*.metadata`: 工具元数据（说明、阶段、主页、警告）
- `parameters.*.min/max`: 参数范围约束（int 类型）

### Server 配置格式（用户编写，保持简单）

用户配置保持扁平结构，使用 snake_case（与 Go 的 PascalCase 自动映射）：

```yaml
subdomain_discovery:
  passive_tools:
    subfinder:
      enabled: true
      timeout: 7200      # 覆盖默认值（自动映射到 Timeout）
      threads: 20        # 自动映射到 Threads
      # 其他参数使用默认值
    
    amass:
      enabled: true      # 所有参数使用默认值
```

**命名映射**:
- 用户配置（snake_case）→ Go Template（PascalCase）
- `timeout` → `Timeout`
- `rate_limit` → `RateLimit`
- `provider_config` → `ProviderConfig`

**映射实现**:
```go
// 自动转换 snake_case 到 PascalCase
func convertKeys(config map[string]any) map[string]any {
    result := make(map[string]any)
    for key, value := range config {
        pascalKey := snakeToPascal(key)  // timeout -> Timeout
        result[pascalKey] = value
    }
    return result
}
```

### 生成的 JSON Schema（config.schema.json）

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Subdomain Discovery Configuration",
  "description": "通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名",
  "version": "1.0.0",
  "type": "object",
  "properties": {
    "subfinder": {
      "type": "object",
      "description": "使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名",
      "x-stage": "passive",
      "x-warning": "建议配置 API Keys（Shodan、Censys、VirusTotal、SecurityTrails）以提高收集效率",
      "properties": {
        "enabled": {
          "type": "boolean",
          "default": false,
          "description": "是否启用该工具"
        },
        "timeout": {
          "type": "integer",
          "default": 3600,
          "minimum": 1,
          "maximum": 86400,
          "description": "扫描超时时间（秒）"
        },
        "rate_limit": {
          "type": "integer",
          "default": 150,
          "minimum": 1,
          "maximum": 10000,
          "description": "每秒请求数限制"
        },
        "threads": {
          "type": "integer",
          "default": 10,
          "minimum": 1,
          "maximum": 100,
          "description": "并发线程数"
        },
        "provider_config": {
          "type": "string",
          "description": "Provider 配置文件路径"
        }
      }
    },
    "subdomain-bruteforce": {
      "type": "object",
      "description": "使用字典对域名进行 DNS 爆破",
      "x-stage": "bruteforce",
      "x-warning": "主动扫描会产生大量 DNS 请求，可能被目标检测",
      "properties": {
        "enabled": {
          "type": "boolean",
          "default": false
        },
        "timeout": {
          "type": "integer",
          "default": 3600,
          "minimum": 1,
          "maximum": 86400,
          "description": "扫描超时时间（秒）"
        },
        "threads": {
          "type": "integer",
          "default": 100,
          "description": "并发线程数"
        },
        "wildcard_tests": {
          "type": "integer",
          "default": 50,
          "description": "泛解析检测测试次数"
        }
      }
    }
  },
  "x-metadata": {
    "stages": [
      {
        "id": "passive",
        "name": "被动收集",
        "description": "使用多个数据源被动收集子域名，不产生主动扫描流量",
        "order": 1,
        "required": true,
        "parallel": true
      },
      {
        "id": "bruteforce",
        "name": "字典爆破",
        "description": "使用字典对域名进行爆破，发现未公开的子域名",
        "order": 2,
        "required": false,
        "parallel": false
      }
    ]
  }
}
```


### 生成的配置文档（config-reference.md）

```markdown
# 子域名发现配置参考

本文档由 Worker 模板自动生成，描述所有可用的配置参数。

## 概述

**名称**: 子域名发现 (subdomain_discovery)  
**版本**: 1.0.0  
**描述**: 通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名  
**支持的目标类型**: domain

## 扫描流程

子域名发现包含 3 个阶段，按顺序执行：

### 阶段 1: 被动收集 (passive) [必需]

**描述**: 使用多个数据源被动收集子域名，不产生主动扫描流量  
**执行方式**: 并行执行  
**依赖**: 无  
**输出**: 子域名列表

### 阶段 2: 字典爆破 (bruteforce) [可选]

**描述**: 使用字典对域名进行爆破，发现未公开的子域名  
**执行方式**: 顺序执行  
**依赖**: 无  
**输出**: 子域名列表

### 阶段 3: 排列组合 (permutation) [可选]

**描述**: 对已发现的子域名进行排列组合，生成新的可能子域名  
**执行方式**: 顺序执行  
**依赖**: passive, bruteforce  
**输出**: 子域名列表

## 工具配置

### subfinder

**描述**: 使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名  
**阶段**: passive  
**警告**: ⚠️ 建议配置 API Keys（Shodan、Censys、VirusTotal、SecurityTrails）以提高收集效率

#### 参数

| 参数 | 类型 | 默认值 | 范围 | 必需 | 描述 |
|------|------|--------|------|------|------|
| timeout | int | 3600 | 1-86400 | 否 | 扫描超时时间（秒） |
| rate_limit | int | 150 | 否 | 每秒请求数限制 |
| threads | int | 10 | 否 | 并发线程数 |
| provider_config | string | - | 否 | Provider 配置文件路径 |

#### 示例

```yaml
subfinder:
  enabled: true
  timeout: 7200
  threads: 20
```

### subdomain-bruteforce

**描述**: 使用字典对域名进行 DNS 爆破  
**阶段**: bruteforce  
**警告**: ⚠️ 主动扫描会产生大量 DNS 请求，可能被目标检测

#### 参数

| 参数 | 类型 | 默认值 | 必需 | 描述 |
|------|------|--------|------|------|
| timeout | int | 3600 | 否 | 扫描超时时间（秒） |
| threads | int | 100 | 否 | 并发线程数 |
| wildcard_tests | int | 50 | 否 | 泛解析检测测试次数 |

#### 示例

```yaml
subdomain-bruteforce:
  enabled: true
  timeout: 86400
  threads: 200
```

## 完整配置示例

```yaml
# 阶段 1: 被动收集（必需）
passive-tools:
  subfinder:
    enabled: true
    timeout: 7200
    threads: 20

# 阶段 2: 字典爆破（可选）
bruteforce:
  enabled: false
  subdomain-bruteforce:
    timeout: 86400
    threads: 200

# 阶段 3: 排列组合（可选）
permutation:
  enabled: true
```
```

## Workflow Code Patterns

### 动态工具发现（模板驱动）

为了支持动态模板加载（文件、URL、数据库），工作流代码必须从模板中动态发现可用工具，而不是硬编码工具名称。

#### ❌ 错误模式：硬编码工具名称

```go
// worker/internal/workflow/subdomain_discovery/stages.go

const (
    toolSubfinder   = "subfinder"
    toolSublist3r   = "sublist3r"
    toolAssetfinder = "assetfinder"
)

var passiveTools = []string{toolSubfinder, toolSublist3r, toolAssetfinder}

func (w *Workflow) runPassiveStage(ctx *workflowContext) stageResult {
    // 硬编码工具列表，模板变化时代码需要修改
    for _, tool := range passiveTools {
        // ...
    }
}
```

**问题**：
- 模板中删除工具时，代码会尝试加载不存在的模板
- 模板中新增工具时，代码不会自动使用
- 工作流代码和模板强耦合，无法支持动态模板

#### ✅ 正确模式：从模板动态发现

```go
// worker/internal/workflow/subdomain_discovery/stages.go

func (w *Workflow) runPassiveStage(ctx *workflowContext) stageResult {
    // 1. 从模板中获取 passive 阶段的所有工具
    tools := w.getToolsByStage("passive")
    
    // 2. 根据用户配置过滤启用的工具
    enabledTools := w.filterEnabledTools(tools, ctx.config)
    
    // 3. 并行执行所有启用的工具
    results := w.executeToolsInParallel(ctx, enabledTools)
    
    return results
}

// getToolsByStage 从模板中获取指定阶段的所有工具
func (w *Workflow) getToolsByStage(stage string) []string {
    templates, err := templateLoader.Load()
    if err != nil {
        pkg.Logger.Error("Failed to load templates", zap.Error(err))
        return nil
    }
    
    var tools []string
    for name, tmpl := range templates {
        if tmpl.Metadata.Stage == stage {
            tools = append(tools, name)
        }
    }
    
    pkg.Logger.Debug("Discovered tools for stage",
        zap.String("stage", stage),
        zap.Strings("tools", tools))
    
    return tools
}

// filterEnabledTools 根据用户配置过滤启用的工具
func (w *Workflow) filterEnabledTools(tools []string, config map[string]any) []string {
    var enabled []string
    for _, tool := range tools {
        if isToolEnabled(config, tool) {
            enabled = append(enabled, tool)
        }
    }
    
    pkg.Logger.Info("Enabled tools",
        zap.Strings("tools", enabled))
    
    return enabled
}

// isToolEnabled 检查工具是否在配置中启用
func isToolEnabled(config map[string]any, toolName string) bool {
    toolConfig, ok := config[toolName].(map[string]any)
    if !ok {
        return false
    }
    
    enabled, ok := toolConfig["enabled"].(bool)
    return ok && enabled
}
```

**优势**：
- ✅ 完全动态：工作流代码不需要知道具体有哪些工具
- ✅ 模板驱动：新增/删除工具只需修改模板，代码不变
- ✅ 向后兼容：旧模板和新模板都能工作
- ✅ 易于扩展：付费功能可以提供不同的模板

### 工作流代码规范

#### 1. 阶段名称可以硬编码

阶段名称相对稳定，可以定义为常量：

```go
const (
    stagePassive     = "passive"
    stageBruteforce  = "bruteforce"
    stagePermutation = "permutation"
)
```

#### 2. 工具名称必须动态获取

工具名称不能硬编码，必须从模板中读取：

```go
// ❌ 错误
const toolSubfinder = "subfinder"

// ✅ 正确
tools := w.getToolsByStage("passive")
```

#### 3. 使用 Stage 字段组织工具

模板中的 `ToolMetadata.Stage` 字段是关键：

```yaml
tools:
  subfinder:
    metadata:
      stage: "passive"  # 工作流代码通过这个字段发现工具
  
  new-tool:
    metadata:
      stage: "passive"  # 新增工具，工作流代码自动发现
```

#### 4. 配置验证

工作流代码应该验证配置的合法性：

```go
func (w *Workflow) validateConfig(config map[string]any) error {
    // 获取所有可用工具
    allTools := w.getAllTools()
    
    // 检查配置中的工具是否都存在于模板中
    for toolName := range config {
        if !contains(allTools, toolName) {
            return fmt.Errorf("unknown tool in config: %s", toolName)
        }
    }
    
    return nil
}
```

### 完整示例

```go
// worker/internal/workflow/subdomain_discovery/stages.go
package subdomain_discovery

import (
    "context"
    "github.com/orbit/worker/internal/activity"
    "github.com/orbit/worker/internal/pkg"
    "go.uber.org/zap"
)

// 阶段名称（相对稳定，可以硬编码）
const (
    stagePassive     = "passive"
    stageBruteforce  = "bruteforce"
    stagePermutation = "permutation"
)

// runAllStages 执行所有阶段（完全动态）
func (w *Workflow) runAllStages(ctx *workflowContext) stageResult {
    var allResults stageResult
    
    // 阶段 1: 被动收集（必需）
    allResults.merge(w.runStage(ctx, stagePassive))
    
    // 阶段 2: 字典爆破（可选）
    if isStageEnabled(ctx.config, stageBruteforce) {
        allResults.merge(w.runStage(ctx, stageBruteforce))
    }
    
    // 阶段 3: 排列组合（可选）
    if isStageEnabled(ctx.config, stagePermutation) && len(allResults.files) > 0 {
        allResults.merge(w.runStage(ctx, stagePermutation))
    }
    
    return allResults
}

// runStage 执行单个阶段（通用方法）
func (w *Workflow) runStage(ctx *workflowContext, stage string) stageResult {
    pkg.Logger.Info("Running stage", zap.String("stage", stage))
    
    // 1. 从模板动态获取工具
    tools := w.getToolsByStage(stage)
    if len(tools) == 0 {
        pkg.Logger.Warn("No tools found for stage", zap.String("stage", stage))
        return stageResult{}
    }
    
    // 2. 过滤启用的工具
    enabledTools := w.filterEnabledTools(tools, ctx.config)
    if len(enabledTools) == 0 {
        pkg.Logger.Info("No enabled tools for stage", zap.String("stage", stage))
        return stageResult{}
    }
    
    // 3. 执行工具
    results := w.executeTools(ctx, enabledTools)
    
    return results
}

// getToolsByStage 从模板中获取指定阶段的所有工具
func (w *Workflow) getToolsByStage(stage string) []string {
    templates, err := templateLoader.Load()
    if err != nil {
        pkg.Logger.Error("Failed to load templates", zap.Error(err))
        return nil
    }
    
    var tools []string
    for name, tmpl := range templates {
        if tmpl.Metadata.Stage == stage {
            tools = append(tools, name)
        }
    }
    
    return tools
}

// filterEnabledTools 根据用户配置过滤启用的工具
func (w *Workflow) filterEnabledTools(tools []string, config map[string]any) []string {
    var enabled []string
    for _, tool := range tools {
        if isToolEnabled(config, tool) {
            enabled = append(enabled, tool)
        }
    }
    return enabled
}

// executeTools 执行工具列表
func (w *Workflow) executeTools(ctx *workflowContext, tools []string) stageResult {
    var results stageResult
    
    for _, tool := range tools {
        // 获取模板
        tmpl, err := templateLoader.Get(tool)
        if err != nil {
            pkg.Logger.Error("Failed to get template",
                zap.String("tool", tool),
                zap.Error(err))
            results.failed = append(results.failed, tool)
            continue
        }
        
        // 构建命令
        cmd, err := buildCommand(tmpl, ctx)
        if err != nil {
            pkg.Logger.Error("Failed to build command",
                zap.String("tool", tool),
                zap.Error(err))
            results.failed = append(results.failed, tool)
            continue
        }
        
        // 执行命令
        if err := executeCommand(ctx.ctx, cmd); err != nil {
            pkg.Logger.Error("Failed to execute command",
                zap.String("tool", tool),
                zap.Error(err))
            results.failed = append(results.failed, tool)
            continue
        }
        
        results.success = append(results.success, tool)
    }
    
    return results
}
```

### 测试策略

#### 单元测试

测试动态工具发现：

```go
func TestGetToolsByStage(t *testing.T) {
    // 创建测试模板
    templates := map[string]CommandTemplate{
        "subfinder": {
            Metadata: ToolMetadata{Stage: "passive"},
        },
        "amass": {
            Metadata: ToolMetadata{Stage: "passive"},
        },
        "puredns": {
            Metadata: ToolMetadata{Stage: "bruteforce"},
        },
    }
    
    // 测试获取 passive 阶段的工具
    tools := getToolsByStage(templates, "passive")
    assert.ElementsMatch(t, []string{"subfinder", "amass"}, tools)
    
    // 测试获取 bruteforce 阶段的工具
    tools = getToolsByStage(templates, "bruteforce")
    assert.ElementsMatch(t, []string{"puredns"}, tools)
    
    // 测试不存在的阶段
    tools = getToolsByStage(templates, "nonexistent")
    assert.Empty(t, tools)
}
```

#### 集成测试

测试模板变化时的兼容性：

```go
func TestWorkflowWithDifferentTemplates(t *testing.T) {
    // 测试 1: 使用完整模板
    fullTemplate := loadTemplate("templates_full.yaml")
    result := runWorkflow(fullTemplate, config)
    assert.NoError(t, result.Error)
    
    // 测试 2: 使用最小模板（只有必需工具）
    minimalTemplate := loadTemplate("templates_minimal.yaml")
    result = runWorkflow(minimalTemplate, config)
    assert.NoError(t, result.Error)
    
    // 测试 3: 使用付费模板（包含额外工具）
    premiumTemplate := loadTemplate("templates_premium.yaml")
    result = runWorkflow(premiumTemplate, config)
    assert.NoError(t, result.Error)
}
```

## Correctness Properties

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 参数默认值支持

*对于任何* 参数定义，系统应该能够为其指定默认值，并在用户未提供配置时使用该默认值

**Validates: Requirements 1.1, 1.2**

### Property 2: 用户配置覆盖默认值

*对于任何* 参数，当用户配置和模板默认值都存在时，系统应该使用用户配置值覆盖默认值

**Validates: Requirements 1.3**

### Property 3: YAML 锚点正确解析

*对于任何* 包含 YAML 锚点和合并键的模板，解析器应该正确识别锚点并合并配置

**Validates: Requirements 2.1, 2.2**

### Property 4: 参数类型验证

*对于任何* 参数值，系统应该验证其类型与定义的类型匹配，类型不匹配时返回错误

**Validates: Requirements 3.2, 3.3**

### Property 5: 类型自动转换

*对于任何* int 类型的参数值，系统应该能够自动转换为 string 以适配命令行参数

**Validates: Requirements 3.5**

### Property 6: 模板验证综合

*对于任何* 包含语法错误、未定义参数引用或类型不匹配的模板，系统应该在加载时检测并报错

**Validates: Requirements 4.2, 4.3, 4.4**

### Property 7: 错误信息完整性

*对于任何* 命令构建失败的情况，错误信息应该包含工具名、参数名、类型信息或占位符列表等上下文

**Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5**

### Property 8: 模板缓存复用

*对于任何* 模板，多次构建命令时应该复用缓存的模板，不重复加载文件

**Validates: Requirements 8.2**

### Property 9: 可选参数处理

*对于任何* 既没有默认值也没有用户配置的可选参数，系统应该不将其添加到最终命令中

**Validates: Requirements 9.4**

### Property 10: 嵌套参数解析

*对于任何* 包含嵌套参数定义的 YAML 模板，系统应该正确加载所有参数属性

**Validates: Requirements 10.3**

### Property 11: 默认值类型支持

*对于任何* string、int、bool 类型的默认值，系统应该能够正确处理

**Validates: Requirements 10.4**

### Property 12: 必需参数验证

*对于任何* 标记为 required 但未提供值的参数，以及缺少必需占位符的配置，系统应该返回明确的错误

**Validates: Requirements 11.5, 11.6**

### Property 13: YAML 锚点冲突处理

*对于任何* 锚点和本地配置冲突的情况，解析器应该优先使用本地配置

**Validates: Requirements 12.4**

### Property 14: 无效锚点引用检测

*对于任何* 引用不存在锚点的 YAML，解析器应该返回明确的错误信息

**Validates: Requirements 12.5**

### Property 15: Schema 生成正确性

*对于任何* Worker 模板，生成的 JSON Schema 应该正确映射所有参数的类型、默认值、描述和范围约束

**Validates: 新增功能 - Schema 生成**

### Property 16: 文档生成完整性

*对于任何* Worker 模板，生成的配置文档应该包含所有参数的完整信息（名称、类型、默认值、范围、描述）以及工具元数据（阶段、主页、警告）

**Validates: 新增功能 - 文档生成**

### Property 17: Server 配置验证

*对于任何* 用户配置，Server 端应该使用 JSON Schema 验证其正确性，并在验证失败时返回详细错误

**Validates: 新增功能 - Server 端验证**

### Property 18: 阶段依赖验证

*对于任何* 启用的阶段，如果它依赖其他阶段，系统应该验证依赖的阶段也已启用

**Validates: 新增功能 - 元数据验证**

### Property 19: 元数据完整性

*对于任何* 工具模板，元数据应该包含必需字段（display_name, description, stage），并且 stage 字段应该引用已定义的阶段

**Validates: 新增功能 - 元数据验证**


## Error Handling

### 错误类型

1. **模板加载错误**
   - YAML 语法错误
   - 文件不存在
   - 权限不足
   - YAML 锚点引用无效
   - Go Template 语法错误

2. **模板验证错误**
   - 参数类型无效（不是 string/int/bool）
   - 默认值类型与参数类型不匹配
   - 必需参数有默认值（逻辑冲突）
   - 引用未定义的锚点
   - 参数定义缺少必需字段
   - Go Template 语法错误（无法解析）

3. **命令构建错误**
   - 工具模板不存在
   - 必需参数缺失（Go Template 自动检测）
   - 参数类型不匹配
   - 模板执行失败
   - 类型转换失败

4. **Schema 生成错误**
   - 模板格式不符合预期
   - 类型映射失败
   - JSON 序列化失败

5. **配置验证错误**（Server 端）
   - 配置不符合 Schema
   - 类型不匹配
   - 必需字段缺失
   - 值超出范围

### 错误信息格式

```go
// 模板验证错误
fmt.Errorf("template validation failed for tool %s: parameter %s has invalid type %s (expected string/int/bool)", 
    toolName, paramName, paramType)

// Go Template 语法错误
fmt.Errorf("template parse failed for tool %s: %w", toolName, err)
// 示例输出：template parse failed for tool subfinder: template: command:1: unexpected "}" in operand

// 必需参数缺失（Go Template 自动检测）
fmt.Errorf("template execution failed for tool %s: %w", toolName, err)
// 示例输出：template execution failed for tool subfinder: template: command:1:15: executing "command" at <.Domain>: map has no entry for key "Domain"

// 类型不匹配
fmt.Errorf("command build failed for tool %s: parameter %s expects type %s but got %s", 
    toolName, paramName, expectedType, actualType)

// YAML 锚点错误
fmt.Errorf("template parsing failed: anchor %s is referenced but not defined", 
    anchorName)

// Schema 验证错误（Server 端）
fmt.Errorf("config validation failed for tool %s: %s", 
    toolName, validationErrors)
```

### Go Template 错误优势

Go Template 提供更详细的错误信息：

```go
// 缺失字段错误（自动检测）
template: command:1:15: executing "command" at <.Domain>: map has no entry for key "Domain"
//                      ↑ 行号:列号              ↑ 字段名                    ↑ 具体错误

// 语法错误
template: command:1: unexpected "}" in operand
//                   ↑ 具体的语法问题

// 函数调用错误
template: command:1:20: executing "command" at <quote .Domain>: error calling quote: invalid argument type
```

### 错误处理策略

1. **启动时快速失败**
   - Worker 启动时验证所有模板，发现错误立即退出
   - Server 启动时加载 Schema，验证失败立即退出
   - 不允许带着无效配置运行

2. **详细的错误上下文**
   - 包含文件名、工具名、参数名
   - 显示期望值和实际值
   - 提供修复建议

3. **日志记录**
   - 错误级别：模板验证失败、命令构建失败
   - 警告级别：使用废弃参数
   - 信息级别：使用默认值、参数覆盖

## Testing Strategy

### 单元测试

测试特定示例、边界情况和错误条件：

1. **TemplateLoader 测试**
   - 加载有效模板
   - 加载包含 YAML 锚点的模板
   - 加载包含语法错误的模板（应失败）
   - 加载包含无效类型的模板（应失败）
   - 加载引用不存在锚点的模板（应失败）

2. **CommandBuilder 测试**
   - 构建只使用默认值的命令
   - 构建部分覆盖默认值的命令
   - 构建缺少必需参数的命令（应失败）
   - 构建包含类型错误的命令（应失败）
   - 构建包含未替换占位符的命令（应失败）

3. **参数合并测试**
   - 用户配置覆盖默认值
   - 可选参数无默认值时不添加
   - YAML 锚点冲突时优先本地配置
   - 类型转换（int → string）

4. **SchemaGenerator 测试**
   - 生成简单工具的 Schema
   - 生成包含所有类型参数的 Schema
   - 验证 Schema 格式正确（JSON 有效）
   - 验证类型映射正确

5. **DocGenerator 测试**
   - 生成简单工具的文档
   - 生成包含所有字段的文档
   - 验证 Markdown 格式正确
   - 验证表格包含所有必需列

6. **ConfigValidator 测试**
   - 验证有效配置
   - 验证类型不匹配的配置（应失败）
   - 验证缺少必需字段的配置（应失败）
   - 验证包含未知字段的配置

### 属性测试

验证通用属性在所有输入下成立（最小 100 次迭代）：

1. **Property 1-2: 默认值和覆盖**
   - 生成随机参数定义和用户配置
   - 验证默认值应用和覆盖逻辑
   - **Feature: worker-command-template-refactor, Property 1: 参数默认值支持**
   - **Feature: worker-command-template-refactor, Property 2: 用户配置覆盖默认值**

2. **Property 3: YAML 锚点**
   - 生成包含锚点的随机 YAML
   - 验证解析和合并正确性
   - **Feature: worker-command-template-refactor, Property 3: YAML 锚点正确解析**

3. **Property 4-5: 类型验证和转换**
   - 生成随机类型的参数值
   - 验证类型检查和转换逻辑
   - **Feature: worker-command-template-refactor, Property 4: 参数类型验证**
   - **Feature: worker-command-template-refactor, Property 5: 类型自动转换**

4. **Property 6: 模板验证**
   - 生成包含各种错误的模板
   - 验证所有错误都能检测
   - **Feature: worker-command-template-refactor, Property 6: 模板验证综合**

5. **Property 7: 错误信息**
   - 触发各种错误场景
   - 验证错误信息包含必要上下文
   - **Feature: worker-command-template-refactor, Property 7: 错误信息完整性**

6. **Property 8: 缓存**
   - 多次构建命令
   - 验证模板只加载一次
   - **Feature: worker-command-template-refactor, Property 8: 模板缓存复用**

7. **Property 9-14: 边界情况**
   - 可选参数处理
   - 嵌套参数解析
   - 必需参数验证
   - YAML 锚点边界情况
   - **Feature: worker-command-template-refactor, Property 9-14**

8. **Property 15-17: 新增功能**
   - Schema 生成正确性
   - 文档生成完整性
   - Server 配置验证
   - **Feature: worker-command-template-refactor, Property 15-17**

### 集成测试

验证完整流程：

1. **Worker 启动验证流程**
   - Worker 启动 → 加载模板 → 验证 → 缓存
   - 验证失败时拒绝启动
   - 验证成功时记录日志

2. **命令构建流程**
   - 接收任务 → 获取模板 → 合并配置 → 构建命令 → 执行
   - 验证最终命令字符串正确
   - 验证参数顺序和格式

3. **Schema 生成流程**
   - 运行 `go generate` → 生成 Schema → 验证 JSON 有效
   - 使用生成的 Schema 验证示例配置
   - 验证 Schema 与模板同步

4. **文档生成流程**
   - 运行 `go generate` → 生成文档 → 验证 Markdown 有效
   - 验证文档包含所有工具和参数
   - 验证示例代码可用

5. **Server 验证流程**
   - Server 启动 → 加载 Schema → 验证用户配置
   - 验证失败时拒绝启动
   - 验证成功时记录日志

6. **错误恢复流程**
   - 构建失败 → 记录错误 → 返回详细信息
   - 验证错误信息可用于调试
   - 验证系统继续运行（不崩溃）

### 测试配置

- 属性测试最小迭代次数：100
- 测试标签格式：`Feature: worker-command-template-refactor, Property {N}: {property_text}`
- 覆盖率目标：核心逻辑 > 90%
- 使用 Go 的 `testing/quick` 包进行属性测试
- 使用 `testify` 进行断言

## Implementation Notes

### go generate 配置

在 `template_loader.go` 文件中添加：

```go
//go:generate go run ./cmd/schema-gen/main.go -input templates.yaml -output config.schema.json
//go:generate go run ./cmd/doc-gen/main.go -input templates.yaml -output config-reference.md
```

运行 `go generate ./...` 自动生成 Schema 和文档。

### 一致性保证

**问题**：如何保证代码中的阶段依赖关系和 templates.yaml 中的元数据定义保持一致？

例如：
- 代码中 `permutation` 阶段依赖 `passive` 和 `bruteforce`
- templates.yaml 中也需要定义相同的依赖关系
- 如果代码改了但忘记更新 YAML，就会出现不一致

**业界最佳实践**（参考 Helm、Terraform）：

#### 方案：代码是唯一真实来源（Single Source of Truth）

**核心思想**：在代码中定义元数据，使用 `go generate` 自动生成 YAML 文件。

```go
// 1. 在 types.go 中定义元数据类型（共享类型定义）
// worker/internal/workflow/types.go
package workflow

// WorkflowMetadata 定义工作流元数据结构
type WorkflowMetadata struct {
    Name         string
    DisplayName  string
    Description  string
    Version      string
    TargetTypes  []string
    Stages       []StageMetadata
}

// StageMetadata 定义阶段元数据结构
type StageMetadata struct {
    ID          string
    Name        string
    Description string
    Order       int
    Required    bool
    Parallel    bool
    DependsOn   []string
    Outputs     []string
}

// 2. 在 stages.go 中定义元数据实例（唯一真实来源）
// worker/internal/workflow/subdomain_discovery/stages.go
package subdomain_discovery

import "github.com/orbit/worker/internal/workflow"

// Metadata 定义子域名发现工作流的元数据
// 这是唯一真实来源，templates.yaml 的 metadata 部分从这里生成
// 放在 stages.go 中是因为它与阶段编排逻辑（runAllStages）高度相关
var Metadata = workflow.WorkflowMetadata{
    Name:        "subdomain_discovery",
    DisplayName: "子域名发现",
    Description: "通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名",
    Version:     "1.0.0",
    TargetTypes: []string{"domain"},
    Stages: []workflow.StageMetadata{
        {
            ID:          "passive",
            Name:        "被动收集",
            Description: "使用多个数据源被动收集子域名，不产生主动扫描流量",
            Order:       1,
            Required:    true,
            Parallel:    true,
            DependsOn:   []string{},
            Outputs:     []string{"subdomains"},
        },
        {
            ID:          "bruteforce",
            Name:        "字典爆破",
            Description: "使用字典对域名进行爆破，发现未公开的子域名",
            Order:       2,
            Required:    false,
            Parallel:    false,
            DependsOn:   []string{},
            Outputs:     []string{"subdomains"},
        },
        {
            ID:          "permutation",
            Name:        "排列组合",
            Description: "对已发现的子域名进行排列组合，生成新的可能子域名",
            Order:       3,
            Required:    false,
            Parallel:    false,
            DependsOn:   []string{"passive", "bruteforce"},  // 在代码中定义依赖
            Outputs:     []string{"subdomains"},
        },
    },
}

// 2. 使用 go generate 从代码生成 YAML
//go:generate go run ../../cmd/gen-metadata/main.go -workflow subdomain_discovery
```

**生成工具实现**：
```go
// worker/cmd/gen-metadata/main.go
package main

import (
    "flag"
    "os"
    "gopkg.in/yaml.v3"
    "worker/internal/workflow/subdomain_discovery"
)

func main() {
    workflow := flag.String("workflow", "", "workflow name")
    flag.Parse()
    
    // 读取代码中的元数据
    metadata := subdomain_discovery.WorkflowMetadata
    
    // 读取现有的 templates.yaml
    data, _ := os.ReadFile("templates.yaml")
    var templates map[string]interface{}
    yaml.Unmarshal(data, &templates)
    
    // 更新 metadata 部分
    templates["metadata"] = metadata
    
    // 写回 templates.yaml
    output, _ := yaml.Marshal(templates)
    os.WriteFile("templates.yaml", output, 0644)
}
```

**工作流程**：
1. 开发者修改 `metadata.go` 中的阶段定义
2. 运行 `go generate ./...` 自动更新 `templates.yaml` 的 metadata 部分
3. CI 检查生成的文件是否最新（防止忘记运行 go generate）
4. 代码和 YAML 始终保持同步

**参考业界方案**：
- **Helm**: 使用 `helm-schema-gen` 从 values.yaml 生成 values.schema.json
  - 原理：从 YAML 推断类型，生成 JSON Schema
  - 链接：https://github.com/karuppiah7890/helm-schema-gen
  
- **Terraform**: 使用 `terraform-plugin-docs` 从 Provider Schema 生成文档
  - 原理：从 Go 代码中的 Schema 定义生成 Markdown 文档
  - 链接：https://github.com/hashicorp/terraform-plugin-docs

**优势**：
- ✅ 代码是唯一真实来源，避免不一致
- ✅ 自动生成配置，减少人工错误
- ✅ 修改代码后自动同步元数据
- ✅ CI 可以验证是否忘记运行 go generate
- ✅ 类型安全（Go 编译器检查）

**CI 验证示例**：
```bash
# .github/workflows/ci.yml
- name: Check generated files are up to date
  run: |
    go generate ./...
    git diff --exit-code || {
      echo "Error: Generated files are out of date."
      echo "Please run 'go generate ./...' and commit the changes."
      exit 1
    }
```

**注意事项**：
- 通用类型（WorkflowMetadata、StageMetadata）定义在 `worker/internal/workflow/types.go`
- 元数据实例（Metadata 变量）定义在各 workflow 的 `workflow.go` 中
- 特殊类型（如果某个 workflow 需要）也定义在各自的 `workflow.go` 中
- 只有 metadata 部分从代码生成，工具模板（tools 部分）仍然手动维护在 templates.yaml 中
- 这样既保证了元数据一致性，又保持了模板的灵活性

### 依赖库

- **YAML 解析**: `gopkg.in/yaml.v3` - 支持 YAML 1.2 和锚点
- **JSON Schema 生成**: `github.com/invopop/jsonschema` - 从 Go struct 生成 JSON Schema
- **JSON Schema 验证**: `github.com/xeipuuv/gojsonschema` - 验证用户配置
- **测试**: `github.com/stretchr/testify` - 断言和 mock
- **属性测试**: `testing/quick` - Go 标准库

### 文件结构

```
worker/
├── internal/
│   ├── activity/
│   │   ├── command_template.go # CommandTemplate, Parameter 结构体
│   │   ├── command_builder.go  # CommandBuilder
│   │   ├── template_loader.go  # TemplateLoader（通用实现，所有 workflow 共享）✅
│   │   └── validator.go        # 参数验证
│   ├── workflow/
│   │   ├── types.go            # 通用类型定义（WorkflowMetadata, StageMetadata）⭐
│   │   └── subdomain_discovery/
│   │       ├── workflow.go         # 工作流入口（Execute、initialize）+ 常量定义
│   │       ├── stages.go           # Metadata 变量 + templateLoader 实例 + 阶段编排逻辑⭐
│   │       ├── templates.yaml      # 工具模板（metadata 部分从 stages.go 生成）
│   │       ├── helpers.go          # buildCommand 辅助函数
│   │       └── stage_*.go          # 各个阶段的具体实现
│   └── ...
├── cmd/
│   ├── schema-gen/
│   │   └── main.go             # Schema 生成工具（待实现）
│   ├── doc-gen/
│   │   └── main.go             # 文档生成工具（待实现）
│   └── gen-metadata/
│       └── main.go             # 元数据生成工具（从 stages.go 生成 templates.yaml）⭐
└── ...

server/
├── configs/
│   └── engines/
│       ├── subdomain_discovery.yaml        # 用户配置
│       └── subdomain_discovery.schema.json # 生成的 Schema（待实现）
├── internal/
│   ├── config/
│   │   └── validator.go        # ConfigValidator（待实现）
│   └── ...
└── ...

docs/
└── config-reference.md         # 生成的配置文档（待实现）
```

**当前实现状态**：
- ✅ `activity/template_loader.go`: 已实现，使用 `embed.FS` + `sync.Once`（通用实现）
- ⏳ 需要增强：添加 `WorkflowMetadata` 支持和阶段依赖验证
- ❌ 待实现：Schema 生成、文档生成、Server 端验证
- ❌ 待删除：`subdomain_discovery/template_loader.go`（多余的封装）

**文件职责说明**：
- `activity/template_loader.go`: 通用模板加载器（所有 workflow 共享）
- `types.go`: 通用类型定义（WorkflowMetadata, StageMetadata）
- `workflow.go`: 工作流入口（Execute、initialize）+ 常量定义
- `stages.go`: Metadata 变量 + templateLoader 实例 + 阶段编排逻辑（runAllStages）
- `templates.yaml`: 工具模板（metadata 部分自动生成）

**为什么不需要 `template_loader.go`？**
- `activity.TemplateLoader` 已经是通用实现
- 每个 workflow 只需在 `stages.go` 中创建实例即可
- 不需要额外的封装层（遵循 YAGNI 原则）

**多 workflow 扩展**：
```
worker/internal/workflow/
├── types.go                    # 通用类型（50-100 行）
├── subdomain_discovery/
│   ├── workflow.go             # 工作流入口 + 常量
│   └── stages.go               # Metadata + templateLoader + 编排逻辑
├── port_scan/
│   ├── workflow.go             # 工作流入口 + 常量
│   └── stages.go               # Metadata + templateLoader + 编排逻辑
└── vulnerability_scan/
    ├── workflow.go             # 工作流入口 + 常量
    └── stages.go               # Metadata + templateLoader + 编排逻辑
```

**stages.go 示例**：
```go
// worker/internal/workflow/subdomain_discovery/stages.go
package subdomain_discovery

import (
    "embed"
    "github.com/orbit/worker/internal/activity"
    "github.com/orbit/worker/internal/workflow"
)

//go:embed templates.yaml
var templatesFS embed.FS

// templateLoader 是子域名发现工作流的模板加载器
var templateLoader = activity.NewTemplateLoader(templatesFS, "templates.yaml")

// Metadata 定义子域名发现工作流的元数据
var Metadata = workflow.WorkflowMetadata{
    Name: "subdomain_discovery",
    // ...
}

// runAllStages 编排所有阶段的执行
func runAllStages(ctx context.Context, config map[string]any) error {
    // 使用 templateLoader.Get("tool_name") 获取模板
    // ...
}
```

### 迁移策略

#### 当前实现状态

**已完成**：
- ✅ `activity/template_loader.go`: 基础模板加载器（使用 `embed.FS` + `sync.Once`）
- ✅ `subdomain_discovery/template_loader.go`: Workflow 特定的加载器封装
- ✅ 基础的模板验证（YAML 语法、Go Template 语法）

**需要增强**：
- ⏳ 添加 `WorkflowMetadata` 支持
- ⏳ 增强验证逻辑（阶段依赖关系）
- ⏳ 添加 `GetMetadata()` 方法

**待实现**：
- ❌ Schema 生成工具
- ❌ 文档生成工具
- ❌ 元数据生成工具
- ❌ Server 端配置验证

#### 阶段 1: 增强 TemplateLoader（1-2 天）

**目标**: 在当前实现基础上添加元数据支持

**任务**:
1. 更新 `worker/internal/activity/template_loader.go`:
   ```go
   type TemplateLoader struct {
       fs       embed.FS
       filename string
       once     sync.Once
       cache    map[string]CommandTemplate
       metadata WorkflowMetadata  // 新增
       err      error
   }
   
   // 新增方法
   func (l *TemplateLoader) GetMetadata() WorkflowMetadata
   func (l *TemplateLoader) ValidateStageDependencies(config map[string]any) error
   ```

2. 更新 `templates.yaml` 格式，添加 metadata 部分:
   ```yaml
   metadata:
     name: "subdomain_discovery"
     display_name: "子域名发现"
     # ...
   
   tools:
     subfinder:
       # ...
   ```

3. 删除多余的封装文件:
   ```bash
   rm worker/internal/workflow/subdomain_discovery/template_loader.go
   ```

4. 在 `stages.go` 中直接创建 templateLoader 实例:
   ```go
   // worker/internal/workflow/subdomain_discovery/stages.go
   //go:embed templates.yaml
   var templatesFS embed.FS
   var templateLoader = activity.NewTemplateLoader(templatesFS, "templates.yaml")
   ```

5. 编写单元测试验证元数据加载

**验证**: 编译通过，单元测试通过，元数据正确加载

#### 阶段 2: 更新数据结构（1-2 天）

**目标**: 将 CommandTemplate 从扁平结构改为嵌套结构

**任务**:
1. 更新 `worker/internal/activity/command_template.go`:
   ```go
   type CommandTemplate struct {
       Metadata    ToolMetadata          `yaml:"metadata"`
       BaseCommand string                `yaml:"base_command"`
       Parameters  map[string]Parameter  `yaml:"parameters"`
   }
   
   type Parameter struct {
       Flag               string      `yaml:"flag"`
       Default            interface{} `yaml:"default"`
       Type               string      `yaml:"type"`
       Required           bool        `yaml:"required"`
       Description        string      `yaml:"description"`
       DeprecationMessage string      `yaml:"deprecation_message,omitempty"`
   }
   ```

2. 更新 `worker/internal/activity/command_builder.go`:
   - 修改 `Build()` 方法处理嵌套参数
   - 实现参数合并逻辑（用户配置 > 默认值）
   - 添加类型验证和转换

3. 编写单元测试验证新结构

**验证**: 编译通过，单元测试通过

#### 阶段 3: 重构模板文件（2-3 天）

**目标**: 将 templates.yaml 改为嵌套结构并使用 YAML 锚点

**任务**:
1. 重构 `worker/internal/workflow/subdomain_discovery/templates.yaml`:
   ```yaml
   # 共享参数定义
   x-common-params: &common-params
     Timeout:
       flag: "-timeout {{.Timeout}}"
       default: 3600
       type: "int"
       required: false
       description: "扫描超时时间（秒）"
   
   # 工具模板
   tools:
     subfinder:
       metadata:
         display_name: "Subfinder"
         description: "使用多个数据源被动收集子域名"
         stage: "passive"
       base_command: "subfinder -d {{.Domain}} -o {{quote .OutputFile}}"
       parameters:
         <<: *common-params
         Threads:
           flag: "-t {{.Threads}}"
           default: 10
           type: "int"
           required: false
           description: "并发线程数"
   ```

2. 更新所有调用点（stage_passive.go, stage_bruteforce.go 等）
3. 编写集成测试验证完整流程

**验证**: 所有现有测试通过，命令构建正确

#### 阶段 4: 实现 Schema 生成（3-4 天）

**目标**: 从 Worker 模板自动生成 JSON Schema

**任务**:
1. 创建 `worker/cmd/schema-gen/main.go`:
   - 读取 templates.yaml
   - 解析为 CommandTemplate 结构
   - 使用 invopop/jsonschema 生成 Schema
   - 输出到 `server/configs/engines/subdomain_discovery.schema.json`

2. 配置 go generate:
   ```go
   //go:generate go run ../../cmd/schema-gen/main.go -input templates.yaml -output ../../../server/configs/engines/subdomain_discovery.schema.json
   ```

3. 编写单元测试验证生成的 Schema 格式正确

**验证**: 运行 `go generate`，生成的 Schema 有效

#### 阶段 5: 实现文档生成（2-3 天）

**目标**: 从 Worker 模板自动生成配置文档

**任务**:
1. 创建 `worker/cmd/doc-gen/main.go`:
   - 读取 templates.yaml
   - 解析为 CommandTemplate 结构
   - 生成 Markdown 表格和示例
   - 输出到 `docs/config-reference.md`

2. 配置 go generate:
   ```go
   //go:generate go run ../../cmd/doc-gen/main.go -input templates.yaml -output ../../../docs/config-reference.md
   ```

3. 编写单元测试验证生成的文档格式正确

**验证**: 运行 `go generate`，生成的文档可读

#### 阶段 6: 实现 Server 端验证（3-4 天）

**目标**: Server 启动时验证用户配置

**任务**:
1. 创建 `server/internal/config/validator.go`:
   ```go
   type ConfigValidator struct {
       schema *jsonschema.Schema
   }
   
   func (v *ConfigValidator) LoadSchema(path string) error
   func (v *ConfigValidator) Validate(config map[string]interface{}) []error
   ```

2. 集成到 Server 启动流程:
   ```go
   // main.go
   validator := config.NewValidator()
   if err := validator.LoadSchema("configs/engines/subdomain_discovery.schema.json"); err != nil {
       log.Fatal(err)
   }
   
   if errs := validator.Validate(engineConfig); len(errs) > 0 {
       for _, err := range errs {
           log.Error(err)
       }
       log.Fatal("Configuration validation failed")
   }
   ```

3. 编写集成测试验证验证流程

**验证**: Server 启动时正确验证配置，错误信息清晰

#### 阶段 7: 迁移现有工具（5-7 天）

**目标**: 将所有现有工具迁移到新格式

**任务**:
1. 逐个迁移工具模板（subfinder, amass, httpx 等）
2. 更新用户配置示例
3. 更新文档说明新的配置格式
4. 运行完整的集成测试

**验证**: 所有工具正常工作，配置验证正确

#### 阶段 8: 清理和优化（2-3 天）

**目标**: 清理旧代码，优化性能

**任务**:
1. 删除旧的命令构建逻辑（如果有）
2. 删除硬编码的默认值
3. 优化模板加载性能
4. 更新所有相关文档
5. 代码审查和重构

**验证**: 代码整洁，性能良好，文档完整

### 重构策略

**内部重构，不需要向后兼容**:
- 这是项目内部的代码重构，不是公开 API
- 可以一次性修改所有使用的地方
- 没有外部依赖需要考虑

**重构步骤**:
1. 增强 `activity.TemplateLoader`（添加元数据支持）
2. 更新 `templates.yaml` 格式
3. 删除多余的 `template_loader.go` 封装文件
4. 在 `stages.go` 中直接使用 `activity.NewTemplateLoader`
5. 更新所有调用点
6. 运行测试验证

### 性能考虑

**模板加载**:
- 使用 `sync.Once` 确保只加载一次
- 缓存解析后的模板结构
- 启动时验证，运行时不重复验证

**Schema 验证**:
- Server 启动时加载 Schema，缓存到内存
- 验证失败快速返回，不影响性能
- 使用高效的 JSON Schema 验证库

**命令构建**:
- 参数合并使用 map 操作，O(n) 复杂度
- 字符串替换使用 strings.ReplaceAll，高效
- 避免不必要的内存分配

### 监控和日志

**关键指标**:
- 模板加载时间
- Schema 验证时间
- 命令构建时间
- 验证失败率

**日志级别**:
- ERROR: 模板加载失败、验证失败
- WARN: 使用废弃参数、参数覆盖
- INFO: 模板加载成功、使用默认值
- DEBUG: 详细的参数合并过程

### 安全考虑

**输入验证**:
- 所有用户输入必须通过 Schema 验证
- 防止命令注入（参数值转义）
- 限制参数值长度和范围

**敏感信息**:
- API key 等敏感参数不记录到日志
- 使用 `***` 替换敏感值
- 考虑使用 secret 管理系统

## Future Extensions

### 外部模板加载支持

**当前实现**: 使用全局配置 + `TemplateSource` 接口

**设计思路**:

#### 1. 全局配置函数（当前实现）

```go
// worker/internal/config/template_config.go
package config

import (
    "embed"
    "net/http"
    "path/filepath"
    "time"
    "github.com/orbit/worker/internal/activity"
    "github.com/spf13/viper"
)

// CreateTemplateLoader 根据全局配置创建模板加载器
func CreateTemplateLoader(embedFS embed.FS, filename string) *activity.TemplateLoader {
    // 从配置文件或环境变量读取
    source := viper.GetString("template.source")  // "embed", "file", "url"
    path := viper.GetString("template.path")
    
    switch source {
    case "file":
        return activity.NewTemplateLoader(&activity.FileSource{
            filepath: filepath.Join(path, filename),
        })
    
    case "url":
        return activity.NewTemplateLoader(&activity.URLSource{
            url:    path + "/" + filename,
            client: &http.Client{Timeout: 10 * time.Second},
        })
    
    default: // "embed" 或空
        return activity.NewTemplateLoader(&activity.EmbedSource{
            fs:       embedFS,
            filename: filename,
        })
    }
}
```

**配置文件**:
```yaml
# worker/config.yaml
template:
  source: "embed"  # 默认使用嵌入文件
  path: ""         # 文件路径或 URL（source 为 embed 时不需要）

# 切换到外部文件
# template:
#   source: "file"
#   path: "/etc/orbit/workflows"

# 切换到远程 URL
# template:
#   source: "url"
#   path: "https://templates.example.com"
```

**环境变量支持**:
```bash
# 使用嵌入文件（默认）
./worker

# 使用外部文件
TEMPLATE_SOURCE=file TEMPLATE_PATH=/etc/orbit/workflows ./worker

# 使用远程 URL
TEMPLATE_SOURCE=url TEMPLATE_PATH=https://templates.example.com ./worker
```

**各个 workflow 使用**:
```go
// worker/internal/workflow/subdomain_discovery/stages.go
package subdomain_discovery

import (
    "embed"
    "github.com/orbit/worker/internal/config"
    "github.com/orbit/worker/internal/workflow"
)

//go:embed templates.yaml
var templatesFS embed.FS

// templateLoader 使用全局配置创建模板加载器
// 自动推断 workflow 名称为 "subdomain_discovery"（从包路径提取）
// 可通过配置文件或环境变量切换数据源（embed/file/url）
var templateLoader = config.CreateTemplateLoader(templatesFS, "templates.yaml")

// Metadata 定义
var Metadata = workflow.WorkflowMetadata{
    Name: "subdomain_discovery",
    // ...
}
```

**自动推断原理**:
- 使用 `runtime.Caller(2)` 获取调用栈（跳过 inferWorkflowName 和 CreateTemplateLoader）
- 从函数名中提取包路径：`github.com/orbit/worker/internal/workflow/subdomain_discovery.init`
- 解析出 workflow 名称：`subdomain_discovery`
- 无需手动传参，完全自动化

**优势**:
- ✅ 所有 workflow 代码一致
- ✅ 切换数据源只需修改配置文件
- ✅ 支持环境变量覆盖
- ✅ 默认使用 embed（向后兼容）

#### 2. 未来演变：数据库配置（付费功能）

当需要支持付费功能、每个 workflow 独立配置时，只需修改 `CreateTemplateLoader` 函数：

```go
// worker/internal/config/template_config.go

// CreateTemplateLoader 根据数据库配置创建模板加载器
func CreateTemplateLoader(embedFS embed.FS, filename string) *activity.TemplateLoader {
    // 自动推断 workflow 名称（从调用栈或包路径）
    workflowName := inferWorkflowName()
    
    // 从 Server API 获取配置
    config := getConfigFromServer(workflowName)
    
    // 根据配置创建对应的 Source（逻辑不变）
    switch config.Source {
    case "file":
        return activity.NewTemplateLoader(&activity.FileSource{
            filepath: config.Path,
        })
    case "url":
        return activity.NewTemplateLoader(&activity.URLSource{
            url:    config.Path,
            client: &http.Client{Timeout: 10 * time.Second},
        })
    case "db":
        return activity.NewTemplateLoader(&activity.DBSource{
            serverURL:    config.ServerURL,
            workflowName: workflowName,
        })
    default:
        return activity.NewTemplateLoader(&activity.EmbedSource{
            fs:       embedFS,
            filename: filename,
        })
    }
}

// inferWorkflowName 自动推断 workflow 名称
func inferWorkflowName() string {
    // 方法 1: 从调用栈获取包路径
    pc, _, _, ok := runtime.Caller(2)
    if !ok {
        return ""
    }
    
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return ""
    }
    
    // 函数名格式：github.com/orbit/worker/internal/workflow/subdomain_discovery.init
    // 提取 subdomain_discovery
    parts := strings.Split(fn.Name(), "/")
    if len(parts) > 0 {
        lastPart := parts[len(parts)-1]
        // 去掉 .init 后缀
        workflowName := strings.Split(lastPart, ".")[0]
        return workflowName
    }
    
    return ""
}

// getConfigFromServer 从 Server API 获取配置
func getConfigFromServer(workflowName string) *WorkflowConfig {
    resp, err := http.Get(fmt.Sprintf("%s/api/workflows/%s/config", 
        viper.GetString("server.url"), workflowName))
    if err != nil {
        // 失败时使用默认配置
        pkg.Logger.Warn("Failed to get workflow config, using default",
            zap.String("workflow", workflowName),
            zap.Error(err))
        return &WorkflowConfig{Source: "embed"}
    }
    defer resp.Body.Close()
    
    var config WorkflowConfig
    json.NewDecoder(resp.Body).Decode(&config)
    return &config
}
```

**各个 workflow 代码完全不需要改**:
```go
// worker/internal/workflow/subdomain_discovery/stages.go

// 完全不变！自动推断 workflow 名称
var templateLoader = config.CreateTemplateLoader(templatesFS, "templates.yaml")
```

**自动推断的实现**:
```go
// worker/internal/config/template_config.go

// inferWorkflowName 自动推断 workflow 名称
func inferWorkflowName() string {
    // 获取调用栈（skip 2: inferWorkflowName -> CreateTemplateLoader -> caller）
    pc, _, _, ok := runtime.Caller(2)
    if !ok {
        return ""
    }
    
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return ""
    }
    
    // 函数名格式：github.com/orbit/worker/internal/workflow/subdomain_discovery.init
    // 提取 subdomain_discovery
    funcName := fn.Name()
    parts := strings.Split(funcName, "/")
    if len(parts) > 0 {
        lastPart := parts[len(parts)-1]
        // 去掉 .init 后缀
        workflowName := strings.Split(lastPart, ".")[0]
        return workflowName
    }
    
    return ""
}
```

**Server 端数据库表**:
```sql
CREATE TABLE workflow_configs (
    id SERIAL PRIMARY KEY,
    workflow_name VARCHAR(100) NOT NULL UNIQUE,
    template_source VARCHAR(20) NOT NULL,  -- 'embed', 'file', 'url', 'db'
    template_path TEXT,
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**前端控制界面**:
- 管理员可以通过前端界面修改 workflow 配置
- 切换免费版/付费版模板
- 上传自定义模板
- Worker 重启后自动应用新配置

**演变路径**:
1. **当前**: 全局配置文件（所有 workflow 共享）
2. **未来**: 数据库配置（每个 workflow 独立，支持付费）
3. **改动**: 只需修改 `config/template_config.go` 一个文件
4. **各个 workflow**: **完全不需要改代码**

**实现时机**: 等需要付费功能时再实现（YAGNI 原则）

**优势**:
- ✅ 架构设计正确，易于演变
- ✅ 改动最小，只改配置层
- ✅ 支持动态配置，前端可控
- ✅ **各个 workflow 代码完全不需要改**
- ✅ 自动推断 workflow 名称，无需手动传参

### 测试策略

**单元测试**:
- CommandTemplate 结构体解析
- CommandBuilder 参数合并
- TemplateLoader 验证逻辑
- SchemaGenerator 生成正确性
- DocGenerator 格式正确性

**集成测试**:
- 完整的命令构建流程
- Schema 生成和验证流程
- 文档生成流程
- Server 启动验证流程

**属性测试**:
- 参数默认值和覆盖
- YAML 锚点解析
- 类型验证和转换
- 错误信息完整性

**性能测试**:
- 模板加载性能
- 命令构建性能
- Schema 验证性能

### 文档更新

**开发者文档**:
- 新的模板格式说明
- Schema 生成流程
- 文档生成流程
- 迁移指南

**用户文档**:
- 配置参考（自动生成）
- 配置示例
- 常见问题
- 故障排查

**API 文档**:
- CommandTemplate 结构
- CommandBuilder 接口
- TemplateLoader 接口
- ConfigValidator 接口
