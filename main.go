package main

import (
	"context"
	"errors"
	"fmt"
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/model"
	"gin-demo/pkg/cron"
	"gin-demo/pkg/logger"
	"gin-demo/pkg/middleware"
	"gin-demo/pkg/queue"
	"gin-demo/router"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 设置Go运行时参数以提高性能
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	config.InitConfig()
	cfg := config.GetConfig()

	// 初始化日志系统
	if err := logger.InitLogger(cfg); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// 生产环境设置
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Info("Application starting",
		zap.String("app_name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.Int("cpu_cores", runtime.NumCPU()),
	)

	// 初始化数据库
	database.InitDB()

	// 自动迁移所有注册的模型
	if err := model.Registry.AutoMigrate(); err != nil {
		logger.Fatal("自动迁移失败", zap.Error(err))
	}

	// 初始化队列管理器
	if err := queue.InitManager(cfg.Queue); err != nil {
		logger.Fatal("Failed to initialize queue manager", zap.Error(err))
	}

	// 注册队列处理器
	if err := queue.RegisterQueueHandlers(queue.GetManager(), cfg.Queue); err != nil {
		logger.Fatal("Failed to register queue handlers", zap.Error(err))
	}

	// 初始化定时任务
	cron.Init()

	// 初始化路由
	r := router.SetupRouter()

	// 添加HTTP/2推送支持检测中间件
	r.Use(func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			c.Header("X-HTTP2-Push", "supported")
			c.Header("X-HTTP2-Protocol", "h2")
		}
		c.Next()
	})

	// 从配置文件读取HTTP/2配置
	h2Config := cfg.Server.HTTP2
	h2s := &http2.Server{
		MaxConcurrentStreams:         h2Config.MaxConcurrentStreams,
		MaxReadFrameSize:             h2Config.MaxReadFrameSize,
		IdleTimeout:                  h2Config.IdleTimeout,
		MaxUploadBufferPerConnection: h2Config.MaxUploadBufferPerConnection,
		MaxUploadBufferPerStream:     h2Config.MaxUploadBufferPerStream,
		PermitProhibitedCipherSuites: h2Config.PermitProhibitedCipherSuites,
	}

	// 根据配置决定是否使用 h2c
	var handler http.Handler = r
	var protocolInfo string

	if cfg.Server.EnableH2C {
		handler = h2c.NewHandler(r, h2s)
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

	// 高性能HTTP服务器配置
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           handler,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       h2Config.IdleTimeout,
		MaxHeaderBytes:    cfg.Server.MaxHeaderBytes,
		ReadHeaderTimeout: h2Config.ReadHeaderTimeout,
	}

	// 启用HTTP/2（对于 HTTPS 连接）
	if err := http2.ConfigureServer(srv, h2s); err != nil {
		logger.Error("Failed to configure HTTP/2", zap.Error(err))
	} else {
		logger.Info("HTTP/2 server configured successfully",
			zap.String("protocol", protocolInfo),
			zap.Uint32("max_concurrent_streams", h2Config.MaxConcurrentStreams),
			zap.Uint32("max_read_frame_size", h2Config.MaxReadFrameSize),
			zap.Duration("idle_timeout", h2Config.IdleTimeout),
			zap.Int32("max_upload_buffer_per_connection", h2Config.MaxUploadBufferPerConnection),
			zap.Int32("max_upload_buffer_per_stream", h2Config.MaxUploadBufferPerStream),
		)
	}

	// 启动服务器
	go func() {
		logger.Info("Server starting",
			zap.String("address", srv.Addr),
			zap.String("protocol", protocolInfo),
			zap.Duration("read_timeout", cfg.Server.ReadTimeout),
			zap.Duration("write_timeout", cfg.Server.WriteTimeout),
			zap.Duration("read_header_timeout", h2Config.ReadHeaderTimeout),
			zap.Bool("h2c_enabled", cfg.Server.EnableH2C),
			zap.String("gin_mode", gin.Mode()),
		)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 停止定时任务
	cron.Stop()

	// 关闭队列管理器
	if queueManager := queue.GetManager(); queueManager != nil {
		if queueErr := queueManager.Close(); queueErr != nil {
			logger.Error("Error closing queue manager", zap.Error(queueErr))
		}
	}

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// 关闭访问日志异步处理器（确保日志 flush 完成）
	middleware.ShutdownAccessLogger()

	// 关闭数据库连接
	if err := database.CloseDB(); err != nil {
		logger.Error("Error closing database", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
