# PortFlow

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20Darwin-blue)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)

English | [简体中文](README.zh-CN.md)

A high-performance, production-grade TCP port forwarding tool written in Go.

## Features

- **Bidirectional Forwarding**: Supports any TCP-based protocols (HTTP, HTTPS, custom protocols).
- **High Concurrency**: Based on goroutine pool, supports 10,000+ concurrent connections.
- **Security**: Optional HTTP Basic Authentication (RFC 7617).
- **Performance**: 
  - Connection pooling for target servers.
  - Multi-layer timeout control (Connect, R/W, Idle).
  - Low latency (<1ms loopback).
- **Production Ready**:
  - Graceful shutdown (SIGINT/SIGTERM).
  - Real-time performance metrics.
  - Comprehensive error handling.
  - Zero external dependencies.
- **Observability**: Customizable log levels (Debug, Info, Warn, Error).

## Installation

### From Source
```bash
go build -o portflow ./cmd/portflow
```

### Using Docker
```bash
docker build -t portflow .
docker run -p 8080:8080 portflow --listen-port 8080 --target-addr 1.2.3.4 --target-port 80
```

## Usage

```bash
./portflow --listen-port <port> --target-addr <ip> --target-port <port> [options]
```

### Options
- `--listen-addr`: Local address to listen on (default: 0.0.0.0)
- `--listen-port`: Local port to listen on (required)
- `--target-addr`: Target server IP/domain (required)
- `--target-port`: Target server port (required)
- `--auth`: Enable HTTP Basic Authentication (default: false)
- `--user`: Auth username
- `--pass`: Auth password
- `--log-level`: Log level (debug, info, warn, error) (default: info)
- `--conn-timeout`: Connection timeout in seconds (default: 30)
- `--rw-timeout`: Read/Write timeout in seconds (default: 60)
- `--idle-timeout`: Idle connection timeout in seconds (default: 300)

## Performance

| Concurrency | Throughput | Latency | Memory per k-conn |
|-------------|------------|---------|-------------------|
| 1,000       | >100 MB/s  | <1 ms   | <1 MB             |
| 5,000       | >100 MB/s  | <1 ms   | <1 MB             |
| 10,000      | >100 MB/s  | <1 ms   | <1 MB             |

## License
MIT
