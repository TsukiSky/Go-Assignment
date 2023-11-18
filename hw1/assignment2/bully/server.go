package bully

import (
	"homework/hw1/logger"
	"math/rand"
	"time"
)

// Role specifies the role of a server
type Role int

const (
	COORDINATOR Role = iota
	WORKER
)

// Status specifies the status of a server
type Status int

const (
	ALIVE Status = iota
	DOWN
)

// Server is the main structure in this implementation
type Server struct {
	id                         int // id of the server
	role                       Role
	status                     Status
	Cluster                    Cluster
	msgChannel                 chan Message // msgChannel is used to receive SynReqMessage, ElectReqMessage, ElectRepMessage and AncMessage
	heartbeatChannel           chan Message // heartbeatChannel is used to receive heartbeat request and reply
	data                       Data         // the data carried in this server
	election                   Election
	heartbeatFrequency         int
	replyTimeout               int
	electionTimeout            int
	syncFrequency              int
	failWhileAnnounce          bool // set to true: this server will fail while announcing itself as the coordinator
	failWhileElection          bool // set to true: this server will fail during the election
	ignorePreviousAnnouncement bool // set to true: this server will ignore the previous announcement
}

// NewServer creates a new server. failWhileAnnounce and failWhileElection are set to false by default
func NewServer(id int, data Data, heartbeatFrequency int, electionTimeout int, replyTimeout int, syncFrequency int) *Server {
	server := Server{
		id:         id,
		role:       WORKER,
		status:     ALIVE,
		msgChannel: make(chan Message),
		data:       data,
		Cluster: Cluster{
			servers:     nil,
			coordinator: nil,
			size:        0,
		},
		election: Election{
			status:        STOP,
			isCoordinator: false,
		},
		heartbeatFrequency:         heartbeatFrequency,
		replyTimeout:               replyTimeout,
		electionTimeout:            electionTimeout,
		syncFrequency:              syncFrequency,
		heartbeatChannel:           make(chan Message),
		failWhileAnnounce:          false,
		failWhileElection:          false,
		ignorePreviousAnnouncement: false,
	}
	return &server
}

// handleMsg handles messages in the msgChannel
func (s *Server) handleMsg(msg Message) {
	if s.status == DOWN {
		return
	}
	switch msg := msg.(type) {
	case *SynReqMessage:
		// Synchronization Request sent by the coordinator
		if s.election.status == RUNNING || s.election.status == PAUSE {
			logger.Logger.Printf("[Server %d] Election is ongoing, data synchronization request is not accepted\n", s.id)
		} else {
			logger.Logger.Printf("[Server %d] Synchronization request from %d, current local time %d, set localtime to %d\n", s.id, msg.content.SenderId, s.data.localTime, msg.data.localTime)
			s.data = msg.data
		}
	case *ElectReqMessage:
		// Election Request sent by some other nodes
		go func() {
			sender := s.Cluster.GetServerById(msg.GetContent().SenderId)
			logger.Logger.Printf("[Server %d] Election request from %d, replying a no message\n", s.id, msg.content.SenderId)
			if !s.sendMsg(NewElectRepMsg(s.id, sender.id), sender.msgChannel) {
				logger.Logger.Printf("[Server %d] Election request from %d, fail to send a no message\n", s.id, msg.content.SenderId)
			}
		}()
		if s.election.status == STOP || s.election.status == PAUSE {
			// this server starts election
			go s.Election(s.electionTimeout)
		}
	case *ElectRepMessage:
		// Election Reply (aka. No) replied by some other nodes
		logger.Logger.Printf("[Server %d] Reply from %d, stop electing %d\n", s.id, msg.content.SenderId, s.id)
		s.election.isCoordinator = false
		s.election.status = PAUSE
	case *AncMessage:
		// Announcement Message sent by the coordinator
		if s.id > msg.GetContent().SenderId || s.ignorePreviousAnnouncement {
			s.ignorePreviousAnnouncement = false
			go s.Election(s.electionTimeout)
		} else {
			logger.Logger.Printf("[Server %d] Announcement received from %d\n", s.id, msg.content.SenderId)
			logger.Logger.Printf("[Server %d] Set coordinator to %d\n", s.id, msg.content.SenderId)
			s.Cluster.SetCoordinator(msg.content.SenderId)
			s.role = WORKER
			s.election.status = STOP
		}
	}
}

