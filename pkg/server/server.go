package server

import (
	"context"
	"errors"
	"fmt"
	"gin-demo/config"
	"gin-demo/pkg/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	srv    *http.Server
	config *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) SetupHTTP2(handler http.Handler) http.Handler {
	h2Config := s.config.Server.HTTP2
	h2s := &http2.Server{
		MaxConcurrentStreams:         h2Config.MaxConcurrentStreams,
		MaxReadFrameSize:             h2Config.MaxReadFrameSize,
		IdleTimeout:                  h2Config.IdleTimeout,
		MaxUploadBufferPerConnection: h2Config.MaxUploadBufferPerConnection,
		MaxUploadBufferPerStream:     h2Config.MaxUploadBufferPerStream,
		PermitProhibitedCipherSuites: h2Config.PermitProhibitedCipherSuites,
	}

	var finalHandler http.Handler = handler
	var protocolInfo string

	if s.config.Server.EnableH2C {
		finalHandler = h2c.NewHandler(handler, h2s)
		protocolInfo = "HTTP/2 Cleartext (h2c)"
		logger.Info("HTTP/2 Cleartext (h2c) enabled",
			zap.Uint32("max_concurrent_streams", h2Config.MaxConcurrentStreams),
			zap.Duration("idle_timeout", h2Config.IdleTimeout),
			zap.Uint32("max_read_frame_size", h2Config.MaxReadFrameSize),
		)
	} else {
		protocolInfo = "HTTP/2 over TLS"
		logger.Info("HTTP/2 over TLS only (h2c disabled)")
	}

	// 创建HTTP服务器
	s.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:           finalHandler,
		ReadTimeout:       s.config.Server.ReadTimeout,
		WriteTimeout:      s.config.Server.WriteTimeout,
		IdleTimeout:       h2Config.IdleTimeout,
		MaxHeaderBytes:    s.config.Server.MaxHeaderBytes,
		ReadHeaderTimeout: h2Config.ReadHeaderTimeout,
	}

	// 启用HTTP/2
	if err := http2.ConfigureServer(s.srv, h2s); err != nil {
		logger.Error("Failed to configure HTTP/2", zap.Error(err))
	} else {
		logger.Info("HTTP/2 server configured successfully",
			zap.String("protocol", protocolInfo),
			zap.Uint32("max_concurrent_streams", h2Config.MaxConcurrentStreams),
			zap.Uint32("max_read_frame_size", h2Config.MaxReadFrameSize),
			zap.Duration("idle_timeout", h2Config.IdleTimeout),
		)
	}

	return finalHandler
}

func (s *Server) Start() error {
	logger.Info("Server starting",
		zap.String("address", s.srv.Addr),
		zap.Duration("read_timeout", s.config.Server.ReadTimeout),
		zap.Duration("write_timeout", s.config.Server.WriteTimeout),
		zap.Bool("h2c_enabled", s.config.Server.EnableH2C),
		zap.String("gin_mode", gin.Mode()),
	)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")
	return s.srv.Shutdown(ctx)
}
