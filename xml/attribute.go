package xml

//Attribute represents XML element attribute
type Attribute struct {
	document *Document
	element  *Element
	index    int
}

//Value returns attribute value
func (a *Attribute) Value() string {
	change, ok := a.document.changes.checkAttributeChanges(a.index)
	if ok {
		return change
	} else {
		return a.document.templateSlice(a.element.startElement.attributeValueSpan(a.index))
	}
}

//Set updates attribute value
func (a *Attribute) Set(value string) {
	a.document.changes.updateAttribute(a.index, value)
}
