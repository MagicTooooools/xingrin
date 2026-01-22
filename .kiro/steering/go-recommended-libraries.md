# Go 推荐库和最佳实践

本文档列出了项目中推荐使用的 Go 第三方库，以及使用这些库的最佳实践。

## 核心原则

1. **优先使用成熟的开源库**：避免重复造轮子
2. **选择活跃维护的项目**：确保长期支持和安全更新
3. **考虑性能和安全性**：特别是在安全扫描场景下
4. **保持依赖最小化**：只引入真正需要的库

## 推荐库列表

### 网络和验证

#### 1. ProjectDiscovery Utils
**包名**：`github.com/projectdiscovery/utils`

**用途**：网络相关的工具函数，专为安全工具优化

**推荐使用的子包**：
- `github.com/projectdiscovery/utils/ip` - IP 地址验证和处理
- `github.com/projectdiscovery/utils/strings` - 字符串处理

**使用示例**：
```go
import (
    iputil "github.com/projectdiscovery/utils/ip"
)

// IP 地址验证
if iputil.IsIP("192.168.1.1") {
    // 处理 IP 地址
}

// IPv4 验证
if iputil.IsIPv4("192.168.1.1") {
    // 处理 IPv4
}

// IPv6 验证
if iputil.IsIPv6("2001:db8::1") {
    // 处理 IPv6
}

// CIDR 验证
if iputil.IsCIDR("192.168.1.0/24") {
    // 处理 CIDR
}

// 内网 IP 检测
if iputil.IsInternal("192.168.1.1") {
    // 处理内网 IP
}
```

**优势**：
- 与 nuclei、subfinder 等工具使用相同的底层库
- 针对大规模扫描场景优化
- 活跃的安全社区维护

**项目中的使用**：
- `worker/internal/pkg/validator/domain.go` - 子域名验证中的 IP 检测

#### 2. govalidator
**包名**：`github.com/asaskevich/govalidator`

**用途**：通用的字符串验证库

**常用函数**：
```go
import "github.com/asaskevich/govalidator"

// DNS 名称验证
if govalidator.IsDNSName("example.com") {
    // 有效的域名
}

// Email 验证
if govalidator.IsEmail("user@example.com") {
    // 有效的邮箱
}

// URL 验证
if govalidator.IsURL("https://example.com") {
    // 有效的 URL
}

// IP 验证（通用）
if govalidator.IsIP("192.168.1.1") {
    // 有效的 IP
}
```

**优势**：
- 成熟稳定，广泛使用
- 支持多种验证类型
- 依赖少，性能好

**项目中的使用**：
- `worker/internal/pkg/validator/domain.go` - 域名格式验证

### 其他 ProjectDiscovery 生态系统库

#### 3. mapcidr
**包名**：`github.com/projectdiscovery/mapcidr`

**用途**：CIDR 操作和 IP 范围处理

**使用场景**：
- 扩展 CIDR 为 IP 列表
- IP 范围计算
- 子网操作

**示例**：
```go
import "github.com/projectdiscovery/mapcidr"

// 扩展 CIDR
ips, err := mapcidr.IPAddresses("192.168.1.0/24")
// 返回 192.168.1.0 到 192.168.1.255 的所有 IP

// 计算 CIDR 中的 IP 数量
count := mapcidr.CountIPsInCIDR("192.168.1.0/24")
```

#### 4. cdncheck
**包名**：`github.com/projectdiscovery/cdncheck`

**用途**：检测 IP 是否属于 CDN 或云服务提供商

**使用场景**：
- 识别 CDN IP
- 过滤云服务 IP
- 优化扫描目标

**示例**：
```go
import (
    "net"
    "github.com/projectdiscovery/cdncheck"
)

client := cdncheck.New()

// 检查是否为 CDN
if client.Check(net.ParseIP("1.1.1.1")) {
    // 这是 CDN IP
}

// 检查域名
if matched, _, err := client.CheckDomainWithFallback("example.com"); matched {
    // 域名使用了 CDN
}
```

#### 5. retryabledns
**包名**：`github.com/projectdiscovery/retryabledns`

**用途**：可重试的 DNS 客户端

