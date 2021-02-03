package webreal

import "sync"

type Channel struct {
	mu          sync.RWMutex
	subscribers map[string]*Client
}

func NewChannel() *Channel {
	return &Channel{
		subscribers: map[string]*Client{},
	}
}

// 添加客户端
func (c *Channel) Add(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.subscribers[client.id] = client
}

// 移除客户端
func (c *Channel) Remove(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.subscribers, client.id)
}

// 遍历客户端
func (c *Channel) Range(f func(client *Client)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, client := range c.subscribers {
		f(client)
	}
}
