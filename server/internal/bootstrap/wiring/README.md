# wiring 规范（模块适配层）

本文用于约束 `server/internal/bootstrap/wiring/*` 下各模块 wiring 的命名与组织，避免后续再次出现风格漂移。

## 目标
- 统一导出函数命名（便于全局搜索与批量替换）
- 统一 adapter 类型命名（降低阅读成本）
- 用编译期断言固定接口边界（重构时第一时间暴露断裂）

## 目录模板
每个子模块（如 `scan/`、`catalog/`）建议包含：

- `exports.go`：模块对外公开的 wiring 入口
- `wiring_<module>_*_adapter.go`：具体适配器实现
- `wiring_<module>_adapter_assertions.go`：接口断言（推荐必备）
- `wiring_<module>_*_service.go`：仅在需要组装应用服务时使用

## 命名规则

### 1) 导出函数（`exports.go`）
- 适配器导出统一：`New<Module><Role>Adapter`
  - 例：`NewScanTaskStoreAdapter`
  - 例：`NewScanTargetLookupAdapter`
- 应用服务导出统一：`New<Module>ApplicationService`
  - 例：`NewScanLogApplicationService`
  - 例：`NewWorkerApplicationService`

> 约束：导出函数优先返回 application/domain 层接口，而不是具体 struct。

### 2) 私有构造函数
- 统一：`new<Module><Role>Adapter`
  - 例：`newScanTaskRuntimeStoreAdapter`
  - 例：`newScanLogLookupAdapter`

### 3) 适配器类型
- 统一：`<module><role>Adapter`
  - 例：`scanTaskStoreAdapter`
  - 例：`catalogEngineStoreAdapter`

## 接口断言规范
每个子模块应有 `wiring_<module>_adapter_assertions.go`，写法固定：

```go
var _ someapp.SomePort = (*someAdapter)(nil)
```

建议至少覆盖：
- `exports.go` 对外返回的全部接口
- 同一 adapter 实现的 Query/Command 复合接口

## 新增 wiring 子模块时的最小清单
- [ ] 有 `exports.go`
- [ ] 导出函数命名符合 `New<Module><Role>Adapter` / `New<Module>ApplicationService`
- [ ] 适配器命名符合 `<module><role>Adapter`
- [ ] 有 `wiring_<module>_adapter_assertions.go`
- [ ] 在 `internal/bootstrap/wiring.go` 中仅通过导出函数注入

## 反例（避免）
- `NewApplicationService`（缺少模块前缀）
- `newCommandStore`（缺少模块前缀）
- `taskRuntimeScanStoreAdapter`（词序和模块命名不一致）

## 说明
该规范优先级低于系统/仓库级指令；若未来确需偏离，请在对应模块目录补充 README 说明原因。
