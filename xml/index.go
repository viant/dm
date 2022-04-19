package xml

import "encoding/xml"

type (
	attribute struct {
		spans [2]span
		index int
	}

	StartElement struct {
		*xml.StartElement
		span

		parent      *StartElement
		children    []*StartElement
		parentIndex int
		elemIndex   int
		vxml        *VirtualXml

		attributeIndex map[string]int
		attributesName []string
		attributes     []*attribute
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

func (s *StartElement) Append(child *StartElement) {
	child.parentIndex = len(s.children)
	s.children = append(s.children, child)
	child.parent = s
}

func (s *StartElement) attrByName(attribute string) (int, bool) {
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

func NewStartElement(element *xml.StartElement, virtualXml *VirtualXml, elemIndex int, startPosition int, attributes []*attribute) *StartElement {
	elem := &StartElement{
		StartElement: element,
		elemIndex:    elemIndex,
		span: span{
			start: startPosition,
		},
		attributes: attributes,
		vxml:       virtualXml,
	}

	elem.init()
	return elem
}

func (s *StartElement) init() {
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

func newAttribute(spans [2]span, counter int) *attribute {
	return &attribute{
		spans: spans,
		index: counter,
	}
}
