package xml

type matcher struct {
	xml     *Xml
	indexes []int
	maxSize []int
	index   int

	selectors []Selector
	currRoot  *startElement
}

func (m *matcher) updateIndexes() {
	if m.index == len(m.selectors) {
		m.index--
		m.indexes[m.index]++
		m.currRoot = m.currRoot.parent
	}

	for !(m.indexes[m.index] < m.maxSize[m.index]) && m.index != 0 {
		m.indexes[m.index] = 0
		m.index--
		m.indexes[m.index] += 1
		m.currRoot = m.currRoot.parent
	}
}

func (m *matcher) match() int {
	m.updateIndexes()

	for m.index < len(m.selectors) {
		if elem, ok := m.matchAny(); ok {
			m.index++
			m.currRoot = elem
			continue
		}

		if m.index == 0 {
			return -1
		}

		m.updateIndexes()
	}

	return m.currRoot.elemIndex
}

func (m *matcher) matchAny() (*startElement, bool) {
	switch actual := m.selectors[m.index].(type) {
	case ElementSelector:
		if len(m.currRoot.children) > mapSize {
			return m.findElement(string(actual))
		}

		m.maxSize[m.index] = len(m.currRoot.children)
		for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
			if element := m.currRoot.children[m.indexes[m.index]]; element.name == string(actual) {
				return element, true
			}
		}

	case AttributeSelector:
		if m.currRoot.childrenAttrSize > mapSize {
			return m.findAttribute(&actual)
		}

		m.maxSize[m.index] = len(m.currRoot.children)
		for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
			startElem := m.currRoot.children[m.indexes[m.index]]
			attrIndex, ok := startElem.attrByName(actual.Name)
			if ok && m.checkAttributeValue(startElem, attrIndex, actual.Value) {
				return startElem, true
			}
		}
	}

	return nil, false
}

func (m *matcher) checkAttributeValue(element *startElement, attr int, attrValue string) bool {
	value, ok := m.xml.checkAttributeChanges(attr)
	if ok {
		return value == attrValue
	}

	return m.xml.templateSlice(element.attributeValueSpan(attr)) == attrValue
}

func (m *matcher) findElement(elementName string) (*startElement, bool) {
	elementsIndexes, ok := m.currRoot.elementsIndex[elementName]
	m.maxSize[m.index] = len(elementsIndexes)
	if !ok || m.sliceIndex() >= len(elementsIndexes) {
		return nil, false
	}

	index := m.sliceIndex()
	element := m.currRoot.children[elementsIndexes[index]]
	return element, true
}

func (m *matcher) sliceIndex() int {
	return m.indexes[m.index]
}

func (m *matcher) findAttribute(a *AttributeSelector) (*startElement, bool) {
	indexes, ok := m.currRoot.attributeChildrenIndex[a.Name]
	m.maxSize[m.index] = len(indexes)

	if !ok {
		return nil, false
	}

	for ; m.indexes[m.index] < len(indexes); m.indexes[m.index]++ {
		index := indexes[m.sliceIndex()]
		attributeOwner := m.currRoot.children[index]
		attrIndex, ok := attributeOwner.attrByName(a.Name)
		if !ok {
			continue
		}

		changed, ok := m.xml.checkAttributeChanges(attributeOwner.attributes[attrIndex].index)
		if ok {
			if changed == a.Value {
				return attributeOwner, true
			}
		} else {
			if m.xml.templateSlice(attributeOwner.attributes[attrIndex].spans[1]) == a.Value {
				return m.currRoot.children[indexes[m.sliceIndex()]], true
			}
		}
	}

	return nil, false
}

func newMatcher(xml *Xml, selectors []Selector) *matcher {
	return &matcher{
		indexes:   make([]int, len(selectors)),
		maxSize:   make([]int, len(selectors)),
		selectors: selectors,
		currRoot:  xml.schema.root,
		xml:       xml,
	}
}
