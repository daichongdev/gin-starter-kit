package middleware

import (
	"context"
	"fmt"
	"gin-demo/database"
	"gin-demo/model"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	once           sync.Once
)

// getDefaultLimiter 获取默认限流器（延迟初始化）
func getDefaultLimiter() *RedisRateLimiter {
	once.Do(func() {
		defaultLimiter = NewRedisRateLimiter(30, time.Minute)
	})
	return defaultLimiter
}

// RateLimitMiddleware 基于Redis的限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !getDefaultLimiter().Allow(c.Request.Context(), ip) {
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
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) bool {
	// 获取Redis客户端
	rdb := database.GetRedis()
	if rdb == nil {
		// Redis未初始化时，允许请求通过（降级策略）
		return true
	}

	// Redis key
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// 当前时间戳（毫秒）
	now := time.Now().UnixMilli()

	// 窗口开始时间
	windowStart := now - rl.window.Milliseconds()

	// 使用Redis Pipeline提高性能
	pipe := rdb.Pipeline()

	// 1. 删除窗口外的记录
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(windowStart, 10))

	// 2. 统计当前窗口内的请求数
	countCmd := pipe.ZCard(ctx, redisKey)

	// 3. 添加当前请求
	pipe.ZAdd(ctx, redisKey, redis.Z{
		Score:  float64(now),
		Member: now, // 使用时间戳作为member确保唯一性
	})

	// 4. 设置过期时间
	pipe.Expire(ctx, redisKey, rl.window+time.Minute) // 多给1分钟缓冲

	// 执行Pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// Redis出错时，允许请求通过（降级策略）
		return true
	}

	// 检查请求数是否超限
	count := countCmd.Val()
	return count < int64(rl.limit)
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
