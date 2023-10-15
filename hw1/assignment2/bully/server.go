package bully

import (
	"homework/hw1/assignment1/logger"
	"homework/hw1/assignment2/bully/util"
	"time"
)

type ServerType int

const (
	COORDINATOR ServerType = iota
	WORKER
)

type Server struct {
	id         int
	serverType ServerType
	channel    chan util.Message
	cluster    *Cluster
	data       Data
	election   Election
}

type ElectionStatus int

const (
	RUNNING ElectionStatus = iota
	STOP
)

type Election struct {
	status        ElectionStatus
	isCoordinator bool
}

func NewServer(id int, data Data) *Server {
	server := Server{
		id:         id,
		serverType: WORKER,
		channel:    make(chan util.Message),
		cluster:    nil,
		data:       data,
		election: Election{
			status:        STOP,
			isCoordinator: false,
		},
	}
	return &server
}

func (s *Server) handleMsg(msg util.Message) {
	switch msg := msg.(type) {
	case *util.SynReqMessage:
		msg.GetMessageType()
		// Syn Request Message
		logger.Logger.Printf("Syn Request received")
	case *util.SynRepMessage:
		// Syn Reply Message
		logger.Logger.Printf("Syn Reply received")
	case *util.ElectReqMessage:
		// Elect Request Message
		logger.Logger.Printf("Elect Request received")
		// Reply No
		sender := s.cluster.GetServerById(msg.GetContent().SenderId)
		sender.channel <- util.NewElectRepMsg(s.id, sender.id, false)
	case *util.ElectRepMessage:
		// Elect Reply Message
		logger.Logger.Printf("Elect Reply received")
		if !msg.IsAgree() {
			s.election.isCoordinator = false
		}
	case *util.AncMessage:
		// Announcement Request Message
		logger.Logger.Printf("Announcement Request received")
		s.election.status = STOP
		s.election.isCoordinator = false
	}
	return
}

func (s *Server) Listen() {
	logger.Logger.Printf("[Server Activate] Server %d starts listening\n", s.id)
	for {
		select {
		case msg := <-s.channel:
			s.handleMsg(msg)
		}
	}
}

func (s *Server) Initialize() {
	go s.Listen()
}

func (s *Server) Elect(timeOut int) {
	// every election starts with a self-voting
	s.election.status = RUNNING
	s.election.isCoordinator = true
	for _, server := range s.cluster.GetAllServersLargerThanId(s.id) {
		msg := util.NewElectReqMsg(s.id, server.id)
		server.channel <- msg
	}

	// timeout
	time.Sleep(time.Duration(timeOut) * time.Second)

	if s.election.status == RUNNING && s.election.isCoordinator {
		s.announce()
	}
}

func (s *Server) announce() {
	for _, server := range s.cluster.GetAllServers() {
		msg := util.NewAncMsg(s.id, server.id)
		server.channel <- msg
	}
	s.election.status = STOP
}
