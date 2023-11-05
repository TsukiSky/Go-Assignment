package lamportclock

import (
	"homework/logger"
	"sync"
	"time"
)

type Client struct {
	id      int
	clock   int
	server  *Server
	channel chan Message
	mu      sync.Mutex
}

func NewClient(id int, server *Server) *Client {
	return &Client{
		id:      id,
		clock:   0,
		server:  server,
		channel: make(chan Message),
	}
}

func (c *Client) Activate(msgInterval int) {
	logger.Logger.Printf("[Client %d] -- Clock %d -- client %d starts listening and sends periodical messages\n", c.id, c.clock, c.id)
	go c.Listen()
	go c.SendMsgWithInterval(msgInterval)
}

// compareAndIncrementClock compares the local clock with the incomingClock, chooses the larger clock and increases it by 1
func (c *Client) compareAndIncrementClock(incomingClock int) {
	c.mu.Lock()
	if c.clock >= incomingClock {
		c.clock += 1
	} else {
		c.clock = incomingClock + 1
	}
	c.mu.Unlock()
}

// Listen to messages sent from the server
func (c *Client) Listen() {
	for {
		select {
		case msg := <-c.channel:
			c.compareAndIncrementClock(msg.clock)
			logger.Logger.Printf("[Client %d] -- Clock %d -- receive server's broadcast message, originally from client %d\n", c.id, c.clock, msg.senderId)
		}
	}
}

// sendMsg to the server
func (c *Client) sendMsg(msg Message) {
	logger.Logger.Printf("[Client %d] -- Clock %d -- send message to server\n", msg.senderId, c.clock)
	c.server.channel <- msg
}

// SendMsgWithInterval to the server
func (c *Client) SendMsgWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		c.compareAndIncrementClock(0)
		msg := Message{senderId: c.id, clock: c.clock}
		c.sendMsg(msg)
	}
}
