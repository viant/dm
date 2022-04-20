package xml

//Buffer hold the current XML value
type Buffer struct {
	buffer []byte
	pos    int
}

func (b *Buffer) appendBytes(bytes []byte) int {
	b.growIfNeeded(len(bytes))
	copied := copy(b.buffer[b.pos:], bytes)
	b.pos += copied

	return copied
}

func (b *Buffer) appendByte(value byte) int {
	b.growIfNeeded(1)
	b.buffer[b.pos] = value
	b.pos++
	return 1
}

func (b *Buffer) growIfNeeded(grow int) {
	if b.pos+grow <= len(b.buffer)-1 {
		return
	}

	newBuffer := make([]byte, b.pos+grow, (b.pos+grow)*2)
	copy(newBuffer, b.buffer)
	b.buffer = newBuffer
}

//String returns String representation of the buffer
func (b *Buffer) String() string {
	return string(b.buffer[:b.pos])
}

//NewBuffer creates new buffer of given size.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		buffer: make([]byte, size),
	}
}
