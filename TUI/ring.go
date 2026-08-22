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
	
	/* TODO more efficient implementation..?
	
	result := make([]Data, ringBuff.Position)
	
	// Calculate where the oldest element is currently located
	start := (ringBuff.Head + ringBuff.Size - ringBuff.Position) % ringBuff.Size

	// Copy in up to two chunks to handle the ring wrap-around
	n := copy(result, ringBuff.Buffer[start:])
	if n < ringBuff.Position {
		copy(result[n:], ringBuff.Buffer[:ringBuff.Position-n])
	}

	return result
	*/


	result:= make([]Data, 0, ringBuff.Size)

	for i := 0 ; i < ringBuff.Position ; i ++ {
		index := (ringBuff.Head + ringBuff.Size - ringBuff.Position + i) % ringBuff.Size
		result = append(result, ringBuff.Buffer[index] )
	}

	return result

}

