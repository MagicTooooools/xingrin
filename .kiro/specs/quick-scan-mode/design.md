# 设计文档

## 概述

本设计实现快速扫描模式，允许用户只扫描指定的目标（而非整个 Target 下的所有资产）。通过新建 `ScanInputTarget` 表存储用户输入，配合 `ScanInputTargetProvider` 和现有快照表，实现阶段间的精确数据传递。

### 核心价值

1. **精确扫描控制**
   - 用户输入 `a.test.com`，只扫描 `a.test.com` 及其发现的子资产
   - 不扫描 `test.com` 下的历史资产（如 `www.test.com`、`api.test.com`）

2. **支持大量输入**
   - 新建 `ScanInputTarget` 表存储用户输入（支持 1 万+ 条）
   - 新建 `ScanInputTargetProvider` 分块迭代读取
   - 复用 `SnapshotTargetProvider`（阶段2+）
   - 复用 6 个已存在的快照表

3. **向后兼容**
   - 默认使用完整扫描模式（`scan_mode='full'`）
   - 现有 API 和定时扫描不受影响

## 架构

### 快速扫描 vs 完整扫描流程对比

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           完整扫描（现有行为）                               │
│                                                                              │
│  用户输入: a.test.com                                                        │
│  创建 Target: test.com (id=1)                                               │
│                                                                              │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐     │
│  │ 子域名发现  │ → │ 端口扫描    │ → │ 网站扫描    │ → │ URL获取     │     │
│  │             │   │             │   │             │   │             │     │
│  │ Database    │   │ Database    │   │ Database    │   │ Database    │     │
│  │ Provider    │   │ Provider    │   │ Provider    │   │ Provider    │     │
│  │ (target=1)  │   │ (target=1)  │   │ (target=1)  │   │ (target=1)  │     │
│  └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘     │
│        ↓                 ↓                 ↓                 ↓              │
│  扫描 test.com     扫描所有子域名    扫描所有端口      扫描所有网站        │
│  下所有子域名      (包括历史数据)    (包括历史数据)    (包括历史数据)      │
│                                                                              │
│  问题：用户只想扫描 a.test.com，但系统扫描了整个 test.com 域                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                           快速扫描（新行为）                                 │
│                                                                              │
│  用户输入: 5000 个目标（域名/IP/URL 混合）                                   │
│  创建 Target: test.com (id=1)                                               │
│  创建 Scan: scan_id=100, scan_mode='quick'                                  │
│  写入 ScanInputTarget 表: 5000 条记录 (scan_id=100)                         │
│                                                                              │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐     │
│  │ 子域名发现  │ → │ 端口扫描    │ → │ 网站扫描    │ → │ URL获取     │     │
│  │             │   │             │   │             │   │             │     │
│  │ ScanInput   │   │ Snapshot    │   │ Snapshot    │   │ Snapshot    │     │
│  │ Target      │   │ Provider    │   │ Provider    │   │ Provider    │     │
│  │ Provider    │   │ (scan=100,  │   │ (scan=100,  │   │ (scan=100,  │     │
│  │ (scan=100)  │   │ type=       │   │ type=       │   │ type=       │     │
│  │             │   │ subdomain)  │   │ host_port)  │   │ website)    │     │
│  └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘     │
│        ↓                 ↓                 ↓                 ↓              │
│        │                 │                 │                 │              │
│        ↓                 ↓                 ↓                 ↓              │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐     │
│  │ Subdomain   │   │ HostPort    │   │ Website     │   │ Endpoint    │     │
│  │ Snapshot    │   │ Mapping     │   │ Snapshot    │   │ Snapshot    │     │
│  │ (scan=100)  │   │ Snapshot    │   │ (scan=100)  │   │ (scan=100)  │     │
│  │             │   │ (scan=100)  │   │             │   │             │     │
│  └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘     │
│                                                                              │
│  结果：只扫描用户输入的目标及其发现的子资产，不扫描历史数据                 │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Provider 选择逻辑

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        initiate_scan_flow                                   │
│                                                                              │
│  读取 Scan.scan_mode，统一创建所有阶段的 Provider                           │
│                                                                              │
│  IF scan_mode == 'quick':                                                   │
│      providers = {                                                          │
│          'subdomain_discovery': ScanInputTargetProvider(scan_id),           │
│          'port_scan': SnapshotTargetProvider(scan_id, type='subdomain'),    │
│          'site_scan': SnapshotTargetProvider(scan_id, type='host_port'),    │
│          'url_fetch': SnapshotTargetProvider(scan_id, type='website'),      │
│          'directory_scan': SnapshotTargetProvider(scan_id, type='website'), │
│      }                                                                      │
│  ELSE:  # scan_mode == 'full'                                               │
│      providers = {                                                          │
│          'subdomain_discovery': DatabaseTargetProvider(target_id),          │
│          'port_scan': DatabaseTargetProvider(target_id),                    │
│          'site_scan': DatabaseTargetProvider(target_id),                    │
│          'url_fetch': DatabaseTargetProvider(target_id),                    │
│          'directory_scan': DatabaseTargetProvider(target_id),               │
│      }                                                                      │
│                                                                              │
│  调用子 Flow 时传入对应的 provider                                          │
└─────────────────────────────────────────────────────────────────────────────┘
                    │
    ┌───────────────┼───────────────┬───────────────┬───────────────┐
    ↓               ↓               ↓               ↓               ↓
