package xml

//Attribute represents XML element attribute
type Attribute struct {
	xml     *Xml
	element *Element
	index   int
}

//Value returns attribute value
func (a *Attribute) Value() string {
	value, ok := a.xml.mutations.attributeValue(a.index)
	if ok {
		return value
	} else {
		return a.element.startElement.Attr[a.index].Value
	}
}

//Set updates attribute value
func (a *Attribute) Set(value string) {
	a.xml.mutations.updateAttribute(a.index, value)
}
