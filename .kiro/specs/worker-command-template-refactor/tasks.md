# Tasks

## 阶段 1: 更新数据结构

### 1.1 更新 CommandTemplate 结构体

更新 `worker/internal/activity/command_template.go`，将扁平结构改为嵌套结构，使用 Go Template 语法，并添加元数据支持：

```go
type CommandTemplate struct {
    Metadata    ToolMetadata          `yaml:"metadata"`
    BaseCommand string                `yaml:"base_command"` // 使用 {{.Var}} 占位符
    Parameters  map[string]Parameter  `yaml:"parameters"`
}

type Parameter struct {
    Flag               string      `yaml:"flag"`        // 使用 {{.Var}} 占位符
    Default            interface{} `yaml:"default"`
    Type               string      `yaml:"type"`        // "string", "int", "bool"
    Required           bool        `yaml:"required"`
    Description        string      `yaml:"description"`
    DeprecationMessage string      `yaml:"deprecation_message,omitempty"`
}

type ToolMetadata struct {
    DisplayName      string   `yaml:"display_name"`
    Description      string   `yaml:"description"`
    Stage            string   `yaml:"stage"`
    Category         string   `yaml:"category"`
    Homepage         string   `yaml:"homepage"`
    RequiresAPIKeys  bool     `yaml:"requires_api_keys"`
    APIProviders     []string `yaml:"api_providers,omitempty"`
    Warning          string   `yaml:"warning,omitempty"`
    DependsOnStages  []string `yaml:"depends_on_stages,omitempty"`
}

type WorkflowMetadata struct {
    Name         string          `yaml:"name"`
    DisplayName  string          `yaml:"display_name"`
    Description  string          `yaml:"description"`
    Version      string          `yaml:"version"`
    TargetTypes  []string        `yaml:"target_types"`
    Stages       []StageMetadata `yaml:"stages"`
}

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

**验收标准**:
- 结构体定义符合业界标准（GitHub Actions/Terraform 风格）
- 支持所有必需字段（flag, default, type, required, description）
- 支持可选字段（deprecation_message）
- 支持工具元数据（display_name, description, stage, homepage, warning 等）
- 支持 Workflow 元数据（stages, version, target_types）
- 注释说明使用 Go Template 语法
- 编译通过

### 1.2 更新 CommandBuilder

更新 `worker/internal/activity/command_builder.go`，使用 Go Template 替代字符串替换：

**需要实现**:
1. 创建 `funcMap` 提供模板函数（quote, default, lower, upper, join）
2. 合并必需参数和可选参数到一个 data map
3. 应用参数覆盖逻辑（用户配置 > 默认值）
4. 构建完整的命令模板字符串（base_command + 启用的参数 flags）
5. 使用 `text/template` 解析和执行模板
6. 配置 `missingkey=error` 自动检测缺失字段
7. 返回最终命令

**代码示例**:
```go
type CommandBuilder struct {
    funcMap template.FuncMap
}

func NewCommandBuilder() *CommandBuilder {
    return &CommandBuilder{
        funcMap: template.FuncMap{
            "quote": func(s string) string {
                return fmt.Sprintf("%q", s)
            },
            "default": func(def, val interface{}) interface{} {
                if val == nil || val == "" {
                    return def
                }
                return val
            },
        },
    }
}

