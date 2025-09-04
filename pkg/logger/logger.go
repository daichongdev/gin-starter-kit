package logger

import (
	"gin-demo/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger       *zap.Logger
	SugarLogger  *zap.SugaredLogger
	AccessLogger *zap.Logger
	DBLogger     *zap.Logger
)

// InitLogger 初始化日志系统
func InitLogger(cfg *config.Config) error {
	// 确保日志目录存在
	if err := os.MkdirAll("./logs", 0755); err != nil {
		return err
	}

	var cores []zapcore.Core

	// 控制台输出
	if cfg.Log.Console.Enabled {
		cores = append(cores, createConsoleCore(cfg))
	}

	// 文件输出
	if cfg.Log.File.Enabled {
		cores = append(cores, createFileCore(*cfg.Log.File))
	}

	// 错误文件输出
	if cfg.Log.ErrorFile.Enabled {
		cores = append(cores, createErrorFileCore(*cfg.Log.ErrorFile))
	}

	// 组合所有核心
	core := zapcore.NewTee(cores...)

	// 配置选项
	var options []zap.Option
	if cfg.Log.EnableCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.Log.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 创建日志器
	Logger = zap.New(core, options...)
	SugarLogger = Logger.Sugar()

	// 初始化访问日志器
	if err := initAccessLogger(cfg); err != nil {
		return err
	}

	// 初始化数据库日志器
	if err := initDBLogger(cfg); err != nil {
		return err
	}

	return nil
}

// createConsoleCore 创建控制台输出核心
func createConsoleCore(cfg *config.Config) zapcore.Core {
	level := getLogLevel(cfg.Log.Console.Level)

	var encoder zapcore.Encoder
	if cfg.Log.Console.Format == "json" {
		encoder = zapcore.NewJSONEncoder(getEncoderConfig(false))
	} else {
		encoder = zapcore.NewConsoleEncoder(getEncoderConfig(cfg.Log.Console.EnableColor))
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// createFileCore 创建文件输出核心
func createFileCore(fileCfg config.FileLogConfig) zapcore.Core {
	level := getLogLevel(fileCfg.Level)

	writer := &lumberjack.Logger{
		Filename:   fileCfg.Filename,
		MaxSize:    fileCfg.MaxSize,
		MaxBackups: fileCfg.MaxBackups,
		MaxAge:     fileCfg.MaxAge,
		Compress:   fileCfg.Compress,
	}

	var encoder zapcore.Encoder
	if fileCfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(getEncoderConfig(false))
	} else {
		encoder = zapcore.NewConsoleEncoder(getEncoderConfig(false))
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
}

// createErrorFileCore 创建错误文件输出核心
func createErrorFileCore(fileCfg config.FileLogConfig) zapcore.Core {
	level := zapcore.ErrorLevel

	writer := &lumberjack.Logger{
		Filename:   fileCfg.Filename,
		MaxSize:    fileCfg.MaxSize,
		MaxBackups: fileCfg.MaxBackups,
		MaxAge:     fileCfg.MaxAge,
		Compress:   fileCfg.Compress,
	}

	encoder := zapcore.NewJSONEncoder(getEncoderConfig(false))
	return zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
}

// initAccessLogger 初始化访问日志器
func initAccessLogger(cfg *config.Config) error {
	if !cfg.Log.Access.Enabled {
		return nil
	}

	writer := &lumberjack.Logger{
		Filename:   cfg.Log.Access.Filename,
		MaxSize:    cfg.Log.Access.MaxSize,
		MaxBackups: cfg.Log.Access.MaxBackups,
		MaxAge:     cfg.Log.Access.MaxAge,
		Compress:   cfg.Log.Access.Compress,
	}

	var encoder zapcore.Encoder
	if cfg.Log.Access.Format == "json" {
		encoder = zapcore.NewJSONEncoder(getAccessEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(getAccessEncoderConfig())
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.InfoLevel)
	AccessLogger = zap.New(core)

	return nil
}

// initDBLogger 初始化数据库日志器 - 支持文件轮转
func initDBLogger(cfg *config.Config) error {
	if !cfg.Log.Database.Enabled {
		return nil
	}

	writer := &lumberjack.Logger{
		Filename:   cfg.Log.Database.Filename,
		MaxSize:    cfg.Log.Database.MaxSize,
		MaxBackups: cfg.Log.Database.MaxBackups,
		MaxAge:     cfg.Log.Database.MaxAge,
		Compress:   cfg.Log.Database.Compress,
	}

	encoder := zapcore.NewJSONEncoder(getEncoderConfig(false))
	level := getLogLevel(cfg.Log.Database.Level)
	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
	DBLogger = zap.New(core, zap.AddCaller())

	return nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoderConfig 获取编码器配置
func getEncoderConfig(enableColor bool) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if enableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return encoderConfig
}

// getAccessEncoderConfig 获取访问日志编码器配置
func getAccessEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

// Sync 同步所有日志器
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
	if AccessLogger != nil {
		AccessLogger.Sync()
	}
	if DBLogger != nil {
		DBLogger.Sync()
	}
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// 便捷方法
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Debugf 结构化日志方法
func Debugf(template string, args ...interface{}) {
	if SugarLogger != nil {
		SugarLogger.Debugf(template, args...)
	}
}

func Infof(template string, args ...interface{}) {
	if SugarLogger != nil {
		SugarLogger.Infof(template, args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if SugarLogger != nil {
		SugarLogger.Warnf(template, args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if SugarLogger != nil {
		SugarLogger.Errorf(template, args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if SugarLogger != nil {
		SugarLogger.Fatalf(template, args...)
	}
}
