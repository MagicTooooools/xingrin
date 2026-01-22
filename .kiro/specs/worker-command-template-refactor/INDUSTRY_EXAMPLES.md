# 业界标准实现示例研究

本文档收集了 GitHub Actions、Terraform、Helm、Kubernetes 等项目的参数定义、验证、文档生成等方面的实现示例，为 Worker 命令模板重构提供参考。

## 目录

1. [GitHub Actions](#1-github-actions)
2. [Terraform](#2-terraform)
3. [Helm Charts](#3-helm-charts)
4. [Kubernetes CRD](#4-kubernetes-crd)
5. [JSON Schema 生成](#5-json-schema-生成)
6. [配置文档生成](#6-配置文档生成)
7. [YAML 锚点与别名](#7-yaml-锚点与别名)
8. [参数默认值管理](#8-参数默认值管理)
9. [Go Template 最佳实践](#9-go-template-最佳实践)
10. [对比分析](#10-对比分析)

---

## 1. GitHub Actions

### 1.1 官方文档

- [Metadata syntax for GitHub Actions](https://docs.github.com/en/actions/sharing-automations/creating-actions/metadata-syntax-for-github-actions)

### 1.2 action.yml 格式

GitHub Actions 使用 `action.yml` 定义 Action 的元数据、输入参数和输出。

**核心特性**：
- 扁平的参数定义结构
- 支持 `deprecationMessage` 标记废弃参数
- 明确的 `required` 和 `default` 字段
- 描述性的 `description` 字段

### 1.3 真实示例：actions/checkout

```yaml
name: 'Checkout'
description: 'Checkout a Git repository at a particular version'

inputs:
  repository:
    description: 'Repository name with owner. For example, actions/checkout'
    default: ${{ github.repository }}
  ref:
    description: 'The branch, tag or SHA to checkout'
  token:
    description: 'Personal access token (PAT) used to fetch the repository'
    default: ${{ github.token }}
  ssh-key:
    description: 'SSH key used to fetch the repository'
    default: ''
  persist-credentials:
    description: 'Whether to configure the token or SSH key with the local git config'
    default: 'true'
  path:
    description: 'Relative path under $GITHUB_WORKSPACE to place the repository'
  clean:
    description: 'Whether to execute `git clean -ffdx && git reset --hard HEAD` before fetching'
    default: 'true'
  fetch-depth:
    description: 'Number of commits to fetch. 0 indicates all history for all branches and tags'
    default: '1'
  fetch-tags:
    description: 'Whether to fetch tags, even if fetch-depth > 0'
    default: 'false'
  lfs:
    description: 'Whether to download Git-LFS files'
    default: 'false'
  submodules:
    description: 'Whether to checkout submodules: `true` to checkout submodules or `recursive` to recursively checkout submodules'
    default: 'false'
```

**来源**: [actions/checkout](https://github.com/actions/checkout)


### 1.4 deprecationMessage 示例

```yaml
inputs:
  num-octocats:
    description: 'Number of Octocats'
    required: false
    default: '1'
    deprecationMessage: 'This input will be removed in the next major version. Use octocat-count instead.'
  octocat-eye-color:
    description: 'Eye color of the Octocats'
    required: true
```

**来源**: [GitHub Actions 官方文档](https://docs.github.com/en/actions/sharing-automations/creating-actions/metadata-syntax-for-github-actions)

### 1.5 最佳实践总结

| 实践 | 说明 |
|------|------|
| **扁平结构** | 所有输入参数在同一层级，易于理解 |
| **明确默认值** | 每个可选参数都有明确的 `default` 值 |
| **废弃标记** | 使用 `deprecationMessage` 提示用户迁移 |
| **类型约束** | 通过描述说明期望的类型（如 `'true'` 表示字符串布尔值）|
| **环境变量** | 输入自动转换为 `INPUT_<NAME>` 环境变量 |

---

## 2. Terraform

### 2.1 官方文档

- [Input Variables](https://www.terraform.io/language/values/variables)
- [Variable Validation](https://developer.hashicorp.com/terraform/language/values/variables#custom-validation-rules)
- [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)

### 2.2 Variable 定义结构

Terraform 使用 HCL 语法定义变量，支持类型约束、验证规则和敏感数据标记。


### 2.3 真实示例：变量验证

```hcl
# 字符串验证 - 只允许特定值
variable "string_only_valid_options" {
  type        = string
  default     = "approved"
  description = "Approval status"
  
  # 使用 regex 验证
  validation {
    condition     = can(regex("^(approved|disapproved)$", var.string_only_valid_options))
    error_message = "Invalid input, options: \"approved\", \"disapproved\"."
  }
  
  # 或使用 contains() 函数
  validation {
    condition     = contains(["approved", "disapproved"], var.string_only_valid_options)
    error_message = "Invalid input, options: \"approved\", \"disapproved\"."
  }
}

# AWS Region 验证
variable "string_like_aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region for deployment"
  
  validation {
    condition     = can(regex("[a-z][a-z]-[a-z]+-[1-9]", var.string_like_aws_region))
    error_message = "Must be valid AWS Region names."
  }
}

# IAM Role 名称验证
variable "string_valid_iam_role_name" {
  type        = string
  default     = "MyCoolRole"
  description = "IAM role name"
  
  validation {
    condition     = can(regex("^[a-zA-Z][a-zA-Z\\-\\_0-9]{1,64}$", var.string_valid_iam_role_name))
    error_message = "IAM role name must start with letter, only contain letters, numbers, dashes, or underscores and must be between 1 and 64 characters."
  }
}


# IPv4 CIDR 验证
variable "string_like_valid_ipv4_cidr" {
  type        = string
  default     = "10.0.0.0/16"
  description = "IPv4 CIDR block"
  
  validation {
    condition     = can(cidrhost(var.string_like_valid_ipv4_cidr, 0))
    error_message = "Must be valid IPv4 CIDR."
  }
}

# 数字范围验证
variable "num_in_range" {
  type        = number
  default     = 1
  description = "Number of instances"
  
  validation {
    condition     = var.num_in_range >= 1 && var.num_in_range <= 16 && floor(var.num_in_range) == var.num_in_range
    error_message = "Accepted values: 1-16."
  }
}

# 敏感数据标记
variable "database_password" {
  type        = string
  description = "Database password"
  sensitive   = true  # 不会在日志中显示
}
```

**来源**: [Terraform Variable Validation Examples](https://dev.to/drewmullen/terraform-variable-validation-with-samples-1ank)

### 2.4 文档自动生成

Terraform 使用 `terraform-plugin-docs` 工具自动生成文档：

```bash
# 生成文档
tfplugindocs generate

# 验证文档
tfplugindocs validate
```

**工作流程**：
1. 从 provider schema 提取信息（`terraform providers schema -json`）
2. 读取 `examples/` 目录中的示例代码
3. 读取 `templates/` 目录中的模板文件
4. 生成 Markdown 文档


**目录结构**：
```
provider/
├── examples/
│   ├── resources/
│   │   └── example_resource/
│   │       └── resource.tf
│   └── data-sources/
│       └── example_data_source/
│           └── data-source.tf
├── templates/
│   ├── resources/
│   │   └── example_resource.md.tmpl
│   └── data-sources/
│       └── example_data_source.md.tmpl
└── docs/
    ├── resources/
    │   └── example_resource.md  # 自动生成
    └── data-sources/
        └── example_data_source.md  # 自动生成
```

**来源**: [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)

### 2.5 最佳实践总结

| 实践 | 说明 |
|------|------|
| **类型系统** | 强类型约束（string, number, bool, list, map, object）|
| **验证规则** | 使用 `validation` 块定义自定义验证逻辑 |
| **清晰错误** | `error_message` 提供明确的错误提示 |
| **敏感数据** | `sensitive = true` 防止敏感信息泄露 |
| **自动文档** | 从 schema 和示例自动生成文档 |
| **版本化** | 支持多版本 schema，平滑迁移 |

---

## 3. Helm Charts

### 3.1 官方文档

- [Helm Values Schema](https://helm.sh/docs/topics/charts/#schema-files)
- [JSON Schema Specification](https://json-schema.org/)

### 3.2 values.schema.json 格式

Helm 3 支持使用 JSON Schema 验证 `values.yaml` 文件。


### 3.3 真实示例：values.schema.json

```json
{
  "$schema": "http://json-schema.org/schema#",
  "type": "object",
  "required": ["image"],
  "properties": {
    "image": {
      "type": "object",
      "required": ["repository", "pullPolicy"],
      "properties": {
        "repository": {
          "type": "string",
          "pattern": "^[a-z0-9-_]+$",
          "description": "Docker image repository"
        },
        "pullPolicy": {
          "type": "string",
          "pattern": "^(Always|Never|IfNotPresent)$",
          "description": "Image pull policy"
        },
        "tag": {
          "type": "string",
          "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$",
          "description": "Image tag (semantic version)"
        }
      }
    },
    "replicaCount": {
      "type": "integer",
      "minimum": 1,
      "maximum": 10,
      "default": 1,
      "description": "Number of replicas"
    },
    "resources": {
      "type": "object",
      "properties": {
        "limits": {
          "type": "object",
          "properties": {
            "cpu": {
              "type": "string",
              "pattern": "^[0-9]+m?$"
            },
            "memory": {
              "type": "string",
              "pattern": "^[0-9]+(Mi|Gi)$"
            }
          }
        }
      }
    }
  }
}
```

**来源**: [Validating Helm Chart Values with JSON Schemas](https://www.arthurkoziel.com/validate-helm-chart-values-with-json-schemas/)


### 3.4 Schema 验证时机

Helm 在以下命令中自动验证 schema：
- `helm install`
- `helm upgrade`
- `helm lint`
- `helm template`

**错误示例**：
```bash
$ helm lint .

==> Linting .
[ERROR] values.yaml: 
- image.repository: Invalid type. Expected: string, given: null
- image.pullPolicy: Does not match pattern '^(Always|Never|IfNotPresent)$'

Error: 1 chart(s) linted, 1 chart(s) failed
```

### 3.5 Schema 生成工具

**方法 1：从 values.yaml 推断**
1. 将 `values.yaml` 转换为 JSON：https://www.json2yaml.com/
2. 使用 JSON Schema 生成器：https://www.jsonschema.net/
3. 手动调整生成的 schema

**方法 2：使用 helm-schema-gen**
```bash
# 安装
go install github.com/karuppiah7890/helm-schema-gen@latest

# 生成 schema
helm-schema-gen values.yaml > values.schema.json
```

### 3.6 最佳实践总结

| 实践 | 说明 |
|------|------|
| **类型验证** | 使用 `type` 字段确保数据类型正确 |
| **模式匹配** | 使用 `pattern` 正则表达式验证格式 |
| **范围限制** | 使用 `minimum`/`maximum` 限制数值范围 |
| **必填字段** | 使用 `required` 数组标记必填字段 |
| **默认值** | 在 schema 中定义 `default` 值 |
| **描述信息** | 使用 `description` 提供字段说明 |

---

## 4. Kubernetes CRD

### 4.1 官方文档

- [Kubebuilder Book - Generating CRDs](https://www.kubebuilder.io/reference/generating-crd)
- [CRD Validation Markers](https://book.kubebuilder.io/reference/markers/crd-validation)

### 4.2 Kubebuilder Markers

Kubernetes 使用 Kubebuilder 的 marker 注释从 Go 代码生成 CRD。

### 4.3 真实示例：CRD 验证

```go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ToySpec defines the desired state of Toy
type ToySpec struct {
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	Knights []string `json:"knights,omitempty"`

	// +kubebuilder:validation:Enum=Lion;Wolf;Dragon
	Alias Alias `json:"alias,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:ExclusiveMaximum=false
	Rank Rank `json:"rank"`
}

// +kubebuilder:validation:Enum=Lion;Wolf;Dragon
type Alias string

// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=3
type Rank int32

// +kubebuilder:printcolumn:name="Alias",type=string,JSONPath=`.spec.alias`
// +kubebuilder:printcolumn:name="Rank",type=integer,JSONPath=`.spec.rank`
// +kubebuilder:subresource:status
type Toy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ToySpec   `json:"spec,omitempty"`
	Status ToyStatus `json:"status,omitempty"`
}
```

**来源**: [Kubebuilder Documentation](https://www.kubebuilder.io/reference/generating-crd)


### 4.4 常用验证 Markers

| Marker | 用途 | 示例 |
|--------|------|------|
| `+kubebuilder:validation:Required` | 必填字段 | `// +kubebuilder:validation:Required` |
| `+kubebuilder:validation:Optional` | 可选字段 | `// +kubebuilder:validation:Optional` |
| `+kubebuilder:validation:Enum` | 枚举值 | `// +kubebuilder:validation:Enum=A;B;C` |
| `+kubebuilder:validation:Minimum` | 最小值 | `// +kubebuilder:validation:Minimum=0` |
| `+kubebuilder:validation:Maximum` | 最大值 | `// +kubebuilder:validation:Maximum=100` |
| `+kubebuilder:validation:MinLength` | 最小长度 | `// +kubebuilder:validation:MinLength=1` |
| `+kubebuilder:validation:MaxLength` | 最大长度 | `// +kubebuilder:validation:MaxLength=255` |
| `+kubebuilder:validation:Pattern` | 正则匹配 | `// +kubebuilder:validation:Pattern="^[a-z]+$"` |
| `+kubebuilder:validation:MinItems` | 数组最小元素数 | `// +kubebuilder:validation:MinItems=1` |
| `+kubebuilder:validation:MaxItems` | 数组最大元素数 | `// +kubebuilder:validation:MaxItems=100` |
| `+kubebuilder:validation:UniqueItems` | 数组元素唯一 | `// +kubebuilder:validation:UniqueItems=true` |
| `+kubebuilder:default` | 默认值 | `// +kubebuilder:default=1` |

### 4.5 CEL 验证规则（Kubernetes 1.25+）

```go
type ReplicaSpec struct {
	// +kubebuilder:validation:XValidation:rule="self.minReplicas <= self.replicas",message="replicas should be in the range minReplicas..maxReplicas"
	MinReplicas int32 `json:"minReplicas"`
	Replicas    int32 `json:"replicas"`
	MaxReplicas int32 `json:"maxReplicas"`
}
```

**验证规则示例**：

| 规则 | 用途 |
|------|------|
| `self.minReplicas <= self.replicas` | 验证整数字段关系 |
| `'Available' in self.stateCounts` | 验证 map 中存在特定键 |
| `self.set1.all(e, !(e in self.set2))` | 验证两个集合不相交 |
| `self == oldSelf` | 验证字段不可变 |
| `self.created + self.ttl < self.expired` | 验证日期关系 |

### 4.6 最佳实践总结

| 实践 | 说明 |
|------|------|
| **代码即文档** | 从 Go 代码生成 OpenAPI Schema |
| **类型安全** | 利用 Go 类型系统确保正确性 |
| **声明式验证** | 使用 marker 注释声明验证规则 |
| **复杂验证** | 使用 CEL 表达式实现复杂逻辑 |
| **版本化** | 支持多版本 CRD，平滑升级 |

---

## 5. JSON Schema 生成

### 5.1 Go 库：invopop/jsonschema

**官方仓库**: [invopop/jsonschema](https://github.com/invopop/jsonschema)

### 5.2 从 Go Struct 生成 Schema

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
)

type Config struct {
	// 基本类型
	Name        string `json:"name" jsonschema:"required,minLength=1,maxLength=100,description=Configuration name"`
	Enabled     bool   `json:"enabled" jsonschema:"default=true,description=Enable this configuration"`
	Timeout     int    `json:"timeout" jsonschema:"minimum=1,maximum=3600,description=Timeout in seconds"`
	
	// 枚举
	Level       string `json:"level" jsonschema:"enum=debug,enum=info,enum=warn,enum=error,description=Log level"`
	
	// 数组
	Tags        []string `json:"tags" jsonschema:"minItems=0,maxItems=10,uniqueItems=true,description=Configuration tags"`
	
	// 嵌套对象
	Database    DatabaseConfig `json:"database" jsonschema:"required,description=Database configuration"`
	
	// 可选字段
	Description *string `json:"description,omitempty" jsonschema:"description=Optional description"`
}

type DatabaseConfig struct {
	Host     string `json:"host" jsonschema:"required,format=hostname,description=Database host"`
	Port     int    `json:"port" jsonschema:"minimum=1,maximum=65535,default=5432,description=Database port"`
	Username string `json:"username" jsonschema:"required,minLength=1,description=Database username"`
	Password string `json:"password" jsonschema:"required,minLength=8,description=Database password"`
	SSLMode  string `json:"sslMode" jsonschema:"enum=disable,enum=require,enum=verify-ca,enum=verify-full,default=require"`
}

func main() {
	// 生成 schema
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(&Config{})
	
	// 输出 JSON
	data, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(data))
}
```


### 5.3 生成的 Schema 示例

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["name", "database"],
  "additionalProperties": false,
  "properties": {
    "name": {
      "type": "string",
      "minLength": 1,
      "maxLength": 100,
      "description": "Configuration name"
    },
    "enabled": {
      "type": "boolean",
      "default": true,
      "description": "Enable this configuration"
    },
    "timeout": {
      "type": "integer",
      "minimum": 1,
      "maximum": 3600,
      "description": "Timeout in seconds"
    },
    "level": {
      "type": "string",
      "enum": ["debug", "info", "warn", "error"],
      "description": "Log level"
    },
    "tags": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 0,
      "maxItems": 10,
      "uniqueItems": true,
      "description": "Configuration tags"
    },
    "database": {
      "type": "object",
      "required": ["host", "username", "password"],
      "description": "Database configuration",
      "properties": {
        "host": {
          "type": "string",
          "format": "hostname",
          "description": "Database host"
        },
        "port": {
          "type": "integer",
          "minimum": 1,
          "maximum": 65535,
          "default": 5432,
          "description": "Database port"
        },
        "username": {
          "type": "string",
          "minLength": 1,
          "description": "Database username"
        },
        "password": {
          "type": "string",
          "minLength": 8,
          "description": "Database password"
        },
        "sslMode": {
          "type": "string",
          "enum": ["disable", "require", "verify-ca", "verify-full"],
          "default": "require"
        }
      }
    }
  }
}
```

### 5.4 支持的 jsonschema Tags

| Tag | 说明 | 示例 |
|-----|------|------|
| `required` | 必填字段 | `jsonschema:"required"` |
| `minLength` | 最小长度 | `jsonschema:"minLength=1"` |
| `maxLength` | 最大长度 | `jsonschema:"maxLength=100"` |
| `minimum` | 最小值 | `jsonschema:"minimum=0"` |
| `maximum` | 最大值 | `jsonschema:"maximum=100"` |
| `enum` | 枚举值 | `jsonschema:"enum=a,enum=b"` |
| `pattern` | 正则表达式 | `jsonschema:"pattern=^[a-z]+$"` |
| `format` | 格式验证 | `jsonschema:"format=email"` |
| `default` | 默认值 | `jsonschema:"default=true"` |
| `description` | 字段描述 | `jsonschema:"description=Field description"` |
| `minItems` | 数组最小元素 | `jsonschema:"minItems=1"` |
| `maxItems` | 数组最大元素 | `jsonschema:"maxItems=10"` |
| `uniqueItems` | 数组元素唯一 | `jsonschema:"uniqueItems=true"` |

### 5.5 最佳实践总结

| 实践 | 说明 |
|------|------|
| **单一数据源** | 从代码生成 schema，避免手动维护 |
| **类型安全** | 利用 Go 类型系统确保正确性 |
| **标签驱动** | 使用 struct tags 声明验证规则 |
| **自动化** | 集成到构建流程，自动生成 schema |
| **版本控制** | 将生成的 schema 纳入版本控制 |

---

## 6. 配置文档生成

### 6.1 terraform-plugin-docs 工作流程

```
┌─────────────────────────────────────────────────────────────┐
│ 1. 提取 Schema                                               │
│    terraform providers schema -json                         │
│    ↓                                                         │
│    {                                                         │
│      "provider_schemas": {                                   │
│        "registry.terraform.io/example/example": {            │
│          "provider": { "block": { "attributes": {...} } },   │
│          "resource_schemas": {...},                          │
│          "data_source_schemas": {...}                        │
│        }                                                     │
│      }                                                       │
│    }                                                         │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. 读取示例代码                                              │
│    examples/resources/example_resource/resource.tf          │
│    examples/data-sources/example_data_source/data-source.tf │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. 读取模板文件                                              │
│    templates/resources/example_resource.md.tmpl             │
│    templates/data-sources/example_data_source.md.tmpl       │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. 生成 Markdown 文档                                        │
│    docs/resources/example_resource.md                       │
│    docs/data-sources/example_data_source.md                 │
└─────────────────────────────────────────────────────────────┘
```

### 6.2 模板文件示例

```markdown
---
page_title: "{{.Type}} {{.Name}} - {{.ProviderName}}"
subcategory: ""
description: |-
  {{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Type}}: {{.Name}}

{{ .Description | trimspace }}

## Example Usage

{{ tffile "examples/resources/example_resource/resource.tf" }}

{{ .SchemaMarkdown | trimspace }}

## Import

Import is supported using the following syntax:

{{ codefile "shell" "examples/resources/example_resource/import.sh" }}
```

### 6.3 自动生成的文档结构

```markdown
# Resource: example_resource

Creates and manages an example resource.

## Example Usage

```hcl
resource "example_resource" "example" {
  name        = "my-resource"
  description = "Example resource"
  enabled     = true
  
  tags = {
    Environment = "production"
  }
}
```

## Argument Reference

- `name` (String, Required) The name of the resource.
- `description` (String, Optional) A description of the resource.
- `enabled` (Boolean, Optional) Whether the resource is enabled. Defaults to `true`.
- `tags` (Map of String, Optional) A map of tags to assign to the resource.

## Attribute Reference

- `id` (String) The ID of the resource.
- `created_at` (String) The timestamp when the resource was created.
- `updated_at` (String) The timestamp when the resource was last updated.

## Import

Resources can be imported using the `id`:

```shell
terraform import example_resource.example resource-id
```
```

### 6.4 最佳实践总结

| 实践 | 说明 |
|------|------|
| **Schema 驱动** | 从 provider schema 自动提取参数信息 |
| **示例优先** | 提供真实可运行的示例代码 |
| **模板化** | 使用模板统一文档格式 |
| **自动化** | 集成到 CI/CD，自动更新文档 |
| **版本同步** | 文档与代码版本保持一致 |

---

## 7. YAML 锚点与别名

### 7.1 基本语法

```yaml
# 定义锚点（使用 &）
default_logging: &default_logging
  level: info
  format: json
  output: stdout

# 引用锚点（使用 *）
service_a:
  name: service-a
  logging: *default_logging

service_b:
  name: service-b
  logging: *default_logging
```

### 7.2 Docker Compose 示例

```yaml
version: '3.8'

# 定义可复用的配置块
x-common-variables: &common-variables
  ENVIRONMENT: production
  LOG_LEVEL: info
  TZ: UTC

x-healthcheck: &healthcheck
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s

x-logging: &logging
  driver: json-file
  options:
    max-size: "10m"
    max-file: "3"

services:
  web:
    image: nginx:latest
    environment:
      <<: *common-variables
      SERVICE_NAME: web
    healthcheck:
      <<: *healthcheck
      test: ["CMD", "curl", "-f", "http://localhost"]
    logging: *logging

  api:
    image: myapp:latest
    environment:
      <<: *common-variables
      SERVICE_NAME: api
      DATABASE_URL: postgres://db:5432/myapp
    healthcheck:
      <<: *healthcheck
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
    logging: *logging

  worker:
    image: myapp:latest
    environment:
      <<: *common-variables
      SERVICE_NAME: worker
    logging: *logging
```

**来源**: [Docker Compose Anchors and Aliases](https://medium.com/@kinghuang/docker-compose-anchors-aliases-extensions-a1e4105d70bd)


### 7.3 Helm Values 示例

```yaml
# values.yaml

# 定义通用配置
_common_resources: &common_resources
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi

_common_security_context: &common_security_context
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
  capabilities:
    drop:
      - ALL

_common_labels: &common_labels
  app.kubernetes.io/managed-by: helm
  app.kubernetes.io/part-of: myapp

# 使用锚点
frontend:
  replicaCount: 3
  image:
    repository: myapp/frontend
    tag: "1.0.0"
  resources: *common_resources
  securityContext: *common_security_context
  labels:
    <<: *common_labels
    app.kubernetes.io/component: frontend

backend:
  replicaCount: 2
  image:
    repository: myapp/backend
    tag: "1.0.0"
  resources: *common_resources
  securityContext: *common_security_context
  labels:
    <<: *common_labels
    app.kubernetes.io/component: backend

worker:
  replicaCount: 1
  image:
    repository: myapp/worker
    tag: "1.0.0"
  resources:
    <<: *common_resources
    limits:
      cpu: 1000m  # 覆盖默认值
      memory: 1Gi
  securityContext: *common_security_context
  labels:
    <<: *common_labels
    app.kubernetes.io/component: worker
```

**来源**: [Streamlining Helm Values Files with YAML Anchors](https://dev.to/pczavre/streamlining-helm-values-files-with-yaml-anchors-bpp)

### 7.4 命名规范

| 规范 | 说明 | 示例 |
|------|------|------|
| **前缀约定** | 使用 `_` 或 `x-` 前缀标识可复用块 | `_common_config`, `x-defaults` |
| **描述性命名** | 使用清晰的名称说明用途 | `&database_config`, `&logging_settings` |
| **分组组织** | 将相关锚点放在一起 | 所有 `x-*` 定义放在文件顶部 |
| **避免嵌套** | 不要过度嵌套锚点引用 | 最多 2-3 层嵌套 |

### 7.5 最佳实践总结

| 实践 | 说明 |
|------|------|
| **减少重复** | 定义一次，多处使用 |
| **提高一致性** | 确保所有服务使用相同配置 |
| **易于维护** | 修改一处，全局生效 |
| **合并覆盖** | 使用 `<<:` 合并并覆盖特定字段 |
| **文档化** | 在锚点定义处添加注释说明用途 |

---

## 8. 参数默认值管理

### 8.1 默认值定义位置

不同系统采用不同的默认值管理策略：

| 系统 | 默认值位置 | 优先级规则 |
|------|-----------|-----------|
| **GitHub Actions** | `action.yml` 中的 `default` 字段 | 用户输入 > default |
| **Terraform** | `variable` 块中的 `default` 字段 | CLI > tfvars > default |
| **Helm** | `values.yaml` + `values.schema.json` | --set > -f values.yaml > chart defaults |
| **Kubernetes** | CRD 中的 `+kubebuilder:default` marker | 用户指定 > default |

### 8.2 多层级默认值示例

**Terraform 优先级**：
```hcl
# 1. 变量定义中的默认值（最低优先级）
variable "instance_type" {
  type    = string
  default = "t2.micro"
}

# 2. terraform.tfvars 文件
instance_type = "t2.small"

# 3. 环境变量
# export TF_VAR_instance_type="t2.medium"

# 4. 命令行参数（最高优先级）
# terraform apply -var="instance_type=t2.large"
```

**Helm 优先级**：
```bash
# 1. Chart 默认值（最低优先级）
# values.yaml: replicaCount: 1

# 2. 用户自定义 values 文件
helm install myapp ./chart -f custom-values.yaml

# 3. 命令行 --set（最高优先级）
helm install myapp ./chart --set replicaCount=3
```

### 8.3 默认值验证

**Terraform 示例**：
```hcl
variable "environment" {
  type        = string
  default     = "development"
  description = "Deployment environment"
  
  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be development, staging, or production."
  }
}
```

**JSON Schema 示例**：
```json
{
  "properties": {
    "logLevel": {
      "type": "string",
      "default": "info",
      "enum": ["debug", "info", "warn", "error"],
      "description": "Logging level"
    }
  }
}
```

### 8.4 最佳实践总结

| 实践 | 说明 |
|------|------|
| **明确默认值** | 所有可选参数都应有明确的默认值 |
| **文档化** | 在文档中说明默认值及其含义 |
| **合理默认** | 默认值应适用于大多数场景 |
| **验证默认值** | 确保默认值通过验证规则 |
| **优先级清晰** | 明确说明多个来源的优先级 |
| **避免硬编码** | 不要在代码中硬编码默认值 |

---

## 9. Go Template 最佳实践

### 9.1 参数验证与默认值

```go
package main

import (
	"bytes"
	"fmt"
	"text/template"
)

// 定义配置结构
type CommandConfig struct {
	Tool    string
	Target  string
	Timeout int
	Verbose bool
	Options map[string]string
}

// 验证配置
func (c *CommandConfig) Validate() error {
	if c.Tool == "" {
		return fmt.Errorf("tool is required")
	}
	if c.Target == "" {
		return fmt.Errorf("target is required")
	}
	if c.Timeout <= 0 {
		c.Timeout = 3600 // 默认值
	}
	return nil
}

// 模板函数
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		// 提供默认值
		"default": func(defaultVal, val interface{}) interface{} {
			if val == nil || val == "" {
				return defaultVal
			}
			return val
		},
		// 条件渲染
		"ifset": func(val interface{}) bool {
			return val != nil && val != ""
		},
		// 引号转义
		"quote": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
	}
}

func main() {
	// 定义模板
	tmplStr := `{{.Tool}} -target {{.Target}} \
  {{- if .Verbose}} -v{{end}} \
  {{- if gt .Timeout 0}} -timeout {{.Timeout}}{{end}} \
  {{- range $key, $val := .Options}} \
  -{{$key}} {{quote $val}} \
  {{- end}}`

	// 解析模板
	tmpl, err := template.New("command").
		Funcs(templateFuncs()).
		Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	// 配置
	config := &CommandConfig{
		Tool:    "subfinder",
		Target:  "example.com",
		Timeout: 1800,
		Verbose: true,
		Options: map[string]string{
			"config": "/etc/config.yaml",
			"output": "/tmp/results.txt",
		},
	}

	// 验证
	if err := config.Validate(); err != nil {
		panic(err)
	}

	// 执行模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
	// 输出: subfinder -target example.com -v -timeout 1800 -config "/etc/config.yaml" -output "/tmp/results.txt"
}
```


### 9.2 错误处理

```go
package main

import (
	"bytes"
	"fmt"
	"text/template"
)

// 自定义错误类型
type TemplateError struct {
	Template string
	Field    string
	Err      error
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("template %s: field %s: %v", e.Template, e.Field, e.Err)
}

// 安全执行模板
func ExecuteTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	
	// 设置错误处理选项
	tmpl = tmpl.Option("missingkey=error")
	
	// 执行模板
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", &TemplateError{
			Template: tmpl.Name(),
			Err:      err,
		}
	}
	
	return buf.String(), nil
}

// 验证模板
func ValidateTemplate(tmplStr string, sampleData interface{}) error {
	tmpl, err := template.New("validation").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	
	_, err = ExecuteTemplate(tmpl, sampleData)
	return err
}

func main() {
	tmplStr := `{{.Tool}} -target {{.Target}} -timeout {{.Timeout}}`
	
	// 测试数据
	validData := map[string]interface{}{
		"Tool":    "subfinder",
		"Target":  "example.com",
		"Timeout": 3600,
	}
	
	invalidData := map[string]interface{}{
		"Tool":   "subfinder",
		"Target": "example.com",
		// 缺少 Timeout
	}
	
	// 验证有效数据
	if err := ValidateTemplate(tmplStr, validData); err != nil {
		fmt.Printf("Valid data error: %v\n", err)
	} else {
		fmt.Println("Valid data: OK")
	}
	
	// 验证无效数据
	if err := ValidateTemplate(tmplStr, invalidData); err != nil {
		fmt.Printf("Invalid data error: %v\n", err)
	}
}
```

### 9.3 类型转换

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

func conversionFuncs() template.FuncMap {
	return template.FuncMap{
		// 字符串转整数
		"atoi": func(s string) (int, error) {
			return strconv.Atoi(s)
		},
		// 整数转字符串
		"itoa": func(i int) string {
			return strconv.Itoa(i)
		},
		// 布尔值转字符串
		"boolStr": func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		},
		// 字符串列表转逗号分隔
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		// 字符串转小写
		"lower": strings.ToLower,
		// 字符串转大写
		"upper": strings.ToUpper,
	}
}
```

### 9.4 最佳实践总结

| 实践 | 说明 |
|------|------|
| **提前验证** | 在执行模板前验证数据完整性 |
| **明确错误** | 提供清晰的错误信息，包含字段名和原因 |
| **类型安全** | 使用结构体而非 map，利用类型检查 |
| **函数库** | 提供常用的模板函数（default, quote, join 等）|
| **错误选项** | 使用 `missingkey=error` 捕获缺失字段 |
| **单元测试** | 为模板编写单元测试 |
| **文档化** | 说明模板期望的数据结构 |

---

## 10. 对比分析

### 10.1 当前实现 vs 业界标准

| 维度 | 当前实现 | GitHub Actions | Terraform | Helm | Kubernetes | 推荐方案 |
|------|---------|---------------|-----------|------|-----------|---------|
| **参数定义** | 嵌套 YAML | 扁平 YAML | HCL 结构体 | JSON Schema | Go Struct + Markers | **Go Struct + Tags** |
| **默认值** | 硬编码在代码中 | `default` 字段 | `default` 字段 | `default` 字段 | `+kubebuilder:default` | **Struct 字段或 Schema** |
| **验证规则** | 运行时检查 | 类型 + required | `validation` 块 | JSON Schema | Marker 注释 | **JSON Schema + 运行时验证** |
| **文档生成** | 手动维护 | 从 action.yml | terraform-plugin-docs | 从 schema | 从 markers | **从 Schema 自动生成** |
| **Schema 生成** | 无 | N/A | 从代码提取 | 手动或工具 | controller-gen | **invopop/jsonschema** |
| **错误信息** | 通用错误 | 描述性 | `error_message` | Schema 错误 | 验证错误 | **清晰的字段级错误** |
| **废弃标记** | 无 | `deprecationMessage` | 文档说明 | 文档说明 | 注释 | **deprecationMessage 字段** |
| **类型安全** | 弱类型（map） | 字符串类型 | 强类型 | JSON 类型 | Go 类型 | **Go 强类型** |

### 10.2 参数结构对比

**当前实现（嵌套结构）**：
```yaml
subdomain_discovery:
  passive_tools:
    subfinder:
      enabled: true
      timeout: 3600
      config:
        api_keys:
          shodan: "xxx"
```

**推荐方案（扁平结构）**：
```yaml
# 方案 1: 完全扁平
subfinder_enabled: true
subfinder_timeout: 3600
subfinder_shodan_api_key: "xxx"

# 方案 2: 工具级分组（推荐）
subfinder:
  enabled: true
  timeout: 3600
  shodan_api_key: "xxx"
  censys_api_id: "xxx"
  censys_api_secret: "xxx"
```

**优势对比**：

| 特性 | 嵌套结构 | 扁平结构 |
|------|---------|---------|
| 可读性 | ⭐⭐ | ⭐⭐⭐⭐ |
| 维护性 | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| 模板复杂度 | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| 验证难度 | ⭐⭐ | ⭐⭐⭐⭐ |
| 文档生成 | ⭐⭐ | ⭐⭐⭐⭐⭐ |


### 10.3 默认值管理对比

**当前实现**：
```go
// 硬编码在代码中
timeout := config.Get("timeout", 3600)  // 默认值分散在各处
```

**推荐方案**：
```go
// 方案 1: 在 Struct 中定义
type SubfinderConfig struct {
    Enabled bool   `json:"enabled" jsonschema:"default=true"`
    Timeout int    `json:"timeout" jsonschema:"default=3600,minimum=1,maximum=86400"`
    APIKey  string `json:"apiKey" jsonschema:"description=Shodan API key"`
}

// 方案 2: 在 Schema 中定义
{
  "properties": {
    "timeout": {
      "type": "integer",
      "default": 3600,
      "minimum": 1,
      "maximum": 86400
    }
  }
}
```

### 10.4 验证机制对比

**当前实现**：
```go
// 运行时检查，错误信息不明确
if timeout < 0 {
    return errors.New("invalid timeout")
}
```

**推荐方案**：
```go
// 方案 1: JSON Schema 验证
schema := generateSchema(&SubfinderConfig{})
if err := schema.Validate(config); err != nil {
    // 返回详细的字段级错误
    // "timeout: must be >= 1 and <= 86400"
}

// 方案 2: Struct 验证
type SubfinderConfig struct {
    Timeout int `json:"timeout" validate:"required,min=1,max=86400"`
}

validate := validator.New()
if err := validate.Struct(config); err != nil {
    // 返回详细的验证错误
}
```

### 10.5 文档生成对比

**当前实现**：
- 手动维护 Markdown 文档
- 代码与文档容易不同步
- 更新成本高

**推荐方案**：
```go
// 1. 定义带文档的 Struct
type SubfinderConfig struct {
    Enabled bool   `json:"enabled" jsonschema:"default=true,description=Enable subfinder tool"`
    Timeout int    `json:"timeout" jsonschema:"default=3600,minimum=1,maximum=86400,description=Timeout in seconds"`
    APIKey  string `json:"apiKey" jsonschema:"description=Shodan API key for enhanced results"`
}

// 2. 生成 JSON Schema
schema := jsonschema.Reflect(&SubfinderConfig{})

// 3. 从 Schema 生成 Markdown 文档
doc := generateMarkdownFromSchema(schema)
```

**生成的文档示例**：
```markdown
## Subfinder Configuration

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `enabled` | boolean | No | `true` | Enable subfinder tool |
| `timeout` | integer | No | `3600` | Timeout in seconds (1-86400) |
| `apiKey` | string | No | - | Shodan API key for enhanced results |
```


### 10.6 错误处理对比

**当前实现**：
```
Error: invalid configuration
```

**推荐方案**：
```
Error: Configuration validation failed:
  - subfinder.timeout: must be between 1 and 86400 (got: 0)
  - subfinder.apiKey: required when enabled=true
  - amass.config_file: file does not exist: /path/to/config.yaml
```

### 10.7 综合推荐方案

基于业界标准分析，推荐采用以下架构：

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Go Struct 定义（单一数据源）                               │
│    - 使用 struct tags 定义验证规则                            │
│    - 使用 jsonschema tags 定义 schema 属性                   │
│    - 类型安全，编译时检查                                     │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. 自动生成 JSON Schema                                      │
│    - 使用 invopop/jsonschema 生成                            │
│    - 包含类型、验证规则、默认值、描述                          │
│    - 可用于前端验证和文档生成                                  │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. 运行时验证                                                │
│    - 使用 JSON Schema 验证用户输入                            │
│    - 提供清晰的字段级错误信息                                  │
│    - 在模板执行前验证                                         │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. 模板执行                                                  │
│    - 使用验证后的结构体                                       │
│    - 类型安全的字段访问                                       │
│    - 清晰的错误处理                                          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. 文档生成                                                  │
│    - 从 JSON Schema 自动生成 Markdown                        │
│    - 包含参数说明、类型、默认值、约束                          │
│    - 集成到 CI/CD 自动更新                                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 11. 实施建议

### 11.1 短期改进（1-2 周）

1. **扁平化参数结构**
   - 将嵌套的 `passive_tools.subfinder.config` 改为 `subfinder.*`
   - 减少模板复杂度

2. **添加参数验证**
   - 在模板执行前验证必填参数
   - 提供清晰的错误信息

3. **文档化默认值**
   - 在代码注释中说明所有默认值
   - 更新用户文档

### 11.2 中期改进（1-2 月）

1. **引入 JSON Schema**
   - 为每个工具定义 JSON Schema
   - 实现 schema 验证

2. **重构配置结构**
   - 使用 Go Struct 定义配置
   - 使用 struct tags 定义验证规则

3. **改进错误处理**
   - 实现字段级错误信息
   - 添加错误上下文

### 11.3 长期改进（3-6 月）

1. **自动化 Schema 生成**
   - 从 Go Struct 生成 JSON Schema
   - 集成到构建流程

2. **自动化文档生成**
   - 从 Schema 生成 Markdown 文档
   - 集成到 CI/CD

3. **废弃机制**
   - 添加 `deprecationMessage` 支持
   - 实现平滑迁移路径

---

## 12. 参考资源

### 12.1 官方文档

- [GitHub Actions Metadata Syntax](https://docs.github.com/en/actions/sharing-automations/creating-actions/metadata-syntax-for-github-actions)
- [Terraform Variable Validation](https://developer.hashicorp.com/terraform/language/values/variables#custom-validation-rules)
- [Helm Values Schema](https://helm.sh/docs/topics/charts/#schema-files)
- [Kubebuilder CRD Validation](https://book.kubebuilder.io/reference/markers/crd-validation)
- [JSON Schema Specification](https://json-schema.org/)

### 12.2 工具和库

- [invopop/jsonschema](https://github.com/invopop/jsonschema) - Go struct to JSON Schema
- [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) - Terraform 文档生成
- [controller-gen](https://book.kubebuilder.io/reference/controller-gen) - Kubernetes CRD 生成
- [go-playground/validator](https://github.com/go-playground/validator) - Go struct 验证

### 12.3 真实项目示例

- [actions/checkout](https://github.com/actions/checkout) - GitHub Actions 示例
- [terraform-provider-aws](https://github.com/hashicorp/terraform-provider-aws) - Terraform Provider 示例
- [bitnami/charts](https://github.com/bitnami/charts) - Helm Charts 示例
- [cert-manager](https://github.com/cert-manager/cert-manager) - Kubernetes Operator 示例

---

## 附录：完整示例代码

### A.1 推荐的配置结构

```go
package config

import "github.com/invopop/jsonschema"

// SubfinderConfig defines configuration for subfinder tool
type SubfinderConfig struct {
	// Enable or disable the tool
	Enabled bool `json:"enabled" jsonschema:"default=true,description=Enable subfinder tool"`
	
	// Timeout in seconds
	Timeout int `json:"timeout" jsonschema:"default=3600,minimum=1,maximum=86400,description=Timeout in seconds"`
	
	// Number of threads
	Threads int `json:"threads" jsonschema:"default=10,minimum=1,maximum=100,description=Number of concurrent threads"`
	
	// API Keys
	ShodanAPIKey  string `json:"shodanApiKey,omitempty" jsonschema:"description=Shodan API key"`
	CensysAPIID   string `json:"censysApiId,omitempty" jsonschema:"description=Censys API ID"`
	CensysSecret  string `json:"censysSecret,omitempty" jsonschema:"description=Censys API Secret"`
	
	// Output options
	OutputFile string `json:"outputFile,omitempty" jsonschema:"description=Output file path"`
	Verbose    bool   `json:"verbose" jsonschema:"default=false,description=Enable verbose output"`
}

// Validate performs custom validation
func (c *SubfinderConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	return nil
}

// GenerateSchema generates JSON Schema for SubfinderConfig
func GenerateSchema() *jsonschema.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	return reflector.Reflect(&SubfinderConfig{})
}
```

### A.2 模板执行示例

```go
package executor

import (
	"bytes"
	"fmt"
	"text/template"
)

// ExecuteCommandTemplate executes a command template with validation
func ExecuteCommandTemplate(tmplStr string, config interface{}) (string, error) {
	// 1. 验证配置
	if validator, ok := config.(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return "", fmt.Errorf("config validation failed: %w", err)
		}
	}
	
	// 2. 解析模板
	tmpl, err := template.New("command").
		Funcs(templateFuncs()).
		Option("missingkey=error").
		Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}
	
	// 3. 执行模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}
	
	return buf.String(), nil
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"quote": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
		"default": func(def, val interface{}) interface{} {
			if val == nil || val == "" {
				return def
			}
			return val
		},
	}
}
```

---

**文档版本**: 1.0  
**最后更新**: 2025-01-05  
**维护者**: Worker 命令模板重构项目组
