package mainmemory

import "math/rand"

const mainMemorySize = 32 * 1024 // 32KB

const PageFrameSize = 4096

var Memory [(mainMemorySize / PageFrameSize)]pageFrame

type pageFrame struct {
	Entries [PageFrameSize]string
	InUse   bool
}

func GetPageFrame() int {
	for pageFrameNumber, pageframe := range Memory {
		if pageframe.InUse == false {
			return pageFrameNumber
		}
	}
	// If no free page frames are present in memory, give out a random pageframe
	return rand.Intn(len(Memory))
}
