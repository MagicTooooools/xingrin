# 配置文件元数据设计

## 概述

为了让配置文件更清晰、更易维护，我们在 Worker 模板中添加元数据，描述：
1. Workflow 的整体流程
2. 每个 Stage 的作用和依赖关系
3. 每个 Tool 的用途和参数说明

## 设计原则

1. **单一数据源**：元数据定义在 Worker 模板中，自动生成到文档
2. **自描述**：配置文件本身包含足够的信息，用户无需查看代码
3. **可验证**：元数据可用于验证配置的正确性
4. **业界标准**：参考 GitHub Actions、Terraform 的元数据格式

## 元数据结构

### Worker 模板元数据

```yaml
# worker/internal/workflow/subdomain_discovery/templates.yaml

# Workflow 元数据
metadata:
  name: "subdomain_discovery"
  display_name: "子域名发现"
  description: "通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名"
  version: "1.0.0"
  target_types: ["domain"]  # 支持的目标类型
  
  # 阶段定义
  stages:
    - id: "passive"
      name: "被动收集"
      description: "使用多个数据源被动收集子域名，不产生主动扫描流量"
      order: 1
      required: true
      parallel: true  # 阶段内工具并行执行
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
      depends_on: ["passive", "bruteforce"]  # 依赖前面阶段的输出
      outputs: ["subdomains"]
      
    - id: "resolve"
      name: "DNS 解析验证"
      description: "验证所有发现的子域名是否可解析"
      order: 4
      required: false
      parallel: false
      depends_on: ["passive", "bruteforce", "permutation"]
      outputs: ["subdomains"]

# 共享参数定义
x-common-params: &common-params
  Timeout:
    flag: "-timeout {{.Timeout}}"
    default: 3600
    type: "int"
    required: false
    description: "扫描超时时间（秒）"
    min: 1
    max: 86400
  
  RateLimit:
    flag: "-rl {{.RateLimit}}"
    default: 150
    type: "int"
    required: false
    description: "每秒请求数限制"
    min: 1
    max: 10000

# 工具模板
tools:
  # 被动收集工具
  subfinder:
    metadata:
      display_name: "Subfinder"
      description: "使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名"
      stage: "passive"
      category: "passive_collection"
      homepage: "https://github.com/projectdiscovery/subfinder"
      requires_api_keys: true
      api_providers: ["shodan", "censys", "virustotal", "securitytrails"]
    
    base_command: "subfinder -d {{.Domain}} -all -o {{quote .OutputFile}} -v"
    
    parameters:
      <<: *common-params
      Threads:
        flag: "-t {{.Threads}}"
        default: 10
        type: "int"
        required: false
        description: "并发线程数"
        min: 1
        max: 100
      
      ProviderConfig:
        flag: "-pc {{quote .ProviderConfig}}"
        default: null
        type: "string"
        required: false
        description: "API 提供商配置文件路径（包含 API keys）"

  sublist3r:
    metadata:
      display_name: "Sublist3r"
      description: "使用搜索引擎（Google、Bing、Yahoo 等）被动收集子域名"
      stage: "passive"
      category: "passive_collection"
      homepage: "https://github.com/aboul3la/Sublist3r"
      requires_api_keys: false
    
    base_command: "python3 '/usr/local/share/Sublist3r/sublist3r.py' -d {{.Domain}} -o {{quote .OutputFile}}"
    
    parameters:
      <<: *common-params
      Threads:
        flag: "-t {{.Threads}}"
        default: 10
        type: "int"
        required: false
        description: "并发线程数"
        min: 1
        max: 100

  assetfinder:
    metadata:
      display_name: "Assetfinder"
      description: "使用多个数据源快速查找子域名"
      stage: "passive"
      category: "passive_collection"
      homepage: "https://github.com/tomnomnom/assetfinder"
      requires_api_keys: false
    
    base_command: "assetfinder --subs-only {{.Domain}} > {{quote .OutputFile}}"
    
    parameters:
      <<: *common-params

  # 主动扫描工具
  subdomain-bruteforce:
    metadata:
      display_name: "Subdomain Bruteforce"
      description: "使用字典对域名进行 DNS 爆破，发现未公开的子域名"
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
        min: 1
        max: 1000
      
      RateLimit:
        flag: "--rate-limit {{.RateLimit}}"
        default: 500
        type: "int"
        required: false
        description: "每秒 DNS 请求数限制"
        min: 1
        max: 10000
      
      WildcardTests:
        flag: "--wildcard-tests {{.WildcardTests}}"
        default: 50
        type: "int"
        required: false
        description: "泛解析检测测试次数"
        min: 1
        max: 1000
      
      WildcardBatch:
        flag: "--wildcard-batch {{.WildcardBatch}}"
        default: 1000000
        type: "int"
        required: false
        description: "泛解析检测批次大小"
        min: 1000
        max: 10000000

  subdomain-permutation-resolve:
    metadata:
      display_name: "Subdomain Permutation + Resolve"
      description: "对已发现的子域名进行排列组合，生成新的可能子域名并验证"
      stage: "permutation"
      category: "permutation"
      homepage: "https://github.com/ProjectAnte/dnsgen"
      requires_api_keys: false
      depends_on_stages: ["passive", "bruteforce"]
    
    base_command: "cat {{quote .InputFile}} | dnsgen - | puredns resolve -r {{quote .Resolvers}} --write {{quote .OutputFile}} --quiet"
    
    parameters:
      <<: *common-params
      Threads:
        flag: "-t {{.Threads}}"
        default: 100
        type: "int"
        required: false
        description: "并发线程数"
        min: 1
        max: 1000
      
      RateLimit:
        flag: "--rate-limit {{.RateLimit}}"
        default: 500
        type: "int"
        required: false
        description: "每秒 DNS 请求数限制"
        min: 1
        max: 10000

  subdomain-resolve:
    metadata:
      display_name: "Subdomain Resolve"
      description: "验证所有发现的子域名是否可解析，过滤无效子域名"
      stage: "resolve"
      category: "validation"
      homepage: "https://github.com/d3mondev/puredns"
      requires_api_keys: false
      depends_on_stages: ["passive", "bruteforce", "permutation"]
    
    base_command: "puredns resolve {{quote .InputFile}} -r {{quote .Resolvers}} --write {{quote .OutputFile}} --quiet"
    
    parameters:
      <<: *common-params
      Threads:
        flag: "-t {{.Threads}}"
        default: 100
        type: "int"
        required: false
        description: "并发线程数"
        min: 1
        max: 1000
      
      RateLimit:
        flag: "--rate-limit {{.RateLimit}}"
        default: 500
        type: "int"
        required: false
        description: "每秒 DNS 请求数限制"
        min: 1
        max: 10000
      
      WildcardTests:
        flag: "--wildcard-tests {{.WildcardTests}}"
        default: 50
        type: "int"
        required: false
        description: "泛解析检测测试次数"
        min: 1
        max: 1000
      
      WildcardBatch:
        flag: "--wildcard-batch {{.WildcardBatch}}"
        default: 1000000
        type: "int"
        required: false
        description: "泛解析检测批次大小"
        min: 1000
        max: 10000000
```

