package webreal

import "sync"

// 订阅中心，存储的所有频道和客户端订阅的对应关系
type SubscriptionHub struct {
	subscribers map[string]map[string]*Client
	mu          sync.RWMutex
}

func NewSubscriptionHub() *SubscriptionHub {
	return &SubscriptionHub{
		subscribers: map[string]map[string]*Client{},
	}
}

// 订阅主题
func (s *SubscriptionHub) Subscribe(channel string, client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.subscribers[channel]; !ok {
		s.subscribers[channel] = map[string]*Client{}
	}

	s.subscribers[channel][client.id] = client
	return true
}

// 退订主题
func (s *SubscriptionHub) Unsubscribe(channel string, client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.subscribers[channel]; !ok {
		return false
	}
	if _, ok := s.subscribers[channel][client.id]; !ok {
		return false
	}

	delete(s.subscribers[channel], client.id)
	return true
}

// 向客户端推送主题消息
func (s *SubscriptionHub) Publish(channel string, msg []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.subscribers[channel]; !ok {
		return
	}

	for _, c := range s.subscribers[channel] {
		c.Write(msg)
	}
}
