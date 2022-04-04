package dm

import "bytes"

type Session struct {
	dom     *DOM
	buffer  *Buffer
	offsets []int
}

func (d *DOM) Session(options ...Option) *Session {
	session := &Session{
		dom:     d,
		offsets: make([]int, len(d.attributes)),
	}

	session.apply(options)
	if session.buffer == nil {
		session.buffer = NewBuffer(d.initialBufferSize)
	}

	session.buffer.appendBytes(d.template)

	return session
}

func (s *Session) SetAttr(tag, attribute, oldValue, newValue []byte, fullMatch bool) {
	currOffset := 0
	currEnd := 0

	for i, attr := range s.dom.attributes {
		s.offsets[i] += currOffset
		if i == 0 {
			if !bytes.Equal(s.buffer.slice(attr.boundaries[0], 0, 0), attribute) || !bytes.Equal(attr.tag, tag) {
				continue
			}
		} else {
			if !bytes.Equal(s.buffer.slice(attr.boundaries[0], s.offsets[i-1], s.offsets[i-1]), attribute) || !bytes.Equal(attr.tag, tag) {
				continue
			}
		}

		currEnd = 0
		if i != 0 && s.offsets[i] != s.offsets[i-1] {
			currEnd = attr.valueEnd() + s.offsets[i] - s.offsets[i-1]
			if !s.match(fullMatch, s.buffer.buffer[attr.valueStart()+s.offsets[i-1]:currEnd], oldValue) {
				continue
			}
			currEnd = attr.valueEnd() - currEnd
		} else {
			if !s.match(fullMatch, s.buffer.buffer[attr.valueStart()+s.offsets[i]:attr.valueEnd()+s.offsets[i]], oldValue) {
				continue
			}
		}

		currOffset += s.buffer.insertBytes(attr.boundaries[1], currOffset+s.offsets[i], currEnd, newValue)
		s.offsets[i] += currOffset
	}
}

func (s *Session) Attribute(tag, attribute []byte) ([]byte, bool) {
	if len(s.dom.attributes) == 0 {
		return nil, false
	}

	var attrVal []byte
	var ok bool
	if attrVal, ok = s.checkByIndex(0, 0, s.offsets[0], tag, attribute); ok {
		return attrVal, ok
	}

	for i := 1; i < len(s.dom.attributes); i++ {
		if attrVal, ok = s.checkByIndex(i, s.offsets[i-1], s.offsets[i], attribute, tag); ok {
			return attrVal, ok
		}
	}

	return nil, false
}

func (s *Session) checkByIndex(index, startOffset, endOffset int, attribute []byte, tag []byte) ([]byte, bool) {
	if !bytes.Equal(s.buffer.slice(s.dom.attributes[index].boundaries[0], startOffset, startOffset), attribute) || !bytes.Equal(s.dom.attributes[index].tag, tag) {
		return nil, false
	}

	return s.buffer.buffer[s.dom.attributes[index].valueStart()+startOffset : s.dom.attributes[index].valueEnd()+endOffset], true
}

func (s *Session) match(fullMatch bool, value []byte, oldValue []byte) bool {
	return (!fullMatch && bytes.HasPrefix(value, oldValue)) || (fullMatch && bytes.Equal(value, oldValue))
}

func (s *Session) Bytes() []byte {
	return s.buffer.bytes()
}

func (s *Session) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case *Buffer:
			s.buffer = actual
		}
	}
}
