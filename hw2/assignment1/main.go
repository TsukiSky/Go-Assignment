package main

import (
	"homework/hw2/assignment1/sharedpriorityqueue"
	"homework/logger"
)

func main() {
	logger.Init("hw2", "assignment_1.log", "assignment 1:")
	cluster := sharedpriorityqueue.NewCluster()
	for i := 0; i < 3; i++ {
		cluster.AddServer(sharedpriorityqueue.NewServer(i))
	}
	cluster.Activate()
	select {}
}
