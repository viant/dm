package html

import (
	"bytes"

	"github.com/viant/dm/option"
)

type (
	//DOM modifies the VirtualDOM
	DOM struct {
		dom    *VirtualDOM
		buffer *Buffer
		mutations
	}

	mutations struct {
		attributesStart []int32
		innerHTMLSize   []int32
		tagsRemoved     []bool
	}
)

// DOM creates new DOM
func (v *VirtualDOM) DOM(options ...option.Option) *DOM {
	session := &DOM{
		dom: v,
		mutations: mutations{
			attributesStart: make([]int32, len(v.attributes)),
			innerHTMLSize:   make([]int32, len(v.tags)),
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

// Select returns ElementIterator to iterate over HTML Elements
func (d *DOM) Select(selectors ...string) *ElementIterator {
	return &ElementIterator{
		iterator: iterator{
			template:  d,
			current:   -1,
			next:      -1,
			selectors: selectors,
		},
		index: -1,
	}
}

// SelectAttributes returns AttributeIterator to iterate over HTML Attributes
func (d *DOM) SelectAttributes(selectors ...string) *AttributeIterator {
	return &AttributeIterator{
		iterator: iterator{
			template:  d,
			current:   -1,
			next:      -1,
			selectors: selectors,
		},
		index: -1,
	}
}

func (d *DOM) attribute(i int32) *attr {
	return d.dom.attributes[i]
}

func (d *DOM) setAttributeByIndex(i int32, value []byte) {
	d.updateAttributeValue(i, value)
}

func (d *DOM) updateAttributeValue(i int32, newValue []byte) {
	currOffset := d.buffer.replaceBytes(d.attribute(i).boundaries[1], d.attributesStart[i], d.offsetDiff(i), newValue)
	d.attributesStart[i] += currOffset
	d.updateAttributesStart(i, currOffset)
}

func (d *DOM) updateAttributesStart(i int32, currOffset int32) {
	for j := i + 1; j < int32(len(d.dom.attributes)); j++ {
		d.attributesStart[j] += currOffset
	}
}

func (d *DOM) nextAttribute(offset int32, selectors ...string) (newOffset int32, index int32) {
	if len(selectors) <= 1 {
		newOffset, index = d.matchAttributeByTag(offset, selectors)
	} else {
		newOffset, index = d.matchAttributeByAttributeName(offset, selectors)
	}

	return newOffset, index
}

func (d *DOM) matchAttributeByTag(offset int32, selectors []string) (int32, int32) {
	if offset == 0 {
		offset = 1
	}
	for i := offset; i < int32(len(d.dom.attributes)); i++ {
		if !d.matchTag(d.attribute(i).tag, selectors) {
			continue
		}
		return i, i
	}
	return -1, -1
}

func (d *DOM) attributeByIndex(i int32) []byte {
	return d.attributeValue(i)
}

// innerHTMLByIndex returns innerHTML of n-th tag
func (d *DOM) innerHTMLByIndex(tagIndex int32) []byte {
	return d.buffer.slice(d.tag(tagIndex).innerHTML, d.tagOffset(tagIndex), d.tagOffset(tagIndex))
}

func (d *DOM) nextMatchingTag(offset int32, selectors []string) (int32, int32) {
	if len(selectors) == 0 {
		i := d.findFirstTag(offset)
		return i, i

	}

	groupIndex := d.dom.index.tagIndex(selectors[0], false)
	if groupIndex == -1 || int32(len(d.dom.tagsGrouped[groupIndex])) <= offset {
		return -1, -1
	}

	tagIndex := d.dom.tagsGrouped[groupIndex][offset]
	for i := tagIndex; i < int32(len(d.dom.tags)); i++ {
		if d.tagsRemoved[i] {
			continue
		}

		if len(selectors) == 1 {
			return i, offset + i - tagIndex
		}

		for j := d.tag(i - 1).attrEnd; j < d.tag(i).attrEnd; j++ {
			if !d.matchAttributeName(j, selectors) {
				continue
			}

			if !d.matchAttributeValue(j, selectors) {
				continue
			}

			return i, offset + i - tagIndex
		}
	}
	return -1, -1
}

func (d *DOM) findFirstTag(offset int32) int32 {
	for i := offset + 1; i < int32(len(d.dom.tags)); i++ {
		if d.tagsRemoved[i] {
			continue
		}

		return i
	}
	return -1
}

func (d *DOM) matchTag(i int32, selectors []string) bool {
	if d.tagsRemoved[i] {
		return false
	}

	return len(selectors) == 0 || bytes.EqualFold(
		d.buffer.slice(d.tag(i).tagName, d.tagOffset(i-1), d.tagOffset(i-1)),
		asBytes(selectors[0]),
	)
}

// Render returns template after VirtualDOM changes
func (d *DOM) Render() string {
	return string(d.buffer.bytes())
}

func (d *DOM) tag(i int32) *tag {
	return d.dom.tags[i]
}

func (d *DOM) offsetDiff(i int32) int32 {
	return d.attributesStart[i] - d.attributesStart[i-1]
}

func (d *DOM) attrByIndex(i int32) *attr {
	return d.dom.attributes[i]
}

func (d *DOM) attrOffset(i int32) int32 {
	return d.attributesStart[i]
}

func (d *DOM) tagOffset(i int32) int32 {
	return d.dom.tags.tagOffset(i, d.attributesStart)
}

func (d *DOM) setInnerHTMLByIndex(tagIndex int32, value []byte) error {
	if err := d.updateInnerHTML(tagIndex, value); err != nil {
		return err
	}
	return nil
}

func (d *DOM) updateInnerHTML(tagIndex int32, newInnerHTML []byte) error {
	diff := d.buffer.replaceBytes(d.tag(tagIndex).innerHTML, d.tagOffset(tagIndex), d.innerHTMLSize[tagIndex], newInnerHTML)
	for i := d.tag(tagIndex).attrEnd; i < int32(len(d.attributesStart)); i++ {
		d.attributesStart[i] += diff
	}

	for i := tagIndex + 1; i < int32(len(d.tagsRemoved)); i++ {
		if d.tag(tagIndex).depth <= d.tag(i).depth {
			break
		}
		d.tagsRemoved[i] = true
	}

	d.innerHTMLSize[tagIndex] += diff
	return nil
}

func (d *DOM) matchAttributeName(i int32, selectors []string) bool {
	if len(selectors) <= 1 {
		return true
	}

	return bytes.EqualFold(
		d.buffer.slice(d.dom.attributes[i].boundaries[0], d.attributesStart[i-1], d.attributesStart[i-1]),
		asBytes(selectors[1]),
	)
}

func (d *DOM) matchAttributeValue(i int32, selectors []string) bool {
	if len(selectors) < 3 {
		return true
	}

	slice := d.buffer.slice(d.attrByIndex(i).boundaries[1], d.attributesStart[i], d.attributesStart[i])
	return bytes.Equal(
		slice,
		asBytes(selectors[2]),
	)
}

func (d *DOM) attributeKey(index int32) []byte {
	return d.buffer.slice(d.attribute(index).boundaries[0], d.attrOffset(index-1), d.attrOffset(index-1))
}

func (d *DOM) attributeValue(i int32) []byte {
	return d.buffer.buffer[d.attribute(i).valueStart()+d.attributesStart[i-1] : d.attribute(i).valueEnd()+d.attributesStart[i]]
}

func (d *DOM) apply(options []option.Option) {
	for _, opt := range options {
		switch actual := opt.(type) {
		case *Buffer:
			d.buffer = actual
		}
	}
}

func (d *DOM) tagLen() int32 {
	return int32(len(d.dom.tags))
}

func (d *DOM) tagAttributes(i int32) attrs {
	return d.dom.attributes[d.tag(i-1).attrEnd:d.tag(i).attrEnd]
}

func (d *DOM) matchAttributeByAttributeName(offset int32, selectors []string) (int32, int32) {
	groupIndex := d.dom.attributeIndex(selectors[1], false)
	if groupIndex == -1 {
		return -1, -1
	}

	for i := offset; i < int32(len(d.dom.attributesGrouped[groupIndex])); i++ {
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

func (d *DOM) addAttribute(tagIndex int32, key string, value string) {
	// adding to the buffer: ` key="value"`
	newAttribute := make([]byte, len(key)+len(value)+4)
	newAttribute[0] = ' '
	offset := 1
	offset += copy(newAttribute[offset:], key)
	offset += copy(newAttribute[offset:], `="`)
	offset += copy(newAttribute[offset:], value)
	offset += copy(newAttribute[offset:], `"`)

	end := d.tag(tagIndex).attrEnd - 1
	d.buffer.insertAfter(d.attribute(end).valueEnd()+1, d.attrOffset(end), newAttribute)

	for i := end; i < int32(len(d.dom.attributes)); i++ {
		d.attributesStart[i] += int32(offset)
	}
}