func (b *CommandBuilder) Build(
    tmpl CommandTemplate,
    params map[string]any,
    config map[string]any,
) (string, error) {
    // 1. 合并数据
    data := mergeParameters(tmpl, params, config)
    
    // 2. 构建完整模板
    cmdTemplate := buildCommandTemplate(tmpl, data)
    
    // 3. 执行 Go Template
    t, err := template.New("command").
        Funcs(b.funcMap).
        Option("missingkey=error").
        Parse(cmdTemplate)
    if err != nil {
        return "", fmt.Errorf("parse template: %w", err)
    }
    
    var buf bytes.Buffer
    if err := t.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("execute template: %w", err)
    }
    
    return buf.String(), nil
}
```

**验收标准**:
- `Build()` 方法使用 Go Template
- 参数覆盖逻辑正确（用户配置 > 默认值）
- 自动检测缺失字段（`missingkey=error`）
- 提供常用模板函数（quote, default 等）
- 错误信息清晰（包含工具名、字段名、位置信息）
- 代码量约 50 行（比自定义替换少 64%）
- 编译通过

### 1.3 编写单元测试

为新的数据结构和 Go Template 实现编写单元测试：

**测试用例**:
1. 解析嵌套参数定义
2. 解析工具元数据（display_name, stage, homepage 等）
3. 解析 Workflow 元数据（stages, version）
4. 参数默认值应用
5. 用户配置覆盖默认值
6. Go Template 占位符替换（`{{.Domain}}`）
7. Go Template 函数调用（`{{quote .OutputFile}}`）
8. Go Template 条件渲染（`{{if .Verbose}}-v{{end}}`）
9. 必需参数验证（自动检测缺失字段）
10. 可选参数处理（无默认值时不添加）
11. 参数范围验证（min/max 约束）
12. 错误信息格式（Go Template 提供详细位置）
13. snake_case 到 PascalCase 转换

**验收标准**:
- 所有测试通过
- 覆盖率 > 90%
- 测试用例清晰易懂
- 验证 Go Template 自动错误检测
- 验证元数据解析正确

## 阶段 2: 重构模板文件

### 2.1 重构 templates.yaml

重构 `worker/internal/workflow/subdomain_discovery/templates.yaml`，使用 Go Template 语法、YAML 锚点和元数据：

**需要实现**:
1. 添加 Workflow 元数据（name, version, stages）
2. 为每个工具添加元数据（display_name, description, stage, homepage, warning）
3. 定义共享参数锚点（`x-common-params`）
4. 将所有工具模板改为 Go Template 语法（`{{.Var}}`）
5. 使用 `<<:` 合并键引用共享参数
6. 添加参数描述和类型定义
7. 添加参数范围约束（min/max）
8. 移除硬编码的默认值到参数定义中
9. 使用 PascalCase 命名占位符

**示例**:
```yaml
# Workflow 元数据
metadata:
  name: "subdomain_discovery"
  display_name: "子域名发现"
  description: "通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名"
  version: "1.0.0"
  target_types: ["domain"]
  stages:
    - id: "passive"
      name: "被动收集"
      description: "使用多个数据源被动收集子域名"
      order: 1
      required: true
      parallel: true
      depends_on: []
      outputs: ["subdomains"]

# 共享参数
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
      category: "passive_collection"
      homepage: "https://github.com/projectdiscovery/subfinder"
      requires_api_keys: true
      api_providers: ["shodan", "censys", "virustotal"]
    
    base_command: "subfinder -d {{.Domain}} -all -o {{quote .OutputFile}} -v"
    parameters:
      <<: *common-params
      Threads:
        flag: "-t {{.Threads}}"
        default: 10
        type: "int"
        required: false
        description: "并发线程数"

  subdomain-bruteforce:
    metadata:
      display_name: "Subdomain Bruteforce"
      description: "使用字典对域名进行 DNS 爆破"
      stage: "bruteforce"
      category: "active_scan"
      homepage: "https://github.com/d3mondev/puredns"
      requires_api_keys: false
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
```

**验收标准**:
- 所有工具模板使用 Go Template 语法
- 占位符使用 PascalCase（`{{.Domain}}`, `{{.OutputFile}}`）
- 共享参数使用 YAML 锚点定义
- 所有参数有明确的类型和描述
- 所有工具有完整的元数据
- Workflow 有完整的元数据（stages 定义）
- 默认值从命令字符串移到参数定义
- 使用 `{{quote}}` 函数处理文件路径
- 使用 `{{if}}` 处理布尔标志
- YAML 语法正确

### 2.2 更新 TemplateLoader

更新 `worker/internal/activity/template_loader.go`，支持新的模板格式、Go Template 验证和元数据：

**需要实现**:
1. 解析 Workflow 元数据（metadata 字段）
2. 解析嵌套参数结构（包含工具元数据）
3. 正确处理 YAML 锚点和合并键
4. 验证参数类型定义（只允许 string, int, bool）
5. 验证默认值类型与参数类型匹配
6. 检查必需参数是否有默认值（逻辑冲突）
7. **验证 Go Template 语法**（尝试解析所有模板字符串）
8. **验证元数据完整性**（必需字段、stage 引用）
9. **验证参数范围约束**（min <= default <= max）
10. **提供阶段依赖验证方法**

**Go Template 验证示例**:
```go
func (l *TemplateLoader) Validate() error {
    // 验证元数据
    if err := l.validateMetadata(); err != nil {
        return err
    }
    
    // 验证工具模板
    for toolName, tmpl := range l.templates {
        // 验证 base_command
        if _, err := template.New("test").Parse(tmpl.BaseCommand); err != nil {
            return fmt.Errorf("tool %s: invalid base_command template: %w", toolName, err)
        }
        
        // 验证每个参数的 flag
        for paramName, param := range tmpl.Parameters {
            if param.Flag != "" {
                if _, err := template.New("test").Parse(param.Flag); err != nil {
                    return fmt.Errorf("tool %s, parameter %s: invalid flag template: %w", 
                        toolName, paramName, err)
                }
            }
            
            // 验证参数范围
            if err := validateParameterRange(param); err != nil {
                return fmt.Errorf("tool %s, parameter %s: %w", toolName, paramName, err)
            }
        }
        
        // 验证工具元数据
        if err := l.validateToolMetadata(toolName, tmpl.Metadata); err != nil {
            return err
        }
    }
    return nil
}

