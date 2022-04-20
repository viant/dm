package xml

type Element struct {
	xml          *Xml
	startElement *StartElement
}

func (e *Element) Value() string {
	return string(e.xml.vXml.template[e.startElement.start:e.startElement.end])
}

func (e *Element) Attribute(attribute string) (*Attribute, bool) {
	index, ok := e.startElement.attrByName(attribute)
	if !ok {
		return nil, false
	}

	return &Attribute{
		xml:     e.xml,
		element: e,
		index:   index,
	}, true
}

func (e *Element) AddElement(value string) {
	e.xml.mutations.appendElement(e.startElement.elemIndex, value)
}
