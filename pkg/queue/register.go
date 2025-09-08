package queue

import (
	"fmt"
	"gin-demo/config"
	"gin-demo/pkg/logger"
)

// RegisterQueueHandlers 注册消息队列服务
func RegisterQueueHandlers(manager *Manager, cfg *config.QueueConfig) error {
	err := emailQueue(manager, cfg)
	if err != nil {
		return err
	}

	logger.Info("消息队列服务启动成功")
	return nil
}

// 邮件消息队列服务
func emailQueue(manager *Manager, cfg *config.QueueConfig) error {
	// 从配置文件读取邮件队列配置
	emailQueueConfig, exists := cfg.Queues[Email]
	if !exists {
		return fmt.Errorf("email queue configuration not found")
	}

	emailHandler := NewEmailHandler(
		Email,
		emailQueueConfig.NumConsumers,
		emailQueueConfig.PrefetchLimit,
	)
	if err := manager.RegisterHandler(emailHandler); err != nil {
		return fmt.Errorf("failed to register email handler: %w", err)
	}

	logger.Info("邮件队列服务注册成功")
	return nil
}
