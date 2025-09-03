package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 主配置结构体
type Config struct {
	Server   *ServerConfig   `mapstructure:"server"`
	Database *DatabaseConfig `mapstructure:"database"`
	App      *AppConfig      `mapstructure:"app"`
	JWT      *JWTConfig      `mapstructure:"jwt"`
	Log      *LogConfig      `mapstructure:"log"`
	Queue    *QueueConfig    `mapstructure:"queue"` // 添加这一行
}

// Cfg 全局配置变量
var Cfg *Config

// InitConfig 初始化配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置所有模块的默认值
	setAllDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatal("Unable to decode config:", err)
	}

	log.Printf("Config loaded: %s v%s", Cfg.App.Name, Cfg.App.Version)
}

// 设置所有模块的默认值
func setAllDefaults() {
	setServerDefaults()
	setDatabaseDefaults()
	setAppDefaults()
	setJWTDefaults()
	setLogDefaults()
	setQueueDefaults() // 添加这一行
}

// GetConfig 获取配置
func GetConfig() *Config {
	return Cfg
}
