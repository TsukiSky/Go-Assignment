package util

import "homework/hw1/assignment2/bully"

type GenericContent struct {
	SenderId   int
	ReceiverId int
}

type MessageType int

const (
	SYN_REQ MessageType = iota
	SYN_REP
	ELE_REQ
	ELE_REP
	ANC_REQ
	//ANC_REP
)

type Message interface {
	GetContent() GenericContent
	GetMessageType() MessageType
}

// SynReqMessage implementation
type SynReqMessage struct {
	messageType MessageType
	content     GenericContent
	data        bully.Data
}

func (m *SynReqMessage) GetContent() GenericContent {
	return m.content
}

func (m *SynReqMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewSynRequestMsg(sender int, receiver int) *SynReqMessage {
	return &SynReqMessage{
		messageType: SYN_REQ,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
	}
}

// SynRepMessage implementation
type SynRepMessage struct {
	messageType MessageType
	content     GenericContent
	synSuccess  bool
}

func (m *SynRepMessage) GetContent() GenericContent {
	return m.content
}

func (m *SynRepMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewSynReplyMsg(sender int, receiver int, synSuccess bool) *SynRepMessage {
	return &SynRepMessage{
		messageType: SYN_REP,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
		synSuccess: synSuccess,
	}
}

// ElectReqMessage implementation
type ElectReqMessage struct {
	messageType MessageType
	content     GenericContent
}

func (m *ElectReqMessage) GetContent() GenericContent {
	return m.content
}

func (m *ElectReqMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewElectReqMsg(sender int, receiver int) *ElectReqMessage {
	return &ElectReqMessage{
		messageType: ELE_REQ,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
	}
}

// ElectRepMessage implementation
type ElectRepMessage struct {
	messageType MessageType
	content     GenericContent
	agree       bool
}

func (m *ElectRepMessage) GetContent() GenericContent {
	return m.content
}

func (m *ElectRepMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewElectRepMsg(sender int, receiver int, agree bool) *ElectRepMessage {
	return &ElectRepMessage{
		messageType: ELE_REP,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
		agree: agree,
	}
}

func (m *ElectRepMessage) IsAgree() bool {
	return m.agree
}

// AncMessage implementation
type AncMessage struct {
	messageType MessageType
	content     GenericContent
}

func (m *AncMessage) GetContent() GenericContent {
	return m.content
}

func (m *AncMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewAncMsg(sender int, receiver int) *AncMessage {
	return &AncMessage{
		messageType: ANC_REQ,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
	}
}

//
//// AncRepMessage implementation
//type AncRepMessage struct {
//	messageType MessageType
//	content     GenericContent
//	agree       bool
//}
//
//func (m *AncRepMessage) GetContent() GenericContent {
//	return m.content
//}
//
//func (m *AncRepMessage) GetMessageType() MessageType {
//	return m.messageType
//}
//
//func NewAncRepMsg(sender int, receiver int, agree bool) *AncRepMessage {
//	return &AncRepMessage{
//		messageType: ANC_REP,
//		content: GenericContent{
//			SenderId:   sender,
//			ReceiverId: receiver,
//		},
//		agree: agree,
//	}
//}
