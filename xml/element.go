package xml

//Element represents XML Element
type Element struct {
	dom          *DOM
	startElement *startElement
}

//Value returns Element value
func (e *Element) Value() string {
	return e.dom.render(e.startElement.elemIndex, e.dom.nextNotDescendant(e.startElement))
}

//Attribute returns Element attribute with given name
func (e *Element) Attribute(attribute string) (*Attribute, bool) {
	index, ok := e.startElement.attrByName(attribute)
	if !ok {
		return nil, false
	}

	return &Attribute{
		xml:     e.dom,
		element: e,
		index:   index,
	}, true
}

//AddElement adds new Element value
func (e *Element) AddElement(value string) {
	e.dom.changes.addElement(e.startElement.elemIndex, value)
}

//AddAttribute adds new Attribute
func (e *Element) AddAttribute(key string, value string) {
	e.dom.changes.addAttribute(e.startElement.elemIndex, key, value)
}

//SetValue updates Element value
func (e *Element) SetValue(value string) {
	e.dom.changes.setValue(e.startElement.elemIndex, value)
}

func (e *Element) InsertBefore(element string) {
	if e.startElement.elemIndex == 0 {
		e.dom.changes.addElement(e.startElement.parent.elemIndex, element)
		return
	}

	e.dom.insertBefore(e.startElement.elemIndex, element)
}

func (e *Element) InsertAfter(element string) {
	e.dom.insertBefore(e.dom.nextNotDescendant(e.startElement), element)
}
