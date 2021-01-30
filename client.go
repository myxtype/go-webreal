package webreal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

const (
	writeWait      = 10 * time.Second    // 写等待
	pongWait       = 60 * time.Second    // 心跳等待
	pingPeriod     = (pongWait * 9) / 10 // 心跳频率
	maxMessageSize = 524288              // 512 kb
)

var (
	clientId int64
)

// 一个连接一个Client，负责处理连接的I/O
type Client struct {
	id        int64
	writeChan chan []byte
	conn      *websocket.Conn
	handler   Handler
	hub       *SubscriptionHub
	req       *http.Request
	mu        sync.Mutex
	channels  map[string]struct{}
}

func NewClient(conn *websocket.Conn, handler Handler, hub *SubscriptionHub, req *http.Request) *Client {
	return &Client{
		id:        atomic.AddInt64(&clientId, 1),
		writeChan: make(chan []byte, 256),
		conn:      conn,
		handler:   handler,
		hub:       hub,
		channels:  map[string]struct{}{},
		req:       req,
	}
}

func (c *Client) Run() {
	// 连接时回调
	c.handler.OnConnect(c)
	// 断开连接时的回调
	c.conn.SetCloseHandler(func(code int, text string) error {
		return c.handler.OnClose(c)
	})
	go c.reader()
	go c.writer()
}

// 读取客户端发过来的内容
func (c *Client) reader() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, d, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			break
		}
		var msg Message
		err = json.Unmarshal(d, &msg)
		if err != nil {
			c.conn.Close()
			break
		}
		c.handler.OnMessage(c, &msg)
	}
}

// 向客户端写入内容
func (c *Client) writer() {
	tik := time.NewTicker(pingPeriod)
	for {
		select {
		case buf := <-c.writeChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				c.conn.Close()
				return
			}
		case <-tik.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				c.conn.Close()
				return
			}
		}
	}
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
}

// 获取客户端ID
func (c *Client) Id() int64 {
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
