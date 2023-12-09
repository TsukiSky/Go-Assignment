package util

// Central Manager Fault Handling
// Methodology: Applying a Heartbeat containing synchronization information
// this information is synchronized between the central manager and the backup central manager

type Heartbeat struct {
	PageTable         CMPageTable
	WritingRequestMap map[int][]Message
}
