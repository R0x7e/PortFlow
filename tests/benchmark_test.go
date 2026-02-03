package tests

import (
	"fmt"
	"io"
	"net"
	"portflow/internal/config"
	"portflow/internal/forwarder"
	"testing"
	"time"
)

func BenchmarkForwarding(b *testing.B) {
	// 1. 启动模拟目标服务器
	targetListener, _ := net.Listen("tcp", "127.0.0.1:0")
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
				io.Copy(io.Discard, c) // Discard all data
			}(conn)
		}
	}()

	// 2. 启动转发器
	cfg := &config.Config{
		ListenAddr: "127.0.0.1",
		ListenPort: 0,
		TargetAddr: "127.0.0.1",
		TargetPort: targetAddr.Port,
		ConnectTimeout: 5,
		RWTimeout:      60,
	}

	// 查找可用端口
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	cfg.ListenPort = l.Addr().(*net.TCPAddr).Port
	l.Close()

	f := forwarder.NewForwarder(cfg)
	go f.Start()
	time.Sleep(100 * time.Millisecond)
	defer f.Stop()

	// 3. 性能测试
	data := make([]byte, 1024*1024) // 1MB data
	b.ResetTimer()
	b.SetBytes(int64(len(data)))

	b.RunParallel(func(pb *testing.PB) {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.ListenPort))
		if err != nil {
			return
		}
		defer conn.Close()

		for pb.Next() {
			_, err := conn.Write(data)
			if err != nil {
				break
			}
		}
	})
}
