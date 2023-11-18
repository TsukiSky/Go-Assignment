package optimizedsharedpriorityqueue

import (
	"homework/logger"
	"sync"
	"time"
)

type Cluster struct {
	Servers   []*Server
	waitGroup *sync.WaitGroup
}

// NewCluster returns a new cluster
func NewCluster() *Cluster {
	return &Cluster{
		Servers: make([]*Server, 0),
	}
}

// SetWaitGroup sets the wait group for the cluster
func (c *Cluster) SetWaitGroup(group *sync.WaitGroup) {
	c.waitGroup = group
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
func (c *Cluster) Activate(numOfPermanentRequester int) {
	for index, server := range c.Servers {
		if index < numOfPermanentRequester {
			server.ActivateAsPermanentRequester()
		} else {
			server.ActivateAsListener()
		}
	}
	select {}
}

// ActivateInPerformanceComparingMode activates all servers in the cluster in performance comparing mode
// It returns the time duration of the whole process from the start until the last access to the critical section
func (c *Cluster) ActivateInPerformanceComparingMode(numOfOneTimeRequesters int) time.Duration {
	if numOfOneTimeRequesters > len(c.Servers) {
		panic("number of one time requesters is larger than the total number of servers")
	}

	start := time.Now()
	if c.waitGroup == nil {
		panic("wait group is not set")
	}
	for index, server := range c.Servers {
		if index < numOfOneTimeRequesters {
			server.SetWaitGroup(c.waitGroup)
			c.waitGroup.Add(1)
			server.ActivateAsOneTimeRequester()
		} else {
			server.ActivateAsListener()
		}
	}
	c.waitGroup.Wait()
	end := time.Now()
	return end.Sub(start)
}
