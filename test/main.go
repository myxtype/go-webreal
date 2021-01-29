package main

import (
	"github.com/myxtype/go-webreal"
	"log"
	"time"
)

type PushingBusiness struct {
}

func (p *PushingBusiness) OnConnect(c *webreal.Client) {
	c.Subscribe("hello")
	log.Printf("New client %d", c.Id())
}

func (p *PushingBusiness) OnMessage(c *webreal.Client, msg *webreal.Message) {
	log.Printf("Client %d Message: %v", c.Id(), msg.Data)
}

func (p *PushingBusiness) OnClose(c *webreal.Client) error {
	defer c.UnsubscribeAll()
	log.Printf("Client %d closed.", c.Id())
	return nil
}

func main() {
	var (
		sh   = webreal.NewSubscriptionHub()
		push = PushingBusiness{}
	)
	go func() {
		for {
			tik := time.NewTicker(time.Second)

			select {
			case <-tik.C:
				sh.Publish("hello", []byte("hello"))
			}
		}
	}()
	server := webreal.NewServer(&push, sh)
	server.Run("127.0.0.1:8080", "/ws")
}