## 生成的配置文档

从上述元数据自动生成用户配置文档：

```markdown
# 子域名发现配置参考

## 概述

**名称**: 子域名发现 (subdomain_discovery)  
**版本**: 1.0.0  
**描述**: 通过被动收集、字典爆破、排列组合等方式发现目标域名的所有子域名  
**支持的目标类型**: domain

## 扫描流程

子域名发现包含 4 个阶段，按顺序执行：

### 阶段 1: 被动收集 (passive) [必需]

**描述**: 使用多个数据源被动收集子域名，不产生主动扫描流量  
**执行方式**: 并行执行  
**依赖**: 无  
**输出**: 子域名列表

**可用工具**:
- subfinder - 使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名
- sublist3r - 使用搜索引擎（Google、Bing、Yahoo 等）被动收集子域名
- assetfinder - 使用多个数据源快速查找子域名

### 阶段 2: 字典爆破 (bruteforce) [可选]

**描述**: 使用字典对域名进行爆破，发现未公开的子域名  
**执行方式**: 顺序执行  
**依赖**: 无  
**输出**: 子域名列表  
**警告**: 主动扫描会产生大量 DNS 请求，可能被目标检测

**可用工具**:
- subdomain-bruteforce - 使用字典对域名进行 DNS 爆破

### 阶段 3: 排列组合 (permutation) [可选]

**描述**: 对已发现的子域名进行排列组合，生成新的可能子域名  
**执行方式**: 顺序执行  
**依赖**: passive, bruteforce  
**输出**: 子域名列表

**可用工具**:
- subdomain-permutation-resolve - 对已发现的子域名进行排列组合并验证

### 阶段 4: DNS 解析验证 (resolve) [可选]

**描述**: 验证所有发现的子域名是否可解析  
**执行方式**: 顺序执行  
**依赖**: passive, bruteforce, permutation  
**输出**: 已验证的子域名列表

**可用工具**:
- subdomain-resolve - 验证所有发现的子域名是否可解析

## 工具配置

### subfinder

**描述**: 使用多个数据源（Shodan、Censys、VirusTotal 等）被动收集子域名  
**阶段**: passive  
**主页**: https://github.com/projectdiscovery/subfinder  
**需要 API Keys**: 是（支持 shodan, censys, virustotal, securitytrails）

#### 参数

| 参数 | 类型 | 默认值 | 必需 | 范围 | 描述 |
|------|------|--------|------|------|------|
| timeout | int | 3600 | 否 | 1-86400 | 扫描超时时间（秒） |
| threads | int | 10 | 否 | 1-100 | 并发线程数 |
| provider_config | string | - | 否 | - | API 提供商配置文件路径 |

#### 示例

```yaml
subfinder:
  enabled: true
  timeout: 7200
  threads: 20
