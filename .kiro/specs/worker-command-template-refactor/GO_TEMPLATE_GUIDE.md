# Go Template 快速参考

本文档提供 Worker 命令模板系统使用 Go Template 的快速参考。

## 为什么选择 Go Template？

### 代码量对比

| 方案 | 代码量 | 功能 |
|------|--------|------|
| 自定义字符串替换 | ~140 行 | 基础替换 |
| **Go Template** | **~50 行** | 替换 + 条件 + 循环 + 函数 |

**减少 64% 代码量，功能更强大！**

### 核心优势

1. ✅ **自动错误检测**: `missingkey=error` 自动检测缺失字段
2. ✅ **详细错误信息**: 提供行号、列号、字段名
3. ✅ **业界标准**: Helm、Kubernetes 都使用
4. ✅ **功能强大**: 支持条件、循环、函数
5. ✅ **减少维护**: 不需要自己实现验证逻辑

## 基础语法

### 占位符

```yaml
# 基础占位符（PascalCase）
base_command: "subfinder -d {{.Domain}} -o {{.OutputFile}}"

# 必需字段
{{.Domain}}      # 目标域名
{{.OutputFile}}  # 输出文件路径
{{.InputFile}}   # 输入文件路径
{{.Target}}      # 目标 URL

# 可选字段
{{.Timeout}}     # 超时时间
{{.Threads}}     # 线程数
{{.Verbose}}     # 详细输出（布尔值）
```

### 函数调用

```yaml
# quote - 自动添加引号
base_command: "subfinder -d {{.Domain}} -o {{quote .OutputFile}}"
# 输出: subfinder -d example.com -o "/tmp/output.txt"

# default - 提供默认值
base_command: "subfinder -d {{.Domain}} -timeout {{default 3600 .Timeout}}"
# 如果 Timeout 未设置，使用 3600

# lower/upper - 大小写转换
base_command: "tool --mode {{lower .Mode}}"
# Mode="DEBUG" -> --mode debug
```

### 条件渲染

```yaml
# if - 条件判断
base_command: "nuclei -u {{.Target}} {{if .Verbose}}-v{{end}}"
# Verbose=true  -> nuclei -u example.com -v
# Verbose=false -> nuclei -u example.com

# if-else
base_command: "tool {{if .Debug}}-debug{{else}}-quiet{{end}}"

# 多条件
base_command: "tool {{if .Verbose}}-v{{end}} {{if .Debug}}-d{{end}}"
```

### 循环（高级）

```yaml
# range - 遍历数组
base_command: "tool {{range .Domains}}-d {{.}} {{end}}"
# Domains=["a.com", "b.com"] -> tool -d a.com -d b.com
```

## 模板文件格式

### 完整示例

```yaml
# 共享参数定义
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
subfinder:
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

# 条件渲染示例
nuclei:
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
```

## 命名规范

### 占位符命名

| 类型 | 格式 | 示例 |
|------|------|------|
| 必需字段 | PascalCase | `{{.Domain}}`, `{{.OutputFile}}` |
| 可选字段 | PascalCase | `{{.Timeout}}`, `{{.Threads}}` |
| 函数调用 | 小写 | `{{quote .Domain}}`, `{{default 3600 .Timeout}}` |

### 参数名映射

用户配置（snake_case）自动映射到 Go Template（PascalCase）：

| 用户配置 | Go Template |
|---------|-------------|
| `timeout` | `{{.Timeout}}` |
| `rate_limit` | `{{.RateLimit}}` |
| `provider_config` | `{{.ProviderConfig}}` |
| `output_file` | `{{.OutputFile}}` |

## 代码实现

### CommandBuilder

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
            "lower": strings.ToLower,
            "upper": strings.ToUpper,
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
        Option("missingkey=error").  // 自动检测缺失字段
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

### 辅助函数

```go
// mergeParameters 合并必需参数、默认值和用户配置
func mergeParameters(
    tmpl CommandTemplate,
    params map[string]any,
    config map[string]any,
) map[string]any {
    result := make(map[string]any)
    
    // 1. 添加必需参数
    for key, value := range params {
        result[key] = value
    }
    
    // 2. 添加可选参数（默认值 + 用户配置）
    for name, param := range tmpl.Parameters {
        if userValue, exists := config[name]; exists {
            result[name] = userValue
        } else if param.Default != nil {
            result[name] = param.Default
        }
    }
    
    return result
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

// convertKeys 转换 snake_case 到 PascalCase
func convertKeys(config map[string]any) map[string]any {
    result := make(map[string]any)
    for key, value := range config {
        pascalKey := snakeToPascal(key)
        result[pascalKey] = value
    }
    return result
}

func snakeToPascal(s string) string {
    parts := strings.Split(s, "_")
    for i, part := range parts {
        if len(part) > 0 {
            parts[i] = strings.ToUpper(part[:1]) + part[1:]
        }
    }
    return strings.Join(parts, "")
}
```

