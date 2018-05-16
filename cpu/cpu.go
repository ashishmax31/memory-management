package cpu

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/mainmemory"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/process"
)

type Cpu struct {
	CurrentExecutionContext *process.Process
	currentAddress          uint16
}

func (c *Cpu) Fetch(addr uint16) (string, int) {
	virtualPageNumber, offset, err := translateVirtualAddressToPhysicalAddress(addr)
	fmt.Printf("---> page: %d \n", virtualPageNumber)
	if err != nil {
		log.Panic("Fatal hardware exception.. CPU burned off! lol \n")
	}
	pageframeNumber, exception := c.CurrentExecutionContext.PgeTble.PageTableLookUp(virtualPageNumber)
	if exception == pagetable.PageFaultException {
		fmt.Println("PageFault!!! Accessing page from disk.")
		c.HandlePageFault(virtualPageNumber, offset)
		// restart the operation on successful page fault handling.
		return "restart", virtualPageNumber
	}
	fmt.Printf("Data word from memory for the address: 0x%x is : %s \n", addr, mainmemory.Memory[pageframeNumber].Entries[offset])
	return "", virtualPageNumber
}

// MMU does this function
func translateVirtualAddressToPhysicalAddress(addr uint16) (int, int, error) {
	binary := fmt.Sprintf("%016b \n", addr)
	pageNumber, err := strconv.ParseInt(binary[0:4], 2, 64)
	if err != nil {
		return -1, -1, err

	}
	offset, err := strconv.ParseInt(binary[4:16], 2, 64)
	if err != nil {
		return -1, -1, err
	}
	return int(pageNumber), int(offset), nil

}

func (c *Cpu) HandlePageFault(virtualPageNumber int, offset int) {
	// time.Sleep(1 * time.Second)

	pageFrameWrittenTo := bringRequiredPageFromDiskToRam(virtualPageNumber, c.CurrentExecutionContext)
	for page, item := range c.CurrentExecutionContext.PgeTble {
		if (c.CurrentExecutionContext.PgeTble[page].PageFrameNumber == pageFrameWrittenTo) && (item.Present) {
			c.CurrentExecutionContext.PgeTble[page].Present = false
		}
	}
	c.CurrentExecutionContext.PgeTble[virtualPageNumber].PageFrameNumber = pageFrameWrittenTo
	c.CurrentExecutionContext.PgeTble[virtualPageNumber].Present = true
}

func bringRequiredPageFromDiskToRam(virtualPageNumber int, currentProcess *process.Process) int {
	pageFrameToWriteTo := mainmemory.GetPageFrame()

	data := currentProcess.BringPageFromDisk(virtualPageNumber) // Get the page from disk

	// Write the entire page from disk to the pageframe in main memory
	for offset := 0; offset < mainmemory.PageFrameSize; offset++ {
		mainmemory.Memory[pageFrameToWriteTo].Entries[offset] = data[offset]
	}
	mainmemory.Memory[pageFrameToWriteTo].InUse = true
	return pageFrameToWriteTo
}
