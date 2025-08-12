package router

import (
	"gin-demo/controller"
	"gin-demo/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(api *gin.RouterGroup) {
	userController := controller.NewUserController()
	userGroup := api.Group("/users")
	
	// 应用JWT认证中间件到所有用户路由
	userGroup.Use(middleware.JWTAuthMiddleware())
	{
		userGroup.POST("/", userController.CreateUser)
		userGroup.GET("/", userController.GetAllUsers)
		userGroup.GET("/:id", userController.GetUser)
	}
}