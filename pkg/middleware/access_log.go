package middleware

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gin-demo/pkg/logger"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogEntry 日志条目结构
type LogEntry struct {
	TraceID      string
	Method       string
	Path         string
	ClientIP     string
	StatusCode   int
	BodySize     int
	Latency      time.Duration
	UserAgent    string
	RequestBody  string
	ResponseBody string
	Errors       []string
	ErrorDetails []interface{}
	SQLLogs      []logger.SQLLog
	Timestamp    time.Time
}

// AsyncLogger 异步日志处理器
type AsyncLogger struct {
	logChan    chan *LogEntry
	workerPool chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

var (
	asyncLogger *AsyncLogger
	once        sync.Once
)

// initAsyncLogger 初始化异步日志处理器
func initAsyncLogger() {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		asyncLogger = &AsyncLogger{
			logChan:    make(chan *LogEntry, 1000),            // 缓冲1000条日志
			workerPool: make(chan struct{}, runtime.NumCPU()), // 工作池大小为CPU核心数
			ctx:        ctx,
			cancel:     cancel,
		}

		// 启动日志处理协程
		for i := 0; i < runtime.NumCPU(); i++ {
			asyncLogger.wg.Add(1)
			go asyncLogger.worker()
		}
	})
}

// worker 日志处理工作协程
func (al *AsyncLogger) worker() {
	defer al.wg.Done()

	for {
		select {
		case entry := <-al.logChan:
			if entry != nil {
				al.processLogEntry(entry)
			}
		case <-al.ctx.Done():
			return
		}
	}
}

// processLogEntry 处理单条日志
func (al *AsyncLogger) processLogEntry(entry *LogEntry) {
	// 构建日志字段（只构建一次）
	fields := []zap.Field{
		zap.String("trace_id", entry.TraceID),
		zap.String("method", entry.Method),
		zap.String("path", entry.Path),
		zap.String("client_ip", entry.ClientIP),
		zap.Int("status_code", entry.StatusCode),
		zap.Int("body_size", entry.BodySize),
		zap.Duration("latency", entry.Latency),
		zap.String("user_agent", entry.UserAgent),
		zap.Time("timestamp", entry.Timestamp),
	}

	// 添加错误信息
	if len(entry.Errors) > 0 {
		fields = append(fields,
			zap.Strings("errors", entry.Errors),
			zap.Any("error_details", entry.ErrorDetails),
		)
	}

	// 添加请求体（如果需要）
	if entry.RequestBody != "" {
		fields = append(fields, zap.String("request_body", entry.RequestBody))
	}

	// 添加响应体（如果需要）
	if entry.ResponseBody != "" {
		fields = append(fields, zap.String("response_body", entry.ResponseBody))
	}

	// 添加SQL日志
	if len(entry.SQLLogs) > 0 {
		fields = append(fields,
			zap.Int("sql_count", len(entry.SQLLogs)),
			zap.Any("sql_details", entry.SQLLogs),
		)
	}

	// 记录访问日志
	if logger.AccessLogger != nil {
		logger.AccessLogger.Info("access", fields...)
	}

	// 根据状态码记录不同级别的日志
	switch {
	case entry.StatusCode >= 500:
		logger.Error("Server Error", fields...)
	case entry.StatusCode == 429:
		logger.Warn("Rate Limited", fields...)
	case entry.StatusCode == 401:
		logger.Info("Authentication Required", fields...)
	case entry.StatusCode == 403:
		logger.Warn("Access Forbidden", fields...)
	case entry.StatusCode == 404:
		logger.Info("Resource Not Found", fields...)
	case entry.StatusCode >= 400:
		logger.Warn("Client Error", fields...)
	}
}

// sendLogEntry 发送日志条目到异步处理队列
func (al *AsyncLogger) sendLogEntry(entry *LogEntry) {
	select {
	case al.logChan <- entry:
		// 成功发送
	default:
		// 队列满了，丢弃日志或同步处理
		logger.Warn("Log queue full, processing synchronously")
		al.processLogEntry(entry)
	}
}

