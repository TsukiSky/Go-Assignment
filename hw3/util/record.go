package util

type Page struct {
	id int
}

type CMPageRecord struct {
	page    *Page
	copySet []int
	owner   int
}

type CMPageTable struct {
	records map[int]*CMPageRecord
}

type ProcessorPageRecord struct {
	page   *Page
	access Access
}

type ProcessorPageTable struct {
	records map[int]*ProcessorPageRecord
}

type Access int

const (
	READ Access = iota
	WRITE
)
