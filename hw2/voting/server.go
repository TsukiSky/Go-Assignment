package voting

import (
	"homework/hw2/logger"
	util2 "homework/hw2/util"
	"sync"
	"time"
)

type Server struct {
	Id               int
	Channel          chan util2.Message
	Connections      map[int]chan util2.Message // map of server id to their message channel
	Queue            util2.MsgPriorityQueue
	ScalarClock      int
	voters           []int
	IsRequesting     bool
	IsExecuting      bool
	IsRescinding     bool
	voteTo           *util2.Message
	archivedRescind  []int
	mu               sync.Mutex
	IsOneTimeRequest bool
	waitGroup        *sync.WaitGroup
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	connection := make(map[int]chan util2.Message)
	channel := make(chan util2.Message)
	connection[id] = channel
	return &Server{
		Id:               id,
		Channel:          channel,
		Connections:      connection,
		Queue:            *util2.NewMsgPriorityQueue(),
		ScalarClock:      0,
		voters:           make([]int, 0),
		IsRequesting:     false,
		IsExecuting:      false,
		IsRescinding:     false,
		voteTo:           nil,
		archivedRescind:  make([]int, 0),
		IsOneTimeRequest: false,
	}
}

// SetWaitGroup sets the wait group
func (s *Server) SetWaitGroup(group *sync.WaitGroup) {
	s.waitGroup = group
}

// onReceiveRequest handles the request message
func (s *Server) onReceiveRequest(msg util2.Message) {
	s.COMPARE_AND_INCREMENT_CLOCK(msg.ScalarClock)
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
func (s *Server) onReceiveVote(msg util2.Message) {
	s.COMPARE_AND_INCREMENT_CLOCK(msg.ScalarClock)
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
func (s *Server) onReceiveRescind(msg util2.Message) {
	s.COMPARE_AND_INCREMENT_CLOCK(msg.ScalarClock)
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
func (s *Server) onReceiveRelease(msg util2.Message, rescind bool) {
	s.COMPARE_AND_INCREMENT_CLOCK(msg.ScalarClock)
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
	msg := util2.Message{
		SenderId:    s.Id,
		ScalarClock: clock,
	}
	if rescind {
		msg.MessageType = util2.RESCIND_RELEASE
		logger.Logger.Printf("[Server %d] Release rescind vote to server %d\n", s.Id, toMachineId)
	} else {
		msg.MessageType = util2.RELEASE
		logger.Logger.Printf("[Server %d] Release vote to server %d\n", s.Id, toMachineId)
	}
	go func() {
		s.Connections[toMachineId] <- msg
	}()
}

// RescindVote rescinds the vote for the given message
func (s *Server) RescindVote(toMachineId int) {
	rescindMsg := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.RESCIND,
		ScalarClock: s.ScalarClock,
	}
	logger.Logger.Printf("[Server %d] Rescind vote to server %d\n", s.Id, toMachineId)
	go func() {
		s.Connections[toMachineId] <- rescindMsg
	}()
}

// Vote votes for the given message
func (s *Server) Vote(toMachineId int) {
	voteMsg := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.VOTE,
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
	if s.IsOneTimeRequest {
		s.waitGroup.Done()
	}
}

// ActivateAsPermanentRequester activates the server as a permanent requester
func (s *Server) ActivateAsPermanentRequester() {
	logger.Logger.Printf("[Server %d] Activated as Permanent Requester\n", s.Id)
	go s.Listen()
	go s.SendRequestWithInterval(5)
}

// ActivateAsListener activates the server as a listener
func (s *Server) ActivateAsListener() {
	logger.Logger.Printf("[Server %d] Activated as Listener\n", s.Id)
	go s.Listen()
}

// ActivateAsOneTimeRequester activates the server as a one-time requester
func (s *Server) ActivateAsOneTimeRequester() {
	logger.Logger.Printf("[Server %d] Activated as One-time Requester\n", s.Id)
	s.IsOneTimeRequest = true
	go s.Listen()
	go s.Request()
}

// Listen listens to the channel
func (s *Server) Listen() {
	for {
		select {
		case msg := <-s.Channel:
			switch msg.MessageType {
			case util2.REQUEST:
				s.onReceiveRequest(msg)
			case util2.VOTE:
				s.onReceiveVote(msg)
			case util2.RESCIND:
				s.onReceiveRescind(msg)
			case util2.RELEASE:
				s.onReceiveRelease(msg, false)
			case util2.RESCIND_RELEASE:
				s.onReceiveRelease(msg, true)
			}
		}
	}
}

// COMPARE_AND_INCREMENT_CLOCK compares the incoming clock with the server's scalar clock and increment the scalar clock
func (s *Server) COMPARE_AND_INCREMENT_CLOCK(incomingClock int) {
	s.mu.Lock()
	if s.ScalarClock >= incomingClock {
		s.ScalarClock++
	} else {
		s.ScalarClock = incomingClock + 1
	}
	s.mu.Unlock()
}

// INCREMENT_CLOCK increase the scalar clock
func (s *Server) INCREMENT_CLOCK() int {
	s.mu.Lock()
	s.ScalarClock++
	newClock := s.ScalarClock
	s.mu.Unlock()
	return newClock
}

// Request sends the request
func (s *Server) Request() {
	clock := s.INCREMENT_CLOCK()
	s.IsRequesting = true
	logger.Logger.Printf("[Server %d] Sent a request to access the critical section\n", s.Id)

	reqMsg := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.REQUEST,
		ScalarClock: clock,
	}

	go func() {
		for _, connection := range s.Connections {
			connection <- reqMsg
		}
	}()
}

// SendRequestWithInterval sends the request with the given interval
func (s *Server) SendRequestWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		if s.IsRequesting || s.IsExecuting {
			continue
		} else {
			// make a new request
			s.Request()
		}
	}
}
