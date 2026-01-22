# Requirements Document

## Introduction

为 Go 后端实现 host_port_mapping 资产的 CRUD API，用于管理目标下的主机-IP-端口映射关系。该 API 主要服务于端口扫描结果（如 naabu）的存储和查询。

## Glossary

- **Host_Port_Mapping_Service**: 处理主机端口映射业务逻辑的服务层组件
- **Host_Port_Mapping_Repository**: 处理主机端口映射数据库操作的数据访问层组件
- **Host_Port_Mapping_Handler**: 处理主机端口映射 HTTP 请求的控制器组件
- **Target**: 扫描目标，可以是域名、IP 或 CIDR
- **Mapping**: 一条 host-ip-port 的映射记录

## Requirements

### Requirement 1: List Host Port Mappings

**User Story:** As a user, I want to list all host-port mappings for a target with pagination and filtering, so that I can view and search the port scan results.

#### Acceptance Criteria

1. WHEN a user requests the list endpoint with a valid target ID, THE Host_Port_Mapping_Handler SHALL return a paginated list of mappings
2. WHEN a user provides filter parameters, THE Host_Port_Mapping_Service SHALL filter results by host, ip, or port
3. WHEN the target does not exist, THE Host_Port_Mapping_Handler SHALL return a 404 error
4. THE Host_Port_Mapping_Repository SHALL support pagination with page and pageSize parameters

### Requirement 2: Export Host Port Mappings

**User Story:** As a user, I want to export all host-port mappings for a target as CSV, so that I can analyze the data externally.

#### Acceptance Criteria

1. WHEN a user requests the export endpoint, THE Host_Port_Mapping_Handler SHALL return a CSV file with all mappings
2. THE Host_Port_Mapping_Service SHALL stream the CSV data to avoid memory issues with large datasets
3. WHEN the target does not exist, THE Host_Port_Mapping_Handler SHALL return a 404 error

### Requirement 3: Bulk Upsert Host Port Mappings

**User Story:** As a scanner, I want to import port scan results in bulk, so that the database stays up-to-date with the latest scan data.

#### Acceptance Criteria

1. WHEN a scanner submits mapping data, THE Host_Port_Mapping_Service SHALL create or update records
2. THE Host_Port_Mapping_Repository SHALL use ON CONFLICT DO NOTHING to handle duplicates (since all fields are in unique constraint)
3. WHEN the target does not exist, THE Host_Port_Mapping_Handler SHALL return a 404 error
4. THE Host_Port_Mapping_Handler SHALL return the count of upserted records
5. THE Host_Port_Mapping_Repository SHALL batch operations to avoid PostgreSQL parameter limits

### Requirement 4: Bulk Delete Host Port Mappings

**User Story:** As a user, I want to delete multiple host-port mappings at once by IP address, so that I can clean up outdated or incorrect data.

#### Acceptance Criteria

1. WHEN a user submits a list of IP addresses, THE Host_Port_Mapping_Service SHALL delete all mapping records for those IPs
2. THE Host_Port_Mapping_Handler SHALL return the count of deleted records
3. IF some IPs do not have any mappings, THE Host_Port_Mapping_Service SHALL delete only the existing ones without error
