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
	e.xml.mutations.addElement(e.startElement.elemIndex, value)
}

func (e *Element) AddAttribute(key string, value string) {
	e.xml.mutations.addAttribute(e.startElement.elemIndex, key, value)
}

func (e *Element) Set(value string) {
	e.xml.mutations.updateValue(e.startElement.elemIndex, value)
}
