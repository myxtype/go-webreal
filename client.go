package webreal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	newline = []byte{'\n'}
)

// 一个连接一个Client，负责处理连接的I/O
type Client struct {
	id        string
	writeChan chan []byte
	conn      *websocket.Conn
	handler   Handler
	hub       *SubscriptionHub
	req       *http.Request
	mu        sync.Mutex
	channels  map[string]struct{}
	conf      *Config
}

func newClient(conn *websocket.Conn, handler Handler, hub *SubscriptionHub, req *http.Request, c *Config) *Client {
	return &Client{
		id:        uuid.NewV4().String(),
		writeChan: make(chan []byte, c.WriteBufferSize),
		conn:      conn,
		handler:   handler,
		hub:       hub,
		channels:  map[string]struct{}{},
		req:       req,
		conf:      c,
	}
}

func (c *Client) run() {
	go c.reader()
	go c.writer()

	c.handler.OnConnect(c)
}

// 读取客户端发过来的内容
func (c *Client) reader() {
	defer func() {
		c.close()
		c.handler.OnClose(c)
	}()

	c.conn.SetReadLimit(c.conf.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.conf.PongWait))
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(c.conf.PongWait)) })

	for {
		_, buf, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var msg Message
		if err = json.Unmarshal(buf, &msg); err != nil {
			return
		}
		c.handler.OnMessage(c, &msg)
	}
}

// 向客户端写入内容
func (c *Client) writer() {
	tik := time.NewTicker(c.conf.PingInterval)
	defer func() {
		c.close()
		tik.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case buf, ok := <-c.writeChan:
			c.conn.SetWriteDeadline(time.Now().Add(c.conf.WriteWait))

			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, buf); err != nil {
				return
			}
		case <-tik.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.conf.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 关闭
func (c *Client) close() {
	c.UnsubscribeAll()
}

// 订阅
func (c *Client) Subscribe(channel string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.channels[channel]
	if found {
		return false
	}
	if c.hub.Subscribe(channel, c) {
		c.channels[channel] = struct{}{}
		return true
	}
	return false
}

// 退订
func (c *Client) Unsubscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hub.Unsubscribe(channel, c) {
		delete(c.channels, channel)
	}
}

// 退订所有
func (c *Client) UnsubscribeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for channel := range c.channels {
		c.hub.Unsubscribe(channel, c)
	}
	c.channels = map[string]struct{}{}
}

// 获取已订阅的主题列表
func (c *Client) Channels() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var channels []string
	for key := range c.channels {
		channels = append(channels, key)
	}
	return channels
}

// 获取已订阅长度
func (c *Client) ChannelsLen() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.channels)
}

// 获取客户端ID
func (c *Client) ID() string {
	return c.id
}

// 向客户端发送数据
func (c *Client) Write(d []byte) {
	c.writeChan <- d
}

// 获取连接参数
func (c *Client) Query() url.Values {
	return c.req.URL.Query()
}

// 获取请求对象
func (c *Client) Request() *http.Request {
	return c.req
}

// 关闭连接对象
func (c *Client) Close() error {
	return c.conn.Close()
}
