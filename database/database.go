package database

import (
	"context"
	"fmt"
	"gin-demo/config"
	"gin-demo/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func InitDB() {
	cfg := config.GetConfig()

	// 初始化MySQL
	initMySQL(cfg)

	// 初始化Redis
	initRedis(cfg)

	logger.Info("All databases initialized successfully")
}

func initMySQL(cfg *config.Config) {
	// 使用我们的GORM日志适配器
	gormLogger := logger.NewGormLogger(cfg.Log.Database)

	// 数据库连接
	dsn := cfg.Database.MySQL.GetDSN()
	logger.Info("Connecting to MySQL",
		logger.String("host", cfg.Database.MySQL.Host),
		logger.Int("port", cfg.Database.MySQL.Port),
		logger.String("database", cfg.Database.MySQL.Database),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		// 禁用外键约束检查（提高性能）
		DisableForeignKeyConstraintWhenMigrating: true,
		// 预编译语句缓存
		PrepareStmt: true,
	})
	if err != nil {
		logger.Fatal("Failed to connect to MySQL", logger.Err(err))
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", logger.Err(err))
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(cfg.Database.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.MySQL.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.MySQL.ConnMaxIdleTime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Fatal("Failed to ping MySQL", logger.Err(err))
	}

	logger.Info("MySQL connected successfully",
		logger.Int("max_idle_conns", cfg.Database.MySQL.MaxIdleConns),
		logger.Int("max_open_conns", cfg.Database.MySQL.MaxOpenConns),
	)
}

func initRedis(cfg *config.Config) {
	RDB = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password:        cfg.Database.Redis.Password,
		DB:              cfg.Database.Redis.DB,
		PoolSize:        cfg.Database.Redis.PoolSize,
		MinIdleConns:    cfg.Database.Redis.MinIdleConns,
		MaxRetries:      cfg.Database.Redis.MaxRetries,
		PoolTimeout:     cfg.Database.Redis.PoolTimeout,
		ConnMaxIdleTime: cfg.Database.Redis.ConnMaxIdleTime,
		DialTimeout:     cfg.Database.Redis.DialTimeout,
		ReadTimeout:     cfg.Database.Redis.ReadTimeout,
		WriteTimeout:    cfg.Database.Redis.WriteTimeout,
	})

	logger.Info("Connecting to Redis",
		logger.String("host", cfg.Database.Redis.Host),
		logger.Int("port", cfg.Database.Redis.Port),
		logger.Int("db", cfg.Database.Redis.DB),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if cfg.Database.Redis.Password != "" {
		if err := RDB.Do(ctx, "AUTH", cfg.Database.Redis.Password).Err(); err != nil {
			logger.Fatal("Redis authentication failed", logger.Err(err))
		}
		logger.Info("Redis authentication successful")
	}

	if _, err := RDB.Ping(ctx).Result(); err != nil {
		logger.Fatal("Failed to ping Redis", logger.Err(err))
	}

	logger.Info("Redis connected successfully",
		logger.Int("pool_size", cfg.Database.Redis.PoolSize),
		logger.Int("min_idle_conns", cfg.Database.Redis.MinIdleConns),
	)
}

func GetRedis() *redis.Client {
	return RDB
}

func CloseDB() error {
	logger.Info("Closing database connections...")

	// 关闭Redis连接
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			logger.Error("Error closing Redis", logger.Err(err))
		} else {
			logger.Info("Redis connection closed")
		}
	}

	// 关闭MySQL连接
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("Error closing MySQL", logger.Err(err))
			return err
		}
		logger.Info("MySQL connection closed")
	}

	return nil
}
