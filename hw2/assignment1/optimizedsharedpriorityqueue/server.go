package optimizedsharedpriorityqueue

import (
	"homework/hw2/assignment1/util"
	"homework/logger"
	"sync"
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
	toReply        []int
	mu             sync.Mutex
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	return &Server{
		Id:             id,
		Channel:        make(chan util.Message),
		Connections:    make(map[int]chan util.Message),
		Queue:          *util.NewMsgPriorityQueue(),
		ScalarClock:    0,
		pendingRequest: nil,
		replyCount:     0,
		toReply:        make([]int, 0),
	}
}

// onReceiveRequest handles the request message
func (s *Server) onReceiveRequest(msg util.Message) {
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
func (s *Server) onReceiveReply(msg util.Message) {
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
	clock := s.INCREMENT_CLOCK()
	logger.Logger.Printf("[Server %d] Executing the critical section\n", s.Id, clock)
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
}

// Check if the server can reply to the request
func (s *Server) canReply(msg util.Message) bool {
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
	reply := util.Message{
		SenderId:    s.Id,
		MessageType: util.REPLY,
		ScalarClock: clock,
	}
	s.Connections[toServerId] <- reply
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
		clock := s.INCREMENT_CLOCK()
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
func (s *Server) INCREMENT_CLOCK() int {
	s.mu.Lock()
	s.ScalarClock++
	newClock := s.ScalarClock
	s.mu.Unlock()
	return newClock
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
