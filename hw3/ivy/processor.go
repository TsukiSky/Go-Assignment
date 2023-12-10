package ivy

import (
	"homework/hw3/logger"
	"homework/hw3/util"
	"math/rand"
	"sync"
	"time"
)

// Processor is a processor in the ivy system
type Processor struct {
	Id                    int
	PageTable             *util.ProcessorPageTable
	MessageChannel        chan util.Message
	CentralManagerChannel chan util.Message
	Connections           map[int]chan util.Message // map of processor id to their message channel
	terminateReadNum      int
	terminateWriteNum     int
	terminated            bool
	WaitGroup             *sync.WaitGroup
}

// NewProcessor creates a new processor
func NewProcessor(id int, terminateReadNum int, terminateWriteNum int, cmChannel chan util.Message) *Processor {
	return &Processor{
		Id:                    id,
		PageTable:             &util.ProcessorPageTable{Records: map[int]*util.ProcessorPageRecord{}},
		MessageChannel:        make(chan util.Message),
		CentralManagerChannel: cmChannel,
		Connections:           map[int]chan util.Message{},
		terminateReadNum:      terminateReadNum,
		terminateWriteNum:     terminateWriteNum,
		terminated:            false,
	}
}

// Activate activates the processor
func (p *Processor) Activate(maxPageNumber int, readRequestInterval int, writeRequestInterval int) {
	go p.listen()
	go p.ReadWithInterval(readRequestInterval, maxPageNumber)   // send read request every readRequestInterval seconds
	go p.WriteWithInterval(writeRequestInterval, maxPageNumber) // send write request every writeRequestInterval seconds
}

// Read reads (request stage) a page
func (p *Processor) Read(pageId int) {
	if p.PageTable.FindPageById(pageId) != nil {
		logger.Logger.Printf("[Processor %d] -- Read Page %d from local page table\n", p.Id, pageId)
	} else {
		p.SendReadReq(pageId)
	}
}

// Write writes (request stage) a page
func (p *Processor) Write(pageId int) {
	p.SendWriteReq(pageId)
}

// ReadWithInterval sends read requests with a given interval
func (p *Processor) ReadWithInterval(interval int, maxPageNumber int) {
	requestTimer := time.NewTimer(time.Duration(interval) * time.Second)
	for {
		select {
		case <-requestTimer.C:
			if p.terminateReadNum > 0 {
				p.terminateReadNum--
				pageId := rand.Intn(maxPageNumber)
				p.Read(pageId)
				requestTimer.Reset(time.Duration(interval) * time.Second)
			} else if p.terminateWriteNum <= 0 && !p.terminated {
				p.terminated = true
				p.WaitGroup.Done()
			}
		}
	}
}

// WriteWithInterval sends write requests with a given interval
func (p *Processor) WriteWithInterval(interval int, maxPageNumber int) {
	requestTimer := time.NewTimer(time.Duration(interval) * time.Second)
	for {
		select {
		case <-requestTimer.C:
			if p.terminateWriteNum > 0 {
				p.terminateWriteNum--
				pageId := rand.Intn(maxPageNumber)
				p.Write(pageId)
				requestTimer.Reset(time.Duration(interval) * time.Second)
			} else if p.terminateReadNum <= 0 && !p.terminated {
				p.terminated = true
				p.WaitGroup.Done()
			}
		}
	}
}

// onReceiveReadForward is called when the processor receives a <READ-FORWARD, toId, pageId> from the central manager
func (p *Processor) onReceiveReadForward(message util.Message) {
	toId := message.ProcessorId
	pageId := message.PageId
	if p.PageTable.FindPageById(pageId) != nil {
		page := *p.PageTable.FindPageById(pageId)
		page = page.Clone()
		go func() {
			p.ForwardPage(page, toId, false)
		}()
	}
}

// onReceiveWriteForward is called when the processor receives a <WRITE-FORWARD, toId, pageId> from the central manager
func (p *Processor) onReceiveWriteForward(message util.Message) {
	toId := message.ProcessorId
	pageId := message.PageId
	if p.PageTable.FindPageById(pageId) != nil {
		page := *p.PageTable.FindPageById(pageId)
		page = page.Clone()
		go func() {
			p.ForwardPageAndInvalidate(page, toId)
		}()
	}
}

