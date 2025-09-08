package queue

import (
	"encoding/json"
	"fmt"
	"gin-demo/pkg/logger"
	"time"

	"go.uber.org/zap"
)

// EmailData 邮件数据结构
type EmailData struct {
	To      []string          `json:"to"`
	CC      []string          `json:"cc,omitempty"`
	BCC     []string          `json:"bcc,omitempty"`
	Subject string            `json:"subject"`
	Body    string            `json:"body"`
	IsHTML  bool              `json:"is_html"`
	Headers map[string]string `json:"headers,omitempty"`
}

// EmailHandler 邮件队列处理器
type EmailHandler struct {
	*BaseHandler
	logger *zap.Logger
}

// NewEmailHandler 创建邮件处理器
func NewEmailHandler(queueName string, numConsumers, prefetchLimit int) *EmailHandler {
	var emailLogger *zap.Logger
	if logger.Logger != nil {
		emailLogger = logger.Logger.Named(queueName)
	} else {
		emailLogger = zap.NewNop()
	}

	return &EmailHandler{
		BaseHandler: NewBaseHandler(queueName, numConsumers, prefetchLimit),
		logger:      emailLogger,
	}
}

// Handle 处理邮件消息
func (h *EmailHandler) Handle(message *Message) error {
	h.logger.Info("Processing email message",
		zap.String("message_id", message.ID),
		zap.String("type", message.Type))

	// 解析邮件数据
	var emailData EmailData
	if err := json.Unmarshal(message.Data, &emailData); err != nil {
		return fmt.Errorf("failed to unmarshal email data: %w", err)
	}

	// 验证邮件数据
	if err := h.validateEmailData(&emailData); err != nil {
		return fmt.Errorf("email data validation failed: %w", err)
	}

	// 发送邮件
	if err := h.sendEmail(&emailData); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	h.logger.Info("Email sent successfully",
		zap.String("message_id", message.ID),
		zap.Strings("to", emailData.To),
		zap.String("subject", emailData.Subject))

	return nil
}

// validateEmailData 验证邮件数据
func (h *EmailHandler) validateEmailData(data *EmailData) error {
	if len(data.To) == 0 {
		return fmt.Errorf("recipient list cannot be empty")
	}

	if data.Subject == "" {
		return fmt.Errorf("email subject cannot be empty")
	}

	if data.Body == "" {
		return fmt.Errorf("email body cannot be empty")
	}

	// 验证邮箱格式（简单验证）
	for _, email := range data.To {
		if !h.isValidEmail(email) {
			return fmt.Errorf("invalid email address: %s", email)
		}
	}

	return nil
}

// isValidEmail 简单的邮箱格式验证
func (h *EmailHandler) isValidEmail(email string) bool {
	// 这里可以使用更复杂的邮箱验证逻辑
	return len(email) > 0 && email != ""
}

// sendEmail 发送邮件（模拟实现）
func (h *EmailHandler) sendEmail(data *EmailData) error {
	// 模拟邮件发送延迟
	time.Sleep(100 * time.Millisecond)

	// 这里应该集成实际的邮件发送服务，如：
	// - SMTP 服务器
	// - 第三方邮件服务（如 SendGrid, AWS SES 等）
	// - 企业邮件服务

	h.logger.Debug("Sending email",
		zap.Strings("to", data.To),
		zap.String("subject", data.Subject),
		zap.Bool("is_html", data.IsHTML))

	// 模拟发送成功
	return nil
}
