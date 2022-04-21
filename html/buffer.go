package html

//Buffer hold the current VirtualDOM value
type Buffer struct {
	buffer []byte
	pos    int
}

func (b *Buffer) appendBytes(template []byte) {
	b.growIfNeeded(len(template))
	b.pos += copy(b.buffer[b.pos:], template)
}

func (b *Buffer) growIfNeeded(grow int) {
	if b.pos+grow <= len(b.buffer)-1 {
		return
	}

	newBuffer := make([]byte, b.pos+grow, (b.pos+grow)*2)
	copy(newBuffer, b.buffer)
	b.buffer = newBuffer
}

func (b *Buffer) replaceBytes(span *span, offset int, actualLenDiff int, value []byte) int {
	increased := len(value) - span.end + span.start - actualLenDiff
	b.insertBetween(span.start, span.end, offset, actualLenDiff, value, increased)
	return increased
}

func (b *Buffer) insertBetween(start, end, shiftedSoFar int, increasedSoFar int, value []byte, lengthDiff int) {
	b.growIfNeeded(lengthDiff)
	copy(b.buffer[end+shiftedSoFar+lengthDiff:], b.buffer[end+shiftedSoFar:])
	copy(b.buffer[start+shiftedSoFar-increasedSoFar:], value)
	b.pos += lengthDiff
}

func (b *Buffer) bytes() []byte {
	return b.buffer[:b.pos]
}

func (b *Buffer) slice(boundary *span, start, end int) []byte {
	return b.buffer[boundary.start+start : boundary.end+end]
}

func (b *Buffer) reset() {
	b.pos = 0
}

func (b *Buffer) insertAfter(i, shiftedSoFar int, attribute []byte) {
	b.insertBetween(i, i, shiftedSoFar, 0, attribute, len(attribute))
}

//NewBuffer creates new buffer of given size.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		buffer: make([]byte, size),
	}
}