// Listen the msgChannel
func (s *Server) Listen() {
	logger.Logger.Printf("[Server Activate] Server %d is activated\n", s.id)
	for {
		select {
		case msg := <-s.msgChannel:
			s.handleMsg(msg)
		}
	}
}

// Activate the server
func (s *Server) Activate() {
	go s.Listen()    // listen to the msgChannel
	go s.Heartbeat() // heartbeat
	go s.Work()      // update local data
}

func (s *Server) PleaseIgnorePreviousAnnouncement() {
	s.ignorePreviousAnnouncement = true
	go s.Listen()    // listen to the msgChannel
	go s.Heartbeat() // heartbeat
	go s.Work()      // update local data
}

func (s *Server) PleaseFailWhileAnnounce() {
	s.failWhileAnnounce = true
	go s.Listen()
	go s.Heartbeat()
	go s.Work()
}

func (s *Server) PleaseFailWhileElection() {
	s.failWhileElection = true
	go s.Listen()
	go s.Heartbeat()
	go s.Work()
}

// Work updates the data stored in the server, in this implementation, the data is a variable called localtime
func (s *Server) Work() {
	syncTimer := time.NewTimer(time.Duration(s.syncFrequency) * time.Second)
	for {
		if s.status == DOWN {
			// Do not proceed if the server status is DOWN
			return
		}
		if s.role == WORKER {
			// introduce some randomness to diverse the data on each server
			if rand.Float64() < 0.5 {
				s.data.localTime += 1
			} else {
				s.data.localTime += 2
			}
		} else {
			select {
			case <-syncTimer.C:
				if s.election.status == RUNNING || s.election.status == PAUSE {
					// Don't synchronization if the election is not over
					logger.Logger.Printf("[Server %d - Coordinator] Election is ongoing, synchronization is not accepted\n", s.id)
				} else {
					currentData := s.data
					for _, server := range s.Cluster.GetAllServersExceptId(s.id) {
						server := server
						go func() {
							logger.Logger.Printf("[Server %d - Coordinator] Synchronize to value %d, sending to %d\n", s.id, currentData.localTime, server.id)
							if !s.sendMsg(NewSynRequestMsg(s.id, server.id, currentData), server.msgChannel) {
								logger.Logger.Printf("[Server %d - Coordinator] Fail to send synchronize message to %d\n", s.id, server.id)
							}
						}()
					}
				}
				syncTimer.Reset(time.Duration(s.syncFrequency) * time.Second) // reset the synchronization clock
			default:
				if rand.Float64() < 0.5 {
					s.data.localTime += 1
				} else {
					s.data.localTime += 2
				}
			}
		}
	}
}

// sendMsg sends a message to a channel
func (s *Server) sendMsg(message Message, messageChannel chan Message) bool {
	select {
	case messageChannel <- message:
		return true
	case <-time.After(time.Duration(s.replyTimeout) * time.Second):
		return false
	}
}

