package TUI



type RingBuffer[Data any] struct {
	Buffer []Data
	Size 	int
	
	Head	int
	Position	int
}

func NewBuffer[Data any](size int) *RingBuffer[Data] {
	return(&RingBuffer[Data]{
		Buffer:	make([]Data, size),
		Size:	size,
	})
}

func (ringBuff *RingBuffer[Data]) Add(logEntry Data)(evicted Data, wasEvicted bool){

	if ringBuff.Position == ringBuff.Size{
		wasEvicted = true
		evicted = ringBuff.Buffer[ringBuff.Head]
	}

	ringBuff.Buffer[ringBuff.Head] = logEntry // Appends the log entry to the head
	ringBuff.Head = (ringBuff.Head + 1) % ringBuff.Size // Resets the index back to zero when .head == .size

	if ringBuff.Position < ringBuff.Size{
		ringBuff.Position++
	}

	return

}

func  (ringBuff *RingBuffer[Data]) GetEntries() []Data{
	

	result := make([]Data, ringBuff.Position)

	start := (ringBuff.Head + ringBuff.Size - ringBuff.Position) % ringBuff.Size

	n := copy(result, ringBuff.Buffer[start:])
	if n < ringBuff.Position{
		copy(result[n:], ringBuff.Buffer[:ringBuff.Position-n]) // faster implementation than previous using copy (memmove)
	}

	return result

}

