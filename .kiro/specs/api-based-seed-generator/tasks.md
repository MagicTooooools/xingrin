# Implementation Plan: API-Based Seed Generator

## Overview

实现一个基于 API 的种子数据生成器，使用 Python 通过 HTTP 请求调用 Go 后端 API 来创建测试数据。项目采用模块化设计，分为 5 个独立文件，总计约 1000 行代码。

## Tasks

- [x] 1. 创建项目结构和基础文件
  - 创建 `tools/seed-api/` 目录
  - 创建 `requirements.txt` 文件（依赖 requests>=2.31.0）
  - 创建 `README.md` 文件（使用说明）
  - 创建空的 Python 模块文件（5 个 .py 文件）
  - _Requirements: 16.3_

- [x] 2. 实现 API Client 模块
  - [x] 2.1 实现 APIClient 类基础结构
    - 实现 `__init__` 方法（初始化 base_url、username、password、Session）
    - 实现 `login` 方法（POST /api/auth/login，获取 JWT token）
    - 实现 `_get_headers` 方法（返回带 Authorization 的请求头）
    - _Requirements: 1.1, 1.2_

  - [x] 2.2 实现 HTTP 请求方法
    - 实现 `post` 方法（发送 POST 请求，自动添加认证头）
    - 实现 `get` 方法（发送 GET 请求，自动添加认证头）
    - 实现 `delete` 方法（发送 DELETE 请求，自动添加认证头）
    - 所有方法使用 30 秒超时
    - _Requirements: 1.2_

  - [x] 2.3 实现 Token 自动刷新
    - 实现 `refresh_token` 方法（POST /api/auth/refresh）
    - 在请求方法中捕获 401 错误，自动调用 refresh_token 后重试
    - _Requirements: 1.3_

  - [x] 2.4 实现错误处理
    - 捕获 requests 异常（Timeout、ConnectionError）
    - 解析 API 错误响应（JSON 格式）
    - 返回统一的错误信息
    - _Requirements: 1.4_

- [x] 3. 实现 Error Handler 模块
  - [x] 3.1 实现 ErrorHandler 类基础结构
    - 实现 `__init__` 方法（初始化 max_retries、retry_delay）
    - 实现 `should_retry` 方法（判断状态码是否应该重试）
    - _Requirements: 12.1, 12.2_

  - [x] 3.2 实现重试逻辑
    - 实现 `handle_error` 方法（根据错误类型决定是否重试）
    - 5xx 错误重试 3 次，每次等待 1 秒
    - 429 错误重试 3 次，每次等待 5 秒
    - 网络超时重试 3 次，每次等待 1 秒
    - _Requirements: 12.1, 12.4_

  - [x] 3.3 实现错误日志
    - 实现 `log_error` 方法（记录错误详情到文件）
    - 日志包含时间戳、错误类型、请求数据、响应数据
    - 日志文件：`seed_errors.log`
    - _Requirements: 12.2, 15.3_

- [x] 4. 实现 Progress Tracker 模块
  - [x] 4.1 实现 ProgressTracker 类基础结构
    - 实现 `__init__` 方法（初始化统计变量）
    - 实现 `start_phase` 方法（开始新阶段，记录阶段名和总数）
    - 实现 `update` 方法（更新当前进度）
    - _Requirements: 13.1_

  - [x] 4.2 实现统计功能
    - 实现 `add_success` 方法（记录成功数量）
    - 实现 `add_error` 方法（记录错误信息）
    - 实现 `finish_phase` 方法（完成当前阶段，显示总结）
    - _Requirements: 13.2, 13.3_

  - [x] 4.3 实现进度显示
    - 使用 emoji 图标显示不同阶段（🏢 🎯 🔗 🌐 📝 等）
    - 显示进度条格式：`[当前/总数]`
    - 显示成功数量和错误数量
    - _Requirements: 13.1, 13.2_

  - [x] 4.4 实现总结报告
    - 实现 `print_summary` 方法（打印最终统计）
    - 显示总耗时、成功记录数、错误记录数
    - _Requirements: 13.4_

