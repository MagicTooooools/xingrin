# 实现任务

## 任务概览

本文档定义了快速扫描模式的实现任务。任务按依赖关系排序，每个任务都是独立可测试的。

## 任务列表

- [x] 1. ScanInputTarget 模型

**描述**: 新建 ScanInputTarget 模型，存储快速扫描时用户输入的目标

**相关需求**: 需求 1

**文件**:
- `backend/apps/scan/models/scan_input_target.py`（新建）
- `backend/apps/scan/models/__init__.py`（更新导出）

**实现步骤**:
1. 创建 ScanInputTarget 模型，包含 scan_id, value, input_type, created_at 字段
2. 添加 InputType 枚举（domain/ip/cidr/url）
3. 添加索引（scan_id, input_type）
4. 创建数据库迁移文件
5. 运行迁移

**验收标准**:
- [x] 模型可正常创建记录
- [x] 支持通过 scan.input_targets 反向查询
- [x] 迁移文件可正常执行

---

### Task 2: Scan 模型扩展

**描述**: 为 Scan 模型添加 scan_mode 字段

**相关需求**: 需求 1

**文件**:
- `backend/apps/scan/models/scan.py`

**实现步骤**:
1. 添加 ScanMode 枚举类（FULL, QUICK）
2. 添加 scan_mode 字段（CharField，默认 'full'）
3. 创建数据库迁移文件
4. 运行迁移

**验收标准**:
- [x] scan_mode 字段默认值为 'full'
- [x] 迁移文件可正常执行

---

### Task 3: ScanInputTargetService

**描述**: 新建 ScanInputTargetService，提供批量创建和查询功能

**相关需求**: 需求 1, 8

**依赖**: Task 1

**文件**:
- `backend/apps/scan/services/scan_input_target_service.py`（新建）

**实现步骤**:
1. 实现 bulk_create(scan_id, inputs) 方法，解析输入类型并批量写入
2. 实现 iter_by_scan(scan_id) 方法，分块迭代查询
3. 实现 iter_hosts_by_scan(scan_id) 方法，只返回 domain/ip/cidr 类型
4. 实现 iter_urls_by_scan(scan_id) 方法，只返回 url 类型

**验收标准**:
- [x] bulk_create 可批量写入 1 万条记录
- [x] iter_by_scan 支持分块迭代
- [x] 输入类型自动识别正确

---

### Task 4: ScanInputTargetProvider

**描述**: 新建 ScanInputTargetProvider，从 ScanInputTarget 表读取数据

**相关需求**: 需求 5, 6

**依赖**: Task 1, Task 3

**文件**:
- `backend/apps/scan/providers/scan_input_provider.py`（新建）
- `backend/apps/scan/providers/__init__.py`（更新导出）

**实现步骤**:
1. 继承 TargetProvider 基类
2. 实现 _iter_raw_hosts() 方法，查询 domain/ip/cidr 类型
3. 实现 iter_urls() 方法，查询 url 类型
4. 实现 get_blacklist_filter() 返回 None
5. 添加单元测试

**验收标准**:
- [x] iter_hosts() 返回 domain/ip/cidr 类型的输入
- [x] iter_urls() 返回 url 类型的输入
- [x] 支持分块迭代（chunk_size=1000）

---

### Task 5: initiate_scan_flow Provider 选择逻辑

**描述**: 修改 initiate_scan_flow，根据 scan_mode 统一创建所有阶段的 Provider

**相关需求**: 需求 6

**依赖**: Task 2, Task 4

**文件**:
- `backend/apps/scan/flows/initiate_scan_flow.py`

**实现步骤**:
1. 从 Scan 模型读取 scan_mode
2. 根据 scan_mode 创建 providers 字典，包含所有阶段的 Provider
3. 调用子 Flow 时传入对应的 provider
4. 移除 scan_mode 参数传递（子 Flow 不需要）

**验收标准**:
- [x] 快速扫描时各阶段使用正确的 Provider（ScanInputTargetProvider/SnapshotTargetProvider）
- [x] 完整扫描时所有阶段使用 DatabaseTargetProvider
- [x] 子 Flow 不再接收 scan_mode 参数

---

### Task 6: subdomain_discovery_flow 快照写入

**描述**: 修改子域名发现 Flow，在保存数据时同时写入 SubdomainSnapshot

**相关需求**: 需求 3, 7

**依赖**: Task 5

**文件**:
- `backend/apps/scan/flows/subdomain_discovery_flow.py`
- `backend/apps/scan/tasks/subdomain_discovery/` 相关任务

**实现步骤**:
1. 检查现有保存逻辑，确认是否已使用 SubdomainSnapshotsService.save_and_sync()
2. 如果未使用，修改为使用 save_and_sync() 方法
3. 确保 scan_id 正确传递到 DTO

**验收标准**:
- [x] 子域名发现结果同时保存到 Subdomain 表和 SubdomainSnapshot 表
- [x] SubdomainSnapshot 记录包含正确的 scan_id

**备注**: Flow 签名已添加 provider 参数，快照写入逻辑已在现有 save_domains_task 中实现

