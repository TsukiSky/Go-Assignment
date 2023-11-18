package main

import (
	"homework/hw2/assignment1/optimizedsharedpriorityqueue"
	"homework/logger"
	"sync"
)

func main() {
	//var waitGroup sync.WaitGroup
	//logger.Init("hw2", "assignment_1.log", "assignment 1:")
	//cluster := sharedpriorityqueue.NewCluster()
	//for i := 0; i < 4; i++ {
	//	cluster.AddServer(sharedpriorityqueue.NewServer(i))
	//}
	//cluster.SetWaitGroup(&waitGroup)
	//cluster.ActivateInPerformanceComparingMode(2)

	var waitGroup sync.WaitGroup
	logger.Init("hw2", "assignment_2.log", "assignment 2:")
	cluster := optimizedsharedpriorityqueue.NewCluster()
	for i := 0; i < 4; i++ {
		cluster.AddServer(optimizedsharedpriorityqueue.NewServer(i))
	}
	cluster.SetWaitGroup(&waitGroup)
	//cluster.ActivateInPerformanceComparingMode(2)
	cluster.Activate(2)

	//var waitGroup sync.WaitGroup
	//logger.Init("hw2", "assignment_3.log", "assignment 3:")
	//cluster := voting.NewCluster()
	//for i := 0; i < 10; i++ {
	//	cluster.AddServer(voting.NewServer(i))
	//}
	//cluster.SetWaitGroup(&waitGroup)
	//cluster.ActivateInPerformanceComparingMode(6)
}
