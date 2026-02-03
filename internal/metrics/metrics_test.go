package metrics

import (
	"testing"
)

func TestMetrics(t *testing.T) {
	// 重置全局指标（简单起见，不真正重置，只检查增量）
	activeBefore, _, _, _, _, _ := GetGlobalMetrics().GetStats()
	
	IncActiveConns()
	IncActiveConns()
	DecActiveConns()
	
	activeAfter, _, _, _, _, _ := GetGlobalMetrics().GetStats()
	if activeAfter != activeBefore+1 {
		t.Errorf("Expected active conns to increase by 1, got %d -> %d", activeBefore, activeAfter)
	}

	IncErrorCount()
	IncAuthFailures()
	AddBytesTransferred(1024)

	_, _, bytes, errors, authFails, uptime := GetGlobalMetrics().GetStats()
	if bytes < 1024 {
		t.Errorf("Expected at least 1024 bytes, got %d", bytes)
	}
	if errors < 1 {
		t.Error("Expected error count to be at least 1")
	}
	if authFails < 1 {
		t.Error("Expected auth failure count to be at least 1")
	}
	if uptime <= 0 {
		t.Error("Expected uptime to be positive")
	}
}
