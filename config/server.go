package config

import (
	"time"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Limit  int           `mapstructure:"limit"`
	Window time.Duration `mapstructure:"window"`
}

// HTTP2Config HTTP/2配置
type HTTP2Config struct {
	MaxConcurrentStreams         uint32        `mapstructure:"max_concurrent_streams"`
	MaxReadFrameSize             uint32        `mapstructure:"max_read_frame_size"`
	IdleTimeout                  time.Duration `mapstructure:"idle_timeout"`
	MaxUploadBufferPerConnection int32         `mapstructure:"max_upload_buffer_per_connection"`
	MaxUploadBufferPerStream     int32         `mapstructure:"max_upload_buffer_per_stream"`
	PermitProhibitedCipherSuites bool          `mapstructure:"permit_prohibited_cipher_suites"`
	ReadHeaderTimeout            time.Duration `mapstructure:"read_header_timeout"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	EnableH2C       bool          `mapstructure:"enable_h2c"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	MaxRequestSize  int64         `mapstructure:"max_request_size"`
	HTTP2           HTTP2Config   `mapstructure:"http2"`
	RateLimit       struct {
		Global RateLimitConfig `mapstructure:"global"`
		Auth   RateLimitConfig `mapstructure:"auth"`
	} `mapstructure:"rate_limit"`
}