subdomain_      port_scan_      site_scan_      url_fetch_      directory_
discovery_flow  flow            flow            flow            scan_flow
(provider)      (provider)      (provider)      (provider)      (provider)

每个子 Flow 只管用传入的 provider，不关心 scan_mode
```

**设计原则**：
- Provider 选择逻辑集中在 initiate_scan_flow（编排层）
- 子 Flow 只负责使用 Provider，职责单一
- 不传递 scan_mode 参数，减少耦合

### 文件结构

```
backend/apps/scan/
├── providers/                        # 已存在
│   ├── __init__.py
│   ├── base.py                       # TargetProvider 抽象基类
│   ├── database_provider.py          # DatabaseTargetProvider
│   ├── list_provider.py              # ListTargetProvider
│   ├── snapshot_provider.py          # SnapshotTargetProvider
│   └── scan_input_provider.py        # 新增：ScanInputTargetProvider
├── models/
│   ├── scan.py                       # 需修改：添加 scan_mode 字段
│   └── scan_input_target.py          # 新增：ScanInputTarget 模型
├── flows/
│   ├── initiate_scan_flow.py         # 需修改：Provider 选择逻辑
│   ├── subdomain_discovery_flow.py   # 需修改：写入快照表
│   ├── port_scan_flow.py             # 需修改：写入快照表
│   ├── site_scan_flow.py             # 需修改：写入快照表
│   └── url_fetch/
│       └── main_flow.py              # 需修改：写入快照表
├── services/
│   ├── quick_scan_service.py         # 已存在，需修改
│   └── scan_input_target_service.py  # 新增：ScanInputTarget 服务
└── views/
    └── scan_views.py                 # 需修改：API 层

backend/apps/asset/
├── models/
│   └── snapshot_models.py            # 已存在（6个快照表）
└── services/
    └── snapshot/                     # 已存在
        ├── subdomain_snapshots_service.py
        ├── website_snapshots_service.py
        ├── endpoint_snapshots_service.py
        └── host_port_mapping_snapshots_service.py
```

## 组件和接口

### 3.1 ScanInputTarget 模型（新增）

```python
# backend/apps/scan/models/scan_input_target.py

class ScanInputTarget(models.Model):
    """
    扫描输入目标表
    
    存储快速扫描时用户输入的目标，支持大量数据（1万+）的分块迭代。
    """
    
    class InputType(models.TextChoices):
        DOMAIN = 'domain', '域名'
        IP = 'ip', 'IP地址'
        CIDR = 'cidr', 'CIDR'
        URL = 'url', 'URL'
    
    id = models.AutoField(primary_key=True)
    scan = models.ForeignKey(
        'scan.Scan',
        on_delete=models.CASCADE,
        related_name='input_targets',
        help_text='所属的扫描任务'
    )
    value = models.CharField(max_length=2000, help_text='用户输入的原始值')
    input_type = models.CharField(
        max_length=10,
        choices=InputType.choices,
        help_text='输入类型'
    )
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        db_table = 'scan_input_target'
        indexes = [
            models.Index(fields=['scan']),
            models.Index(fields=['input_type']),
        ]
```

### 3.2 Scan 模型扩展

```python
# backend/apps/scan/models/scan.py

