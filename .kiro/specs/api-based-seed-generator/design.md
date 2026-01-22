# Design Document

## Overview

本文档描述基于 API 的种子数据生成器的设计。该生成器是一个 Python 脚本，通过调用 Go 后端的 REST API 来创建测试数据。设计重点是独立性、可靠性和易用性。

## Architecture

### 整体架构

```
┌─────────────────────────────────────────────────────┐
│           Python Seed Generator                      │
│                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │ CLI Parser   │  │ API Client   │  │ Data      │ │
│  │              │  │              │  │ Generator │ │
│  └──────┬───────┘  └──────┬───────┘  └─────┬─────┘ │
│         │                 │                 │       │
│         └─────────────────┼─────────────────┘       │
│                           │                         │
└───────────────────────────┼─────────────────────────┘
                            │ HTTP/JSON
                            ▼
                  ┌──────────────────┐
                  │   Go Backend     │
                  │   (Port 8888)    │
                  └──────────────────┘
```

### 模块划分

| 模块 | 职责 | 文件 |
|------|------|------|
| CLI Parser | 解析命令行参数 | `seed_generator.py` (main) |
| API Client | HTTP 请求封装、认证管理 | `api_client.py` |
| Data Generator | 生成随机测试数据 | `data_generator.py` |
| Progress Tracker | 进度显示和统计 | `progress.py` |
| Error Handler | 错误处理和重试 | `error_handler.py` |

## Components and Interfaces

### 1. API Client

**职责：** 封装所有 HTTP 请求，管理认证 token

**接口：**

```python
class APIClient:
    def __init__(self, base_url: str, username: str, password: str):
        """初始化 API 客户端"""
        
    def login(self) -> str:
        """登录并获取 JWT token"""
        
    def refresh_token(self) -> str:
        """刷新过期的 token"""
        
    def post(self, endpoint: str, data: dict) -> dict:
        """发送 POST 请求"""
        
    def get(self, endpoint: str, params: dict = None) -> dict:
        """发送 GET 请求"""
        
    def delete(self, endpoint: str) -> None:
        """发送 DELETE 请求"""
```

**实现细节：**

- 使用 `requests.Session` 保持连接
- 自动在请求头中添加 `Authorization: Bearer {token}`
- Token 过期时自动调用 `refresh_token()`
- 所有请求使用 30 秒超时
- 返回解析后的 JSON 数据

### 2. Data Generator

**职责：** 生成随机但合理的测试数据

**接口：**

```python
class DataGenerator:
    @staticmethod
    def generate_organization(index: int) -> dict:
        """生成组织数据"""
        
    @staticmethod
    def generate_targets(count: int, target_type_ratios: dict) -> list[dict]:
        """生成目标数据（域名70%、IP20%、CIDR10%）"""
        
    @staticmethod
    def generate_websites(target: dict, count: int) -> list[dict]:
        """为目标生成 Website 数据"""
        
    @staticmethod
    def generate_subdomains(target: dict, count: int) -> list[dict]:
        """为域名目标生成 Subdomain 数据"""
        
    @staticmethod
    def generate_endpoints(target: dict, count: int) -> list[dict]:
        """为目标生成 Endpoint 数据"""
        
    @staticmethod
    def generate_directories(target: dict, count: int) -> list[dict]:
        """为目标生成 Directory 数据"""
        
    @staticmethod
    def generate_host_ports(target: dict, count: int) -> list[dict]:
        """为目标生成 HostPort 数据"""
        
    @staticmethod
    def generate_vulnerabilities(target: dict, count: int) -> list[dict]:
        """为目标生成 Vulnerability 数据"""
```

**数据模板：**

- 组织名称：从预定义列表中选择 + 随机后缀
- 域名：`{env}.{company}-{suffix}.{tld}` 格式
- IP：随机生成合法的 IPv4 地址
- CIDR：随机生成 /8、/16、/24 网段
- URL：根据目标类型生成合理的 URL
- 技术栈：从常见技术中随机选择

### 3. Progress Tracker

**职责：** 显示生成进度和统计信息

**接口：**

```python
class ProgressTracker:
    def __init__(self):
        """初始化进度跟踪器"""
        
    def start_phase(self, phase_name: str, total: int):
        """开始新阶段"""
        
    def update(self, count: int):
        """更新进度"""
        
    def add_success(self, count: int):
        """记录成功数量"""
        
    def add_error(self, error: str):
        """记录错误"""
        
    def finish_phase(self):
        """完成当前阶段"""
        
    def print_summary(self):
        """打印总结"""
```

**显示格式：**

