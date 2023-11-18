package main

import (
	"homework/hw2/assignment1/voting"
	"homework/logger"
)

func main() {
	//logger.Init("hw2", "assignment_1.log", "assignment 1:")
	//cluster := sharedpriorityqueue.NewCluster()
	//for i := 0; i < 3; i++ {
	//	cluster.AddServer(sharedpriorityqueue.NewServer(i))
	//}
	//cluster.Activate()
	//select {}

	//logger.Init("hw2", "assignment_2.log", "assignment 2:")
	//cluster := optimizedsharedpriorityqueue.NewCluster()
	//for i := 0; i < 3; i++ {
	//	cluster.AddServer(optimizedsharedpriorityqueue.NewServer(i))
	//}
	//cluster.Activate()
	//select {}

	logger.Init("hw2", "assignment_3.log", "assignment 3:")
	cluster := voting.NewCluster()
	for i := 0; i < 3; i++ {
		cluster.AddServer(voting.NewServer(i))
	}
	cluster.Activate()
	select {}
}
