package queue

import (
	"context"
	"fmt"
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/pkg/logger"
	"sync"
	"time"

	"github.com/adjust/rmq/v5"
	"go.uber.org/zap"
)

// 全局队列管理器实例
var (
	globalManager *Manager
	managerOnce   sync.Once
)

// Manager 队列管理器
type Manager struct {
	config     *config.QueueConfig
	connection rmq.Connection
	queues     map[string]rmq.Queue
	handlers   map[string]MessageHandler
	logger     *zap.Logger
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// InitManager 初始化全局队列管理器
func InitManager(cfg *config.QueueConfig) error {
	var err error
	managerOnce.Do(func() {
		globalManager, err = NewManager(cfg)
	})
	return err
}

// GetManager 获取全局队列管理器实例
func GetManager() *Manager {
	return globalManager
}

// NewManager 创建队列管理器（复用现有Redis客户端）
func NewManager(cfg *config.QueueConfig) (*Manager, error) {
	// 获取现有的Redis客户端
	rdb := database.GetRedis()
	if rdb == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection test failed: %w", err)
	}

	// 创建RMQ连接（复用现有Redis客户端）
	connection, err := rmq.OpenConnectionWithRedisClient(cfg.RMQ.Tag, rdb, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open rmq connection: %w", err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	m := &Manager{
		config:     cfg,
		connection: connection,
		queues:     make(map[string]rmq.Queue),
		handlers:   make(map[string]MessageHandler),
		logger:     logger.Logger.Named("queue"),
		ctx:        ctx,
		cancel:     cancel,
	}

	return m, nil
}

// RegisterHandler 注册消息处理器
func (m *Manager) RegisterHandler(handler MessageHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	queueName := handler.GetQueueName()
	if _, exists := m.handlers[queueName]; exists {
		return fmt.Errorf("handler for queue %s already registered", queueName)
	}

	// 获取队列配置
	queueConfig, exists := m.config.Queues[queueName]
	if !exists {
		return fmt.Errorf("queue config for %s not found", queueName)
	}

	// 打开队列：优先使用配置的实际队列名（如 email_queue），为空则回退逻辑名（如 email）
	openName := queueConfig.Name
	if openName == "" {
		openName = queueName
	}
	queue, err := m.connection.OpenQueue(openName)
	if err != nil {
		return fmt.Errorf("failed to open queue %s: %w", openName, err)
	}

	// 启动消费者（这里设置prefetch limit和poll duration）
	if err := queue.StartConsuming(int64(queueConfig.PrefetchLimit), m.config.RMQ.PollDuration); err != nil {
		return fmt.Errorf("failed to start consuming queue %s: %w", queueName, err)
	}

	// 创建消费者
	consumer := NewConsumer(handler, m.logger, m.config, m)

	// 添加消费者
	for i := 0; i < queueConfig.NumConsumers; i++ {
		consumerName := fmt.Sprintf("%s-consumer-%d", queueName, i)
		if _, err := queue.AddConsumer(consumerName, consumer); err != nil {
			return fmt.Errorf("failed to add consumer %s: %w", consumerName, err)
		}
	}

	// 保存队列和处理器
	m.queues[queueName] = queue
	m.handlers[queueName] = handler

	m.logger.Info("Handler registered successfully",
		zap.String("queue", queueName),
		zap.Int("consumers", queueConfig.NumConsumers),
		zap.Int("prefetch_limit", queueConfig.PrefetchLimit))

	return nil
}

// Publish 发布消息
func (m *Manager) Publish(queueName string, message *Message) error {
	m.mu.RLock()
	queue, exists := m.queues[queueName]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("queue %s not found", queueName)
	}

	data, err := message.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := queue.Publish(string(data)); err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %w", queueName, err)
	}

	m.logger.Debug("Message published",
		zap.String("queue", queueName),
		zap.String("message_id", message.ID),
		zap.String("message_type", message.Type))

	return nil
}

// PublishData 发布数据（便捷方法）
func (m *Manager) PublishData(queueName string, msgType string, data interface{}) error {
	message, err := NewMessage(msgType, data)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return m.Publish(queueName, message)
}

// GetQueueStats 获取队列统计信息
func (m *Manager) GetQueueStats(queueName string) (rmq.Stats, error) {
	m.mu.RLock()
	_, exists := m.queues[queueName]
	m.mu.RUnlock()

	if !exists {
		return rmq.Stats{}, fmt.Errorf("queue %s not found", queueName)
	}

	// 使用 Connection.CollectStats 方法获取统计信息
	stats, err := m.connection.CollectStats([]string{queueName})
	if err != nil {
		return rmq.Stats{}, fmt.Errorf("failed to collect stats for queue %s: %w", queueName, err)
	}

	return stats, nil
}

// Close 关闭队列管理器
func (m *Manager) Close() error {
	m.cancel()

	m.mu.Lock()
	defer m.mu.Unlock()

	// 停止所有队列
	for name, queue := range m.queues {
		stoppedChan := queue.StopConsuming()
		// 等待队列停止消费
		select {
		case <-stoppedChan:
			m.logger.Info("Queue stopped consuming", zap.String("queue", name))
		case <-time.After(10 * time.Second):
			m.logger.Warn("Timeout waiting for queue to stop consuming", zap.String("queue", name))
		}
	}

	// 停止所有消费者（使用 Connection 的 StopAllConsuming 方法）
	finishedChan := m.connection.StopAllConsuming()

	// 等待所有消费者停止（可选，根据需要设置超时）
	select {
	case <-finishedChan:
		m.logger.Info("All consumers stopped successfully")
	case <-time.After(30 * time.Second):
		m.logger.Warn("Timeout waiting for consumers to stop")
	}

	// 清理映射
	m.queues = make(map[string]rmq.Queue)
	m.handlers = make(map[string]MessageHandler)

	m.logger.Info("Queue manager closed")
	return nil
}
