package TUI

//"sync"

type RingBuffer[Data any] struct {
	Buffer []Data
	Size 	int
	//muteX 	sync.Mutex
	Head	int
	Count	int
}

func NewBuffer[Data any](size int) *RingBuffer[Data] {
	return(&RingBuffer[Data]{
		Buffer:	make([]Data, size),
		Size:	size,
	})
}

func (ringBuff *RingBuffer[Data]) Add (logEntry Data)(int){
	
	//ringBuff.muteX.Lock()
	//defer ringBuff.muteX.Unlock() // Locks the buffer to prevent race conditions

	ringBuff.Buffer[ringBuff.Head] = logEntry // Appends the log entry to the head
	ringBuff.Head = (ringBuff.Head + 1) % ringBuff.Size // Resets the index back to zero when .head == .size

	if ringBuff.Count < ringBuff.Size{
		ringBuff.Count++
	}

	return 0 //

}

