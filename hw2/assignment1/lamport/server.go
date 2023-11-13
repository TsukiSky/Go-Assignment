package lamport

import (
	"homework/hw2/assignment1/lamport/util"
	"homework/logger"
	"time"
)

// Server represents a server in the distributed system
type Server struct {
	Id             int
	Channel        chan util.Message
	Connections    map[int]chan util.Message // map of server id to their message channel
	Queue          util.MsgPriorityQueue
	ScalarClock    int
	pendingRequest *util.Message
	replyCount     int
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	scalarClock := 0
	return &Server{
		Id:             id,
		Channel:        make(chan util.Message),
		Connections:    make(map[int]chan util.Message),
		Queue:          *util.NewMsgPriorityQueue(),
		ScalarClock:    scalarClock,
		pendingRequest: nil,
		replyCount:     0,
	}
}

// onReceiveRequest handles the request message
func (s *Server) onReceiveRequest(msg util.Message) {
	s.INCREMENT_CLOCK()
	s.Queue.Push(msg) // push the request to the queue
	logger.Logger.Printf("[Server %d] Received a request from server %d\n", s.Id, msg.SenderId)
	if s.canReply(msg) {
		// reply to the server
		s.reply(msg)
	}
}

// onReceiveReply handles the reply message
func (s *Server) onReceiveReply(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received reply from %d\n", s.Id, msg.SenderId)
	if s.pendingRequest != nil {
		s.replyCount++
		if s.replyCount == len(s.Connections) {
			peek := s.Queue.Peek()
			if peek != nil && peek.Equal(*s.pendingRequest) {
				s.executeAndRelease()
			}
		}
	}
}

// onReceiveRelease handles the release message
func (s *Server) onReceiveRelease(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received release from %d\n", s.Id, msg.SenderId)
	s.Queue.Pop()
	peek := s.Queue.Peek()
	if peek != nil && s.pendingRequest != nil && peek.Equal(*s.pendingRequest) {
		s.executeAndRelease()
	}
}

// Execute the critical section
func (s *Server) execute() {
	s.INCREMENT_CLOCK()
	time.Sleep(1 * time.Second)
}

// Execute the critical section and release the critical section
func (s *Server) executeAndRelease() {
	logger.Logger.Printf("[Server %d] Executing the critical section\n", s.Id)
	s.execute()
	s.Queue.Pop()
	s.ResetRequest()
	s.release()
}

// Release the critical section
func (s *Server) release() {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Released the critical section\n", s.Id)
	release := util.Message{
		SenderId:    s.Id,
		MessageType: util.RELEASE,
		ScalarClock: s.ScalarClock,
	}
	for _, outChannel := range s.Connections {
		outChannel <- release
	}
}

// Check if the server can reply to the request
func (s *Server) canReply(msg util.Message) bool {
	// check if the request is at the top of the queue
	if s.pendingRequest == nil || s.pendingRequest.IsLargerThan(msg) {
		return true
	}
	if s.Id > msg.SenderId {
		return true
	}
	return false
}

// Increment the scalar clock and reply to the server
func (s *Server) reply(msg util.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Replied to server %d\n", s.Id, msg.SenderId)
	clock := s.ScalarClock
	reply := util.Message{
		SenderId:    s.Id,
		MessageType: util.REPLY,
		ScalarClock: clock,
	}
	s.Connections[msg.SenderId] <- reply
}

// Listen listens to the channel and handles the incoming message
func (s *Server) Listen() {
	for {
		select {
		case msg := <-s.Channel:
			switch msg.MessageType {
			case util.REQUEST:
				s.onReceiveRequest(msg)
			case util.REPLY:
				s.onReceiveReply(msg)
			case util.RELEASE:
				s.onReceiveRelease(msg)
			}
		}
	}
}

// Activate activates the server
func (s *Server) Activate() {
	logger.Logger.Printf("[Server %d] Activated\n", s.Id)
	go s.Listen()                   // start listening to new messages
	go s.SendRequestWithInterval(5) // send request with interval
}

// SendRequestWithInterval sends the request with the given interval
func (s *Server) SendRequestWithInterval(second int) {
	for {
		time.Sleep(time.Duration(second) * time.Second)
		if s.hasOngoingRequest() {
			continue
		}

		// make a new request
		s.INCREMENT_CLOCK()
		clock := s.ScalarClock
		msg := util.Message{
			SenderId:    s.Id,
			MessageType: util.REQUEST,
			ScalarClock: clock,
		}
		s.pendingRequest = &msg
		s.Queue.Push(msg)
		s.replyCount = 0
		logger.Logger.Printf("[Server %d] Sent a request to access the critical section\n", s.Id)
		for _, msgChannel := range s.Connections {
			msgChannel <- msg
		}
	}
}

// INCREMENT_CLOCK increase the scalar clock
func (s *Server) INCREMENT_CLOCK() {
	s.ScalarClock++
}

// ResetRequest resets the pending request and reply count
func (s *Server) ResetRequest() {
	//s.pendingRequest = nil
	//s.replyCount = 0
}

// return true if the server has pending request, otherwise return false
func (s *Server) hasOngoingRequest() bool {
	if s.pendingRequest != nil {
		return true
	}
	return false
}
