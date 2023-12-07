package ivy

import "homework/hw3/util"

type CentralManager struct {
	PageTable      *util.CMPageTable
	Processors     []*Processor
	messageChannel chan util.Message
	Connections    map[int]chan util.Message // map of processor id to their message channel
}

func (c *CentralManager) onReceiveReadReq(message util.Message) {
	processorId := message.ProcessorId
	pageId := message.PageId
	c.PageTable.Records[pageId].CopySet = append(c.PageTable.Records[pageId].CopySet, processorId)
	ownerId := c.PageTable.Records[pageId].Owner
	// ask the owner to forward a copy of the page
	c.Connections[ownerId] <- util.Message{
		Type:        util.READ_FORWARD,
		PageId:      pageId,
		ProcessorId: processorId,
	}
}

func (c *CentralManager) onReceiveWriteReq(message util.Message) {

}

func (c *CentralManager) onReceiveWriteAck(message util.Message) {
	pageId := message.PageId
	c.PageTable.Records[pageId].OwnerIsWriting = false
}

func (c *CentralManager) forwardPage(pageId int, toProcessorId int) {
	c.PageTable.Records[pageId].Owner = toProcessorId
	c.PageTable.Records[pageId].OwnerIsWriting = true
	page := c.PageTable.Records[pageId].Page.Clone()
	c.Connections[toProcessorId] <- util.Message{
		Type: util.PAGE,
		Page: page,
	}
}
