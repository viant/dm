package xml

//Buffer hold the current VirtualDOM value
type Buffer struct {
	buffer []byte
	pos    int
}

func (b *Buffer) appendBytes(template []byte) int {
	b.growIfNeeded(len(template))
	copied := copy(b.buffer[b.pos:], template)
	b.pos += copied

	return copied
}

func (b *Buffer) growIfNeeded(grow int) {
	if b.pos+grow <= len(b.buffer)-1 {
		return
	}

	newBuffer := make([]byte, b.pos+grow, (b.pos+grow)*2)
	copy(newBuffer, b.buffer)
	b.buffer = newBuffer
}

func (b *Buffer) String() string {
	return string(b.buffer[:b.pos])
}

//NewBuffer creates new buffer of given size.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		buffer: make([]byte, size),
	}
}
