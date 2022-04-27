package xml

type (
	//Document modifies the VirtualDOM
	Document struct {
		dom    *DOM
		buffer *Buffer
		changes
	}
)

//Document returns new *Document
func (d *DOM) Document() *Document {
	return &Document{
		dom:     d,
		buffer:  NewBuffer(d.bufferSize),
		changes: newMutations(d.dom),
	}
}

//Render returns XML value after changes
func (d *Document) Render() string {
	return d.render(0, len(d.dom.elements))
}

func (d *Document) render(lowerBound, upperBound int) string {
	d.buffer.pos = 0
	elemChanges, _ := d.elementMutations(lowerBound)
	elem := d.dom.elements[lowerBound]

	prevEnd, i := d.renderElem(false, elem, elemChanges)
	for i < upperBound {
		elem = d.dom.elements[i]
		d.buffer.appendBytes(d.dom.template[prevEnd:elem.tag.start])
		elemChanges, _ = d.elementMutations(i)

		prevEnd, i = d.renderElem(true, elem, elemChanges)
	}

	if lowerBound == 0 && upperBound == len(d.dom.elements) {
		d.buffer.appendBytes(d.dom.template[prevEnd:])
	} else if prevEnd < d.dom.elements[lowerBound].value.end {
		d.buffer.appendBytes(d.dom.template[prevEnd:d.dom.elements[lowerBound].value.end])
	}

	return d.buffer.String()
}

func (d *Document) renderElem(withTag bool, elem *startElement, changes *elementChanges) (nextEnd int, nextElemIndex int) {
	if changes != nil && changes.replacedChanged {
		d.buffer.appendBytes([]byte(changes.replacedWith))
		return elem.tag.end, d.nextNotDescendant(elem)
	}

	d.renderElementsBefore(changes)

	if withTag {
		d.renderTag(elem.tag.start, elem, changes)
	}

	return d.renderValue(elem.value.start, elem, changes, !withTag)
}

func (d *Document) nextNotDescendant(elem *startElement) int {
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
func (d *Document) Select(selectors ...Selector) *Iterator {
	return newIterator(d, selectors)
}

//SelectFirst returns first Element that matches selectors
func (d *Document) SelectFirst(selectors ...Selector) (*Element, bool) {
	selMatcher := newMatcher(d, selectors)
	matched := selMatcher.match()
	if matched == -1 {
		return nil, false
	}

	return &Element{
		document:     d,
		startElement: selMatcher.currRoot,
	}, true
}

func (d *Document) renderAttributeValue(attribute *attribute) {
	value, ok := d.changes.checkAttributeChanges(attribute.index)
	if ok {
		d.buffer.appendBytes([]byte(value))
	} else {
		d.buffer.appendBytes(d.dom.template[attribute.valueStart():attribute.valueEnd()])
	}
}

func (d *Document) renderNewAttributes(elem *startElement) {
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

func (d *Document) templateSlice(span *span) string {
	return d.dom.templateSlice(span)
}

func (d *Document) insertBefore(index int, value string) {
	mutations, ok := d.elementMutations(index)

	if !ok {
		mutations = &elementChanges{index: index}
		d.addElementChanges(mutations)
	}

	mutations.elementsBefore = append(mutations.elementsBefore, value)
}

func (d *Document) renderTag(prevEnd int, elem *startElement, elemChanges *elementChanges) int {
	d.buffer.appendBytes(d.dom.template[prevEnd:elem.tag.start])

	prevEnd = elem.tag.start
	for _, attr := range elem.attributes {
		d.buffer.appendBytes(d.dom.template[prevEnd:attr.valueStart()])
		d.renderAttributeValue(attr)
		prevEnd = attr.valueEnd()
	}

	if elemChanges != nil {
		d.buffer.appendByte(d.dom.template[prevEnd])
		prevEnd++
		for _, newAttr := range elemChanges.newAttributes {
			d.buffer.appendByte(' ')
			d.buffer.appendBytes([]byte(newAttr.key))
			d.buffer.appendBytes([]byte(`="`))
			d.buffer.appendBytes([]byte(newAttr.value))
			d.buffer.appendByte('"')
		}
	}

	d.buffer.appendBytes(d.dom.template[prevEnd:elem.value.start])
	return elem.value.start
}

func (d *Document) renderValue(prevEnd int, elem *startElement, elemChanges *elementChanges, valueOnly bool) (int, int) {
	d.renderNewElements(elemChanges)

	if elemChanges != nil && elemChanges.valueChanged {
		d.buffer.appendBytes([]byte(elemChanges.value))
		return elem.value.end, d.nextNotDescendant(elem)
	}

	if len(elem.children) > 0 {
		d.buffer.appendBytes(d.dom.template[prevEnd:elem.children[0].tag.start])
		return elem.children[0].tag.start, elem.elemIndex + 1
	}

	if valueOnly {
		d.buffer.appendBytes(d.dom.template[elem.value.start:elem.value.end])
		return elem.value.end, elem.elemIndex + 1
	}

	d.buffer.appendBytes(d.dom.template[elem.value.start:elem.tag.end])

	return elem.tag.end, elem.elemIndex + 1
}

func (d *Document) renderElementsBefore(elemChanges *elementChanges) {
	if elemChanges == nil {
		return
	}

	for _, newElem := range elemChanges.elementsBefore {
		d.buffer.appendBytes([]byte(newElem))
	}
}

func (d *Document) renderNewElements(elemChanges *elementChanges) {
	if elemChanges == nil {
		return
	}

	for _, element := range elemChanges.newElements {
		d.buffer.appendBytes([]byte(element.value))
	}
}

func (d *Document) replace(element *startElement, newElement string) {
	mutations := &elementChanges{
		index:           element.elemIndex,
		replacedChanged: true,
		replacedWith:    newElement,
	}

	d.addElementChanges(mutations)
}

func (d *Document) renderReplacement(elemChanges *elementChanges) {
	d.buffer.appendBytes([]byte(elemChanges.replacedWith))
}

func (d *Document) wasReplaced(element *startElement) bool {
	mutations, ok := d.elementMutations(element.elemIndex)
	if ok && mutations.replacedChanged {
		return true
	}

	return false
}
