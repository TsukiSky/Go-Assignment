package ivy

import "homework/hw3/util"

type CentralManager struct {
	PageTable      *util.CMPageTable
	Processors     []*Processor
	messageChannel chan int
}
