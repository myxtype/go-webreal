package webreal

import "encoding/json"

// 客户端发来的消息格式
type Message struct {
	Type string           `json:"type"`
	Data *json.RawMessage `json:"data"`
}

func (m *Message) ReadMessageData(x interface{}) error {
	if x == nil {
		return nil
	}
	return json.Unmarshal(*m.Data, x)
}
