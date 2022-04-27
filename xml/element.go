package xml

//Element represents XML Element
type Element struct {
	document     *Document
	startElement *startElement
}

//Value returns Element value
func (e *Element) Value() string {
	return e.document.render(e.startElement.elemIndex, e.document.nextNotDescendant(e.startElement))
}

//Attribute returns Element attribute with given name
func (e *Element) Attribute(attribute string) (*Attribute, bool) {
	if e.document.wasReplaced(e.startElement) {
		return nil, false
	}

	attrIndex, ok := e.startElement.attrByName(attribute)
	if !ok {
		return nil, false
	}

	return &Attribute{
		document: e.document,
		element:  e,
		index:    attrIndex,
	}, true
}

//AddElement adds new Element value
func (e *Element) AddElement(value string) {
	if e.document.wasReplaced(e.startElement) {
		return
	}

	e.document.changes.addElement(e.startElement.elemIndex, value)
}

//AddAttribute adds new Attribute
func (e *Element) AddAttribute(key string, value string) {
	if e.document.wasReplaced(e.startElement) {
		return
	}

	e.document.changes.addAttribute(e.startElement.elemIndex, key, value)
}

//SetValue updates Element value
func (e *Element) SetValue(value string) {
	e.document.changes.setValue(e.startElement.elemIndex, value)
}

//InsertBefore inserts new element before receiver
func (e *Element) InsertBefore(element string) {
	if e.startElement.elemIndex == 0 {
		e.document.changes.addElement(e.startElement.parent.elemIndex, element)
		return
	}

	e.document.insertBefore(e.startElement.elemIndex, element)
}

//InsertAfter inserts new element after receiver
func (e *Element) InsertAfter(element string) {
	e.document.insertBefore(e.document.nextNotDescendant(e.startElement), element)
}

//SetAttribute update attribute value, or creates new one if it doesn't exist.
func (e *Element) SetAttribute(key, value string) {
	if e.document.wasReplaced(e.startElement) {
		return
	}

	attrIndex, ok := e.startElement.attrByName(key)
	if ok {
		e.document.updateAttribute(attrIndex, value)
		return
	}

	e.AddAttribute(key, value)
}

func (e *Element) ReplaceWith(newElement string) {
	e.document.replace(e.startElement, newElement)
}
