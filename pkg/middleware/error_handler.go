package middleware

import (
	"fmt"
	"gin-demo/model/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ErrorHandlerMiddleware 全局错误处理中间件（包含panic恢复）
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer捕获panic
		defer func() {
			if r := recover(); r != nil {
				handlePanic(c, r)
			}
		}()

		c.Next()

		// 检查是否有错误发生
		if len(c.Errors) > 0 {
			handleError(c, c.Errors.Last())
		}
	}
}

// handleError 处理普通错误
func handleError(c *gin.Context, err *gin.Error) {
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

// handlePanic 处理panic恢复
func handlePanic(c *gin.Context, r interface{}) {
	// 分析panic类型，提供更友好的错误信息
	errorMessage := analyzePanicType(r)
	// 确保响应没有被写入
	if !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, tool.ErrorResponse(errorMessage))
	}

	// 中止请求处理
	c.Abort()
}

// analyzePanicType 分析panic类型，返回用户友好的错误信息
func analyzePanicType(r interface{}) string {
	panicStr := fmt.Sprintf("%v", r)

	// 常见panic类型分析
	switch {
	case strings.Contains(panicStr, "nil pointer dereference"):
		return "服务器内部错误：空指针访问"
	case strings.Contains(panicStr, "index out of range"):
		return "服务器内部错误：数组越界"
	case strings.Contains(panicStr, "slice bounds out of range"):
		return "服务器内部错误：切片越界"
	case strings.Contains(panicStr, "invalid memory address"):
		return "服务器内部错误：内存访问异常"
	case strings.Contains(panicStr, "concurrent map"):
		return "服务器内部错误：并发访问冲突"
	case strings.Contains(panicStr, "close of closed channel"):
		return "服务器内部错误：通道操作异常"
	default:
		return "服务器内部错误，请稍后重试"
	}
}

// NotFoundHandler 404处理器
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, tool.ErrorResponse("请求的接口不存在"))
	}
}

// MethodNotAllowedHandler 405方法不允许处理器
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, tool.ErrorResponse("HTTP方法不被允许"))
	}
}
