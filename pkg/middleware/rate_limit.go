package middleware

import (
	"gin-demo/config"
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

// RateLimitMiddleware 简单的限流中间件（从配置文件读取参数）
func RateLimitMiddleware() gin.HandlerFunc {
	// 从配置文件读取全局限流参数
	limiter := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    config.Cfg.Server.RateLimit.Global.Limit,
		window:   config.Cfg.Server.RateLimit.Global.Window,
	}
	
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

// AuthRateLimitMiddleware 认证接口限流中间件（从配置文件读取参数）
func AuthRateLimitMiddleware() gin.HandlerFunc {
	// 从配置文件读取认证接口限流参数
	limiter := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    config.Cfg.Server.RateLimit.Auth.Limit,
		window:   config.Cfg.Server.RateLimit.Auth.Window,
	}
	
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