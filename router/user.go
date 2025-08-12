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
		userGroup.POST("/", userController.CreateUser)                    // 创建用户
		userGroup.GET("/", userController.GetAllUsers)                   // 获取所有用户（不分页）
		userGroup.GET("/paginated", userController.GetUsersWithPagination) // 分页获取用户列表
		userGroup.GET("/:id", userController.GetUser)                    // 获取单个用户
		userGroup.PUT("/:id", userController.UpdateUser)                 // 更新用户
		userGroup.DELETE("/:id", userController.DeleteUser)              // 删除用户
	}
}