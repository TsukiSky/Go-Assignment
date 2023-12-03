package util

type Message struct {
	Type   MessageType
	PageId int
}

type MessageType int

const (
	READ_REQUEST MessageType = iota // read
	READ_FORWARD
	CONTENT
	WRITE_REQUEST // write
	INVALIDATE
	INVALIDATE_ACK
	WRITE_FORWARD
	WRITE_ACK
)
