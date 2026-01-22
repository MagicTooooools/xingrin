# 设计文档

## 概述

补全 Go 后端的所有数据模型，确保与 Django 模型完全兼容。

## 模型设计原则

1. **表名一致** - TableName() 返回与 Django db_table 相同的值
2. **字段名一致** - gorm column tag 使用 snake_case
3. **JSON 输出** - json tag 使用 camelCase
4. **类型映射** - PostgreSQL 数组用 pq.StringArray/Int64Array，JSONB 用 datatypes.JSON

## 模型文件组织

```
go-backend/internal/model/
├── organization.go      # ✅ 已有
├── target.go            # ✅ 已有
├── scan.go              # ⚠️ 需补充字段
├── subdomain.go         # ✅ 已有
├── website.go           # ⚠️ 需补充字段
├── worker_node.go       # ✅ 已有
├── scan_engine.go       # ✅ 已有
├── endpoint.go          # 新增
├── directory.go         # 新增
├── host_port_mapping.go # 新增
├── vulnerability.go     # 新增
├── screenshot.go        # 新增
├── subdomain_snapshot.go    # 新增
├── website_snapshot.go      # 新增
├── endpoint_snapshot.go     # 新增
├── directory_snapshot.go    # 新增
├── host_port_mapping_snapshot.go # 新增
├── vulnerability_snapshot.go    # 新增
├── screenshot_snapshot.go       # 新增
├── scan_log.go          # 新增
├── scan_input_target.go # 新增
├── scheduled_scan.go    # 新增
├── subfinder_provider_settings.go # 新增
├── wordlist.go          # 新增
├── nuclei_template_repo.go # 新增
├── notification.go      # 新增
├── notification_settings.go # 新增
├── blacklist_rule.go    # 新增
├── asset_statistics.go  # 新增
├── statistics_history.go # 新增
├── user.go              # 新增
├── session.go           # 新增
└── model_test.go        # 更新测试
```

## 正确性属性

### Property 1: 表名映射正确性
*对于任意* Go 模型，其 TableName() 方法返回的表名应与 Django 模型的 db_table 一致。
**验证: 需求 9.1**

### Property 2: JSON 字段名转换正确性
*对于任意* Go 模型序列化为 JSON，所有字段名应为 camelCase 格式。
**验证: 需求 9.5**

## 测试策略

- 更新 model_test.go，添加所有新模型的表名测试
- 验证 JSON 序列化输出格式
