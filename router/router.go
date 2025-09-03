package router

import (
	"gin-demo/pkg/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置主路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 使用我们的访问日志中间件替代默认的日志中间件
	r.Use(middleware.AccessLogMiddleware())
	r.Use(gin.Recovery())

	// 添加全局错误处理中间件
	r.Use(middleware.ErrorHandlerMiddleware())

	// 添加限流中间件
	r.Use(middleware.RateLimitMiddleware())

	// GZIP压缩中间件
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 性能监控 (仅在开发环境)
	if gin.Mode() != gin.ReleaseMode {
		pprof.Register(r)
	}

	// 健康检查路由
	SetupHealthRoutes(r)

	// API路由组
	api := r.Group("/api")
	{
		SetupEmailRoutes(api)

		// 认证路由（无需JWT认证，但有更严格的限流）
		auth := api.Group("/auth")
		auth.Use(middleware.CustomRateLimitMiddleware(20, time.Minute)) // 20次/分钟
		SetupAuthRoutes(auth)

		// 用户路由（需要JWT认证）
		SetupUserRoutes(api)
	}

	// 设置404和405处理器
	r.NoRoute(middleware.NotFoundHandler())
	r.NoMethod(middleware.MethodNotAllowedHandler())

	return r
}
