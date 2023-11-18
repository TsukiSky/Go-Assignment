package optimizedsharedpriorityqueue

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
	replyCount       int
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
		replyCount:       0,
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
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received a request from server %d\n", s.Id, msg.SenderId)
	if s.canReply(msg) {
		// reply to the server
		s.reply(msg.SenderId)
	} else {
		s.toReply = append(s.toReply, msg.SenderId) // add the sender id to the list of servers to reply after executing the critical section
	}
}

// onReceiveReply handles the reply message
func (s *Server) onReceiveReply(msg util2.Message) {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Received reply from %d\n", s.Id, msg.SenderId)
	if s.pendingRequest != nil {
		s.replyCount++
		if s.replyCount == len(s.Connections) {
			go s.executeAndReply()
		}
	}
}

// Execute the critical section
func (s *Server) execute() {
	s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Executing the critical section\n", s.Id)
	time.Sleep(1 * time.Second)
	logger.Logger.Printf("[Server %d] Finished executing the critical section\n", s.Id)
}

// Execute the critical section and release the critical section
func (s *Server) executeAndReply() {
	s.execute()
	s.Queue.Pop()
	s.ResetRequest()
	for _, serverId := range s.toReply {
		s.reply(serverId)
	}
	s.toReply = make([]int, 0) // reset the list of servers to reply
	if s.IsOneTimeRequest {
		s.waitGroup.Done()
	}
}

// Check if the server can reply to the request
func (s *Server) canReply(msg util2.Message) bool {
	// check if the request is at the top of the queue
	if s.pendingRequest == nil {
		return true
	} else {
		return s.pendingRequest.IsLargerThan(msg)
	}
}

// Increment the scalar clock and reply to the server
func (s *Server) reply(toServerId int) {
	clock := s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Replied to server %d\n", s.Id, toServerId)
	reply := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.REPLY,
		ScalarClock: clock,
	}
	go func() {
		s.Connections[toServerId] <- reply
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
			}
		}
	}
}

// ActivateAsPermanentRequester activates the server as a permanent requester
func (s *Server) ActivateAsPermanentRequester() {
	logger.Logger.Printf("[Server %d] Activated as Permanent Requester\n", s.Id)
	go s.Listen()                   // start listening to new messages
	go s.SendRequestWithInterval(5) // send request with interval
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

// Request sends the request
func (s *Server) Request() {
	clock := s.INCREMENT_CLOCK()
	msg := util2.Message{
		SenderId:    s.Id,
		MessageType: util2.REQUEST,
		ScalarClock: clock,
	}
	s.pendingRequest = &msg
	s.Queue.Push(msg)
	s.replyCount = 0
	logger.Logger.Printf("[Server %d] Sent a request to access the critical section\n", s.Id)

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

// ResetRequest resets the pending request and reply count
func (s *Server) ResetRequest() {
	s.pendingRequest = nil
	s.replyCount = 0
}

// return true if the server has pending request, otherwise return false
func (s *Server) hasOngoingRequest() bool {
	if s.pendingRequest != nil {
		return true
	}
	return false
}
