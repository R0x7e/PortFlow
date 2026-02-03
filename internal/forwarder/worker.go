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
func (f *Forwarder) doAuth(conn net.Conn, reader *bufio.Reader) (bool, []byte, error) {
	// 设置认证超时
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var readData []byte

	// 1. 读取第一行，判断是否是 HTTP 请求
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return false, nil, err
	}
	readData = append(readData, []byte(firstLine)...)

	isHTTP := strings.Contains(firstLine, "HTTP/")
	
	// 2. 检查第一行是否直接包含认证信息（兼容某些非标准客户端）
	trimmedFirstLine := strings.TrimSpace(firstLine)
	if f.auth.Verify(trimmedFirstLine) {
		return true, readData, nil
	}

	// 3. 如果是 HTTP 请求，继续读取头部寻找 Authorization 字段
	if isHTTP {
		// 循环读取每一行头部
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			readData = append(readData, []byte(line)...)
			
			line = strings.TrimSpace(line)
			if line == "" {
				// 头部读取结束，仍未找到认证信息
				break
			}

			if strings.HasPrefix(line, "Authorization: Basic ") {
				authBase64 := strings.TrimPrefix(line, "Authorization: Basic ")
				if f.auth.Verify(authBase64) {
					return true, readData, nil
				}
			}
			if strings.HasPrefix(line, "Proxy-Authorization: Basic ") {
				authBase64 := strings.TrimPrefix(line, "Proxy-Authorization: Basic ")
				if f.auth.Verify(authBase64) {
					return true, readData, nil
				}
			}
		}

		// 4. 仍未认证成功，发送 401 挑战
		challenge := fmt.Sprintf("HTTP/1.1 401 Unauthorized\r\n" +
			"WWW-Authenticate: Basic realm=\"%s\"\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: 26\r\n\r\n" +
			"Authentication Required\r\n", f.config.AuthRealm)
		conn.Write([]byte(challenge))
	}

	return false, readData, nil
}




