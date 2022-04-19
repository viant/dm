package xml

import "encoding/xml"

type (
	StartElement struct {
		*xml.StartElement
		span

		parent      *StartElement
		children    []*StartElement
		parentIndex int
		elemIndex   int

		attributeIndex map[string]int
		attributesName []string
	}

	span struct {
		start int
		end   int
	}
)

func (s *StartElement) Append(child *StartElement) {
	child.parentIndex = len(s.children)
	s.children = append(s.children, child)
	child.parent = s
}

func (s *StartElement) attrByName(attribute string) (string, bool) {
	if s.attributeIndex != nil {
		value, ok := s.attributeIndex[attribute]
		if !ok {
			return "", false
		}
		return s.Attr[value].Value, ok
	}

	for i, attr := range s.attributesName {
		if attr == attribute {
			return s.Attr[i].Value, true
		}
	}

	return "", false
}

func NewStartElement(element *xml.StartElement, elemIndex int, startPosition int) *StartElement {

	elem := &StartElement{
		StartElement: element,
		elemIndex:    elemIndex,
		span: span{
			start: startPosition,
		},
	}

	elem.init()
	return elem
}

func (s *StartElement) init() {
	if s.StartElement == nil {
		return
	}

	s.attributesName = make([]string, len(s.Attr))

	for i, attr := range s.Attr {
		attributeName := attr.Name.Local
		if attr.Name.Space != "" {
			attributeName = attr.Name.Space + ":" + attr.Name.Local
		}

		s.attributesName[i] = attributeName
	}

	if len(s.Attr) > 5 {
		s.attributeIndex = map[string]int{}
		for i, attr := range s.attributesName {
			s.attributeIndex[attr] = i
		}
	}
}
