package main

import (
	"github.com/myxtype/go-webreal"
	"log"
	"time"
)

type PushingHandler struct {
}

func (p *PushingHandler) OnConnect(c *webreal.Client) {
	c.Subscribe("hello")
	log.Printf("New client %d", c.Id())
}

func (p *PushingHandler) OnMessage(c *webreal.Client, msg *webreal.Message) {
	log.Printf("Client %d Message: %v", c.Id(), msg.Data)
}

func (p *PushingHandler) OnClose(c *webreal.Client) error {
	defer c.UnsubscribeAll()
	log.Printf("Client %d closed.", c.Id())
	return nil
}

func main() {
	hub := webreal.NewSubscriptionHub()

	go func() {
		tik := time.NewTicker(time.Second)

		for {
			select {
			case <-tik.C:
				hub.Publish("hello", []byte("hello"))
			}
		}
	}()
	server := webreal.NewServer(&PushingHandler{}, hub)
	server.Run("127.0.0.1:8080", "/ws")
}
