package bully

type GenericContent struct {
	SenderId   int
	ReceiverId int
}

type MessageType int

const (
	SYN_REQ MessageType = iota // synchronization request
	ELE_REQ                    // election start request
	ELE_REP                    // election reply
	ANC_REQ                    // announcement
	HBT_REQ                    // heartbeat request
	HBT_REP                    // heartbeat reply
)

type Message interface {
	GetContent() GenericContent
	GetMessageType() MessageType
}

// SynReqMessage implementation
type SynReqMessage struct {
	messageType MessageType
	content     GenericContent
	data        Data
}

func (m *SynReqMessage) GetContent() GenericContent {
	return m.content
}

func (m *SynReqMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewSynRequestMsg(sender int, receiver int, data Data) *SynReqMessage {
	return &SynReqMessage{
		messageType: SYN_REQ,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
		data: data,
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
}

func (m *ElectRepMessage) GetContent() GenericContent {
	return m.content
}

func (m *ElectRepMessage) GetMessageType() MessageType {
	return m.messageType
}

func NewElectRepMsg(sender int, receiver int) *ElectRepMessage {
	return &ElectRepMessage{
		messageType: ELE_REP,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
	}
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

type Heartbeat interface {
	GetBeater() int
	GetAsker() int
}

// HeartbeatReq implementation
type HeartbeatReq struct {
	messageType MessageType
	content     GenericContent
}

func (m *HeartbeatReq) GetContent() GenericContent {
	return m.content
}

func (m *HeartbeatReq) GetMessageType() MessageType {
	return m.messageType
}

func (m *HeartbeatReq) GetBeater() int {
	return m.content.ReceiverId
}

func (m *HeartbeatReq) GetAsker() int {
	return m.content.SenderId
}

func NewHeartbeatReq(sender int, receiver int) *HeartbeatReq {
	return &HeartbeatReq{
		messageType: HBT_REQ,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
	}
}

// HeartbeatRep implementation
type HeartbeatRep struct {
	messageType MessageType
	content     GenericContent
}

func (m *HeartbeatRep) GetContent() GenericContent {
	return m.content
}

func (m *HeartbeatRep) GetMessageType() MessageType {
	return m.messageType
}

func (m *HeartbeatRep) GetBeater() int {
	return m.content.SenderId
}

func (m *HeartbeatRep) GetAsker() int {
	return m.content.ReceiverId
}

func NewHeartbeatRep(sender int, receiver int) *HeartbeatRep {
	return &HeartbeatRep{
		messageType: HBT_REP,
		content: GenericContent{
			SenderId:   sender,
			ReceiverId: receiver,
		},
	}
}
