package xml

import (
	"encoding/xml"
)

type builder struct {
	root         *StartElement
	elements     []*StartElement
	indexesStack []int
	counter      int
}

func (b *builder) addElement(actual xml.StartElement, offset int) {
	element := NewStartElement(&actual, b.counter, offset)
	b.addStartElement(element)
}

func (b *builder) addStartElement(element *StartElement) {
	b.appendElementIfNeeded(element)
	b.elements = append(b.elements, element)
	b.indexesStack = append(b.indexesStack, b.counter)

	b.counter++
}

func (b *builder) appendElementIfNeeded(element *StartElement) {
	if len(b.indexesStack) == 0 {
		return
	}

	parent := b.elements[b.indexesStack[len(b.indexesStack)-1]]
	parent.Append(element)
}

func (b *builder) closeElement() {
	b.indexesStack = b.indexesStack[:len(b.indexesStack)-1]
}

func (b *builder) addCharData(offset int) {
	currentElem := b.indexesStack[len(b.indexesStack)-1]
	element := b.elements[currentElem]
	element.end = offset

	if element.parent != nil {
		element.parent.end = offset
	}
}

func newBuilder() *builder {
	element := NewStartElement(nil, 0, 0)
	b := &builder{
		root: element,
	}

	b.addStartElement(element)
	return b
}

func (b *builder) allElements() []*StartElement {
	return b.elements[1:]
}
