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
		heartbeatChannel:   make(chan Heartbeat),
	}
	return &server
}

func (s *Server) handleMsg(msg Message) {
	switch msg := msg.(type) {
	case *SynReqMessage:
		msg.GetMessageType()
		// Syn Request Message
		fmt.Printf("[Server %d] Synchronization request from %d\n", s.id, msg.content.SenderId)
	case *SynRepMessage:
		// Syn Reply Message
		fmt.Printf("[Server %d] Synchronization reply from %d\n", s.id, msg.content.SenderId)
	case *ElectReqMessage:
		// Elect Request Message
		fmt.Printf("[Server %d] Election request from %d\n", s.id, msg.content.SenderId)
		// Reply No
		sender := s.cluster.GetServerById(msg.GetContent().SenderId)
		sender.msgChannel <- NewElectRepMsg(s.id, sender.id, false)
		if s.election.status == STOP {
			// start election
			s.Elect(s.electionTimeout)
		}
	case *ElectRepMessage:
		// Elect Reply Message
		fmt.Printf("[Server %d] Election reply from %d\n", s.id, msg.content.SenderId)
		if !msg.IsAgree() {
			s.election.isCoordinator = false
		}
	case *AncMessage:
		// Announcement Request Message
		fmt.Printf("[Server %d] Announcement received from %d\n", s.id, msg.content.SenderId)
		fmt.Printf("[Server %d] Set coordinator to %d\n", s.id, msg.content.SenderId)
		s.election.status = STOP
		s.election.isCoordinator = false
		s.cluster.SetCoordinator(msg.content.SenderId)
		s.serverStatus = WORKER
	}
}

func (s *Server) Listen() {
	fmt.Printf("[Server Activate] Server %d is activated\n", s.id)
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
	heartbeatReplied := true
	for {
		select {
		case heartbeat := <-s.heartbeatChannel:
			switch heartbeat := heartbeat.(type) {
			case *HeartbeatReq:
				// heartbeat request
				reply := NewHeartbeatRep(s.id, heartbeat.GetAsker())
				fmt.Printf("[Server %d] Heartbeat check from server %d\n", s.id, heartbeat.GetAsker())
				s.cluster.GetServerById(heartbeat.GetAsker()).heartbeatChannel <- reply
			case *HeartbeatRep:
				fmt.Printf("[Server %d] Heartbeat from server %d\n", s.id, heartbeat.GetBeater())
				heartbeatReplied = true
			}
		case <-heartbeatTimer.C:
			if s.cluster.GetCoordinator() != nil && s.serverStatus == WORKER {
				// send heartbeat request if the last time request has been replied
				if heartbeatReplied {
					heartbeatReplied = false
					fmt.Printf("[Server %d] Heartbeat check to server %d\n", s.id, s.cluster.GetCoordinator().id)
					s.cluster.GetCoordinator().heartbeatChannel <- NewHeartbeatReq(s.id, s.cluster.GetCoordinator().id)
				} else {
					// coordinator is down
					if s.election.status == STOP {
						go s.Elect(s.electionTimeout)
					}
				}
			} else if s.serverStatus != COORDINATOR {
				// coordinator is nil
				if s.election.status == STOP {
					go s.Elect(s.electionTimeout) // start election
				}
			}
			heartbeatTimer.Reset(time.Duration(s.heartbeatFrequency) * time.Second)
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
	s.election.status = STOP
	s.serverStatus = COORDINATOR
	s.cluster.SetCoordinator(s.id)
	for _, server := range s.cluster.GetAllServers() {
		if server.id != s.id {
			msg := NewAncMsg(s.id, server.id)
			server.msgChannel <- msg
		}
	}
}

func (s *Server) SetCluster(cluster Cluster) {
	s.cluster = cluster
}
