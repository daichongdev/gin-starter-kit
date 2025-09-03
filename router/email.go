package router

import (
	"gin-demo/controller"

	"github.com/gin-gonic/gin"
)

// SetupEmailRoutes 设置邮件相关路由
func SetupEmailRoutes(api *gin.RouterGroup) {
	emailController := controller.NewEmailController()

	emailGroup := api.Group("/email_queue")
	emailGroup.POST("/send", emailController.SendEmail)

	// 发送测试邮件（GET请求，方便浏览器测试）
	emailGroup.GET("/test", emailController.SendTestEmail)

	// 获取邮件队列状态
	emailGroup.GET("/status", emailController.GetEmailStatus)
}
