package dm

import (
	"bytes"
)

type (
	//Session modifies the DOM
	Session struct {
		dom               *DOM
		buffer            *Buffer
		attributesOffsets []int
		innerIncreased    []int
		skipped           []bool
		removedTags       map[int]bool
	}
)

//Session creates new Session
func (d *DOM) Session(options ...Option) *Session {
	session := &Session{
		dom:               d,
		attributesOffsets: make([]int, len(d.builder.attributes)),
		innerIncreased:    make([]int, len(d.builder.tags)),
		skipped:           make([]bool, len(d.builder.tags)),
	}

	session.apply(options)
	if session.buffer == nil {
		session.buffer = NewBuffer(d.initialBufferSize)
	}

	session.buffer.appendBytes(d.template)

	return session
}

//SetAttribute set attribute value matches the selectors and returns attribute index, -1 if not found
func (s *Session) SetAttribute(offset int, newValue []byte, selectors ...[]byte) int {
	if len(selectors) == 0 {
		return -1
	}

	for i := offset + 1; i < len(s.dom.builder.attributes); i++ { // s.dom.attributes[0] is a sentinel
		if !s.matchTag(s.dom.builder.attributes[i].tag, selectors) {
			continue
		}

		if !s.matchAttribute(i, selectors) {
			continue
		}

		if !s.matchAttributeValue(i, selectors) {
			continue
		}

		s.updateAttributeValue(i, newValue)
		return i
	}

	return -1
}

func (s *Session) updateAttributeValue(i int, newValue []byte) {
	currOffset := s.buffer.replaceBytes(s.attribute(i).boundaries[1], s.attributesOffsets[i], s.offsetDiff(i), newValue)
	s.attributesOffsets[i] += currOffset
	for j := i + 1; j < len(s.dom.builder.attributes); j++ {
		s.attributesOffsets[j] += currOffset
	}
}

func (s *Session) attribute(i int) *attr {
	return s.dom.builder.attributes[i]
}

//SetAttributeByIndex sets attribute value for n-th attribute from the Root DOM element
func (s *Session) SetAttributeByIndex(i int, value []byte) {
	s.updateAttributeValue(i, value)
}

//Attribute returns attribute value and index matches selectors, -1 if not found
func (s *Session) Attribute(offset int, selectors ...[]byte) ([]byte, int) {
	for i := offset + 1; i < len(s.dom.builder.attributes); i++ { // s.dom.attributes[0] is a sentinel
		if !s.matchTag(s.dom.builder.attributes[i].tag, selectors) {
			continue
		}

		if !s.matchAttribute(i, selectors) {
			continue
		}

		if !s.matchAttributeValue(i, selectors) {
			continue
		}

		return s.attributeValue(i), i
	}
	return nil, -1
}

//AttributeByIndex returns value of n-th attribute from the Root DOM element
func (s *Session) AttributeByIndex(i int) []byte {
	return s.attributeValue(i)
}

func (s *Session) attributeValue(i int) []byte {
	return s.buffer.buffer[s.attribute(i).valueStart()+s.attributesOffsets[i-1] : s.attribute(i).valueEnd()+s.attributesOffsets[i]]
}

//InnerHTML returns InnerHTML of element matching selectors and its index, -1 if not found
func (s *Session) InnerHTML(offset int, selectors ...[]byte) ([]byte, int) {
	tagIndex := s.findMatchingTag(offset, selectors)
	if tagIndex == -1 {
		return nil, tagIndex
	}

	return s.InnerHTMLByIndex(tagIndex), tagIndex
}

//InnerHTMLByIndex returns InnerHTML of n-th tag
func (s *Session) InnerHTMLByIndex(tagIndex int) []byte {
	return s.buffer.slice(s.tag(tagIndex).innerHTML, s.tagOffset(tagIndex), s.tagOffset(tagIndex))
}

