# Requirements Document: 业界标准研究与示例收集

## Introduction

本 spec 的目标是研究业界领先项目（GitHub Actions, Terraform, Helm, Kubernetes）如何实现命令模板、配置管理、Schema 验证等功能，并收集具体的代码示例和最佳实践，用于指导我们的重构工作。

## Glossary

- **GitHub Actions**: GitHub 的 CI/CD 平台，使用 YAML 定义工作流
- **Terraform**: HashiCorp 的基础设施即代码工具
- **Helm**: Kubernetes 的包管理器
- **JSON Schema**: 用于验证 JSON/YAML 数据结构的标准
- **Schema Generation**: 从代码或模板自动生成验证 Schema
- **Documentation Generation**: 从代码或 Schema 自动生成文档

## Requirements

### Requirement 1: GitHub Actions 实现研究

**User Story:** 作为开发者，我希望了解 GitHub Actions 如何定义 action 的输入参数，这样可以学习业界标准的参数定义方式。

#### Acceptance Criteria

1. THE 研究 SHALL 找到 GitHub Actions 的 action.yml 文件格式规范
2. THE 研究 SHALL 收集至少 3 个真实 GitHub Actions 项目的 action.yml 示例
3. THE 研究 SHALL 分析参数定义的结构（inputs, description, required, default）
4. THE 研究 SHALL 记录参数命名规范（kebab-case, camelCase）
5. THE 研究 SHALL 记录 deprecationMessage 的使用方式

### Requirement 2: Terraform Provider 实现研究

**User Story:** 作为开发者，我希望了解 Terraform Provider 如何定义变量和验证规则，这样可以学习参数类型系统和验证机制。

#### Acceptance Criteria

1. THE 研究 SHALL 找到 Terraform variable 定义的代码示例
2. THE 研究 SHALL 收集 validation block 的使用示例
3. THE 研究 SHALL 分析 sensitive 参数的处理方式
4. THE 研究 SHALL 记录 terraform-plugin-docs 的文档生成机制
5. THE 研究 SHALL 收集至少 2 个真实 Terraform Provider 的代码示例

### Requirement 3: Helm Chart 实现研究

**User Story:** 作为开发者，我希望了解 Helm 如何管理 values.yaml 和 values.schema.json，这样可以学习配置验证和 Schema 生成。

#### Acceptance Criteria

1. THE 研究 SHALL 找到 Helm values.schema.json 的格式规范
2. THE 研究 SHALL 收集 helm-schema-gen 插件的使用示例
3. THE 研究 SHALL 分析 values.yaml 和 values.schema.json 的关系
4. THE 研究 SHALL 记录 Schema 验证的时机（lint, install, upgrade）
5. THE 研究 SHALL 收集至少 2 个真实 Helm Chart 的 Schema 示例

### Requirement 4: Kubernetes Operator 实现研究

**User Story:** 作为开发者，我希望了解 Kubernetes Operator 如何定义 CRD 和验证规则，这样可以学习声明式配置和验证机制。

#### Acceptance Criteria

1. THE 研究 SHALL 找到 CRD (CustomResourceDefinition) 的定义示例
2. THE 研究 SHALL 分析 OpenAPI v3 Schema 在 CRD 中的使用
3. THE 研究 SHALL 记录 validation rules 的定义方式
4. THE 研究 SHALL 收集 kubebuilder 或 operator-sdk 的代码生成示例
5. THE 研究 SHALL 分析 CRD 的版本管理机制

### Requirement 5: 命令模板模式研究

**User Story:** 作为开发者，我希望了解业界如何实现命令模板和参数替换，这样可以选择最佳的实现方式。

#### Acceptance Criteria

1. THE 研究 SHALL 搜索 "command template pattern" 的实现方式
2. THE 研究 SHALL 比较字符串替换 vs 模板引擎（text/template, Jinja2）
3. THE 研究 SHALL 收集参数占位符的命名规范（{param}, {{param}}, $param）
4. THE 研究 SHALL 分析参数类型转换的处理方式
5. THE 研究 SHALL 记录错误处理和验证的最佳实践

### Requirement 6: 配置文档生成研究

**User Story:** 作为开发者，我希望了解如何自动生成配置文档，这样可以减少手动维护文档的工作量。

#### Acceptance Criteria

1. THE 研究 SHALL 找到 terraform-plugin-docs 的实现原理
2. THE 研究 SHALL 分析文档生成的输入（Schema, 注释, 示例）
3. THE 研究 SHALL 收集文档模板的格式（Markdown, HTML）
4. THE 研究 SHALL 记录文档生成的触发时机（go generate, CI/CD）
5. THE 研究 SHALL 收集至少 2 个文档生成工具的示例

### Requirement 7: JSON Schema 生成研究

