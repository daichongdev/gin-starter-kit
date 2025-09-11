package app

import (
	"context"
	"fmt"
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/model"
	"gin-demo/pkg/cron"
	"gin-demo/pkg/logger"
	"gin-demo/pkg/middleware"
	"gin-demo/pkg/queue"
	"gin-demo/pkg/server"
	"gin-demo/router"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	config *config.Config
	server *server.Server
}

func New() *App {
	return &App{}
}

func (a *App) Initialize() error {
	// 设置Go运行时参数
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置
	config.InitConfig()
	a.config = config.GetConfig()

	// 初始化日志系统
	if err := logger.InitLogger(a.config); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 生产环境设置
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Info("Application starting",
		zap.String("app_name", a.config.App.Name),
		zap.String("version", a.config.App.Version),
		zap.Int("cpu_cores", runtime.NumCPU()),
	)

	// 初始化数据库
	database.InitDB()

	// 自动迁移
	if err := model.Registry.AutoMigrate(); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	// 初始化队列
	if err := a.initQueue(); err != nil {
		return err
	}

	// 初始化定时任务
	cron.Init()

	// 创建服务器
	a.server = server.New(a.config)

	return nil
}

func (a *App) initQueue() error {
	if err := queue.InitManager(a.config.Queue); err != nil {
		return fmt.Errorf("failed to initialize queue manager: %w", err)
	}

	if err := queue.RegisterQueueHandlers(queue.GetManager(), a.config.Queue); err != nil {
		return fmt.Errorf("failed to register queue handlers: %w", err)
	}

	return nil
}

func (a *App) setupRouter() *gin.Engine {
	r := router.SetupRouter()

	// 添加HTTP/2推送支持检测中间件
	r.Use(func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			c.Header("X-HTTP2-Push", "supported")
			c.Header("X-HTTP2-Protocol", "h2")
		}
		c.Next()
	})

	return r
}

func (a *App) Run() error {
	// 设置路由
	r := a.setupRouter()

	// 配置HTTP/2
	a.server.SetupHTTP2(r)

	// 启动服务器
	go func() {
		if err := a.server.Start(); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	a.waitForShutdown()

	return nil
}

func (a *App) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.shutdown()
}

func (a *App) shutdown() {
	// 停止定时任务
	cron.Stop()

	// 关闭队列管理器
	if queueManager := queue.GetManager(); queueManager != nil {
		if err := queueManager.Close(); err != nil {
			logger.Error("Error closing queue manager", zap.Error(err))
		}
	}

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// 关闭访问日志异步处理器
	middleware.ShutdownAccessLogger()

	// 关闭数据库连接
	if err := database.CloseDB(); err != nil {
		logger.Error("Error closing database", zap.Error(err))
	}

	// 同步日志
	logger.Sync()

	logger.Info("Server exited gracefully")
}