```
🏢 Creating organizations... [15/15] ✓ 15 created
🎯 Creating targets... [225/225] ✓ 225 created (domains: 157, IPs: 45, CIDRs: 23)
🔗 Linking targets to organizations... [225/225] ✓ 225 links created
🌐 Creating websites... [3375/3375] ✓ 3375 created
📝 Creating subdomains... [2355/2355] ✓ 2355 created (157 domain targets)
...

✅ Test data generation completed!
   Total time: 45.2s
   Success: 12,000 records
   Errors: 3 records
```

### 4. Error Handler

**职责：** 处理 API 错误和重试逻辑

**接口：**

```python
class ErrorHandler:
    def __init__(self, max_retries: int = 3, retry_delay: float = 1.0):
        """初始化错误处理器"""
        
    def should_retry(self, status_code: int) -> bool:
        """判断是否应该重试"""
        
    def handle_error(self, error: Exception, context: dict) -> bool:
        """处理错误，返回是否应该继续"""
        
    def log_error(self, error: str, request_data: dict = None, response_data: dict = None):
        """记录错误详情"""
```

**重试策略：**

| 状态码 | 行为 | 重试次数 |
|--------|------|----------|
| 5xx | 自动重试 | 3 次 |
| 429 | 等待后重试 | 3 次 |
| 401 | 刷新 token 后重试 | 1 次 |
| 4xx (其他) | 记录错误，跳过 | 0 次 |
| 网络超时 | 重试 | 3 次 |

## Data Models

### JSON 请求格式

所有 JSON 使用 **camelCase** 字段名（符合前端规范）：

**创建组织：**
```python
{
    "name": "Acme Corporation - Global (5123-0)",
    "description": "A leading technology company..."
}
```

**批量创建目标：**
```python
{
    "targets": [
        {"name": "example.com", "type": "domain"},
        {"name": "192.168.1.1", "type": "ip"},
        {"name": "10.0.0.0/8", "type": "cidr"}
    ]
}
```

**关联目标到组织：**
```python
{
    "targetIds": [1, 2, 3, 4, 5]
}
```

**批量创建 Website：**
```python
{
    "websites": [
        {
            "url": "https://www.example.com",
            "title": "Welcome - Dashboard",
            "statusCode": 200,
            "contentLength": 1500,
            "contentType": "text/html; charset=utf-8",
            "webserver": "nginx/1.24.0",
            "tech": ["nginx", "PHP", "MySQL"],
            "vhost": false
        }
    ]
}
```

**批量创建 Subdomain：**
```python
{
    "subdomains": [
        {"name": "www.example.com"},
        {"name": "api.example.com"}
    ]
}
```

**批量创建 Endpoint：**
```python
{
    "endpoints": [
        {
            "url": "https://api.example.com/v1/users",
            "title": "User Service",
            "statusCode": 200,
            "contentLength": 500,
            "contentType": "application/json",
            "webserver": "nginx/1.24.0",
            "tech": ["nginx", "Node.js", "Express"],
            "matchedGfPatterns": ["cors", "ssrf"],
            "vhost": false
        }
    ]
}
```

**批量创建 Directory：**
```python
{
    "directories": [
        {
            "url": "https://www.example.com/admin/",
            "status": 403,
            "contentLength": 1200,
            "contentType": "text/html",
            "duration": 55
        }
    ]
}
```

**批量创建 HostPort：**
```python
{
    "hostPorts": [
        {
            "host": "www.example.com",
            "ip": "192.168.1.10",
            "port": 443
        }
    ]
}
```

**批量创建 Vulnerability：**
```python
{
    "vulnerabilities": [
        {
            "url": "https://www.example.com/login",
            "vulnType": "SQL Injection",
            "severity": "critical",
            "source": "nuclei",
            "cvssScore": 9.8,
            "description": "A SQL injection vulnerability was found..."
        }
    ]
}
```

## Correctness Properties

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的正式陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 认证 Token 有效性

*对于任何* API 请求，如果 token 有效，则请求应该成功；如果 token 过期，则应该自动刷新后重试成功

**Validates: Requirements 1.2, 1.3**

### Property 2: 批量创建幂等性

*对于任何* 批量创建请求，重复发送相同的数据应该不会创建重复记录（由于 API 的 ON CONFLICT DO NOTHING）

**Validates: Requirements 3.4, 5.4, 6.4, 7.4, 8.4, 9.4, 10.4**

### Property 3: 目标类型分布

*对于任何* 目标生成请求，生成的目标类型分布应该接近指定的比例（域名70%、IP20%、CIDR10%，误差±5%）

**Validates: Requirements 3.2**

### Property 4: 资产归属验证

*对于任何* 资产（Website、Subdomain、Endpoint 等），其 URL/名称应该与所属目标匹配（域名匹配、IP 匹配、CIDR 范围内）

