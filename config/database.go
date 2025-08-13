package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL *MySQLConfig `mapstructure:"mysql"`
	Redis *RedisConfig `mapstructure:"redis"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	DialTimeout     time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxRetries      int           `mapstructure:"max_retries"`
	PoolTimeout     time.Duration `mapstructure:"pool_timeout"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	DialTimeout     time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
}

// 设置数据库默认值
func setDatabaseDefaults() {
	// MySQL 默认值 - 不再从环境变量读取
	viper.SetDefault("database.mysql.host", "127.0.0.1")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.username", "root")
	viper.SetDefault("database.mysql.password", "")
	viper.SetDefault("database.mysql.database", "daka_dev")
	viper.SetDefault("database.mysql.charset", "utf8mb4")
	viper.SetDefault("database.mysql.parse_time", true)
	viper.SetDefault("database.mysql.loc", "Local")
	viper.SetDefault("database.mysql.max_idle_conns", 20)
	viper.SetDefault("database.mysql.max_open_conns", 100)
	viper.SetDefault("database.mysql.conn_max_lifetime", "3600s")
	viper.SetDefault("database.mysql.conn_max_idle_time", "1800s")
	viper.SetDefault("database.mysql.dial_timeout", "10s")
	viper.SetDefault("database.mysql.read_timeout", "30s")
	viper.SetDefault("database.mysql.write_timeout", "30s")

	// Redis 默认值
	viper.SetDefault("database.redis.host", "localhost")
	viper.SetDefault("database.redis.port", 6379)
	viper.SetDefault("database.redis.password", "")
	viper.SetDefault("database.redis.db", 0)
	viper.SetDefault("database.redis.pool_size", 20)
	viper.SetDefault("database.redis.min_idle_conns", 5)
	viper.SetDefault("database.redis.max_retries", 3)
	viper.SetDefault("database.redis.pool_timeout", "30s")
	viper.SetDefault("database.redis.conn_max_idle_time", "1800s")
	viper.SetDefault("database.redis.dial_timeout", "10s")
	viper.SetDefault("database.redis.read_timeout", "10s")
	viper.SetDefault("database.redis.write_timeout", "10s")
}

// GetDSN 获取MySQL DSN
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc,
		c.DialTimeout, c.ReadTimeout, c.WriteTimeout)
}

// GetAddr 获取Redis地址
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
