package config

import (
	"time"

	"github.com/spf13/viper"
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

// 设置服务器默认值
func setServerDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.shutdown_timeout", "10s")
	viper.SetDefault("server.max_header_bytes", 1048576) // 1MB
	viper.SetDefault("server.max_request_size", 10485760) // 10MB
	// 限流默认值
	viper.SetDefault("server.rate_limit.global.limit", 100)
	viper.SetDefault("server.rate_limit.global.window", "1m")
	viper.SetDefault("server.rate_limit.auth.limit", 20)
	viper.SetDefault("server.rate_limit.auth.window", "1m")
}
