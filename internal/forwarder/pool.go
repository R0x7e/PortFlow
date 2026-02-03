package forwarder

import (
	"fmt"
	"net"
	"time"
)

// ConnPool 简单的目标连接池
type ConnPool struct {
	targetAddr string
	targetPort int
	timeout    time.Duration
}

// NewConnPool 创建目标连接池
func NewConnPool(addr string, port int, timeoutSec int) *ConnPool {
	return &ConnPool{
		targetAddr: addr,
		targetPort: port,
		timeout:    time.Duration(timeoutSec) * time.Second,
	}
}

// Get 获取一个到目标服务器的连接
func (p *ConnPool) Get() (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", p.targetAddr, p.targetPort)
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	if err != nil {
		return nil, err
	}
	
	// 设置初始超时
	// conn.SetDeadline(time.Now().Add(p.timeout))
	
	return conn, nil
}
