# 设计文档

## 概述

本设计实现扫描目标提供者的策略模式，将数据源抽象为统一的 `TargetProvider` 接口。这使得扫描任务可以灵活地从不同来源获取目标（数据库、内存列表、快照表、管道输出），同时保持向后兼容。

### 核心价值

1. **解耦扫描范围和结果归属**
   - `DatabaseTargetProvider`: target_id 决定扫描什么（查询数据库所有资产）
   - `ListTargetProvider`: targets 决定扫描什么，target_id 只用于保存结果
   - `SnapshotTargetProvider`: scan_id 决定扫描什么（只扫描本次发现的资产）

2. **支持精确扫描**
   - 快速扫描：用户输入 `a.test.com`，只扫描 `a.test.com` 及其发现的子资产（不扫描整个 test.com 域）
   - 完整扫描：扫描 Target 下的所有资产（test.com + 所有子域名）

3. **代码复用和可测试性**
   - 查询逻辑封装在 Provider 中，避免重复代码
   - 测试时用 ListProvider，不需要数据库

### 快照表方案

快照表用于解决快速扫描的精确控制问题：

**问题场景**：
```
用户输入: a.test.com
创建 Target: test.com (id=1)
    ↓
阶段1: 子域名发现
  发现: b.a.test.com, c.a.test.com
  保存到: Subdomain(target_id=1)
    ↓
阶段2: 端口扫描
  问题: 如何只扫描 b.a.test.com, c.a.test.com？
  
  ❌ 使用 DatabaseTargetProvider(target_id=1)
     → 会扫描 target_id=1 下的所有子域名（包括历史数据 www.test.com, api.test.com）
  
  ✅ 使用 SnapshotTargetProvider(scan_id=100, snapshot_type="subdomain")
     → 只扫描本次扫描（scan_id=100）发现的子域名
```

**快照表流程**：
```
阶段1: 子域名发现
  输入: ListTargetProvider(targets=["a.test.com"])
  输出: b.a.test.com, c.a.test.com
  保存: SubdomainSnapshot(scan_id=100) + Subdomain(target_id=1)
    ↓
阶段2: 端口扫描
  输入: SnapshotTargetProvider(scan_id=100, snapshot_type="subdomain")
  输出: b.a.test.com, c.a.test.com（只读取本次扫描的快照）
  保存: HostPortMappingSnapshot(scan_id=100) + HostPortMapping(target_id=1)
    ↓
阶段3: 网站扫描
  输入: SnapshotTargetProvider(scan_id=100, snapshot_type="host_port")
  输出: 本次扫描发现的 IP:Port
  保存: WebsiteSnapshot(scan_id=100) + Website(target_id=1)
    ↓
阶段4: 端点扫描
  输入: SnapshotTargetProvider(scan_id=100, snapshot_type="website")
  输出: 本次扫描发现的网站 URL
  保存: EndpointSnapshot(scan_id=100) + Endpoint(target_id=1)
```

**快照表优势**：
- ✅ 天然隔离（通过 scan_id）
- ✅ 内存友好（数据在数据库，按需查询）
- ✅ 可追溯历史扫描
- ✅ 支持扫描重放
- ✅ 易于清理旧数据

### 使用场景

| 场景 | Provider | target_id 用途 | 扫描范围 |
|------|----------|---------------|----------|
| **快速扫描（阶段1）** | `ListTargetProvider` + context | 保存结果 | 只扫描用户指定的目标 |
| **快速扫描（阶段2+）** | `SnapshotTargetProvider` | 保存结果 | 只扫描本次扫描发现的资产 |
| **完整扫描** | `DatabaseTargetProvider` | 查询数据 + 保存结果 | 扫描 Target 下所有资产 |
| **定时扫描** | `DatabaseTargetProvider` | 查询数据 + 保存结果 | 扫描 Target 下所有资产 |
| **临时测试** | `ListTargetProvider`（无 context） | 不保存 | 只扫描指定目标 |
| **管道模式（预留）** | `PipelineTargetProvider` | 保存结果 | 使用上一阶段输出 |

