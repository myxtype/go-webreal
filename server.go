package webreal

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Server struct {
	sh       *SubscriptionHub
	bs       Business
	upgrader websocket.Upgrader
}

func NewServer(bs Business, sh *SubscriptionHub) *Server {
	return &Server{
		sh: sh,
		bs: bs,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	NewClient(conn, s.bs, s.sh, r).Run()
}

func (s *Server) Run(addr string, path string) error {
	http.Handle(path, s)
	return http.ListenAndServe(addr, nil)
}
