package config

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
	Debug   bool   `mapstructure:"debug"`
}

// IsProduction 判断是否为生产环境
func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment 判断是否为开发环境
func (c *AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}
