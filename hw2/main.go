package main

import (
	"homework/hw2/logger"
	"homework/hw2/optimizedsharedpriorityqueue"
	"homework/hw2/sharedpriorityqueue"
	"homework/hw2/voting"
	"strconv"
	"sync"
	"time"
)

const (
	runningMode     = PERFORMANCE_COMPARING_MODE
	algorithm       = OPTIMIZED_SHARED_PRIORITY_QUEUE
	numOfServers    = 15
	numOfRequesters = 10
)

type RunningMode int

type Algorithm int

func (a Algorithm) String() string {
	switch a {
	case SHARED_PRIORITY_QUEUE:
		return "Shared Priority Queue"
	case OPTIMIZED_SHARED_PRIORITY_QUEUE:
		return "Optimized Shared Priority Queue (Ricart and Agrawala’s Optimization)"
	case VOTING:
		return "Voting Protocol"
	default:
		return "Unknown Algorithm"
	}
}

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

func main() {
	if runningMode == PERFORMANCE_COMPARING_MODE {
		// run all algorithms in performance comparing mode
		runInPerformanceMode(SHARED_PRIORITY_QUEUE)
		runInPerformanceMode(OPTIMIZED_SHARED_PRIORITY_QUEUE)
		runInPerformanceMode(VOTING)
		logger.PerformanceLogger.Println("########################################################################################")
	} else if runningMode == SINGLE_PERFORMANCE_MODE {
		// run only one algorithm in performance comparing mode
		runInPerformanceMode(algorithm)
		logger.PerformanceLogger.Println("########################################################################################")
	} else {
		// run only one algorithm in permanent mode
		runInPermanentMode(algorithm)
	}
}

func runInPerformanceMode(algorithm Algorithm) {
	logger.InitPerformanceLog("hw2")
	logger.PerformanceLogger.Println("########################################################################################")
	logger.PerformanceLogger.Println("[Algorithm]: " + algorithm.String())
	logger.PerformanceLogger.Println("[Number of Servers]: " + strconv.Itoa(numOfServers))
	logger.PerformanceLogger.Println("[Number of Requesters]: " + strconv.Itoa(numOfRequesters))

	var runningTime time.Duration
	var waitGroup sync.WaitGroup
	switch algorithm {
	case SHARED_PRIORITY_QUEUE:
		// run shared priority queue algorithm
		logger.Init("hw2", "shared_priority_queue.log", "shared priority queue:")
		cluster := sharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(sharedpriorityqueue.NewServer(i))
		}
		cluster.SetWaitGroup(&waitGroup)
		runningTime = cluster.ActivateInPerformanceComparingMode(numOfRequesters)

	case OPTIMIZED_SHARED_PRIORITY_QUEUE:
		// run optimized shared priority queue (Ricart and Agrawala's optimization) algorithm
		logger.Init("hw2", "optimized_shared_priority_queue.log", "optimized shared priority queue:")
		cluster := optimizedsharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(optimizedsharedpriorityqueue.NewServer(i))
		}
		cluster.SetWaitGroup(&waitGroup)
		runningTime = cluster.ActivateInPerformanceComparingMode(numOfRequesters)

	case VOTING:
		// run voting algorithm
		logger.Init("hw2", "voting_algorithm.log", "voting algorithm:")
		cluster := voting.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(voting.NewServer(i))
		}
		cluster.SetWaitGroup(&waitGroup)
		runningTime = cluster.ActivateInPerformanceComparingMode(numOfRequesters)
	}
	logger.PerformanceLogger.Println("[Time (s)]: " + runningTime.String())
}

func runInPermanentMode(algorithm Algorithm) {
	switch algorithm {
	case SHARED_PRIORITY_QUEUE:
		// run shared priority queue algorithm
		logger.Init("hw2", "shared_priority_queue.log", " shared priority queue:")
		cluster := sharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(sharedpriorityqueue.NewServer(i))
		}
		cluster.ActivateInPerformanceComparingMode(numOfRequesters)

	case OPTIMIZED_SHARED_PRIORITY_QUEUE:
		// run optimized shared priority queue (Ricart and Agrawala's optimization) algorithm
		logger.Init("hw2", "optimized_shared_priority_queue.log", "optimized shared priority queue:")
		cluster := optimizedsharedpriorityqueue.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(optimizedsharedpriorityqueue.NewServer(i))
		}
		cluster.Activate(numOfRequesters)

	case VOTING:
		// run voting algorithm
		logger.Init("hw2", "voting_algorithm.log", "voting algorithm:")
		cluster := voting.NewCluster()
		for i := 0; i < numOfServers; i++ {
			cluster.AddServer(voting.NewServer(i))
		}
		cluster.Activate(numOfRequesters)
	}
}
