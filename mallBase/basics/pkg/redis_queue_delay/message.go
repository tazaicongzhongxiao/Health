package redis_queue_delay

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/queue"
	"MyTestMall/mallBase/basics/pkg/unique"
	"encoding/json"
	"time"
)

type Message struct {
	Id          int64     `json:"id"`
	CreateTime  time.Time `json:"createTime"`
	ConsumeTime time.Time `json:"consumeTime"`
	Body        []byte    `json:"body"`
}

// NewMessage 创建消息实体
func NewMessage(body queue.Message, delay int32) *Message {
	return &Message{
		Id:          unique.ID(),
		CreateTime:  time.Now(),
		ConsumeTime: time.Now().Add(time.Duration(delay) * time.Second),
		Body:        app.Struct2Json(body),
	}
}

func (m *Message) GetScore() float64 {
	return float64(m.ConsumeTime.Unix())
}

func (m *Message) GetId() int64 {
	return m.Id
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
