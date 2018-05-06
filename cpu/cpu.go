package cpu

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/mainmemory"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/process"
)

type Cpu struct {
	CurrentExecutionContext *process.Process
	currentAddress          uint16
}

func (c *Cpu) Fetch(addr uint16) (string, int) {
	// fmt.Printf("Accessing virtual address: %x by process %d \n", addr, c.CurrentExecutionContext.Pid)
	virtualPageNumber, offset, err := translateVirtualAddressToPhysicalAddress(addr)
	// fmt.Printf("Has virtual pageNumber: %d, offset: %d \n", virtualPageNumber, offset)
	if err != nil {
		log.Panic("Fatal hardware exception.. CPU burned off! lol \n")
	}

	pageframeNumber, exception := c.CurrentExecutionContext.PgeTble.PageTableLookUp(virtualPageNumber)
	if exception == pagetable.PageFaultException {
		fmt.Println("PageFault!!! Accessing page from disk.")
		err := c.HandlePageFault(virtualPageNumber, offset)
		// println("Page fault")
		if err != nil {
			log.Panic("Couldnt free a page frame from memory :( \n")
		}
		// restart the operation on successful page fault handling.
		return "restart", virtualPageNumber
	}
	// fmt.Printf("Got page frame number: %d \n", pageframeNumber)
	fmt.Printf("Data word from memory for the address: %x is : %s \n", addr, mainmemory.Memory[pageframeNumber].Entries[offset])
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

func (c *Cpu) HandlePageFault(virtualPageNumber int, offset int) (err error) {
	time.Sleep(3 * time.Second)
	pageFrameWrittenTo := bringRequiredPageFromDiskToRam(virtualPageNumber, c.CurrentExecutionContext)
	for page, item := range c.CurrentExecutionContext.PgeTble {
		if (c.CurrentExecutionContext.PgeTble[page].PageFrameNumber == pageFrameWrittenTo) && (item.Present) {
			println("marked false")
			c.CurrentExecutionContext.PgeTble[page].Present = false
		}
	}
	c.CurrentExecutionContext.PgeTble[virtualPageNumber].PageFrameNumber = pageFrameWrittenTo
	c.CurrentExecutionContext.PgeTble[virtualPageNumber].Present = true
	return nil
}

func bringRequiredPageFromDiskToRam(virtualPageNumber int, currentProcess *process.Process) int {
	pageFrameToWriteTo := mainmemory.GetPageFrame()
	for page, programData := range currentProcess.ProgramText {
		if page == virtualPageNumber {
			for offset, data := range programData {
				mainmemory.Memory[pageFrameToWriteTo].Entries[offset] = data
			}
		}
	}
	mainmemory.Memory[pageFrameToWriteTo].InUse = true
	return pageFrameToWriteTo
}
