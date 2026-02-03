package forwarder

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"portflow/internal/metrics"
	"strings"
	"time"
)

// copyData 在两个连接之间双向拷贝数据
func (f *Forwarder) copyData(dst io.Writer, src io.Reader, timeout time.Duration, errChan chan error) {
	buf := make([]byte, 32*1024) // 32KB 缓冲区
	for {
		// 如果 src 是 net.Conn，设置超时
		if conn, ok := src.(net.Conn); ok && timeout > 0 {
			conn.SetReadDeadline(time.Now().Add(timeout))
		}
		// 如果 dst 是 net.Conn，设置超时
		if conn, ok := dst.(net.Conn); ok && timeout > 0 {
			conn.SetWriteDeadline(time.Now().Add(timeout))
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
func (f *Forwarder) doAuth(conn net.Conn, reader *bufio.Reader) (bool, error) {
	// 设置认证超时
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// 尝试读取第一行数据（Peek 不会消耗缓冲区数据，但 ReadString 会）
	// 为了简单起见，我们直接 ReadString，之后在 handleConnection 中使用这个 reader
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	trimmedLine := strings.TrimSpace(line)

	// 1. 检查是否是 HTTP 请求并包含 Authorization 头部
	if strings.Contains(line, "Authorization: Basic ") {
		parts := strings.Split(line, " ")
		for i, p := range parts {
			if strings.HasPrefix(p, "Basic") && i+1 < len(parts) {
				if f.auth.Verify(parts[i+1]) {
					return true, nil
				}
			}
		}
	}

	// 2. 尝试直接验证整行（兼容某些简单客户端直接发送 base64）
	if f.auth.Verify(trimmedLine) {
		return true, nil
	}

	// 3. 如果是 HTTP 请求但认证失败，发送 401 挑战
	if strings.Contains(line, "HTTP/") {
		challenge := fmt.Sprintf("HTTP/1.1 401 Unauthorized\r\n" +
			"WWW-Authenticate: Basic realm=\"%s\"\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: 26\r\n\r\n" +
			"Authentication Required\r\n", f.config.AuthRealm)
		conn.Write([]byte(challenge))
	}

	return false, nil
}


