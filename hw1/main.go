package main

import (
	"homework/hw1/assignment1"
)

// Consider some client-server architecture as follows.
// Several clients are registered to the server. Periodically, each client sends message to the server.
// Upon receiving a message, the server flips a coin and decides to either forward the message to all other registered clients (excluding the original sender of the message) or drops the message altogether.
//
// To solve this question, you will do the following:
// 1. Simulate the behaviour of both the server and the registered clients via GO routines.
// 2. Use Lamport’s logical clock to determine a total order of all the messages received at all the registered clients. Subsequently, present (i.e., print) this order for all registered clients to know the order in which the messages should be read.
// 3. Use Vector clock to redo the assignment. Implement the detection of causality violation and print any such detected causality violation.
//

const (
	numOfClients = 5
	timeInterval = 2 // seconds
)

//var wg sync.WaitGroup

func main() {
	server := assignment1.NewServer()
	server.Initialize()

	for i := 0; i < numOfClients; i++ {
		client := assignment1.NewClient(i, server)
		server.Register(client)
	}

	for _, client := range server.GetClients() {
		client.Activate(timeInterval)
	}

	for {
	}
}
