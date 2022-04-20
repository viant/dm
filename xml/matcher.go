package xml

type matcher struct {
	xml     *Xml
	indexes []int
	index   int

	selectors []Selector
	currRoot  *StartElement
}

func (m *matcher) updateIndex() {
	if m.index == len(m.selectors) {
		m.index--
		m.indexes[m.index]++
	}

	for !(m.indexes[m.index] < len(m.currRoot.children)-1) && !(m.index == 0) {
		m.indexes[m.index] = 0
		m.index--
		m.indexes[m.index] += 1
	}
}

func (m *matcher) match() int {
	m.updateIndex()

	for m.index < len(m.selectors) {
		if elem, ok := m.matchAny(); ok {
			m.index++
			m.currRoot = elem
			continue
		}

		if m.index == 0 {
			return -1
		}

		m.updateIndex()
	}

	return m.currRoot.elemIndex
}

func (m *matcher) matchAny() (*StartElement, bool) {
	switch actual := m.selectors[m.index].(type) {
	case ElementSelector:
		for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
			if element := m.currRoot.children[m.indexes[m.index]]; element.Name.Local == string(actual) {
				return element, true
			}
		}
	case AttributeSelector:
		for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
			startElement := m.currRoot.children[m.indexes[m.index]]
			attrIndex, ok := startElement.attrByName(actual.Name)
			if ok && m.checkAttributeValue(startElement, attrIndex, actual.Value) {
				return startElement, true
			}
		}
	}

	return nil, false
}

func (m *matcher) checkAttributeValue(element *StartElement, attr int, attrValue string) bool {
	value, ok := m.xml.attributeValue(attr)
	if ok {
		return value == attrValue
	}

	return element.Attr[attr].Value == attrValue
}

func newMatcher(xml *Xml, selectors []Selector) *matcher {
	return &matcher{
		indexes:   make([]int, len(selectors)),
		selectors: selectors,
		currRoot:  xml.vXml.root,
		xml:       xml,
	}
}
