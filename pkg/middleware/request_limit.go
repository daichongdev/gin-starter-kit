package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestSizeLimitMiddleware 限制请求体大小
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	})
}
