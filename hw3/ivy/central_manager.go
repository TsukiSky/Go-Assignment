package ivy

import (
	"homework/hw3/logger"
	"homework/hw3/util"
)

// CentralManager is the central manager of the ivy system
type CentralManager struct {
	PageTable         *util.CMPageTable
	Processors        []int
	MessageChannel    chan util.Message
	Connections       map[int]chan util.Message // map of processor id to their message channel
	WritingRequestMap map[int][]util.Message    // map of page id to write request
}

// NewCentralManager creates a new central manager
func NewCentralManager() *CentralManager {
	return &CentralManager{
		PageTable:         &util.CMPageTable{Records: map[int]*util.CMPageRecord{}},
		Processors:        []int{},
		MessageChannel:    make(chan util.Message),
		Connections:       map[int]chan util.Message{},
		WritingRequestMap: map[int][]util.Message{},
	}
}

// onReceiveReadReq handles the read request from a processor
func (c *CentralManager) onReceiveReadReq(message util.Message) {
	processorId := message.ProcessorId
	pageId := message.PageId
	c.PageTable.Records[pageId].CopySet = append(c.PageTable.Records[pageId].CopySet, processorId)
	ownerId := 0
	if c.PageTable.Records[pageId].HasOwner {
		ownerId = c.PageTable.Records[pageId].Owner
		// ask the owner to forward a copy of the page
		c.Connections[ownerId] <- util.Message{
			Type:        util.READ_FORWARD,
			PageId:      pageId,
			ProcessorId: processorId,
		}
	} else {
		ownerId = processorId
		c.forwardPage(pageId, processorId, false)
	}
}

// onReceiveWriteReq handles the write request from a processor
func (c *CentralManager) onReceiveWriteReq(message util.Message) {
	pageId := message.PageId
	processorId := message.ProcessorId
	if c.PageTable.Records[pageId].HasOwner {
		if c.PageTable.Records[pageId].OwnerIsWriting {
			// wait for the owner to finish writing
			c.WritingRequestMap[pageId] = append(c.WritingRequestMap[pageId], message)
		} else {
			c.handleWriteForward(pageId, processorId)
		}
	} else {
		// forward page
		c.forwardPage(pageId, processorId, true)
	}
}

// onReceiveWriteAck handles the write ack from a processor
func (c *CentralManager) onReceiveWriteAck(message util.Message) {
	pageId := message.PageId
	c.PageTable.Records[pageId].OwnerIsWriting = false
	if len(c.WritingRequestMap[pageId]) > 0 {
		// handle write forward
		request := c.WritingRequestMap[pageId][0]
		c.handleWriteForward(pageId, request.ProcessorId)
		c.WritingRequestMap[pageId] = c.WritingRequestMap[pageId][1:]
	}
}

// handleWriteForward handles the write forward action when it is decided that the page should be <Write Forward>
func (c *CentralManager) handleWriteForward(pageId int, processorId int) {
	// send out write forward
	c.Connections[c.PageTable.Records[pageId].Owner] <- util.Message{
		Type:        util.WRITE_FORWARD,
		PageId:      pageId,
		ProcessorId: processorId,
	}

	c.PageTable.Records[pageId].Owner = processorId
	c.PageTable.Records[pageId].OwnerIsWriting = true
	c.PageTable.Records[pageId].HasOwner = true
	// invalidate all caches
	go func() {
		// invalidate all copy set
		for _, processor := range c.PageTable.Records[pageId].CopySet {
			if processor != processorId {
				c.Connections[processor] <- util.Message{
					Type:   util.INVALIDATE,
					PageId: pageId,
				}
			}
		}
		c.PageTable.Records[pageId].ClearCopies() // clear copy set
	}()
}

// forwardPage forwards a page to a processor
func (c *CentralManager) forwardPage(pageId int, toProcessorId int, isWriteForward bool) {
	c.PageTable.Records[pageId].HasOwner = true
	c.PageTable.Records[pageId].Owner = toProcessorId
	c.PageTable.Records[pageId].OwnerIsWriting = isWriteForward

	page := c.PageTable.Records[pageId].Page.Clone()
	c.Connections[toProcessorId] <- util.Message{
		Type:           util.PAGE,
		Page:           page,
		PageId:         page.Id,
		IsWriteForward: isWriteForward,
	}
}

// Register registers a processor to the central manager
func (c *CentralManager) Register(processor *Processor) {
	c.Processors = append(c.Processors, processor.Id)
	c.Connections[processor.Id] = processor.MessageChannel
}

// Activate activates the central manager
func (c *CentralManager) Activate() {
	logger.Logger.Printf("[Central Manager] Central Manager activated\n")
	go c.listen()
}

// listen listens to the message channel
func (c *CentralManager) listen() {
	for {
		message := <-c.MessageChannel
		switch message.Type {
		case util.READ_REQUEST:
			logger.Logger.Printf("[Central Manager] Receive <<<Read Request>>> for Page %d from Processor %d\n", message.PageId, message.ProcessorId)
			c.onReceiveReadReq(message)
		case util.WRITE_REQUEST:
			logger.Logger.Printf("[Central Manager] Receive <<<Write Request>>> for Page %d from Processor %d\n", message.PageId, message.ProcessorId)
			c.onReceiveWriteReq(message)
		case util.WRITE_ACK:
			logger.Logger.Printf("[Central Manager] Receive <<<Write Ack>>> for Page %d from Processor %d\n", message.PageId, message.ProcessorId)
			c.onReceiveWriteAck(message)
		}
	}
}