## 架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         调用层 (Tasks/Flows)                         │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  export_hosts_task / port_scan_flow / site_scan_flow ...    │   │
│  │                                                              │   │
│  │  # 直接创建具体的 Provider                                    │   │
│  │  provider = SnapshotTargetProvider(scan_id=100, ...)        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                    ↓                                │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                  TargetProvider (抽象基类)                    │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │   │
│  │  │  Database   │ │    List     │ │  Snapshot   │           │   │
│  │  │  Provider   │ │  Provider   │ │  Provider   │           │   │
│  │  │ (target_id) │ │  (targets)  │ │  (scan_id)  │           │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘           │   │
│  │                                                              │   │
│  │  ┌─────────────┐                                             │   │
│  │  │  Pipeline   │  ← 为 Phase 2 管道模式预留                   │   │
│  │  │  Provider   │                                             │   │
│  │  │(StageOutput)│                                             │   │
│  │  └─────────────┘                                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### 文件结构

```
backend/apps/scan/
├── providers/                        # 新增目录
│   ├── __init__.py                   # 导出公共接口
│   ├── base.py                       # TargetProvider 抽象基类 + ProviderContext
│   ├── database_provider.py          # DatabaseTargetProvider
│   ├── list_provider.py              # ListTargetProvider
│   ├── snapshot_provider.py          # SnapshotTargetProvider（新增）
│   └── pipeline_provider.py          # PipelineTargetProvider (预留)
├── services/
│   └── target_export_service.py      # 重构：复用 Provider
└── tasks/
    ├── port_scan/
    │   └── export_hosts_task.py      # 改造：支持 provider 参数
    └── site_scan/
        └── export_site_urls_task.py  # 改造：支持 provider 参数
```

## 组件和接口

### 3.1 ProviderContext 数据类

```python
# backend/apps/scan/providers/base.py

from dataclasses import dataclass
from typing import Optional

@dataclass
class ProviderContext:
    """
    Provider 上下文，携带元数据
    
    Attributes:
        target_id: 关联的 Target ID（用于结果保存），None 表示临时扫描（不保存）
        scan_id: 扫描任务 ID
    
    判断是否保存结果：
        - target_id 不为 None：保存到数据库
        - target_id 为 None：临时扫描，不保存
    """
    target_id: Optional[int] = None
    scan_id: Optional[int] = None
```

### 3.2 TargetProvider 抽象基类

```python
# backend/apps/scan/providers/base.py

from abc import ABC, abstractmethod
from typing import Iterator, Optional, TYPE_CHECKING
import ipaddress

if TYPE_CHECKING:
    from apps.common.utils import BlacklistFilter

class TargetProvider(ABC):
    """
    扫描目标提供者抽象基类
    
    职责：
    - 提供扫描目标（域名、IP、URL 等）的迭代器
    - 提供黑名单过滤器
    - 携带上下文信息（target_id, scan_id 等）
    - 提供 CIDR 展开的通用逻辑
    
    使用方式：
        provider = create_target_provider(target_id=123)
        for host in provider.iter_hosts():
            print(host)
    """
    
    def __init__(self, context: Optional[ProviderContext] = None):
        """
        初始化 Provider
        
        Args:
            context: Provider 上下文，None 时创建默认上下文
        """
        self.context = context or ProviderContext()
    
    @staticmethod
    def _expand_host(host: str) -> Iterator[str]:
        """
        展开主机（如果是 CIDR 则展开为多个 IP，否则直接返回）
        
        这是一个通用的辅助方法，所有子类都可以使用。
        
        Args:
            host: 主机字符串（IP/域名/CIDR）
            
        Yields:
            str: 单个主机（IP 或域名）
            
        示例：
            "192.168.1.0/30" → "192.168.1.1", "192.168.1.2"
            "192.168.1.1" → "192.168.1.1"
            "example.com" → "example.com"
        """
        # 尝试解析为 CIDR
        try:
            network = ipaddress.ip_network(host, strict=False)
            # 如果是单个 IP（/32 或 /128）
            if network.num_addresses == 1:
                yield str(network.network_address)
            else:
                # 展开 CIDR 为多个主机 IP
                for ip in network.hosts():
                    yield str(ip)
        except ValueError:
            # 不是有效的 IP/CIDR，直接返回（可能是域名）
            yield host
    
    @abstractmethod
    def iter_hosts(self) -> Iterator[str]:
        """
        迭代主机列表（域名/IP）
        
        注意：如果输入包含 CIDR，子类应该使用 _expand_host() 展开
        
        Yields:
            str: 主机名或 IP 地址（单个，不包含 CIDR）
        """
        pass
    
    @abstractmethod
    def iter_urls(self) -> Iterator[str]:
        """
        迭代 URL 列表
        
        Yields:
            str: URL 字符串
        """
        pass
    
    @abstractmethod
    def get_blacklist_filter(self) -> Optional['BlacklistFilter']:
        """
        获取黑名单过滤器
        
        Returns:
            BlacklistFilter: 黑名单过滤器实例，或 None（不过滤）
        """
        pass
    
    @property
    def target_id(self) -> Optional[int]:
        """返回关联的 target_id，None 表示临时扫描（不保存）"""
        return self.context.target_id
    
    @property
    def scan_id(self) -> Optional[int]:
        """返回关联的 scan_id"""
        return self.context.scan_id
```

