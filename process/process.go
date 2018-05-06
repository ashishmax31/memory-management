package process

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ashishmax31/memory-management/virtual-memory-simulation/pagetable"
)

const programLength = 100000

var randSrc = rand.NewSource(time.Now().UnixNano())
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Process struct {
	Pid         uint
	PgeTble     *pagetable.Pagetable
	ProgramText Entry
}

type Entry map[int]Value

type Value map[int]string

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (p *Process) GenerateVirtualAddressess() (memoryReferences []uint16) {
	rnd := rand.New(randSrc)
	for i := 0; i < programLength; i++ {
		addr := uint16(rnd.Uint64())
		memoryReferences = append(memoryReferences, addr)
		addressString := fmt.Sprintf("%016b", addr)
		pageNumber, _ := strconv.ParseInt(addressString[0:4], 2, 64)
		offset, _ := strconv.ParseInt(addressString[4:16], 2, 64)
		p.ProgramText[int(pageNumber)][int(offset)] = randSeq(8)
	}
	return memoryReferences
}