// Shutdown 优雅关闭异步日志处理器
func (al *AsyncLogger) Shutdown() {
	al.cancel()
	close(al.logChan)
	al.wg.Wait()
}

// AccessLogMiddleware 优化后的访问日志中间件
func AccessLogMiddleware() gin.HandlerFunc {
	// 初始化异步日志处理器
	initAsyncLogger()

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 生成高性能链路追踪ID
		traceID := generateFastTraceID()

		// 设置链路追踪上下文（使用context而非全局变量）
		ctx := context.WithValue(c.Request.Context(), "trace_id", traceID)
		c.Request = c.Request.WithContext(ctx)

		// 条件性读取请求体
		var requestBody string
		if shouldLogRequestBody(c) {
			requestBody = readRequestBodySafely(c)
		}

		// 处理请求
		c.Next()

		// 异步记录日志
		go func() {
			entry := &LogEntry{
				TraceID:     traceID,
				Method:      c.Request.Method,
				Path:        buildFullPath(path, raw),
				ClientIP:    c.ClientIP(),
				StatusCode:  c.Writer.Status(),
				BodySize:    c.Writer.Size(),
				Latency:     time.Since(start),
				UserAgent:   c.Request.UserAgent(),
				RequestBody: requestBody,
				Timestamp:   start,
			}

			// 收集错误信息
			collectErrors(c, entry)

			// 获取SQL日志（从context中获取）
			if sqlLogs := getSQLLogsFromContext(ctx); len(sqlLogs) > 0 {
				entry.SQLLogs = sqlLogs
			}

			// 发送到异步处理队列
			asyncLogger.sendLogEntry(entry)
		}()
	}
}

// generateFastTraceID 高性能链路追踪ID生成
func generateFastTraceID() string {
	// 使用更高效的方式生成ID
	now := time.Now()
	timestamp := now.UnixNano()

	// 生成8字节随机数
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	return fmt.Sprintf("%d%s", timestamp, hex.EncodeToString(randomBytes))
}

// shouldLogRequestBody 判断是否需要记录请求体
func shouldLogRequestBody(c *gin.Context) bool {
	// 只对特定方法和内容类型记录请求体
	if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
		return false
	}

	// 检查Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/json") &&
		!strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return false
	}

	// 检查Content-Length，避免记录过大的请求体
	if c.Request.ContentLength > 10240 { // 10KB限制
		return false
	}

	return true
}

// readRequestBodySafely 安全读取请求体
func readRequestBodySafely(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	// 限制读取大小
	limitedReader := io.LimitReader(c.Request.Body, 10240) // 10KB限制
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return ""
	}

	// 重新设置请求体
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

// buildFullPath 构建完整路径
func buildFullPath(path, raw string) string {
	if raw != "" {
		return path + "?" + raw
	}
	return path
}

// collectErrors 收集错误信息
func collectErrors(c *gin.Context, entry *LogEntry) {
	var errorMessages []string
	var errorDetails []interface{}

	// 获取Gin框架中的错误
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			errorMessages = append(errorMessages, err.Error())
			errorDetails = append(errorDetails, map[string]interface{}{
				"error": err.Error(),
				"type":  err.Type,
				"meta":  err.Meta,
			})
		}
	}

	// 根据状态码添加错误信息
	if entry.StatusCode >= 400 && len(errorMessages) == 0 {
		errorMessages = append(errorMessages, getStatusText(entry.StatusCode))
	}

	entry.Errors = errorMessages
	entry.ErrorDetails = errorDetails
}

// getSQLLogsFromContext 从context获取SQL日志
func getSQLLogsFromContext(ctx context.Context) []logger.SQLLog {
	// 这里需要修改logger包来支持context-based的SQL日志收集
	// 暂时返回空，需要配合logger包的修改
	return nil
}

// getStatusText 根据状态码获取描述文本
func getStatusText(statusCode int) string {
	switch statusCode {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 409:
		return "Conflict"
	case 422:
		return "Unprocessable Entity"
	case 429:
		return "Too Many Requests"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	case 504:
		return "Gateway Timeout"
	default:
		return "Unknown Error"
	}
}
