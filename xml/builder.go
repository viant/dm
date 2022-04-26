package xml

import (
	"encoding/xml"
	"github.com/viant/dm/option"
)

type builder struct {
	schema       *VirtualDOM
	root         *startElement
	elements     []*startElement
	indexesStack []int

	attributeCounter int
	filters          *option.Filters
	skipped          int
}

func (b *builder) addElement(actual xml.StartElement, valueStart int, raw []byte, offset int) error {
	if b.skipped > 0 {
		b.skipped++
		return nil
	}

	var attributeFilter *option.Filter
	if b.filters != nil {
		var ok bool
		attributeFilter, ok = b.filters.ElementFilter(name(&actual))
		if !ok {
			b.skipped++
			return nil
		}
	}

	attributesSpan, err := extractAttributes(offset, raw)
	if err != nil {
		return err
	}

	attributes := make([]*attribute, len(attributesSpan))
	counter := 0
	for i, spans := range attributesSpan {
		if attributeFilter != nil && !attributeFilter.Matches(b.schema.templateSlice(spans[0])) {
			continue
		}

		attributes[i] = attributeOf(attributesSpan[i], b.attributeCounter)
		b.attributeCounter++
		counter++
	}

	attributes = attributes[:counter]
	element := newStartElement(&actual, b.schema, len(b.elements), valueStart, attributes, offset)
	b.addStartElement(element)
	return nil
}

func (b *builder) addStartElement(element *startElement) {
	b.appendElementIfNeeded(element)
	b.indexesStack = append(b.indexesStack, len(b.elements))
	b.elements = append(b.elements, element)
}

func (b *builder) appendElementIfNeeded(element *startElement) {
	if len(b.indexesStack) == 0 {
		return
	}

	parent := b.elements[b.indexesStack[len(b.indexesStack)-1]]
	parent.append(element)
}

func (b *builder) closeElement(offset int) {
	if b.skipped > 0 {
		b.skipped--
		return
	}

	b.elements[b.indexesStack[len(b.indexesStack)-1]].tag.end = offset
	b.indexesStack = b.indexesStack[:len(b.indexesStack)-1]
}

func (b *builder) addCharData(offset int, actual xml.CharData) {
	elemIndex := b.indexesStack[len(b.indexesStack)-1]
	element := b.elements[elemIndex]

	if element.indent == nil {
		element.indent = actual.Copy()
	}

	element.value.end = offset

	if element.parent != nil {
		element.parent.value.end = offset
	}
}

func newBuilder(vxml *VirtualDOM) *builder {
	element := newStartElement(nil, vxml, 0, 0, []*attribute{}, 0)
	b := &builder{
		root:   element,
		schema: vxml,
	}

	b.addStartElement(element)
	return b
}

func (b *builder) allElements() []*startElement {
	return b.elements[1:]
}

func name(element *xml.StartElement) string {
	result := element.Name.Local
	if element.Name.Space != "" {
		result = element.Name.Space + ":" + result
	}

	return result
}
