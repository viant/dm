package xml

type (
	changes struct {
		attributes      []*attributeChanges
		attributesIndex map[int]int

		elements      []*elementChanges
		elementsIndex map[int]int
	}

	attributeChanges struct {
		newValue string
		index    int
	}

	elementChanges struct {
		index           int
		value           string
		newElements     []*newElement
		newAttributes   []*newAttribute
		valueChanged    bool
		elementsBefore  []string
		replacedWith    string
		replacedChanged bool
	}

	newElement struct {
		value string
	}

	newAttribute struct {
		key   string
		value string
	}
)

func (m *changes) updateAttribute(index int, value string) {
	if m.attributesIndex == nil {
		m.attributes[index] = &attributeChanges{newValue: value, index: index}
	} else {
		sliceIndex, ok := m.attributesIndex[index]
		if !ok {
			m.attributes = append(m.attributes, &attributeChanges{newValue: value, index: index})
			m.attributesIndex[index] = len(m.attributes) - 1
		} else {
			m.attributes[sliceIndex].newValue = value
		}
	}
}

func (m *changes) checkAttributeChanges(index int) (string, bool) {
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

func (m *changes) addElement(index int, value string) {
	mutation, ok := m.elementMutations(index)
	if ok {
		mutation.newElements = append(mutation.newElements, elementOf(value))
		return
	}

	elemMutation := &elementChanges{newElements: []*newElement{elementOf(value)}, index: index}
	m.addElementChanges(elemMutation)
}

func elementOf(value string) *newElement {
	return &newElement{
		value: value,
	}
}

func (m *changes) elementMutations(index int) (*elementChanges, bool) {
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

func (m *changes) addAttribute(elemIndex int, key string, value string) {
	elemMutation, ok := m.elementMutations(elemIndex)
	newAttr := &newAttribute{
		key:   key,
		value: value,
	}

	if ok {
		elemMutation.newAttributes = append(elemMutation.newAttributes, newAttr)
		return
	}

	elemMutation = &elementChanges{
		index:         elemIndex,
		newAttributes: []*newAttribute{newAttr},
	}

	m.addElementChanges(elemMutation)
}

func (m *changes) addElementChanges(mutation *elementChanges) {
	if m.elementsIndex == nil {
		m.elements[mutation.index] = mutation
	} else {
		index, ok := m.elementsIndex[mutation.index]
		if ok {
			m.elements[index] = mutation
		} else {
			m.elementsIndex[mutation.index] = len(m.elements)
			m.elements = append(m.elements, mutation)
		}
	}
}

func (m *changes) setValue(elemIndex int, value string) {
	mutation, ok := m.elementMutations(elemIndex)
	if ok {
		mutation.value = value
		mutation.newElements = nil
		mutation.valueChanged = true
		return
	}

	mutation = &elementChanges{
		index:        elemIndex,
		value:        value,
		valueChanged: true,
	}

	m.addElementChanges(mutation)
}

func newMutations(schema *DOM) changes {
	var attributesIndex map[int]int
	var attributesMutations []*attributeChanges
	if schema.builder.attributeCounter < schema.attributesChangesSize {
		attributesMutations = make([]*attributeChanges, schema.builder.attributeCounter)
	} else {
		attributesIndex = map[int]int{}
	}

	var elementsMutations []*elementChanges
	var elementsMutationsIndex map[int]int
	if len(schema.elements) < schema.elementsChangesSize {
		elementsMutations = make([]*elementChanges, len(schema.elements))
	} else {
		elementsMutationsIndex = map[int]int{}
	}

	return changes{
		attributes:      attributesMutations,
		attributesIndex: attributesIndex,
		elements:        elementsMutations,
		elementsIndex:   elementsMutationsIndex,
	}
}