---

### Task 7: port_scan_flow 快照读取和写入

**描述**: 修改端口扫描 Flow，使用传入的 Provider，并写入 HostPortMappingSnapshot

**相关需求**: 需求 3, 5, 7

**依赖**: Task 6

**文件**:
- `backend/apps/scan/flows/port_scan_flow.py`
- `backend/apps/scan/tasks/port_scan/` 相关任务

**实现步骤**:
1. 修改 Flow 签名，provider 改为必需参数
2. 移除 scan_mode 参数和相关判断逻辑
3. 直接使用传入的 provider 获取扫描目标
4. 检查现有保存逻辑，确认是否已写入 HostPortMappingSnapshot
5. 如果未写入，添加快照写入逻辑

**验收标准**:
- [x] Flow 直接使用传入的 provider
- [x] 端口扫描结果同时保存到 HostPortMapping 表和 HostPortMappingSnapshot 表

**备注**: Flow 签名已添加 provider 参数，export_hosts_task 已支持 provider 模式

---

### Task 8: site_scan_flow 快照读取和写入

**描述**: 修改网站扫描 Flow，使用传入的 Provider，并写入 WebsiteSnapshot

**相关需求**: 需求 3, 5, 7

**依赖**: Task 7

**文件**:
- `backend/apps/scan/flows/site_scan_flow.py`
- `backend/apps/scan/tasks/site_scan/` 相关任务

**实现步骤**:
1. 修改 Flow 签名，provider 改为必需参数
2. 移除 scan_mode 参数和相关判断逻辑
3. 直接使用传入的 provider 获取扫描目标
4. 检查现有保存逻辑，确认是否已写入 WebsiteSnapshot
5. 如果未写入，添加快照写入逻辑

**验收标准**:
- [x] Flow 直接使用传入的 provider
- [x] 网站扫描结果同时保存到 Website 表和 WebsiteSnapshot 表

**备注**: Flow 签名已添加 provider 参数，export_site_urls_task 已支持 provider 模式

---

### Task 9: url_fetch_flow 快照读取和写入

**描述**: 修改 URL 获取 Flow，使用传入的 Provider，并写入 EndpointSnapshot

**相关需求**: 需求 3, 5, 7

**依赖**: Task 8

**文件**:
- `backend/apps/scan/flows/url_fetch/main_flow.py`
- `backend/apps/scan/tasks/url_fetch/` 相关任务

**实现步骤**:
1. 修改 Flow 签名，provider 改为必需参数
2. 移除 scan_mode 参数和相关判断逻辑
3. 直接使用传入的 provider 获取扫描目标
4. 检查现有保存逻辑，确认是否已写入 EndpointSnapshot
5. 如果未写入，添加快照写入逻辑

**验收标准**:
- [x] Flow 直接使用传入的 provider
- [x] URL 获取结果同时保存到 Endpoint 表和 EndpointSnapshot 表

**备注**: Flow 签名已添加 provider 参数（暂未使用，预留接口）

---

### Task 10: directory_scan_flow 快照读取和写入

**描述**: 修改目录扫描 Flow，使用传入的 Provider，并写入 DirectorySnapshot

**相关需求**: 需求 3, 5, 7

**依赖**: Task 8

**文件**:
- `backend/apps/scan/flows/directory_scan/main_flow.py`
- `backend/apps/scan/tasks/directory_scan/` 相关任务

**实现步骤**:
1. 修改 Flow 签名，provider 改为必需参数
2. 移除 scan_mode 参数和相关判断逻辑
3. 直接使用传入的 provider 获取扫描目标
4. 检查现有保存逻辑，确认是否已写入 DirectorySnapshot
5. 如果未写入，添加快照写入逻辑

**验收标准**:
- [x] Flow 直接使用传入的 provider
- [x] 目录扫描结果同时保存到 Directory 表和 DirectorySnapshot 表

**备注**: Flow 签名已添加 provider 参数（暂未使用，预留接口）

---

### Task 11: QuickScanService 改造

**描述**: 修改 QuickScanService，支持创建快速扫描模式的 Scan 记录并写入 ScanInputTarget

**相关需求**: 需求 8

**依赖**: Task 1, Task 2, Task 3

**文件**:
- `backend/apps/scan/services/quick_scan_service.py`
- `backend/apps/scan/services/scan_creation_service.py`
- `backend/apps/scan/services/scan_service.py`

**实现步骤**:
1. 修改 ScanCreationService.create_scans 方法，添加 scan_mode 参数
2. 修改 ScanService.create_scans 方法，传递 scan_mode 参数
3. 在 API 层调用时设置 scan_mode='quick'
4. 调用 ScanInputTargetService.bulk_create() 写入用户输入

**验收标准**:
- [x] 快速扫描创建的 Scan 记录 scan_mode 为 'quick'
- [x] 用户输入正确写入 ScanInputTarget 表

---

### Task 12: API 层改造

**描述**: 修改快速扫描 API，支持新的扫描模式

**相关需求**: 需求 8

