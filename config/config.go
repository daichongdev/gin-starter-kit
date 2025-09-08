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

	// 如果存在环境变量 go_app_config，则优先从该变量解析 YAML 配置
	if content := os.Getenv("go_app_config"); content != "" {
		viper.SetConfigType("yaml")
		if err := viper.ReadConfig(strings.NewReader(content)); err != nil {
			log.Printf("Error reading config from env 'go_app_config': %v", err)
		} else {
			log.Printf("Config loaded from env 'go_app_config'")
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("YAML config file not found, trying JSON format...")
			// yaml找不到，尝试json格式
			viper.SetConfigType("json")
			if err := viper.ReadInConfig(); err != nil {
				log.Fatalf("Config file not found in both YAML and JSON formats: %v", err)
			} else {
				log.Printf("Using JSON config file: %s", viper.ConfigFileUsed())
			}
		} else {
			log.Printf("Using YAML config file: %s", viper.ConfigFileUsed())
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
