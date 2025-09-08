package config

import (
	"github.com/spf13/viper"
	"time"
)

// QueueConfig 队列配置
type QueueConfig struct {
	RMQ    RMQConfig                  `mapstructure:"rmq" yaml:"rmq"`
	Queues map[string]QueueItemConfig `mapstructure:"queues" yaml:"queues"`
}

// RMQConfig RMQ配置
type RMQConfig struct {
	Tag             string        `mapstructure:"tag" yaml:"tag"`
	NumConsumers    int           `mapstructure:"num_consumers" yaml:"num_consumers"`
	PrefetchLimit   int           `mapstructure:"prefetch_limit" yaml:"prefetch_limit"`
	PollDuration    time.Duration `mapstructure:"poll_duration" yaml:"poll_duration"`
	ReportBatchSize int           `mapstructure:"report_batch_size" yaml:"report_batch_size"`
	RetryLimit      int           `mapstructure:"retry_limit" yaml:"retry_limit"`
	RetryDelay      time.Duration `mapstructure:"retry_delay" yaml:"retry_delay"`
}

// QueueItemConfig 单个队列配置
type QueueItemConfig struct {
	Name          string `mapstructure:"name" yaml:"name"`
	NumConsumers  int    `mapstructure:"num_consumers" yaml:"num_consumers"`
	PrefetchLimit int    `mapstructure:"prefetch_limit" yaml:"prefetch_limit"`
}

// setQueueDefaults 设置队列默认值
func setQueueDefaults() {
	// RMQ 默认配置
	viper.SetDefault("queue.rmq.tag", "gin-demo-queue")
	viper.SetDefault("queue.rmq.num_consumers", 10)
	viper.SetDefault("queue.rmq.prefetch_limit", 1000)
	viper.SetDefault("queue.rmq.poll_duration", "100ms")
	viper.SetDefault("queue.rmq.report_batch_size", 100)
	viper.SetDefault("queue.rmq.retry_limit", 3)
	viper.SetDefault("queue.rmq.retry_delay", "5s")

	// 预定义队列默认配置
	viper.SetDefault("queue.queues.email.name", "email_queue")
	viper.SetDefault("queue.queues.email.num_consumers", 5)
	viper.SetDefault("queue.queues.email.prefetch_limit", 100)

	viper.SetDefault("queue.queues.notification.name", "notification_queue")
	viper.SetDefault("queue.queues.notification.num_consumers", 3)
	viper.SetDefault("queue.queues.notification.prefetch_limit", 50)

	viper.SetDefault("queue.queues.data_processing.name", "data_processing_queue")
	viper.SetDefault("queue.queues.data_processing.num_consumers", 8)
	viper.SetDefault("queue.queues.data_processing.prefetch_limit", 200)
}
