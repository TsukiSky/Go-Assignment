package util

import "container/heap"

/**
 * I don't know why Golang designs its "interface" in such an ugly way.
 * Sometimes development in Golang has no choice but to adapt to the weird way of its interface.
 * An example:
 * msgPriorityQueue implements the interface "heap.Interface"
 * However, it does not really "inherit" the interface. The developer has to call heap.Init(), heap.Push(), heap.Pop() to maintain a heapified structure.
 * To avoid exposing heap methods to the outside world, I wrapped up the priority queue in a struct MsgPriorityQueue.
 */

// MsgPriorityQueue wraps up the priority queue for messages
type MsgPriorityQueue struct {
	queue msgPriorityQueue
}

func NewMsgPriorityQueue() *MsgPriorityQueue {
	queue := make(msgPriorityQueue, 0)
	heap.Init(&queue)
	return &MsgPriorityQueue{
		queue: queue,
	}
}

func (q *MsgPriorityQueue) Len() int {
	return q.queue.Len()
}

func (q *MsgPriorityQueue) Push(msg Message) {
	heap.Push(&q.queue, msg)
}

func (q *MsgPriorityQueue) Pop() Message {
	return heap.Pop(&q.queue).(Message)
}

// Peek returns the message with the highest priority
func (q *MsgPriorityQueue) Peek() Message {
	return q.queue[0]
}

// msgPriorityQueue is a priority queue for messages
type msgPriorityQueue []Message

// Len returns the length of the priority queue
func (q *msgPriorityQueue) Len() int {
	return len(*q)
}

// Less defines the priority of messages in the priority queue
func (q *msgPriorityQueue) Less(i, j int) bool {
	return !(*q)[i].LargerThan((*q)[j])
}

// Swap defines the swapping rule of two messages in the priority queue
func (q *msgPriorityQueue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

// Push adds a message to the priority queue
func (q *msgPriorityQueue) Push(x any) {
	*q = append(*q, x.(Message))
}

// Pop removes the message with the highest priority from the priority queue
func (q *msgPriorityQueue) Pop() any {
	n := len(*q)
	x := (*q)[n-1]
	*q = (*q)[:n-1]
	return x
}
