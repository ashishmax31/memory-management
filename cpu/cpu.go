package cpu

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/mainmemory"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/process"
)

type Cpu struct {
	CurrentExecutionContext *process.Process
	currentAddress          uint16
}

func (c *Cpu) Fetch(addr uint16) (status string) {
	// fmt.Printf("Accessing virtual address: %x by process %d \n", addr, c.CurrentExecutionContext.Pid)
	virtualPageNumber, offset, err := translateVirtualAddressToPhysicalAddress(addr)
	fmt.Printf("Has virtual pageNumber: %d, offset: %d \n", virtualPageNumber, offset)
	if err != nil {
		log.Panic("Fatal hardware exception.. CPU burned off! lol \n")
	}
	pageframeNumber, exception := c.CurrentExecutionContext.PgeTble.PageTableLookUp(virtualPageNumber)
	if exception == pagetable.PageFaultException {
		err := c.HandlePageFault(virtualPageNumber)
		println("Page fault")
		if err != nil {
			log.Panic("Couldnt free a page frame from memory :( \n")
		}
		// restart the operation on successful page fault handling.
		return "restart"
	}
	fmt.Printf("Got page frame number: %d \n", pageframeNumber)
	println(mainmemory.Memory[pageframeNumber].Entries[offset])
	return ""
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

func (c *Cpu) HandlePageFault(virtualPageNumber int) (err error) {
	var pageForReplacement int
	for {
		n := rand.Intn(len(c.CurrentExecutionContext.PgeTble))
		if virtualPageNumber != n {
			pageForReplacement = n
			break
		}
	}
	// Mark the selected page for replacement as in use(Present = false) and move the corresponding page-frame to disk.
	pageFrameWrittenTo := bringRequiredPageFromDiskToRam(virtualPageNumber, c.CurrentExecutionContext, pageForReplacement)
	c.CurrentExecutionContext.PgeTble[virtualPageNumber].PageFrameNumber = pageFrameWrittenTo
	c.CurrentExecutionContext.PgeTble[virtualPageNumber].Present = true
	return nil
}

func bringRequiredPageFromDiskToRam(virtualPageNumber int, currentProcess *process.Process, pageForReplacement int) int {
	var pageFrameToWriteTo int
	if currentProcess.PgeTble[pageForReplacement].Present {
		pageFrameToWriteTo = currentProcess.PgeTble[pageForReplacement].PageFrameNumber
		currentProcess.PgeTble[pageForReplacement].Present = false
	} else {
		pageFrameToWriteTo = rand.Intn(len(mainmemory.Memory))
	}
	var entriesTobeWrittenToRam []process.Entry
	for _, program_data := range currentProcess.ProgramText {
		if program_data.Page == virtualPageNumber {
			entriesTobeWrittenToRam = append(entriesTobeWrittenToRam, program_data)
		}
	}
	for i, entry := range entriesTobeWrittenToRam {
		mainmemory.Memory[pageFrameToWriteTo].Entries[i] = entry.Data
	}
	mainmemory.Memory[pageFrameToWriteTo].Free = false
	return pageFrameToWriteTo

}