### 3.3 DatabaseTargetProvider

```python
# backend/apps/scan/providers/database_provider.py

from typing import Iterator, Optional
from .base import TargetProvider, ProviderContext

class DatabaseTargetProvider(TargetProvider):
    """
    数据库目标提供者 - 从 Target 表及关联资产表查询
    
    这是现有行为的封装，保持向后兼容。
    
    数据来源：
    - iter_hosts(): 根据 Target 类型返回域名/IP
      - DOMAIN: 根域名 + Subdomain 表
      - IP: 直接返回 IP
      - CIDR: 使用 _expand_host() 展开为所有主机 IP
    - iter_urls(): WebSite/Endpoint 表，带回退链
    
    使用方式：
        provider = DatabaseTargetProvider(target_id=123)
        for host in provider.iter_hosts():
            scan(host)
    """
    
    def __init__(self, target_id: int, context: Optional[ProviderContext] = None):
        """
        初始化数据库目标提供者
        
        Args:
            target_id: 目标 ID（必需）
            context: Provider 上下文
        """
        ctx = context or ProviderContext()
        ctx.target_id = target_id
        super().__init__(ctx)
        self._target_id = target_id
        self._blacklist_filter = None  # 延迟加载
    
    def iter_hosts(self) -> Iterator[str]:
        """
        从数据库查询主机列表
        
        根据 Target 类型决定数据来源：
        - DOMAIN: 根域名 + Subdomain 表
        - IP: 直接返回 target.name
        - CIDR: 使用 _expand_host() 展开 CIDR 范围
        """
        from apps.targets.services import TargetService
        from apps.targets.models import Target
        from apps.asset.services.asset.subdomain_service import SubdomainService
        
        target = TargetService().get_target(self._target_id)
        if not target:
            return
        
        blacklist = self.get_blacklist_filter()
        
        if target.type == Target.TargetType.DOMAIN:
            # 先返回根域名
            if not blacklist or blacklist.is_allowed(target.name):
                yield target.name
            
            # 再返回子域名
            subdomain_service = SubdomainService()
            for domain in subdomain_service.iter_subdomain_names_by_target(
                target_id=self._target_id,
                chunk_size=1000
            ):
                if domain != target.name:  # 避免重复
                    if not blacklist or blacklist.is_allowed(domain):
                        yield domain
        
        elif target.type == Target.TargetType.IP:
            if not blacklist or blacklist.is_allowed(target.name):
                yield target.name
        
        elif target.type == Target.TargetType.CIDR:
            # 使用基类的 _expand_host() 展开 CIDR
            for ip_str in self._expand_host(target.name):
                if not blacklist or blacklist.is_allowed(ip_str):
                    yield ip_str
    
    def iter_urls(self) -> Iterator[str]:
        """
        从数据库查询 URL 列表
        
        使用现有的回退链逻辑：Endpoint → WebSite → Default
        """
        from apps.scan.services.target_export_service import (
            _iter_urls_with_fallback, DataSource
        )
        
        blacklist = self.get_blacklist_filter()
        
        for url, source in _iter_urls_with_fallback(
            target_id=self._target_id,
            sources=[DataSource.ENDPOINT, DataSource.WEBSITE, DataSource.DEFAULT],
            blacklist_filter=blacklist
        ):
            yield url
    
    def get_blacklist_filter(self):
        """获取黑名单过滤器（延迟加载）"""
        if self._blacklist_filter is None:
            from apps.common.services import BlacklistService
            from apps.common.utils import BlacklistFilter
            rules = BlacklistService().get_rules(self._target_id)
            self._blacklist_filter = BlacklistFilter(rules)
        return self._blacklist_filter
```

