package dm

import (
	"bytes"
)

type (
	//Template modifies the DOM
	Template struct {
		dom               *DOM
		buffer            *Buffer
		attributesOffsets []int
		innerIncreased    []int
		skipped           []bool
		removedTags       map[int]bool
	}
)

//Template creates new Template
func (d *DOM) Template(options ...Option) *Template {
	session := &Template{
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

func (t *Template) SelectTags(selectors ...[]byte) *TagIterator {
	return &TagIterator{
		iterator: iterator{
			template:  t,
			current:   0,
			next:      0,
			selectors: selectors,
		},
	}
}

func (t *Template) SelectAttributes(selectors ...[]byte) *AttributeIterator {
	return &AttributeIterator{
		iterator: iterator{
			template:  t,
			current:   0,
			next:      0,
			selectors: selectors,
		},
	}
}

func (t *Template) attribute(i int) *attr {
	return t.dom.builder.attributes[i]
}

func (t *Template) setAttributeByIndex(i int, value []byte) {
	t.updateAttributeValue(i, value)
}

func (t *Template) updateAttributeValue(i int, newValue []byte) {
	currOffset := t.buffer.replaceBytes(t.attribute(i).boundaries[1], t.attributesOffsets[i], t.offsetDiff(i), newValue)
	t.attributesOffsets[i] += currOffset
	for j := i + 1; j < len(t.dom.builder.attributes); j++ {
		t.attributesOffsets[j] += currOffset
	}
}

func (t *Template) nextAttribute(offset int, selectors ...[]byte) int {
	for i := offset + 1; i < len(t.dom.builder.attributes); i++ { // t.dom.attributes[0] is a sentinel
		if !t.matchTag(t.dom.builder.attributes[i].tag, selectors) {
			continue
		}

		if !t.matchAttributeName(i, selectors) {
			continue
		}

		if !t.matchAttributeValue(i, selectors) {
			continue
		}

		return i
	}
	return -1
}

func (t *Template) attributeByIndex(i int) []byte {
	return t.attributeValue(i)
}

//innerHTMLByIndex returns innerHTML of n-th tag
func (t *Template) innerHTMLByIndex(tagIndex int) []byte {
	return t.buffer.slice(t.tag(tagIndex).innerHTML, t.tagOffset(tagIndex), t.tagOffset(tagIndex))
}

func (t *Template) findMatchingTag(offset int, selectors [][]byte) int {
	for i := offset + 1; i < len(t.dom.builder.tags); i++ {
		if !t.matchTag(i, selectors) {
			continue
		}

		for j := t.dom.builder.tags[i-1].attrEnd; j < t.dom.builder.tags[i].attrEnd; j++ {
			if !t.matchAttributeName(j, selectors) {
				continue
			}

			if !t.matchAttributeValue(j, selectors) {
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

func (t *Template) matchTag(i int, selectors [][]byte) bool {
	if t.skipped[i] {
		return false
	}

	return len(selectors) == 0 || bytes.Equal(t.buffer.slice(t.tag(i).tagName, t.tagOffset(i-1), t.tagOffset(i-1)), selectors[0])
}

//Bytes returns template after DOM changes
func (t *Template) Bytes() []byte {
	return t.buffer.bytes()
}

func (t *Template) tag(i int) *tag {
	return t.dom.builder.tags[i]
}

func (t *Template) offsetDiff(i int) int {
	return t.attributesOffsets[i-1] - t.attributesOffsets[i]
}

func (t *Template) attrByIndex(i int) *attr {
	return t.dom.builder.attributes[i]
}

func (t *Template) attrOffset(i int) int {
	return t.attributesOffsets[i]
}

func (t *Template) tagOffset(i int) int {
	return t.dom.builder.tags.tagOffset(i, t.attributesOffsets)
}

func (t *Template) setInnerHTMLByIndex(tagIndex int, value []byte) error {
	if err := t.updateInnerHTML(tagIndex, value); err != nil {
		return err
	}
	return nil
}

func (t *Template) updateInnerHTML(tagIndex int, newInnerHTML []byte) error {
	diff := t.buffer.replaceBytes(t.tag(tagIndex).innerHTML, t.tagOffset(tagIndex), t.innerIncreased[tagIndex], newInnerHTML)
	for i := t.tag(tagIndex).attrEnd; i < len(t.attributesOffsets); i++ {
		t.attributesOffsets[i] += diff
	}

	for i := tagIndex + 1; i < len(t.skipped); i++ {
		if t.tag(tagIndex).depth <= t.tag(i).depth {
			break
		}

		t.skipped[i] = true
	}
	return nil
}

func (t *Template) matchAttributeName(i int, selectors [][]byte) bool {
	if len(selectors) < 1 {
		return true
	}

	return bytes.Equal(t.buffer.slice(t.dom.builder.attributes[i].boundaries[0], t.attributesOffsets[i-1], t.attributesOffsets[i-1]), selectors[1])
}

func (t *Template) matchAttributeValue(i int, selectors [][]byte) bool {
	if len(selectors) < 3 {
		return true
	}

	if t.attributesOffsets[i] != t.attributesOffsets[i-1] {
		return bytes.Equal(t.buffer.slice(t.attrByIndex(i).boundaries[1], t.attributesOffsets[i-1], -(t.attributesOffsets[i-1]-t.attributesOffsets[i])), selectors[2])
	}

	return bytes.Equal(t.buffer.slice(t.attrByIndex(i).boundaries[1], t.attributesOffsets[i], t.attributesOffsets[i]), selectors[2])
}

func (t *Template) attributeKey(index int) []byte {
	return t.buffer.slice(t.attribute(index).boundaries[0], t.attrOffset(index-1), t.attrOffset(index-1))
}

func (t *Template) attributeValue(i int) []byte {
	return t.buffer.buffer[t.attribute(i).valueStart()+t.attributesOffsets[i-1] : t.attribute(i).valueEnd()+t.attributesOffsets[i]]
}

func (t *Template) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case *Buffer:
			t.buffer = actual
		}
	}
}

func (t *Template) tagLen() int {
	return len(t.dom.builder.tags)
}

func (t *Template) tagAttributes(i int) attrs {
	return t.dom.builder.attributes[t.tag(i-1).attrEnd:t.tag(i).attrEnd]
}
