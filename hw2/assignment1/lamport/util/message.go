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
// If the vector clocks are equal or concurrent, the message with the larger sender id is larger
func (m *Message) LargerThan(message Message) bool {
	hasLarger := false
	hasLess := false
	for i, clock := range m.VectorClock {
		if clock < message.VectorClock[i] {
			hasLess = true
		} else if clock > message.VectorClock[i] {
			hasLarger = true
		}
	}
	if hasLarger && !hasLess {
		return true
	} else if !hasLarger && hasLess {
		return false
	} else {
		return m.SenderId > message.SenderId
	}
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