### 3.4 ListTargetProvider

```python
# backend/apps/scan/providers/list_provider.py

from typing import Iterator, Optional, List
from .base import TargetProvider, ProviderContext

class ListTargetProvider(TargetProvider):
    """
    列表目标提供者 - 直接使用内存中的列表
    
    用于快速扫描、临时扫描等场景，只扫描用户指定的目标。
    
    特点：
    - 不查询数据库
    - 不应用黑名单过滤（用户明确指定的目标）
    - 通过 context 关联 target_id（用于保存结果）
    - 自动检测输入类型（URL/域名/IP/CIDR）
    - 自动展开 CIDR
    
    与 DatabaseTargetProvider 的区别：
    - DatabaseTargetProvider: target_id 决定扫描什么（查询数据库）
    - ListTargetProvider: targets 决定扫描什么，target_id 只用于保存结果
    
    使用方式：
        # 场景1: 快速扫描（需要保存结果）
        # 用户输入: a.test.com
        # 创建 Target: test.com (id=1)
        context = ProviderContext(
            target_id=1,      # 关联到 test.com，用于保存结果
            scan_id=scan.id
        )
        provider = ListTargetProvider(
            targets=["a.test.com"],  # 只扫描用户指定的
            context=context          # 携带 target_id
        )
        for host in provider.iter_hosts():
            scan(host)  # 只扫描 a.test.com
            # 保存结果到 target_id=1
        
        # 场景2: 临时测试（不保存结果）
        provider = ListTargetProvider(targets=["example.com"])
        # target_id=None，扫描任务会跳过保存
        for host in provider.iter_hosts():
            check_reachable(host)  # 不保存结果
    """
    
    def __init__(
        self,
        targets: Optional[List[str]] = None,
        context: Optional[ProviderContext] = None
    ):
        """
        初始化列表目标提供者
        
        Args:
            targets: 目标列表（自动识别类型：URL/域名/IP/CIDR）
            context: Provider 上下文
        """
        from apps.common.validators import detect_input_type
        
        ctx = context or ProviderContext()
        super().__init__(ctx)
        
        # 自动分类目标
        self._hosts = []
        self._urls = []
        
        if targets:
            for target in targets:
                target = target.strip()
                if not target:
                    continue
                
                try:
                    input_type = detect_input_type(target)
                    if input_type == 'url':
                        self._urls.append(target)
                    else:
                        # domain/ip/cidr 都作为 host
                        self._hosts.append(target)
                except ValueError:
                    # 无法识别类型，默认作为 host
                    self._hosts.append(target)
    
    def _iter_raw_hosts(self) -> Iterator[str]:
        """迭代原始主机列表（可能包含 CIDR）"""
        yield from self._hosts
    
    def iter_urls(self) -> Iterator[str]:
        """迭代 URL 列表"""
        yield from self._urls
    
    def get_blacklist_filter(self):
        """列表模式不使用黑名单过滤"""
        return None
```

### 3.5 SnapshotTargetProvider（新增）