func (l *TemplateLoader) validateMetadata() error {
    if l.metadata.Name == "" {
        return fmt.Errorf("workflow metadata: name is required")
    }
    if l.metadata.Version == "" {
        return fmt.Errorf("workflow metadata: version is required")
    }
    // 验证 stages 定义...
    return nil
}

func (l *TemplateLoader) validateToolMetadata(toolName string, meta ToolMetadata) error {
    if meta.DisplayName == "" {
        return fmt.Errorf("tool %s: display_name is required", toolName)
    }
    if meta.Stage == "" {
        return fmt.Errorf("tool %s: stage is required", toolName)
    }
    // 验证 stage 引用的阶段存在
    if !l.stageExists(meta.Stage) {
        return fmt.Errorf("tool %s: stage %s not defined in workflow metadata", toolName, meta.Stage)
    }
    return nil
}

func validateParameterRange(param Parameter) error {
    // 移除范围验证 - 由工具自己验证
    return nil
}
```

**验收标准**:
- 正确解析嵌套结构
- 正确解析 Workflow 元数据
- 正确解析工具元数据
- YAML 锚点正确合并
- 验证逻辑完整
- **Go Template 语法验证**（启动时检测模板错误）
- **元数据完整性验证**（必需字段、stage 引用）
- 错误信息清晰
- 编译通过

### 2.3 更新调用点

更新所有调用 `buildCommand()` 的地方，适配 Go Template 和新的参数格式：

**文件列表**:
- `worker/internal/workflow/subdomain_discovery/stage_passive.go`
- `worker/internal/workflow/subdomain_discovery/stage_bruteforce.go`
- `worker/internal/workflow/subdomain_discovery/stage_merge.go`
- 其他使用命令构建的文件

**需要修改**:
1. 参数类型从 `map[string]string` 改为 `map[string]any`
2. 参数名从 snake_case 改为 PascalCase（或添加转换函数）
3. 传递正确的必需参数（Domain, OutputFile 等）

**示例**:
```go
// 旧代码
params := map[string]string{
    "domain":      domain,
    "output-file": outputFile,
}
config := map[string]any{
    "timeout": 3600,
    "threads": 10,
}

