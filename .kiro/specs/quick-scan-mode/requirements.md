# 需求文档

## 简介

本文档定义了快速扫描模式的实现需求。快速扫描模式允许用户只扫描指定的目标（而非整个 Target 下的所有资产），通过快照表实现阶段间的精确数据传递。

## 术语表

- **快速扫描（Quick Scan）**: 只扫描用户指定的目标，不扫描 Target 下的历史资产
- **完整扫描（Full Scan）**: 扫描 Target 下的所有资产（现有行为）
- **扫描模式（Scan Mode）**: 区分快速扫描和完整扫描的标识
- **快照表（Snapshot Table）**: 存储单次扫描发现的资产，用于阶段间数据传递（已存在）
- **SubdomainSnapshot**: 子域名快照表（已存在）
- **HostPortMappingSnapshot**: 主机端口映射快照表（已存在）
- **WebsiteSnapshot**: 网站快照表（已存在）
- **DirectorySnapshot**: 目录快照表（已存在）
- **EndpointSnapshot**: 端点快照表（已存在）
- **VulnerabilitySnapshot**: 漏洞快照表（已存在）
- **SnapshotTargetProvider**: 从快照表读取数据的 Provider（已存在，需完善）
- **ListTargetProvider**: 从内存列表读取数据的 Provider（已存在）
- **DatabaseTargetProvider**: 从数据库查询数据的 Provider（已存在）
- **scan_id**: 扫描任务唯一标识
- **target_id**: 目标唯一标识
- **用户输入目标（User Input Targets）**: 用户在快速扫描时指定的目标列表

## 需求

### 需求 1: 扫描模式标识

**用户故事:** 作为系统，我希望能够区分快速扫描和完整扫描模式，以便根据不同模式选择正确的数据源。

#### 验收标准

1. THE Scan 模型 SHALL 包含 scan_mode 字段，支持 'full' 和 'quick' 两种值
2. WHEN 创建扫描任务时, THE 系统 SHALL 根据请求参数设置 scan_mode
3. THE scan_mode 字段 SHALL 默认为 'full'（向后兼容）
4. WHEN scan_mode 为 'quick' 时, THE Scan 模型 SHALL 存储用户输入的目标列表

### 需求 2: 快照表数据模型（已存在）

**用户故事:** 作为系统，我希望有快照表来存储单次扫描发现的资产，以便在快速扫描模式下实现阶段间的精确数据传递。

#### 验收标准

1. THE SubdomainSnapshot 表 SHALL 包含 scan_id、name 字段（已存在）
2. THE HostPortMappingSnapshot 表 SHALL 包含 scan_id、host、ip、port 字段（已存在）
3. THE WebsiteSnapshot 表 SHALL 包含 scan_id、url、host 字段（已存在）
4. THE DirectorySnapshot 表 SHALL 包含 scan_id、url 字段（已存在）
5. THE EndpointSnapshot 表 SHALL 包含 scan_id、url、host 字段（已存在）
6. THE VulnerabilitySnapshot 表 SHALL 包含 scan_id、url、vuln_type、severity 字段（已存在）
7. THE 快照表 SHALL 通过 scan_id 建立索引以支持高效查询（已存在）
8. THE 快照表 SHALL 支持级联删除（删除 Scan 时自动删除关联快照）（已存在）

### 需求 3: 快照保存服务

**用户故事:** 作为开发者，我希望有统一的服务来保存快照数据，以便在各个 Flow 中复用。

#### 验收标准

1. WHEN 子域名发现完成时, THE 系统 SHALL 同时保存到 Subdomain 表和 SubdomainSnapshot 表
2. WHEN 端口扫描完成时, THE 系统 SHALL 同时保存到 HostPortMapping 表和 HostPortMappingSnapshot 表
3. WHEN 站点扫描完成时, THE 系统 SHALL 同时保存到 Website 表和 WebsiteSnapshot 表
4. WHEN URL 获取完成时, THE 系统 SHALL 同时保存到 Endpoint 表和 EndpointSnapshot 表
5. THE 快照保存 SHALL 使用批量插入以优化性能
6. IF target_id 为 None, THEN THE 系统 SHALL 跳过保存到主表（只保存快照）

### 需求 4: 快照查询服务

**用户故事:** 作为开发者，我希望有统一的服务来查询快照数据，以便 SnapshotTargetProvider 使用。

#### 验收标准

1. THE SubdomainSnapshotsService SHALL 提供 iter_subdomain_names_by_scan(scan_id) 方法
2. THE HostPortMappingSnapshotsService SHALL 提供 iter_by_scan(scan_id) 方法
3. THE WebsiteSnapshotsService SHALL 提供 iter_by_scan(scan_id) 方法
4. THE EndpointSnapshotsService SHALL 提供 iter_by_scan(scan_id) 方法
5. THE 查询方法 SHALL 支持分块迭代以优化内存使用
6. THE 查询方法 SHALL 返回迭代器而非列表

