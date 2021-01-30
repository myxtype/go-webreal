package webreal

import "encoding/json"

// 消息
// 客户端需要按照此结构体发送消息到服务端
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func (m *Message) ReadMessageData(x interface{}) error {
	if x == nil {
		return nil
	}
	return json.Unmarshal(m.Data, x)
}
