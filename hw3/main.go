package main

import (
	"homework/hw3/ivy"
	"homework/hw3/logger"
	"homework/hw3/util"
)

const (
	numOfProcessor       = 2
	numOfPage            = 5
	readRequestInterval  = 2 // seconds
	writeRequestInterval = 5 // seconds
)

func main() {
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
	select {}
}
