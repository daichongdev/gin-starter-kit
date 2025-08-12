package router

import (
	"gin-demo/controller"
	"gin-demo/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(api *gin.RouterGroup) {
	authController := controller.NewAuthController()
	authGroup := api.Group("/auth")
	{
		// 公开路由（无需认证）
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)

		// 需要认证的路由
		protected := authGroup.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.GET("/profile", authController.GetProfile)
			protected.POST("/refresh", authController.RefreshToken)
		}
	}
}
