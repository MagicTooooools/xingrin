# Requirements Document

## Introduction

本文档定义了 Worker 命令模板系统重构的需求。当前系统存在默认值硬编码、配置重复、缺少验证等问题，导致系统不够生产级、易出错、难维护。重构目标是建立一个清晰、可靠、易维护的命令模板系统。

## Glossary

- **CommandTemplate**: 命令模板结构体，定义工具的基础命令和可选参数
- **CommandBuilder**: 命令构建器，根据模板和配置生成最终的命令字符串
- **TemplateLoader**: 模板加载器，从 YAML 文件加载和缓存命令模板
- **Worker**: 执行扫描任务的临时容器
- **YAML 锚点 (Anchor)**: YAML 语法特性，用于定义可重用的配置片段（`&anchor`）
- **YAML 合并键 (Merge Key)**: YAML 语法特性，用于引用和合并锚点定义的配置（`<<: *anchor`）
- **默认值 (Default Value)**: 参数未在用户配置中指定时使用的预设值
- **参数覆盖 (Parameter Override)**: 用户配置值覆盖默认值的机制
- **Schema 验证**: 在加载时检查配置结构和类型的正确性

## Requirements

### Requirement 1: 默认值管理

**User Story:** 作为开发者，我希望默认参数值在模板中明确定义，而不是硬编码在命令字符串中，这样我可以轻松查看和修改默认值。

#### Acceptance Criteria

1. WHEN 模板定义包含参数时，THEN 系统 SHALL 支持为每个参数指定默认值
2. WHEN 用户配置未提供参数值时，THEN 系统 SHALL 使用模板中定义的默认值
3. WHEN 用户配置提供了参数值时，THEN 系统 SHALL 使用用户配置值覆盖默认值
4. THE 系统 SHALL 在模板文件中集中管理所有默认值，不在命令字符串中硬编码
5. WHEN 查看模板文件时，THEN 开发者 SHALL 能够清晰看到每个参数的默认值

### Requirement 2: 配置复用

**User Story:** 作为开发者，我希望使用 YAML 锚点消除重复配置，这样修改共享配置时只需改一处。

#### Acceptance Criteria

1. THE 系统 SHALL 支持在模板文件中使用 YAML 锚点定义共享配置
2. THE 系统 SHALL 支持使用 YAML 合并键引用共享配置
3. WHEN 多个工具共享相同参数时，THEN 系统 SHALL 允许通过锚点定义一次，多处引用
4. WHEN 修改锚点定义的配置时，THEN 所有引用该锚点的工具 SHALL 自动继承修改
5. THE 模板文件 SHALL 包含清晰的注释说明锚点的用途和使用方式

### Requirement 3: 参数类型定义

**User Story:** 作为开发者，我希望参数类型明确定义，这样可以避免类型错误。

#### Acceptance Criteria

1. THE 系统 SHALL 为每个参数定义明确的类型（string, int, bool）
2. WHEN 加载模板时，THEN 系统 SHALL 验证参数类型定义的正确性
3. WHEN 构建命令时，THEN 系统 SHALL 验证参数值与定义的类型匹配
4. IF 参数类型不匹配，THEN 系统 SHALL 返回清晰的错误信息，指明参数名和期望类型
5. THE 系统 SHALL 支持类型转换（如 int 转 string）以适配命令行参数

### Requirement 4: 模板验证

**User Story:** 作为开发者，我希望模板加载时自动验证，这样可以在启动时就发现配置错误，而不是运行时。

#### Acceptance Criteria

1. WHEN 系统启动时，THEN 系统 SHALL 加载并验证所有命令模板
2. WHEN 模板包含语法错误时，THEN 系统 SHALL 拒绝启动并返回详细错误信息
3. WHEN 模板引用未定义的参数时，THEN 系统 SHALL 在加载时报错
4. WHEN 模板的默认值类型与参数类型不匹配时，THEN 系统 SHALL 在加载时报错
5. WHEN 模板验证失败时，THEN 错误信息 SHALL 包含文件名、工具名、参数名和具体错误原因
6. THE 系统 SHALL 在所有模板验证通过后才允许接收扫描任务

### Requirement 5: 错误信息改进

**User Story:** 作为开发者，我希望得到清晰的错误信息，这样可以快速定位问题。

#### Acceptance Criteria

1. WHEN 命令构建失败时，THEN 错误信息 SHALL 包含工具名称
2. WHEN 参数缺失时，THEN 错误信息 SHALL 明确指出缺失的参数名
3. WHEN 参数类型错误时，THEN 错误信息 SHALL 显示期望类型和实际类型
4. WHEN 模板未找到时，THEN 错误信息 SHALL 列出所有可用的模板名称
5. WHEN 占位符未替换时，THEN 错误信息 SHALL 显示未替换的占位符列表

### Requirement 6: 配置文档化

**User Story:** 作为运维人员，我希望配置文件有清晰的注释和示例，这样我可以理解每个参数的作用。

#### Acceptance Criteria

1. THE 模板文件 SHALL 包含文件级注释，说明整体结构和使用方式
2. THE 模板文件 SHALL 为每个工具添加注释，说明工具用途
3. THE 模板文件 SHALL 为共享配置锚点添加注释，说明适用范围
4. THE 模板文件 SHALL 为每个参数添加注释，说明参数含义、单位和取值范围
5. THE 服务端配置文件 SHALL 包含示例，展示如何覆盖默认值

### Requirement 7: 简化配置

**User Story:** 作为运维人员，我希望只需配置 `enabled: true` 就能使用默认值启动工具，这样可以快速开始使用。