// 新代码
params := map[string]any{
    "Domain":     domain,
    "OutputFile": outputFile,
}
config := convertKeys(map[string]any{  // snake_case -> PascalCase
    "timeout": 3600,
    "threads": 10,
})
```

**验收标准**:
- 所有调用点更新
- 参数名使用 PascalCase
- 传递正确的参数类型
- 编译通过

### 2.4 编写集成测试

编写集成测试验证完整流程：

**测试场景**:
1. 加载模板 → 构建命令 → 验证命令字符串
2. 使用默认值构建命令
3. 用户配置覆盖默认值
4. YAML 锚点正确合并
5. Go Template 占位符正确替换
6. Go Template 函数正确调用（quote）
7. Go Template 条件正确渲染（if）
8. 错误场景（缺少必需参数、Go Template 语法错误等）
9. snake_case 到 PascalCase 自动转换

**验收标准**:
- 所有集成测试通过
- 覆盖主要使用场景
- 错误场景正确处理
- 验证 Go Template 自动错误检测

## 阶段 3: 实现 Schema 生成

### 3.1 创建 SchemaGenerator

创建 `worker/cmd/schema-gen/main.go`，从 Worker 模板生成 JSON Schema：

**需要实现**:
1. 读取 templates.yaml 文件
2. 解析为 CommandTemplate 结构（包含元数据）
3. 遍历所有工具和参数
4. 生成 JSON Schema 结构：
   - 类型映射：string → string, int → integer, bool → boolean
   - Required 参数添加到 required 数组
   - Description 映射到 description 字段
   - Default 映射到 default 字段
   - **使用工具元数据生成 x-stage, x-homepage, x-warning 等扩展字段**
   - **使用 Workflow 元数据生成 title, description, x-metadata 字段**
5. 输出 JSON 文件

**验收标准**:
- 生成的 Schema 符合 JSON Schema Draft 7 规范
- 所有参数正确映射
- 类型、默认值、描述完整
- 工具元数据映射到扩展字段
- Workflow 元数据映射到顶层字段
- JSON 格式正确

### 3.2 配置 go generate

在 `worker/internal/workflow/subdomain_discovery/template_loader.go` 添加：

```go
//go:generate go run ../../cmd/schema-gen/main.go -input templates.yaml -output ../../../server/configs/engines/subdomain_discovery.schema.json
```

**验收标准**:
- `go generate` 命令正确执行
- Schema 文件生成到正确位置
- 生成的 Schema 有效

### 3.3 编写单元测试

为 SchemaGenerator 编写单元测试：

**测试用例**:
1. 简单工具的 Schema 生成
2. 包含所有类型参数的 Schema 生成
3. 必需参数映射到 required 数组
4. 默认值正确映射
5. 描述正确映射
6. JSON 格式验证

**验收标准**:
- 所有测试通过
- 生成的 Schema 可以被 JSON Schema 验证器加载

## 阶段 4: 实现文档生成

### 4.1 创建 DocGenerator

创建 `worker/cmd/doc-gen/main.go`，从 Worker 模板生成配置文档：

**需要实现**:
1. 读取 templates.yaml 文件
2. 解析为 CommandTemplate 结构（包含元数据）
3. 生成 Markdown 文档：
   - **使用 Workflow 元数据生成概述章节**（名称、版本、描述）
   - **使用 Stage 元数据生成扫描流程章节**（阶段顺序、依赖关系）
   - 每个工具一个章节，**使用工具元数据**（display_name, description, stage, homepage）
   - 参数表格（名称、类型、默认值、必需、描述）
   - 标记废弃的参数
   - **显示警告信息**（如主动扫描警告）
   - **显示 API Keys 需求**
   - 包含使用示例
4. 输出 Markdown 文件

**验收标准**:
- 生成的文档格式正确
- 包含 Workflow 概述章节
- 包含扫描流程章节（阶段顺序、依赖）
- 表格包含所有必需列
- 显示工具元数据（阶段、主页、警告）
- 示例代码可用
- Markdown 语法正确

### 4.2 配置 go generate

在 `worker/internal/workflow/subdomain_discovery/template_loader.go` 添加：

```go
//go:generate go run ../../cmd/doc-gen/main.go -input templates.yaml -output ../../../docs/config-reference.md
```

**验收标准**:
- `go generate` 命令正确执行
- 文档文件生成到正确位置
- 生成的文档可读

### 4.3 编写单元测试

为 DocGenerator 编写单元测试：

**测试用例**:
1. 简单工具的文档生成
2. 包含所有字段的文档生成
3. 表格格式正确
4. 示例代码格式正确
5. Markdown 语法验证

**验收标准**:
- 所有测试通过
- 生成的文档可以被 Markdown 渲染器正确显示

## 阶段 5: 实现 Server 端验证

### 5.1 创建 ConfigValidator

创建 `server/internal/config/validator.go`，实现配置验证器：

**需要实现**:
1. `LoadSchema(schemaPath, metadataPath string) error` - 加载 JSON Schema 和元数据
2. `Validate(config map[string]interface{}) []error` - 验证用户配置
3. `ValidateStageDependencies(config map[string]interface{}) []error` - 验证阶段依赖
4. 使用 JSON Schema 验证库（如 `github.com/xeipuuv/gojsonschema`）
5. 返回详细的验证错误列表

**阶段依赖验证示例**:
```go
func (v *ConfigValidator) ValidateStageDependencies(config map[string]interface{}) []error {
    var errors []error
    
    // 获取所有启用的工具及其所属阶段
    enabledStages := make(map[string]bool)
    for toolName, toolConfig := range config {
        if enabled, ok := toolConfig.(map[string]interface{})["enabled"].(bool); ok && enabled {
            // 从元数据获取工具所属阶段
            stage := v.getToolStage(toolName)
            if stage != "" {
                enabledStages[stage] = true
            }
        }
    }
    
    // 验证每个启用的阶段的依赖
    for stage := range enabledStages {
        stageMeta := v.getStageMetadata(stage)
        for _, dep := range stageMeta.DependsOn {
            if !enabledStages[dep] {
                errors = append(errors, fmt.Errorf(
                    "stage %s depends on %s, but no tools in stage %s are enabled",
                    stage, dep, dep))
            }
        }
    }
    
    return errors
}
```

**验收标准**:
- 正确加载 JSON Schema
- 正确加载元数据
- 验证逻辑正确
- **阶段依赖验证正确**
- 错误信息清晰（包含字段名、期望类型、实际值）
- 编译通过

### 5.2 集成到 Server 启动流程

在 Server 启动代码中集成配置验证：

**需要实现**:
1. 在 Server 启动时加载 Schema
2. 验证用户配置
3. 验证失败时拒绝启动
4. 记录详细的验证错误

**示例**:
```go
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

