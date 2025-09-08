package config

// JWTConfig JWT配置
type JWTConfig struct {
	Secret              string `mapstructure:"secret"`
	ExpiresHours        int    `mapstructure:"expires_hours"`
	RefreshExpiresHours int    `mapstructure:"refresh_expires_hours"`
	Issuer              string `mapstructure:"issuer"`
}
