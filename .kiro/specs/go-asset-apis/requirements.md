# Requirements Document

## Introduction

本文档定义了 Go 后端资产 API 的需求，包括 Subdomain（子域名）、Endpoint（URL 端点）和 Directory（目录）三类资产的 CRUD 操作。这些 API 将替代现有的 Python Django 后端实现，保持与前端的兼容性。

## Glossary

- **Asset_API**: Go 后端资产管理 API 系统
- **Subdomain**: 子域名资产，关联到特定 Target
- **Endpoint**: URL 端点资产，包含 HTTP 响应元数据
- **Directory**: 目录扫描结果资产，包含 HTTP 状态和内容信息
- **Target**: 扫描目标，是所有资产的父级实体
- **Filter**: 智能过滤查询字符串，支持字段搜索和纯文本搜索
- **Bulk_Operation**: 批量操作，支持批量创建和批量删除

## Requirements

### Requirement 1: Subdomain 列表查询

**User Story:** As a security analyst, I want to list subdomains for a target, so that I can review discovered subdomains.

#### Acceptance Criteria

1. WHEN a user requests subdomains for a target, THE Asset_API SHALL return a paginated list of subdomains belonging to that target
2. WHEN pagination parameters are provided, THE Asset_API SHALL return the specified page with the specified page size
3. WHEN a filter parameter is provided, THE Asset_API SHALL filter subdomains by name containing the filter text
4. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error with message "Target not found"
5. THE Asset_API SHALL return subdomains sorted by creation time in descending order

### Requirement 2: Subdomain 批量创建

**User Story:** As a security analyst, I want to bulk create subdomains for a target, so that I can import discovered subdomains efficiently.

#### Acceptance Criteria

1. WHEN a user submits a list of subdomain names, THE Asset_API SHALL create subdomains that don't already exist
2. WHEN duplicate subdomain names are submitted, THE Asset_API SHALL skip duplicates and continue processing
3. THE Asset_API SHALL return the count of successfully created subdomains
4. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error
5. WHEN subdomain names are invalid (empty or whitespace only), THE Asset_API SHALL skip them silently
6. IF the target type is not "domain", THEN THE Asset_API SHALL return a 400 error with message "Invalid target type"
7. WHEN a subdomain does not match the target domain (not equal to or ending with .target), THE Asset_API SHALL skip it silently

### Requirement 3: Subdomain 批量删除

**User Story:** As a security analyst, I want to bulk delete subdomains, so that I can clean up unwanted data efficiently.

#### Acceptance Criteria

1. WHEN a user submits a list of subdomain IDs, THE Asset_API SHALL delete all specified subdomains
2. THE Asset_API SHALL return the count of successfully deleted subdomains
3. WHEN some IDs don't exist, THE Asset_API SHALL delete existing ones and ignore non-existent IDs

### Requirement 4: Subdomain 导出

**User Story:** As a security analyst, I want to export subdomains as a text file, so that I can use them with external tools.

#### Acceptance Criteria

1. WHEN a user requests subdomain export for a target, THE Asset_API SHALL return a text file with one subdomain per line
2. THE Asset_API SHALL set appropriate Content-Type and Content-Disposition headers for file download
3. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error

### Requirement 5: Endpoint 列表查询

**User Story:** As a security analyst, I want to list endpoints for a target, so that I can review discovered URLs.

#### Acceptance Criteria

1. WHEN a user requests endpoints for a target, THE Asset_API SHALL return a paginated list of endpoints belonging to that target
2. WHEN pagination parameters are provided, THE Asset_API SHALL return the specified page with the specified page size
3. WHEN a filter parameter is provided, THE Asset_API SHALL filter endpoints by URL, host, or title containing the filter text
4. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error with message "Target not found"
5. THE Asset_API SHALL return endpoints sorted by creation time in descending order

### Requirement 6: Endpoint 详情查询

**User Story:** As a security analyst, I want to view endpoint details, so that I can analyze HTTP response information.

#### Acceptance Criteria

1. WHEN a user requests an endpoint by ID, THE Asset_API SHALL return the complete endpoint details
2. THE Asset_API SHALL include all HTTP metadata fields (statusCode, contentLength, contentType, tech, etc.)
3. IF the endpoint does not exist, THEN THE Asset_API SHALL return a 404 error with message "Endpoint not found"

### Requirement 7: Endpoint 批量创建

**User Story:** As a security analyst, I want to bulk create endpoints for a target, so that I can import discovered URLs efficiently.

