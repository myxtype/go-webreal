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
	server := webreal.NewServer(&PushingBusiness{}, hub)
	server.Run("127.0.0.1:8080", "/ws")
}
