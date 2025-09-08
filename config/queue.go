package config

import (
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
