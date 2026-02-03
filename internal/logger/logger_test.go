package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	Init("debug", buf)

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	output := buf.String()
	if !strings.Contains(output, "[DEBUG] debug message") {
		t.Error("Debug message not found in log")
	}
	if !strings.Contains(output, "[INFO] info message") {
		t.Error("Info message not found in log")
	}
	if !strings.Contains(output, "[WARN] warn message") {
		t.Error("Warn message not found in log")
	}
	if !strings.Contains(output, "[ERROR] error message") {
		t.Error("Error message not found in log")
	}

	// 测试级别过滤
	buf.Reset()
	Init("warn", buf)
	Debug("should not show")
	Info("should not show")
	Warn("should show")
	
	output = buf.String()
	if strings.Contains(output, "should not show") {
		t.Error("Log level filtering failed")
	}
	if !strings.Contains(output, "should show") {
		t.Error("Warn message should be shown")
	}
}
