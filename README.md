# UOS Squid Exporter for Prometheus
## 项目简介

Squid 代理服务器监控工具。它连接到 Squid 代理服务器，收集性能指标并以 Prometheus 格式导出，包括客户端和服务器 HTTP 统计信息、缓存性能、服务时间和各种系统资源使用指标。

## 安装

### 从源码编译

```bash
git clone https://gitee.com/deepin-community/uos-squid-exporter.git
cd uos-squid-exporter
go build
```

### 二进制安装

从 [发布页面](https://gitee.com/deepin-community/uos-squid-exporter/releases) 下载适用于您系统的最新二进制文件。

## 使用方法

### 基本使用

```bash
./uos-squid-exporter
```

### 指定 Squid 服务器参数

```bash
./uos-squid-exporter --squid.hostname localhost --squid.port 3128
```

### 配置参数

导出器支持以下命令行参数：

```bash
--squid.hostname       Squid 服务器主机名 (默认: "localhost")
--squid.port           Squid 服务器端口 (默认: 3128)
--squid.login          Squid 服务器登录用户名 (如需认证)
--squid.password       Squid 服务器登录密码 (如需认证)
--squid.extractTimes   是否提取服务时间指标 (默认: true)
```

### YAML 配置文件

也可以使用 YAML 配置文件：

```yaml
address: "127.0.0.1"  # 导出器监听地址
port: 8090            # 导出器监听端口
metricsPath: "/metrics"

# Squid 配置
squid:
  hostname: "localhost"
  port: 3128
  login: ""
  password: ""
  extractTimes: true
```

## 监控指标

### 客户端/服务器 HTTP 指标

- 客户端 HTTP 请求总数
- 客户端 HTTP 命中总数
- 客户端 HTTP 错误总数
- 服务器 HTTP 请求总数
- 服务器 HTTP 错误总数
- 更多...

### 服务时间指标

- HTTP 请求服务时间
- 缓存命中服务时间
- 缓存未命中服务时间
- 近似命中服务时间
- DNS 查找服务时间

### 系统信息

- 访问缓存的客户端数量
- CPU 使用率
- 内存使用情况
- 文件描述符使用情况
- 存储指标
- 更多...

## Prometheus 配置

在您的 `prometheus.yaml` 中添加以下配置：

```yaml
scrape_configs:
  - job_name: "uos-squid"
    static_configs:
      - targets: ["localhost:8090"]
```

## Squid 配置

为了允许导出器查询 Squid 指标，请在您的 squid.conf 中添加：

```
# 允许来自本地主机的缓存管理器访问
acl prometheus src 127.0.0.1
http_access allow manager prometheus
```
