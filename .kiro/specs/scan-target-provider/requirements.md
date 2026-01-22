# 需求文档

## 简介

本文档定义了扫描管道架构中目标提供者策略模式的实现需求。目标是抽象扫描目标的数据源，在保持向后兼容的同时支持灵活的输入方式。

## 术语表

- **Target（目标）**: 需要扫描的域名、IP 地址或 URL
- **Provider（提供者）**: 从各种来源提供扫描目标的组件
- **TargetProvider（目标提供者）**: 定义目标提供者接口的抽象基类
- **DatabaseTargetProvider（数据库目标提供者）**: 从数据库查询目标的提供者
- **ListTargetProvider（列表目标提供者）**: 使用内存列表的目标提供者
- **FileTargetProvider（文件目标提供者）**: 从文件读取目标的提供者
- **PipelineTargetProvider（管道目标提供者）**: 使用上一阶段输出的提供者
- **BlacklistFilter（黑名单过滤器）**: 过滤黑名单目标的组件
- **ProviderContext（提供者上下文）**: 提供者配置的元数据容器
- **Scan（扫描）**: 对一个或多个目标执行的安全扫描操作
- **Task（任务）**: 执行特定扫描操作的 Prefect 任务
- **Flow（流程）**: 编排多个任务的 Prefect 流程

## 需求

### 需求 1: 抽象目标提供者接口

**用户故事:** 作为开发者，我希望有一个统一的接口来提供扫描目标，以便在不修改扫描逻辑的情况下轻松切换不同的数据源。

#### 验收标准

1. THE TargetProvider SHALL 定义包含迭代主机和 URL 方法的抽象接口
2. THE TargetProvider SHALL 提供获取黑名单过滤器的方法
3. THE TargetProvider SHALL 携带包括 target_id 和 scan_id 的上下文信息
4. THE TargetProvider SHALL 暴露一个属性来指示是否应将结果保存到数据库
5. WHERE 提供了 ProviderContext, THE TargetProvider SHALL 使用它进行配置
6. THE TargetProvider SHALL 提供通用的 CIDR 展开辅助方法供所有子类使用

### 需求 2: 数据库目标提供者

**用户故事:** 作为系统，我希望保持与现有基于数据库的目标查询的向后兼容性，以便当前的扫描功能无需修改即可继续工作。

#### 验收标准

1. WHEN 提供了 target_id, THE DatabaseTargetProvider SHALL 从 Subdomain 表查询主机
2. WHEN 提供了 target_id, THE DatabaseTargetProvider SHALL 从 WebSite 和 Endpoint 表查询 URL
3. THE DatabaseTargetProvider SHALL 检索并应用指定目标的黑名单规则
4. THE DatabaseTargetProvider SHALL 复用现有的 TargetExportService 逻辑
5. THE DatabaseTargetProvider SHALL 在其上下文中设置 target_id

### 需求 3: 列表目标提供者

**用户故事:** 作为用户，我希望对我提供的特定目标执行快速扫描，以便只扫描我感兴趣的 URL 或主机，而不扫描目标下的所有资产。

#### 验收标准

1. WHEN 以列表形式提供目标, THE ListTargetProvider SHALL 自动识别每个目标的类型（URL/域名/IP/CIDR）
2. THE ListTargetProvider SHALL 将 URL 类型的目标通过 iter_urls() 返回
3. THE ListTargetProvider SHALL 将非 URL 类型的目标（域名/IP/CIDR）通过 iter_hosts() 返回
4. THE ListTargetProvider SHALL 默认不应用黑名单过滤器
5. WHEN 未提供目标, THE ListTargetProvider SHALL 返回空迭代器
6. WHEN 目标列表包含 CIDR, THE ListTargetProvider SHALL 自动展开为单个 IP 地址
7. THE ListTargetProvider SHALL 使用 detect_input_type() 进行类型检测
8. WHEN 提供了 ProviderContext, THE ListTargetProvider SHALL 通过 target_id 属性暴露 context.target_id
9. THE ListTargetProvider SHALL 支持通过 context 关联 target_id（用于保存扫描结果）

