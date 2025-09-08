package service

import (
	"fmt"
	"gin-demo/pkg/queue"
)

// EmailService 邮件服务
type EmailService struct {
	// 移除 queueManager 字段，直接使用全局实例
}

// NewEmailService 创建邮件服务（无需参数）
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendEmail 发送邮件（异步）
func (s *EmailService) SendEmail(to []string, subject, body string, isHTML bool) error {
	emailData := &queue.EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	// 直接使用全局队列管理器
	manager := queue.GetManager()
	if manager == nil {
		return fmt.Errorf("queue manager not initialized")
	}

	return manager.PublishData(queue.Email, "email.send", emailData)
}
