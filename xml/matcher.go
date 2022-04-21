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
	if !(len(m.currRoot.children) > mapSize) { //maxSize will be updated in the m.findByMap
		m.maxSize[m.index] = len(m.currRoot.children)
	}

	switch actual := m.selectors[m.index].(type) {
	case ElementSelector:
		if len(m.currRoot.children) > mapSize {
			return m.findByMap(string(actual))
		}

		for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
			if element := m.currRoot.children[m.indexes[m.index]]; element.Name.Local == string(actual) {
				return element, true
			}
		}

	case AttributeSelector:
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
	value, ok := m.xml.attributeValue(attr)
	if ok {
		return value == attrValue
	}

	return element.Attr[attr].Value == attrValue
}

func (m *matcher) findByMap(elementName string) (*startElement, bool) {
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

func newMatcher(xml *Xml, selectors []Selector) *matcher {
	return &matcher{
		indexes:   make([]int, len(selectors)),
		maxSize:   make([]int, len(selectors)),
		selectors: selectors,
		currRoot:  xml.vXml.root,
		xml:       xml,
	}
}