**使用场景**：
- 高可靠性 DNS 查询
- 自定义 DNS 解析器
- 批量 DNS 查询

**示例**：
```go
import "github.com/projectdiscovery/retryabledns"

// 创建客户端
client := retryabledns.New([]string{"8.8.8.8:53", "1.1.1.1:53"}, 3)

// DNS 查询
ips, err := client.Resolve("example.com")
```

### 日志和配置

#### 6. zap
**包名**：`go.uber.org/zap`

**用途**：高性能结构化日志

**项目中的使用**：
```go
import "go.uber.org/zap"

// 记录日志
pkg.Logger.Info("Operation completed",
    zap.String("operation", "scan"),
    zap.Int("count", 100),
    zap.Duration("duration", time.Second))

// 错误日志
pkg.Logger.Error("Operation failed",
    zap.String("operation", "scan"),
    zap.Error(err))
```

**优势**：
- 高性能，零内存分配
- 结构化日志，易于解析
- 类型安全

#### 7. viper
**包名**：`github.com/spf13/viper`

**用途**：配置管理

**使用场景**：
- 读取配置文件（YAML、JSON、TOML）
- 环境变量管理
- 配置热更新

## 使用规范

### 1. 导入别名

当包名可能冲突时，使用别名：

```go
import (
    iputil "github.com/projectdiscovery/utils/ip"
    strutil "github.com/projectdiscovery/utils/strings"
)
```

### 2. 错误处理

始终检查错误，不要忽略：

```go
// ✅ 正确
if iputil.IsIP(s) {
    // 处理 IP
}

// ❌ 错误 - 不要自己实现已有的功能
func isIPLike(s string) bool {
    // 自定义实现...
}
```

### 3. 性能考虑

在高频调用的场景下，选择性能优化的库：

```go
// ProjectDiscovery 的库针对大规模扫描优化
for _, subdomain := range millions {
    if iputil.IsIP(subdomain) {
        continue
    }
    // 处理子域名
}
```

### 4. 依赖管理

添加新依赖后，运行：

```bash
go mod tidy
```

确保 `go.mod` 和 `go.sum` 正确更新。

## 避免使用的模式

### ❌ 不要重复造轮子

```go
// ❌ 错误 - 自己实现 IP 验证
func isIPLike(s string) bool {
    parts := strings.Split(s, ".")
    if len(parts) != 4 {
        return false
    }
    // ... 更多代码
}

// ✅ 正确 - 使用现有库
if iputil.IsIP(s) {
    // 处理 IP
}
```

### ❌ 不要忽略错误

```go
// ❌ 错误
_ = file.Close()

// ✅ 正确
if err := file.Close(); err != nil {
    log.Warn("Failed to close file", zap.Error(err))
}
```

### ❌ 不要使用过时的库

在选择库时，检查：
- 最后更新时间
- GitHub stars 和 forks
- Issue 响应速度
- 社区活跃度

## 库版本管理

### 更新依赖

定期检查和更新依赖：

```bash
# 查看可更新的依赖
go list -u -m all

# 更新特定依赖
go get -u github.com/projectdiscovery/utils

# 更新所有依赖
go get -u ./...
```

### 安全更新

关注安全公告，及时更新有漏洞的依赖：

```bash
# 使用 govulncheck 检查漏洞
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## 参考资源

### ProjectDiscovery 生态系统
- 官方文档：https://docs.projectdiscovery.io/
- GitHub 组织：https://github.com/projectdiscovery
- 工具列表：nuclei, subfinder, httpx, dnsx, naabu, katana

### Go 标准库
- 官方文档：https://pkg.go.dev/std
- 优先使用标准库，除非有特殊需求

### 库选择指南
1. 检查 GitHub stars（> 1000 为佳）
2. 查看最近更新时间（< 6 个月为佳）
3. 阅读文档和示例
4. 检查 Issue 和 PR 活跃度
5. 查看依赖数量（越少越好）

## 总结

- **优先使用 ProjectDiscovery 生态系统的库**：专为安全工具优化
- **使用成熟的通用库**：如 govalidator、zap
- **避免重复造轮子**：充分利用开源社区的成果
- **保持依赖更新**：定期检查安全更新
- **遵循最佳实践**：错误处理、性能优化、代码可读性
