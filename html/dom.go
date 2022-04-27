package html

import (
	"bytes"
	"github.com/viant/dm/option"
)

type (
	//DOM modifies the VirtualDOM
	DOM struct {
		vdom   *VirtualDOM
		buffer *Buffer
		mutations
	}

	mutations struct {
		attrValueEnd  []int
		innerHTMLSize []int
		tagsRemoved   []bool
	}
)

//DOM creates new DOM
func (v *VirtualDOM) DOM(options ...option.Option) *DOM {
	session := &DOM{
		vdom: v,
		mutations: mutations{
			attrValueEnd:  make([]int, v.attributeCounter),
			innerHTMLSize: make([]int, len(v.tags)),
			tagsRemoved:   make([]bool, len(v.tags)),
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
		matcher: newElementMatcher(d, selectors),
		index:   -1,
	}
}

func (d *DOM) SelectFirst(selectors ...string) (*Element, bool) {
	next := newElementMatcher(d, selectors).match()
	if next == -1 {
		return nil, false
	}

	return &Element{
		template: d,
		tag:      d.tag(next),
	}, true
}

//SelectAttributes returns AttributeIterator to iterate over HTML Attributes
func (d *DOM) SelectAttributes(selectors ...string) *AttributeIterator {
	return &AttributeIterator{
		dom:     d,
		matcher: newAttributeMatcher(d, selectors),
	}
}

func (d *DOM) setAttribute(anAttr *attr, newValue []byte) {
	start := d.attrValueEnd[anAttr.index]
	currOffset := d.buffer.replaceBytes(anAttr.boundaries[1], start, d.offsetDiff(anAttr.index), newValue)

	if currOffset == 0 {
		return
	}

	for i := anAttr.index; i < len(d.attrValueEnd); i++ {
		d.attrValueEnd[i] += currOffset
	}
}

func (d *DOM) innerHTML(aTag *tag) []byte {
	start := d.tagValueOffset(aTag)
	slice := d.buffer.slice(aTag.innerHTML, start, start+d.innerHTMLSize[aTag.index])
	return slice
}

func (d *DOM) tagValueOffset(aTag *tag) int {
	if aTag.attrEnd >= 0 {
		return d.attrValueEnd[aTag.attrEnd]
	}

	return 0
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
	return d.vdom.tags[i]
}

func (d *DOM) offsetDiff(i int) int {
	if i == 0 {
		return 0
	}

	return d.attrValueEnd[i] - d.attrValueEnd[i-1]
}

func (d *DOM) attrOffset(i int) int {
	return d.attrValueEnd[i]
}

func (d *DOM) tagOffset(i int) int {
	return d.vdom.tags.tagOffset(i, d.attrValueEnd)
}

func (d *DOM) setInnerHTML(aTag *tag, value []byte) error {
	if err := d.updateInnerHTML(aTag, value); err != nil {
		return err
	}
	return nil
}

func (d *DOM) updateInnerHTML(aTag *tag, newInnerHTML []byte) error {
	diff := d.buffer.replaceBytes(aTag.innerHTML, d.tagValueOffset(aTag), d.innerHTMLSize[aTag.index], newInnerHTML)
	for i := aTag.attrEnd + 1; i < len(d.attrValueEnd); i++ {
		d.attrValueEnd[i] += diff
	}

	for i := aTag.index + 1; i < len(d.tagsRemoved); i++ {
		if aTag.depth <= d.tag(i).depth {
			break
		}
		d.tagsRemoved[i] = true
	}

	d.innerHTMLSize[aTag.index] += diff
	return nil
}

func (d *DOM) matchAttributeValue(anAttr *attr, selectors []string) bool {
	if len(selectors) < 3 {
		return true
	}

	slice := d.buffer.slice(anAttr.boundaries[1], d.attrValueOffset(anAttr), d.attrValueEnd[anAttr.index])
	return bytes.Equal(
		slice,
		asBytes(selectors[2]),
	)
}

func (d *DOM) attrValueOffset(anAttr *attr) int {
	if anAttr.index > 0 {
		return d.attrValueEnd[anAttr.index-1]
	}

	return 0
}

func (d *DOM) attributeKey(anAttr *attr) []byte {
	return d.buffer.slice(anAttr.boundaries[0], d.attrOffset(anAttr.index-1), d.attrOffset(anAttr.index-1))
}

func (d *DOM) attributeValue(anAttr *attr) []byte {
	start := d.attrValueOffset(anAttr)
	return d.buffer.buffer[anAttr.valueStart()+start : anAttr.valueEnd()+d.attrValueEnd[anAttr.index]]
}

func (d *DOM) apply(options []option.Option) {
	for _, opt := range options {
		switch actual := opt.(type) {
		case *Buffer:
			d.buffer = actual
		}
	}
}

func (d *DOM) tagLen() int {
	return len(d.vdom.tags)
}

func (d *DOM) addAttribute(aTag *tag, key string, value string) {
	// adding to the buffer: ` key="value"`
	newAttribute := make([]byte, len(key)+len(value)+4)
	newAttribute[0] = ' '
	offset := 1
	offset += copy(newAttribute[offset:], key)
	offset += copy(newAttribute[offset:], `="`)
	offset += copy(newAttribute[offset:], value)
	offset += copy(newAttribute[offset:], `"`)

	var end int
	if len(aTag.attrs) > 0 {
		end = aTag.attrs[len(aTag.attrs)-1].valueEnd()
	} else {
		end = aTag.tagName.end
	}

	start := d.tagValueOffset(aTag)
	d.buffer.insertAfter(end+1, start, newAttribute)
	for i := zeroIfNegative(aTag.attrEnd); i < len(d.attrValueEnd); i++ {
		d.attrValueEnd[i] += offset
	}
}

func zeroIfNegative(value int) int {
	if value < 0 {
		return 0
	}

	return value
}
