package ivyfaulttolerant

import (
	"homework/hw3/logger"
	"homework/hw3/util"
	"time"
)

// CentralManager is the central manager of the ivy system
type CentralManager struct {
	PageTable            *util.CMPageTable
	Processors           []int
	MessageChannel       chan util.Message
	BackupMessageChannel chan util.Message
	Connections          map[int]chan util.Message // map of processor id to their message channel
	WritingRequestMap    map[int][]util.Message    // map of page id to write request
	IsPrimary            bool
	Downtime             int
	Type                 string
}

// NewPrimaryCentralManager creates a new central manager
func NewPrimaryCentralManager(primaryDownTime int) *CentralManager {
	return &CentralManager{
		PageTable:         &util.CMPageTable{Records: map[int]*util.CMPageRecord{}},
		Processors:        []int{},
		MessageChannel:    make(chan util.Message),
		Connections:       map[int]chan util.Message{},
		WritingRequestMap: map[int][]util.Message{},
		IsPrimary:         true,
		Downtime:          primaryDownTime,
		Type:              "Primary",
	}
}

// NewBackupCentralManager creates a new backup central manager
func NewBackupCentralManager() *CentralManager {
	return &CentralManager{
		PageTable:         &util.CMPageTable{Records: map[int]*util.CMPageRecord{}},
		Processors:        []int{},
		MessageChannel:    make(chan util.Message),
		Connections:       map[int]chan util.Message{},
		WritingRequestMap: map[int][]util.Message{},
		IsPrimary:         false,
		Type:              "Backup",
	}
}

