package assignment1

import (
	"fmt"
	"math/rand"
)

type Server struct {
	clients []*Client
	channel chan Message
}

func NewServer() *Server {
	return &Server{
		clients: make([]*Client, 0),
		channel: make(chan Message),
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
	fmt.Printf("[Server] Server starts listen\n")
	for {
		select {
		case msg := <-s.channel:
			s.handleMsg(msg)
		}
	}
}

func (s *Server) handleMsg(msg Message) {
	if flipCoin() {
		// broadcast msg
		fmt.Printf("[Broadcast] message from client %d is broadcast\n", msg.senderId)
		s.Broadcast(msg)
	} else {
		// discard msg
		fmt.Printf("[Discard] message from client %d is discarded\n", msg.senderId)
	}
	return
}

func flipCoin() bool {
	if rand.Float64() < 0.5 {
		return true
	}
	return false
}