- [x] 5. 实现 Data Generator 模块
  - [x] 5.1 实现组织数据生成
    - 实现 `generate_organization` 方法
    - 使用预定义的组织名称列表 + 随机后缀
    - 生成合理的描述文本
    - 返回 Python 字典（camelCase 字段名）
    - _Requirements: 2.2, 15.1_

  - [x] 5.2 实现目标数据生成
    - 实现 `generate_targets` 方法
    - 按比例生成域名（70%）、IP（20%）、CIDR（10%）
    - 域名格式：`{env}.{company}-{suffix}.{tld}`
    - IP 格式：随机合法 IPv4
    - CIDR 格式：随机 /8、/16、/24 网段
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 3.2, 15.1_

  - [x] 5.3 实现 Website 数据生成
    - 实现 `generate_websites` 方法
    - 根据目标类型生成合理的 URL
    - 生成 title、statusCode、contentLength、tech 等字段
    - tech 字段为数组（如 ["nginx", "PHP", "MySQL"]）
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 5.2, 5.3, 15.1_

  - [x] 5.4 实现 Subdomain 数据生成
    - 实现 `generate_subdomains` 方法
    - 仅为域名类型目标生成
    - 使用常见子域名前缀（www、api、admin 等）
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 6.2, 6.3, 15.1_

  - [x] 5.5 实现 Endpoint 数据生成
    - 实现 `generate_endpoints` 方法
    - 生成 API 路径（/api/v1/users、/login 等）
    - 生成 tech、matchedGfPatterns 数组字段
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 7.2, 7.3, 15.1_

  - [x] 5.6 实现 Directory 数据生成
    - 实现 `generate_directories` 方法
    - 生成常见目录路径（/admin/、/backup/ 等）
    - 生成 status、contentLength、duration 字段
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 8.2, 8.3, 15.1_

  - [x] 5.7 实现 HostPort 数据生成
    - 实现 `generate_host_ports` 方法
    - 生成 host、ip、port 字段
    - 使用常见端口（80、443、22、3306 等）
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 9.2, 9.3, 15.1_

  - [x] 5.8 实现 Vulnerability 数据生成
    - 实现 `generate_vulnerabilities` 方法
    - 生成 vulnType、severity、cvssScore 等字段
    - 按比例生成不同严重级别
    - 返回 Python 字典列表（camelCase 字段名）
    - _Requirements: 10.2, 10.3, 15.1_

- [x] 6. 实现主程序（seed_generator.py）
  - [x] 6.1 实现命令行参数解析
    - 使用 argparse 解析参数
    - 支持 --api-url、--username、--password
    - 支持 --orgs、--targets-per-org、--assets-per-target
    - 支持 --clear、--batch-size、--verbose
    - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6_

  - [x] 6.2 实现数据清理功能
    - 调用批量删除 API 清空数据
    - 按正确顺序删除（先资产，再目标，最后组织）
    - 显示清理进度和结果
    - _Requirements: 14.1, 14.2, 14.4_

  - [x] 6.3 实现组织和目标创建流程
    - 调用 API Client 登录获取 token
    - 调用 POST /api/organizations 创建组织
    - 调用 POST /api/targets/batch_create 批量创建目标
    - 调用 POST /api/organizations/:id/link_targets 关联目标
    - 使用 Progress Tracker 显示进度
    - 使用 Error Handler 处理错误和重试
    - _Requirements: 2.1, 2.3, 3.1, 3.4, 4.1, 4.2, 4.4_

  - [x] 6.4 实现资产创建流程
    - 调用 POST /api/targets/:id/websites/bulk-upsert 创建 Website
    - 调用 POST /api/targets/:id/subdomains/bulk-create 创建 Subdomain
    - 调用 POST /api/targets/:id/endpoints/bulk-upsert 创建 Endpoint
    - 调用 POST /api/targets/:id/directories/bulk-upsert 创建 Directory
    - 调用 POST /api/targets/:id/host-ports/bulk-upsert 创建 HostPort
    - 调用 POST /api/targets/:id/vulnerabilities/bulk-create 创建 Vulnerability
    - 每个资产类型分批发送（batch_size=100）
    - 使用 Progress Tracker 显示进度
    - 使用 Error Handler 处理错误和重试
    - _Requirements: 5.1, 5.4, 6.1, 6.4, 7.1, 7.4, 8.1, 8.4, 9.1, 9.4, 10.1, 10.4_

  - [x] 6.5 实现主流程编排
    - 检查 Python 版本（需要 3.8+）
    - 检查 requests 库是否已安装
    - 按顺序执行：清理 → 组织 → 目标 → 关联 → 资产
    - 捕获所有异常，显示友好的错误信息
    - 最后显示总结报告
    - _Requirements: 16.1, 16.4_

- [x] 7. 编写文档和测试
  - [x] 7.1 编写 README.md
    - 说明项目用途和功能
    - 列出依赖和安装步骤
    - 提供使用示例和命令行参数说明
    - 说明常见问题和解决方法
    - _Requirements: 16.3_

  - [x] 7.2 编写单元测试
    - 测试 Data Generator 的数据生成逻辑
    - 测试 API Client 的请求构造
    - 测试 Error Handler 的重试逻辑
    - 使用 pytest 框架
    - _Requirements: 所有需求_

  - [x] 7.3 执行集成测试
    - 启动 Go 后端
    - 运行种子生成器（小规模：2 个组织，10 个目标）
    - 验证数据是否正确创建
    - 验证 JSON 字段命名（camelCase）
    - _Requirements: 15.4_

- [x] 8. Checkpoint - 确保所有功能正常
  - 运行完整的数据生成流程
  - 验证所有 API 调用成功
  - 验证错误处理和重试机制
  - 验证进度显示和统计
  - 如有问题，询问用户

## Notes

- 所有任务都是必做的，确保从一开始就保证质量
- 每个任务引用了具体的需求编号，便于追溯
- Checkpoint 任务确保增量验证
- Python 代码使用 camelCase 构造 JSON，但变量名使用 snake_case（Python 规范）
- 所有 API 请求使用独立构造的 Python 字典，不依赖 Go 后端代码