**User Story:** 作为开发者，我希望了解如何从代码或配置自动生成 JSON Schema，这样可以实现配置验证。

#### Acceptance Criteria

1. THE 研究 SHALL 找到从 Go struct 生成 JSON Schema 的工具
2. THE 研究 SHALL 找到从 YAML 生成 JSON Schema 的工具（helm-schema-gen）
3. THE 研究 SHALL 分析 JSON Schema 的版本选择（draft-04, draft-07, 2020-12）
4. THE 研究 SHALL 记录 Schema 验证库的选择（Go, Python, JavaScript）
5. THE 研究 SHALL 收集 Schema 生成的代码示例

### Requirement 8: 参数默认值管理研究

**User Story:** 作为开发者，我希望了解业界如何管理参数默认值，这样可以避免硬编码问题。

#### Acceptance Criteria

1. THE 研究 SHALL 分析 GitHub Actions 的 default 字段使用方式
2. THE 研究 SHALL 分析 Terraform 的 default 字段使用方式
3. THE 研究 SHALL 分析 Helm 的 values.yaml 默认值管理
4. THE 研究 SHALL 记录默认值的优先级规则
5. THE 研究 SHALL 收集默认值覆盖的代码示例

### Requirement 9: YAML 锚点和合并键研究

**User Story:** 作为开发者，我希望了解 YAML 锚点在实际项目中的使用方式，这样可以有效消除配置重复。

#### Acceptance Criteria

1. THE 研究 SHALL 找到 Docker Compose 中使用锚点的示例
2. THE 研究 SHALL 找到 Kubernetes 配置中使用锚点的示例
3. THE 研究 SHALL 找到 GitHub Actions 中使用锚点的示例
4. THE 研究 SHALL 记录锚点的命名规范（x-*, _*）
5. THE 研究 SHALL 分析合并键的优先级规则

### Requirement 10: 错误处理和验证研究

**User Story:** 作为开发者，我希望了解业界如何提供清晰的错误信息，这样可以改进我们的错误处理。

#### Acceptance Criteria

1. THE 研究 SHALL 收集 Terraform 的错误信息示例
2. THE 研究 SHALL 收集 Helm 的验证错误信息示例
3. THE 研究 SHALL 分析错误信息的结构（位置, 原因, 建议）
4. THE 研究 SHALL 记录错误信息的国际化处理
5. THE 研究 SHALL 收集错误恢复和建议的最佳实践

### Requirement 11: 真实项目代码收集

**User Story:** 作为开发者，我希望收集真实项目的代码示例，这样可以看到完整的实现。

#### Acceptance Criteria

1. THE 研究 SHALL 找到至少 3 个流行的 GitHub Actions 项目
2. THE 研究 SHALL 找到至少 2 个流行的 Terraform Provider 项目
3. THE 研究 SHALL 找到至少 2 个流行的 Helm Chart 项目
4. THE 研究 SHALL 收集这些项目的关键代码片段
5. THE 研究 SHALL 分析这些项目的目录结构和文件组织

### Requirement 12: 最佳实践总结

**User Story:** 作为开发者，我希望总结业界的最佳实践，这样可以指导我们的重构设计。

#### Acceptance Criteria

1. THE 研究 SHALL 总结参数定义的最佳实践
2. THE 研究 SHALL 总结配置验证的最佳实践
3. THE 研究 SHALL 总结文档生成的最佳实践
4. THE 研究 SHALL 总结错误处理的最佳实践
5. THE 研究 SHALL 创建对比表格，比较不同项目的实现方式

### Requirement 13: 示例代码编写

**User Story:** 作为开发者，我希望基于研究结果编写示例代码，这样可以验证设计的可行性。

#### Acceptance Criteria

1. THE 研究 SHALL 编写 Go 代码示例：从 struct 生成 JSON Schema
2. THE 研究 SHALL 编写 Go 代码示例：使用 JSON Schema 验证 YAML
3. THE 研究 SHALL 编写 Go 代码示例：生成配置文档
4. THE 研究 SHALL 编写 YAML 示例：使用锚点定义共享配置
5. THE 研究 SHALL 编写完整的示例项目结构

### Requirement 14: 文档更新

**User Story:** 作为开发者，我希望将研究结果更新到设计文档中，这样可以指导后续的实现。

#### Acceptance Criteria

1. THE 研究 SHALL 更新 design.md，添加业界示例章节
2. THE 研究 SHALL 更新 design.md，添加代码示例章节
3. THE 研究 SHALL 创建 INDUSTRY_EXAMPLES.md，详细记录所有示例
4. THE 研究 SHALL 创建对比表格，比较当前实现 vs 业界标准
5. THE 研究 SHALL 提供重构建议和优先级

