# Requirements Document

## Introduction

本文档定义了基于 API 的种子数据生成器的需求。该生成器将通过调用 Go 后端的 REST API 来创建测试数据，而不是直接操作数据库。这种方式能够测试完整的 API 流程，包括路由、中间件、验证、序列化等，更接近真实的生产环境。

**实现语言：Python**

选择 Python 的原因：
1. 最适合脚本任务，代码简洁
2. JSON 操作自然（Python 字典 = JSON）
3. 完全独立于 Go 后端代码，真正测试 API
4. requests 库成熟，错误处理简单
5. 调试方便，易于维护

## Glossary

- **Seed_Generator**: 种子数据生成器，用于创建测试数据的 Python 脚本
- **API_Client**: HTTP 客户端，使用 Python requests 库发送 API 请求
- **JWT_Token**: JSON Web Token，用于 API 认证
- **Batch_API**: 批量操作 API，支持一次创建多条记录
- **Test_Data**: 测试数据，包括组织、目标、资产等
- **Asset**: 资产，包括 Website、Subdomain、Endpoint、Directory、HostPort 等
- **Independent_JSON**: 独立构造的 JSON，不依赖后端 DTO 结构体

## Requirements

### Requirement 1: 认证管理

**User Story:** 作为种子数据生成器，我需要通过 API 进行身份认证，以便访问受保护的 API 端点。

#### Acceptance Criteria

1. WHEN 生成器启动时，THE Seed_Generator SHALL 调用 `/api/auth/login` 获取 JWT token
2. WHEN token 获取成功后，THE Seed_Generator SHALL 在所有后续请求中携带 Authorization header
3. IF token 过期，THEN THE Seed_Generator SHALL 自动调用 `/api/auth/refresh` 刷新 token
4. WHEN 认证失败时，THE Seed_Generator SHALL 返回明确的错误信息并退出

### Requirement 2: 组织数据生成

**User Story:** 作为测试人员，我需要生成多个组织数据，以便测试多租户场景。

#### Acceptance Criteria

1. WHEN 用户指定组织数量时，THE Seed_Generator SHALL 调用 `/api/organizations` POST 接口创建组织
2. WHEN 创建组织时，THE Seed_Generator SHALL 生成随机但合理的组织名称和描述
3. WHEN 组织创建成功后，THE Seed_Generator SHALL 保存组织 ID 用于后续关联
4. WHEN API 返回错误时，THE Seed_Generator SHALL 记录错误详情并继续处理其他数据

### Requirement 3: 目标数据生成

**User Story:** 作为测试人员，我需要生成不同类型的目标（域名、IP、CIDR），以便测试各种扫描场景。

#### Acceptance Criteria

1. WHEN 用户指定目标数量时，THE Seed_Generator SHALL 调用 `/api/targets/batch_create` 批量创建目标
2. WHEN 生成目标时，THE Seed_Generator SHALL 按比例生成域名（70%）、IP（20%）、CIDR（10%）
3. WHEN 目标创建成功后，THE Seed_Generator SHALL 保存目标 ID 用于后续资产创建
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 4: 目标与组织关联

**User Story:** 作为测试人员，我需要将目标关联到组织，以便测试组织的目标管理功能。

#### Acceptance Criteria

1. WHEN 目标和组织都创建完成后，THE Seed_Generator SHALL 调用 `/api/organizations/:id/link_targets` 关联目标
2. WHEN 关联目标时，THE Seed_Generator SHALL 平均分配目标到各个组织
3. WHEN 关联失败时，THE Seed_Generator SHALL 记录失败的组织和目标 ID
4. WHEN 批量关联时，THE Seed_Generator SHALL 每批次不超过 50 个目标

### Requirement 5: Website 资产生成

**User Story:** 作为测试人员，我需要为每个目标生成 Website 资产，以便测试 Website 列表和导出功能。

#### Acceptance Criteria

1. WHEN 目标创建完成后，THE Seed_Generator SHALL 调用 `/api/targets/:id/websites/bulk-upsert` 创建 Website
2. WHEN 生成 Website 时，THE Seed_Generator SHALL 根据目标类型生成合理的 URL
3. WHEN 生成 Website 时，THE Seed_Generator SHALL 包含 title、statusCode、tech 等字段
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 6: Subdomain 资产生成

**User Story:** 作为测试人员，我需要为域名类型的目标生成 Subdomain 资产，以便测试子域名发现功能。

#### Acceptance Criteria

1. WHEN 目标类型为 domain 时，THE Seed_Generator SHALL 调用 `/api/targets/:id/subdomains/bulk-create` 创建 Subdomain
2. WHEN 目标类型不是 domain 时，THE Seed_Generator SHALL 跳过该目标的 Subdomain 生成
3. WHEN 生成 Subdomain 时，THE Seed_Generator SHALL 使用常见的子域名前缀（www、api、admin 等）
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 7: Endpoint 资产生成

**User Story:** 作为测试人员，我需要为每个目标生成 Endpoint 资产，以便测试端点发现和分析功能。

#### Acceptance Criteria

1. WHEN 目标创建完成后，THE Seed_Generator SHALL 调用 `/api/targets/:id/endpoints/bulk-upsert` 创建 Endpoint
2. WHEN 生成 Endpoint 时，THE Seed_Generator SHALL 包含 URL、statusCode、tech、matchedGFPatterns 等字段
3. WHEN 生成 Endpoint 时，THE Seed_Generator SHALL 使用常见的 API 路径（/api/v1/users、/login 等）
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 8: Directory 资产生成

**User Story:** 作为测试人员，我需要为每个目标生成 Directory 资产，以便测试目录扫描功能。

