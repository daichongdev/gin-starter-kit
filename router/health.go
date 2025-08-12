package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes 设置健康检查路由
func SetupHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})
	
	// 可以添加更多健康检查相关的路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}