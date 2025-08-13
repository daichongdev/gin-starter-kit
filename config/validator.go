package config

import (
	"errors"
	"fmt"
)

// Validate 验证配置
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config error: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config error: %w", err)
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config error: %w", err)
	}

	return nil
}

// Validate 验证服务器配置
func (c *ServerConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("invalid port number")
	}
	return nil
}

// Validate 验证数据库配置
func (c *DatabaseConfig) Validate() error {
	if c.MySQL.Host == "" {
		return errors.New("mysql host is required")
	}
	if c.MySQL.Database == "" {
		return errors.New("mysql database name is required")
	}
	return nil
}

// Validate 验证JWT配置
func (c *JWTConfig) Validate() error {
	if len(c.Secret) < 32 {
		return errors.New("jwt secret must be at least 32 characters")
	}
	if c.ExpiresHours <= 0 {
		return errors.New("jwt expires hours must be positive")
	}
	return nil
}
