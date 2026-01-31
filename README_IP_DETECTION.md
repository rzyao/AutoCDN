# IP类型自动检测功能

## 概述

程序现在支持根据配置文件中指定的IP文件自动检测IP类型（IPv4或IPv6），并执行相应的测速和DNS更新操作。

## 功能特点

1. **自动检测IP类型**：程序会读取`config.yaml`中的`ip_file`指定的文件，自动判断是IPv4还是IPv6
2. **智能路由**：根据检测到的IP类型，自动选择相应的测速和DNS更新逻辑
3. **全面支持**：IPv4和IPv6都支持测试所有IP（通过`test_all_ip`配置）
4. **错误检查**：如果检测到IPv6文件但未配置IPv6域名，或检测到IPv4文件但未配置IPv4域名，程序会给出明确的错误提示

## 使用方法

### 1. 配置IP文件

在`config.yaml`中设置要测速的IP文件：

```yaml
speed_test:
  ip_file: "ip.txt"      # 指定IP段数据文件
  output: "result.csv"    # 输出结果文件
```

### 2. 配置域名

根据你要测速的IP类型，配置相应的域名：

**IPv4测速**：
```yaml
cloudflare:
  domains:                # IPv4域名列表
    - "aaa.ayaoblog.space"
    - "bbb.ayaoblog.space"
    - "ccc.ayaoblog.space"
```

**IPv6测速**：
```yaml
cloudflare:
  domainipv6s:           # IPv6域名列表
    - "666.ayaoblog.space"
    - "777.ayaoblog.space"
    - "888.ayaoblog.space"
```

### 3. 运行程序

程序会自动：
1. 读取指定的IP文件
2. 检测IP类型（IPv4或IPv6）
3. 执行相应的测速逻辑
4. 更新相应的DNS记录

## IP类型检测规则

- **IPv6检测**：文件中包含冒号（:）的IP地址
- **IPv4检测**：文件中包含点（.）的IP地址
- **默认**：如果无法确定，默认为IPv4

## 示例

### IPv4测速示例

1. 创建`ip.txt`文件，包含IPv4地址：
```
103.21.244.0/22
103.22.200.0/22
104.16.0.0/13
```

2. 配置`config.yaml`：
```yaml
speed_test:
  ip_file: "ip.txt"
  output: "result.csv"

cloudflare:
  domains:
    - "aaa.ayaoblog.space"
    - "bbb.ayaoblog.space"
```

3. 运行程序，会自动检测为IPv4并执行IPv4测速

### IPv6测速示例

1. 创建`ipv6.txt`文件，包含IPv6地址：
```
2606:4700::/32
2606:4700:10::/48
2606:4700:20::/48
```

2. 配置`config.yaml`：
```yaml
speed_test:
  ip_file: "ipv6.txt"
  output: "result.csv"

cloudflare:
  domainipv6s:
    - "666.ayaoblog.space"
    - "777.ayaoblog.space"
```

3. 运行程序，会自动检测为IPv6并执行IPv6测速

## 错误处理

程序会在以下情况给出错误提示：

1. **未指定IP文件**：`未指定IP文件，请在配置文件中设置 ip_file`
2. **IPv6文件但未配置IPv6域名**：`检测到IPv6文件，但未配置IPv6域名，请在配置文件中设置 domainipv6s`
3. **IPv4文件但未配置IPv4域名**：`检测到IPv4文件，但未配置IPv4域名，请在配置文件中设置 domains`
4. **文件读取失败**：`检测IP类型失败: [错误信息]`

## 注意事项

1. 确保IP文件格式正确（CIDR格式或单个IP）
2. 确保配置了相应类型的域名
3. 程序只会处理一种类型的IP，不会同时处理IPv4和IPv6
4. 输出文件名建议使用通用名称（如`result.csv`），因为程序会根据IP类型自动处理 