package vectorclock

import (
	"homework/hw1/assignment1/logger"
	"sync"
	"time"
)

type Client struct {
	Id          int
	vectorClock []int
	server      *Server
	channel     chan Message
	mu          sync.Mutex // mutex lock used to protect vectorClock
}

func NewClient(id int, server *Server) *Client {
	return &Client{
		Id:          id,
		vectorClock: make([]int, 0),
		server:      server,
		channel:     make(chan Message),
	}
}

func (c *Client) Activate(msgInterval int) {
	logger.Logger.Printf("[Client Activate] -- Clock %v -- client %2d starts listening and sends periodical messages\n", c.vectorClock, c.Id)
	go c.Listen()
	go c.SendMsgWithInterval(msgInterval)
}

func (c *Client) MadlyActivate(msgInterval int) {
	logger.Logger.Printf("[Mad Client Activate] -- Clock %v -- client %2d madly starts listening and sends periodical messages\n", c.vectorClock, c.Id)
	go c.Listen()
	go c.SendMadMsgWithInterval(msgInterval)
}

// Listen to messages sent from the server
func (c *Client) Listen() {
	for {
		select {
		case msg := <-c.channel:
			c.mu.Lock()
			if c.isCausalityViolation(msg.vectorClock) {
				// detect causality violation
				logger.Logger.Printf("\t##############################################################################################################\n"+
					"\t\t\t\t\t\t\t[Potential Causality Violation Detected on Client %2d when receiving %2d's message]\n"+
					"\t\t\t\t\t\t\t-- Vector Clock on client %2d--%v\n"+
					"\t\t\t\t\t\t\t-- Vector Clock from server --%v\n"+
					"\t\t\t\t\t##############################################################################################################\n",
					c.Id, msg.senderId, c.Id, c.vectorClock, msg.vectorClock)
			}
			c.compareAndIncrementClock(msg.vectorClock)
			logger.Logger.Printf("[Client %2d] -- Clock %v -- receive %2d's message\n", c.Id, c.vectorClock, msg.senderId)
			c.mu.Unlock()
		}
	}
}

// sendMsg sends a message to the server
func (c *Client) sendMsg(msg Message) {
	logger.Logger.Printf("[Client %2d] -- Clock %v -- send message to server\n", msg.senderId, msg.vectorClock)
	c.server.channels[c.Id-1] <- msg
}

// SendMsgWithInterval sends periodical messages to the server
func (c *Client) SendMsgWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		c.mu.Lock()
		c.incrementClock()
		clock := make([]int, len(c.vectorClock))
		copy(clock, c.vectorClock)
		msg := Message{senderId: c.Id, vectorClock: clock}
		c.mu.Unlock()
		c.sendMsg(msg)
	}
}

// SendMadMsgWithInterval sends malicious causality violated periodical messages to the server
func (c *Client) SendMadMsgWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		c.mu.Lock()
		c.incrementClock()
		clockSmall := make([]int, len(c.vectorClock))
		copy(clockSmall, c.vectorClock)
		msgSmall := Message{senderId: c.Id, vectorClock: clockSmall} // construct a message with a smaller vector clock

		c.incrementClock()
		clockLarge := make([]int, len(c.vectorClock))
		copy(clockLarge, c.vectorClock)
		msgLarge := Message{senderId: c.Id, vectorClock: clockLarge} // construct a message with a larger vector clock
		c.mu.Unlock()

		// send the messages in a wrong order
		c.sendMsg(msgLarge)
		c.sendMsg(msgSmall)
	}
}

func (c *Client) incrementClock() {
	c.vectorClock[c.Id] += 1
}

func (c *Client) compareAndIncrementClock(incomingClock []int) {
	for index, clockValue := range incomingClock {
		if c.vectorClock[index] < clockValue {
			c.vectorClock[index] = clockValue
		}
	}
	c.vectorClock[c.Id] += 1
}

// check if there is any potential causality violation
func (c *Client) isCausalityViolation(incomingClock []int) bool {
	atLeastOneLarger := false
	for index, clockVal := range incomingClock {
		if c.vectorClock[index] < clockVal {
			return false
		}
		if !atLeastOneLarger && c.vectorClock[index] > clockVal {
			atLeastOneLarger = true
		}
	}
	if atLeastOneLarger {
		return true
	} else {
		return false
	}
}