**验收标准**:
- Server 启动时正确验证配置
- 验证失败时拒绝启动
- 错误信息清晰
- 不影响启动性能

### 5.3 编写集成测试

编写集成测试验证验证流程：

**测试场景**:
1. 有效配置验证通过
2. 类型不匹配的配置验证失败
3. 缺少必需字段的配置验证失败
4. 包含未知字段的配置验证失败
5. **阶段依赖验证**：启用的阶段依赖未启用的阶段（应失败）
6. **阶段依赖验证**：所有依赖都启用（应通过）

**验收标准**:
- 所有集成测试通过
- 覆盖主要验证场景
- 覆盖阶段依赖验证场景
- 错误信息清晰

## 阶段 6: 迁移现有工具

### 6.1 迁移被动收集工具

迁移以下工具到新格式：
- subfinder
- sublist3r
- assetfinder
- amass

**验收标准**:
- 所有工具模板使用嵌套结构
- 参数定义完整（类型、默认值、描述）
- 命令构建正确
- 测试通过

### 6.2 迁移主动扫描工具

迁移以下工具到新格式：
- subdomain-bruteforce
- subdomain-resolve
- subdomain-permutation-resolve

**验收标准**:
- 所有工具模板使用嵌套结构
- 参数定义完整
- 命令构建正确
- 测试通过

### 6.3 迁移其他工具

迁移其他扫描工具（如 httpx, nuclei 等）到新格式。

**验收标准**:
- 所有工具模板使用嵌套结构
- 参数定义完整
- 命令构建正确
- 测试通过

### 6.4 更新用户配置示例

更新 `server/configs/engines/subdomain_discovery.yaml` 中的示例和注释：

**需要更新**:
1. 添加参数说明注释
2. 提供完整的配置示例
3. 说明默认值
4. 说明可选参数

**验收标准**:
- 示例配置清晰易懂
- 注释完整
- 用户可以快速上手

### 6.5 更新文档

更新以下文档：
- `docs/config-reference.md` - 配置参考（自动生成）
- `README.md` - 项目说明
- 其他相关文档

**验收标准**:
- 文档完整准确
- 示例代码可用
- 说明清晰

### 6.6 运行完整测试

运行所有测试验证迁移正确：

**测试类型**:
1. 单元测试
2. 集成测试
3. 属性测试
4. 端到端测试

**验收标准**:
- 所有测试通过
- 覆盖率 > 90%
- 性能无明显下降

## 阶段 7: 清理和优化

### 7.1 清理旧代码

删除不再使用的旧代码：

**需要清理**:
1. 旧的命令构建逻辑（如果有）
2. 硬编码的默认值
3. 废弃的辅助函数
4. 临时的兼容代码

**验收标准**:
- 代码整洁
- 无死代码
- 编译通过

### 7.2 性能优化

优化系统性能：

**优化点**:
1. 模板加载性能（使用 sync.Once）
2. Schema 验证性能（缓存 Schema）
3. 命令构建性能（减少内存分配）
4. 日志性能（异步日志）

**验收标准**:
- 模板加载时间 < 100ms
- Schema 验证时间 < 10ms
- 命令构建时间 < 1ms
- 内存占用不超过重构前的 120%

### 7.3 代码审查

进行代码审查：

**审查内容**:
1. 代码风格一致性
2. 错误处理完整性
3. 日志记录合理性
4. 注释清晰度
5. 测试覆盖率

**验收标准**:
- 代码符合 Go 规范
- 错误处理完整
- 日志合理
- 注释清晰
- 测试覆盖率 > 90%

### 7.4 更新文档

更新所有相关文档：

**文档列表**:
1. 开发者文档（架构、设计、实现）
2. 用户文档（配置、使用、故障排查）
3. API 文档（接口、数据结构）
4. 迁移指南

