package webreal

import "sync"

// 订阅中心，存储的所有频道和客户端订阅的对应关系
type SubscriptionHub struct {
	mu       sync.Mutex
	channels map[string]*Channel
}

func NewSubscriptionHub() *SubscriptionHub {
	return &SubscriptionHub{
		channels: map[string]*Channel{},
	}
}

// 订阅主题
func (s *SubscriptionHub) Subscribe(channel string, client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.channels[channel]; !found {
		s.channels[channel] = NewChannel()
	}

	s.channels[channel].Add(client)
}

// 退订主题
func (s *SubscriptionHub) Unsubscribe(channel string, client *Client) {
	if _, found := s.channels[channel]; !found {
		return
	}

	s.channels[channel].Remove(client)
}

// 向客户端推送主题消息
func (s *SubscriptionHub) Publish(channel string, msg []byte) {
	if _, found := s.channels[channel]; !found {
		return
	}

	s.channels[channel].Range(func(c *Client) {
		c.Write(msg)
	})
}
