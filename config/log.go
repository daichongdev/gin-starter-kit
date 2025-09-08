package config

// LogConfig 日志配置
type LogConfig struct {
	Level            string             `mapstructure:"level"`
	Format           string             `mapstructure:"format"`
	EnableColor      bool               `mapstructure:"enable_color"`
	EnableCaller     bool               `mapstructure:"enable_caller"`
	EnableStacktrace bool               `mapstructure:"enable_stacktrace"`
	Console          *ConsoleLogConfig  `mapstructure:"console"`
	File             *FileLogConfig     `mapstructure:"file"`
	ErrorFile        *FileLogConfig     `mapstructure:"error_file"`
	Database         *DatabaseLogConfig `mapstructure:"database"`
	Access           *AccessLogConfig   `mapstructure:"access"`
}

// ConsoleLogConfig 控制台日志配置
type ConsoleLogConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	EnableColor bool   `mapstructure:"enable_color"`
}

// FileLogConfig 文件日志配置
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

// DatabaseLogConfig 数据库日志配置
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

// AccessLogConfig 访问日志配置
type AccessLogConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Format     string `mapstructure:"format"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}
