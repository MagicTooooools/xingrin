# 需求文档

## 简介

补全 Go 后端的所有数据模型，确保与 Django 模型完全兼容。包括新增缺失模型和修复已有模型的缺失字段。

## 术语表

- **Model**: Go 结构体，映射到数据库表
- **GORM**: Go 的 ORM 库
- **Snapshot**: 快照表，记录某次扫描的结果

## 模型清单

### 已有模型（需检查字段完整性）

| 模型 | 表名 | 状态 |
|------|------|------|
| Organization | organization | ✅ 完整 |
| Target | target | ✅ 完整 |
| Scan | scan | ⚠️ 缺少字段 |
| Subdomain | subdomain | ✅ 完整 |
| WebSite | website | ⚠️ 缺少字段 |
| WorkerNode | worker_node | ✅ 完整 |
| ScanEngine | scan_engine | ✅ 完整 |

### 缺失模型

| 模型 | 表名 | 分类 |
|------|------|------|
| Endpoint | endpoint | 资产 |
| Directory | directory | 资产 |
| HostPortMapping | host_port_mapping | 资产 |
| Vulnerability | vulnerability | 资产 |
| SubdomainSnapshot | subdomain_snapshot | 快照 |
| WebsiteSnapshot | website_snapshot | 快照 |
| EndpointSnapshot | endpoint_snapshot | 快照 |
| DirectorySnapshot | directory_snapshot | 快照 |
| HostPortMappingSnapshot | host_port_mapping_snapshot | 快照 |
| VulnerabilitySnapshot | vulnerability_snapshot | 快照 |
| ScreenshotSnapshot | screenshot_snapshot | 快照 |
| Screenshot | screenshot | 资产 |
| ScanLog | scan_log | 扫描 |
| ScanInputTarget | scan_input_target | 扫描 |
| ScheduledScan | scheduled_scan | 扫描 |
| SubfinderProviderSettings | subfinder_provider_settings | 配置 |
| Wordlist | wordlist | 引擎 |
| NucleiTemplateRepo | nuclei_template_repo | 引擎 |
| Notification | notification | 通知 |
| NotificationSettings | notification_settings | 通知 |
| BlacklistRule | blacklist_rule | 配置 |
| AssetStatistics | asset_statistics | 统计 |
| StatisticsHistory | statistics_history | 统计 |
| User | auth_user | 认证 |
| Session | django_session | 认证 |

## 需求

### 需求 1: 修复已有模型

**用户故事:** 作为开发者，我希望已有的 Go 模型字段与 Django 完全一致。

#### 验收标准

1. THE Scan 模型 SHALL 添加缺失字段（cached_directories_count, cached_screenshots_count, cached_vulns_critical/high/medium/low, stats_updated_at）
2. THE WebSite 模型 SHALL 添加缺失字段（location, response_body, content_type, vhost）

### 需求 2: 资产模型

**用户故事:** 作为开发者，我希望 Go 后端有完整的资产模型。

#### 验收标准

1. THE Go_Model SHALL 实现 Endpoint 模型
2. THE Go_Model SHALL 实现 Directory 模型
3. THE Go_Model SHALL 实现 HostPortMapping 模型
4. THE Go_Model SHALL 实现 Vulnerability 模型
5. THE Go_Model SHALL 实现 Screenshot 模型

### 需求 3: 快照模型

**用户故事:** 作为开发者，我希望 Go 后端有完整的快照模型。

#### 验收标准

1. THE Go_Model SHALL 实现 SubdomainSnapshot 模型
2. THE Go_Model SHALL 实现 WebsiteSnapshot 模型
3. THE Go_Model SHALL 实现 EndpointSnapshot 模型
4. THE Go_Model SHALL 实现 DirectorySnapshot 模型
5. THE Go_Model SHALL 实现 HostPortMappingSnapshot 模型
6. THE Go_Model SHALL 实现 VulnerabilitySnapshot 模型
7. THE Go_Model SHALL 实现 ScreenshotSnapshot 模型

### 需求 4: 扫描相关模型

**用户故事:** 作为开发者，我希望 Go 后端有完整的扫描相关模型。

#### 验收标准

1. THE Go_Model SHALL 实现 ScanLog 模型
2. THE Go_Model SHALL 实现 ScanInputTarget 模型
3. THE Go_Model SHALL 实现 ScheduledScan 模型
4. THE Go_Model SHALL 实现 SubfinderProviderSettings 模型（单例）

### 需求 5: 引擎相关模型

**用户故事:** 作为开发者，我希望 Go 后端有完整的引擎相关模型。

#### 验收标准

1. THE Go_Model SHALL 实现 Wordlist 模型
2. THE Go_Model SHALL 实现 NucleiTemplateRepo 模型

### 需求 6: 通知和配置模型

**用户故事:** 作为开发者，我希望 Go 后端有完整的通知和配置模型。

#### 验收标准

1. THE Go_Model SHALL 实现 Notification 模型
2. THE Go_Model SHALL 实现 NotificationSettings 模型（单例）
3. THE Go_Model SHALL 实现 BlacklistRule 模型

### 需求 7: 统计模型

**用户故事:** 作为开发者，我希望 Go 后端有完整的统计模型。

#### 验收标准

1. THE Go_Model SHALL 实现 AssetStatistics 模型（单例）
2. THE Go_Model SHALL 实现 StatisticsHistory 模型

### 需求 8: 认证模型

**用户故事:** 作为开发者，我希望 Go 后端有用户模型，以便实现认证。

#### 验收标准

1. THE Go_Model SHALL 实现 User 模型（兼容 Django auth_user 表）
2. THE Go_Model SHALL 实现 Session 模型（兼容 Django django_session 表）

### 需求 9: 模型一致性

**用户故事:** 作为开发者，我希望所有 Go 模型与 Django 模型完全一致。

#### 验收标准

1. THE Go_Model SHALL 使用相同的表名（TableName() 方法）
2. THE Go_Model SHALL 使用相同的字段名（gorm column tag）
3. THE Go_Model SHALL 正确处理 PostgreSQL 数组类型
4. THE Go_Model SHALL 正确处理 JSONB 类型
5. WHEN 序列化为 JSON 时，THE Go_Model SHALL 输出 camelCase 字段名
