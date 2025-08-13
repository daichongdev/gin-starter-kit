package config

import (
	"github.com/spf13/viper"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret              string `mapstructure:"secret"`
	ExpiresHours        int    `mapstructure:"expires_hours"`
	RefreshExpiresHours int    `mapstructure:"refresh_expires_hours"`
	Issuer              string `mapstructure:"issuer"`
}

// 设置JWT默认值
func setJWTDefaults() {
	// 直接设置默认值，不再从环境变量读取
	viper.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-this-in-production-2024")
	viper.SetDefault("jwt.expires_hours", 168)
	viper.SetDefault("jwt.refresh_expires_hours", 720)
	viper.SetDefault("jwt.issuer", "gin-demo")
}
