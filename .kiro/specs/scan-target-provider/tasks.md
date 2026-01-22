# 实现计划: 扫描目标提供者策略模式

## 概述

实现 TargetProvider 策略模式，让扫描任务可以从不同数据源获取目标（数据库、列表、文件），同时保持向后兼容。

## 任务

- [x] 1. 创建 providers 模块基础结构
  - 创建 `backend/apps/scan/providers/` 目录
  - 创建 `__init__.py` 导出公共接口
  - 创建 `base.py` 定义 ProviderContext 和 TargetProvider 抽象基类
  - 实现 `_expand_host()` 静态方法用于 CIDR 展开
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.6, 8.1, 8.2, 8.3_

- [x] 2. 实现具体 Provider 类
  - [x] 2.1 实现 ListTargetProvider
    - 创建 `list_provider.py`
    - 实现 iter_hosts(), iter_urls(), get_blacklist_filter()
    - 在 iter_hosts() 中使用 _expand_host() 展开 CIDR
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

  - [x] 2.2 编写 ListTargetProvider 属性测试
    - **Property 1: ListTargetProvider Round-Trip**
    - **Validates: Requirements 3.1, 3.2**

  - [x] 2.3 实现 DatabaseTargetProvider
    - 创建 `database_provider.py`
    - 实现从 Subdomain/WebSite/Endpoint 表查询
    - 实现黑名单过滤器集成
    - 复用现有 TargetExportService 逻辑
    - 在 CIDR 类型处理中使用 _expand_host() 展开
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 9.1, 9.2, 9.3_

  - [x] 2.4 编写 DatabaseTargetProvider 属性测试
    - **Property 6: DatabaseTargetProvider Blacklist Application**
    - **Validates: Requirements 2.3, 9.1, 9.2, 9.3_

  - [x] 2.5 实现 PipelineTargetProvider（预留）
    - 创建 `pipeline_provider.py`
    - 实现从 StageOutput 读取数据
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

  - [x] 2.6 编写 PipelineTargetProvider 属性测试
    - **Property 2: PipelineTargetProvider Round-Trip**
    - **Validates: Requirements 4.1, 4.2**

- [x] 3. Checkpoint - 确保所有 Provider 测试通过
  - 运行所有测试，确保 Provider 模块功能正确
  - 如有问题请询问用户

- [x] 4. 改造现有 Task（向后兼容）
  - [x] 4.1 改造 export_hosts_task
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 实现兼容逻辑：无 provider 时用 target_id 创建 DatabaseTargetProvider
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [x] 4.2 编写 export_hosts_task 向后兼容测试
    - **Property 8: Task Backward Compatibility**
    - **Validates: Requirements 6.1, 6.2, 6.4, 6.5**

  - [x] 4.3 改造 export_site_urls_task
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 实现兼容逻辑
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 5. 编写通用属性测试
  - [x] 5.1 编写 Context Propagation 属性测试
    - **Property 4: Context Propagation**
    - **Validates: Requirements 1.3, 1.5, 7.4, 7.5**

  - [x] 5.2 编写 Non-Database Provider Blacklist Filter 属性测试
    - **Property 5: Non-Database Provider Blacklist Filter**
    - **Validates: Requirements 3.4, 9.4, 9.5**

  - [x] 5.3 编写 CIDR Expansion Consistency 属性测试
    - **Property 7: CIDR Expansion Consistency**
    - **Validates: Requirements 1.6, 3.6**

- [x] 6. Final Checkpoint - 确保所有测试通过
  - 运行完整测试套件
  - 验证现有扫描功能不受影响（回归测试）
  - 如有问题请询问用户

- [ ] 7. 改造剩余 Export Task（Phase 2）
  - [x] 7.1 改造 url_fetch/export_sites_task
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 实现兼容逻辑：无 provider 时用 target_id 创建 DatabaseTargetProvider
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [x] 7.2 改造 directory_scan/export_sites_task
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 实现兼容逻辑
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [x] 7.3 改造 vuln_scan/export_endpoints_task
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 实现兼容逻辑
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [x] 7.4 改造 fingerprint_detect/export_urls_task
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 实现兼容逻辑
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 8. Checkpoint - 确保所有 Export Task 改造完成
  - 运行所有测试
  - 验证向后兼容性
  - 如有问题请询问用户

- [ ] 9. 改造 Screenshot Flow（Phase 2 续）
  - [x] 9.1 改造 screenshot_flow 支持 TargetProvider
    - 添加 provider 参数
    - 保留 target_id 参数（向后兼容）
    - 当有 provider 时，使用 provider.iter_urls() 获取 URL
    - 当无 provider 时，使用现有 get_urls_with_fallback()
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

## 备注

- 每个任务都引用了具体的需求以便追溯
- Checkpoint 任务用于增量验证
- 属性测试验证通用正确性属性
- 单元测试验证具体示例和边界情况
