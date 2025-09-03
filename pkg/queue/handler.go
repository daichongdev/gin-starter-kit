package queue

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// MessageHandler 消息处理器接口
type MessageHandler interface {
	Handle(message *Message) error
	GetQueueName() string
	GetNumConsumers() int
	GetPrefetchLimit() int
}

// Message 消息结构体
type Message struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Data       []byte                 `json:"data"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	RetryCount int                    `json:"retry_count"`
}

// NewMessage 创建新消息
func NewMessage(msgType string, data interface{}) (*Message, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:         uuid.New().String(),
		Type:       msgType,
		Data:       dataBytes,
		Metadata:   make(map[string]interface{}),
		Timestamp:  time.Now(),
		RetryCount: 0,
	}, nil
}

// ToJSON 将消息转换为JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从JSON创建消息
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// BaseHandler 基础处理器，提供默认实现
type BaseHandler struct {
	queueName     string
	numConsumers  int
	prefetchLimit int
}

func NewBaseHandler(queueName string, numConsumers, prefetchLimit int) *BaseHandler {
	return &BaseHandler{
		queueName:     queueName,
		numConsumers:  numConsumers,
		prefetchLimit: prefetchLimit,
	}
}

func (h *BaseHandler) GetQueueName() string {
	return h.queueName
}

func (h *BaseHandler) GetNumConsumers() int {
	return h.numConsumers
}

func (h *BaseHandler) GetPrefetchLimit() int {
	return h.prefetchLimit
}
