package middleware

import (
	"gin-demo/model/tool"
	"gin-demo/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误发生
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 记录错误日志
			logger.Error("Request error occurred",
				logger.String("path", c.Request.URL.Path),
				logger.String("method", c.Request.Method),
				logger.String("error", err.Error()),
				logger.String("client_ip", c.ClientIP()))

			// 根据错误类型返回不同的HTTP状态码
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, tool.ErrorResponse("请求格式无效: "+err.Error()))
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, tool.ErrorResponse("服务器内部错误"))
			default:
				c.JSON(http.StatusInternalServerError, tool.ErrorResponse("服务器内部错误"))
			}
		}
	}
}

// NotFoundHandler 404处理器
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Route not found",
			logger.String("path", c.Request.URL.Path),
			logger.String("method", c.Request.Method),
			logger.String("client_ip", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()))

		c.JSON(http.StatusNotFound, tool.ErrorResponse("请求的接口不存在"))
	}
}

// MethodNotAllowedHandler 405方法不允许处理器
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Warn("Method not allowed",
			logger.String("path", c.Request.URL.Path),
			logger.String("method", c.Request.Method),
			logger.String("client_ip", c.ClientIP()))

		c.JSON(http.StatusMethodNotAllowed, tool.ErrorResponse("HTTP方法不被允许"))
	}
}