// ForwardPageAndInvalidate forwards a page to another processor and invalidate the page in the local page table
func (p *Processor) ForwardPageAndInvalidate(page util.Page, toProcessorId int) {
	p.ForwardPage(page, toProcessorId, true)
	p.PageTable.InvalidatePage(page.Id)
}

// ForwardPage forwards a page to another processor
func (p *Processor) ForwardPage(page util.Page, toProcessorId int, isWriteForward bool) {
	p.Connections[toProcessorId] <- util.Message{
		Type:           util.PAGE,
		PageId:         page.Id,
		ProcessorId:    toProcessorId,
		Page:           page,
		IsWriteForward: isWriteForward,
	}
}

// writing simulates the writing time
func (p *Processor) writing(pageId int) {
	// This function is used to simulate the writing time
	writingTimer := time.NewTimer(time.Duration(2) * time.Second)
	defer writingTimer.Stop()
	<-writingTimer.C
	if p.PageTable.FindPageById(pageId) != nil {
		p.PageTable.Records[pageId].Access = util.READ
		p.CentralManagerChannel <- util.Message{
			Type:        util.WRITE_ACK,
			PageId:      pageId,
			ProcessorId: p.Id,
		}
	}
}

// onReceivePage is called when the processor receives a <PAGE, page> from another processor
func (p *Processor) onReceivePage(message util.Message) {
	page := message.Page
	if message.IsWriteForward {
		p.PageTable.Records[page.Id] = &util.ProcessorPageRecord{
			Page:   &page,
			Access: util.WRITE,
		}
		go func() {
			p.writing(page.Id)
		}()
	} else {
		p.PageTable.Records[page.Id] = &util.ProcessorPageRecord{
			Page:   &page,
			Access: util.READ,
		}
	}
}

// onReceiveInvalidate is called when the processor receives a <INVALIDATE, pageId> from another processor
func (p *Processor) onReceiveInvalidate(message util.Message) {
	pageId := message.PageId
	p.PageTable.InvalidatePage(pageId)
}

// SendReadReq sends a <READ-REQUEST, pageId> to the central manager
func (p *Processor) SendReadReq(pageId int) {
	logger.Logger.Printf("[Processor %d] -- Send <<<Read Request>>> for Page %d\n", p.Id, pageId)
	p.CentralManagerChannel <- util.Message{
		Type:        util.READ_REQUEST,
		PageId:      pageId,
		ProcessorId: p.Id,
	}
}

// SendWriteReq sends a <WRITE-REQUEST, pageId> to the central manager
func (p *Processor) SendWriteReq(pageId int) {
	logger.Logger.Printf("[Processor %d] -- Send <<<Write Request>>> for Page %d\n", p.Id, pageId)
	p.CentralManagerChannel <- util.Message{
		Type:        util.WRITE_REQUEST,
		PageId:      pageId,
		ProcessorId: p.Id,
	}
}

// listen listens to the message channel
func (p *Processor) listen() {
	for {
		select {
		case message := <-p.MessageChannel:
			switch message.Type {
			case util.READ_FORWARD:
				logger.Logger.Printf("[Processor %d] -- Receive <<<Read Forward>>> for Page %d to Processor %d\n", p.Id, message.PageId, message.ProcessorId)
				p.onReceiveReadForward(message)
			case util.WRITE_FORWARD:
				if message.ProcessorId != p.Id {
					// only forward if the message is not sent by itself
					logger.Logger.Printf("[Processor %d] -- Receive <<<Write Forward>>> for Page %d to Processor %d\n", p.Id, message.PageId, message.ProcessorId)
					p.onReceiveWriteForward(message)
				} else {
					// if the write message is sent by itself, it can start writing on the page
					go func() {
						p.writing(message.PageId)
					}()
				}
			case util.PAGE:
				logger.Logger.Printf("[Processor %d] -- Receive <<<Page>>> Page %d\n", p.Id, message.PageId)
				p.onReceivePage(message)
			case util.INVALIDATE:
				logger.Logger.Printf("[Processor %d] -- Receive <<<Invalidate>>> Page %d\n", p.Id, message.PageId)
				p.onReceiveInvalidate(message)
			}
		}
	}
}
