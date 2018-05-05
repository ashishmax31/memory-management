package pagetable

const pageSize = 4096     //bytes(4KB)
const osAddressSpace = 16 //bits
const PageFaultException = "PageFault"

// Pagetable... has 2^OsAddressSpace/pagesize number of entries.
type Pagetable [16]Page
type Page struct {
	PageFrameNumber int
	Present         bool
	Referenced      bool
	Dirty           bool
}

func (pgtble *Pagetable) PageTableLookUp(virtualPageNumber int) (pageFrameNumber int, exception string) {
	page := pgtble[virtualPageNumber]
	if page.Present {
		page.Referenced = true
		return page.PageFrameNumber, ""
	} else {
		return -1, PageFaultException
	}
}