```

### subdomain-bruteforce

**描述**: 使用字典对域名进行 DNS 爆破，发现未公开的子域名  
**阶段**: bruteforce  
**主页**: https://github.com/d3mondev/puredns  
**警告**: ⚠️ 主动扫描会产生大量 DNS 请求，可能被目标检测

#### 参数

| 参数 | 类型 | 默认值 | 必需 | 范围 | 描述 |
|------|------|--------|------|------|------|
| timeout | int | 3600 | 否 | 1-86400 | 扫描超时时间（秒） |
| threads | int | 100 | 否 | 1-1000 | 并发线程数 |
| rate_limit | int | 500 | 否 | 1-10000 | 每秒 DNS 请求数限制 |
| wildcard_tests | int | 50 | 否 | 1-1000 | 泛解析检测测试次数 |
| wildcard_batch | int | 1000000 | 否 | 1000-10000000 | 泛解析检测批次大小 |

#### 示例

```yaml
subdomain-bruteforce:
  enabled: true
  timeout: 86400
  threads: 200
  rate_limit: 1000
```

## 完整配置示例

```yaml
# 阶段 1: 被动收集（必需）
passive-tools:
  subfinder:
    enabled: true
    timeout: 7200
    threads: 20
  
  sublist3r:
    enabled: true
    timeout: 3600
  
  assetfinder:
    enabled: true

# 阶段 2: 字典爆破（可选）
bruteforce:
  enabled: false
  subdomain-bruteforce:
    timeout: 86400
    threads: 200

# 阶段 3: 排列组合（可选）
permutation:
  enabled: true
  subdomain-permutation-resolve:
    timeout: 86400

# 阶段 4: DNS 解析验证（可选）
resolve:
  enabled: true
  subdomain-resolve:
    timeout: 86400
```
```

## Server 配置文件增强

在 Server 配置文件中添加注释，引用元数据：

