package vhtml

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

func (b *Buffer) insertBytes(span *Span, offset int, end int, value []byte) int {
	increased := len(value) - span.End + span.Start + end
	b.growIfNeeded(increased)

	copy(b.buffer[span.End+offset+increased:], b.buffer[span.End+offset:])
	copy(b.buffer[span.Start+offset+end:], value)
	b.pos += increased
	return increased
}

func (b *Buffer) bytes() []byte {
	return b.buffer[:b.pos]
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		buffer: make([]byte, size),
	}
}
