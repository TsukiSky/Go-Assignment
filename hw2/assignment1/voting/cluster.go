package voting

import "homework/logger"

type Cluster struct {
	Servers []*Server
}

// NewCluster returns a new cluster
func NewCluster() *Cluster {
	return &Cluster{
		Servers: make([]*Server, 0),
	}
}

// AddServer adds a server to the cluster
func (c *Cluster) AddServer(server *Server) {
	for _, s := range c.Servers {
		s.Connections[server.Id] = server.Channel
		server.Connections[s.Id] = s.Channel
	}
	server.ScalarClock = 0
	c.Servers = append(c.Servers, server)
	logger.Logger.Printf("[Cluster ] Server %d added to the cluster\n", server.Id)
}

// Activate activates all servers in the cluster
func (c *Cluster) Activate() {
	for _, server := range c.Servers {
		server.Activate()
	}
}
