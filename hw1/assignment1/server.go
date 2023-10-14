package assignment1

import (
	"fmt"
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

func (s *Server) Initialize() {
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
	fmt.Printf("[Server Activate] -- Clock %d -- Server starts listenning\n", s.clock)
	for {
		select {
		case msg := <-s.channel:
			s.compareAndIncrementClock(msg.clock)
			fmt.Printf("[Server Receive] -- Clock %d -- message from client %d is received\n", s.clock, msg.senderId)
			s.handleMsg(msg)
		}
	}
}

func (s *Server) handleMsg(msg Message) {
	s.incrementClock()
	if flipCoin() {
		// broadcast msg
		fmt.Printf("[Server Broadcast] -- Clock %d -- message from client %d is broadcast\n", s.clock, msg.senderId)
		s.Broadcast(msg)
	} else {
		// discard msg
		fmt.Printf("[Server Discard] -- Clock %d -- message from client %d is discarded\n", s.clock, msg.senderId)
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
