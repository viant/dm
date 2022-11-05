package html

import "bytes"

type (
	//Element represents HTML Element
	Element struct {
		template *DOM
		tag      *tag
		attrs    attrs

		index           int32
		attributeOffset int32
	}

	//Attribute represents HTML Element attribute
	Attribute struct {
		template *DOM
		index    int32

		parent  int32
		_parent *Element
	}
)

// InnerHTML returns InnerHTML of Element
func (e *Element) InnerHTML() string {
	return string(e.template.innerHTMLByIndex(e.tag.index)) //not converting unsafe.Pointers to not make mutable string, if you change source slice, string will also change
}

// SetInnerHTML updates InnerHTML of Element
func (e *Element) SetInnerHTML(innerHTML string) error {
	return e.template.setInnerHTMLByIndex(e.index, asBytes(innerHTML))
}

// AttributesLen returns number of Element Attributes
func (e *Element) AttributesLen() int32 {
	return int32(len(e.attrs))
}

// AttributeByIndex returns n-th Attribute
func (e *Element) AttributeByIndex(i int32) *Attribute {
	return &Attribute{
		template: e.template,
		index:    e.attributeOffset + i,
		parent:   e.index,
		_parent:  e,
	}
}

func (e *Element) HasAttribute(name string) bool {
	for i := range e.attrs {
		if bytes.EqualFold(e.template.attributeKey(e.attributeOffset+int32(i)), asBytes(name)) {
			return true
		}
	}
	return false
}

// Attribute returns matched attribute, true or nil, false
func (e *Element) Attribute(name string) (*Attribute, bool) {
	for i := range e.attrs {
		if bytes.EqualFold(e.template.attributeKey(e.attributeOffset+int32(i)), asBytes(name)) {
			return e.AttributeByIndex(int32(i)), true
		}
	}
	return nil, false
}

// MatchAttribute returns an attribute that matches the supplied attribute name and value
func (e *Element) MatchAttribute(name, value string) (*Attribute, bool) {
	for i := range e.attrs {
		if bytes.EqualFold(e.template.attributeKey(e.attributeOffset+int32(i)), asBytes(name)) && bytes.EqualFold(e.template.attributeValue(e.attributeOffset+int32(i)), asBytes(value)) {
			return e.AttributeByIndex(int32(i)), true
		}
	}
	return nil, false
}

func (e *Element) AddAttribute(key, value string) {
	e.template.addAttribute(e.index, key, value)
}

// Name returns Attribute Key
func (a *Attribute) Name() string {
	return string(a.template.attributeKey(a.index)) //not converting unsafe.Pointers to not make mutable string, if you change source slice, string will also change
}

// Value returns Attribute Value
func (a *Attribute) Value() string {
	return string(a.template.attributeValue(a.index)) //not converting unsafe.Pointers to not make mutable string, if you change source slice, string will also change
}

// Set updates Attribute value
func (a *Attribute) Set(newValue string) {
	a.template.setAttributeByIndex(a.index, asBytes(newValue))
}

// Parent returns Attribute parent Element
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
