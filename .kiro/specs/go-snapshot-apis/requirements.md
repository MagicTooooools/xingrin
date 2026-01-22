# 需求文档

## 简介

本文档定义了 Go 后端网站快照（WebsiteSnapshot）API 的需求。这是快照 API 系列的第一个实现，将作为其他快照类型（Subdomain、Endpoint、Directory 等）的参考模板。

网站快照记录了每次扫描时发现的网站及其响应信息，作为扫描（Scan）的嵌套资源提供访问。

**核心设计原则**：扫描结果通过快照 API 写入，内部自动同步到资产表，Worker 只需调用一个接口。

## 术语表

- **Snapshot_API**: 快照 API 服务，提供对扫描快照数据的访问和写入接口
- **Scan**: 扫描记录，快照的父资源
- **WebsiteSnapshot**: 网站快照，记录扫描发现的网站及其响应信息（URL、标题、状态码、技术栈等）
- **Website**: 网站资产，从快照同步而来的去重资产记录
- **Filter_Query**: 过滤查询字符串，支持字段级别的模糊匹配
- **CSV_Export**: CSV 格式导出功能，支持流式输出
- **Save_And_Sync**: 保存并同步操作，同时写入快照表和资产表

## 需求

### 需求 1: 网站快照批量写入（扫描结果导入）

**用户故事:** 作为扫描 Worker，我希望通过一个接口提交扫描发现的网站，系统自动保存快照并同步到资产表。

#### 验收标准

1. WHEN Worker 请求 POST /api/scans/{scan_id}/websites/bulk-upsert THEN Snapshot_API SHALL 保存网站快照到 website_snapshot 表
2. WHEN 快照保存成功后 THEN Snapshot_API SHALL 自动同步数据到 website 资产表（upsert 模式）
3. WHEN scan_id 不存在或已被软删除 THEN Snapshot_API SHALL 返回 404 Not Found 错误
4. WHEN 请求体包含重复的 URL THEN Snapshot_API SHALL 基于唯一约束（scan_id + url）去重，忽略冲突
5. 资产表同步 SHALL 使用 upsert 策略：新记录插入，已存在记录更新（保留 created_at）
6. 响应 SHALL 返回成功写入的记录数

### 需求 2: 网站快照列表查询

**用户故事:** 作为安全分析师，我希望查看特定扫描中发现的网站，以便分析 Web 服务及其技术栈。

#### 验收标准

1. WHEN 用户请求 GET /api/scans/{scan_id}/websites/ THEN Snapshot_API SHALL 返回该扫描的分页 WebsiteSnapshot 记录列表
2. WHEN scan_id 不存在 THEN Snapshot_API SHALL 返回 404 Not Found 错误及相应错误信息
3. WHEN 提供 page 和 pageSize 查询参数 THEN Snapshot_API SHALL 返回指定页的结果
4. WHEN 未提供 page 参数 THEN Snapshot_API SHALL 默认返回第 1 页
5. WHEN 未提供 pageSize 参数 THEN Snapshot_API SHALL 默认每页返回 20 条记录
6. 响应 SHALL 包含分页元数据：总数、当前页、每页大小
7. 响应 SHALL 包含所有网站字段：id, scanId, url, host, title, statusCode, contentLength, location, webserver, contentType, tech, responseBody, vhost, responseHeaders, createdAt

### 需求 3: 网站快照过滤查询

**用户故事:** 作为安全分析师，我希望按多个字段过滤网站快照，以便快速找到相关记录。

#### 验收标准

1. WHEN 提供 filter 查询参数 THEN Snapshot_API SHALL 使用模糊匹配（LIKE）过滤 url、host、title、webserver 字段
2. WHEN filter 包含数字值 THEN Snapshot_API SHALL 同时使用精确匹配过滤 statusCode 字段
3. WHEN filter 参数为空或未提供 THEN Snapshot_API SHALL 返回该扫描的所有记录
4. 文本字段的过滤匹配 SHALL 不区分大小写

### 需求 4: 网站快照排序

**用户故事:** 作为安全分析师，我希望按不同字段对网站快照排序，以便根据分析需求组织数据。

#### 验收标准

1. WHEN 提供 ordering 查询参数 THEN Snapshot_API SHALL 按指定字段排序结果
2. Snapshot_API SHALL 支持按以下字段排序：url, host, title, statusCode, createdAt
3. WHEN ordering 以 "-" 前缀开头 THEN Snapshot_API SHALL 按降序排序
4. WHEN 未提供 ordering 参数 THEN Snapshot_API SHALL 默认按 createdAt 降序排序

### 需求 5: 网站快照导出

**用户故事:** 作为安全分析师，我希望将网站快照导出为 CSV，以便进行离线分析或与团队分享。

#### 验收标准

1. WHEN 用户请求 GET /api/scans/{scan_id}/websites/export/ THEN Snapshot_API SHALL 返回包含该扫描所有 WebsiteSnapshot 记录的 CSV 文件
2. WHEN scan_id 不存在 THEN Snapshot_API SHALL 返回 404 Not Found 错误
3. CSV 文件 SHALL 包含以下列：url, host, location, title, status_code, content_length, content_type, webserver, tech, response_body, response_headers, vhost, created_at
4. tech 数组字段 SHALL 在 CSV 中格式化为逗号分隔的值
5. Snapshot_API SHALL 使用流式响应以高效处理大数据集
6. 响应 Content-Type 头 SHALL 设置为 text/csv
7. 响应 Content-Disposition 头 SHALL 包含文件名，格式为：scan-{scan_id}-websites.csv

### 需求 6: API 响应格式

**用户故事:** 作为前端开发者，我希望 API 响应格式一致，以便轻松与后端集成。

#### 验收标准

1. Snapshot_API SHALL 为列表端点返回结构一致的 JSON 响应
2. 列表响应 SHALL 遵循格式：{ "results": [...], "total": number, "page": number, "pageSize": number, "totalPages": number }
3. 错误响应 SHALL 遵循格式：{ "error": { "code": string, "message": string } }
4. Snapshot_API SHALL 返回适当的 HTTP 状态码：200 表示成功，400 表示请求错误，404 表示未找到，500 表示服务器错误

### 需求 7: 扫描存在性验证

**用户故事:** 作为系统，我希望在查询或写入快照前验证扫描是否存在，以便提供有意义的错误信息并防止孤立数据。

#### 验收标准

1. WHEN 处理任何快照请求 THEN Snapshot_API SHALL 首先验证扫描是否存在
2. WHEN 扫描不存在 THEN Snapshot_API SHALL 返回 404 并提示 "Scan not found"
3. WHEN 扫描已被软删除（deleted_at 不为空）THEN Snapshot_API SHALL 将其视为不存在
4. WHEN 扫描已删除但 Worker 仍在提交结果 THEN Snapshot_API SHALL 拒绝写入并返回 404
