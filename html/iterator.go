package html

import "fmt"

type (
	iterator struct {
		template  *Document
		current   int
		next      int
		selectors []string
	}

	//ElementIterator iterates over matching tags
	ElementIterator struct {
		iterator
		matcher *elementMatcher
		index   int
	}

	//AttributeIterator iterates over matching attributes
	AttributeIterator struct {
		nextAttr    *attr
		currentAttr *attr
		matcher     *attributeMatcher
		dom         *Document
	}
)

//Has returns true if there are more tags matching given selectors
func (it *ElementIterator) Has() bool {
	if it.current < it.next {
		return true
	}

	if it.next == -1 && it.current != -1 {
		return false
	}

	it.next = it.matcher.match()

	return it.next != -1
}

//Next returns Element matching given selectors
//Note: it is needed to call Has before calling Next
func (it *ElementIterator) Next() (*Element, error) {
	if it.current == it.next {
		return nil, fmt.Errorf("it is needed to call Has, before Next is called")
	}

	if it.next == -1 {
		return nil, fmt.Errorf("next was called but there are no elements left")
	}

	it.current = it.next
	return &Element{
		template: it.template,
		tag:      it.template.tag(it.current),
	}, nil
}

//Has returns true if there are more attributes matching given selectors
func (at *AttributeIterator) Has() bool {
	if at.nextAttr == nil && at.currentAttr != nil {
		return false
	}

	if at.nextAttr != nil && at.currentAttr == nil {
		return true
	}

	at.nextAttr = at.matcher.match()
	at.currentAttr = nil

	return at.nextAttr != nil
}

//Next returns Attribute matching given selectors
//Note: it is needed to call Has before calling Next
func (at *AttributeIterator) Next() (*Attribute, error) {
	if at.nextAttr == nil {
		return nil, fmt.Errorf("it is needed to call Has, before Next is called")
	}

	result := &Attribute{
		template: at.dom,
		attr:     at.nextAttr,
	}

	at.currentAttr = at.nextAttr
	return result, nil
}
