package util

// Page is a page in the ivy system
type Page struct {
	Id int
}

// Clone clones a page
func (p *Page) Clone() Page {
	return Page{p.Id}
}

// CMPageRecord is a page record in the central manager
type CMPageRecord struct {
	Page           *Page
	CopySet        []int
	Owner          int
	HasOwner       bool
	OwnerIsWriting bool
}

// ClearCopies clears the copy set of a page
func (c *CMPageRecord) ClearCopies() {
	c.CopySet = c.CopySet[:0]
}

// Clone clones a page record
func (c *CMPageRecord) Clone() CMPageRecord {
	page := c.Page.Clone()
	return CMPageRecord{
		Page:           &page,
		CopySet:        c.CopySet[:],
		Owner:          c.Owner,
		HasOwner:       c.HasOwner,
		OwnerIsWriting: c.OwnerIsWriting,
	}
}

// CMPageTable is the page table in the central manager
type CMPageTable struct {
	Records map[int]*CMPageRecord
}

// FindPageById finds a page by its id
func (c *CMPageTable) findPageById(id int) Page {
	return *c.Records[id].Page
}

// Init initializes the page table in the central manager
func (c *CMPageTable) Init(numOfPage int) {
	for i := 0; i < numOfPage; i++ {
		c.Records[i] = &CMPageRecord{
			Page:           &Page{Id: i},
			CopySet:        []int{},
			Owner:          0,
			HasOwner:       false,
			OwnerIsWriting: false,
		}
	}
}

// Clone clones a page table
func (c *CMPageTable) Clone() CMPageTable {
	records := map[int]*CMPageRecord{}
	for k, v := range c.Records {
		record := v.Clone()
		records[k] = &record
	}
	return CMPageTable{Records: records}
}

// ProcessorPageRecord is a page record in a processor
type ProcessorPageRecord struct {
	Page   *Page
	Access Access
}

// ProcessorPageTable is the page table in a processor
type ProcessorPageTable struct {
	Records map[int]*ProcessorPageRecord
}

// FindPageById finds a page by its id
func (p *ProcessorPageTable) FindPageById(id int) *Page {
	if _, ok := p.Records[id]; !ok {
		return nil
	}
	return p.Records[id].Page
}

// InvalidatePage invalidates a page in the page table
func (p *ProcessorPageTable) InvalidatePage(id int) {
	delete(p.Records, id)
}

// Access is the access type owned by a processor for a page
type Access int

const (
	READ Access = iota
	WRITE
)
