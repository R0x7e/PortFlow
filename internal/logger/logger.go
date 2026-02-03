package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Level 定义日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var (
	currentLevel = InfoLevel
	logger       = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

// Init 初始化日志模块
func Init(levelStr string, out io.Writer) {
	if out == nil {
		out = os.Stdout
	}
	logger = log.New(out, "", log.LstdFlags|log.Lshortfile)
	
	switch strings.ToLower(levelStr) {
	case "debug":
		currentLevel = DebugLevel
	case "info":
		currentLevel = InfoLevel
	case "warn":
		currentLevel = WarnLevel
	case "error":
		currentLevel = ErrorLevel
	default:
		currentLevel = InfoLevel
	}
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	if currentLevel <= DebugLevel {
		logger.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
	}
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if currentLevel <= InfoLevel {
		logger.Output(2, fmt.Sprintf("[INFO] "+format, v...))
	}
}

// Warn 记录警告日志
func Warn(format string, v ...interface{}) {
	if currentLevel <= WarnLevel {
		logger.Output(2, fmt.Sprintf("[WARN] "+format, v...))
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if currentLevel <= ErrorLevel {
		logger.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
	}
}
