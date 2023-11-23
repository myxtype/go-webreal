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
	log.Printf("client %v connected", c.ID())
}

func (p *PushingHandler) OnMessage(c *webreal.Client, msg *webreal.Message) {
	log.Printf("client %v Message: %v", c.ID(), msg.Data)

	switch msg.Type {
	case "close":
		c.Close()
	}
}

func (p *PushingHandler) OnClose(c *webreal.Client) {
	log.Printf("client %v closed.", c.ID())
}

func main() {
	hub := webreal.NewSubscriptionHub()

	go func() {
		tik := time.NewTicker(time.Second)

		for {
			select {
			case <-tik.C:
				for i := 0; i < 10; i++ {
					hub.Publish("hello", []byte("hellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello"))
				}
			}
		}
	}()
	server := webreal.NewServer(&PushingHandler{}, hub, webreal.DefaultConfig())
	server.Run("127.0.0.1:8080", "/ws")
}