// Heartbeat is the alive-checking mechanism, for servers to know that their coordinator is still alive
func (s *Server) Heartbeat() {
	heartbeatTimer := time.NewTimer(time.Duration(s.heartbeatFrequency) * time.Second)
	heartbeatReplied := true
	for {
		select {
		case heartbeat := <-s.heartbeatChannel:
			if s.status == DOWN {
				return
			}
			switch heartbeat := heartbeat.(type) {
			case *HeartbeatReq: // Heartbeat Request, it is sent by server to the coordinator, "Do a heartbeat so that I can be sure that you are still alive"
				go func() {
					reply := NewHeartbeatRep(s.id, heartbeat.GetAsker())
					logger.Logger.Printf("[Server %d - Coordinator] Sending Heartbeat to %d\n", s.id, heartbeat.GetAsker())
					if !s.sendMsg(reply, s.Cluster.GetServerById(heartbeat.GetAsker()).heartbeatChannel) {
						logger.Logger.Printf("[Server %d - Coordinator] Fail to send heartbeat to %d", s.id, heartbeat.GetAsker())
					}
				}()
			case *HeartbeatRep: // Heartbeat reply, it is sent by the coordinator to the heartbeat requester
				heartbeatReplied = true
			}
		case <-heartbeatTimer.C:
			if s.status == DOWN {
				return
			}
			if (s.election.status == RUNNING || s.election.status == PAUSE) && s.role != COORDINATOR {
				logger.Logger.Printf("[Server %d] Election is ongoing, heartbeat checking is ceased\n", s.id)
			} else {
				if s.Cluster.GetCoordinator() != nil && s.role == WORKER {
					// send heartbeat request if the last time request has been replied
					if heartbeatReplied {
						heartbeatReplied = false
						go func() {
							logger.Logger.Printf("[Server %d] Ask Heartbeat from %d\n", s.id, s.Cluster.GetCoordinator().id)
							if !s.sendMsg(NewHeartbeatReq(s.id, s.Cluster.GetCoordinator().id), s.Cluster.GetCoordinator().heartbeatChannel) {
								logger.Logger.Printf("[Server %d] Fail to Ask Heartbeat from %d, restart election\n", s.id, s.Cluster.GetCoordinator().id)
								go s.Election(s.electionTimeout)
								heartbeatReplied = true
							}
						}()
					} else {
						// coordinator is down, start an election if there is no ongoing election
						go func() {
							if s.election.status == STOP {
								logger.Logger.Printf("[Server %d] Fail to Get Heartbeat from %d, retart election\n", s.id, s.Cluster.GetCoordinator().id)
								go s.Election(s.electionTimeout)
								heartbeatReplied = true
							}
						}()
					}
				} else if s.role != COORDINATOR {
					// coordinator is nil, start an election if there is no ongoing election
					if s.election.status == STOP {
						go s.Election(s.electionTimeout)
					}
				}
			}
			heartbeatTimer.Reset(time.Duration(s.heartbeatFrequency) * time.Second) // reset heartbeat timer
		}

	}
}

// Election elects the server itself as the coordinator
func (s *Server) Election(timeOut int) {
	s.election.status = RUNNING
	s.election.isCoordinator = true // every election starts with a self-voting
	for _, server := range s.Cluster.GetAllServersLargerThanId(s.id) {
		if s.election.status == RUNNING {
			msg := NewElectReqMsg(s.id, server.id)
			logger.Logger.Printf("[Server %d] Sending election message to %d\n", s.id, server.id)
			go s.sendMsg(msg, server.msgChannel)
		}
		if s.failWhileElection {
			s.status = DOWN
			logger.Logger.Printf("[Server %d] Ops... I am DOWN\n", s.id)
			break
		}
	}

	time.Sleep(time.Duration(timeOut) * time.Second) // wait other's reply to the self-election request
	if s.election.status == RUNNING && s.election.isCoordinator {
		s.announce()
	}
}

// announce the server as the coordinator
func (s *Server) announce() {
	if s.status == DOWN {
		return
	}
	s.election.status = STOP
	s.role = COORDINATOR
	s.Cluster.SetCoordinator(s.id)
	logger.Logger.Printf("[Server %d - Coordinator] Announcing Coordinator\n", s.id)
	if s.failWhileAnnounce {
		server := s.Cluster.GetAllServersExceptId(s.id)[0]
		msg := NewAncMsg(s.id, server.id)
		go func() {
			if !s.sendMsg(msg, server.msgChannel) {
				logger.Logger.Printf("[Server %d - Coordinator] Fail to announce to %d\n", s.id, server.id)
			}
		}()
		// Fail
		s.status = DOWN
		logger.Logger.Printf("[Server %d - Coordinator] Opps.. I am DOWN\n", s.id)
	} else {
		for _, server := range s.Cluster.GetAllServersExceptId(s.id) {
			server := server
			go func() {
				msg := NewAncMsg(s.id, server.id)
				if !s.sendMsg(msg, server.msgChannel) {
					logger.Logger.Printf("[Server %d - Coordinator] Fail to announce to %d\n", s.id, server.id)
				}
			}()
		}
	}
}

func (s *Server) SetCluster(cluster Cluster) {
	s.Cluster = cluster
}
