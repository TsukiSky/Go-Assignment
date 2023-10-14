package bully

type Data struct {
	users []user
}

type user struct {
	id   int
	name string
}

type Cluster struct {
	servers     []*Server
	coordinator *Server
}
