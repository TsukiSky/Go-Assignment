package voting

import (
	"homework/hw2/assignment1/util"
	"homework/logger"
	"sync"
	"time"
)

type Server struct {
	Id              int
	Channel         chan util.Message
	Connections     map[int]chan util.Message // map of server id to their message channel
	Queue           util.MsgPriorityQueue
	ScalarClock     int
	voters          []int
	IsRequesting    bool
	IsExecuting     bool
	IsRescinding    bool
	voteTo          *util.Message
	archivedRescind []int
	mu              sync.Mutex
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	connection := make(map[int]chan util.Message)
	channel := make(chan util.Message)
	connection[id] = channel
	return &Server{
		Id:              id,
		Channel:         channel,
		Connections:     connection,
		Queue:           *util.NewMsgPriorityQueue(),
		ScalarClock:     0,
		voters:          make([]int, 0),
		IsRequesting:    false,
		IsExecuting:     false,
		IsRescinding:    false,
		voteTo:          nil,
		archivedRescind: make([]int, 0),
	}
}

// onReceiveRequest handles the request message
func (s *Server) onReceiveRequest(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a vote request from server %d\n", s.Id, msg.SenderId)
	if s.voteTo == nil {
		// msg is the only request in the queue
		s.voteTo = &msg
		s.Vote(msg.SenderId)
	} else if s.voteTo.IsLargerThan(msg) {
		// msg is the smallest request in the queue
		s.RescindVote(s.voteTo.SenderId) // rescind and wait for release
		s.Queue.Push(msg)                // push the request to the queue
	} else {
		// msg is not the smallest request in the queue
		s.Queue.Push(msg) // push the request to the queue
	}

}

// onReceiveVote handles the vote message
func (s *Server) onReceiveVote(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a vote from server %d\n", s.Id, msg.SenderId)
	if s.IsRequesting {
		// check if the vote is rescinded
		for i, sender := range s.archivedRescind {
			if sender == msg.SenderId {
				s.archivedRescind[i] = s.archivedRescind[len(s.archivedRescind)-1]
				s.archivedRescind = s.archivedRescind[:len(s.archivedRescind)-1]
				s.ReleaseVote(msg.SenderId, true)
				return
			}
		}

		s.voters = append(s.voters, msg.SenderId)
		if len(s.voters) > len(s.Connections)/2 {
			s.ExecuteAndRelease()
		}
	} else if s.IsExecuting {
		s.voters = append(s.voters, msg.SenderId)
	} else {
		s.ReleaseVote(msg.SenderId, false)
	}
}

// onReceiveRescind handles the rescind vote message
func (s *Server) onReceiveRescind(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a rescind request from server %d\n", s.Id, msg.SenderId)

	if !s.IsRequesting || s.IsExecuting {
		logger.Logger.Printf("[Server %d] Rescind request from server %d is ignored\n", s.Id, msg.SenderId)
		return
	}

	for _, voter := range s.voters {
		if voter == msg.SenderId {
			s.ReleaseVote(msg.SenderId, true)
			return
		}
	}

	s.archivedRescind = append(s.archivedRescind, msg.SenderId) // add the rescind to the list of archived rescinds
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
	clock := s.INCREMENT_CLOCK()
	s.voters = removeVoter(s.voters, toMachineId)
	msg := util.Message{
		SenderId:    s.Id,
		ScalarClock: clock,
	}
	if rescind {
		msg.MessageType = util.RESCIND_RELEASE
		logger.Logger.Printf("[Server %d] Release rescind vote to server %d\n", s.Id, toMachineId)
	} else {
		msg.MessageType = util.RELEASE
		logger.Logger.Printf("[Server %d] Release vote to server %d\n", s.Id, toMachineId)
	}
	go func() {
		s.Connections[toMachineId] <- msg
	}()
}

// RescindVote rescinds the vote for the given message
func (s *Server) RescindVote(toMachineId int) {
	rescindMsg := util.Message{
		SenderId:    s.Id,
		MessageType: util.RESCIND,
		ScalarClock: s.ScalarClock,
	}
	logger.Logger.Printf("[Server %d] Rescind vote to server %d\n", s.Id, toMachineId)
	go func() {
		s.Connections[toMachineId] <- rescindMsg
	}()
}

// Vote votes for the given message
func (s *Server) Vote(toMachineId int) {
	voteMsg := util.Message{
		SenderId:    s.Id,
		MessageType: util.VOTE,
		ScalarClock: s.ScalarClock,
	}
	logger.Logger.Printf("[Server %d] Vote for server %d\n", s.Id, toMachineId)
	go func() {
		s.Connections[toMachineId] <- voteMsg
	}()
}

// Execute the critical section
func (s *Server) Execute() {
	s.IsRequesting = false
	s.IsExecuting = true
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Executing the critical section\n", s.Id)
	time.Sleep(1 * time.Second)
	s.IsExecuting = false
}

func (s *Server) ExecuteAndRelease() {
	s.Execute()
	s.ReleaseAllVotes()
}

// Activate activates the server
func (s *Server) Activate() {
	logger.Logger.Printf("[Server %d] Activated\n", s.Id)
	go s.Listen()
	go s.RequestWithInterval(5)
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

func (s *Server) Request() {
	clock := s.INCREMENT_CLOCK()
	s.IsRequesting = true
	logger.Logger.Printf("[Server %d] Sent a request to access the critical section\n", s.Id)

	reqMsg := util.Message{
		SenderId:    s.Id,
		MessageType: util.REQUEST,
		ScalarClock: clock,
	}

	go func() {
		for _, connection := range s.Connections {
			connection <- reqMsg
		}
	}()
}

// RequestWithInterval sends the request with the given interval
func (s *Server) RequestWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		if s.IsRequesting || s.IsExecuting {
			continue
		} else {
			// make a new request
			s.Request()
		}
	}
	//time.Sleep(time.Duration(second) * time.Second)
	// make a new request
	//s.Request()
}
