# 实现计划: Go 模型补全

## 概述

补全所有 Go 数据模型，确保与 Django 模型完全兼容。

## 任务

- [x] 1. 修复已有模型
  - [x] 1.1 补充 Scan 模型缺失字段
    - 添加 cached_directories_count, cached_screenshots_count
    - 添加 cached_vulns_critical/high/medium/low
    - 添加 stats_updated_at
    - _需求: 1.1_
  - [x] 1.2 补充 WebSite 模型缺失字段
    - 添加 location, response_body, content_type, vhost
    - _需求: 1.2_

- [x] 2. 实现资产模型
  - [x] 2.1 实现 Endpoint 模型
    - _需求: 2.1_
  - [x] 2.2 实现 Directory 模型
    - _需求: 2.2_
  - [x] 2.3 实现 HostPortMapping 模型
    - _需求: 2.3_
  - [x] 2.4 实现 Vulnerability 模型
    - _需求: 2.4_
  - [x] 2.5 实现 Screenshot 模型
    - _需求: 2.5_

- [x] 3. 实现快照模型
  - [x] 3.1 实现 SubdomainSnapshot 模型
    - _需求: 3.1_
  - [x] 3.2 实现 WebsiteSnapshot 模型
    - _需求: 3.2_
  - [x] 3.3 实现 EndpointSnapshot 模型
    - _需求: 3.3_
  - [x] 3.4 实现 DirectorySnapshot 模型
    - _需求: 3.4_
  - [x] 3.5 实现 HostPortMappingSnapshot 模型
    - _需求: 3.5_
  - [x] 3.6 实现 VulnerabilitySnapshot 模型
    - _需求: 3.6_
  - [x] 3.7 实现 ScreenshotSnapshot 模型
    - _需求: 3.7_

- [x] 4. 实现扫描相关模型
  - [x] 4.1 实现 ScanLog 模型
    - _需求: 4.1_
  - [x] 4.2 实现 ScanInputTarget 模型
    - _需求: 4.2_
  - [x] 4.3 实现 ScheduledScan 模型
    - _需求: 4.3_
  - [x] 4.4 实现 SubfinderProviderSettings 模型
    - _需求: 4.4_

- [x] 5. 实现引擎相关模型
  - [x] 5.1 实现 Wordlist 模型
    - _需求: 5.1_
  - [x] 5.2 实现 NucleiTemplateRepo 模型
    - _需求: 5.2_

- [x] 6. 实现通知和配置模型
  - [x] 6.1 实现 Notification 模型
    - _需求: 6.1_
  - [x] 6.2 实现 NotificationSettings 模型
    - _需求: 6.2_
  - [x] 6.3 实现 BlacklistRule 模型
    - _需求: 6.3_

- [x] 7. 实现统计模型
  - [x] 7.1 实现 AssetStatistics 模型
    - _需求: 7.1_
  - [x] 7.2 实现 StatisticsHistory 模型
    - _需求: 7.2_

- [x] 8. 实现认证模型
  - [x] 8.1 实现 User 模型
    - _需求: 8.1_
  - [x] 8.2 实现 Session 模型
    - _需求: 8.2_

- [x] 9. 更新测试
  - [x] 9.1 更新 model_test.go
    - 添加所有新模型的表名测试
    - _需求: 9.1_

- [x] 10. 检查点 - 验证模型完整性
  - 所有测试通过 ✅
  - 模型数量与 Django 一致 ✅ (33 个模型)
