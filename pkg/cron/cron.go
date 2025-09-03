package cron

import (
	"gin-demo/pkg/logger"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var scheduler *cron.Cron

// Init 初始化定时任务
func Init() {
	scheduler = cron.New()

	// 添加定时任务
	if err := addJobs(); err != nil {
		logger.Error("添加定时任务失败", zap.Error(err))
		return
	}

	// 启动调度器
	scheduler.Start()
	logger.Info("定时任务启动成功")
}

// Stop 停止定时任务
func Stop() {
	if scheduler != nil {
		scheduler.Stop()
		logger.Info("定时任务已停止")
	}
}

// addJobs 添加定时任务
func addJobs() error {
	// 每天凌晨2点清理日志
	if _, err := scheduler.AddFunc("0 2 * * *", cleanLogs); err != nil {
		logger.Error("添加清理日志任务失败", zap.Error(err))
		return err
	}

	// 每小时统计数据
	if _, err := scheduler.AddFunc("0 * * * *", statistics); err != nil {
		logger.Error("添加统计数据任务失败", zap.Error(err))
		return err
	}

	// 每5分钟健康检查
	if _, err := scheduler.AddFunc("*/5 * * * *", healthCheck); err != nil {
		logger.Error("添加健康检查任务失败", zap.Error(err))
		return err
	}

	return nil
}

// cleanLogs 清理日志任务
func cleanLogs() {
	logger.Info("开始清理日志文件")
	// 这里添加清理逻辑
	logger.Info("日志清理完成")
}

// statistics 统计任务
func statistics() {
	logger.Info("开始统计数据")
	// 这里添加统计逻辑
	logger.Info("数据统计完成")
}

// healthCheck 健康检查任务
func healthCheck() {
	logger.Info("系统健康检查", zap.String("time", time.Now().Format("15:04:05")))
}
