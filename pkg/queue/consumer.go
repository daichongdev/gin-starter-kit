package queue

import (
	"gin-demo/config"
	"time"

	"github.com/adjust/rmq/v5"
	"go.uber.org/zap"
)

// Consumer RMQ消费者实现
type Consumer struct {
	handler MessageHandler
	logger  *zap.Logger
	config  *config.QueueConfig
	manager *Manager // 添加Manager引用
}

// NewConsumer 创建新的消费者
func NewConsumer(handler MessageHandler, logger *zap.Logger, config *config.QueueConfig, manager *Manager) *Consumer {
	return &Consumer{
		handler: handler,
		logger:  logger,
		config:  config,
		manager: manager,
	}
}

// Consume 实现rmq.Consumer接口
func (c *Consumer) Consume(delivery rmq.Delivery) {
	// 解析消息
	message, err := FromJSON([]byte(delivery.Payload()))
	if err != nil {
		c.logger.Error("Failed to parse message",
			zap.String("queue", c.handler.GetQueueName()),
			zap.String("payload", delivery.Payload()),
			zap.Error(err))
		delivery.Reject()
		return
	}

	// 处理消息
	if err := c.processMessage(message, delivery); err != nil {
		c.logger.Error("Failed to process message",
			zap.String("queue", c.handler.GetQueueName()),
			zap.String("message_id", message.ID),
			zap.Error(err))

		// 检查是否需要重试
		if c.shouldRetry(message) {
			message.RetryCount++
			c.retryMessage(message, delivery)
		} else {
			delivery.Reject()
		}
		return
	}

	// 确认消息处理成功
	delivery.Ack()
	c.logger.Debug("Message processed successfully",
		zap.String("queue", c.handler.GetQueueName()),
		zap.String("message_id", message.ID))
}

// processMessage 处理消息
func (c *Consumer) processMessage(message *Message, delivery rmq.Delivery) error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Panic occurred while processing message",
				zap.String("queue", c.handler.GetQueueName()),
				zap.String("message_id", message.ID),
				zap.Any("panic", r))
		}
	}()

	return c.handler.Handle(message)
}

// shouldRetry 判断是否应该重试
func (c *Consumer) shouldRetry(message *Message) bool {
	return message.RetryCount < c.config.RMQ.RetryLimit
}

// retryMessage 重试消息
func (c *Consumer) retryMessage(message *Message, delivery rmq.Delivery) {
	// 延迟重试
	go func() {
		time.Sleep(c.config.RMQ.RetryDelay)

		// 通过Manager重新发布消息
		if err := c.manager.Publish(c.handler.GetQueueName(), message); err != nil {
			c.logger.Error("Failed to retry message",
				zap.String("message_id", message.ID),
				zap.Error(err))
		} else {
			c.logger.Info("Message retried successfully",
				zap.String("message_id", message.ID),
				zap.Int("retry_count", message.RetryCount))
		}
	}()

	delivery.Ack() // 确认当前消息，避免重复处理
}
