package config

import (
	"log"
	"os"
	"strings"

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
	// 设置所有模块的默认值
	setAllDefaults()

	// 如果存在环境变量 test-config，则优先从该变量解析 YAML 配置
	if content := os.Getenv("test-config"); content != "" {
		viper.SetConfigType("yaml")
		if err := viper.ReadConfig(strings.NewReader(content)); err != nil {
			log.Printf("Error reading config from env 'test-config': %v", err)
		} else {
			log.Printf("Config loaded from env 'test-config'")
		}
	} else {
		// 未设置 test-config 时，保持原有文件查找逻辑（本地开发）
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		// 兼容容器路径与本地路径
		viper.AddConfigPath("/app/config")
		viper.AddConfigPath("/app")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reading config file: %v", err)
		} else {
			log.Printf("Using config file: %s", viper.ConfigFileUsed())
		}
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
