package middleware

import (
	"bytes"
	"gin-demo/pkg/logger"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AccessLogMiddleware 访问日志中间件，支持链路追踪
func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 生成链路追踪ID
		traceID := generateTraceID()
		
		// 设置当前goroutine的链路追踪上下文
		logger.SetTraceContext(traceID)
		defer logger.ClearTraceContext()

		// 读取请求体（如果需要记录）
		var requestBody string
		if c.Request.Body != nil && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH") {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// 重新设置请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取请求方法
		method := c.Request.Method

		// 获取状态码
		statusCode := c.Writer.Status()

		// 获取响应大小
		bodySize := c.Writer.Size()

		// 获取User-Agent
		userAgent := c.Request.UserAgent()

		// 构建完整路径
		if raw != "" {
			path = path + "?" + raw
		}

		// 构建日志字段
		fields := []zap.Field{
			zap.String("trace_id", traceID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Int("status_code", statusCode),
			zap.Int("body_size", bodySize),
			zap.Duration("latency", latency),
			zap.String("user_agent", userAgent),
		}

		// 添加请求体（如果有且不为空）
		if requestBody != "" && len(requestBody) < 1000 { // 限制长度避免日志过大
			fields = append(fields, zap.String("request_body", requestBody))
		}

		// 获取SQL执行日志
		sqlLogs := logger.GetSQLLogs()
		if len(sqlLogs) > 0 {
			fields = append(fields,
				zap.Int("sql_count", len(sqlLogs)),
				zap.Any("sql_details", sqlLogs),
			)
		}

		// 记录访问日志
		if logger.AccessLogger != nil {
			logger.AccessLogger.Info("access", fields...)
		}

		// 根据状态码类型记录不同级别的日志
		switch {
		case statusCode >= 500:
			// 5xx 服务器错误 - 记录为ERROR级别
			logger.Error("Server Error", fields...)
		case statusCode == 401:
			// 401 认证失败 - 记录为INFO级别（正常业务逻辑）
			logger.Info("Authentication Required", fields...)
		case statusCode == 403:
			// 403 权限不足 - 记录为WARN级别
			logger.Warn("Access Forbidden", fields...)
		case statusCode == 404:
			// 404 资源未找到 - 记录为INFO级别
			logger.Info("Resource Not Found", fields...)
		case statusCode >= 400:
			// 其他4xx客户端错误 - 记录为WARN级别
			logger.Warn("Client Error", fields...)
		}
	}
}

// generateTraceID 生成链路追踪ID
func generateTraceID() string {
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "") + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// GetTraceID 获取当前请求的链路追踪ID（保持向后兼容）
func GetTraceID(c *gin.Context) string {
	return logger.GetCurrentTraceID()
}
