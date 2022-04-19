package xml

type Attribute struct {
	xml     *Xml
	element *Element
	index   int
}

func (a *Attribute) Value() string {
	return a.element.startElement.Attr[a.index].Value
}

func (a *Attribute) Set(value string) {
	//a.xml.
}
