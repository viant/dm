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
	}

	attributeMutation struct {
		newValue string
		index    int
	}
)

func (m *mutations) update(index int, value string) {
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

func (m *mutations) value(index int) (string, bool) {
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

	return mutations{
		attributes:      attributesMutations,
		attributesIndex: attributesIndex,
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
	value, ok := x.mutations.value(attribute.index)
	if ok {
		x.buffer.appendBytes([]byte(value))
	} else {
		x.buffer.appendBytes(x.vXml.template[attribute.valueStart():attribute.valueEnd()])
	}
}