```python
# backend/apps/scan/providers/snapshot_provider.py

from typing import Iterator, Optional, Literal
from .base import TargetProvider, ProviderContext

SnapshotType = Literal["subdomain", "website", "endpoint", "host_port"]

class SnapshotTargetProvider(TargetProvider):
    """
    快照目标提供者 - 从快照表读取本次扫描的数据
    
    用于快速扫描的阶段间数据传递，解决精确扫描控制问题。
    
    核心价值：
    - 只返回本次扫描（scan_id）发现的资产
    - 避免扫描历史数据（DatabaseTargetProvider 会扫描所有历史资产）
    
    特点：
    - 通过 scan_id 过滤快照表
    - 不应用黑名单过滤（数据已在上一阶段过滤）
    - 支持多种快照类型（subdomain/website/endpoint/host_port）
    
    使用场景：
        # 快速扫描流程
        用户输入: a.test.com
        创建 Target: test.com (id=1)
        创建 Scan: scan_id=100
        
        # 阶段1: 子域名发现
        provider = ListTargetProvider(
            targets=["a.test.com"],
            context=ProviderContext(target_id=1, scan_id=100)
        )
        # 发现: b.a.test.com, c.a.test.com
        # 保存: SubdomainSnapshot(scan_id=100) + Subdomain(target_id=1)
        
        # 阶段2: 端口扫描
        provider = SnapshotTargetProvider(
            scan_id=100,
            snapshot_type="subdomain",
            context=ProviderContext(target_id=1, scan_id=100)
        )
        # 只返回: b.a.test.com, c.a.test.com（本次扫描发现的）
        # 不返回: www.test.com, api.test.com（历史数据）
        
        # 阶段3: 网站扫描
        provider = SnapshotTargetProvider(
            scan_id=100,
            snapshot_type="host_port",
            context=ProviderContext(target_id=1, scan_id=100)
        )
        # 只返回本次扫描发现的 IP:Port
    """
    
    def __init__(
        self,
        scan_id: int,
        snapshot_type: SnapshotType,
        context: Optional[ProviderContext] = None
    ):
        """
        初始化快照目标提供者
        
        Args:
            scan_id: 扫描任务 ID（必需）
            snapshot_type: 快照类型
                - "subdomain": 子域名快照（SubdomainSnapshot）
                - "website": 网站快照（WebsiteSnapshot）
                - "endpoint": 端点快照（EndpointSnapshot）
                - "host_port": 主机端口映射快照（HostPortMappingSnapshot）
            context: Provider 上下文
        """
        ctx = context or ProviderContext()
        ctx.scan_id = scan_id
        super().__init__(ctx)
        self._scan_id = scan_id
        self._snapshot_type = snapshot_type
    
    def _iter_raw_hosts(self) -> Iterator[str]:
        """
        从快照表迭代主机列表
        
        根据 snapshot_type 选择不同的快照表：
        - subdomain: SubdomainSnapshot.name
        - host_port: HostPortMappingSnapshot.host
        """
        if self._snapshot_type == "subdomain":
            from apps.asset.services.snapshot import SubdomainSnapshotsService
            service = SubdomainSnapshotsService()
            yield from service.iter_subdomain_names_by_scan(
                scan_id=self._scan_id,
                chunk_size=1000
            )
        
        elif self._snapshot_type == "host_port":
            from apps.asset.services.snapshot import HostPortMappingSnapshotsService
            service = HostPortMappingSnapshotsService()
            # 返回 host:port 格式（用于网站扫描）
            for mapping in service.iter_by_scan(scan_id=self._scan_id, chunk_size=1000):
                yield f"{mapping.host}:{mapping.port}"
        
        else:
            # 其他类型暂不支持 iter_hosts
            return
    
    def iter_urls(self) -> Iterator[str]:
        """
        从快照表迭代 URL 列表
        
        根据 snapshot_type 选择不同的快照表：
        - website: WebsiteSnapshot.url
        - endpoint: EndpointSnapshot.url
        """
        if self._snapshot_type == "website":
            from apps.asset.services.snapshot import WebsiteSnapshotsService
            service = WebsiteSnapshotsService()
            for website in service.iter_by_scan(scan_id=self._scan_id, chunk_size=1000):
                yield website.url
        
        elif self._snapshot_type == "endpoint":
            from apps.asset.services.snapshot import EndpointSnapshotsService
            service = EndpointSnapshotsService()
            for endpoint in service.iter_by_scan(scan_id=self._scan_id, chunk_size=1000):
                yield endpoint.url
        
        else:
            # 其他类型暂不支持 iter_urls
            return
    
    def get_blacklist_filter(self) -> None:
        """快照数据已在上一阶段过滤过了"""
        return None
```

### 3.6 PipelineTargetProvider（预留）