**Validates: Requirements 5.2, 6.2, 7.2, 8.2, 9.2**

### Property 5: JSON 字段命名一致性

*对于任何* API 响应，所有字段名应该使用 camelCase 格式，不应该出现 snake_case

**Validates: Requirements 15.4**

### Property 6: 错误重试收敛性

*对于任何* 5xx 错误或网络超时，重试次数应该不超过 3 次，且最终要么成功要么记录失败

**Validates: Requirements 12.1, 12.3**

### Property 7: 进度显示单调性

*对于任何* 生成阶段，显示的进度数字应该单调递增，且最终等于总数

**Validates: Requirements 13.1, 13.2**

### Property 8: 批量操作分批一致性

*对于任何* 批量操作，如果总数超过批次大小，则应该分批发送，且所有批次的总和等于原始总数

**Validates: Requirements 3.4, 5.4, 6.4, 7.4, 8.4, 9.4, 10.4**

### Property 9: 数据清理顺序正确性

*对于任何* 清理操作，删除顺序应该遵循外键约束（先删除子表，再删除父表），不应该出现外键冲突错误

**Validates: Requirements 14.2**

### Property 10: 组织目标分配均匀性

*对于任何* 目标关联操作，每个组织分配的目标数量应该大致相等（误差不超过 1）

**Validates: Requirements 4.2**

## Error Handling

### 错误分类

| 错误类型 | HTTP 状态码 | 处理策略 |
|----------|-------------|----------|
| 认证失败 | 401 | 刷新 token 后重试 1 次 |
| 权限不足 | 403 | 记录错误，终止程序 |
| 资源不存在 | 404 | 记录错误，跳过该记录 |
| 请求格式错误 | 400 | 记录详细错误（包含请求 JSON），跳过该记录 |
| 资源冲突 | 409 | 记录警告，跳过该记录（可能是重复数据） |
| 限流 | 429 | 等待 5 秒后重试，最多 3 次 |
| 服务器错误 | 5xx | 等待 1 秒后重试，最多 3 次 |
| 网络超时 | Timeout | 等待 1 秒后重试，最多 3 次 |
| 连接失败 | ConnectionError | 等待 2 秒后重试，最多 3 次 |

### 错误日志格式

```python
{
    "timestamp": "2026-01-14T10:30:45Z",
    "error_type": "API_ERROR",
    "status_code": 400,
    "endpoint": "/api/targets/1/websites/bulk-create",
    "request": {
        "websites": [...]
    },
    "response": {
        "error": {
            "code": "VALIDATION_ERROR",
            "message": "Invalid URL format"
        }
    },
    "retry_count": 0
}
```

### 错误恢复

- **部分失败策略：** 批量操作中，单条记录失败不影响其他记录
- **断点续传：** 记录已成功创建的数据 ID，失败后可以从断点继续
- **回滚机制：** 提供 `--clear` 参数清空所有数据，重新开始

## Testing Strategy

### 单元测试

**测试范围：**
- Data Generator 的数据生成逻辑
- API Client 的请求构造
- Error Handler 的重试逻辑
- Progress Tracker 的统计计算

**测试工具：** `pytest`

**示例测试：**

```python
def test_generate_domain_target():
    """测试域名目标生成"""
    target = DataGenerator.generate_target("domain", 0)
    assert target["type"] == "domain"
    assert "." in target["name"]
    assert not target["name"].startswith(".")

def test_api_client_auto_refresh_token(mock_api):
    """测试 token 自动刷新"""
    client = APIClient("http://localhost:8888", "admin", "admin")
    # 模拟 token 过期
    mock_api.set_token_expired()
    # 应该自动刷新 token 并重试
    response = client.post("/api/targets", {"name": "test.com", "type": "domain"})
    assert response["id"] > 0
    assert mock_api.refresh_called

def test_error_handler_retry_on_5xx():
    """测试 5xx 错误重试"""
    handler = ErrorHandler(max_retries=3)
    assert handler.should_retry(500) == True
    assert handler.should_retry(503) == True
    assert handler.should_retry(400) == False
```

### 集成测试

**测试范围：**
- 完整的数据生成流程
- API 调用的正确性
- 错误处理和重试
- 进度显示

**测试环境：** 本地 Go 后端（端口 8888）

**测试步骤：**

1. 启动 Go 后端
2. 运行种子生成器（小规模：2 个组织，10 个目标）
3. 验证数据是否正确创建
4. 验证 JSON 字段命名（camelCase）
5. 验证错误处理（模拟网络错误）

### 手动测试

**测试场景：**

