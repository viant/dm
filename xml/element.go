package xml

type Element struct {
	xml          *Xml
	startElement *StartElement
}

func (e *Element) Value() string {
	return string(e.xml.vXml.template[e.startElement.start:e.startElement.end])
}

func (e *Element) Attribute(attribute string) (string, bool) {
	return e.startElement.attrByName(attribute)
}
