# Go 后端迁移策略

## 迁移目标

本项目正在从 Python (Django) 完全迁移到 Go。`server/` 是新的主后端，`backend/` (Python) 将逐步废弃。

## 迁移原则

1. **全新数据库 schema** - Go 后端拥有独立的数据库 schema 管理权，不依赖 Django migrations
2. **SQL 迁移文件** - 使用 golang-migrate 管理数据库迁移，手写 SQL 文件
3. **版本控制** - 每次 schema 变更创建新的迁移文件，支持回滚

## 数据库迁移

### 当前策略

- 使用 `golang-migrate` 执行 SQL 迁移文件
- 迁移文件位于 `server/cmd/server/migrations/`
- 启动时自动执行所有待执行的迁移
- 外键约束已启用，使用 `ON DELETE CASCADE`

### 迁移文件命名规范

```
000001_init_schema.up.sql      # 创建初始表结构
000001_init_schema.down.sql    # 回滚初始表结构
000002_add_gin_indexes.up.sql  # 添加 GIN 索引
000002_add_gin_indexes.down.sql
```

### 创建新迁移

```bash
# 手动创建迁移文件
touch server/cmd/server/migrations/000003_add_new_table.up.sql
touch server/cmd/server/migrations/000003_add_new_table.down.sql
```

### 迁移 API

```go
// 执行所有待执行的迁移
database.RunMigrations(sqlDB)

// 回滚最后一个迁移
database.MigrateDown(sqlDB)

// 迁移到指定版本
database.MigrateToVersion(sqlDB, 1)

// 获取当前版本
version, dirty, err := database.GetMigrationVersion(sqlDB)
```

### GIN 索引

PostgreSQL 数组字段需要 GIN 索引以支持高效查询：

```sql
-- 创建 GIN 索引
CREATE INDEX idx_website_tech_gin ON website USING GIN (tech);

-- 查询示例（使用 @> 包含操作符）
SELECT * FROM website WHERE tech @> ARRAY['nginx'];
```

已添加 GIN 索引的字段：
- `website.tech`
- `endpoint.tech`
- `endpoint.matched_gf_patterns`
- `scan.engine_ids`
- `scan.container_ids`

## 模型定义规范

参考 `项目信息以及如何快速自动化操作.md` 中的 Go 后端部分。

**注意**: 模型定义仅用于 GORM 查询，不再用于自动迁移。Schema 变更必须通过 SQL 迁移文件。

## 迁移进度

- [x] 基础模型（Target, Scan）
- [x] 资产模型（Subdomain, Host, Website）
- [x] 快照模型（各类 Snapshot）
- [x] GIN 索引（数组字段）
- [ ] API 接口迁移
- [ ] 扫描引擎迁移
- [ ] 任务队列（硬删除 + CASCADE 清理）

## 待完善功能

### 删除策略

当前状态：
- ✅ 软删除已实现（设置 `deleted_at`）
- ❌ 硬删除未实现

计划：
- 等任务队列完善后，实现后台硬删除任务
- 硬删除时使用数据库 CASCADE 自动清理关联数据（Website、Subdomain 等）
- 流程：用户删除 → 软删除 → 后台任务 → 硬删除 + CASCADE

### 软删除注意事项

统计关联数据时必须排除已软删除的记录：

```go
// ❌ 错误 - 会统计已删除的 target
SELECT COUNT(*) FROM organization_target 
WHERE organization_id = ?

// ✅ 正确 - 排除已软删除的 target
SELECT COUNT(*) FROM organization_target 
INNER JOIN target ON target.id = organization_target.target_id 
WHERE organization_id = ? AND target.deleted_at IS NULL
```

查询关联数据时同理，始终添加 `deleted_at IS NULL` 条件。
