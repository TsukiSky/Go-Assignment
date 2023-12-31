package sharedpriorityqueue

import (
	"homework/hw2/logger"
	util2 "homework/hw2/util"
	"sync"
	"time"
)

// Server represents a server in the distributed system
type Server struct {
	Id               int
	Channel          chan util2.Message
	Connections      map[int]chan util2.Message // map of server id to their message channel
	Queue            util2.MsgPriorityQueue
	ScalarClock      int
	pendingRequest   *util2.Message
	repliedServerIds []int
	toReply          []int
	mu               sync.Mutex
	IsOneTimeRequest bool
	waitGroup        *sync.WaitGroup
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	return &Server{
		Id:               id,
		Channel:          make(chan util2.Message),
		Connections:      make(map[int]chan util2.Message),
		Queue:            *util2.NewMsgPriorityQueue(),
		ScalarClock:      0,
		pendingRequest:   nil,
		repliedServerIds: make([]int, 0),
		toReply:          make([]int, 0),
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
	s.Queue.Push(msg) // push the request to the queue
	logger.Logger.Printf("[Server %d] Received a request from server %d\n", s.Id, msg.SenderId)
	if s.canReply(msg) {
		// reply to the server
		s.reply(msg)
	} else {
		for _, id := range s.repliedServerIds {
			if id == msg.SenderId {
				s.reply(msg)
				return
			}
		}
		s.toReply = append(s.toReply, msg.SenderId) // add the sender id to the list of servers to reply after receiving its reply
	}
}

// onReceiveReply handles the reply message
func (s *Server) onReceiveReply(msg util2.Message) {
	s.COMPARE_AND_INCREMENT_CLOCK(msg.ScalarClock)
	logger.Logger.Printf("[Server %d] Received reply from %d\n", s.Id, msg.SenderId)
	if s.pendingRequest != nil {
		for i, id := range s.toReply {
			if id == msg.SenderId {
				s.reply(msg) // reply to the server
				s.toReply[i] = s.toReply[len(s.toReply)-1]
				s.toReply = s.toReply[:len(s.toReply)-1]
				break
			}
		}

		s.repliedServerIds = append(s.repliedServerIds, msg.SenderId)
		if len(s.repliedServerIds) == len(s.Connections) {
			peek := s.Queue.Peek()
			if peek != nil && peek.Equal(*s.pendingRequest) {
				go s.executeAndRelease()
			}
		}
	}
}

// onReceiveRelease handles the release message
func (s *Server) onReceiveRelease(msg util2.Message) {
	s.COMPARE_AND_INCREMENT_CLOCK(msg.ScalarClock)
	logger.Logger.Printf("[Server %d] Received release from %d\n", s.Id, msg.SenderId)
	s.Queue.Pop()
	peek := s.Queue.Peek()
	if peek != nil && s.pendingRequest != nil && peek.Equal(*s.pendingRequest) {
		go s.executeAndRelease()
	}
}

// Execute the critical section
func (s *Server) execute() {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Executing the critical section\n", s.Id)
	time.Sleep(1 * time.Second)
}

// Execute the critical section and release the critical section
func (s *Server) executeAndRelease() {
	s.execute()
	s.Queue.Pop()
	s.ResetRequest()
	s.release()
	if s.IsOneTimeRequest {
		s.waitGroup.Done()
	}
}

// Release the critical section
func (s *Server) release() {
	clock := s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Released the critical section\n", s.Id)
	release := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.RELEASE,
		ScalarClock: clock,
	}
	for _, outChannel := range s.Connections {
		outChannel <- release
	}
}

// Check if the server can reply to the request
func (s *Server) canReply(msg util2.Message) bool {
	// check if the request is at the top of the queue
	if s.pendingRequest == nil || s.pendingRequest.IsLargerThan(msg) {
		return true
	} else {
		for _, id := range s.repliedServerIds {
			if id == msg.SenderId {
				return true
			}
		}
		return false
	}
}

// Increment the scalar clock and reply to the server
func (s *Server) reply(msg util2.Message) {
	clock := s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Replied to server %d\n", s.Id, msg.SenderId)
	reply := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.REPLY,
		ScalarClock: clock,
	}
	go func() {
		s.Connections[msg.SenderId] <- reply
	}()
}

// Listen listens to the channel and handles the incoming message
func (s *Server) Listen() {
	for {
		select {
		case msg := <-s.Channel:
			switch msg.MessageType {
			case util2.REQUEST:
				s.onReceiveRequest(msg)
			case util2.REPLY:
				s.onReceiveReply(msg)
			case util2.RELEASE:
				s.onReceiveRelease(msg)
			}
		}
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

// SendRequestWithInterval sends the request with the given interval
func (s *Server) SendRequestWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		if s.hasOngoingRequest() {
			continue
		} else {
			// make a new request
			s.Request()
		}
	}
}

func (s *Server) Request() {
	// make a new request
	clock := s.INCREMENT_CLOCK()
	msg := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.REQUEST,
		ScalarClock: clock,
	}
	logger.Logger.Printf("[Server %d] Sent a request to access the critical section\n", s.Id)
	s.pendingRequest = &msg
	s.Queue.Push(msg)
	s.repliedServerIds = make([]int, 0)
	go func() {
		for _, msgChannel := range s.Connections {
			msgChannel <- msg
		}
	}()
}

// INCREMENT_CLOCK increase the scalar clock
func (s *Server) INCREMENT_CLOCK() int {
	s.mu.Lock()
	s.ScalarClock++
	newClock := s.ScalarClock
	s.mu.Unlock()
	return newClock
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

// ResetRequest resets the pending request and reply count
func (s *Server) ResetRequest() {
	s.pendingRequest = nil
}

// return true if the server has pending request, otherwise return false
func (s *Server) hasOngoingRequest() bool {
	if s.pendingRequest != nil {
		return true
	}
	return false
}
