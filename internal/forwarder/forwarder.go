package forwarder

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"portflow/internal/auth"
	"portflow/internal/config"
	"portflow/internal/logger"
	"portflow/internal/metrics"
	"sync"
	"time"
)

// Forwarder 核心转发器结构体
type Forwarder struct {
	config     *config.Config
	auth       *auth.Authenticator
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	connPool   *ConnPool
	workerSem  chan struct{} // 用于限制最大并发协程数
}

// NewForwarder 创建一个新的转发器实例
func NewForwarder(cfg *config.Config) *Forwarder {
	ctx, cancel := context.WithCancel(context.Background())
	f := &Forwarder{
		config:    cfg,
		ctx:       ctx,
		cancel:    cancel,
		workerSem: make(chan struct{}, 10000), // 默认最大 10000 并发
	}

	if cfg.EnableAuth {
		f.auth = auth.NewAuthenticator(cfg.AuthUser, cfg.AuthPass)
	}

	f.connPool = NewConnPool(cfg.TargetAddr, cfg.TargetPort, cfg.ConnectTimeout)
	return f
}

// Start 启动转发服务
func (f *Forwarder) Start() error {
	addr := fmt.Sprintf("%s:%d", f.config.ListenAddr, f.config.ListenPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Failed to listen on %s: %v", addr, err)
		return err
	}
	defer listener.Close()

	logger.Info("PortFlow started, listening on %s, forwarding to %s:%d", 
		addr, f.config.TargetAddr, f.config.TargetPort)

	// 监听 context 关闭信号
	go func() {
		<-f.ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-f.ctx.Done():
				return nil
			default:
				logger.Error("Accept error: %v", err)
				metrics.IncErrorCount()
				continue
			}
		}

		f.wg.Add(1)
		go f.handleConnection(conn)
	}
}

// Stop 停止转发服务并等待连接关闭
func (f *Forwarder) Stop() {
	logger.Info("Stopping PortFlow...")
	f.cancel()
	
	// 等待现有连接完成，设置超时时间
	done := make(chan struct{})
	go func() {
		f.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("All connections closed.")
	case <-time.After(10 * time.Second):
		logger.Warn("Shutdown timeout, some connections might be forced to close.")
	}
}

// handleConnection 处理单个传入连接
func (f *Forwarder) handleConnection(localConn net.Conn) {
	defer f.wg.Done()
	defer localConn.Close()

	metrics.IncActiveConns()
	defer metrics.DecActiveConns()

	// 限制并发
	f.workerSem <- struct{}{}
	defer func() { <-f.workerSem }()

	var localReader io.Reader = localConn

	// 认证检查
	if f.config.EnableAuth {
		reader := bufio.NewReader(localConn)
		ok, err := f.doAuth(localConn, reader)
		if !ok {
			metrics.IncAuthFailures()
			if err != nil && err != io.EOF {
				logger.Error("Auth error from %s: %v", localConn.RemoteAddr(), err)
			} else {
				logger.Warn("Authentication failed for connection from %s", localConn.RemoteAddr())
			}
			return
		}
		// 认证成功，确保已读取的数据（在 reader 缓冲区中）能被转发
		localReader = reader
	}

	// 获取目标连接
	targetConn, err := f.connPool.Get()
	if err != nil {
		logger.Error("Failed to connect to target: %v", err)
		metrics.IncErrorCount()
		return
	}
	defer targetConn.Close()

	// 设置读写超时
	rwTimeout := time.Duration(f.config.RWTimeout) * time.Second

	// 双向转发
	errChan := make(chan error, 2)
	go f.copyData(targetConn, localReader, rwTimeout, errChan)
	go f.copyData(localConn, targetConn, rwTimeout, errChan)

	// 等待任一方向结束或错误
	select {
	case <-f.ctx.Done():
		return
	case <-errChan:
		return
	}
}

