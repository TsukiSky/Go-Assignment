package main

import (
	"homework/hw2/assignment1/lamport"
	"homework/logger"
)

func main() {
	logger.Init("hw2", "assignment_1.log", "assignment 1:")
	cluster := lamport.NewCluster()
	for i := 0; i < 3; i++ {
		cluster.AddServer(lamport.NewServer(i))
	}
	cluster.Activate()
	select {}
}