### 需求 5: SnapshotTargetProvider 完善

**用户故事:** 作为开发者，我希望 SnapshotTargetProvider 能够从快照表读取数据，以便在快速扫描的后续阶段使用。

#### 验收标准

1. WHEN snapshot_type 为 'subdomain' 时, THE SnapshotTargetProvider SHALL 从 SubdomainSnapshot 表读取主机列表
2. WHEN snapshot_type 为 'host_port' 时, THE SnapshotTargetProvider SHALL 从 HostPortMappingSnapshot 表读取主机端口列表
3. WHEN snapshot_type 为 'website' 时, THE SnapshotTargetProvider SHALL 从 WebsiteSnapshot 表读取 URL 列表
4. WHEN snapshot_type 为 'endpoint' 时, THE SnapshotTargetProvider SHALL 从 EndpointSnapshot 表读取 URL 列表
5. THE SnapshotTargetProvider SHALL 不应用黑名单过滤（数据已在上一阶段过滤）
6. WHEN 快照表为空时, THE SnapshotTargetProvider SHALL 返回空迭代器

### 需求 6: initiate_scan_flow 改造

**用户故事:** 作为系统，我希望 initiate_scan_flow 能够根据扫描模式选择正确的 Provider，以便实现精确扫描。

#### 验收标准

1. WHEN scan_mode 为 'quick' 时, THE initiate_scan_flow SHALL 为第一个阶段创建 ListTargetProvider
2. WHEN scan_mode 为 'quick' 时, THE initiate_scan_flow SHALL 为后续阶段创建 SnapshotTargetProvider
3. WHEN scan_mode 为 'full' 时, THE initiate_scan_flow SHALL 为所有阶段创建 DatabaseTargetProvider
4. THE initiate_scan_flow SHALL 将 Provider 传递给所有子 Flow
5. THE initiate_scan_flow SHALL 从 Scan 模型读取用户输入的目标列表

### 需求 7: 各 Flow 保存逻辑改造

**用户故事:** 作为系统，我希望各个 Flow 在保存数据时同时写入快照表，以便支持快速扫描模式。

#### 验收标准

1. WHEN subdomain_discovery_flow 保存子域名时, THE 系统 SHALL 同时写入 SubdomainSnapshot
2. WHEN port_scan_flow 保存端口映射时, THE 系统 SHALL 同时写入 HostPortMappingSnapshot
3. WHEN site_scan_flow 保存网站时, THE 系统 SHALL 同时写入 WebsiteSnapshot
4. WHEN url_fetch_flow 保存端点时, THE 系统 SHALL 同时写入 EndpointSnapshot
5. THE 快照写入 SHALL 使用 scan_id 关联
6. THE 快照写入 SHALL 不影响现有的主表保存逻辑

### 需求 8: API 层改造

**用户故事:** 作为用户，我希望通过 API 发起快速扫描，以便只扫描我指定的目标。

#### 验收标准

1. THE /api/scans/quick/ 接口 SHALL 接收用户输入的目标列表
2. THE /api/scans/quick/ 接口 SHALL 将目标列表存储到 Scan 模型
3. THE /api/scans/quick/ 接口 SHALL 设置 scan_mode 为 'quick'
4. THE /api/scans/initiate/ 接口 SHALL 保持 scan_mode 为 'full'（向后兼容）
5. WHEN 快速扫描完成时, THE 系统 SHALL 返回扫描结果摘要

### 需求 9: 向后兼容性

**用户故事:** 作为系统维护者，我希望快速扫描改造不影响现有的完整扫描功能。

#### 验收标准

1. WHEN scan_mode 未指定时, THE 系统 SHALL 默认使用完整扫描模式
2. THE 现有 API 接口 SHALL 保持向后兼容
3. THE 现有 Flow 和 Task SHALL 在完整扫描模式下行为不变
4. THE 快照表写入 SHALL 对完整扫描模式透明（不影响性能）
5. THE 定时扫描 SHALL 继续使用完整扫描模式

### 需求 10: 快照数据清理

**用户故事:** 作为系统管理员，我希望能够清理过期的快照数据，以便节省存储空间。

#### 验收标准

1. WHEN Scan 记录被删除时, THE 系统 SHALL 级联删除关联的快照数据
2. THE 系统 SHALL 提供清理指定时间之前快照数据的方法
3. THE 快照清理 SHALL 不影响主表数据
4. THE 快照清理 SHALL 支持批量删除以优化性能

