package xml

type (
	Xml struct {
		vXml   *VirtualXml
		buffer *Buffer
		mutations
	}

	mutations struct {
		attributes      []*attributeMutation
		attributesIndex map[int]int

		elements      []*elementMutation
		elementsIndex map[int]int
	}

	attributeMutation struct {
		newValue string
		index    int
	}

	elementMutation struct {
		index       int
		value       string
		newElements []*newElement
	}

	newElement struct {
		value string
	}
)

func (m *mutations) updateAttribute(index int, value string) {
	if m.attributesIndex == nil {
		m.attributes[index] = &attributeMutation{newValue: value, index: index}
	} else {
		sliceIndex, ok := m.attributesIndex[index]
		if !ok {
			m.attributes = append(m.attributes, &attributeMutation{newValue: value, index: index})
			m.attributesIndex[index] = len(m.attributes) - 1
		} else {
			m.attributes[sliceIndex].newValue = value
		}
	}
}

func (m *mutations) attributeValue(index int) (string, bool) {
	if m.attributesIndex != nil {
		if len(m.attributes) < 5 {
			for _, mutation := range m.attributes {
				if mutation != nil && mutation.index == index {
					return mutation.newValue, true
				}
			}

			return "", false
		}

		sliceIndex, ok := m.attributesIndex[index]
		if !ok {
			return "", false
		}
		return m.attributes[sliceIndex].newValue, true
	}

	mutation := m.attributes[index]

	if mutation == nil {
		return "", false
	}

	return mutation.newValue, mutation != nil
}

func (m *mutations) appendElement(index int, value string) {
	if m.elementsIndex == nil {
		m.updateElementUsingSlice(index, value)
	} else {
		m.updateElementUsingMap(index, value)
	}
}

func (m *mutations) updateElementUsingMap(index int, value string) {
	sliceIndex, ok := m.elementsIndex[index]
	if ok {
		m.elements[sliceIndex].newElements = append(m.elements[sliceIndex].newElements, elementOf(value))
	} else {
		m.elements = append(m.elements, &elementMutation{
			index:       index,
			newElements: []*newElement{elementOf(value)},
		})

		m.elementsIndex[index] = len(m.elements) - 1
	}
}

func elementOf(value string) *newElement {
	return &newElement{
		value: value,
	}
}

func (m *mutations) updateElementUsingSlice(index int, value string) {
	mutation := m.elements[index]
	if mutation != nil {
		mutation.newElements = append(mutation.newElements, elementOf(value))
	} else {
		m.elements[index] = &elementMutation{
			index:       index,
			newElements: []*newElement{elementOf(value)},
		}
	}
}

func (m *mutations) elementMutations(index int) (*elementMutation, bool) {
	if m.elementsIndex != nil {
		if len(m.elements) < 5 {
			for _, element := range m.elements {
				if element.index == index {
					return element, true
				}
			}

			return nil, false
		}

		sliceIndex, ok := m.elementsIndex[index]
		if !ok {
			return nil, false
		}
		return m.elements[sliceIndex], true
	}

	return m.elements[index], m.elements[index] != nil
}

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

	var prevStart int
	var elem *StartElement
	for _, elem = range x.vXml.allElements() {
		if len(elem.attributes) == 0 {
			x.buffer.appendBytes(x.vXml.template[prevStart:elem.start])
		} else {
			x.buffer.appendBytes(x.vXml.template[prevStart:elem.attributes[0].valueStart()])
			x.renderAttributeValue(elem.attributes[0])

			if len(elem.attributes) > 1 {
				x.buffer.appendBytes(x.vXml.template[elem.attributes[0].valueEnd():elem.attributes[1].keyStart()])
			}

			for i := 1; i < len(elem.attributes); i++ {
				x.renderAttribute(elem, i)
			}

			x.buffer.appendBytes(x.vXml.template[elem.attributes[len(elem.attributes)-1].valueEnd():elem.start])
		}

		x.renderNewElements(elem)
		prevStart = elem.start
	}

	if elem != nil {
		x.buffer.appendBytes(x.vXml.template[elem.start:elem.end])
		x.buffer.appendBytes(x.vXml.template[elem.end:])
	}

	return x.buffer.String()
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
