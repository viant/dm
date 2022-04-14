package dm

import (
	"bytes"
)

type (
	//DOM modifies the VirtualDOM
	DOM struct {
		dom    *VirtualDOM
		buffer *Buffer
		mutations
	}

	mutations struct {
		attributesStart []int
		innerHTMLSize   []int
		tagsRemoved     []bool
	}
)

//DOM creates new DOM
func (v *VirtualDOM) DOM(options ...Option) *DOM {
	session := &DOM{
		dom: v,
		mutations: mutations{
			attributesStart: make([]int, len(v.attributes)),
			innerHTMLSize:   make([]int, len(v.tags)),
			tagsRemoved:     make([]bool, len(v.tags)),
		},
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
			current:   -1,
			next:      -1,
			selectors: selectors,
		},
	}
}

//SelectAttributes returns AttributeIterator to iterate over HTML Attributes
func (d *DOM) SelectAttributes(selectors ...string) *AttributeIterator {
	return &AttributeIterator{
		iterator: iterator{
			template:  d,
			current:   -1,
			next:      -1,
			selectors: selectors,
		},
	}
}

func (d *DOM) attribute(i int) *attr {
	return d.dom.attributes[i]
}

func (d *DOM) setAttributeByIndex(i int, value []byte) {
	d.updateAttributeValue(i, value)
}

func (d *DOM) updateAttributeValue(i int, newValue []byte) {
	currOffset := d.buffer.replaceBytes(d.attribute(i).boundaries[1], d.attributesStart[i], d.offsetDiff(i), newValue)
	d.attributesStart[i] += currOffset
	for j := i + 1; j < len(d.dom.attributes); j++ {
		d.attributesStart[j] += currOffset
	}
}

func (d *DOM) nextAttribute(offset int, selectors ...string) (newOffset int, index int) {
	if len(selectors) <= 1 {
		newOffset, index = d.matchAttributeByTag(offset, selectors)
	} else {
		newOffset, index = d.matchAttributeByAttributeName(offset, selectors)
	}

	return newOffset, index
}

func (d *DOM) matchAttributeByTag(offset int, selectors []string) (int, int) {
	if offset == 0 {
		offset = 1
	}
	for i := offset; i < len(d.dom.attributes); i++ {
		if !d.matchTag(d.attribute(i).tag, selectors) {
			continue
		}
		return i, i
	}
	return -1, -1
}

func (d *DOM) attributeByIndex(i int) []byte {
	return d.attributeValue(i)
}

//innerHTMLByIndex returns innerHTML of n-th tag
func (d *DOM) innerHTMLByIndex(tagIndex int) []byte {
	return d.buffer.slice(d.tag(tagIndex).innerHTML, d.tagOffset(tagIndex), d.tagOffset(tagIndex))
}

func (d *DOM) nextMatchingTag(offset int, selectors []string) int {
	if len(selectors) == 0 {
		return offset + 1
	}

	groupIndex := d.dom.index.tagIndex(selectors[0], false)
	if groupIndex == -1 || len(d.dom.tagsGrouped[groupIndex]) <= offset {
		return -1
	}

	tagIndex := d.dom.tagsGrouped[groupIndex][offset]
	for i := tagIndex; i < len(d.dom.tags); i++ {
		if d.tagsRemoved[i] {
			continue
		}

		if len(selectors) == 1 {
			return i
		}

		for j := d.tag(i - 1).attrEnd; j < d.tag(i).attrEnd; j++ {
			if !d.matchAttributeName(j, selectors) {
				continue
			}

			if !d.matchAttributeValue(j, selectors) {
				continue
			}

			return i
		}
	}
	return -1
}

func (d *DOM) matchTag(i int, selectors []string) bool {
	if d.tagsRemoved[i] {
		return false
	}

	return len(selectors) == 0 || bytes.EqualFold(
		d.buffer.slice(d.tag(i).tagName, d.tagOffset(i-1), d.tagOffset(i-1)),
		asBytes(selectors[0]),
	)
}

//Render returns template after VirtualDOM changes
func (d *DOM) Render() string {
	return string(d.buffer.bytes())
}

func (d *DOM) tag(i int) *tag {
	return d.dom.tags[i]
}

func (d *DOM) offsetDiff(i int) int {
	return d.attributesStart[i] - d.attributesStart[i-1]
}

func (d *DOM) attrByIndex(i int) *attr {
	return d.dom.attributes[i]
}

func (d *DOM) attrOffset(i int) int {
	return d.attributesStart[i]
}

func (d *DOM) tagOffset(i int) int {
	return d.dom.tags.tagOffset(i, d.attributesStart)
}

func (d *DOM) setInnerHTMLByIndex(tagIndex int, value []byte) error {
	if err := d.updateInnerHTML(tagIndex, value); err != nil {
		return err
	}
	return nil
}

func (d *DOM) updateInnerHTML(tagIndex int, newInnerHTML []byte) error {
	diff := d.buffer.replaceBytes(d.tag(tagIndex).innerHTML, d.tagOffset(tagIndex), d.innerHTMLSize[tagIndex], newInnerHTML)
	for i := d.tag(tagIndex).attrEnd; i < len(d.attributesStart); i++ {
		d.attributesStart[i] += diff
	}

	for i := tagIndex + 1; i < len(d.tagsRemoved); i++ {
		if d.tag(tagIndex).depth <= d.tag(i).depth {
			break
		}

		d.tagsRemoved[i] = true
	}
	return nil
}

func (d *DOM) matchAttributeName(i int, selectors []string) bool {
	if len(selectors) <= 1 {
		return true
	}

	return bytes.EqualFold(
		d.buffer.slice(d.dom.attributes[i].boundaries[0], d.attributesStart[i-1], d.attributesStart[i-1]),
		asBytes(selectors[1]),
	)
}

func (d *DOM) matchAttributeValue(i int, selectors []string) bool {
	if len(selectors) < 3 {
		return true
	}

	if d.attributesStart[i] != d.attributesStart[i-1] {
		return bytes.EqualFold(
			d.buffer.slice(d.attrByIndex(i).boundaries[1], d.attributesStart[i-1], d.offsetDiff(i)),
			asBytes(selectors[2]),
		)
	}

	return bytes.EqualFold(
		d.buffer.slice(d.attrByIndex(i).boundaries[1], d.attributesStart[i], d.attributesStart[i]),
		asBytes(selectors[2]),
	)
}

func (d *DOM) attributeKey(index int) []byte {
	return d.buffer.slice(d.attribute(index).boundaries[0], d.attrOffset(index-1), d.attrOffset(index-1))
}

func (d *DOM) attributeValue(i int) []byte {
	return d.buffer.buffer[d.attribute(i).valueStart()+d.attributesStart[i-1] : d.attribute(i).valueEnd()+d.attributesStart[i]]
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
	return len(d.dom.tags)
}

func (d *DOM) tagAttributes(i int) attrs {
	return d.dom.attributes[d.tag(i-1).attrEnd:d.tag(i).attrEnd]
}

func (d *DOM) matchAttributeByAttributeName(offset int, selectors []string) (int, int) {
	groupIndex := d.dom.attributeIndex(selectors[1], false)
	if groupIndex == -1 {
		return -1, -1
	}

	for i := offset; i < len(d.dom.attributesGrouped[groupIndex]); i++ {
		attrIndex := d.dom.attributesGrouped[groupIndex][i]
		if !d.matchTag(d.attribute(attrIndex).tag, selectors) {
			continue
		}

		if len(selectors) < 2 && !d.matchAttributeValue(attrIndex, selectors) {
			continue
		}
		return i, attrIndex
	}

	return -1, -1
}
