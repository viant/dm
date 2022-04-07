package dm

import (
	"bytes"
)

type (
	Session struct {
		dom          *DOM
		buffer       *Buffer
		offsets      []int
		innerChanged bool
	}
)

func (d *DOM) Session(options ...Option) *Session {
	session := &Session{
		dom:     d,
		offsets: make([]int, len(d.builder.attributes)),
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
		return -1
	}

	var attr *attr
	for i := offset + 1; i < len(s.dom.builder.attributes); i++ { // s.dom.attributes[0] is a sentinel
		attr = s.dom.builder.attributes[i]
		if len(selectors) > 0 && !bytes.Equal(s.attributeTag(i), selectors[0]) {
			continue
		}

		if len(selectors) > 1 && !bytes.Equal(s.buffer.slice(attr.boundaries[0], s.offsets[i-1], s.offsets[i-1]), selectors[1]) {
			continue
		}

		if len(selectors) > 2 && !s.attrValueMatches(i, selectors[2]) {
			continue
		}

		currOffset := s.buffer.insertBytes(attr.boundaries[1], s.offsets[i], s.offsetDiff(i), newValue)
		s.offsets[i] += currOffset
		for j := i + 1; j < len(s.dom.builder.attributes); j++ {
			s.offsets[j] += currOffset
		}
		return i
	}

	return -1
}

func (s *Session) Attribute(offset int, selectors ...[]byte) ([]byte, int, bool) {
	for i := offset + 1; i < len(s.dom.builder.attributes); i++ { // s.dom.attributes[0] is a sentinel
		if len(selectors) > 0 && !bytes.Equal(s.attributeTag(i), selectors[0]) {
			continue
		}

		if len(selectors) > 1 && !bytes.Equal(s.buffer.slice(s.dom.builder.attributes[i].boundaries[0], s.offsets[i-1], s.offsets[i-1]), selectors[1]) {
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
	return s.buffer.buffer[s.dom.builder.attributes[i].valueStart()+s.offsets[i-1] : s.dom.builder.attributes[i].valueEnd()+s.offsets[i]]
}

func (s *Session) InnerHTML(offset int, selectors ...[]byte) ([]byte, int) {
	tagIndex := s.findMatchingTag(offset, selectors)
	if tagIndex == -1 {
		return nil, tagIndex
	}

	return s.buffer.slice(s.tag(tagIndex).InnerHTML, s.tagOffset(tagIndex), s.tagOffset(tagIndex)), tagIndex
}

func (s *Session) findMatchingTag(offset int, selectors [][]byte) int {
	for i := offset + 1; i < len(s.dom.builder.tags); i++ {
		slice := s.buffer.slice(s.tag(i).TagName, s.tagOffset(i-1), s.tagOffset(i-1))
		if len(selectors) > 0 && !bytes.Equal(slice, selectors[0]) {
			continue
		}

		for j := s.dom.builder.tags[i-1].AttrEnd; j < s.dom.builder.tags[i].AttrEnd; j++ {
			slice := s.buffer.slice(s.dom.builder.attributes[j].boundaries[0], s.attrOffset(j-1), s.attrOffset(j-1))
			if len(selectors) > 1 && !bytes.Equal(slice, selectors[1]) {
				continue
			}

			if len(selectors) > 2 && !bytes.Equal(s.attributeValue(j), selectors[2]) {
				continue
			}

			return i
		}

		if len(selectors) == 1 {
			return i
		}
	}
	return -1
}

func (s *Session) Bytes() []byte {
	return s.buffer.bytes()
}

func (s *Session) tag(i int) *tag {
	return s.dom.builder.tags[i]
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
	return s.dom.builder.attributes[i]
}

func (s *Session) attributeTag(i int) []byte {
	return s.buffer.slice(s.dom.builder.tags[s.dom.builder.attributes[i].tag].TagName, s.attrOffset(i-1), s.attrOffset(i-1))
}

func (s *Session) attrOffset(i int) int {
	return s.offsets[i]
}

func (s *Session) tagOffset(i int) int {
	return s.dom.builder.tags.tagOffset(i, s.offsets)
}

func (s *Session) SetInnerHTML(offset int, value []byte, selectors ...[]byte) (int, error) {
	tagIndex := s.findMatchingTag(offset, selectors)
	if tagIndex == -1 {
		return tagIndex, nil
	}

	if err := s.rebuildDOM(tagIndex, value); err != nil {
		return 0, err
	}
	s.offsets = make([]int, len(s.dom.builder.attributes))
	return tagIndex, nil
}

func (s *Session) rebuildDOM(tagIndex int, newInnerHTML []byte) error {
	dom := s.dom
	if !s.innerChanged {
		dom = innerDom(s, dom, tagIndex, s.buffer.bytes())
		s.innerChanged = true
	} else {
		dom.template = s.buffer.bytes()
	}

	_, err := dom.rebuildTemplate(newInnerHTML)
	if err != nil {
		return err
	}
	s.dom = dom
	s.buffer.reset()
	s.buffer.appendBytes(s.dom.template)
	return nil
}

func innerDom(s *Session, dom *DOM, index int, template []byte) *DOM {
	return newInnerDOM(dom, &elementsBuilder{
		attributes: dom.builder.attributes.sliceTo(s.offsets, s.tag(index).AttrEnd),
		tags:       dom.builder.tags.sliceTo(s.offsets, index),
		tagIndexes: make([]int, 0),
		tagCounter: index - 1,
		offset:     s.tag(index).InnerHTML.Start + s.tagOffset(index),
		depth:      dom.builder.tags[index].Depth,
	}, template, dom)
}
