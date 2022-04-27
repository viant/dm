package html

type (
	elementMatcher struct {
		dom   *Document
		index int

		selectors  []string
		groupIndex int

		offset  int
		nextTag int
	}

	attributeMatcher struct {
		elemMatcher *elementMatcher

		attributeOffset int
		selectors       []string
		currTag         *tag
	}
)

func (m *elementMatcher) match() int {
	m.nextTag = m.matchNextTag()
	if m.nextTag == -1 {
		return -1
	}

	m.offset++

	return m.currentGroup()[m.nextTag]
}

func (m *elementMatcher) matchNextTag() int {
	if m.groupIndex == -1 || m.offset >= len(m.dom.dom.tagsGrouped[m.groupIndex]) {
		return -1
	}

	for ; m.offset < len(m.dom.dom.tagsGrouped[m.groupIndex]); m.offset++ {
		tagIndex := m.currentGroup()[m.offset]
		if m.dom.tagsRemoved[tagIndex] {
			continue
		}

		if m.matchTagAttributes(tagIndex) {
			return m.offset
		}
	}

	return -1
}

func (m *elementMatcher) currentGroup() []int {
	return m.dom.dom.tagsGrouped[m.groupIndex]
}

func newElementMatcher(d *Document, selectors []string) *elementMatcher {
	m := &elementMatcher{
		dom:       d,
		selectors: selectors,
	}

	m.init()
	return m
}

func (m *elementMatcher) init() {
	if len(m.selectors) == 0 {
		return
	}

	m.groupIndex = m.dom.dom.index.tagIndex(m.selectors[0], false)
}

func (m *elementMatcher) matchTagAttributes(tagIndex int) bool {
	if len(m.selectors) < 2 {
		return true
	}

	attrByName, ok := m.dom.tag(tagIndex).attributeByName(m.selectors[1])
	if !ok {
		return false
	}

	return m.dom.matchAttributeValue(attrByName, m.selectors)
}

func newAttributeMatcher(dom *Document, selectors []string) *attributeMatcher {
	elemNameSelector := selectors
	if len(selectors) > 1 {
		elemNameSelector = selectors[:1]
	}

	attrMatcher := &attributeMatcher{
		elemMatcher:     newElementMatcher(dom, elemNameSelector),
		attributeOffset: 0,
		selectors:       selectors,
	}

	attrMatcher.init()
	return attrMatcher
}

func (a *attributeMatcher) match() *attr {
	for {
		if a.currTag == nil {
			return nil
		}

		if len(a.selectors) <= 1 {
			if a.attributeOffset < len(a.currTag.attrs)-1 {
				anAttr := a.currTag.attrs[a.attributeOffset]
				a.attributeOffset++
				return anAttr
			} else {
				a.attributeOffset = 0
				a.matchNextElem()
				continue
			}
		}

		attrByName, ok := a.currTag.attributeByName(a.selectors[1])
		ok = ok && a.elemMatcher.dom.matchAttributeValue(attrByName, a.selectors)

		a.matchNextElem()
		if ok {
			return attrByName
		}
	}
}

func (a *attributeMatcher) init() {
	a.matchNextElem()
}

func (a *attributeMatcher) matchNextElem() {
	tagIndex := a.elemMatcher.match()
	a.attributeOffset = 0

	if tagIndex == -1 {
		a.currTag = nil
		return
	}

	a.currTag = a.elemMatcher.dom.tag(tagIndex)
}