class Scan(models.Model):
    """扫描任务模型"""
    
    class ScanMode(models.TextChoices):
        FULL = 'full', '完整扫描'
        QUICK = 'quick', '快速扫描'
    
    # ... 现有字段 ...
    
    # 新增字段
    scan_mode = models.CharField(
        max_length=10,
        choices=ScanMode.choices,
        default=ScanMode.FULL,
        help_text='扫描模式：full=完整扫描，quick=快速扫描'
    )
    # 注意：用户输入存储在 ScanInputTarget 表，通过 input_targets 关联
```

### 3.3 ScanInputTargetProvider（新增）

```python
# backend/apps/scan/providers/scan_input_provider.py

class ScanInputTargetProvider(TargetProvider):
    """
    扫描输入目标提供者 - 从 ScanInputTarget 表读取用户输入
    
    用于快速扫描的第一阶段，支持大量输入的分块迭代。
    
    特点：
    - 通过 scan_id 查询 ScanInputTarget 表
    - 按 input_type 分类返回 hosts/urls
    - 支持分块迭代（chunk_size=1000）
    - 不应用黑名单过滤（用户明确指定的目标）
    """
    
    def __init__(self, scan_id: int, context: Optional[ProviderContext] = None):
        ctx = context or ProviderContext()
        ctx.scan_id = scan_id
        super().__init__(ctx)
        self._scan_id = scan_id
    
    def _iter_raw_hosts(self) -> Iterator[str]:
        """迭代 domain/ip/cidr 类型的输入"""
        from apps.scan.models import ScanInputTarget
        queryset = ScanInputTarget.objects.filter(
            scan_id=self._scan_id,
            input_type__in=['domain', 'ip', 'cidr']
        )
        for item in queryset.iterator(chunk_size=1000):
            yield item.value
    
    def iter_urls(self) -> Iterator[str]:
        """迭代 url 类型的输入"""
        from apps.scan.models import ScanInputTarget
        queryset = ScanInputTarget.objects.filter(
            scan_id=self._scan_id,
            input_type='url'
        )
        for item in queryset.iterator(chunk_size=1000):
            yield item.value
    
    def get_blacklist_filter(self) -> None:
        """用户输入不使用黑名单过滤"""
        return None
```

### 3.4 initiate_scan_flow 改造

```python
# backend/apps/scan/flows/initiate_scan_flow.py

@flow(name='initiate_scan')
def initiate_scan_flow(
    scan_id: int,
    target_name: str,
    target_id: int,
    scan_workspace_dir: str,
    engine_name: str,
    scheduled_scan_name: str | None = None,
) -> dict:
    """
    初始化扫描任务
    
    根据 scan_mode 统一创建所有阶段的 Provider，然后传给各子 Flow。
    子 Flow 只负责使用 Provider，不关心 scan_mode。
    """
    # ... 现有代码 ...
    
    # ==================== 创建 Provider ====================
    from apps.scan.models import Scan
    scan = Scan.objects.get(id=scan_id)
    
    provider_context = ProviderContext(target_id=target_id, scan_id=scan_id)
    
    if scan.scan_mode == Scan.ScanMode.QUICK:
        # 快速扫描：各阶段使用不同的 Provider
        providers = {
            'subdomain_discovery': ScanInputTargetProvider(
                scan_id=scan_id,
                context=provider_context
            ),
            'port_scan': SnapshotTargetProvider(
                scan_id=scan_id,
                snapshot_type='subdomain',
                context=provider_context
            ),
            'site_scan': SnapshotTargetProvider(
                scan_id=scan_id,
                snapshot_type='host_port',
                context=provider_context
            ),
            'url_fetch': SnapshotTargetProvider(
                scan_id=scan_id,
                snapshot_type='website',
                context=provider_context
            ),
            'directory_scan': SnapshotTargetProvider(
                scan_id=scan_id,
                snapshot_type='website',
                context=provider_context
            ),
        }
        logger.info(f"✓ 快速扫描模式 - 创建各阶段 Provider")
    else:
        # 完整扫描：所有阶段使用 DatabaseTargetProvider
        db_provider = DatabaseTargetProvider(target_id=target_id, context=provider_context)
        providers = {
            'subdomain_discovery': db_provider,
            'port_scan': db_provider,
            'site_scan': db_provider,
            'url_fetch': db_provider,
            'directory_scan': db_provider,
        }
        logger.info(f"✓ 完整扫描模式 - 使用 DatabaseTargetProvider")
    
    # 调用子 Flow 时传入对应的 provider
    # flow_kwargs['provider'] = providers[scan_type]
    # ... 后续代码 ...
