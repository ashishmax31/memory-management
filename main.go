package main

import (
	"fmt"
	"time"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/cpu"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/mainmemory"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/process"
)

func main() {

	// Initialize hardware
	var processor cpu.Cpu
	var proc process.Process
	var pgtable pagetable.Pagetable
	processor.CurrentExecutionContext = &proc
	proc.PgeTble = &pgtable
	proc.ProgramText = make(process.Entry)
	for i := 0; i < 16; i++ {
		proc.ProgramText[i] = make(process.Value)
	}
	// Get list of program generated addresses.
	addresses := proc.GenerateVirtualAddressess()
	//var pagefault []int

	for _, address := range addresses {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Fetching address %x \n", address)
		status, pageNum := processor.Fetch(address)
		if status == "restart" {
			// Page fault has occured need to restart the current instruction
			processor.Fetch(address)
			// pagefault = append(pagefault, pageNum)
		}
		fmt.Printf("\n\n")
	}

	// Overview
	fmt.Printf("Memory: \n")
	for _, item := range mainmemory.Memory {
		fmt.Printf("%v", item.InUse)
	}
	fmt.Printf("Pagetable: \n")
	for _, item := range proc.PgeTble {
		// if item.Referenced {
		fmt.Printf("%+v \n", item)
		// }
	}
}
