package util

type Message struct {
	SenderId    int
	MessageType MessageType
	VectorClock []int
}

type MessageType int

const (
	REQUEST MessageType = iota // critical section request
	REPLY                      // critical section reply
	RELEASE                    // critical section release
)

func (m *Message) SetClock(clock []int) {
	m.VectorClock = clock
}

// LargerThan returns true if the message's vector clock is larger than a given message
func (m *Message) LargerThan(message Message) bool {
	for i, clock := range m.VectorClock {
		if clock < message.VectorClock[i] {
			return false
		}
	}
	return true
}

// Equal returns true if the message equals to a given message
func (m *Message) Equal(message Message) bool {
	if (m.SenderId != message.SenderId) || (m.MessageType != message.MessageType || len(m.VectorClock) != len(message.VectorClock)) {
		return false
	}
	for i, clock := range m.VectorClock {
		if clock != message.VectorClock[i] {
			return false
		}
	}
	return true
}
