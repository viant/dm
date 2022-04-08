package dm

import "fmt"

type (
	iterator struct {
		template  *DOM
		current   int
		next      int
		selectors []string
	}

	//ElementIterator iterates over matching tags
	ElementIterator struct {
		iterator
	}

	//AttributeIterator iterates over matching attributes
	AttributeIterator struct {
		iterator
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

	if it.current == -1 {
		it.next = it.template.nextMatchingTag(0, it.selectors)
	} else {
		it.next = it.template.nextMatchingTag(it.current+1, it.selectors)
	}

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
		template:        it.template,
		tag:             it.template.tag(it.current),
		attributeOffset: it.template.tag(it.current - 1).attrEnd,
		attrs:           it.template.tagAttributes(it.current),
		index:           it.current,
	}, nil
}

//Has returns true if there are more attributes matching given selectors
func (at *AttributeIterator) Has() bool {
	if at.current < at.next {
		return true
	}

	if at.next == -1 {
		return false
	}

	at.next = at.template.nextAttribute(at.current, at.selectors...)
	return at.next != -1
}

//Next returns Attribute matching given selectors
//Note: it is needed to call Has before calling Next
func (at *AttributeIterator) Next() (*Attribute, error) {
	if at.current == at.next {
		return nil, fmt.Errorf("it is needed to call Has, before Next is called")
	}

	if at.next == -1 {
		return nil, fmt.Errorf("next was called but there are no elements left")
	}

	result := &Attribute{
		template: at.template,
		index:    at.next,
		parent:   at.template.attribute(at.current).tag,
	}

	at.current = at.template.tag(at.template.attribute(at.current).tag).attrEnd
	at.next = at.current
	return result, nil
}