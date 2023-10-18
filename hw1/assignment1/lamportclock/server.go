package lamportclock

import (
	"homework/hw1/assignment1/logger"
	"math/rand"
)

type Server struct {
	clients []*Client
	channel chan Message
	clock   int
}

func NewServer() *Server {
	return &Server{
		clients: make([]*Client, 0),
		channel: make(chan Message),
		clock:   0,
	}
}

func (s *Server) Activate() {
	go s.Listen()
}

func (s *Server) GetClients() []*Client {
	return s.clients
}

// Register a client
func (s *Server) Register(client *Client) {
	s.clients = append(s.clients, client)
	return
}

// Broadcast a message to all clients except the sender
func (s *Server) Broadcast(msg Message) {
	for _, client := range s.clients {
		if client.id != msg.senderId {
			client.channel <- msg
		}
	}
	return
}

func (s *Server) Listen() {
	logger.Logger.Printf("[ Server ] -- Clock %d -- Server starts listening\n", s.clock)
	for {
		select {
		case msg := <-s.channel:
			s.compareAndIncrementClock(msg.clock)
			logger.Logger.Printf("[ Server ] -- Clock %d -- receive message from client %d\n", s.clock, msg.senderId)
			s.handleMsg(msg)
		}
	}
}

func (s *Server) handleMsg(msg Message) {
	s.incrementClock()
	if flipCoin() {
		// broadcast msg
		logger.Logger.Printf("[ Server ] -- Clock %d -- broadcast message from client %d\n", s.clock, msg.senderId)
		broadcastMsg := Message{
			senderId: msg.senderId,
			clock:    s.clock,
		}
		s.Broadcast(broadcastMsg)
	} else {
		// discard msg
		logger.Logger.Printf("[ Server ] -- Clock %d -- discard message from client %d\n", s.clock, msg.senderId)
	}
	return
}

func flipCoin() bool {
	if rand.Float64() < 0.5 {
		return true
	}
	return false
}

// incrementClock increases the clock by 1
func (s *Server) incrementClock() {
	s.clock += 1
}

// compareAndIncrementClock compares the local clock with the incomingClock, chooses the larger clock and increases it by 1
func (s *Server) compareAndIncrementClock(incomingClock int) {
	if s.clock >= incomingClock {
		s.clock += 1
	} else {
		s.clock = incomingClock + 1
	}
}
