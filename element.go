package dm

import "bytes"

type (
	//Element represents HTML Element
	Element struct {
		template *Template
		tag      *tag
		attrs    attrs

		index           int
		attributeOffset int
	}

	//Attribute represents HTML Element attribute
	Attribute struct {
		template *Template
		index    int

		parent  int
		_parent *Element
	}
)

//InnerHTML returns InnerHTML of Element
func (e *Element) InnerHTML() []byte {
	return e.template.innerHTMLByIndex(e.tag.index)
}

//SetInnerHTML updates InnerHTML of Element
func (e *Element) SetInnerHTML(innerHTML []byte) error {
	return e.template.setInnerHTMLByIndex(e.index, innerHTML)
}

//AttributesLen returns number of Element Attributes
func (e *Element) AttributesLen() int {
	return len(e.attrs)
}

//AttributeByIndex returns n-th Attribute
func (e *Element) AttributeByIndex(i int) *Attribute {
	return &Attribute{
		template: e.template,
		index:    e.attributeOffset + i,
		parent:   e.index,
		_parent:  e,
	}
}

//Attribute returns Attribute that matches given Selectors
func (e *Element) Attribute(attrName, attrValue []byte) (*Attribute, bool) {
	for i := range e.attrs {
		if bytes.Equal(e.template.attributeKey(e.attributeOffset+i), attrName) && bytes.Equal(e.template.attributeValue(e.attributeOffset+i), attrValue) {
			return e.AttributeByIndex(i), true
		}
	}

	return nil, false
}

//Name returns Attribute Key
func (a *Attribute) Name() []byte {
	return a.template.attributeKey(a.index)
}

//Value returns Attribute Value
func (a *Attribute) Value() []byte {
	return a.template.attributeValue(a.index)
}

//Set updates Attribute value
func (a *Attribute) Set(newValue []byte) {
	a.template.setAttributeByIndex(a.index, newValue)
}

//Parent returns Attribute parent Element
func (a *Attribute) Parent() *Element {
	if a._parent != nil {
		return a._parent
	}

	element := &Element{
		template:        a.template,
		tag:             a.template.tag(a.parent),
		attrs:           a.template.tagAttributes(a.parent),
		index:           a.parent,
		attributeOffset: a.template.tag(a.parent - 1).attrEnd,
	}
	a._parent = element
	return element
}
