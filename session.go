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
	matches := false

	var attr *attr
	for i := 1; i < len(s.dom.attributes); i++ { // s.dom.attributes[0] is a sentinel
		attr = s.dom.attributes[i]
		s.offsets[i] += currOffset
		if !bytes.Equal(s.buffer.slice(attr.boundaries[0], s.offsets[i-1], s.offsets[i-1]), attribute) || !bytes.Equal(attr.tag, tag) {
			continue
		}

		currEnd, matches = s.attrValueMatches(i, oldValue, fullMatch)
		if !matches {
			continue
		}

		currOffset += s.buffer.insertBytes(attr.boundaries[1], currOffset+s.offsets[i], currEnd, newValue)
		s.offsets[i] += currOffset
	}
}

func (s *Session) Attribute(tag, attribute []byte) ([]byte, bool) {
	for i := 1; i < len(s.dom.attributes); i++ { // s.dom.attributes[0] is a sentinel
		if !bytes.Equal(s.buffer.slice(s.dom.attributes[i].boundaries[0], s.offsets[i-1], s.offsets[i-1]), attribute) || !bytes.Equal(s.dom.attributes[i].tag, tag) {
			continue
		}

		return s.buffer.buffer[s.dom.attributes[i].valueStart()+s.offsets[i-1] : s.dom.attributes[i].valueEnd()+s.offsets[i]], true
	}
	return nil, false
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

func (s *Session) attrValueMatches(i int, oldValue []byte, fullMatch bool) (int, bool) {
	if s.offsets[i] != s.offsets[i-1] && s.match(fullMatch, s.buffer.slice(s.attrByIndex(i).boundaries[1], s.offsets[i-1], -(s.offsets[i-1]-s.offsets[i])), oldValue) {
		return s.offsets[i-1] - s.offsets[i], true
	}

	if s.match(fullMatch, s.buffer.slice(s.attrByIndex(i).boundaries[1], s.offsets[i], s.offsets[i]), oldValue) {
		return 0, true
	}

	return 0, false
}

func (s *Session) attrByIndex(i int) *attr {
	return s.dom.attributes[i]
}
