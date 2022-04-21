package xml

import "encoding/xml"

type (
	startElement struct {
		span

		name        string
		parent      *startElement
		children    []*startElement
		parentIndex int
		elemIndex   int
		nextSibling int
		schema      *Schema
		indent      []byte

		elementsIndex          map[string][]int
		attributeIndex         map[string]int
		attributeChildrenIndex map[string][]int

		attributesName   []string
		attributes       []*attribute
		childrenAttrSize int
	}

	attribute struct {
		spans [2]*span
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

	s.elementsIndex[child.name] = append(s.elementsIndex[child.name], len(s.children))
	s.indexChildAttributes(child)

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

func newStartElement(element *xml.StartElement, schema *Schema, elemIndex int, startPosition int, attributes []*attribute) *startElement {
	var elemName string
	if element != nil {
		elemName = element.Name.Local
		if element.Name.Space != "" {
			elemName = element.Name.Local + ":" + element.Name.Space
		}
	}

	elem := &startElement{
		elemIndex: elemIndex,
		span: span{
			start: startPosition,
		},
		name:                   elemName,
		attributes:             attributes,
		schema:                 schema,
		nextSibling:            -1,
		elementsIndex:          map[string][]int{},
		attributeChildrenIndex: map[string][]int{},
	}

	elem.init()
	return elem
}

func (s *startElement) init() {
	s.indexAttributes()
}

func (s *startElement) indexAttributes() {
	s.attributesName = make([]string, len(s.attributes))
	for i, attr := range s.attributes {
		attributeName := string(s.schema.template[attr.keyStart():attr.keyEnd()])
		s.attributesName[i] = attributeName
	}

	if len(s.attributes) > mapSize {
		s.attributeIndex = map[string]int{}
		for i, attr := range s.attributesName {
			s.attributeIndex[attr] = i
		}
	}
}

func (s *startElement) attributeValueSpan(attributeIndex int) *span {
	return s.attributes[attributeIndex].spans[1]
}

func (s *startElement) indexChildAttributes(child *startElement) {
	index := len(s.children)
	for _, attr := range child.attributes {
		s.attributeChildrenIndex[s.schema.templateSlice(attr.spans[0])] = append(s.attributeChildrenIndex[s.schema.templateSlice(attr.spans[0])], index)
	}

	s.childrenAttrSize += len(child.attributes)
}

func attributeOf(spans [2]*span, counter int) *attribute {
	return &attribute{
		spans: spans,
		index: counter,
	}
}
