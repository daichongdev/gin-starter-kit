package main

import (
	"context"
	"fmt"
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/model"
	"gin-demo/pkg/cron"
	"gin-demo/pkg/logger"
	"gin-demo/router"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

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

	// 初始化定时任务
	cron.Init()

	// 初始化路由
	r := router.SetupRouter()

	// 高性能HTTP服务器配置
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		logger.Info("Server starting",
			zap.String("address", srv.Addr),
			zap.Duration("read_timeout", cfg.Server.ReadTimeout),
			zap.Duration("write_timeout", cfg.Server.WriteTimeout),
		)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// 关闭数据库连接
	if err := database.CloseDB(); err != nil {
		logger.Error("Error closing database", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
