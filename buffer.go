package dm

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
	b.growIfNeeded(increased)

	copy(b.buffer[span.end+offset+increased:], b.buffer[span.end+offset:])
	copy(b.buffer[span.start+offset-actualLenDiff:], value)
	b.pos += increased
	return increased
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

//NewBuffer creates new buffer of given size.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		buffer: make([]byte, size),
	}
}
