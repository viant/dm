package html

// Buffer hold the current VirtualDOM value
type Buffer struct {
	buffer []byte
	pos    int32
}

func (b *Buffer) appendBytes(template []byte) {
	b.growIfNeeded(int32(len(template)))
	b.pos += int32(copy(b.buffer[b.pos:], template))
}

func (b *Buffer) growIfNeeded(grow int32) {
	newSize := int(b.pos + grow)
	if newSize <= len(b.buffer)-1 {
		return
	}

	newBuffer := make([]byte, newSize)
	copy(newBuffer, b.buffer)
	b.buffer = newBuffer
}

func (b *Buffer) replaceBytes(span *span, offset int32, actualLenDiff int32, value []byte) int32 {
	increased := int32(len(value)-span.end+span.start) - actualLenDiff
	b.insertBetween(int32(span.start), int32(span.end), offset, actualLenDiff, value, increased)
	return increased
}

func (b *Buffer) insertBetween(start, end, shiftedSoFar int32, increasedSoFar int32, value []byte, lengthDiff int32) {
	b.growIfNeeded(lengthDiff)
	copy(b.buffer[end+shiftedSoFar+lengthDiff:], b.buffer[end+shiftedSoFar:])
	copy(b.buffer[start+shiftedSoFar-increasedSoFar:], value)
	b.pos += lengthDiff
}

func (b *Buffer) bytes() []byte {
	return b.buffer[:b.pos]
}

func (b *Buffer) slice(boundary *span, start, end int32) []byte {
	return b.buffer[int32(boundary.start)+start : int32(boundary.end)+end]
}

func (b *Buffer) reset() {
	b.pos = 0
}

func (b *Buffer) insertAfter(i, shiftedSoFar int32, attribute []byte) {
	b.insertBetween(i, i, shiftedSoFar, 0, attribute, int32(len(attribute)))
}

// NewBuffer creates new buffer of given size.
func NewBuffer(size int32) *Buffer {
	return &Buffer{
		buffer: make([]byte, size),
	}
}
