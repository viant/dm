package xml

type matcher struct {
	document *Document
	indexes  []int
	maxSize  []int
	index    int

	selectors []Selector
	currRoots []*startElement

	next *startElement
}

func (m *matcher) updateIndexes() {
	if m.index == len(m.selectors) {
		if m.index == 0 {
			return
		}

		m.index--
		m.indexes[m.index]++
		m.currRoots[0] = m.currRoots[0].parent
	}

	for !(m.indexes[m.index] < m.maxSize[m.index]) && m.index != 0 {
		m.indexes[m.index] = 0
		m.index--
		m.indexes[m.index] += 1
		m.currRoots[0] = m.currRoots[0].parent
	}
}

func (m *matcher) match() int {
	if len(m.currRoots) == 0 {
		return -1
	}

	m.updateIndexes()

	for m.index < len(m.selectors) {
		if elem, ok := m.checkIfMatches(); ok {
			m.index++
			m.currRoots[0] = elem
			continue
		}

		if m.index == 0 {
			if len(m.currRoots) <= 1 {
				return -1
			} else {
				m.currRoots = m.currRoots[1:]
				continue
			}
		}

		m.updateIndexes()
	}

	m.next = m.currRoots[0]
	m.currRoots = m.currRoots[1:]

	return m.next.elemIndex
}

func (m *matcher) checkIfMatches() (*startElement, bool) {
	if len(m.currRoots[0].children) > mapSize {
		return m.findElement()
	}

	m.maxSize[m.index] = len(m.currRoots[0].children)
	for ; m.indexes[m.index] < len(m.currRoots[0].children); m.indexes[m.index]++ {
		element := m.currRoots[0].children[m.indexes[m.index]]
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
		attrIndex, ok := elem.attrByName(attributeSelector.Name)
		if !ok || !m.checkAttributeValue(elem, attrIndex, &attributeSelector) {
			return false
		}
	}

	return true
}

func (m *matcher) checkAttributeValue(element *startElement, attr int, attrSelector *AttributeSelector) bool {
	value, ok := m.document.checkAttributeChanges(attr)
	if ok {
		return m.compare(value, attrSelector)
	}

	return m.compare(m.document.templateSlice(element.attributeValueSpan(attr)), attrSelector)
}

func (m *matcher) findElement() (*startElement, bool) {
	elementsIndexes, ok := m.currRoots[0].elementsIndex[m.selectors[m.index].Name]
	m.maxSize[m.index] = len(elementsIndexes)
	if !ok || m.sliceIndex() >= len(elementsIndexes) {
		return nil, false
	}

	index := m.sliceIndex()
	element := m.currRoots[0].children[elementsIndexes[index]]
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

func newMatcher(document *Document, selectors []Selector) *matcher {
	roots := []*startElement{document.dom.root}
	if len(selectors) > 0 && selectors[0].MatchAny {
		ints := document.dom.groups[selectors[0].Name]
		roots = make([]*startElement, len(ints))

		for i, elemIndex := range ints {
			roots[i] = document.dom.elements[elemIndex]
		}

		selectors = selectors[1:]
	}

	return &matcher{
		indexes:   make([]int, len(selectors)),
		maxSize:   make([]int, len(selectors)),
		selectors: selectors,
		currRoots: roots,
		document:  document,
	}
}