```

### 3.5 子 Flow 接口

子 Flow 只接收 provider 参数，直接使用，不关心 scan_mode：

```python
# 示例：port_scan_flow.py

@flow(name='port_scan')
def port_scan_flow(
    scan_id: int,
    target_name: str,
    target_id: int,
    scan_workspace_dir: str,
    provider: TargetProvider,  # 必需参数，由 initiate_scan_flow 传入
    enabled_tools: dict = None,
) -> dict:
    """
    端口扫描 Flow
    
    直接使用传入的 provider 获取扫描目标，不关心 scan_mode。
    """
    # 使用 provider 获取主机列表
    for host in provider.iter_hosts():
        # 扫描逻辑...
        pass
```

### 3.6 快照写入逻辑

各 Flow 在保存数据时，需要同时写入快照表。已有的快照服务提供了 `save_and_sync` 方法：

```python
# 示例：subdomain_discovery_flow.py 中的保存逻辑

from apps.asset.services.snapshot import SubdomainSnapshotsService
from apps.asset.dtos import SubdomainSnapshotDTO

# 构建快照 DTO
snapshot_dtos = [
    SubdomainSnapshotDTO(
        scan_id=scan_id,
        target_id=target_id,
        name=subdomain_name
    )
    for subdomain_name in discovered_subdomains
]

# 保存到快照表并同步到资产表
snapshot_service = SubdomainSnapshotsService()
snapshot_service.save_and_sync(snapshot_dtos)
```

### 3.7 API 层改造

```python
# backend/apps/scan/views/scan_views.py

class QuickScanView(APIView):
    """快速扫描 API"""
    
    def post(self, request):
        """
        发起快速扫描
        
        请求体:
        {
            "inputs": ["a.com", "b.com", "192.168.1.1", ...],  # 支持大量输入
            "engineId": 1,
            "configuration": "..."
        }
        
        设计原则：每个 Target 创建一个 Scan
        - a.com 和 b.com 是不同的根域名，应该是不同的 Target
        - 每个 Target 有独立的 Scan 记录，语义清晰
        - 可以独立查看每个 Target 的扫描结果
        """
        inputs = request.data.get('inputs', [])
        engine_id = request.data.get('engine_id')
        configuration = request.data.get('configuration', '')
        
        # 1. 解析输入，创建 Target 和资产
        quick_scan_service = QuickScanService()
        result = quick_scan_service.process_quick_scan(
            inputs=inputs,
            engine_id=engine_id,
            create_scan=True,
            yaml_configuration=configuration
        )
        
        # 2. 返回结果（包含多个 Scan）
        return Response({
            'count': len(result['scans']),
            'scans': [ScanSerializer(s).data for s in result['scans']],
            'target_stats': result['target_stats'],
            'asset_stats': result['asset_stats'],
            'errors': result['errors']
        })