func (s *Session) findMatchingTag(offset int, selectors [][]byte) int {
	for i := offset + 1; i < len(s.dom.builder.tags); i++ {
		if !s.matchTag(i, selectors) {
			continue
		}

		for j := s.dom.builder.tags[i-1].attrEnd; j < s.dom.builder.tags[i].attrEnd; j++ {
			if !s.matchAttribute(j, selectors) {
				continue
			}

			if !s.matchAttributeValue(j, selectors) {
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

func (s *Session) matchTag(i int, selectors [][]byte) bool {
	if s.skipped[i] {
		return false
	}

	return len(selectors) == 0 || bytes.Equal(s.buffer.slice(s.tag(i).tagName, s.tagOffset(i-1), s.tagOffset(i-1)), selectors[0])
}

//Bytes returns template after DOM changes
func (s *Session) Bytes() []byte {
	return s.buffer.bytes()
}

func (s *Session) tag(i int) *tag {
	return s.dom.builder.tags[i]
}

func (s *Session) offsetDiff(i int) int {
	return s.attributesOffsets[i-1] - s.attributesOffsets[i]
}

func (s *Session) attrByIndex(i int) *attr {
	return s.dom.builder.attributes[i]
}

func (s *Session) attrOffset(i int) int {
	return s.attributesOffsets[i]
}

func (s *Session) tagOffset(i int) int {
	return s.dom.builder.tags.tagOffset(i, s.attributesOffsets)
}

//SetInnerHTML updates InnerHTML of first element that matches selectors, returns the index of matched element
func (s *Session) SetInnerHTML(offset int, value []byte, selectors ...[]byte) (int, error) {
	tagIndex := s.findMatchingTag(offset, selectors)
	if tagIndex == -1 {
		return tagIndex, nil
	}

	return s.setInnerHTML(tagIndex, value)
}

func (s *Session) setInnerHTML(tagIndex int, value []byte) (int, error) {
	if err := s.updateInnerHTML(tagIndex, value); err != nil {
		return -1, err
	}
	s.attributesOffsets = make([]int, len(s.dom.builder.attributes))
	return tagIndex, nil
}

//SetInnerHTMLByIndex updates InnerHTML of n-th element from the DOM Root element
func (s *Session) SetInnerHTMLByIndex(tagIndex int, value []byte) (int, error) {
	return s.setInnerHTML(tagIndex, value)
}

func (s *Session) updateInnerHTML(tagIndex int, newInnerHTML []byte) error {
	diff := s.buffer.replaceBytes(s.tag(tagIndex).innerHTML, s.tagOffset(tagIndex), s.innerIncreased[tagIndex], newInnerHTML)
	for i := s.tag(tagIndex).attrEnd; i < len(s.attributesOffsets); i++ {
		s.attributesOffsets[i] += diff
	}

	for i := tagIndex + 1; i < len(s.skipped); i++ {
		if s.tag(tagIndex).depth <= s.tag(i).depth {
			break
		}

		s.skipped[i] = true
	}
	return nil
}

func (s *Session) innerBounds(tagIndex int) (int, int) {
	return s.tagOffset(tagIndex), s.tagOffset(tagIndex)
}

func (s *Session) matchAttribute(i int, selectors [][]byte) bool {
	if len(selectors) < 1 {
		return true
	}

	return bytes.Equal(s.buffer.slice(s.dom.builder.attributes[i].boundaries[0], s.attributesOffsets[i-1], s.attributesOffsets[i-1]), selectors[1])
}

func (s *Session) matchAttributeValue(i int, selectors [][]byte) bool {
	if len(selectors) < 3 {
		return true
	}

	if s.attributesOffsets[i] != s.attributesOffsets[i-1] {
		return bytes.Equal(s.buffer.slice(s.attrByIndex(i).boundaries[1], s.attributesOffsets[i-1], -(s.attributesOffsets[i-1]-s.attributesOffsets[i])), selectors[2])
	}

	return bytes.Equal(s.buffer.slice(s.attrByIndex(i).boundaries[1], s.attributesOffsets[i], s.attributesOffsets[i]), selectors[2])
}

func (s *Session) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case *Buffer:
			s.buffer = actual
		}
	}
}
