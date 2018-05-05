package mainmemory

const mainMemorySize = 32 * 1024 // 32KB

const pageFrameSize = 4096

var Memory [(mainMemorySize / pageFrameSize)]pageFrame

type pageFrame struct {
	Entries [pageFrameSize - 1]string
	Free    bool
}

func init() {
	// In the begining, mark all the page frames as free.

	for _, pageFrame := range Memory {
		pageFrame.Free = true
	}
}
