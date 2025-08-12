package logger

import (
	"context"
	"gin-demo/config"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger GORM日志适配器
type GormLogger struct {
	ZapLogger                 *zap.Logger
	SlowThreshold             time.Duration
	LogLevel                  logger.LogLevel
	LogAllSQL                 bool
	IgnoreRecordNotFoundError bool
}

// NewGormLogger 创建GORM日志器
func NewGormLogger(cfg *config.DatabaseLogConfig) *GormLogger {
	var logLevel logger.LogLevel
	switch cfg.Level {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Info
	}

	return &GormLogger{
		ZapLogger:                 DBLogger,
		SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
		LogLevel:                  logLevel,
		LogAllSQL:                 cfg.LogAllSQL,
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFound,
	}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info && l.ZapLogger != nil {
		fields := []zap.Field{zap.String("message", msg)}
		if traceID := GetCurrentTraceID(); traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}
		l.ZapLogger.Info("GORM Info", fields...)
	}
}

// Warn 警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn && l.ZapLogger != nil {
		fields := []zap.Field{zap.String("message", msg)}
		if traceID := GetCurrentTraceID(); traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}
		l.ZapLogger.Warn("GORM Warn", fields...)
	}
}

// Error 错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error && l.ZapLogger != nil {
		fields := []zap.Field{zap.String("message", msg)}
		if traceID := GetCurrentTraceID(); traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}
		l.ZapLogger.Error("GORM Error", fields...)
	}
}

// Trace SQL跟踪日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 获取当前链路追踪ID
	traceID := GetCurrentTraceID()

	// 构建基础字段
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	// 添加链路追踪ID
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// 记录到数据库日志文件
	if l.ZapLogger != nil {
		switch {
		case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || err.Error() != "record not found"):
			l.ZapLogger.Error("SQL Error", append(fields, zap.Error(err))...)
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			l.ZapLogger.Warn("Slow SQL", append(fields, zap.Duration("threshold", l.SlowThreshold))...)
		case l.LogAllSQL && l.LogLevel >= logger.Info:
			l.ZapLogger.Info("SQL", fields...)
		}
	}

	// 添加到链路追踪日志
	if traceID != "" {
		AddSQLLog(sql, elapsed.String(), rows, err)
	}
}
