package xml

type matcher struct {
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

	m.updateTrail()
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
	for ; m.indexes[m.index] < len(m.currRoot.children); m.indexes[m.index]++ {
		switch actual := m.selectors[m.index].(type) {
		case ElementSelector:
			if element := m.currRoot.children[m.indexes[m.index]]; element.Name.Local == string(actual) {
				return element, true
			}
		}
	}

	return nil, false
}

func (m *matcher) updateTrail() {
	for !(m.indexes[m.index] < len(m.currRoot.children)-1) && !(m.index == 0) {
		m.indexes[m.index] = 0
		m.index--
		m.indexes[m.index] += 1
	}
}

func newMatcher(root *StartElement, selectors []Selector) *matcher {
	return &matcher{
		indexes:   make([]int, len(selectors)),
		selectors: selectors,
		currRoot:  root,
	}
}
