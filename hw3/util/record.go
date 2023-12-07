package util

type Page struct {
	Id int
}

func (p *Page) Clone() Page {
	return Page{p.Id}
}

type CMPageRecord struct {
	Page           *Page
	CopySet        []int
	Owner          int
	OwnerIsWriting bool
}

type CMPageTable struct {
	Records map[int]*CMPageRecord
}

func (c *CMPageTable) findPageById(id int) Page {
	return *c.Records[id].Page
}

type ProcessorPageRecord struct {
	Page   *Page
	Access Access
}

type ProcessorPageTable struct {
	Records map[int]*ProcessorPageRecord
}

func (p *ProcessorPageTable) FindPageById(id int) Page {
	return *p.Records[id].Page
}

func (p *ProcessorPageTable) InvalidatePage(id int) {
	delete(p.Records, id)
}

type Access int

const (
	READ Access = iota
	WRITE
)
