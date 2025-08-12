package config

import "github.com/spf13/viper"

// 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Debug   bool   `mapstructure:"debug"`
}

// 设置应用默认值
func setAppDefaults() {
	viper.SetDefault("app.name", "gin-demo")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.debug", true)
}

// 判断是否为生产环境
func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}

// 判断是否为开发环境
func (c *AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}