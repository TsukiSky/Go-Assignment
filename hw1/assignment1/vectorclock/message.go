package vectorclock

type Message struct {
	senderId    int
	vectorClock []int
}

func (m *Message) SetClock(clock []int) {
	m.vectorClock = clock
}
