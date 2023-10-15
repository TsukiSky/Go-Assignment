package bully

import "homework/hw1/assignment1/logger"

type Data struct {
	users []user
}

type user struct {
	id   int
	name string
}

type Cluster struct {
	servers     []*Server
	coordinator *Server
	size        int
	index       int
}

func NewCluster() *Cluster {
	return &Cluster{
		servers:     make([]*Server, 0),
		coordinator: nil,
		size:        0,
		index:       0,
	}
}

// AddServer adds a server to this cluster
func (c *Cluster) AddServer(s *Server) {
	for _, server := range c.servers {
		if server == s {
			logger.Logger.Printf("[WARNING] server %d has been added to the cluster", s.id)
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
	returnServers := make([]*Server, len(c.servers)-1)
	for _, server := range c.servers {
		if server.id != id {
			returnServers = append(returnServers, server)
		}
	}
	return returnServers
}

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