### 需求 4: 管道目标提供者

**用户故事:** 作为开发者，我希望使用上一个管道阶段的输出作为后续阶段的输入，以便在扫描阶段之间实现高效的数据流。

#### 验收标准

1. WHEN 提供了 StageOutput, THE PipelineTargetProvider SHALL 迭代该输出中的主机
2. WHEN 提供了 StageOutput, THE PipelineTargetProvider SHALL 迭代该输出中的 URL
3. THE PipelineTargetProvider SHALL 不应用黑名单过滤器（数据已被过滤）
4. THE PipelineTargetProvider SHALL 可选地将结果与 target_id 关联
5. WHEN StageOutput 为空, THE PipelineTargetProvider SHALL 返回空迭代器

### 需求 5: 提供者工厂

**用户故事:** 作为开发者，我希望有一个工厂函数根据可用参数创建适当的提供者，以便不需要手动确定要实例化哪个提供者。

#### 验收标准

1. WHEN 提供了 previous_output, THE factory SHALL 创建 PipelineTargetProvider
2. WHEN 提供了 targets 列表, THE factory SHALL 创建 ListTargetProvider
3. WHEN 仅提供了 target_id, THE factory SHALL 创建 DatabaseTargetProvider
4. IF 未提供有效参数, THEN THE factory SHALL 抛出 ValueError
5. THE factory SHALL 按照优先级顺序选择 Provider 类型

### 需求 6: 向后兼容的任务重构

**用户故事:** 作为系统维护者，我希望现有任务支持新的提供者模式同时保持向后兼容性，以便现有代码无需修改即可继续工作。

#### 验收标准

1. WHEN 任务接收到 provider 参数, THE task SHALL 使用该提供者进行数据访问
2. WHEN 任务仅接收到 target_id 参数, THE task SHALL 创建 DatabaseTargetProvider
3. IF 既未提供 provider 也未提供 target_id, THEN THE task SHALL 抛出 ValueError
4. THE task SHALL 在过滤目标时使用提供者的黑名单过滤器
5. THE task SHALL 在保存结果时使用提供者的 target_id 属性

### 需求 7: 提供者上下文管理

**用户故事:** 作为开发者，我希望通过提供者传递元数据，以便扫描操作能够访问必要的上下文信息。

#### 验收标准

1. THE ProviderContext SHALL 包含可选的 target_id 字段
2. THE ProviderContext SHALL 包含可选的 scan_id 字段
3. WHEN target_id 为 None, THE 扫描任务 SHALL 跳过保存结果到数据库
4. WHEN 未提供 ProviderContext, THE TargetProvider SHALL 创建默认上下文
5. THE TargetProvider SHALL 通过其接口暴露上下文属性

### 需求 8: 基于迭代器的数据访问

**用户故事:** 作为开发者，我希望提供者使用迭代器进行数据访问，以便在处理大型数据集时保持内存使用效率。

#### 验收标准

1. THE iter_hosts 方法 SHALL 返回主机字符串的迭代器
2. THE iter_urls 方法 SHALL 返回 URL 字符串的迭代器
3. THE 迭代器 SHALL 支持惰性求值以最小化内存使用
4. THE 迭代器 SHALL 优雅地处理空数据集
5. THE 迭代器 SHALL 不将所有数据一次性加载到内存中

### 需求 9: 黑名单过滤器集成

**用户故事:** 作为用户，我希望在扫描数据库目标时自动应用黑名单规则，以便将黑名单中的主机和 URL 排除在扫描之外。

#### 验收标准

1. WHEN 使用 DatabaseTargetProvider, THE provider SHALL 检索目标的黑名单规则
2. THE get_blacklist_filter 方法 SHALL 返回 BlacklistFilter 实例或 None
3. THE BlacklistFilter SHALL 从通过 BlacklistService 检索的规则创建
4. WHEN 使用 ListTargetProvider, THE provider SHALL 为黑名单过滤器返回 None
5. WHEN 使用 PipelineTargetProvider, THE provider SHALL 返回 None（数据已被过滤）
