package controller

import (
	"gin-demo/model"
	"gin-demo/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmailController struct {
	emailService *service.EmailService
}

func NewEmailController() *EmailController {
	return &EmailController{
		emailService: service.NewEmailService(),
	}
}

// SendEmailRequest 发送邮件请求结构
type SendEmailRequest struct {
	To      []string `json:"to" binding:"required"`
	CC      []string `json:"cc,omitempty"`
	BCC     []string `json:"bcc,omitempty"`
	Subject string   `json:"subject" binding:"required"`
	Body    string   `json:"body" binding:"required"`
	IsHTML  bool     `json:"is_html"`
}

// SendEmail 发送邮件接口
func (c *EmailController) SendEmail(ctx *gin.Context) {
	var req SendEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(err.Error()))
		return
	}

	// 发送邮件到队列
	err := c.emailService.SendEmail(req.To, req.Subject, req.Body, req.IsHTML)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("邮件已加入发送队列", nil))
}

// SendTestEmail 发送测试邮件接口
func (c *EmailController) SendTestEmail(ctx *gin.Context) {
	// 获取查询参数
	to := ctx.DefaultQuery("to", "test@example.com")
	subject := ctx.DefaultQuery("subject", "测试邮件")
	body := ctx.DefaultQuery("body", "这是一封测试邮件")
	isHTMLStr := ctx.DefaultQuery("is_html", "false")

	isHTML, _ := strconv.ParseBool(isHTMLStr)

	// 发送测试邮件
	err := c.emailService.SendEmail([]string{to}, subject, body, isHTML)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("测试邮件已加入发送队列", gin.H{
		"to":      to,
		"subject": subject,
		"body":    body,
		"is_html": isHTML,
	}))
}

// GetEmailStatus 获取邮件队列状态（用于调试）
func (c *EmailController) GetEmailStatus(ctx *gin.Context) {
	// 这里可以添加队列状态查询逻辑
	ctx.JSON(http.StatusOK, model.SuccessResponse("邮件队列状态", gin.H{
		"status": "running",
		"queue":  "email_queue",
	}))
}