#### Acceptance Criteria

1. WHEN 目标创建完成后，THE Seed_Generator SHALL 调用 `/api/targets/:id/directories/bulk-upsert` 创建 Directory
2. WHEN 生成 Directory 时，THE Seed_Generator SHALL 包含 URL、status、contentLength 等字段
3. WHEN 生成 Directory 时，THE Seed_Generator SHALL 使用常见的目录路径（/admin/、/backup/ 等）
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 9: HostPort 资产生成

**User Story:** 作为测试人员，我需要为每个目标生成 HostPort 映射，以便测试端口扫描功能。

#### Acceptance Criteria

1. WHEN 目标创建完成后，THE Seed_Generator SHALL 调用 `/api/targets/:id/host-ports/bulk-upsert` 创建 HostPort
2. WHEN 生成 HostPort 时，THE Seed_Generator SHALL 包含 host、ip、port 字段
3. WHEN 生成 HostPort 时，THE Seed_Generator SHALL 使用常见的端口（80、443、22、3306 等）
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 10: Vulnerability 数据生成

**User Story:** 作为测试人员，我需要为每个目标生成漏洞数据，以便测试漏洞管理功能。

#### Acceptance Criteria

1. WHEN 目标创建完成后，THE Seed_Generator SHALL 调用 `/api/targets/:id/vulnerabilities/bulk-create` 创建漏洞
2. WHEN 生成漏洞时，THE Seed_Generator SHALL 包含 vulnType、severity、cvssScore 等字段
3. WHEN 生成漏洞时，THE Seed_Generator SHALL 按比例生成不同严重级别的漏洞
4. WHEN 批量创建时，THE Seed_Generator SHALL 每批次不超过 100 条记录

### Requirement 11: 命令行参数

**User Story:** 作为测试人员，我需要通过命令行参数控制生成的数据量，以便灵活调整测试规模。

#### Acceptance Criteria

1. THE Seed_Generator SHALL 支持 `-orgs` 参数指定组织数量
2. THE Seed_Generator SHALL 支持 `-targets-per-org` 参数指定每个组织的目标数量
3. THE Seed_Generator SHALL 支持 `-assets-per-target` 参数指定每个目标的资产数量
4. THE Seed_Generator SHALL 支持 `-clear` 参数清空现有数据
5. THE Seed_Generator SHALL 支持 `-api-url` 参数指定 API 地址
6. THE Seed_Generator SHALL 支持 `-username` 和 `-password` 参数指定登录凭据

### Requirement 12: 错误处理和重试

**User Story:** 作为测试人员，我需要生成器能够处理 API 错误并重试，以便在网络不稳定时也能完成数据生成。

#### Acceptance Criteria

1. WHEN API 返回 5xx 错误时，THE Seed_Generator SHALL 自动重试最多 3 次
2. WHEN API 返回 4xx 错误时，THE Seed_Generator SHALL 记录错误并跳过该条记录
3. WHEN 重试次数耗尽时，THE Seed_Generator SHALL 记录失败详情并继续处理其他数据
4. WHEN 遇到网络超时时，THE Seed_Generator SHALL 等待 1 秒后重试

### Requirement 13: 进度显示

**User Story:** 作为测试人员，我需要看到数据生成的进度，以便了解生成过程的状态。

#### Acceptance Criteria

1. WHEN 生成器运行时，THE Seed_Generator SHALL 显示当前正在处理的数据类型
2. WHEN 每个阶段完成时，THE Seed_Generator SHALL 显示成功创建的记录数量
3. WHEN 发生错误时，THE Seed_Generator SHALL 显示错误数量和详情
4. WHEN 生成完成时，THE Seed_Generator SHALL 显示总耗时和统计信息

### Requirement 14: 数据清理

**User Story:** 作为测试人员，我需要在生成新数据前清空旧数据，以便从干净的状态开始测试。

#### Acceptance Criteria

1. WHEN 用户指定 `-clear` 参数时，THE Seed_Generator SHALL 调用批量删除 API 清空数据
2. WHEN 清理数据时，THE Seed_Generator SHALL 按正确的顺序删除（先删除资产，再删除目标，最后删除组织）
3. WHEN 清理失败时，THE Seed_Generator SHALL 显示错误信息并询问是否继续
4. WHEN 清理完成时，THE Seed_Generator SHALL 显示清理的记录数量

### Requirement 15: 独立 JSON 构造

**User Story:** 作为开发人员，我需要种子生成器独立构造 JSON 请求，而不是复用后端 DTO 结构体，以便发现 API 序列化问题。

#### Acceptance Criteria

1. THE Seed_Generator SHALL 使用 Python 字典独立构造请求 JSON
2. THE Seed_Generator SHALL NOT 导入或依赖 Go 后端的任何代码
3. WHEN API 返回错误时，THE Seed_Generator SHALL 显示原始 JSON 请求和响应
4. THE Seed_Generator SHALL 验证响应字段是否符合预期的 camelCase 格式
5. WHEN 字段命名不一致时，THE Seed_Generator SHALL 记录详细的错误信息

### Requirement 16: Python 环境依赖

**User Story:** 作为测试人员，我需要简单的环境配置，以便快速运行种子生成器。

#### Acceptance Criteria

1. THE Seed_Generator SHALL 使用 Python 3.8+ 版本
2. THE Seed_Generator SHALL 仅依赖 `requests` 库（通过 pip 安装）
3. THE Seed_Generator SHALL 提供 `requirements.txt` 文件列出依赖
4. THE Seed_Generator SHALL 在脚本开头检查依赖是否已安装
