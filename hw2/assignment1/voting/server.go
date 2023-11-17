package voting

import (
	"homework/hw2/assignment1/util"
	"homework/logger"
	"sync"
	"time"
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
	IsRescinding bool
	voteTo       *util.Message
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
		IsRescinding: false,
		voteTo:       nil,
	}
}

// onReceiveRequest handles the request message
func (s *Server) onReceiveRequest(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a vote request from server %d\n", s.Id, msg.SenderId)
	if s.voteTo == nil {
		// msg is the only request in the queue
		s.Vote(msg.SenderId)
	} else if s.voteTo.IsLargerThan(msg) {
		// msg is the smallest request in the queue
		s.RescindVote(s.voteTo.SenderId) // rescind and wait for release
	}
	s.Queue.Push(msg) // push the request to the queue
}

// onReceiveVote handles the vote message
func (s *Server) onReceiveVote(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a vote from server %d\n", s.Id, msg.SenderId)
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
	logger.Logger.Printf("[Server %d] Received a rescind request from server %d\n", s.Id, msg.SenderId)
	if s.IsRequesting {
		s.ReleaseVote(msg.SenderId, true)
	}
}

// onReceiveRelease handles the release vote message
func (s *Server) onReceiveRelease(msg util.Message, rescind bool) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a release from server %d\n", s.Id, msg.SenderId)
	if !rescind {
		// received normal release message
		if msg.SenderId == s.voteTo.SenderId {
			s.voteTo = nil
			if s.Queue.Peek() != nil {
				msg := s.Queue.Pop()
				s.voteTo = &msg
				s.Vote(msg.SenderId)
			}
		}
	} else {
		// received rescind release message
		if msg.SenderId == s.voteTo.SenderId {
			s.Queue.Push(*s.voteTo)
			msg := s.Queue.Pop()
			s.voteTo = &msg
			s.Vote(msg.SenderId)
		}
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
	for _, voter := range s.voters {
		s.ReleaseVote(voter, false)
	}
}

// ReleaseVote releases the vote for the given machine id
func (s *Server) ReleaseVote(toMachineId int, rescind bool) {
	s.voters = removeVoter(s.voters, toMachineId)
	msg := util.Message{
		SenderId:    s.Id,
		ScalarClock: s.ScalarClock,
	}
	if rescind {
		msg.MessageType = util.RESCIND_RELEASE
		logger.Logger.Printf("[Server %d] Release rescind vote for server %d\n", s.Id, toMachineId)
	} else {
		msg.MessageType = util.RELEASE
		logger.Logger.Printf("[Server %d] Release vote for server %d\n", s.Id, toMachineId)
	}
	s.Connections[toMachineId] <- msg
}

// RescindVote rescinds the vote for the given message
func (s *Server) RescindVote(toMachineId int) {
	rescindMsg := util.Message{
		SenderId:    s.Id,
		MessageType: util.RESCIND,
		ScalarClock: s.ScalarClock,
	}
	logger.Logger.Printf("[Server %d] Rescind vote for server %d\n", s.Id, toMachineId)
	s.Connections[toMachineId] <- rescindMsg
}

// Vote votes for the given message
func (s *Server) Vote(toMachineId int) {
	voteMsg := util.Message{
		SenderId:    s.Id,
		MessageType: util.VOTE,
		ScalarClock: s.ScalarClock,
	}
	logger.Logger.Printf("[Server %d] Vote for server %d\n", s.Id, toMachineId)
	s.Connections[toMachineId] <- voteMsg
}

// Execute the critical section
func (s *Server) Execute() {
	clock := s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Executing the critical section\n", s.Id, clock)
	time.Sleep(1 * time.Second)
}

// Activate activates the server
func (s *Server) Activate() {
	logger.Logger.Printf("[Server %d] Activated\n", s.Id)
	go s.Listen()
}

// Listen listens to the channel
func (s *Server) Listen() {
	for {
		select {
		case msg := <-s.Channel:
			switch msg.MessageType {
			case util.REQUEST:
				s.onReceiveRequest(msg)
			case util.VOTE:
				s.onReceiveVote(msg)
			case util.RESCIND:
				s.onReceiveRescind(msg)
			case util.RELEASE:
				s.onReceiveRelease(msg, false)
			case util.RESCIND_RELEASE:
				s.onReceiveRelease(msg, true)
			}
		}
	}
}

// INCREMENT_CLOCK increase the scalar clock
func (s *Server) INCREMENT_CLOCK() int {
	s.mu.Lock()
	s.ScalarClock++
	newClock := s.ScalarClock
	s.mu.Unlock()
	return newClock
}
