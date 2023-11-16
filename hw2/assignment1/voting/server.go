package voting

import (
	"homework/hw2/assignment1/util"
	"sync"
)

type Server struct {
	Id           int
	Channel      chan util.Message
	Connections  map[int]chan util.Message // map of server id to their message channel
	Queue        util.MsgPriorityQueue
	ScalarClock  int
	voters       []int
	IsRequesting bool
	IsExecuting  bool
	mu           sync.Mutex
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	return &Server{
		Id:           id,
		Channel:      make(chan util.Message),
		Connections:  make(map[int]chan util.Message),
		Queue:        *util.NewMsgPriorityQueue(),
		ScalarClock:  0,
		voters:       make([]int, 0),
		IsRequesting: false,
		IsExecuting:  false,
	}
}

// onReceiveRequest handles the request message
func (s *Server) onReceiveRequest(msg util.Message) {
	s.INCREMENT_CLOCK()
	peek := s.Queue.Peek()
	if peek == nil {
		// msg is the only request in the queue
		s.Vote(msg)
	} else if peek.IsLargerThan(msg) {
		// msg is the smallest request in the queue
		s.RescindVote(*peek)
		s.Vote(msg)
	}
	s.Queue.Push(msg) // push the request to the queue
}

// onReceiveVote handles the vote message
func (s *Server) onReceiveVote(msg util.Message) {
	s.INCREMENT_CLOCK()
	if s.IsRequesting {
		s.voters = append(s.voters, msg.SenderId)
		if len(s.voters) >= len(s.Connections)/2 {
			s.Execute()
		}
	}
}

// onReceiveRescind handles the rescind vote message
func (s *Server) onReceiveRescind(msg util.Message) {
	s.INCREMENT_CLOCK()
	if s.IsRequesting {
		s.voters = removeVoter(s.voters, msg.SenderId)
	}
}

// removeVoter removes the voter from the list of voters
func removeVoter(voters []int, voter int) []int {
	for i, v := range voters {
		if v == voter {
			voters[i] = voters[len(voters)-1]
			return voters[:len(voters)-1]
		}
	}
	return voters
}

// ReleaseAllVotes releases all votes
func (s *Server) ReleaseAllVotes() {
	// TODO: implement
}

// ReleaseVote releases the vote for the given machine id
func (s *Server) ReleaseVote(toMachineId int) {
	// TODO: implement
}

// RescindVote rescinds the vote for the given message
func (s *Server) RescindVote(msg util.Message) {
	// TODO: implement
}

// Vote votes for the given message
func (s *Server) Vote(msg util.Message) {
	// TODO: implement
}

// Execute the critical section
func (s *Server) Execute() {
	// TODO: implement
}

func (s *Server) Activate() {
	// TODO: implement
}

func (s *Server) Listen() {
	// TODO: implement
}

// INCREMENT_CLOCK increase the scalar clock
func (s *Server) INCREMENT_CLOCK() int {
	s.mu.Lock()
	s.ScalarClock++
	newClock := s.ScalarClock
	s.mu.Unlock()
	return newClock
}
