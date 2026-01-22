# 设计文档: Go CRUD API

## 概述

实现 RESTful CRUD API，采用 Handler -> Service -> Repository 三层架构。

## 架构

```
HTTP Request → Handler → Service → Repository → Database
                 ↓
            Validation
```

## 目录结构

```
internal/
├── handler/
│   ├── user_handler.go
│   ├── organization_handler.go
│   ├── target_handler.go
│   └── engine_handler.go
├── service/
│   ├── user_service.go
│   ├── organization_service.go
│   ├── target_service.go
│   └── engine_service.go
├── repository/
│   ├── user_repository.go
│   ├── organization_repository.go
│   ├── target_repository.go
│   └── engine_repository.go
└── dto/
    ├── request.go      # 请求 DTO
    ├── response.go     # 响应 DTO
    └── pagination.go   # 分页
```

## API 设计

### 统一响应格式

成功:
```json
{
    "data": { ... },
    "message": "success"
}
```

列表:
```json
{
    "data": [ ... ],
    "total": 100,
    "page": 1,
    "pageSize": 20
}
```

错误:
```json
{
    "error": "error message"
}
```

### 分页参数

- `page`: 页码，默认 1
- `pageSize`: 每页数量，默认 20，最大 100

### API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 创建用户 |
| GET | /api/users | 用户列表 |
| PUT | /api/users/:id/password | 修改密码 |
| GET | /api/organizations | 组织列表 |
| POST | /api/organizations | 创建组织 |
| GET | /api/organizations/:id | 获取组织 |
| PUT | /api/organizations/:id | 更新组织 |
| DELETE | /api/organizations/:id | 删除组织 |
| GET | /api/targets | 目标列表 |
| POST | /api/targets | 创建目标 |
| GET | /api/targets/:id | 获取目标 |
| PUT | /api/targets/:id | 更新目标 |
| DELETE | /api/targets/:id | 删除目标 |
| GET | /api/engines | 引擎列表 |
| POST | /api/engines | 创建引擎 |
| GET | /api/engines/:id | 获取引擎 |
| PUT | /api/engines/:id | 更新引擎 |
| DELETE | /api/engines/:id | 删除引擎 |

## 软删除

Organization 和 Target 使用软删除：
- 删除时设置 `deleted_at` 字段
- 查询时默认过滤已删除记录
