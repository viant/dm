package xml

type (
	//DOM modifies the VirtualDOM
	DOM struct {
		vdom   *VirtualDOM
		buffer *Buffer
		changes
	}
)

//DOM returns new *DOM
func (s *VirtualDOM) DOM() *DOM {
	return &DOM{
		vdom:    s,
		buffer:  NewBuffer(s.bufferSize),
		changes: newMutations(s.schema),
	}
}

//Render returns XML value after changes
func (d *DOM) Render() string {
	return d.render(0, len(d.vdom.elements))
}

func (d *DOM) render(lowerBound, upperBound int) string {
	d.buffer.pos = 0
	elemChanges, _ := d.elementMutations(lowerBound)
	elem := d.vdom.elements[lowerBound]
	prevEnd, valueChanged := d.renderValue(d.vdom.elements[lowerBound].value.start, elem, elemChanges, true)
	i := lowerBound
	if valueChanged {
		i = d.nextNotDescendant(elem)
	} else {
		i += 1
	}

	for i < upperBound {
		elem = d.vdom.elements[i]
		elemChanges, _ = d.elementMutations(i)

		d.buffer.appendBytes(d.vdom.template[prevEnd:elem.tag.start])
		d.renderElementsBefore(elemChanges)
		prevEnd = d.renderTag(elem.tag.start, elem, elemChanges)
		prevEnd, valueChanged = d.renderValue(prevEnd, elem, elemChanges, false)

		if valueChanged {
			i = d.nextNotDescendant(elem)
			continue
		}

		i++
	}

	if lowerBound == 0 && upperBound == len(d.vdom.elements) {
		d.buffer.appendBytes(d.vdom.template[prevEnd:])
	} else {
		d.buffer.appendBytes(d.vdom.template[prevEnd:d.vdom.elements[lowerBound].value.end])
	}

	return d.buffer.String()
}

func (d *DOM) nextNotDescendant(elem *startElement) int {
	if elem.nextSibling != -1 {
		return elem.nextSibling
	}

	lastChild := elem
	for len(lastChild.children) > 0 {
		lastChild = lastChild.children[len(lastChild.children)-1]
	}

	return lastChild.elemIndex + 1
}

//Select returns Iterator over matching Elements
func (d *DOM) Select(selectors ...Selector) *Iterator {
	return newIterator(d, selectors)
}

func (d *DOM) renderAttributeValue(attribute *attribute) {
	value, ok := d.changes.checkAttributeChanges(attribute.index)
	if ok {
		d.buffer.appendBytes([]byte(value))
	} else {
		d.buffer.appendBytes(d.vdom.template[attribute.valueStart():attribute.valueEnd()])
	}
}

func (d *DOM) renderNewAttributes(elem *startElement) {
	mutation, ok := d.changes.elementMutations(elem.elemIndex)
	if !ok {
		return
	}

	for _, newAttr := range mutation.newAttributes {
		d.buffer.appendByte(' ')
		d.buffer.appendBytes([]byte(newAttr.key))
		d.buffer.appendBytes([]byte(`="`))
		d.buffer.appendBytes([]byte(newAttr.value))
		d.buffer.appendByte('"')
	}
}

func (d *DOM) templateSlice(span *span) string {
	return d.vdom.templateSlice(span)
}

func (d *DOM) insertBefore(index int, value string) {
	mutations, ok := d.elementMutations(index)

	if !ok {
		mutations = &elementChanges{index: index}
		d.addElementChanges(mutations)
	}

	mutations.elementsBefore = append(mutations.elementsBefore, value)
}

func (d *DOM) renderTag(prevEnd int, elem *startElement, elemChanges *elementChanges) int {
	d.buffer.appendBytes(d.vdom.template[prevEnd:elem.tag.start])

	prevEnd = elem.tag.start
	for _, attr := range elem.attributes {
		d.buffer.appendBytes(d.vdom.template[prevEnd:attr.valueStart()])
		d.renderAttributeValue(attr)
		prevEnd = attr.valueEnd()
	}

	if elemChanges != nil {
		d.buffer.appendByte(d.vdom.template[prevEnd])
		prevEnd++
		for _, newAttr := range elemChanges.newAttributes {
			d.buffer.appendByte(' ')
			d.buffer.appendBytes([]byte(newAttr.key))
			d.buffer.appendBytes([]byte(`="`))
			d.buffer.appendBytes([]byte(newAttr.value))
			d.buffer.appendByte('"')
		}
	}

	d.buffer.appendBytes(d.vdom.template[prevEnd:elem.value.start])
	return elem.value.start
}

func (d *DOM) renderValue(prevEnd int, elem *startElement, elemChanges *elementChanges, valueOnly bool) (int, bool) {
	d.renderNewElements(elemChanges)

	if elemChanges != nil && elemChanges.valueChanged {
		d.buffer.appendBytes([]byte(elemChanges.value))
		return elem.value.end, true
	}

	if len(elem.children) > 0 {
		d.buffer.appendBytes(d.vdom.template[prevEnd:elem.children[0].tag.start])
		return elem.children[0].tag.start, false
	}

	if valueOnly {
		d.buffer.appendBytes(d.vdom.template[elem.value.start:elem.value.end])
		return elem.value.end, false
	}

	d.buffer.appendBytes(d.vdom.template[elem.value.start:elem.tag.end])

	return elem.tag.end, false
}

func (d *DOM) renderElementsBefore(elemChanges *elementChanges) {
	if elemChanges == nil {
		return
	}

	for _, newElem := range elemChanges.elementsBefore {
		d.buffer.appendBytes([]byte(newElem))
	}
}

func (d *DOM) renderNewElements(elemChanges *elementChanges) {
	if elemChanges == nil {
		return
	}

	for _, element := range elemChanges.newElements {
		d.buffer.appendBytes([]byte(element.value))
	}
}
