package xml

//Attribute represents XML element attribute
type Attribute struct {
	xml     *DOM
	element *Element
	index   int
}

//Value returns attribute value
func (a *Attribute) Value() string {
	change, ok := a.xml.changes.checkAttributeChanges(a.index)
	if ok {
		return change
	} else {
		return a.xml.templateSlice(a.element.startElement.attributeValueSpan(a.index))
	}
}

//Set updates attribute value
func (a *Attribute) Set(value string) {
	a.xml.changes.updateAttribute(a.index, value)
}
