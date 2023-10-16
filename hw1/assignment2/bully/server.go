package bully

import (
	"fmt"
	"time"
)

type ServerStatus int

const (
	COORDINATOR ServerStatus = iota
	WORKER
)

type Server struct {
	id                 int
	serverStatus       ServerStatus
	msgChannel         chan Message
	cluster            Cluster
	data               Data
	election           Election
	heartbeatChannel   chan Heartbeat
	heartbeatFrequency int
	electionTimeout    int
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

func NewServer(id int, data Data, heartbeatFrequency int, electionTimeout int) *Server {
	server := Server{
		id:           id,
		serverStatus: WORKER,
		msgChannel:   make(chan Message),
		cluster: Cluster{
			servers:     nil,
			coordinator: nil,
			size:        0,
		},
		data: data,
		election: Election{
			status:        STOP,
			isCoordinator: false,
		},
		heartbeatFrequency: heartbeatFrequency,
		electionTimeout:    electionTimeout,
	}
	return &server
}

func (s *Server) handleMsg(msg Message) {
	switch msg := msg.(type) {
	case *SynReqMessage:
		msg.GetMessageType()
		// Syn Request Message
		fmt.Printf("%d Syn Request received\n", s.id)
	case *SynRepMessage:
		// Syn Reply Message
		fmt.Printf("%d Syn Reply received\n", s.id)
	case *ElectReqMessage:
		// Elect Request Message
		fmt.Printf("%d Elect Request received\n", s.id)
		// Reply No
		sender := s.cluster.GetServerById(msg.GetContent().SenderId)
		sender.msgChannel <- NewElectRepMsg(s.id, sender.id, false)
		if s.election.status == STOP {
			s.Elect(s.electionTimeout)
		}
	case *ElectRepMessage:
		// Elect Reply Message
		fmt.Printf("%d Elect Reply received\n", s.id)
		if !msg.IsAgree() {
			s.election.isCoordinator = false
		}
	case *AncMessage:
		// Announcement Request Message
		fmt.Printf("%d Announcement Request received\n", s.id)
		s.election.status = STOP
		s.election.isCoordinator = false
		s.cluster.SetCoordinator(msg.content.SenderId)
		fmt.Printf("%d Set coordinator to %d\n", s.id, s.cluster.coordinator.id)
	}
}

func (s *Server) Listen() {
	fmt.Printf("[Server Activate] Server %d starts listening\n", s.id)
	for {
		select {
		case msg := <-s.msgChannel:
			s.handleMsg(msg)
		}
	}
}

func (s *Server) Activate() {
	go s.Listen()
	go s.Heartbeat()
}

func (s *Server) Heartbeat() {
	heartbeatTimer := time.NewTimer(time.Duration(s.heartbeatFrequency) * time.Second)
	for {
		select {
		case heartbeat := <-s.heartbeatChannel:
			switch heartbeat := heartbeat.(type) {
			case *HeartbeatReq:
				// heartbeat request
				fmt.Printf("%d heartbeat request from %d received\n", s.id, s.cluster.GetCoordinator().id)
				reply := NewHeartbeatRep(s.id, heartbeat.GetAsker())
				s.cluster.GetServerById(heartbeat.GetAsker()).heartbeatChannel <- reply
			case *HeartbeatRep:
				heartbeatTimer.Reset(time.Duration(s.heartbeatFrequency) * time.Second)
			}
		case <-heartbeatTimer.C:
			// coordinator might be down
			if s.cluster.GetCoordinator() != nil && s.cluster.GetCoordinator().id != s.id {
				fmt.Printf("%d heartbeat request sent to %d\n", s.id, s.cluster.GetCoordinator().id)
				s.cluster.GetCoordinator().heartbeatChannel <- NewHeartbeatReq(s.id, s.cluster.GetCoordinator().id)
			} else {
				// coordinator does not exist
				if s.election.status == STOP {
					go s.Elect(s.electionTimeout)
				}
			}
		}
	}
}

func (s *Server) Elect(timeOut int) {
	// every election starts with a self-voting
	s.election.status = RUNNING
	s.election.isCoordinator = true
	for _, server := range s.cluster.GetAllServersLargerThanId(s.id) {
		msg := NewElectReqMsg(s.id, server.id)
		server.msgChannel <- msg
	}

	// timeout
	time.Sleep(time.Duration(timeOut) * time.Second)

	if s.election.status == RUNNING && s.election.isCoordinator {
		s.announce()
	}
}

func (s *Server) announce() {
	for _, server := range s.cluster.GetAllServers() {
		msg := NewAncMsg(s.id, server.id)
		server.msgChannel <- msg
	}
	s.election.status = STOP
	s.serverStatus = COORDINATOR
}

func (s *Server) SetCluster(cluster Cluster) {
	s.cluster = cluster
}