#### Acceptance Criteria

1. WHEN 工具配置只包含 `enabled: true` 时，THEN 系统 SHALL 使用所有参数的默认值
2. WHEN 工具配置省略某个参数时，THEN 系统 SHALL 使用该参数的默认值
3. THE 系统 SHALL 允许用户选择性覆盖部分参数，其余参数使用默认值
4. THE 默认值 SHALL 适用于大多数常见场景，无需额外配置即可正常工作
5. THE 文档 SHALL 明确说明哪些参数有默认值，哪些参数必须配置

### Requirement 8: 性能保持

**User Story:** 作为系统架构师，我希望重构不影响性能，这样系统响应速度保持不变。

#### Acceptance Criteria

1. THE 系统 SHALL 继续使用 `sync.Once` 缓存已加载的模板
2. WHEN 多次构建命令时，THEN 系统 SHALL 复用缓存的模板，不重复加载
3. THE 模板验证 SHALL 只在启动时执行一次，不在每次命令构建时执行
4. THE 命令构建时间 SHALL 不超过重构前的 110%
5. THE 内存占用 SHALL 不超过重构前的 120%

### Requirement 9: 参数覆盖优先级

**User Story:** 作为开发者，我希望参数覆盖优先级清晰明确，这样可以预测最终使用的参数值。

#### Acceptance Criteria

1. THE 系统 SHALL 按以下优先级应用参数值：用户配置 > 模板默认值
2. WHEN 用户配置和模板默认值都存在时，THEN 系统 SHALL 使用用户配置值
3. WHEN 用户配置不存在但模板默认值存在时，THEN 系统 SHALL 使用模板默认值
4. WHEN 用户配置和模板默认值都不存在时，THEN 系统 SHALL 不添加该参数到命令
5. THE 文档 SHALL 明确说明参数覆盖的优先级规则

### Requirement 10: 模板结构扩展

**User Story:** 作为开发者，我希望模板结构采用业界标准的嵌套格式，这样可以集中管理每个参数的所有属性。

#### Acceptance Criteria

1. THE CommandTemplate 结构体 SHALL 使用 `parameters` 字段（map），每个参数包含所有属性
2. THE Parameter 结构体 SHALL 包含 `flag`, `default`, `type`, `required`, `description` 字段
3. WHEN 解析 YAML 模板时，THEN 系统 SHALL 正确加载嵌套的参数定义
4. THE `default` 字段 SHALL 支持 string、int、bool 类型的值
5. THE `type` 字段 SHALL 使用字符串表示类型（"string", "int", "bool"）
6. THE 结构 SHALL 符合 GitHub Actions 和 Terraform 的命名规范（使用单数形式）

### Requirement 11: 命令构建逻辑增强

**User Story:** 作为开发者，我希望命令构建器支持从嵌套参数定义中获取值，这样可以自动处理参数覆盖逻辑。

#### Acceptance Criteria

1. WHEN 构建命令时，THEN CommandBuilder SHALL 遍历所有参数定义
2. WHEN 用户配置存在时，THEN CommandBuilder SHALL 用用户配置覆盖参数默认值
3. WHEN 参数值需要类型转换时，THEN CommandBuilder SHALL 自动执行转换
4. WHEN 参数类型不匹配且无法转换时，THEN CommandBuilder SHALL 返回类型错误
5. WHEN 参数标记为 Required 但未提供值时，THEN CommandBuilder SHALL 返回缺失参数错误
6. THE CommandBuilder SHALL 在构建前验证所有必需的占位符都已提供

### Requirement 12: 解析器支持 YAML 锚点

**User Story:** 作为开发者，我希望 YAML 解析器正确处理锚点和合并键，这样可以使用 YAML 的高级特性。

#### Acceptance Criteria

1. THE 系统 SHALL 使用支持 YAML 1.2 规范的解析器
2. WHEN 模板包含锚点定义时，THEN 解析器 SHALL 正确识别和存储锚点
3. WHEN 模板使用合并键引用锚点时，THEN 解析器 SHALL 正确合并配置
4. WHEN 锚点和本地配置冲突时，THEN 解析器 SHALL 优先使用本地配置
5. IF 引用的锚点不存在，THEN 解析器 SHALL 返回明确的错误信息

### Requirement 13: 日志记录增强

**User Story:** 作为运维人员，我希望系统记录详细的日志，这样可以追踪参数来源和命令构建过程。

#### Acceptance Criteria

1. WHEN 加载模板时，THEN 系统 SHALL 记录加载的模板数量和文件路径
2. WHEN 应用默认值时，THEN 系统 SHALL 记录使用默认值的参数名和值
3. WHEN 用户配置覆盖默认值时，THEN 系统 SHALL 记录被覆盖的参数名和新旧值
4. WHEN 构建命令成功时，THEN 系统 SHALL 记录最终的命令字符串
5. WHEN 构建命令失败时，THEN 系统 SHALL 记录详细的错误堆栈和上下文信息

### Requirement 14: 测试覆盖

**User Story:** 作为开发者，我希望重构后的代码有完整的测试覆盖，这样可以确保功能正确性。

#### Acceptance Criteria

1. THE 系统 SHALL 为 CommandTemplate 结构体编写单元测试
2. THE 系统 SHALL 为 CommandBuilder 编写单元测试，覆盖默认值合并逻辑
3. THE 系统 SHALL 为 TemplateLoader 编写单元测试，覆盖验证逻辑
4. THE 系统 SHALL 编写集成测试，验证完整的命令构建流程
5. THE 测试 SHALL 覆盖边界情况：空配置、类型错误、缺失参数等
