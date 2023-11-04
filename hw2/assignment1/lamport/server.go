package lamport

type server struct {
	id      int
	channel chan Message
	queue   MsgPriorityQueue
}
