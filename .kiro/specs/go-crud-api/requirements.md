# 需求文档: Go CRUD API

## 简介

实现基础的 CRUD API，包括用户管理、组织、目标、扫描引擎。

## 需求

### 需求 1: 用户管理

**用户故事:** 作为管理员，我希望能创建用户和修改密码。

#### 验收标准

1. THE API SHALL 提供 `POST /api/users` 创建用户
2. THE API SHALL 提供 `PUT /api/users/:id/password` 修改密码
3. THE API SHALL 提供 `GET /api/users` 获取用户列表
4. THE API SHALL 使用 bcrypt 加密密码

### 需求 2: 组织管理

**用户故事:** 作为用户，我希望能管理组织。

#### 验收标准

1. THE API SHALL 提供 `GET /api/organizations` 获取组织列表（支持分页）
2. THE API SHALL 提供 `POST /api/organizations` 创建组织
3. THE API SHALL 提供 `GET /api/organizations/:id` 获取单个组织
4. THE API SHALL 提供 `PUT /api/organizations/:id` 更新组织
5. THE API SHALL 提供 `DELETE /api/organizations/:id` 软删除组织

### 需求 3: 目标管理

**用户故事:** 作为用户，我希望能管理扫描目标。

#### 验收标准

1. THE API SHALL 提供 `GET /api/targets` 获取目标列表（支持分页、筛选）
2. THE API SHALL 提供 `POST /api/targets` 创建目标
3. THE API SHALL 提供 `GET /api/targets/:id` 获取单个目标
4. THE API SHALL 提供 `PUT /api/targets/:id` 更新目标
5. THE API SHALL 提供 `DELETE /api/targets/:id` 软删除目标
6. THE API SHALL 自动检测目标类型（domain/ip/cidr）

### 需求 4: 扫描引擎管理

**用户故事:** 作为用户，我希望能管理扫描引擎配置。

#### 验收标准

1. THE API SHALL 提供 `GET /api/engines` 获取引擎列表
2. THE API SHALL 提供 `POST /api/engines` 创建引擎
3. THE API SHALL 提供 `GET /api/engines/:id` 获取单个引擎
4. THE API SHALL 提供 `PUT /api/engines/:id` 更新引擎
5. THE API SHALL 提供 `DELETE /api/engines/:id` 删除引擎

### 需求 5: 通用功能

#### 验收标准

1. THE API SHALL 支持分页参数 `page` 和 `pageSize`
2. THE API SHALL 返回统一的响应格式
3. THE API SHALL 所有接口需要认证（除了登录）
4. THE API SHALL 返回 camelCase JSON 字段名
