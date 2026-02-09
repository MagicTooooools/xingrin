# scan/repository

scan 模块 repository 规范：

- `scan.go` + `scan_query.go` + `scan_command.go`：Scan 仓储三分结构。
- `scan_adapter.go`：domain port 适配器（`domain <-> model` 映射）。
- `scan_log.go` + `scan_log_query.go` + `scan_log_command.go`：日志仓储三分结构。
- `scan_task.go` + `scan_task_query.go` + `scan_task_command.go` + `scan_task_sql.go`：任务仓储与 SQL 常量。

约束：

- 禁止使用 `*_mutation.go` 命名。
- 禁止使用泛名 `types.go`。
- `*_query.go` 不得出现写操作方法；`*_command.go` 不得出现查询方法。
