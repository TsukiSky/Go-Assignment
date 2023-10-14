package assignment1

type Message struct {
	senderId int
	clock    int
}

func (m *Message) SetClock(clock int) {
	m.clock = clock
}
