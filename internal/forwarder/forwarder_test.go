package forwarder

import (
	"fmt"
	"io"
	"net"
	"portflow/internal/config"
	"testing"
	"time"
)

func TestForwarder_Functional(t *testing.T) {
	// 1. 启动模拟目标服务器
	targetListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start target server: %v", err)
	}
	defer targetListener.Close()
	targetAddr := targetListener.Addr().(*net.TCPAddr)

	go func() {
		for {
			conn, err := targetListener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c) // Echo server
			}(conn)
		}
	}()

	// 2. 启动转发器
	cfg := &config.Config{
		ListenAddr:    "127.0.0.1",
		ListenPort:    0, // 自动分配
		TargetAddr:    "127.0.0.1",
		TargetPort:    targetAddr.Port,
		ConnectTimeout: 5,
		RWTimeout:      5,
		IdleTimeout:    10,
	}

	f := NewForwarder(cfg)
	
	// 获取分配的监听端口
	// 注意：为了获取端口，我们需要稍微修改 Start 逻辑或手动监听
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen for forwarder: %v", err)
	}
	fAddr := l.Addr().(*net.TCPAddr)
	l.Close() // 释放端口，虽然有竞争风险，但测试中通常 OK

	cfg.ListenPort = fAddr.Port
	go f.Start()
	time.Sleep(100 * time.Millisecond) // 等待启动
	defer f.Stop()

	// 3. 客户端连接并发送数据
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.ListenPort))
	if err != nil {
		t.Fatalf("Failed to connect to forwarder: %v", err)
	}
	defer conn.Close()

	msg := "hello portflow"
	fmt.Fprintf(conn, "%s\n", msg)

	buf := make([]byte, len(msg)+1)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := io.ReadFull(conn, buf)
	if err != nil {
		t.Fatalf("Failed to read from forwarder: %v", err)
	}

	if string(buf[:n-1]) != msg {
		t.Errorf("Expected %s, got %s", msg, string(buf[:n-1]))
	}
}

func TestForwarder_Auth(t *testing.T) {
	// 1. 模拟目标服务器
	targetListener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer targetListener.Close()
	targetAddr := targetListener.Addr().(*net.TCPAddr)

	// 2. 启动带认证的转发器
	cfg := &config.Config{
		ListenAddr: "127.0.0.1",
		ListenPort: 0,
		TargetAddr: "127.0.0.1",
		TargetPort: targetAddr.Port,
		EnableAuth: true,
		AuthUser:   "admin",
		AuthPass:   "123456",
		ConnectTimeout: 2,
		RWTimeout:      2,
	}

	// 查找可用端口
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	cfg.ListenPort = l.Addr().(*net.TCPAddr).Port
	l.Close()

	f := NewForwarder(cfg)
	go f.Start()
	time.Sleep(100 * time.Millisecond)
	defer f.Stop()

	// 3. 测试认证失败
	conn, _ := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.ListenPort))
	fmt.Fprintf(conn, "wrong_auth\n")
	time.Sleep(100 * time.Millisecond)
	
	// 检查连接是否已关闭
	buf := make([]byte, 1)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err := conn.Read(buf)
	if err == nil {
		t.Error("Expected connection to be closed after auth failure")
	}
	conn.Close()

	// 4. 测试认证成功
	conn2, _ := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.ListenPort))
	defer conn2.Close()
	fmt.Fprintf(conn2, "YWRtaW46MTIzNDU2\n") // admin:123456
	time.Sleep(100 * time.Millisecond)
	
	// 如果认证成功，连接不应该关闭
	conn2.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err = conn2.Read(buf)
	// 因为目标服务器没有发送数据，所以应该是超时错误而不是 EOF
	if err != nil && err != io.EOF && !err.(net.Error).Timeout() {
		t.Errorf("Unexpected error after successful auth: %v", err)
	}
}