```yaml
# server/configs/engines/subdomain_discovery.yaml

# 子域名发现配置
# 版本: 1.0.0
# 文档: docs/config-reference.md#subdomain_discovery
#
# 扫描流程:
#   1. 被动收集 (passive) - 必需，并行执行
#   2. 字典爆破 (bruteforce) - 可选
#   3. 排列组合 (permutation) - 可选，依赖前面阶段
#   4. DNS 解析验证 (resolve) - 可选，依赖前面阶段

# ============================================================
# 阶段 1: 被动收集 (必需)
# ============================================================
# 使用多个数据源被动收集子域名，不产生主动扫描流量
# 工具并行执行，互不影响

passive-tools:
  # Subfinder - 使用多个数据源（Shodan、Censys 等）
  # 需要 API Keys 以获得最佳效果
  subfinder:
    enabled: true
    timeout: 3600  # 1 小时
    # threads: 10  # 可选，默认 10
  
  # Sublist3r - 使用搜索引擎
  sublist3r:
    enabled: true
    timeout: 3600
  
  # Assetfinder - 快速查找
  assetfinder:
    enabled: true
    timeout: 3600

# ============================================================
# 阶段 2: 字典爆破 (可选)
# ============================================================
# 使用字典对域名进行 DNS 爆破
# ⚠️ 警告: 主动扫描会产生大量 DNS 请求

bruteforce:
  enabled: false  # 默认禁用
  subdomain-bruteforce:
    timeout: 86400  # 24 小时
    wordlist-name: subdomains-top1million-110000.txt
    # threads: 100  # 可选，默认 100
    # rate_limit: 500  # 可选，默认 500

# ============================================================
# 阶段 3: 排列组合 (可选)
# ============================================================
# 对已发现的子域名进行排列组合
# 依赖: passive, bruteforce

permutation:
  enabled: true
  subdomain-permutation-resolve:
    timeout: 86400

# ============================================================
# 阶段 4: DNS 解析验证 (可选)
# ============================================================
# 验证所有发现的子域名是否可解析
# 依赖: passive, bruteforce, permutation

resolve:
  enabled: true
  subdomain-resolve:
    timeout: 86400
```

## 元数据的用途

### 1. 自动生成文档

从元数据生成：
- 配置参考文档（Markdown）
- JSON Schema（用于验证）
- API 文档（Swagger/OpenAPI）

### 2. 配置验证

```go
// 验证阶段依赖
func ValidateStageDependencies(config map[string]any, metadata Metadata) error {
    for _, stage := range metadata.Stages {
        if !isStageEnabled(config, stage.ID) {
            continue
        }
        
        // 检查依赖的阶段是否启用
        for _, dep := range stage.DependsOn {
            if !isStageEnabled(config, dep) {
                return fmt.Errorf("stage %s depends on %s, but %s is not enabled", 
                    stage.ID, dep, dep)
            }
        }
    }
    return nil
}
```

### 3. UI 展示

前端可以使用元数据：
- 显示阶段流程图
- 显示工具说明和链接
- 显示参数范围和默认值
- 显示警告信息

### 4. 参数范围验证

```go
// 验证参数范围
func ValidateParameterRange(value int, param Parameter) error {
    if param.Min != nil && value < *param.Min {
        return fmt.Errorf("value %d is less than minimum %d", value, *param.Min)
    }
    if param.Max != nil && value > *param.Max {
        return fmt.Errorf("value %d is greater than maximum %d", value, *param.Max)
    }
    return nil
}
```

## 实施建议

### 阶段 1: 添加基础元数据（1-2 天）

1. 在 Worker 模板中添加 `metadata` 和 `tools.*.metadata` 字段
2. 为每个工具添加基本元数据（display_name, description, stage）
3. 更新 TemplateLoader 解析元数据

### 阶段 2: 生成文档（1-2 天）

1. 更新 DocGenerator 使用元数据生成文档
2. 生成包含阶段流程的文档
3. 生成包含工具说明的文档

### 阶段 3: 配置验证（2-3 天）

1. 实现阶段依赖验证
2. 实现参数范围验证
3. 集成到 Server 启动流程

### 阶段 4: UI 集成（可选，3-5 天）

1. 前端读取元数据
2. 显示阶段流程图
3. 显示工具说明和参数范围

## 总结

添加元数据的好处：

1. ✅ **自描述**: 配置文件本身包含足够的信息
2. ✅ **自动文档**: 从元数据生成文档，保持同步
3. ✅ **更好的验证**: 验证阶段依赖和参数范围
4. ✅ **更好的 UI**: 前端可以展示流程图和说明
5. ✅ **易于维护**: 元数据集中管理，修改一处即可

需要我更新 design.md 和 tasks.md 来包含元数据功能吗？
