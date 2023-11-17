package util

type Message struct {
	SenderId    int
	MessageType MessageType
	ScalarClock int
}

type MessageType int

const (
	REQUEST         MessageType = iota // critical section request
	REPLY                              // critical section reply
	RELEASE                            // critical section release
	RESCIND                            // vote rescind (only in voting protocol)
	RESCIND_RELEASE                    // rescind release (only in voting protocol)
	VOTE                               // vote (only in voting protocol)
)

// SetClock sets the scalar clock of the message
func (m *Message) SetClock(clock int) {
	m.ScalarClock = clock
}

// IsLargerThan returns true if the message's scalar clock is larger than a given message
// If the scalar clocks are equal or concurrent, the message with the larger sender id is larger
func (m *Message) IsLargerThan(message Message) bool {
	if m.ScalarClock == message.ScalarClock {
		return m.SenderId > message.SenderId
	} else {
		return m.ScalarClock > message.ScalarClock
	}
}

// Equal returns true if the message equals to a given message
func (m *Message) Equal(message Message) bool {
	if (m.SenderId != message.SenderId) || (m.MessageType != message.MessageType || m.ScalarClock != message.ScalarClock) {
		return false
	}
	return true
}
