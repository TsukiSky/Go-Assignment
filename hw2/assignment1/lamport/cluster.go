package lamport

type Cluster struct {
	Servers []*Server
}

func NewCluster() *Cluster {
	return &Cluster{
		Servers: make([]*Server, 0),
	}
}

// AddServer adds a server to the cluster
func (c *Cluster) AddServer(server *Server) {
	for _, s := range c.Servers {
		s.Connections[server.Id] = server.Channel
		s.VectorClock = append(s.VectorClock, 0)
		server.Connections[s.Id] = s.Channel
		server.VectorClock = append(server.VectorClock, 0)
	}
	c.Servers = append(c.Servers, server)
}

func (c *Cluster) Activate() {
	for _, server := range c.Servers {
		server.Activate()
	}
}
