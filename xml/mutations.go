package xml

type (
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
		index         int
		value         string
		newElements   []*newElement
		newAttributes []*newAttribute
		valueChanged  bool
	}

	newElement struct {
		value string
	}

	newAttribute struct {
		key   string
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
		if len(m.attributes) < mapSize {
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

	return mutation.newValue, true
}

func (m *mutations) addElement(index int, value string) {
	mutation, ok := m.elementMutations(index)
	if ok {
		mutation.newElements = append(mutation.newElements, elementOf(value))
		return
	}

	elemMutation := &elementMutation{newElements: []*newElement{elementOf(value)}, index: index}
	m.addElementMutation(elemMutation)
}

func elementOf(value string) *newElement {
	return &newElement{
		value: value,
	}
}

func (m *mutations) elementMutations(index int) (*elementMutation, bool) {
	if m.elementsIndex != nil {
		if len(m.elements) < mapSize {
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

func (m *mutations) addAttribute(elemIndex int, key string, value string) {
	elemMutation, ok := m.elementMutations(elemIndex)
	newAttr := &newAttribute{
		key:   key,
		value: value,
	}

	if ok {
		elemMutation.newAttributes = append(elemMutation.newAttributes, newAttr)
		return
	}

	elemMutation = &elementMutation{
		index:         elemIndex,
		newAttributes: []*newAttribute{newAttr},
	}

	m.addElementMutation(elemMutation)
}

func (m *mutations) addElementMutation(mutation *elementMutation) {
	if m.elementsIndex == nil {
		m.elements[mutation.index] = mutation
	} else {
		m.elementsIndex[mutation.index] = len(m.elements)
		m.elements = append(m.elements, mutation)
	}
}

func (m *mutations) setValue(elemIndex int, value string) {
	mutation, ok := m.elementMutations(elemIndex)
	if ok {
		mutation.value = value
		mutation.newElements = nil
		mutation.valueChanged = true
		return
	}

	mutation = &elementMutation{
		index:        elemIndex,
		value:        value,
		valueChanged: true,
	}
	m.addElementMutation(mutation)
}

func newMutations(vxml *Schema) mutations {
	var attributesIndex map[int]int
	var attributesMutations []*attributeMutation
	if vxml.builder.attributeCounter < prealocateSize {
		attributesMutations = make([]*attributeMutation, vxml.builder.attributeCounter)
	} else {
		attributesIndex = map[int]int{}
	}

	var elementsMutations []*elementMutation
	var elementsMutationsIndex map[int]int
	if len(vxml.elements) < prealocateSize {
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
