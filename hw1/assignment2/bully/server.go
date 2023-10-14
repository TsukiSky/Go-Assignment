package bully

import "homework/hw1/assignment2/bully/util"

type ServerType int

const (
	COORDINATOR ServerType = iota
	WORKER
)

type Server struct {
	id            int
	serverType    ServerType
	coordinatorId int
	channel       chan util.Message
	cluster       Cluster
	data          Data
}
