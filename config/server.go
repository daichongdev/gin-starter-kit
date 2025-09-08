package config

import (
	"time"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Limit  int           `mapstructure:"limit"`
	Window time.Duration `mapstructure:"window"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	MaxRequestSize  int64         `mapstructure:"max_request_size"`
	RateLimit       struct {
		Global RateLimitConfig `mapstructure:"global"`
		Auth   RateLimitConfig `mapstructure:"auth"`
	} `mapstructure:"rate_limit"`
}
