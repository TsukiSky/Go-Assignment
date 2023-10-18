package main

import (
	"homework/hw1/assignment2/bully"
	"homework/hw1/logger"
)

// Use Bully algorithm to implement a working version of replica synchronization.
// You may assume that
//  1. each replica maintains some data structure (that may diverge for arbitrary reasons),
//     which are periodically synchronized with the coordinator.
//  2. The coordinator initiates the synchronization process by sending message to all other machines.
//  3. Upon receiving the message from the coordinator, each machine updates its local version of the
//     data structure with the coordinator’s version.
//  4. The coordinator, being an arbitrary machine in the network, is subject to fault. Thus,
//     4.1 a new coordinator is chosen by the Bully algorithm.
//     4.2	You can assume a fixed timeout to simulate the behaviour of detecting a fault
//     The objective is to have a consensus across all machines (simulated by GO routines)
//     in terms of the newly elected coordinator.

const (
	numOfServers       = 5  // number of servers in the cluster
	heartbeatFrequency = 10 // intervals between heartbeat (seconds)
	electionTimeout    = 8  // intervals from request -> announce (seconds)
	replyTimeout       = 4
	syncFrequency      = 10 // intervals of synchronization (seconds)
)

func main() {
	logger.Init("assignment_2.log", "assignment 2: ")
	servers := make([]*bully.Server, 0)
	// create servers
	for i := 0; i < numOfServers; i++ {
		servers = append(servers, bully.NewServer(i, bully.NewData(), heartbeatFrequency, electionTimeout, replyTimeout, syncFrequency))
	}
	// initialize servers
	for _, server := range servers {
		server.SetCluster(bully.NewCluster(servers))
	}

	// 1. Normal simulation
	for _, server := range servers {
		server.Activate()
	}

	// 2.1 Simulate worst-case
	// The worst-case is: the node with the smallest id start the election.
	// Here, we will simulate it by activating server 1 after the first announcement.
	// It means that server 1 has to join the cluster by emitting a new round of election.
	//for index, server := range servers {
	//	if index != 0 {
	//		server.Activate()
	//	}
	//}
	//
	//for {
	//	if servers[1].Cluster.GetCoordinator() != nil {
	//		servers[0].PleaseIgnorePreviousAnnouncement()
	//		break
	//	}
	//}

	// 2.2 Simulate best-case
	// The best-case is: the node with the highest id start the election.
	// Here, we will simulate it by activating the server with the highest id after the first announcement.
	// It means that server highest_id has to join the cluster by emitting a new round of election.
	//for index, server := range servers {
	//	if index != len(servers)-1 {
	//		server.Activate()
	//	}
	//}
	//
	//for {
	//	if servers[0].Cluster.GetCoordinator() != nil {
	//		servers[len(servers)-1].Activate()
	//		break
	//	}
	//}

	// 3.a Simulate newly elected coordinator fails while announcing that it has won election to all nodes
	// This case could be well-handled by the heartbeat mechanism in this implementation.
	// To simulate this, activate the servers by calling PleaseFailWhileAnnounce()
	//for index, server := range servers {
	//	if index < len(servers)-1 {
	//		server.Activate()
	//	} else {
	//		server.PleaseFailWhileAnnounce()
	//	}
	//}

	// 3.b Simulate a node fails while election, the failed node is not the newly elected coordinator
	// This case could be well-handled by the election timeout mechanism in this implementation.
	// To simulate this, activate the servers by calling PleaseFailWhileElection()
	//for index, server := range servers {
	//	if index != 1 {
	//		server.Activate()
	//	} else {
	//		server.PleaseFailWhileElection()
	//	}
	//}

	// 4. Multiple GO routines start the election process simultaneously
	// This is embedded in the system, Multiple GO routines starts the election process simultaneously in this implementation

	// 5. An arbitrary node silently leaves the network
	// This case could be well-handled by the election timeout mechanism in this implementation.
	// It is explained, and almost the same as 3
	select {}
}
