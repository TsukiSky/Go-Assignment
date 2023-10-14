package assignment2

// Use Bully algorithm to implement a working version of replica synchronization.
// You may assume that
// 1. each replica maintains some data structure (that may diverge for arbitrary reasons),
// 		which are periodically synchronized with the coordinator.
// 2. The coordinator initiates the synchronization process by sending message to all other machines.
// 3. Upon receiving the message from the coordinator, each machine updates its local version of the
// 		data structure with the coordinator’s version.
// 4. The coordinator, being an arbitrary machine in the network, is subject to fault. Thus,
//		4.1 a new coordinator is chosen by the Bully algorithm.
//		4.2	You can assume a fixed timeout to simulate the behaviour of detecting a fault
//			The objective is to have a consensus across all machines (simulated by GO routines)
//			in terms of the newly elected coordinator.

func main() {

}
