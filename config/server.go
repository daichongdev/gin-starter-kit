package config

import (
	"time"

	"github.com/spf13/viper"
)

// 服务器配置
type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
}

// 设置服务器默认值
func setServerDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.shutdown_timeout", "10s")
	viper.SetDefault("server.max_header_bytes", 1048576) // 1MB
}