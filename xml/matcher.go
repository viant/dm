package xml

type matcher struct {
	xml     *DOM
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
	if len(m.currRoot.children) > mapSize {
		return m.findElement()
	}

	m.maxSize[m.index] = len(m.currRoot.children)
	for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
		element := m.currRoot.children[m.indexes[m.index]]
		if !m.matches(element) {
			continue
		}

		return element, true
	}

	return nil, false
}

func (m *matcher) matches(elem *startElement) bool {
	if m.selectors[m.index].Name != elem.name {
		return false
	}

	for _, attributeSelector := range m.selectors[m.index].Attributes {
		byName, ok := elem.attrByName(attributeSelector.Name)
		if !ok || !m.checkAttributeValue(elem, byName, &attributeSelector) {
			return false
		}
	}

	return true
}

func (m *matcher) checkAttributeValue(element *startElement, attr int, attrSelector *AttributeSelector) bool {
	value, ok := m.xml.checkAttributeChanges(attr)
	if ok {
		return m.compare(value, attrSelector)
	}

	return m.compare(m.xml.templateSlice(element.attributeValueSpan(attr)), attrSelector)
}

func (m *matcher) findElement() (*startElement, bool) {
	elementsIndexes, ok := m.currRoot.elementsIndex[m.selectors[m.index].Name]
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

func (m *matcher) compare(currentValue string, attrSelector *AttributeSelector) bool {
	switch attrSelector.Compare {
	case EQ:
		return currentValue == attrSelector.Value
	case NEQ:
		return currentValue != attrSelector.Value
	default:
		return currentValue == attrSelector.Value
	}
}

func newMatcher(dom *DOM, selectors []Selector) *matcher {
	return &matcher{
		indexes:   make([]int, len(selectors)),
		maxSize:   make([]int, len(selectors)),
		selectors: selectors,
		currRoot:  dom.vdom.root,
		xml:       dom,
	}
}
