package TUI

//"sync"

type RingBuffer[Data any] struct {
	Buffer []Data
	Size 	int
	//muteX 	sync.Mutex
	Head	int
	Position	int
}

func NewBuffer[Data any](size int) *RingBuffer[Data] {
	return(&RingBuffer[Data]{
		Buffer:	make([]Data, size),//TODO is the program trying to write to the buffer when I haven constructed it yet?
		Size:	size,
	})
}

func (ringBuff *RingBuffer[Data]) Add (logEntry Data){
	
	//ringBuff.muteX.Lock()
	//defer ringBuff.muteX.Unlock() // Locks the buffer to prevent race conditions

	ringBuff.Buffer[ringBuff.Head] = logEntry // Appends the log entry to the head
	ringBuff.Head = (ringBuff.Head + 1) % ringBuff.Size // Resets the index back to zero when .head == .size

	if ringBuff.Position < ringBuff.Size{
		ringBuff.Position++
	}


}

func  (ringBuff *RingBuffer[Data]) GetEntries() []Data{
	
	result:= make([]Data, 0, ringBuff.Size)

	for i := 0 ; i < ringBuff.Position ; i ++ {
		index := (ringBuff.Head + ringBuff.Size - ringBuff.Position + i) % ringBuff.Size
		result = append(result, ringBuff.Buffer[index])
	}

	return result

}

