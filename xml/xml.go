package xml

type (
	Xml struct {
		vXml   *VirtualXml
		buffer *Buffer
		mutations
	}
)

func (v *VirtualXml) Xml() *Xml {
	return &Xml{
		vXml:      v,
		buffer:    NewBuffer(v.bufferSize),
		mutations: newMutations(v.vxml),
	}
}

func newMutations(vxml *VirtualXml) mutations {
	var attributesIndex map[int]int
	var attributesMutations []*attributeMutation
	if vxml.builder.attributeCounter < 30 {
		attributesMutations = make([]*attributeMutation, vxml.builder.attributeCounter)
	} else {
		attributesIndex = map[int]int{}
	}

	var elementsMutations []*elementMutation
	var elementsMutationsIndex map[int]int
	if len(vxml.elements) < 30 {
		elementsMutations = make([]*elementMutation, len(vxml.elements))
	} else {
		elementsMutationsIndex = map[int]int{}
	}

	return mutations{
		attributes:      attributesMutations,
		attributesIndex: attributesIndex,
		elements:        elementsMutations,
		elementsIndex:   elementsMutationsIndex,
	}
}

func (x *Xml) Render() string {
	x.buffer.pos = 0

	var prevEnd int
	var valueChanged bool
	for i := 1; i < len(x.vXml.elements); {
		elem := x.vXml.elements[i]
		if len(elem.attributes) > 0 {
			x.buffer.appendBytes(x.vXml.template[prevEnd:elem.attributes[0].valueStart()])
			x.renderAttributeValue(elem.attributes[0])

			if len(elem.attributes) > 1 {
				x.buffer.appendBytes(x.vXml.template[elem.attributes[0].valueEnd():elem.attributes[1].keyStart()])
			}

			for i := 1; i < len(elem.attributes); i++ {
				x.renderAttribute(elem, i)
			}

			x.buffer.appendByte(x.vXml.template[elem.attributes[len(elem.attributes)-1].valueEnd()])
			x.renderNewAttributes(elem)
			x.buffer.appendBytes(x.vXml.template[elem.attributes[len(elem.attributes)-1].valueEnd()+1 : elem.start])
		} else {
			x.buffer.appendBytes(x.vXml.template[prevEnd:elem.start])
		}

		x.renderNewElements(elem)
		prevEnd, valueChanged = x.renderElemValue(elem)
		if valueChanged {
			i = x.nextNotChildIndex(elem)
			continue
		}

		i++
	}

	x.buffer.appendBytes(x.vXml.template[prevEnd:])

	return x.buffer.String()
}

func (x *Xml) nextNotChildIndex(elem *StartElement) int {
	if elem.nextSibling != -1 {
		return elem.nextSibling
	}

	lastChild := elem
	for len(lastChild.children) > 0 {
		lastChild = lastChild.children[len(lastChild.children)-1]
	}

	return lastChild.elemIndex + 1
}

func (x *Xml) renderAttribute(elem *StartElement, attributeIndex int) {
	x.buffer.appendBytes(x.vXml.template[elem.attributes[attributeIndex].keyStart():elem.attributes[attributeIndex].valueStart()])
	x.renderAttributeValue(elem.attributes[attributeIndex])

	if attributeIndex < len(elem.attributes)-1 {
		x.buffer.appendBytes(x.vXml.template[elem.attributes[attributeIndex].valueEnd():elem.attributes[attributeIndex+1].keyStart()])
	}
}

func (x *Xml) Select(selectors ...Selector) *Iterator {
	return NewIterator(x, selectors)
}

func (x *Xml) renderAttributeValue(attribute *attribute) {
	value, ok := x.mutations.attributeValue(attribute.index)
	if ok {
		x.buffer.appendBytes([]byte(value))
	} else {
		x.buffer.appendBytes(x.vXml.template[attribute.valueStart():attribute.valueEnd()])
	}
}

func (x *Xml) renderNewElements(elem *StartElement) {
	mutation, ok := x.mutations.elementMutations(elem.elemIndex)
	if !ok {
		return
	}

	for _, element := range mutation.newElements {
		if elem.indent != nil {
			x.buffer.appendBytes(elem.indent)
		}
		x.buffer.appendBytes([]byte(element.value))
	}
}

func (x *Xml) renderNewAttributes(elem *StartElement) {
	mutation, ok := x.mutations.elementMutations(elem.elemIndex)
	if !ok {
		return
	}

	for _, newAttr := range mutation.newAttributes {
		x.buffer.appendByte(' ')
		x.buffer.appendBytes([]byte(newAttr.key))
		x.buffer.appendBytes([]byte(`="`))
		x.buffer.appendBytes([]byte(newAttr.value))
		x.buffer.appendByte('"')
	}
}

func (x *Xml) renderElemValue(elem *StartElement) (int, bool) {
	elemMutation, ok := x.mutations.elementMutations(elem.elemIndex)
	elemStart := elem.start
	if ok && elemMutation.valueChanged {
		x.buffer.appendBytes([]byte(elemMutation.value))
		elemStart = elem.end
	}

	end := len(x.vXml.template)
	if elem.elemIndex < len(x.vXml.elements)-1 {
		nextElem := elem.elemIndex + 1
		end = x.vXml.elements[nextElem].start
		if len(x.vXml.elements[nextElem].attributes) > 0 {
			end = x.vXml.elements[nextElem].attributes[0].keyStart()
		}
	}

	if elemStart > end {
		return elemStart, true
	}
	x.buffer.appendBytes(x.vXml.template[elemStart:end])
	return end, elemStart == elem.end
}
