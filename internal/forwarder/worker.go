package forwarder

import (
	"bufio"
	"net"
	"portflow/internal/metrics"
	"strings"
	"time"
)

// copyData 在两个连接之间双向拷贝数据
func (f *Forwarder) copyData(dst, src net.Conn, timeout time.Duration, errChan chan error) {
	buf := make([]byte, 32*1024) // 32KB 缓冲区
	for {
		// 每次读写前更新超时时间，防止死连接
		if timeout > 0 {
			src.SetReadDeadline(time.Now().Add(timeout))
			dst.SetWriteDeadline(time.Now().Add(timeout))
		}

		n, err := src.Read(buf)
		if n > 0 {
			_, wErr := dst.Write(buf[:n])
			if wErr != nil {
				errChan <- wErr
				return
			}
			metrics.AddBytesTransferred(int64(n))
		}
		if err != nil {
			errChan <- err
			return
		}
	}
}

// doAuth 执行认证逻辑
func (f *Forwarder) doAuth(conn net.Conn) bool {
	// 设置认证超时
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	
	reader := bufio.NewReader(conn)
	// 读取第一行或第一段数据
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	// 检查是否是 HTTP Basic 认证格式
	// 例如: "Proxy-Authorization: Basic <base64>" 或直接 "<base64>"
	line = strings.TrimSpace(line)
	
	// 尝试作为 HTTP Header 处理
	if strings.Contains(line, "Basic ") {
		parts := strings.Split(line, " ")
		for i, p := range parts {
			if strings.HasPrefix(p, "Basic") && i+1 < len(parts) {
				return f.auth.Verify(parts[i+1])
			}
		}
	}

	// 尝试直接作为 Base64 处理
	return f.auth.Verify(line)
}
