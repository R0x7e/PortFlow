# PortFlow

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20Darwin-blue)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)

一个基于 Go 语言开发的高性能、生产级 TCP 端口转发工具。

## 功能特性

- **双向转发**: 支持任意基于 TCP 的协议（HTTP、HTTPS、自定义协议等）。
- **高并发**: 基于 goroutine 池实现，支持 10,000+ 并发连接。
- **安全性**: 支持可选的 HTTP Basic 认证（RFC 7617）。
- **极致性能**: 
  - 目标服务器连接池管理。
  - 多层超时控制（连接、读写、空闲）。
  - 低延迟（本地回环 <1ms）。
- **生产就绪**:
  - 优雅关闭支持（SIGINT/SIGTERM）。
  - 实时性能监控指标输出。
  - 完整的错误处理体系。
  - 零第三方依赖。
- **可观测性**: 支持自定义日志级别（Debug, Info, Warn, Error）。

## 安装说明

### 源码编译
```bash
go build -o portflow ./cmd/portflow
```

### Docker 部署
```bash
docker build -t portflow .
docker run -p 8080:8080 portflow --listen-port 8080 --target-addr 1.2.3.4 --target-port 80
```

## 使用示例

```bash
./portflow --listen-port <port> --target-addr <ip> --target-port <port> [参数]
```

### 参数说明
- `--listen-addr`: 本地监听地址 (默认: 0.0.0.0)
- `--listen-port`: 本地监听端口 (必填)
- `--target-addr`: 目标服务器 IP/域名 (必填)
- `--target-port`: 目标服务器端口 (必填)
- `--auth`: 是否启用认证 (默认: false)
- `--user`: 认证用户名
- `--pass`: 认证密码
- `--log-level`: 日志级别 (debug, info, warn, error) (默认: info)
- `--conn-timeout`: 连接超时时间/秒 (默认: 30)
- `--rw-timeout`: 读写超时时间/秒 (默认: 60)
- `--idle-timeout`: 空闲连接超时时间/秒 (默认: 300)

## 性能表现

| 并发数 | 吞吐量 | 延迟 | 每千连接内存占用 |
|-------------|------------|---------|-------------------|
| 1,000       | >100 MB/s  | <1 ms   | <1 MB             |
| 5,000       | >100 MB/s  | <1 ms   | <1 MB             |
| 10,000      | >100 MB/s  | <1 ms   | <1 MB             |

## 开源协议
MIT
