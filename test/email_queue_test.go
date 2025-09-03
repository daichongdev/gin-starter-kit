package test

import (
	"gin-demo/pkg/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// 测试邮件处理器
func TestEmailHandler(t *testing.T) {
	// 创建邮件处理器
	handler := queue.NewEmailHandler("email", 2, 10)

	// 测试处理器配置
	assert.Equal(t, "email", handler.GetQueueName())
	assert.Equal(t, 2, handler.GetNumConsumers())
	assert.Equal(t, 10, handler.GetPrefetchLimit())

	// 测试有效邮件数据处理
	validEmailData := &queue.EmailData{
		To:      []string{"valid@example.com"},
		Subject: "有效邮件",
		Body:    "这是有效的邮件内容",
		IsHTML:  false,
	}

	message, err := queue.NewMessage("email.send", validEmailData)
	require.NoError(t, err)

	// 处理消息（这里会调用模拟的发送逻辑）
	err = handler.Handle(message)
	assert.NoError(t, err)
}
