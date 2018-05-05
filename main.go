package main

import (
	"fmt"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/cpu"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/process"
)

func main() {
	var processor cpu.Cpu
	var process process.Process
	var pgtable pagetable.Pagetable
	processor.CurrentExecutionContext = &process
	process.PgeTble = &pgtable
	addresses := process.GenerateVirtualAddressess()
	for _, address := range addresses {
		fmt.Printf("Fetching address %x \n", address)
		status := processor.Fetch(address)
		if status == "restart" {
			println("restarting")
			processor.Fetch(address)
		}
	}
}
