package dm

import "fmt"

type (
	iterator struct {
		template  *Template
		current   int
		next      int
		selectors [][]byte
	}

	//TagIterator iterates over matching tags
	TagIterator struct {
		iterator
	}

	//AttributeIterator iterates over matching attributes
	AttributeIterator struct {
		iterator
	}
)

//HasMore returns true if there are more tags matching given selectors
func (it *TagIterator) HasMore() bool {
	if it.current < it.next {
		return true
	}

	if it.next == -1 {
		return false
	}

	for i := it.current + 1; i < it.template.tagLen(); i++ {
		if it.template.matchTag(i, it.selectors) {
			it.next = i
			return true
		}
	}

	it.next = -1
	return false
}

//Next returns Element matching given selectors
//Note: it is needed to call HasMore before calling Next
func (it *TagIterator) Next() (*Element, error) {
	if it.current == it.next {
		return nil, fmt.Errorf("it is needed to call HasMore, before Next is called")
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

//HasMore returns true if there are more attributes matching given selectors
func (at *AttributeIterator) HasMore() bool {
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
//Note: it is needed to call HasMore before calling Next
func (at *AttributeIterator) Next() (*Attribute, error) {
	if at.current == at.next {
		return nil, fmt.Errorf("it is needed to call HasMore, before Next is called")
	}

	if at.next == -1 {
		return nil, fmt.Errorf("next was called but there are no elements left")
	}

	at.current = at.next
	return &Attribute{
		template: at.template,
		index:    at.current,
		parent:   at.template.attribute(at.current).tag,
	}, nil
}
