package xml

import (
	"encoding/xml"
)

type builder struct {
	vxml         *VirtualXml
	root         *StartElement
	elements     []*StartElement
	indexesStack []int

	attributeCounter int
}

func (b *builder) addElement(actual xml.StartElement, valueStart int, raw []byte, offset int) error {
	attributesSpan, err := extractAttributes(offset, raw)
	if err != nil {
		return err
	}

	attributes := make([]*attribute, len(attributesSpan))
	for i := range attributesSpan {
		attributes[i] = newAttribute(attributesSpan[i], b.attributeCounter)
		b.attributeCounter++
	}

	element := NewStartElement(&actual, b.vxml, len(b.elements), valueStart, attributes)
	b.addStartElement(element)
	return nil
}

func (b *builder) addStartElement(element *StartElement) {
	b.appendElementIfNeeded(element)
	b.indexesStack = append(b.indexesStack, len(b.elements))
	b.elements = append(b.elements, element)
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

func newBuilder(vxml *VirtualXml) *builder {
	element := NewStartElement(nil, vxml, 0, 0, []*attribute{})
	b := &builder{
		root: element,
		vxml: vxml,
	}

	b.addStartElement(element)
	return b
}

func (b *builder) allElements() []*StartElement {
	return b.elements[1:]
}