```

## 数据模型

### 4.1 ScanInputTarget 表（新增）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | AutoField | 主键 |
| scan_id | ForeignKey | 关联的 Scan ID |
| value | CharField(2000) | 用户输入的原始值 |
| input_type | CharField(10) | 输入类型：domain/ip/cidr/url |
| created_at | DateTimeField | 创建时间 |

索引：`scan_id`, `input_type`

### 4.2 Scan 模型新增字段

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| scan_mode | CharField(10) | 'full' | 扫描模式：full/quick |

注意：用户输入存储在 `ScanInputTarget` 表，通过 `scan.input_targets` 关联访问。

### 4.2 快照表（已存在）

| 快照表 | 用途 | 关键字段 |
|--------|------|----------|
| SubdomainSnapshot | 子域名快照 | scan_id, name |
| HostPortMappingSnapshot | 主机端口映射快照 | scan_id, host, ip, port |
| WebsiteSnapshot | 网站快照 | scan_id, url, host |
| DirectorySnapshot | 目录快照 | scan_id, url |
| EndpointSnapshot | 端点快照 | scan_id, url, host |
| VulnerabilitySnapshot | 漏洞快照 | scan_id, url, vuln_type |

### 4.3 阶段间数据传递映射

| 阶段 | 输入 Provider | 输入 snapshot_type | 输出快照表 |
|------|---------------|-------------------|------------|
| 子域名发现 | ListTargetProvider | - | SubdomainSnapshot |
| 端口扫描 | SnapshotTargetProvider | subdomain | HostPortMappingSnapshot |
| 网站扫描 | SnapshotTargetProvider | host_port | WebsiteSnapshot |
| URL获取 | SnapshotTargetProvider | website | EndpointSnapshot |
| 目录扫描 | SnapshotTargetProvider | website | DirectorySnapshot |
| 漏洞扫描 | SnapshotTargetProvider | endpoint | VulnerabilitySnapshot |

## 正确性属性

### Property 1: 扫描模式隔离

*For any* 快速扫描任务，系统只扫描用户输入的目标及其发现的子资产，不扫描 Target 下的历史资产。

**验证方法**：
- 创建 Target，添加历史子域名
- 发起快速扫描，指定新的子域名
- 验证扫描结果只包含新子域名及其发现的资产

**Validates: Requirements 1, 6**

### Property 2: 快照数据完整性

*For any* 扫描任务，每个阶段发现的资产都会写入对应的快照表，且 scan_id 正确关联。

**验证方法**：
- 发起扫描任务
- 验证各快照表中的记录都有正确的 scan_id
- 验证快照数量与扫描结果一致

**Validates: Requirements 2, 3, 7**

### Property 3: Provider 选择正确性

*For any* 扫描任务：
- scan_mode='full' 时，所有阶段使用 DatabaseTargetProvider
- scan_mode='quick' 时，阶段1 使用 ListTargetProvider，阶段2+ 使用 SnapshotTargetProvider

**验证方法**：
- 分别发起完整扫描和快速扫描
- 验证各阶段使用的 Provider 类型

**Validates: Requirements 5, 6**

### Property 4: 向后兼容性

*For any* 现有 API 调用（不指定 scan_mode），系统默认使用完整扫描模式，行为与改造前一致。

**验证方法**：
- 使用现有 API 发起扫描
- 验证 scan_mode 默认为 'full'
- 验证扫描行为与改造前一致

**Validates: Requirements 9**

### Property 5: 快照级联删除

*For any* Scan 记录被删除时，关联的所有快照数据都会被级联删除。

**验证方法**：
- 创建扫描任务，生成快照数据
- 删除 Scan 记录
- 验证所有关联的快照表记录都被删除

**Validates: Requirements 10**

## 错误处理

### 6.1 输入验证错误

| 错误场景 | 处理方式 |
|----------|----------|
| user_input_targets 为空 | 返回 400 错误，提示"请输入扫描目标" |
| 输入格式无效 | 返回解析错误列表，继续处理有效输入 |
| 所有输入都无效 | 返回 400 错误，提示"没有有效的扫描目标" |

### 6.2 扫描执行错误

| 错误场景 | 处理方式 |
|----------|----------|
| 快照表为空 | SnapshotTargetProvider 返回空迭代器，阶段跳过 |
| 快照服务异常 | 记录错误日志，继续执行后续阶段 |
| Provider 创建失败 | 回退到 DatabaseTargetProvider |

## 测试策略

### 7.1 单元测试

| 测试目标 | 测试内容 |
|----------|----------|
| Scan 模型 | scan_mode 字段默认值、user_input_targets 序列化 |
| Provider 选择逻辑 | 根据 scan_mode 返回正确的 Provider 类型 |
| 快照写入 | 各 Flow 正确写入对应的快照表 |

### 7.2 集成测试

| 测试场景 | 验证内容 |
|----------|----------|
| 完整扫描流程 | 使用 DatabaseTargetProvider，扫描所有资产 |
| 快速扫描流程 | 使用 ListTargetProvider + SnapshotTargetProvider，只扫描指定目标 |
| 混合场景 | 同一 Target 下同时存在完整扫描和快速扫描 |

### 7.3 测试文件结构

```
backend/apps/scan/tests/
├── test_scan_model.py              # Scan 模型测试
├── test_initiate_scan_flow.py      # Provider 选择逻辑测试
├── flows/
│   ├── test_subdomain_discovery_flow.py  # 快照写入测试
│   ├── test_port_scan_flow.py
│   └── ...
└── integration/
    ├── test_full_scan.py           # 完整扫描集成测试
    └── test_quick_scan.py          # 快速扫描集成测试
```
