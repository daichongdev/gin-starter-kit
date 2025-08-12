package config

import "github.com/spf13/viper"

// 日志配置
type LogConfig struct {
	Level            string            `mapstructure:"level"`
	Format           string            `mapstructure:"format"`
	EnableColor      bool              `mapstructure:"enable_color"`
	EnableCaller     bool              `mapstructure:"enable_caller"`
	EnableStacktrace bool              `mapstructure:"enable_stacktrace"`
	Console          *ConsoleLogConfig `mapstructure:"console"`
	File             *FileLogConfig    `mapstructure:"file"`
	ErrorFile        *FileLogConfig    `mapstructure:"error_file"`
	Database         *DatabaseLogConfig `mapstructure:"database"`
	Access           *AccessLogConfig  `mapstructure:"access"`
}

// 控制台日志配置
type ConsoleLogConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	EnableColor bool   `mapstructure:"enable_color"`
}

// 文件日志配置
type FileLogConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// 数据库日志配置
type DatabaseLogConfig struct {
	Enabled              bool   `mapstructure:"enabled"`
	Level                string `mapstructure:"level"`
	SlowThreshold        int    `mapstructure:"slow_threshold"`
	LogAllSQL            bool   `mapstructure:"log_all_sql"`
	IgnoreRecordNotFound bool   `mapstructure:"ignore_record_not_found"`
	Filename             string `mapstructure:"filename"`
	MaxSize              int    `mapstructure:"max_size"`
	MaxBackups           int    `mapstructure:"max_backups"`
	MaxAge               int    `mapstructure:"max_age"`
	Compress             bool   `mapstructure:"compress"`
}

// 访问日志配置
type AccessLogConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Format     string `mapstructure:"format"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// 设置日志默认值
func setLogDefaults() {
	// 全局日志默认值
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.enable_color", true)
	viper.SetDefault("log.enable_caller", false)
	viper.SetDefault("log.enable_stacktrace", false)

	// 控制台日志默认值
	viper.SetDefault("log.console.enabled", true)
	viper.SetDefault("log.console.level", "info")
	viper.SetDefault("log.console.format", "text")
	viper.SetDefault("log.console.enable_color", true)

	// 文件日志默认值
	viper.SetDefault("log.file.enabled", true)
	viper.SetDefault("log.file.level", "info")
	viper.SetDefault("log.file.format", "json")
	viper.SetDefault("log.file.filename", "./logs/app.log")
	viper.SetDefault("log.file.max_size", 10)
	viper.SetDefault("log.file.max_backups", 30)
	viper.SetDefault("log.file.max_age", 7)
	viper.SetDefault("log.file.compress", true)

	// 错误文件日志默认值
	viper.SetDefault("log.error_file.enabled", true)
	viper.SetDefault("log.error_file.level", "error")
	viper.SetDefault("log.error_file.format", "json")
	viper.SetDefault("log.error_file.filename", "./logs/error.log")
	viper.SetDefault("log.error_file.max_size", 10)
	viper.SetDefault("log.error_file.max_backups", 20)
	viper.SetDefault("log.error_file.max_age", 7)
	viper.SetDefault("log.error_file.compress", true)

	// 数据库日志默认值
	viper.SetDefault("log.database.enabled", true)
	viper.SetDefault("log.database.level", "info")
	viper.SetDefault("log.database.slow_threshold", 200)
	viper.SetDefault("log.database.log_all_sql", true)
	viper.SetDefault("log.database.ignore_record_not_found", true)
	viper.SetDefault("log.database.filename", "./logs/database.log")
	viper.SetDefault("log.database.max_size", 10)
	viper.SetDefault("log.database.max_backups", 30)
	viper.SetDefault("log.database.max_age", 7)
	viper.SetDefault("log.database.compress", true)

	// 访问日志默认值
	viper.SetDefault("log.access.enabled", true)
	viper.SetDefault("log.access.format", "json")
	viper.SetDefault("log.access.filename", "./logs/access.log")
	viper.SetDefault("log.access.max_size", 10)
	viper.SetDefault("log.access.max_backups", 30)
	viper.SetDefault("log.access.max_age", 7)
	viper.SetDefault("log.access.compress", true)
}