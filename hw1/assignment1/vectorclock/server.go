package vectorclock

import (
	"homework/hw1/assignment1/logger"
	"math/rand"
)

type Server struct {
	clients     []*Client
	channel     chan Message
	vectorClock []int
}

func NewServer() *Server {
	server := Server{
		clients:     make([]*Client, 0),
		channel:     make(chan Message),
		vectorClock: make([]int, 0),
	}
	server.vectorClock = append(server.vectorClock, 0) // register the vectorClock of the server itself
	return &server
}

func (s *Server) Initialize() {
	go s.Listen()
}

func (s *Server) GetClients() []*Client {
	return s.clients
}

// Register a client
func (s *Server) Register(client *Client) {
	client.id = len(s.clients) + 1
	s.clients = append(s.clients, client)
	s.vectorClock = append(s.vectorClock, 0)
	for i := range s.clients {
		s.clients[i].vectorClock = make([]int, len(s.vectorClock))
	}
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
	logger.Logger.Printf("[Server Activate] -- Clock %v -- Server starts listening\n", s.vectorClock)
	for {
		select {
		case msg := <-s.channel:
			// check causality violation
			if s.isCausalityViolation(msg.vectorClock) {
				logger.Logger.Printf("[Potential Causality Violation Detected on Server when receiving %2d's message]\n"+
					"-- Vector Clock on Server -- %v\n"+
					"-- Vector Clock from client %2d-- %v\n", msg.senderId, s.vectorClock, msg.senderId, msg.vectorClock)
			}

			// compare and increase server clock
			s.compareAndIncrementClock(msg.vectorClock)
			logger.Logger.Printf("[ Server  ] -- Clock %v -- receive message from client %2d\n", s.vectorClock, msg.senderId)
			s.handleMsg(msg)
		}
	}
}

func (s *Server) handleMsg(msg Message) {
	s.incrementClock()
	if flipCoin() {
		// broadcast msg
		msg.SetClock(s.vectorClock)
		logger.Logger.Printf("[ Server  ] -- Clock %v -- broadcast message from client %2d\n", s.vectorClock, msg.senderId)
		s.Broadcast(msg)
	} else {
		// discard msg
		logger.Logger.Printf("[ Server  ] -- Clock %v -- discard message from client %d\n", s.vectorClock, msg.senderId)
	}
	return
}

func flipCoin() bool {
	if rand.Float64() < 0.5 {
		return true
	}
	return false
}

func (s *Server) incrementClock() {
	s.vectorClock[0] += 1
}

func (s *Server) compareAndIncrementClock(incomingClock []int) {
	for index, clockValue := range incomingClock {
		if s.vectorClock[index] < clockValue {
			s.vectorClock[index] = clockValue
		}
	}
	s.incrementClock()
}

// check if there is any potential causality violation
func (s *Server) isCausalityViolation(incomingClock []int) bool {
	atLeastOneLarger := false
	for index, clockVal := range incomingClock {
		if s.vectorClock[index] < clockVal {
			return false
		}
		if !atLeastOneLarger && s.vectorClock[index] > clockVal {
			atLeastOneLarger = true
		}
	}
	if atLeastOneLarger {
		return true
	} else {
		return false
	}
}