```python
# backend/apps/scan/providers/pipeline_provider.py

from typing import Iterator, Optional, TYPE_CHECKING
from .base import TargetProvider, ProviderContext

if TYPE_CHECKING:
    from apps.scan.pipeline.data import StageOutput

class PipelineTargetProvider(TargetProvider):
    """
    管道目标提供者 - 使用上一阶段的输出
    
    用于 Phase 2 管道模式的阶段间数据传递。
    
    特点：
    - 不查询数据库
    - 不应用黑名单过滤（数据已在上一阶段过滤）
    - 直接使用 StageOutput 中的数据
    
    使用方式（Phase 2）：
        stage1_output = stage1.run(input)
        provider = PipelineTargetProvider(
            previous_output=stage1_output,
            target_id=123
        )
        for host in provider.iter_hosts():
            stage2.scan(host)
    """
    
    def __init__(
        self,
        previous_output: 'StageOutput',
        target_id: Optional[int] = None,
        context: Optional[ProviderContext] = None
    ):
        """
        初始化管道目标提供者
        
        Args:
            previous_output: 上一阶段的输出
            target_id: 可选，关联到某个 Target（用于保存结果）
            context: Provider 上下文
        """
        ctx = context or ProviderContext(target_id=target_id)
        super().__init__(ctx)
        self._previous_output = previous_output
    
    def _iter_raw_hosts(self) -> Iterator[str]:
        """迭代上一阶段输出的原始主机（可能包含 CIDR）"""
        yield from self._previous_output.hosts
    
    def iter_urls(self) -> Iterator[str]:
        """迭代上一阶段输出的 URL"""
        yield from self._previous_output.urls
    
    def get_blacklist_filter(self) -> None:
        """管道传递的数据已经过滤过了"""
        return None
```

### 3.7 模块导出

```python
# backend/apps/scan/providers/__init__.py

"""
扫描目标提供者模块

提供统一的目标获取接口，支持多种数据源：
- DatabaseTargetProvider: 从数据库查询（完整扫描）
- ListTargetProvider: 使用内存列表（快速扫描阶段1）
- SnapshotTargetProvider: 从快照表读取（快速扫描阶段2+）
- PipelineTargetProvider: 使用管道输出（Phase 2）

使用方式：
    from apps.scan.providers import (
        DatabaseTargetProvider,
        ListTargetProvider,
        SnapshotTargetProvider,
        ProviderContext
    )
    
    # 数据库模式（完整扫描）
    provider = DatabaseTargetProvider(target_id=123)
    
    # 列表模式（快速扫描阶段1）
    context = ProviderContext(target_id=1, scan_id=100)
    provider = ListTargetProvider(
        targets=["a.test.com"],
        context=context
    )
    
    # 快照模式（快速扫描阶段2+）
    context = ProviderContext(target_id=1, scan_id=100)
    provider = SnapshotTargetProvider(
        scan_id=100,
        snapshot_type="subdomain",
        context=context
    )
    
    # 使用 Provider
    for host in provider.iter_hosts():
        scan(host)
"""

from .base import TargetProvider, ProviderContext
from .database_provider import DatabaseTargetProvider
from .list_provider import ListTargetProvider
from .snapshot_provider import SnapshotTargetProvider, SnapshotType
from .pipeline_provider import PipelineTargetProvider, StageOutput

__all__ = [
    'TargetProvider',
    'ProviderContext',
    'DatabaseTargetProvider',
    'ListTargetProvider',
    'SnapshotTargetProvider',
    'SnapshotType',
    'PipelineTargetProvider',
    'StageOutput',
]
```

## 数据模型

### 4.1 ProviderContext

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| target_id | Optional[int] | None | 关联的 Target ID，None 表示临时扫描（不保存） |
| scan_id | Optional[int] | None | 扫描任务 ID |

### 4.2 StageOutput（Phase 2 预留）

```python
# backend/apps/scan/pipeline/data.py（Phase 2 实现）

@dataclass
class StageOutput:
    """阶段输出数据"""
    hosts: List[str] = field(default_factory=list)
    urls: List[str] = field(default_factory=list)
    new_targets: List[str] = field(default_factory=list)
    stats: Dict[str, Any] = field(default_factory=dict)
    success: bool = True
    error: Optional[str] = None
```



## 正确性属性

*正确性属性是系统在所有有效执行中应保持为真的特征或行为——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: ListTargetProvider Round-Trip

*For any* 主机列表和 URL 列表，创建 ListTargetProvider 后迭代 iter_hosts() 和 iter_urls() 应该返回与输入相同的元素（顺序相同）。

**Validates: Requirements 3.1, 3.2**

### Property 2: PipelineTargetProvider Round-Trip

*For any* StageOutput 对象，PipelineTargetProvider 的 iter_hosts() 和 iter_urls() 应该返回与 StageOutput 中 hosts 和 urls 列表相同的元素。

**Validates: Requirements 4.1, 4.2**

### Property 3: Factory Provider Type Selection

