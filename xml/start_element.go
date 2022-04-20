package xml

import "encoding/xml"

type (
	startElement struct {
		*xml.StartElement
		span

		parent      *startElement
		children    []*startElement
		parentIndex int
		elemIndex   int
		nextSibling int
		vxml        *Schema
		indent      []byte

		attributeIndex map[string]int
		attributesName []string
		attributes     []*attribute
	}

	attribute struct {
		spans [2]span
		index int
	}
)

func (a *attribute) valueStart() int {
	return a.spans[1].start
}

func (a *attribute) valueEnd() int {
	return a.spans[1].end
}

func (a *attribute) keyStart() int {
	return a.spans[0].start
}

func (a *attribute) keyEnd() int {
	return a.spans[0].end
}

func (s *startElement) append(child *startElement) {
	child.parentIndex = len(s.children)
	if len(s.children) > 0 {
		s.children[len(s.children)-1].nextSibling = child.elemIndex
	}

	s.children = append(s.children, child)
	child.parent = s
}

func (s *startElement) attrByName(attribute string) (int, bool) {
	if s.attributeIndex != nil {
		value, ok := s.attributeIndex[attribute]
		return value, ok
	}

	for i, attr := range s.attributesName {
		if attr == attribute {
			return i, true
		}
	}

	return -1, false
}

func newStartElement(element *xml.StartElement, virtualXml *Schema, elemIndex int, startPosition int, attributes []*attribute) *startElement {
	elem := &startElement{
		StartElement: element,
		elemIndex:    elemIndex,
		span: span{
			start: startPosition,
		},
		attributes:  attributes,
		vxml:        virtualXml,
		nextSibling: -1,
	}

	elem.init()
	return elem
}

func (s *startElement) init() {
	if s.StartElement == nil {
		return
	}

	s.attributesName = make([]string, len(s.Attr))

	for i, attr := range s.attributes {
		attributeName := string(s.vxml.template[attr.keyStart():attr.keyEnd()])
		s.attributesName[i] = attributeName
	}

	if len(s.Attr) > 5 {
		s.attributeIndex = map[string]int{}
		for i, attr := range s.attributesName {
			s.attributeIndex[attr] = i
		}
	}
}

func attributeOf(spans [2]span, counter int) *attribute {
	return &attribute{
		spans: spans,
		index: counter,
	}
}
