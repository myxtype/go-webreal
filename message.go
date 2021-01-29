package webreal

// 客户端发来的消息格式
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
