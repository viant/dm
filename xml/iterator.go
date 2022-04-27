package xml

import "fmt"

//Iterator iterates over matching Elements
type Iterator struct {
	document *Document
	current  int
	next     int
	matcher  *matcher
}

//Has returns true if there are more matching elements
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

//Next returns next matching Element
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
		document:     i.document,
	}, nil
}

func newIterator(xml *Document, selectors []Selector) *Iterator {
	return &Iterator{
		document: xml,
		current:  -1,
		next:     -1,
		matcher:  newMatcher(xml, selectors),
	}
}
