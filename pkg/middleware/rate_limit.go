package middleware

import (
	"gin-demo/model"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 简单的内存限流器
type rateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int           // 限制次数
	window   time.Duration // 时间窗口
}

var limiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100,                // 默认100次/分钟
	window:   time.Minute,        // 1分钟窗口
}

// RateLimitMiddleware 简单的限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, model.ErrorResponse("请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// CustomRateLimitMiddleware 自定义限流参数的中间件
func CustomRateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	customLimiter := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !customLimiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, model.ErrorResponse("请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// allow 检查是否允许请求
func (rl *rateLimiter) allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	
	// 获取该IP的请求记录
	requests := rl.requests[key]
	
	// 清理过期的请求记录
	var validRequests []time.Time
	for _, req := range requests {
		if now.Sub(req) < rl.window {
			validRequests = append(validRequests, req)
		}
	}
	
	// 检查是否超过限制
	if len(validRequests) >= rl.limit {
		rl.requests[key] = validRequests
		return false
	}
	
	// 添加当前请求
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests
	
	return true
}