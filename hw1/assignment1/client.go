package assignment1

import (
	"fmt"
	"time"
)

type Client struct {
	id      int
	server  *Server
	channel chan Message
}

func NewClient(id int, server *Server) *Client {
	return &Client{
		id:      id,
		server:  server,
		channel: make(chan Message),
	}
}

func (c *Client) Activate(msgInterval int) {
	fmt.Printf("[Activate] client %d starts listening and sends periodical messages\n", c.id)
	go c.Listen()
	go c.SendMsgWithInterval(msgInterval)
}

func (c *Client) Listen() {
	for {
		select {
		case msg := <-c.channel:
			fmt.Printf("[Client %d] receive %d's message\n", c.id, msg.senderId)
		}
	}
}

func (c *Client) sendMsg(msg Message) {
	fmt.Printf("[Client %d] send message to server\n", msg.senderId)
	c.server.channel <- msg
}

func (c *Client) SendMsgWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		msg := Message{senderId: c.id}
		c.sendMsg(msg)
	}
}
