package xml

//Attribute represents XML element attribute
type Attribute struct {
	xml     *Xml
	element *Element
	index   int
}

//Value returns attribute value
func (a *Attribute) Value() string {
	change, ok := a.xml.mutations.checkAttributeChanges(a.index)
	if ok {
		return change
	} else {
		return a.xml.templateSlice(a.element.startElement.attributeValueSpan(a.index))
	}
}

//Set updates attribute value
func (a *Attribute) Set(value string) {
	a.xml.mutations.updateAttribute(a.index, value)
}
