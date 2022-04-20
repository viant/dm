package xml

type Attribute struct {
	xml     *Xml
	element *Element
	index   int
}

func (a *Attribute) Value() string {
	value, ok := a.xml.mutations.attributeValue(a.index)
	if ok {
		return value
	} else {
		return a.element.startElement.Attr[a.index].Value
	}
}

func (a *Attribute) Set(value string) {
	a.xml.mutations.updateAttribute(a.index, value)
}