func (c *CentralManager) SetBackupManagerChannel(backupManagerChannel chan util.Message) {
	c.BackupMessageChannel = backupManagerChannel
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
		go func() {
			c.Connections[ownerId] <- util.Message{
				Type:        util.READ_FORWARD,
				PageId:      pageId,
				ProcessorId: processorId,
			}
		}()
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
func (c *CentralManager) Activate(heartbeatInterval int) {
	if c.IsPrimary {
		logger.Logger.Printf("[Primary Central Manager] Central Manager activated\n")
		downTimerForHeartbeat := time.NewTimer(time.Duration(c.Downtime) * time.Second)
		downTimerForListening := time.NewTimer(time.Duration(c.Downtime) * time.Second)
		go c.listenAsPrimary(downTimerForListening)
		go c.SendHeartbeatWithInterval(heartbeatInterval, downTimerForHeartbeat)
	} else {
		logger.Logger.Printf("[Backup Central Manager] Backup Central Manager activated\n")
		go c.listenAsBackup(heartbeatInterval)
	}

}

// SendHeartbeat sends a heartbeat to the backup central manager
func (c *CentralManager) SendHeartbeat() {
	writingRequestMap := map[int][]util.Message{}
	// send heartbeat to the backup central manager
	for key, requests := range c.WritingRequestMap {
		writingRequestMap[key] = requests[:]
	}
	go func() {
		c.BackupMessageChannel <- util.Message{
			Type: util.HEARTBEAT,
			Heartbeat: util.Heartbeat{
				PageTable:         c.PageTable.Clone(),
				WritingRequestMap: writingRequestMap,
			},
		}
	}()
}

// SendHeartbeatWithInterval sends a heartbeat to the backup central manager with a certain interval
func (c *CentralManager) SendHeartbeatWithInterval(interval int, downTimer *time.Timer) {
	heartbeatSendingTimer := time.NewTimer(time.Duration(interval) * time.Second)
	for {
		select {
		case <-heartbeatSendingTimer.C:
			logger.Logger.Printf("[Primary Central Manager] Send <<<Heartbeat>>> to Backup Central Manager\n")
			c.SendHeartbeat()
			heartbeatSendingTimer.Reset(time.Duration(interval) * time.Second)
		case <-downTimer.C:
			// primary central manager is down
			return
		}

	}
}

// onReceiveHeartbeat handles the heartbeat from the primary central manager
func (c *CentralManager) onReceiveHeartbeat(message util.Message) {
	// update page table
	c.PageTable = &message.Heartbeat.PageTable
	// update writing request map
	c.WritingRequestMap = message.Heartbeat.WritingRequestMap
}

// listenAsPrimary listens to the message channel as a primary central manager
func (c *CentralManager) listenAsPrimary(downTimer *time.Timer) {
	if downTimer != nil {
		for {
			select {
			case <-downTimer.C:
				// primary central manager is down
				logger.Logger.Printf("[Primary Central Manager] Primary is DOWN\n")
				//return
			case message := <-c.MessageChannel:
				switch message.Type {
				case util.READ_REQUEST:
					logger.Logger.Printf("[%s Central Manager] Receive <<<Read Request>>> for Page %d from Processor %d\n", c.Type, message.PageId, message.ProcessorId)
					c.onReceiveReadReq(message)
				case util.WRITE_REQUEST:
					logger.Logger.Printf("[%s Central Manager] Receive <<<Write Request>>> for Page %d from Processor %d\n", c.Type, message.PageId, message.ProcessorId)
					c.onReceiveWriteReq(message)
				case util.WRITE_ACK:
					logger.Logger.Printf("[%s Central Manager] Receive <<<Write Ack>>> for Page %d from Processor %d\n", c.Type, message.PageId, message.ProcessorId)
					c.onReceiveWriteAck(message)
				}
			}
		}
	} else {
		for {
			message := <-c.MessageChannel
			switch message.Type {
			case util.READ_REQUEST:
				logger.Logger.Printf("[%s Central Manager] Receive <<<Read Request>>> for Page %d from Processor %d\n", c.Type, message.PageId, message.ProcessorId)
				c.onReceiveReadReq(message)
			case util.WRITE_REQUEST:
				logger.Logger.Printf("[%s Central Manager] Receive <<<Write Request>>> for Page %d from Processor %d\n", c.Type, message.PageId, message.ProcessorId)
				c.onReceiveWriteReq(message)
			case util.WRITE_ACK:
				logger.Logger.Printf("[%s Central Manager] Receive <<<Write Ack>>> for Page %d from Processor %d\n", c.Type, message.PageId, message.ProcessorId)
				c.onReceiveWriteAck(message)
			}
		}
	}
}

// BroadcastPrimaryCMDown broadcasts the primary central manager is down
func (c *CentralManager) BroadcastPrimaryCMDown() {
	for _, processorId := range c.Processors {
		c.Connections[processorId] <- util.Message{
			Type: util.PRIMARY_DOWN,
		}
	}
}

func (c *CentralManager) PromoteToPrimary() {
	c.IsPrimary = true
	go c.listenAsPrimary(nil)
	go c.BroadcastPrimaryCMDown()
}

// listenAsBackup listens to the message channel as a backup central manager
func (c *CentralManager) listenAsBackup(heartbeatInterval int) {
	heartbeatCheckingTimer := time.NewTimer(time.Duration(heartbeatInterval) * time.Second * 2)
	for {
		select {
		case <-heartbeatCheckingTimer.C:
			// no heartbeat received, the primary central manager is down
			// promote the backup central manager to primary central manager
			logger.Logger.Printf("[Backup Central Manager] Primary Central Manager is down, promote Backup Central Manager to Primary Central Manager\n")
			c.PromoteToPrimary()
			return
		case message := <-c.MessageChannel:
			switch message.Type {
			case util.HEARTBEAT:
				heartbeatCheckingTimer.Reset(time.Duration(heartbeatInterval) * time.Second * 2)
				logger.Logger.Printf("[Backup Central Manager] Receive <<<Heartbeat>>> from Primary Central Manager\n")
				c.onReceiveHeartbeat(message)
			}
		}
	}
}
