package ivy

import (
	"homework/hw3/util"
)

type Processor struct {
	Id                           int
	PageTable                    *util.ProcessorPageTable
	MessageChannel               chan util.Message
	CentralManagerMessageChannel chan util.Message
	Connections                  map[int]chan util.Message // map of processor id to their message channel
}

// onReceiveReadForward is called when the processor receives a <READ-FORWARD, toId, pageId> from the central manager
func (p *Processor) onReceiveReadForward(message util.Message) {
	toId := message.ProcessorId
	pageId := message.PageId
	page := p.PageTable.FindPageById(pageId)
	page = page.Clone()
	p.ForwardPage(page, toId)
}

// onReceiveWriteForward is called when the processor receives a <WRITE-FORWARD, toId, pageId> from the central manager
func (p *Processor) onReceiveWriteForward(message util.Message) {
	toId := message.ProcessorId
	pageId := message.PageId
	page := p.PageTable.FindPageById(pageId)
	page = page.Clone()
	p.ForwardPage(page, toId)
	p.PageTable.InvalidatePage(pageId)
}

// ForwardPage forwards a page to another processor
func (p *Processor) ForwardPage(page util.Page, toProcessorId int) {
	p.Connections[toProcessorId] <- util.Message{
		Type:        util.PAGE,
		PageId:      page.Id,
		ProcessorId: toProcessorId,
		Page:        page,
	}
}

// onReceivePage is called when the processor receives a <PAGE, page> from another processor
func (p *Processor) onReceivePage(message util.Message) {
	page := message.Page
	p.PageTable.Records[page.Id] = &util.ProcessorPageRecord{
		Page:   &page,
		Access: util.READ, // TODO: should this be READ or WRITE?
	}
}

// SendReadReq sends a <READ-REQUEST, pageId> to the central manager
func (p *Processor) SendReadReq(pageId int) {
	p.CentralManagerMessageChannel <- util.Message{
		Type:        util.READ_REQUEST,
		PageId:      pageId,
		ProcessorId: p.Id,
	}
}

// SendWriteReq sends a <WRITE-REQUEST, pageId> to the central manager
func (p *Processor) SendWriteReq(pageId int) {
	p.CentralManagerMessageChannel <- util.Message{
		Type:        util.WRITE_REQUEST,
		PageId:      pageId,
		ProcessorId: p.Id,
	}
}
