package dm

import "bytes"

type (
	Session struct {
		dom     *DOM
		buffer  *Buffer
		offsets []int
	}

	SelectorOffset struct {
		Tag       int
		Attribute int
	}
)

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

func (s *Session) SetAttr(offset int, newValue []byte, selectors ...[]byte) int {
	if len(selectors) == 0 {
		return 0
	}
	currOffset := 0

	var attr *attr
	for i := offset + 1; i < len(s.dom.attributes); i++ { // s.dom.attributes[0] is a sentinel
		attr = s.dom.attributes[i]
		s.offsets[i] += currOffset
		if len(selectors) > 0 && !bytes.Equal(s.attributeTag(i), selectors[0]) {
			continue
		}

		if len(selectors) > 1 && !bytes.Equal(s.buffer.slice(attr.boundaries[0], s.offsets[i-1], s.offsets[i-1]), selectors[1]) {
			continue
		}

		if len(selectors) > 2 && !s.attrValueMatches(i, selectors[2]) {
			continue
		}

		currOffset += s.buffer.insertBytes(attr.boundaries[1], currOffset+s.offsets[i], s.offsetDiff(i), newValue)
		s.offsets[i] += currOffset
		return i
	}

	return 0
}

func (s *Session) Attribute(offset int, selectors ...[]byte) ([]byte, int, bool) {
	for i := offset + 1; i < len(s.dom.attributes); i++ { // s.dom.attributes[0] is a sentinel
		if len(selectors) > 0 && !bytes.Equal(s.attributeTag(i), selectors[0]) {
			continue
		}

		if len(selectors) > 1 && !bytes.Equal(s.buffer.slice(s.dom.attributes[i].boundaries[0], s.offsets[i-1], s.offsets[i-1]), selectors[1]) {
			continue
		}

		if len(selectors) > 2 && !bytes.Equal(s.attributeValue(i), selectors[2]) {
			continue
		}

		return s.attributeValue(i), i, true
	}
	return nil, 0, false
}

func (s *Session) attributeValue(i int) []byte {
	return s.buffer.buffer[s.dom.attributes[i].valueStart()+s.offsets[i-1] : s.dom.attributes[i].valueEnd()+s.offsets[i]]
}

func (s *Session) InnerHTML(offset int, selectors ...[]byte) ([]byte, int) {
	for i := offset + 1; i < len(s.dom.tags); i++ {
		if len(selectors) > 0 && !bytes.Equal(s.buffer.slice(s.tag(i).TagName, s.tagOffset(i-1), s.tagOffset(i-1)), selectors[0]) {
			continue
		}

		for j := s.dom.tags[i-1].AttrEnd; j < s.dom.tags[i].AttrEnd; j++ {
			slice := s.buffer.slice(s.dom.attributes[j].boundaries[0], s.attrOffset(j-1), s.attrOffset(j-1))
			if len(selectors) > 1 && !bytes.Equal(slice, selectors[1]) {
				continue
			}

			if len(selectors) > 2 && !bytes.Equal(s.attributeValue(j), selectors[2]) {
				continue
			}

			return s.buffer.slice(s.tag(i).InnerHTML, s.tagOffset(i), s.tagOffset(i)), i
		}

		if len(selectors) == 1 {
			return s.buffer.slice(s.tag(i).InnerHTML, s.tagOffset(i), s.tagOffset(i)), i
		}
	}
	return nil, 0
}

func (s *Session) Bytes() []byte {
	return s.buffer.bytes()
}

func (s *Session) tag(i int) *tag {
	return s.dom.tags[i]
}

func (s *Session) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case *Buffer:
			s.buffer = actual
		}
	}
}

func (s *Session) attrValueMatches(i int, oldValue []byte) bool {
	if s.offsets[i] != s.offsets[i-1] {
		return bytes.Equal(s.buffer.slice(s.attrByIndex(i).boundaries[1], s.offsets[i-1], -(s.offsets[i-1]-s.offsets[i])), oldValue)
	}

	return bytes.Equal(s.buffer.slice(s.attrByIndex(i).boundaries[1], s.offsets[i], s.offsets[i]), oldValue)
}

func (s *Session) offsetDiff(i int) int {
	return s.offsets[i-1] - s.offsets[i]
}

func (s *Session) attrByIndex(i int) *attr {
	return s.dom.attributes[i]
}

func (s *Session) attributeTag(i int) []byte {
	return s.buffer.slice(s.dom.tags[s.dom.attributes[i].tag].TagName, s.attrOffset(i-1), s.attrOffset(i-1))
}

func (s *Session) attrOffset(i int) int {
	return s.offsets[i]
}

func (s *Session) tagOffset(i int) int {
	i = s.dom.tags[i].AttrEnd - 1
	if i == -1 {
		return 0
	}

	return s.offsets[i]
}
