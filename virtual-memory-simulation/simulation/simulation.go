package simulation

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/cpu"
	"github.com/ashishmax31/memory-management/virtual-memory-simulation/process"
)

var Processor cpu.Cpu
var clock *time.Ticker

var processes []process.Process
var counter uint64

func StartSimulation() {
	// In this simulation just assume pagetable of the respective process is brought in to the main memory.
	InitializeProcesses()
	// Run the scheduler for the first time.
	runScheduler()
	for {
		<-clock.C
		counter++
		if (counter % 100) == 0 {
			// Run the scheduler after 100 clock cycles
			runScheduler()
		}
		memoryAccessRequest := Processor.CurrentExecutionContext.GenerateProgramVirtualAddress()
		fmt.Printf("Fetching address 0x%x for process %d ", memoryAccessRequest, Processor.CurrentExecutionContext.Pid)
		status, _ := Processor.Fetch(memoryAccessRequest)
		if status == "restart" {
			// Page fault has occured need to restart the current instruction
			Processor.Fetch(memoryAccessRequest)
		}

	}
}

func init() {
	// Set clock signal interval to 10ms
	clock = time.NewTicker(10 * time.Millisecond)
}

func runScheduler() {
	pidToRun := rand.Intn(len(processes))
	Processor.CurrentExecutionContext = &processes[pidToRun]
}

func InitializeProcesses() {
	for i := 0; i < 2; i++ {
		programLocationOnDisk, err := os.Open(fmt.Sprintf("program_%d", i+1))
		if err != nil {
			panic("Cant find the program source file on the disk!")
		}
		processes = append(processes, process.Process{Pid: uint(i), PgeTble: &pagetable.Pagetable{}, ProgramLocation: programLocationOnDisk})
	}

}
