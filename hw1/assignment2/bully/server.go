package bully

import (
	"homework/hw1/logger"
	"math/rand"
	"time"
)

type Role int

const (
	COORDINATOR Role = iota
	WORKER
)

type Status int

const (
	ALIVE Status = iota
	DOWN
)

type Server struct {
	id                 int
	role               Role
	status             Status
	Cluster            Cluster
	msgChannel         chan Message
	heartbeatChannel   chan Message
	data               Data
	election           Election
	heartbeatFrequency int
	replyTimeout       int
	electionTimeout    int
	syncFrequency      int
	failWhileAnnounce  bool
	failWhileElection  bool
}

type ElectionStatus int

const (
	RUNNING ElectionStatus = iota
	STOP
	PAUSE
)

type Election struct {
	status        ElectionStatus
	isCoordinator bool
}

func NewServer(id int, data Data, heartbeatFrequency int, electionTimeout int, replyTimeout int, syncFrequency int) *Server {
	server := Server{
		id:         id,
		role:       WORKER,
		status:     ALIVE,
		msgChannel: make(chan Message),
		Cluster: Cluster{
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
		replyTimeout:       replyTimeout,
		electionTimeout:    electionTimeout,
		syncFrequency:      syncFrequency,
		heartbeatChannel:   make(chan Message),
		failWhileAnnounce:  false,
		failWhileElection:  false,
	}
	return &server
}

func (s *Server) handleMsg(msg Message) {
	if s.status == DOWN {
		return
	}
	switch msg := msg.(type) {
	case *SynReqMessage:
		if s.election.status == RUNNING || s.election.status == PAUSE {
			logger.Logger.Printf("[Server %d] Election is ongoing, cease data synchronization\n", s.id)
		} else {
			logger.Logger.Printf("[Server %d] Synchronization Request from %d, current local time %d, set localtime to %d\n", s.id, msg.content.SenderId, s.data.localTime, msg.data.localTime)
			s.data = msg.data
		}
	case *ElectReqMessage:
		go func() {
			sender := s.Cluster.GetServerById(msg.GetContent().SenderId)
			logger.Logger.Printf("[Server %d] Election Request from %d, Replying No\n", s.id, msg.content.SenderId)
			if !s.sendMsg(NewElectRepMsg(s.id, sender.id, false), sender.msgChannel) {
				logger.Logger.Printf("[Server %d] Election Request from %d, fail to send no\n", s.id, msg.content.SenderId)
			}
		}()
		if s.election.status == STOP || s.election.status == PAUSE {
			// start election
			go s.Election(s.electionTimeout)
		}
	case *ElectRepMessage:
		logger.Logger.Printf("[Server %d] Disagree Reply from %d, Stop self-election\n", s.id, msg.content.SenderId)
		if !msg.IsAgree() {
			s.election.isCoordinator = false
			s.election.status = PAUSE
		}
	case *AncMessage:
		// Announcement Request Message
		if s.id > msg.GetContent().SenderId {
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

func (s *Server) Listen() {
	logger.Logger.Printf("[Server Activate] Server %d is activated\n", s.id)
	for {
		if s.status == ALIVE {
			select {
			case msg := <-s.msgChannel:
				s.handleMsg(msg)
			}
		}
	}
}

func (s *Server) Activate() {
	go s.Listen()
	go s.Heartbeat()
	go s.Work()
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

func (s *Server) Work() {
	syncTimer := time.NewTimer(time.Duration(s.syncFrequency) * time.Second)
	for {
		if s.status == DOWN {
			return
		}
		if s.role == WORKER {
			if rand.Float64() < 0.5 {
				s.data.localTime += 1
			} else {
				s.data.localTime += 2
			}
		} else {
			select {
			case <-syncTimer.C:
				if s.election.status == RUNNING || s.election.status == PAUSE {
					logger.Logger.Printf("[Server %d - Coordinator] Election ongoing, cease synchronization\n", s.id)
				} else {
					currentData := s.data
					for _, server := range s.Cluster.GetAllServersExceptId(s.id) {
						server := server
						go func() {
							logger.Logger.Printf("[Server %d - Coordinator] Synchronize to value %d, sending to %d\n", s.id, s.data.localTime, server.id)
							if !s.sendMsg(NewSynRequestMsg(s.id, server.id, currentData), server.msgChannel) {
								logger.Logger.Printf("[Server %d - Coordinator] Fail to send synchronize message to %d\n", s.id, server.id)
							}
						}()
					}
				}
				syncTimer.Reset(time.Duration(s.syncFrequency) * time.Second)
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

func (s *Server) sendMsg(message Message, messageChannel chan Message) bool {
	select {
	case messageChannel <- message:
		return true
	case <-time.After(time.Duration(s.replyTimeout) * time.Second):
		return false
	}
}

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
			case *HeartbeatReq:
				// heartbeat request
				go func() {
					reply := NewHeartbeatRep(s.id, heartbeat.GetAsker())
					logger.Logger.Printf("[Server %d - Coordinator] Sending Heartbeat to %d\n", s.id, heartbeat.GetAsker())
					if !s.sendMsg(reply, s.Cluster.GetServerById(heartbeat.GetAsker()).heartbeatChannel) {
						logger.Logger.Printf("[Server %d - Coordinator] Fail to send heartbeat to %d", s.id, heartbeat.GetAsker())
					}
				}()
			case *HeartbeatRep:
				heartbeatReplied = true
			}
		case <-heartbeatTimer.C:
			if s.status == DOWN {
				return
			}
			if (s.election.status == RUNNING || s.election.status == PAUSE) && s.role != COORDINATOR {
				logger.Logger.Printf("[Server %d] Election is ongoing, cease heartbeat check\n", s.id)
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
						// coordinator is down
						go func() {
							if s.election.status == STOP {
								logger.Logger.Printf("[Server %d] Fail to Get Heartbeat from %d, retart election\n", s.id, s.Cluster.GetCoordinator().id)
								go s.Election(s.electionTimeout)
								heartbeatReplied = true
							}
						}()
					}
				} else if s.role != COORDINATOR {
					// coordinator is nil
					if s.election.status == STOP {
						go s.Election(s.electionTimeout) // start election
					}
				}
			}
			heartbeatTimer.Reset(time.Duration(s.heartbeatFrequency) * time.Second)
		}

	}
}

func (s *Server) Election(timeOut int) {
	// every election starts with a self-voting
	s.election.status = RUNNING
	s.election.isCoordinator = true
	for _, server := range s.Cluster.GetAllServersLargerThanId(s.id) {
		if s.election.status == RUNNING {
			msg := NewElectReqMsg(s.id, server.id)
			logger.Logger.Printf("[Server %d] Sending Election Message to %d\n", s.id, server.id)
			go s.sendMsg(msg, server.msgChannel)
		}
		if s.failWhileElection {
			s.status = DOWN
			logger.Logger.Printf("[Server %d] Opps.. I am DOWN\n", s.id)
			break
		}
	}

	// timeout
	time.Sleep(time.Duration(timeOut) * time.Second)

	if s.election.status == RUNNING && s.election.isCoordinator {
		s.announce()
	}
}

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