*For any* 参数组合，create_target_provider 工厂函数应该根据优先级规则返回正确类型的 Provider：
- previous_output 存在 → PipelineTargetProvider
- targets 存在 → ListTargetProvider
- 仅 target_id 存在 → DatabaseTargetProvider

**Validates: Requirements 5.1, 5.2, 5.3, 5.4**

### Property 4: Context Propagation

*For any* ProviderContext，传入 Provider 构造函数后，Provider 的 target_id 和 scan_id 属性应该与 context 中的值一致。

**Validates: Requirements 1.3, 1.5, 7.4, 7.5**

### Property 5: Non-Database Provider Blacklist Filter

*For any* ListTargetProvider 或 PipelineTargetProvider 实例，get_blacklist_filter() 方法应该返回 None。

**Validates: Requirements 3.4, 9.4, 9.5**

### Property 6: DatabaseTargetProvider Blacklist Application

*For any* 带有黑名单规则的 target_id，DatabaseTargetProvider 的 iter_hosts() 和 iter_urls() 应该过滤掉匹配黑名单规则的目标。

**Validates: Requirements 2.3, 9.1, 9.2, 9.3**

### Property 7: CIDR Expansion Consistency

*For any* CIDR 字符串（如 "192.168.1.0/24"），所有 Provider（DatabaseTargetProvider、ListTargetProvider）的 iter_hosts() 方法应该将其展开为相同的单个 IP 地址列表。

**Validates: Requirements 1.1, 3.6**

### Property 8: Task Backward Compatibility

*For any* 任务调用，当仅提供 target_id 参数时，任务应该创建 DatabaseTargetProvider 并使用它进行数据访问，行为与改造前一致。

**Validates: Requirements 6.1, 6.2, 6.4, 6.5**

## 错误处理

### 6.1 工厂函数错误

| 错误场景 | 异常类型 | 错误消息 |
|----------|----------|----------|
| 未提供任何有效参数 | ValueError | "必须提供以下参数之一: target_id, targets, previous_output" |

### 6.2 DatabaseTargetProvider 错误

| 错误场景 | 处理方式 |
|----------|----------|
| target_id 不存在 | iter_hosts() 返回空迭代器，不抛出异常 |
| 数据库连接失败 | 抛出 Django 数据库异常 |

### 6.3 Task 错误

| 错误场景 | 异常类型 | 错误消息 |
|----------|----------|----------|
| 既未提供 provider 也未提供 target_id | ValueError | "必须提供 target_id 或 provider" |

## 测试策略

### 7.1 测试框架

- **单元测试框架**: pytest
- **属性测试框架**: hypothesis
- **Mock 框架**: pytest-mock / unittest.mock

### 7.2 测试类型

#### 单元测试

| 测试目标 | 测试内容 |
|----------|----------|
| ProviderContext | 默认值、字段赋值 |
| ListTargetProvider | 空列表、单元素、多元素、类型自动识别 |
| PipelineTargetProvider | 空 StageOutput、正常 StageOutput |
| create_target_provider | 各种参数组合、优先级验证 |

#### 属性测试

| Property | 测试策略 |
|----------|----------|
| Property 1 | 生成随机字符串列表，验证 round-trip |
| Property 2 | 生成随机 StageOutput，验证 round-trip |
| Property 3 | 生成随机参数组合，验证返回类型 |
| Property 4 | 生成随机 ProviderContext，验证属性传递 |
| Property 5 | 对所有非数据库 Provider 验证返回 None |
| Property 6 | 需要数据库 fixture，验证黑名单过滤 |
| Property 7 | 验证 CIDR 展开一致性 |
| Property 8 | 需要 mock，验证向后兼容 |

### 7.3 测试配置

```python
# pytest.ini 或 pyproject.toml
[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]

# hypothesis 配置
[tool.hypothesis]
max_examples = 100
```

### 7.4 测试文件结构

```
backend/tests/scan/providers/
├── __init__.py
├── test_base.py                  # ProviderContext, TargetProvider 基类测试
├── test_list_provider.py         # ListTargetProvider 单元测试 + 属性测试
├── test_pipeline_provider.py     # PipelineTargetProvider 单元测试 + 属性测试
├── test_database_provider.py     # DatabaseTargetProvider 集成测试
├── test_factory.py               # create_target_provider 单元测试 + 属性测试
└── conftest.py                   # 共享 fixtures
```
