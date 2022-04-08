package dm

import (
	"bytes"
)

type (
	//DOM modifies the VirtualDOM
	DOM struct {
		dom               *VirtualDOM
		buffer            *Buffer
		attributesOffsets []int
		innerIncreased    []int
		skipped           []bool
		removedTags       map[int]bool
	}
)

//DOM creates new DOM
func (v *VirtualDOM) DOM(options ...Option) *DOM {
	session := &DOM{
		dom:               v,
		attributesOffsets: make([]int, len(v.builder.attributes)),
		innerIncreased:    make([]int, len(v.builder.tags)),
		skipped:           make([]bool, len(v.builder.tags)),
	}

	session.apply(options)
	if session.buffer == nil {
		session.buffer = NewBuffer(v.initialBufferSize)
	}

	session.buffer.appendBytes(v.template)

	return session
}

//Select returns ElementIterator to iterate over HTML Elements
func (d *DOM) Select(selectors ...string) *ElementIterator {
	return &ElementIterator{
		iterator: iterator{
			template:  d,
			current:   0,
			next:      0,
			selectors: selectors,
		},
	}
}

//SelectAttributes returns AttributeIterator to iterate over HTML Attributes
func (d *DOM) SelectAttributes(selectors ...string) *AttributeIterator {
	return &AttributeIterator{
		iterator: iterator{
			template:  d,
			current:   0,
			next:      0,
			selectors: selectors,
		},
	}
}

func (d *DOM) attribute(i int) *attr {
	return d.dom.builder.attributes[i]
}

func (d *DOM) setAttributeByIndex(i int, value []byte) {
	d.updateAttributeValue(i, value)
}

func (d *DOM) updateAttributeValue(i int, newValue []byte) {
	currOffset := d.buffer.replaceBytes(d.attribute(i).boundaries[1], d.attributesOffsets[i], d.offsetDiff(i), newValue)
	d.attributesOffsets[i] += currOffset
	for j := i + 1; j < len(d.dom.builder.attributes); j++ {
		d.attributesOffsets[j] += currOffset
	}
}

func (d *DOM) nextAttribute(offset int, selectors ...string) int {
	for i := offset + 1; i < len(d.dom.builder.attributes); i++ { // d.dom.attributes[0] is a sentinel
		if !d.matchTag(d.dom.builder.attributes[i].tag, selectors) {
			continue
		}

		if !d.matchAttributeName(i, selectors) {
			continue
		}

		if !d.matchAttributeValue(i, selectors) {
			continue
		}

		return i
	}
	return -1
}

func (d *DOM) attributeByIndex(i int) []byte {
	return d.attributeValue(i)
}

//innerHTMLByIndex returns innerHTML of n-th tag
func (d *DOM) innerHTMLByIndex(tagIndex int) []byte {
	return d.buffer.slice(d.tag(tagIndex).innerHTML, d.tagOffset(tagIndex), d.tagOffset(tagIndex))
}

func (d *DOM) findMatchingTag(offset int, selectors []string) int {
	for i := offset + 1; i < len(d.dom.builder.tags); i++ {
		if !d.matchTag(i, selectors) {
			continue
		}

		for j := d.dom.builder.tags[i-1].attrEnd; j < d.dom.builder.tags[i].attrEnd; j++ {
			if !d.matchAttributeName(j, selectors) {
				continue
			}

			if !d.matchAttributeValue(j, selectors) {
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

func (d *DOM) matchTag(i int, selectors []string) bool {
	if d.skipped[i] {
		return false
	}

	return len(selectors) == 0 || bytes.Equal(
		d.buffer.slice(d.tag(i).tagName, d.tagOffset(i-1), d.tagOffset(i-1)),
		asBytes(selectors[0]),
	)
}

//Render returns template after VirtualDOM changes
func (d *DOM) Render() string {
	return string(d.buffer.bytes())
}

func (d *DOM) tag(i int) *tag {
	return d.dom.builder.tags[i]
}

func (d *DOM) offsetDiff(i int) int {
	return d.attributesOffsets[i] - d.attributesOffsets[i-1]
}

func (d *DOM) attrByIndex(i int) *attr {
	return d.dom.builder.attributes[i]
}

func (d *DOM) attrOffset(i int) int {
	return d.attributesOffsets[i]
}

func (d *DOM) tagOffset(i int) int {
	return d.dom.builder.tags.tagOffset(i, d.attributesOffsets)
}

func (d *DOM) setInnerHTMLByIndex(tagIndex int, value []byte) error {
	if err := d.updateInnerHTML(tagIndex, value); err != nil {
		return err
	}
	return nil
}

func (d *DOM) updateInnerHTML(tagIndex int, newInnerHTML []byte) error {
	diff := d.buffer.replaceBytes(d.tag(tagIndex).innerHTML, d.tagOffset(tagIndex), d.innerIncreased[tagIndex], newInnerHTML)
	for i := d.tag(tagIndex).attrEnd; i < len(d.attributesOffsets); i++ {
		d.attributesOffsets[i] += diff
	}

	for i := tagIndex + 1; i < len(d.skipped); i++ {
		if d.tag(tagIndex).depth <= d.tag(i).depth {
			break
		}

		d.skipped[i] = true
	}
	return nil
}

func (d *DOM) matchAttributeName(i int, selectors []string) bool {
	if len(selectors) < 1 {
		return true
	}

	return bytes.Equal(
		d.buffer.slice(d.dom.builder.attributes[i].boundaries[0], d.attributesOffsets[i-1], d.attributesOffsets[i-1]),
		asBytes(selectors[1]),
	)
}

func (d *DOM) matchAttributeValue(i int, selectors []string) bool {
	if len(selectors) < 3 {
		return true
	}

	if d.attributesOffsets[i] != d.attributesOffsets[i-1] {
		return bytes.Equal(
			d.buffer.slice(d.attrByIndex(i).boundaries[1], d.attributesOffsets[i-1], d.offsetDiff(i)),
			asBytes(selectors[2]),
		)
	}

	return bytes.Equal(
		d.buffer.slice(d.attrByIndex(i).boundaries[1], d.attributesOffsets[i], d.attributesOffsets[i]),
		asBytes(selectors[2]),
	)
}

func (d *DOM) attributeKey(index int) []byte {
	return d.buffer.slice(d.attribute(index).boundaries[0], d.attrOffset(index-1), d.attrOffset(index-1))
}

func (d *DOM) attributeValue(i int) []byte {
	return d.buffer.buffer[d.attribute(i).valueStart()+d.attributesOffsets[i-1] : d.attribute(i).valueEnd()+d.attributesOffsets[i]]
}

func (d *DOM) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case *Buffer:
			d.buffer = actual
		}
	}
}

func (d *DOM) tagLen() int {
	return len(d.dom.builder.tags)
}

func (d *DOM) tagAttributes(i int) attrs {
	return d.dom.builder.attributes[d.tag(i-1).attrEnd:d.tag(i).attrEnd]
}
