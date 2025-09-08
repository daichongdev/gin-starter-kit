package test

import (
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/pkg/logger"
	"gin-demo/pkg/queue"
	"sync"
	"testing"
)

var (
	initOnce    sync.Once
	cleanupOnce sync.Once
	initErr     error
	manager     *queue.Manager
)

// InitTestEnvironment 初始化测试环境（只执行一次）
func InitTestEnvironment() error {
	initOnce.Do(func() {
		// 初始化配置
		config.InitConfig()
		cfg := config.GetConfig()

		// 初始化日志系统
		if err := logger.InitLogger(cfg); err != nil {
			initErr = err
			return
		}

		// 初始化数据库
		database.InitDB()

		// 初始化队列管理器
		if err := queue.InitManager(cfg.Queue); err != nil {
			initErr = err
			return
		}

		// 获取管理器实例用于后续清理
		manager = queue.GetManager()

		// 注册队列处理器
		if err := queue.RegisterQueueHandlers(manager, cfg.Queue); err != nil {
			initErr = err
			return
		}
	})

	return initErr
}

// CleanupTestEnvironment 清理测试环境
func CleanupTestEnvironment() {
	cleanupOnce.Do(func() {
		// 关闭队列管理器
		if manager != nil {
			manager.Close()
		}

		// 同步日志
		logger.Sync()

		// 关闭数据库连接
		database.CloseDB()
	})
}

// SetupTest 为单个测试设置环境（推荐使用）
func SetupTest(t *testing.T) func() {
	err := InitTestEnvironment()
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %v", err)
	}

	// 返回清理函数
	return func() {
		CleanupTestEnvironment()
	}
}
