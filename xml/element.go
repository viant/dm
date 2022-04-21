package xml

//Element represents XML Element
type Element struct {
	xml          *Xml
	startElement *startElement
}

//Value returns Element value
func (e *Element) Value() string {
	return e.xml.render(e.startElement.elemIndex, e.xml.nextNotDescendant(e.startElement), true)
}

//Attribute returns Element attribute with given name
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

//AddElement adds new Element value
func (e *Element) AddElement(value string) {
	e.xml.mutations.addElement(e.startElement.elemIndex, value)
}

//AddAttribute adds new Attribute
func (e *Element) AddAttribute(key string, value string) {
	e.xml.mutations.addAttribute(e.startElement.elemIndex, key, value)
}

//SetValue updates Element value
func (e *Element) SetValue(value string) {
	e.xml.mutations.setValue(e.startElement.elemIndex, value)
}
