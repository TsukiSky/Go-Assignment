package lamport

import "homework/hw2/assignment1/lamport/util"

// Server represents a server in the distributed system
type Server struct {
	Id             int
	Connections    map[int]Connection
	Queue          util.MsgPriorityQueue
	VectorClock    []int
	pendingRequest *util.Message
	replyCount     int
}

// NewServer returns a new server with the given id
func NewServer(id int) *Server {
	return &Server{
		Id:             id,
		Connections:    make(map[int]Connection),
		Queue:          *util.NewMsgPriorityQueue(),
		VectorClock:    make([]int, 0),
		pendingRequest: nil,
	}
}

func (s *Server) onReceiveRequest(msg util.Message) {
	s.INCREMENT_CLOCK()
	s.Queue.Push(msg) // push the request to the queue
	if s.canReply(msg) {
		// reply to the server
		s.reply(msg)
	}
}

func (s *Server) onReceiveReply(msg util.Message) {
	s.INCREMENT_CLOCK()
	// TODO: implement this
}

// Check if the server can reply to the request
func (s *Server) canReply(msg util.Message) bool {
	// check if the request is at the top of the queue
	if s.pendingRequest == nil || s.pendingRequest.LargerThan(msg) {
		return true
	}
	return false
}

// Increment the vector clock and reply to the server
func (s *Server) reply(msg util.Message) {
	s.INCREMENT_CLOCK()
	clock := make([]int, len(s.VectorClock))
	copy(clock, s.VectorClock)
	reply := util.Message{
		SenderId:    s.Id,
		MessageType: util.REPLY,
		VectorClock: clock,
	}
	s.Connections[msg.SenderId].outChannel <- reply
}

// INCREMENT_CLOCK increase the vector clock
func (s *Server) INCREMENT_CLOCK() {
	s.VectorClock[s.Id]++
}

// Connection represents a connection between two servers
type Connection struct {
	inChannel  chan util.Message
	outChannel chan util.Message
}