**验收标准**:
- 文档完整准确
- 示例代码可用
- 说明清晰
- 格式统一

### 7.5 发布准备

准备发布新版本：

**准备工作**:
1. 更新 CHANGELOG
2. 更新版本号
3. 打 tag
4. 构建 Docker 镜像
5. 更新部署文档

**验收标准**:
- CHANGELOG 完整
- 版本号正确
- Docker 镜像可用
- 部署文档准确

## 属性测试任务

### Property 1-2: 默认值和覆盖

编写属性测试验证默认值应用和覆盖逻辑：

**测试策略**:
- 生成随机参数定义和用户配置
- 验证默认值应用正确
- 验证用户配置覆盖默认值

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 1: 参数默认值支持`

### Property 3: YAML 锚点

编写属性测试验证 YAML 锚点解析：

**测试策略**:
- 生成包含锚点的随机 YAML
- 验证解析和合并正确性

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 3: YAML 锚点正确解析`

### Property 4-5: 类型验证和转换

编写属性测试验证类型检查和转换：

**测试策略**:
- 生成随机类型的参数值
- 验证类型检查正确
- 验证类型转换正确

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 4: 参数类型验证`

### Property 6: 模板验证

编写属性测试验证模板验证逻辑：

**测试策略**:
- 生成包含各种错误的模板
- 验证所有错误都能检测

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 6: 模板验证综合`

### Property 7: 错误信息

编写属性测试验证错误信息完整性：

**测试策略**:
- 触发各种错误场景
- 验证错误信息包含必要上下文

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 7: 错误信息完整性`

### Property 8: 缓存

编写属性测试验证模板缓存：

**测试策略**:
- 多次构建命令
- 验证模板只加载一次

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 8: 模板缓存复用`

### Property 9-14: 边界情况

编写属性测试验证边界情况：

**测试场景**:
- 可选参数处理
- 嵌套参数解析
- 必需参数验证
- YAML 锚点边界情况

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 9-14`

### Property 15-19: 新增功能

编写属性测试验证新增功能：

**测试场景**:
- Schema 生成正确性（包含元数据）
- 文档生成完整性（包含阶段流程和工具元数据）
- Server 配置验证
- 阶段依赖验证
- 元数据完整性验证

**验收标准**:
- 最小 100 次迭代
- 所有测试通过
- 标签格式：`Feature: worker-command-template-refactor, Property 15-19`

## 总结

**预计时间**: 2-3 周（比原计划减少 1 周，因为 Go Template 代码量更少）

**关键里程碑**:
1. 阶段 1-2 完成：新的数据结构、Go Template 实现和元数据支持可用
2. 阶段 3-4 完成：Schema 和文档自动生成（使用元数据）
3. 阶段 5 完成：Server 端验证集成（包含阶段依赖验证）
4. 阶段 6 完成：所有工具迁移完成
5. 阶段 7 完成：代码清理和优化，准备发布

**Go Template 方案优势**:
1. **代码量更少**: ~50 行 vs ~140 行（减少 64%）
2. **自动错误检测**: `missingkey=error` 自动检测缺失字段
3. **更强大**: 支持条件、循环、函数等高级特性
4. **业界标准**: Helm、Kubernetes 都使用 Go Template
5. **更好的错误信息**: 提供详细的行号、列号、字段名
6. **减少维护成本**: 不需要自己实现替换、验证、类型转换逻辑

**元数据功能优势**:
1. **自描述**: 配置文件本身包含足够的信息（阶段流程、工具说明）
2. **自动文档**: 从元数据生成文档，保持同步
3. **更好的验证**: 验证阶段依赖和参数范围
4. **更好的 UI**: 前端可以展示流程图和说明
5. **易于维护**: 元数据集中管理，修改一处即可
6. **用户友好**: 清晰的警告信息、API Keys 需求说明

**风险和缓解**:
1. **风险**: Go Template 语法学习曲线
   - **缓解**: 提供详细的文档和示例，语法简单（主要是 `{{.Var}}`）
2. **风险**: YAML 锚点解析复杂
   - **缓解**: 使用成熟的 YAML 库（gopkg.in/yaml.v3）
3. **风险**: 元数据维护成本
   - **缓解**: 元数据验证确保完整性，自动生成文档减少手动维护
4. **风险**: 性能下降
   - **缓解**: 使用缓存，Go Template 性能优秀
5. **风险**: 文档不同步
   - **缓解**: 自动生成文档，集成到 CI/CD
