package bully

import (
	"homework/logger"
)

// Data to be stored in the server
type Data struct {
	localTime int
}

func NewData() Data {
	return Data{localTime: 0}
}

// ElectionStatus and Election
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

// Cluster contains a group of servers
type Cluster struct {
	servers     []*Server
	coordinator *Server
	size        int
}

func NewCluster(servers []*Server) Cluster {
	return Cluster{
		servers:     servers,
		coordinator: nil,
		size:        len(servers),
	}
}

// AddServer adds a server to this Cluster
func (c *Cluster) AddServer(s *Server) {
	for _, server := range c.servers {
		if server == s {
			logger.Logger.Printf("[WARNING] server %d has been added to the Cluster", s.id)
			return
		}
	}
	c.servers = append(c.servers, s)
}

func (c *Cluster) GetAllServers() []*Server {
	returnServers := make([]*Server, len(c.servers))
	copy(returnServers, c.servers)
	return returnServers
}

func (c *Cluster) GetServerById(id int) *Server {
	for _, server := range c.servers {
		if server.id == id {
			return server
		}
	}
	return nil
}

func (c *Cluster) GetAllServersExceptId(id int) []*Server {
	returnServers := make([]*Server, 0)
	for _, server := range c.servers {
		if server.id != id {
			returnServers = append(returnServers, server)
		}
	}
	return returnServers
}

// GetAllServersLargerThanId gets all servers that has a larger id than the input id
func (c *Cluster) GetAllServersLargerThanId(id int) []*Server {
	returnServers := make([]*Server, 0)
	for _, server := range c.servers {
		if server.id > id {
			returnServers = append(returnServers, server)
		}
	}
	return returnServers
}

func (c *Cluster) GetCoordinator() *Server {
	return c.coordinator
}

func (c *Cluster) SetCoordinator(coordinatorId int) *Server {
	for _, server := range c.servers {
		if server.id == coordinatorId {
			c.coordinator = server
			return server
		}
	}
	return nil
}
