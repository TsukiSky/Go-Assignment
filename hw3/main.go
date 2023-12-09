package main

import (
	"homework/hw3/ivy"
	"homework/hw3/ivyfaulttolerant"
	"homework/hw3/logger"
	"homework/hw3/util"
)

const (
	numOfProcessor       = 2
	numOfPage            = 5
	readRequestInterval  = 2 // seconds
	writeRequestInterval = 5 // seconds
	isFaulty             = true
	syncInterval         = 4 // seconds
	primaryDownTime      = 5 // seconds
)

func main() {
	if !isFaulty {
		// simulate the no-fault scenario
		logger.Init("hw3", "assignment_1.log", "assignment 1: ")
		centralManager := ivy.NewCentralManager() // create a central manager
		var processors []*ivy.Processor
		// create processors
		for i := 0; i < numOfProcessor; i++ {
			processor := ivy.NewProcessor(i, centralManager.MessageChannel)
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

		var processors []*ivyfaulttolerant.Processor
		// create processors
		for i := 0; i < numOfProcessor; i++ {
			processor := ivyfaulttolerant.NewProcessor(i, primaryCentralManager.MessageChannel, backupCentralManager.MessageChannel)
			processors = append(processors, processor)
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
	}
	select {}
}
