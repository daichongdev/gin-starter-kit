package middleware

import (
	"context"
	"fmt"
	"gin-demo/config"
	"gin-demo/database"
	"gin-demo/model"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RedisRateLimiter Redis限流器
type RedisRateLimiter struct {
	limit  int           // 限制次数
	window time.Duration // 时间窗口
}

// NewRedisRateLimiter 创建Redis限流器
func NewRedisRateLimiter(limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		limit:  limit,
		window: window,
	}
}

// 默认限流器
var (
	defaultLimiter *RedisRateLimiter
	authLimiter    *RedisRateLimiter
	limitOnce      sync.Once
)

// getDefaultLimiter 获取默认限流器（从配置文件读取参数）
func getDefaultLimiter() *RedisRateLimiter {
	limitOnce.Do(func() {
		defaultLimiter = NewRedisRateLimiter(
			config.Cfg.Server.RateLimit.Global.Limit,
			config.Cfg.Server.RateLimit.Global.Window,
		)
		authLimiter = NewRedisRateLimiter(
			config.Cfg.Server.RateLimit.Auth.Limit,
			config.Cfg.Server.RateLimit.Auth.Window,
		)
	})
	return defaultLimiter
}

// getAuthLimiter 获取认证限流器
func getAuthLimiter() *RedisRateLimiter {
	limitOnce.Do(func() {
		defaultLimiter = NewRedisRateLimiter(
			config.Cfg.Server.RateLimit.Global.Limit,
			config.Cfg.Server.RateLimit.Global.Window,
		)
		authLimiter = NewRedisRateLimiter(
			config.Cfg.Server.RateLimit.Auth.Limit,
			config.Cfg.Server.RateLimit.Auth.Window,
		)
	})
	return authLimiter
}

// RateLimitMiddleware 简单的限流中间件（从配置文件读取参数）
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getDefaultLimiter()

		if !limiter.Allow(c.Request.Context(), ip) {
			c.JSON(http.StatusTooManyRequests, model.ErrorResponse("请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware 认证接口限流中间件（从配置文件读取参数）
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getAuthLimiter()

		if !limiter.Allow(c.Request.Context(), ip) {
			c.JSON(http.StatusTooManyRequests, model.ErrorResponse("请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CustomRateLimitMiddleware 自定义限流参数的中间件
func CustomRateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRedisRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(c.Request.Context(), ip) {
			c.JSON(http.StatusTooManyRequests, model.ErrorResponse("请求过于频繁，请稍后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow 检查是否允许请求（使用Redis滑动窗口算法）
// 优化限流算法，使用滑动窗口
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) bool {
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	// 使用Lua脚本保证原子性
	luaScript := `
		local key = KEYS[1]
		local window_start = ARGV[1]
		local now = ARGV[2]
		local limit = tonumber(ARGV[3])
		
		-- 清理过期数据
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
		
		-- 获取当前计数
		local current = redis.call('ZCARD', key)
		
		if current < limit then
			-- 添加当前请求
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, 3600)
			return 1
		else
			return 0
		end
	`

	result, err := database.GetRedis().Eval(ctx, luaScript, []string{key},
		windowStart, now, rl.limit).Result()
	if err != nil {
		return false
	}

	return result.(int64) == 1
}

// GetRemainingRequests 获取剩余请求次数
func (rl *RedisRateLimiter) GetRemainingRequests(ctx context.Context, key string) int {
	// 获取Redis客户端
	rdb := database.GetRedis()
	if rdb == nil {
		return rl.limit // Redis未初始化时返回最大值
	}

	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// 当前时间戳（毫秒）
	now := time.Now().UnixMilli()
	windowStart := now - rl.window.Milliseconds()

	// 清理过期记录并统计
	pipe := rdb.Pipeline()
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(windowStart, 10))
	countCmd := pipe.ZCard(ctx, redisKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return rl.limit // 出错时返回最大值
	}

	used := int(countCmd.Val())
	remaining := rl.limit - used
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

// RateLimitInfoMiddleware 添加限流信息到响应头
func RateLimitInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getDefaultLimiter()

		// 获取剩余请求次数
		remaining := limiter.GetRemainingRequests(c.Request.Context(), ip)

		// 添加到响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(limiter.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Window", limiter.window.String())

		c.Next()
	}
}
