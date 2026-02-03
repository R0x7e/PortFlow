package metrics

import (
	"sync/atomic"
	"time"
)

// Metrics 存储运行时性能指标
type Metrics struct {
	ActiveConns    int64 // 当前活跃连接数
	TotalConns     int64 // 总连接数
	BytesTransferred int64 // 总传输字节数
	ErrorCount     int64 // 错误总数
	AuthFailures   int64 // 认证失败次数
	startTime      time.Time
}

var globalMetrics = &Metrics{
	startTime: time.Now(),
}

// GetGlobalMetrics 返回全局指标实例
func GetGlobalMetrics() *Metrics {
	return globalMetrics
}

// IncActiveConns 增加活跃连接数
func IncActiveConns() {
	atomic.AddInt64(&globalMetrics.ActiveConns, 1)
	atomic.AddInt64(&globalMetrics.TotalConns, 1)
}

// DecActiveConns 减少活跃连接数
func DecActiveConns() {
	atomic.AddInt64(&globalMetrics.ActiveConns, -1)
}

// AddBytesTransferred 增加传输字节数
func AddBytesTransferred(n int64) {
	atomic.AddInt64(&globalMetrics.BytesTransferred, n)
}

// IncErrorCount 增加错误计数
func IncErrorCount() {
	atomic.AddInt64(&globalMetrics.ErrorCount, 1)
}

// IncAuthFailures 增加认证失败计数
func IncAuthFailures() {
	atomic.AddInt64(&globalMetrics.AuthFailures, 1)
}

// GetUptime 返回运行时间
func (m *Metrics) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

// GetStats 返回当前统计信息的快照
func (m *Metrics) GetStats() (active, total, bytes, errors, authFails int64, uptime time.Duration) {
	active = atomic.LoadInt64(&m.ActiveConns)
	total = atomic.LoadInt64(&m.TotalConns)
	bytes = atomic.LoadInt64(&m.BytesTransferred)
	errors = atomic.LoadInt64(&m.ErrorCount)
	authFails = atomic.LoadInt64(&m.AuthFailures)
	uptime = time.Since(m.startTime)
	return
}