| 场景 | 命令 | 预期结果 |
|------|------|----------|
| 小规模生成 | `python seed_generator.py --orgs 2` | 快速完成，数据正确 |
| 大规模生成 | `python seed_generator.py --orgs 50` | 进度显示正常，无内存问题 |
| 清空数据 | `python seed_generator.py --clear` | 所有数据被删除 |
| 网络中断 | 生成过程中断开网络 | 自动重试，显示错误 |
| 认证失败 | 使用错误的密码 | 显示认证错误，退出 |
| API 不可用 | 后端未启动 | 显示连接错误，退出 |

## Implementation Notes

### 项目位置

脚本放在 `tools/seed-api/` 目录下，与 Go 后端分离：

```
项目根目录/
├── go-backend/           # Go 后端
├── backend/              # Python 后端（旧）
├── frontend/             # 前端
└── tools/                # 工具脚本
    └── seed-api/         # API 种子数据生成器 ⭐
        ├── seed_generator.py      # 主程序入口
        ├── api_client.py          # API 客户端
        ├── data_generator.py      # 数据生成器
        ├── progress.py            # 进度跟踪
        ├── error_handler.py       # 错误处理
        ├── requirements.txt       # Python 依赖
        └── README.md              # 使用说明
```

**选择 `tools/seed-api/` 的原因：**
1. ✅ 独立于后端代码，不会被误认为是后端的一部分
2. ✅ 与其他工具脚本放在一起，便于管理
3. ✅ 可以独立运行，不依赖 Go 后端的构建
4. ✅ 便于版本控制和分发

### 代码组织

**模块化设计 - 每个文件一个职责：**

| 文件 | 行数估计 | 职责 |
|------|----------|------|
| `seed_generator.py` | ~150 行 | 主程序：CLI 参数解析、流程编排 |
| `api_client.py` | ~200 行 | API 客户端：HTTP 请求、认证管理 |
| `data_generator.py` | ~400 行 | 数据生成：生成各类测试数据 |
| `progress.py` | ~100 行 | 进度跟踪：显示进度和统计 |
| `error_handler.py` | ~150 行 | 错误处理：重试逻辑、错误日志 |

**总计：** ~1000 行代码，分 5 个文件

**不要写成单文件的原因：**
- ❌ 单文件 1000+ 行难以维护
- ❌ 职责不清晰，难以测试
- ❌ 难以复用（比如 API Client 可以用于其他脚本）
- ✅ 模块化便于单元测试
- ✅ 每个文件可以独立理解和修改

### 命令行参数

```bash
python seed_generator.py [OPTIONS]

Options:
  --api-url URL          API 地址 (默认: http://localhost:8888)
  --username USER        用户名 (默认: admin)
  --password PASS        密码 (默认: admin)
  --orgs N               组织数量 (默认: 15)
  --targets-per-org N    每个组织的目标数量 (默认: 15)
  --assets-per-target N  每个目标的资产数量 (默认: 15)
  --clear                清空现有数据
  --batch-size N         批量操作的批次大小 (默认: 100)
  --verbose              显示详细日志
  --help                 显示帮助信息
```

### 性能优化

1. **批量操作：** 使用批量 API（bulk-create、bulk-upsert）减少请求次数
2. **连接复用：** 使用 `requests.Session` 复用 HTTP 连接
3. **并发控制：** 单线程顺序执行（避免复杂性，性能已足够）
4. **内存优化：** 分批生成数据，避免一次性加载所有数据到内存

### 依赖管理

**requirements.txt:**
```
requests>=2.31.0
```

**安装命令：**
```bash
pip install -r requirements.txt
```

### 使用示例

```bash
# 1. 进入工具目录
cd tools/seed-api

# 2. 安装依赖
pip install -r requirements.txt

# 3. 启动 Go 后端（另一个终端）
cd ../../go-backend
make run

# 4. 生成测试数据（默认配置）
python seed_generator.py

# 5. 生成大规模测试数据
python seed_generator.py --orgs 50 --targets-per-org 20

# 6. 清空数据后重新生成
python seed_generator.py --clear --orgs 10

# 7. 使用自定义 API 地址
python seed_generator.py --api-url http://192.168.1.100:8888
```

### 模块导入关系

```python
# seed_generator.py (主程序)
from api_client import APIClient
from data_generator import DataGenerator
from progress import ProgressTracker
from error_handler import ErrorHandler

# api_client.py (独立模块)
import requests

# data_generator.py (独立模块)
import random
from typing import List, Dict

# progress.py (独立模块)
from datetime import datetime

# error_handler.py (独立模块)
import json
from typing import Optional
```

**模块间依赖：**
- `seed_generator.py` 依赖所有其他模块
- 其他模块相互独立，便于测试和复用
