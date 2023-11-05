package main

import (
	"homework/hw2/assignment1/lamport"
)

func main() {
	cluster := lamport.NewCluster()
	for i := 0; i < 3; i++ {
		cluster.AddServer(lamport.NewServer(i))
	}
	cluster.Activate()
	select {}
}
