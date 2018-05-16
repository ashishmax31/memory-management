package process

import (
	"bufio"
	"math/rand"
	"os"
	"time"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
)

const programLength = 64000

var cachedProgram []string

var randSrc = rand.NewSource(time.Now().UnixNano())
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Process struct {
	Pid             uint
	PgeTble         *pagetable.Pagetable
	ProgramLocation *os.File
}

func (p *Process) GenerateProgramVirtualAddress() uint16 {
	program_word := rand.Intn(programLength)
	addr := uint16(program_word)
	return addr
	// memoryReferences = append(memoryReferences, addr)

}

func (p *Process) BringPageFromDisk(pagenumber int) []string {
	startingAddress := pagenumber * 4096
	endingAddress := startingAddress + 4096
	if len(cachedProgram) == 0 {
		text := readProgramFromdisk(pagenumber)
		return text[startingAddress:endingAddress]
	} else {
		return cachedProgram[startingAddress:endingAddress]
	}
}

func readProgramFromdisk(pagenumber int) []string {
	file, err := os.Open("./program_1")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cachedProgram = append(cachedProgram, scanner.Text())
	}
	return cachedProgram
}
