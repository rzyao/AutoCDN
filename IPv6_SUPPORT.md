# IPv6 支持说明

## 修改概述

本次修改为AutoCDN项目添加了IPv6域名DNS记录管理功能，使其能够同时处理IPv4和IPv6域名的DNS更新。

## 主要修改内容

### 1. 配置文件修改 (`config/config.go`)

- 在 `CloudflareConfig` 结构体中添加了 `DomainIPv6s []string` 字段
- 支持从 `config.yaml` 中读取 `domainipv6s` 配置项

### 2. CDN处理逻辑修改 (`cdn/cdn.go`)

#### 新增函数：
- `UpdateDNSRecordsIPv6()` - 更新IPv6 DNS记录（AAAA记录）
- `CreateDNSRecordIPv6()` - 创建新的IPv6 DNS记录（AAAA记录）
- `GetRecordListWithType()` - 获取包含记录类型的DNS记录列表
- `HandleDNSRecordsIPv6()` - 处理IPv6 DNS记录的更新或创建
- `HandleAllDNSRecords()` - 同时处理IPv4和IPv6 DNS记录
- `GetIPListForIPv6Domains()` - 为IPv6域名获取IP列表

#### 修改函数：
- `GetRecordList()` - 添加了对记录类型的过滤，只处理A记录和AAAA记录
- `UpdateDNSRecords()` - 保持原有IPv4功能不变

### 3. 主程序修改 (`main.go`)

- 修改了 `main()` 函数，使其能够分别处理IPv4和IPv6域名
- 为IPv6域名处理添加了独立的测速流程
- 支持使用不同的IP文件和输出文件（`ipv6.txt` 和 `result6.csv`）

## 配置文件示例

```yaml
cloudflare:
  api_key: "your_api_key"
  zone_id: "your_zone_id"
  zone_name: "your_domain.com"
  email: "your_email@example.com"
  domains:                    # IPv4域名列表
    - "aaa.your_domain.com"
    - "bbb.your_domain.com"
    - "ccc.your_domain.com"
  domainipv6s:               # IPv6域名列表
    - "666.your_domain.com"
    - "777.your_domain.com"
    - "888.your_domain.com"
```

## 使用方法

1. 在 `config.yaml` 中配置IPv6域名列表
2. 确保 `ipv6.txt` 文件包含IPv6地址段
3. 运行程序，它会自动：
   - 先处理IPv4域名（使用 `ip.txt` 和 `result.csv`）
   - 再处理IPv6域名（使用 `ipv6.txt` 和 `result6.csv`）

## 工作流程

1. **IPv4处理流程**：
   - 使用 `ip.txt` 中的IPv4地址段进行测速
   - 将测速结果应用到 `domains` 列表中的域名
   - 创建或更新A记录

2. **IPv6处理流程**：
   - 使用 `ipv6.txt` 中的IPv6地址段进行测速
   - 将测速结果应用到 `domainipv6s` 列表中的域名
   - 创建或更新AAAA记录

## 注意事项

- IPv4和IPv6域名可以同时配置，程序会分别处理
- 如果只配置了IPv4域名，程序只会处理IPv4
- 如果只配置了IPv6域名，程序只会处理IPv6
- 确保Cloudflare API有足够的权限来创建和更新DNS记录
- IPv6地址段文件 `ipv6.txt` 需要包含有效的IPv6 CIDR格式

## 兼容性

- 完全向后兼容，现有的IPv4配置无需修改
- 新增的IPv6功能不会影响现有的IPv4功能
- 可以逐步迁移到IPv6，也可以同时使用两种协议 