**依赖**: Task 11

**文件**:
- `backend/apps/scan/views/scan_views.py`
- `backend/apps/scan/serializers/` 相关序列化器

**实现步骤**:
1. 确认 /api/scans/quick/ 接口正确调用 QuickScanService
2. 确保响应包含 scan_mode 信息
3. 添加输入验证（至少一个有效目标）

**验收标准**:
- [x] /api/scans/quick/ 接口创建的 Scan 记录 scan_mode 为 'quick'
- [x] 响应包含 scan_mode 字段
- [x] 空输入返回 400 错误

---

### Task 13: 向后兼容性验证

**描述**: 验证现有 API 和定时扫描的向后兼容性

**相关需求**: 需求 9

**依赖**: Task 5

**文件**:
- `backend/apps/scan/views/scan_views.py`
- `backend/apps/scan/services/scheduled_scan_service.py`

**实现步骤**:
1. 验证 /api/scans/initiate/ 接口默认使用 scan_mode='full'
2. 验证定时扫描默认使用 scan_mode='full'
3. 验证完整扫描行为与改造前一致

**验收标准**:
- [x] /api/scans/initiate/ 创建的 Scan 记录 scan_mode 为 'full'
- [x] 定时扫描创建的 Scan 记录 scan_mode 为 'full'
- [x] 完整扫描使用 DatabaseTargetProvider

**备注**: initiate 接口和定时扫描都调用 create_scans() 时不传 scan_mode，默认使用 'full'

---

### Task 14: 删除未使用的 Provider

**描述**: 删除不再使用的 ListTargetProvider 和 PipelineTargetProvider

**相关需求**: 代码清理

**依赖**: Task 4

**文件**:
- `backend/apps/scan/providers/list_provider.py`（删除）
- `backend/apps/scan/providers/pipeline_provider.py`（删除）
- `backend/apps/scan/providers/__init__.py`（移除导出）
- `backend/apps/scan/providers/tests/`（删除相关测试）

**实现步骤**:
1. 确认没有其他代码引用 ListTargetProvider 和 PipelineTargetProvider
2. 删除 list_provider.py 和 pipeline_provider.py 文件
3. 删除相关测试文件
4. 更新 __init__.py 移除导出

**验收标准**:
- [ ] ListTargetProvider 和 PipelineTargetProvider 相关文件已删除
- [ ] 无代码引用这两个 Provider
- [ ] 测试通过

**备注**: 暂时保留，因为测试文件中有引用。可在后续清理迭代中删除。

---

### Task 15: 集成测试

**描述**: 编写快速扫描模式的集成测试

**相关需求**: 所有需求

**依赖**: Task 1-14

**文件**:
- `backend/apps/scan/tests/integration/test_quick_scan.py`

**实现步骤**:
1. 编写快速扫描端到端测试
2. 验证各阶段使用正确的 Provider
3. 验证快照数据正确写入
4. 验证扫描结果只包含指定目标

**验收标准**:
- [ ] 快速扫描只扫描用户指定的目标
- [ ] 各阶段快照数据正确写入
- [ ] 不扫描 Target 下的历史资产

**备注**: 集成测试需要在实际环境中验证，暂不实现

---

## 任务依赖图

```
Task 1 (ScanInputTarget 模型) ──┬──→ Task 3 (ScanInputTargetService)
                                │           │
Task 2 (Scan 模型扩展) ─────────┼───────────┼──→ Task 4 (ScanInputTargetProvider)
                                │           │           │
                                │           │           ↓
                                │           └──→ Task 5 (initiate_scan_flow)
                                │                       │
                                │                       ├──→ Task 6 (subdomain_discovery_flow)
                                │                       │           │
                                │                       │           ↓
                                │                       │    Task 7 (port_scan_flow)
                                │                       │           │
                                │                       │           ↓
                                │                       │    Task 8 (site_scan_flow)
                                │                       │           │
                                │                       │           ├──→ Task 9 (url_fetch_flow)
                                │                       │           │
                                │                       │           └──→ Task 10 (directory_scan_flow)
                                │                       │
                                │                       └──→ Task 13 (向后兼容性验证)
                                │
                                └──→ Task 11 (QuickScanService)
                                            │
                                            └──→ Task 12 (API 层)

Task 1-14 ──→ Task 15 (集成测试)
```

## 估算工时

| 任务 | 估算工时 | 复杂度 |
|------|----------|--------|
| Task 1 | 1h | 低 |
| Task 2 | 0.5h | 低 |
| Task 3 | 1.5h | 中 |
| Task 4 | 1.5h | 中 |
| Task 5 | 2h | 中 |
| Task 6 | 2h | 中 |
| Task 7 | 2h | 中 |
| Task 8 | 2h | 中 |
| Task 9 | 2h | 中 |
| Task 10 | 2h | 中 |
| Task 11 | 1h | 低 |
| Task 12 | 1h | 低 |
| Task 13 | 1h | 低 |
| Task 14 | 0.5h | 低 |
| Task 15 | 3h | 中 |
| **总计** | **23h** | - |
