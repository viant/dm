package vhtml

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
	//currOffset := 0
	//currEnd := 0
	//
	//for i, attr := range s.dom.attributes {
	//	s.offsets[i] += currOffset
	//	if !bytes.Equal(attr.Key, attribute) || !bytes.Equal(attr.Tag, tag) {
	//		continue
	//	}
	//
	//	currEnd = 0
	//
	//	if i != 0 && s.offsets[i] != s.offsets[i-1] {
	//		currEnd = attr.ValueEnd() + s.offsets[i] - s.offsets[i-1]
	//		if !s.match(fullMatch, s.buffer.buffer[attr.ValueStart()+s.offsets[i-1]:currEnd], oldValue) {
	//			continue
	//		}
	//		currEnd = attr.ValueEnd() - currEnd
	//	} else {
	//		if !s.match(fullMatch, s.buffer.buffer[attr.ValueStart()+s.offsets[i]:attr.ValueEnd()+s.offsets[i]], oldValue) {
	//			continue
	//		}
	//	}
	//
	//	currOffset += s.buffer.insertBytes(attr.Boundaries[1], currOffset+s.offsets[i], currEnd, newValue)
	//	s.offsets[i] += currOffset
	//}
	currOffset := 0
	currEnd := 0

	for i, attr := range s.dom.attributes {
		s.offsets[i] += currOffset
		a := s.buffer.buffer[attr.KeyStart()+s.offsets[i] : attr.KeyEnd()+s.offsets[i]]
		if !bytes.Equal(a, attribute) || !bytes.Equal(attr.Tag, tag) {
			continue
		}

		currEnd = 0

		if i != 0 && s.offsets[i] != s.offsets[i-1] {
			currEnd = attr.ValueEnd() + s.offsets[i] - s.offsets[i-1]
			if !s.match(fullMatch, s.buffer.buffer[attr.ValueStart()+s.offsets[i-1]:currEnd], oldValue) {
				continue
			}
			currEnd = attr.ValueEnd() - currEnd
		} else {
			if !s.match(fullMatch, s.buffer.buffer[attr.ValueStart()+s.offsets[i]:attr.ValueEnd()+s.offsets[i]], oldValue) {
				continue
			}
		}

		currOffset += s.buffer.insertBytes(attr.Boundaries[1], currOffset+s.offsets[i], currEnd, newValue)
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
	if !bytes.Equal(s.buffer.buffer[s.dom.attributes[index].KeyStart()+startOffset:s.dom.attributes[index].KeyEnd()+startOffset], attribute) || !bytes.Equal(s.dom.attributes[index].Tag, tag) {
		return nil, false
	}

	return s.buffer.buffer[s.dom.attributes[index].ValueStart()+startOffset : s.dom.attributes[index].ValueEnd()+endOffset], true
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
