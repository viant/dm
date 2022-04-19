package xml

import "fmt"

type Iterator struct {
	xml     *Xml
	current int
	next    int
	matcher *matcher
}

func (i *Iterator) Has() bool {
	if i.current < i.next {
		return true
	}

	if i.next < i.current {
		return false
	}

	i.next = i.matcher.match()

	return i.next != -1
}

func (i *Iterator) Next() (*Element, error) {
	if i.next == i.current {
		return nil, fmt.Errorf("it is needed to call Has, before Next is called")
	}

	if i.next == -1 {
		return nil, fmt.Errorf("next was called but there are no elements left")
	}

	i.current = i.next
	return &Element{
		startElement: i.matcher.currRoot,
		xml:          i.xml,
	}, nil
}

func NewIterator(xml *Xml, selectors []Selector) *Iterator {
	return &Iterator{
		xml:     xml,
		current: -1,
		next:    -1,
		matcher: newMatcher(xml.vXml.root, selectors),
	}
}
