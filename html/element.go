package html

import "bytes"

type (
	//Element represents HTML Element
	Element struct {
		template *Document
		tag      *tag
	}

	//Attribute represents HTML Element attribute
	Attribute struct {
		template *Document
		attr     *attr
		_parent  *Element
	}
)

//InnerHTML returns InnerHTML of Element
func (e *Element) InnerHTML() string {
	return string(e.template.innerHTML(e.tag))
}

//SetInnerHTML updates InnerHTML of Element
func (e *Element) SetInnerHTML(innerHTML string) error {
	return e.template.setInnerHTML(e.tag, asBytes(innerHTML))
}

//AttributesLen returns number of Element Attributes
func (e *Element) AttributesLen() int {
	return len(e.tag.attrs)
}

//AttributeByIndex returns n-th Attribute
func (e *Element) AttributeByIndex(i int) *Attribute {
	return &Attribute{
		template: e.template,
		attr:     e.tag.attrs[i],
		_parent:  e,
	}
}

func (e *Element) HasAttribute(name string) bool {
	for i := range e.tag.attrs {
		if bytes.EqualFold(e.template.attributeKey(e.tag.attrs[i]), asBytes(name)) {
			return true
		}
	}
	return false
}

//Attribute returns matched attribute, true or nil, false
func (e *Element) Attribute(name string) (*Attribute, bool) {
	for i := range e.tag.attrs {
		if bytes.EqualFold(e.template.attributeKey(e.tag.attrs[i]), asBytes(name)) {
			return e.AttributeByIndex(i), true
		}
	}
	return nil, false
}

//MatchAttribute returns an attribute that matches the supplied attribute name and value
func (e *Element) MatchAttribute(name, value string) (*Attribute, bool) {
	for i := range e.tag.attrs {
		if bytes.EqualFold(e.template.attributeKey(e.tag.attrs[i]), asBytes(name)) && bytes.EqualFold(e.template.attributeValue(e.tag.attrs[i]), asBytes(value)) {
			return e.AttributeByIndex(i), true
		}
	}
	return nil, false
}

func (e *Element) AddAttribute(key, value string) {
	e.template.addAttribute(e.tag, key, value)
}

//Name returns Attribute Key
func (a *Attribute) Name() string {
	return string(a.template.dom.template[a.attr.keyStart():a.attr.keyEnd()])
}

//Value returns Attribute Value
func (a *Attribute) Value() string {
	return string(a.template.attributeValue(a.attr))
}

//Set updates Attribute value
func (a *Attribute) Set(newValue string) {
	a.template.setAttribute(a.attr, asBytes(newValue))
}

//Parent returns Attribute parent Element
func (a *Attribute) Parent() *Element {
	if a._parent != nil {
		return a._parent
	}

	element := &Element{
		template: a.template,
		tag:      a.attr.tag,
	}

	a._parent = element
	return element
}
