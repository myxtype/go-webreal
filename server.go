package webreal

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Server struct {
	hub      *SubscriptionHub
	handler  Handler
	upgrader websocket.Upgrader
}

func NewServer(handler Handler, hub *SubscriptionHub) *Server {
	return &Server{
		hub:     hub,
		handler: handler,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// 实现http.Handler接口
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	NewClient(conn, s.handler, s.hub, r).Run()
}

// 使用默认的http启动监听服务
func (s *Server) Run(addr string, path string) error {
	http.Handle(path, s)
	return http.ListenAndServe(addr, nil)
}