#### Acceptance Criteria

1. WHEN a user submits a list of URLs, THE Asset_API SHALL create endpoints that don't already exist
2. WHEN duplicate URLs are submitted, THE Asset_API SHALL skip duplicates and continue processing
3. THE Asset_API SHALL return the count of successfully created endpoints
4. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error
5. WHEN URLs are invalid (empty or whitespace only), THE Asset_API SHALL skip them silently
6. WHEN a URL hostname does not match the target (domain suffix, IP equality, or CIDR range), THE Asset_API SHALL skip it silently

### Requirement 8: Endpoint 单个删除

**User Story:** As a security analyst, I want to delete a single endpoint, so that I can remove unwanted data.

#### Acceptance Criteria

1. WHEN a user requests to delete an endpoint by ID, THE Asset_API SHALL delete the endpoint
2. THE Asset_API SHALL return 204 No Content on successful deletion
3. IF the endpoint does not exist, THEN THE Asset_API SHALL return a 404 error

### Requirement 9: Endpoint 批量删除

**User Story:** As a security analyst, I want to bulk delete endpoints, so that I can clean up unwanted data efficiently.

#### Acceptance Criteria

1. WHEN a user submits a list of endpoint IDs, THE Asset_API SHALL delete all specified endpoints
2. THE Asset_API SHALL return the count of successfully deleted endpoints
3. WHEN some IDs don't exist, THE Asset_API SHALL delete existing ones and ignore non-existent IDs

### Requirement 10: Endpoint 导出

**User Story:** As a security analyst, I want to export endpoints as a text file, so that I can use them with external tools.

#### Acceptance Criteria

1. WHEN a user requests endpoint export for a target, THE Asset_API SHALL return a text file with one URL per line
2. THE Asset_API SHALL set appropriate Content-Type and Content-Disposition headers for file download
3. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error

### Requirement 11: Directory 列表查询

**User Story:** As a security analyst, I want to list directories for a target, so that I can review discovered paths.

#### Acceptance Criteria

1. WHEN a user requests directories for a target, THE Asset_API SHALL return a paginated list of directories belonging to that target
2. WHEN pagination parameters are provided, THE Asset_API SHALL return the specified page with the specified page size
3. WHEN a filter parameter is provided, THE Asset_API SHALL filter directories by URL containing the filter text
4. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error with message "Target not found"
5. THE Asset_API SHALL return directories sorted by creation time in descending order

### Requirement 12: Directory 批量创建

**User Story:** As a security analyst, I want to bulk create directories for a target, so that I can import discovered paths efficiently.

#### Acceptance Criteria

1. WHEN a user submits a list of directory URLs, THE Asset_API SHALL create directories that don't already exist
2. WHEN duplicate URLs are submitted, THE Asset_API SHALL skip duplicates and continue processing
3. THE Asset_API SHALL return the count of successfully created directories
4. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error
5. WHEN URLs are invalid (empty or whitespace only), THE Asset_API SHALL skip them silently
6. WHEN a URL hostname does not match the target (domain suffix, IP equality, or CIDR range), THE Asset_API SHALL skip it silently

### Requirement 13: Directory 批量删除

**User Story:** As a security analyst, I want to bulk delete directories, so that I can clean up unwanted data efficiently.

#### Acceptance Criteria

1. WHEN a user submits a list of directory IDs, THE Asset_API SHALL delete all specified directories
2. THE Asset_API SHALL return the count of successfully deleted directories
3. WHEN some IDs don't exist, THE Asset_API SHALL delete existing ones and ignore non-existent IDs

### Requirement 14: Directory 导出

**User Story:** As a security analyst, I want to export directories as a text file, so that I can use them with external tools.

#### Acceptance Criteria

1. WHEN a user requests directory export for a target, THE Asset_API SHALL return a text file with one URL per line
2. THE Asset_API SHALL set appropriate Content-Type and Content-Disposition headers for file download
3. IF the target does not exist, THEN THE Asset_API SHALL return a 404 error

### Requirement 15: 分页响应格式一致性

**User Story:** As a frontend developer, I want consistent pagination response format, so that I can reuse pagination components.

#### Acceptance Criteria

1. THE Asset_API SHALL return pagination responses with fields: results, total, page, pageSize, totalPages
2. THE Asset_API SHALL use camelCase for JSON field names
3. WHEN results are empty, THE Asset_API SHALL return an empty array instead of null
4. THE Asset_API SHALL calculate totalPages correctly based on total and pageSize
