package xml

import (
	"encoding/xml"
	"github.com/viant/dm/option"
)

type builder struct {
	dom          *DOM
	root         *startElement
	elements     []*startElement
	indexesStack []int

	attributeCounter int
	filters          *option.Filters

	groups  map[string][]int
	skipped int
}

func (b *builder) addElement(actual xml.StartElement, valueStart int, raw []byte, offset int) error {
	if b.skipped > 0 {
		b.skipped++
		return nil
	}

	var attributeFilter *option.Filter
	if b.filters != nil {
		var ok bool
		attributeFilter, ok = b.filters.ElementFilter(elementFullName(&actual), true)
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
		if attributeFilter != nil && !attributeFilter.Matches(b.dom.templateSlice(spans[0])) {
			continue
		}

		attributes[counter] = attributeOf(attributesSpan[i], b.attributeCounter)
		b.attributeCounter++
		counter++
	}

	attributes = attributes[:counter]
	element := newStartElement(&actual, b.dom, len(b.elements), valueStart, attributes, offset)
	b.addStartElement(element)
	return nil
}

func (b *builder) addStartElement(element *startElement) {
	b.groups[element.name] = append(b.groups[element.name], len(b.elements))

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

func newBuilder(dom *DOM) *builder {
	element := newStartElement(nil, dom, 0, 0, []*attribute{}, 0)
	b := &builder{
		root:   element,
		dom:    dom,
		groups: map[string][]int{},
	}

	b.addStartElement(element)
	return b
}

func (b *builder) allElements() []*startElement {
	return b.elements[1:]
}

func elementFullName(element *xml.StartElement) string {
	result := element.Name.Local
	if element.Name.Space != "" {
		result = element.Name.Space + ":" + result
	}

	return result
}
