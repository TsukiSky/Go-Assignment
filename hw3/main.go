package main

import (
	"homework/hw3/ivy"
	"homework/hw3/ivyfaulttolerant"
	"homework/hw3/logger"
	"homework/hw3/util"
	"sync"
	"time"
)

const (
	numOfProcessor         = 10
	numOfPage              = 50
	readRequestInterval    = 1 // seconds
	writeRequestInterval   = 3 // seconds
	isFaulty               = false
	syncInterval           = 2 // seconds
	primaryDownTime        = 5 // seconds
	primaryFailCount       = 1
	primaryRestartInterval = 8
	terminateReadNum       = 16 // read: write = 4:1
	terminateWriteNum      = 4
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(numOfProcessor)
	logger.InitPerformanceLog("hw3")
	startTime := time.Now()
	logger.StartTime = startTime

	if !isFaulty {
		// simulate the no-fault scenario
		logger.Init("hw3", "assignment_1.log", "assignment 1: ")
		centralManager := ivy.NewCentralManager() // create a central manager
		var processors []*ivy.Processor
		// create processors
		for i := 0; i < numOfProcessor; i++ {
			processor := ivy.NewProcessor(i, terminateReadNum, terminateWriteNum, centralManager.MessageChannel)
			processor.WaitGroup = &wg
			processors = append(processors, processor)
			centralManager.Register(processor)
		}
		// setup connections
		connections := map[int]chan util.Message{}
		for _, processor := range processors {
			connections[processor.Id] = processor.MessageChannel
		}
		for _, processor := range processors {
			processor.Connections = connections
		}
		// initialize page table
		centralManager.PageTable.Init(numOfPage)
		// activate central manager
		centralManager.Activate()
		// activate processors
		for _, processor := range processors {
			processor.Activate(numOfPage, readRequestInterval, writeRequestInterval)
		}
	} else {
		// simulate the faulty scenario
		logger.Init("hw3", "assignment_2.log", "assignment 2: ")
		primaryCentralManager := ivyfaulttolerant.NewPrimaryCentralManager(primaryDownTime) // create a primary central manager
		backupCentralManager := ivyfaulttolerant.NewBackupCentralManager()                  // create a backup central manager
		primaryCentralManager.SetBackupManagerChannel(backupCentralManager.MessageChannel)
		backupCentralManager.SetBackupManagerChannel(primaryCentralManager.MessageChannel)

		var processors []*ivyfaulttolerant.Processor
		// create processors
		for i := 0; i < numOfProcessor; i++ {
			processor := ivyfaulttolerant.NewProcessor(i, terminateReadNum, terminateWriteNum, primaryCentralManager.MessageChannel, backupCentralManager.MessageChannel)
			processors = append(processors, processor)
			processor.WaitGroup = &wg
			primaryCentralManager.Register(processor)
			backupCentralManager.Register(processor)
		}
		// setup connections
		connections := map[int]chan util.Message{}
		for _, processor := range processors {
			connections[processor.Id] = processor.MessageChannel
		}
		for _, processor := range processors {
			processor.Connections = connections
		}
		// initialize page table
		primaryCentralManager.PageTable.Init(numOfPage)
		backupCentralManager.PageTable.Init(numOfPage)

		// activate central managers
		primaryCentralManager.Activate(syncInterval)
		backupCentralManager.Activate(syncInterval)
		// activate processors
		for _, processor := range processors {
			processor.Activate(numOfPage, readRequestInterval, writeRequestInterval)
		}

		pDownTime := 0
		for i := 0; i < primaryFailCount; i++ {
			pDownTime += primaryDownTime
			time.Sleep(time.Duration(pDownTime+primaryRestartInterval) * time.Second)
			primaryCentralManager.RestartPrimary()
		}
	}
	wg.Wait()
	endTime := time.Now()
	executionTime := endTime.Sub(startTime)
	logger.PerformanceLogger.Println("Number of Processors:", numOfProcessor)
	logger.PerformanceLogger.Println("Number of Pages:", numOfPage)
	logger.PerformanceLogger.Println("Total Number of Read Requests:", terminateReadNum*numOfProcessor)
	logger.PerformanceLogger.Println("Total Number of Write Requests:", terminateWriteNum*numOfProcessor)
	logger.PerformanceLogger.Println("Reached termination condition! Total execution time:", executionTime)
}
