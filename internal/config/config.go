package config

import (
	"flag"
	"fmt"
	"os"
)

// Config 存储程序的配置信息
type Config struct {
	ListenAddr    string // 本地监听地址
	ListenPort    int    // 本地监听端口
	TargetAddr    string // 目标服务器地址
	TargetPort    int    // 目标服务器端口
	EnableAuth    bool   // 是否启用认证
	AuthUser      string // 认证用户名
	AuthPass      string // 认证密码
	LogLevel      string // 日志级别
	ConnectTimeout int    // 连接建立超时时间（秒）
	RWTimeout      int    // 读写超时时间（秒）
	IdleTimeout    int    // 空闲连接超时时间（秒）
}

// LoadConfig 从命令行参数加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.ListenAddr, "listen-addr", "0.0.0.0", "Local listen address")
	flag.IntVar(&cfg.ListenPort, "listen-port", 0, "Local listen port")
	flag.StringVar(&cfg.TargetAddr, "target-addr", "", "Target server address")
	flag.IntVar(&cfg.TargetPort, "target-port", 0, "Target server port")
	flag.BoolVar(&cfg.EnableAuth, "auth", false, "Enable HTTP Basic authentication")
	flag.StringVar(&cfg.AuthUser, "user", "", "Authentication username")
	flag.StringVar(&cfg.AuthPass, "pass", "", "Authentication password")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.IntVar(&cfg.ConnectTimeout, "conn-timeout", 30, "Connection timeout in seconds")
	flag.IntVar(&cfg.RWTimeout, "rw-timeout", 60, "Read/Write timeout in seconds")
	flag.IntVar(&cfg.IdleTimeout, "idle-timeout", 300, "Idle connection timeout in seconds")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of PortFlow:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// 校验必填参数
	if cfg.ListenPort <= 0 || cfg.ListenPort > 65535 {
		return nil, fmt.Errorf("invalid local listen port: %d", cfg.ListenPort)
	}
	if cfg.TargetAddr == "" {
		return nil, fmt.Errorf("target server address is required")
	}
	if cfg.TargetPort <= 0 || cfg.TargetPort > 65535 {
		return nil, fmt.Errorf("invalid target server port: %d", cfg.TargetPort)
	}
	if cfg.EnableAuth && (cfg.AuthUser == "" || cfg.AuthPass == "") {
		return nil, fmt.Errorf("username and password are required when authentication is enabled")
	}

	return cfg, nil
}
