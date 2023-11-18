package main

import (
	"homework/hw2/assignment1/optimizedsharedpriorityqueue"
	"homework/hw2/assignment1/sharedpriorityqueue"
	"homework/hw2/assignment1/voting"
	"homework/logger"
	"sync"
)

type RunningMode int

type Algorithm int

const (
	PERFORMANCE_COMPARING_MODE RunningMode = iota
	SINGLE_PERFORMANCE_MODE
	SINGLE_RUNNING_MODE
)

const (
	SHARED_PRIORITY_QUEUE Algorithm = iota
	OPTIMIZED_SHARED_PRIORITY_QUEUE
	VOTING
)

const (
	runningMode     = PERFORMANCE_COMPARING_MODE
	algorithm       = OPTIMIZED_SHARED_PRIORITY_QUEUE
	numOfServers    = 15
	numOfRequesters = 10
)

func main() {
	if runningMode == PERFORMANCE_COMPARING_MODE {
		// run all algorithms in performance comparing mode
		runInPerformanceMode(SHARED_PRIORITY_QUEUE)
		runInPerformanceMode(OPTIMIZED_SHARED_PRIORITY_QUEUE)
		runInPerformanceMode(VOTING)
	} else if runningMode == SINGLE_PERFORMANCE_MODE {
		// run only one algorithm in performance comparing mode
		runInPerformanceMode(algorithm)
	} else {
		// run only one algorithm in permanent mode
		runInPermanentMode(algorithm)
	}
}

func runInPerformanceMode(algorithm Algorithm) {
	var waitGroup sync.WaitGroup
	switch algorithm {
	case SHARED_PRIORITY_QUEUE:
		// run shared priority queue algorithm
		logger.Init("hw2", "assignment_1.log", "assignment 1:")
		cluster := sharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(sharedpriorityqueue.NewServer(i))
		}
		cluster.SetWaitGroup(&waitGroup)
		cluster.ActivateInPerformanceComparingMode(numOfRequesters)

	case OPTIMIZED_SHARED_PRIORITY_QUEUE:
		// run optimized shared priority queue (Ricart and Agrawala's optimization) algorithm
		logger.Init("hw2", "assignment_2.log", "assignment 2:")
		cluster := optimizedsharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(optimizedsharedpriorityqueue.NewServer(i))
		}
		cluster.SetWaitGroup(&waitGroup)
		cluster.ActivateInPerformanceComparingMode(numOfRequesters)

	case VOTING:
		// run voting algorithm
		logger.Init("hw2", "assignment_3.log", "assignment 3:")
		cluster := voting.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(voting.NewServer(i))
		}
		cluster.SetWaitGroup(&waitGroup)
		cluster.ActivateInPerformanceComparingMode(numOfRequesters)
	}
}

func runInPermanentMode(algorithm Algorithm) {
	switch algorithm {
	case SHARED_PRIORITY_QUEUE:
		// run shared priority queue algorithm
		logger.Init("hw2", "assignment_1.log", "assignment 1:")
		cluster := sharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(sharedpriorityqueue.NewServer(i))
		}
		cluster.ActivateInPerformanceComparingMode(numOfRequesters)

	case OPTIMIZED_SHARED_PRIORITY_QUEUE:
		// run optimized shared priority queue (Ricart and Agrawala's optimization) algorithm
		logger.Init("hw2", "assignment_2.log", "assignment 2:")
		cluster := optimizedsharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(optimizedsharedpriorityqueue.NewServer(i))
		}
		cluster.Activate(numOfRequesters)

	case VOTING:
		// run voting algorithm
		logger.Init("hw2", "assignment_3.log", "assignment 3:")
		cluster := voting.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(voting.NewServer(i))
		}
		cluster.Activate(numOfRequesters)
	}
}