## 错误处理

### Go Template 自动错误检测

```go
// 配置 missingkey=error
t := template.New("command").
    Option("missingkey=error").  // 自动检测缺失字段
    Parse(cmdTemplate)
```

### 错误信息示例

```
# 缺失字段错误
template: command:1:15: executing "command" at <.Domain>: map has no entry for key "Domain"
//                      ↑ 行号:列号              ↑ 字段名                    ↑ 具体错误

# 语法错误
template: command:1: unexpected "}" in operand
//                   ↑ 具体的语法问题

# 函数调用错误
template: command:1:20: executing "command" at <quote .Domain>: error calling quote: invalid argument type
```

## 测试示例

### 单元测试

```go
func TestCommandBuilder_Build(t *testing.T) {
    builder := NewCommandBuilder()
    
    tmpl := CommandTemplate{
        BaseCommand: "subfinder -d {{.Domain}} -o {{quote .OutputFile}}",
        Parameters: map[string]Parameter{
            "Timeout": {
                Flag:    "-timeout {{.Timeout}}",
                Default: 3600,
                Type:    "int",
            },
        },
    }
    
    params := map[string]any{
        "Domain":     "example.com",
        "OutputFile": "/tmp/out.txt",
    }
    
    config := map[string]any{
        "Timeout": 7200,
    }
    
    cmd, err := builder.Build(tmpl, params, config)
    assert.NoError(t, err)
    assert.Equal(t, `subfinder -d example.com -o "/tmp/out.txt" -timeout 7200`, cmd)
}

func TestCommandBuilder_MissingField(t *testing.T) {
    builder := NewCommandBuilder()
    
    tmpl := CommandTemplate{
        BaseCommand: "subfinder -d {{.Domain}}",
        Parameters:  map[string]Parameter{},
    }
    
    params := map[string]any{} // 缺少 Domain
    config := map[string]any{}
    
    _, err := builder.Build(tmpl, params, config)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "Domain")  // 错误信息包含字段名
}
```

## 迁移指南

### 从自定义替换迁移到 Go Template

**步骤 1**: 更新占位符格式

```yaml
# 旧格式
base_command: "subfinder -d {domain} -o '{output-file}'"

# 新格式
base_command: "subfinder -d {{.Domain}} -o {{quote .OutputFile}}"
```

**步骤 2**: 更新参数名

```yaml
# 旧格式（kebab-case）
{domain}
{output-file}
{provider-config}

# 新格式（PascalCase）
{{.Domain}}
{{.OutputFile}}
{{.ProviderConfig}}
```

**步骤 3**: 使用函数替代手动引号

```yaml
# 旧格式
base_command: "subfinder -d {domain} -o '{output-file}'"

# 新格式（使用 quote 函数）
base_command: "subfinder -d {{.Domain}} -o {{quote .OutputFile}}"
```

**步骤 4**: 使用条件替代可选标志

```yaml
# 旧格式（需要在代码中判断）
base_command: "nuclei -u {target}"
# 代码中: if verbose { cmd += " -v" }

# 新格式（模板内条件）
base_command: "nuclei -u {{.Target}} {{if .Verbose}}-v{{end}}"
```

## 常见问题

### Q: 为什么使用 PascalCase 而不是 snake_case？

A: Go Template 中的字段名必须是导出的（首字母大写），所以使用 PascalCase。用户配置仍然可以使用 snake_case，会自动转换。

### Q: 如何处理文件路径中的空格？

A: 使用 `{{quote .OutputFile}}` 函数自动添加引号。

### Q: 如何处理布尔标志？

A: 使用条件渲染：`{{if .Verbose}}-v{{end}}`

### Q: 如何调试模板错误？

A: Go Template 提供详细的错误信息，包含行号、列号和字段名。启动时会验证所有模板语法。

### Q: 性能如何？

A: Go Template 性能优秀，模板会被缓存，执行速度快。比自定义字符串替换更快。

## 参考资源

- [Go text/template 官方文档](https://pkg.go.dev/text/template)
- [Helm Template 指南](https://helm.sh/docs/chart_template_guide/)
- [Kubernetes Go Template](https://kubernetes.io/docs/reference/kubectl/jsonpath/)

---

**版本**: 1.0  
**最后更新**: 2026-01-17
