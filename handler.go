package webreal

// 业务处理程序
type Handler interface {
	// 连接时触发
	OnConnect(client *Client)
	// 有新消息时
	OnMessage(client *Client, msg *Message)
	// 连接关闭时
	OnClose(client *Client)
}
