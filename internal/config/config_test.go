package config

import (
	"flag"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 重置 flag
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	
	os.Args = []string{"cmd", "-listen-port", "8080", "-target-addr", "127.0.0.1", "-target-port", "9090"}
	
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.ListenPort != 8080 {
		t.Errorf("Expected listen port 8080, got %d", cfg.ListenPort)
	}
	if cfg.TargetAddr != "127.0.0.1" {
		t.Errorf("Expected target addr 127.0.0.1, got %s", cfg.TargetAddr)
	}
}

func TestInvalidConfig(t *testing.T) {
	// 测试无效端口
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-listen-port", "70000"}
	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error for invalid port, got nil")
	}

	// 测试缺失目标地址
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd", "-listen-port", "8080"}
	_, err = LoadConfig()
	if err == nil {
		t.Error("Expected error for missing target addr, got nil")
	}
}
