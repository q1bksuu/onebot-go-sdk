package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ServerConfig 通用服务器配置字段.
type ServerConfig struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// BaseServer 封装通用的 http.Server 启动与关闭逻辑.
type BaseServer struct {
	Srv *http.Server
}

// NewBaseServer 创建基础服务器.
func NewBaseServer(cfg ServerConfig, handler http.Handler) *BaseServer {
	return &BaseServer{
		Srv: &http.Server{
			Addr:              cfg.Addr,
			Handler:           handler,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
	}
}

// Start 启动服务器并等待上下文取消.
// onShutdown 是可选的回调，在调用 http.Server.Shutdown 之前执行.
func (s *BaseServer) Start(ctx context.Context, onShutdown func(context.Context) error) error {
	errCh := make(chan error, 1)

	go func() {
		err := s.Srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if onShutdown != nil {
			err := onShutdown(shutdownCtx)
			if err != nil {
				return err
			}
		}

		err := s.Srv.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		return nil
	case err := <-errCh:
		return fmt.Errorf("server listen failed: %w", err)
	}
}

// Shutdown 直接关闭服务器.
func (s *BaseServer) Shutdown(ctx context.Context) error {
	err := s.Srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
