package xml

import "github.com/viant/parsly"

//AttributeSelector matches Element by Attribute name and value
type AttributeSelector struct {
	Name    string
	Value   string
	Compare ComparisonToken
}

//Selector matches Element by name
type Selector struct {
	Name       string
	Attributes []AttributeSelector
}

//ElementSelector returns new Selector with only Element Name.
func ElementSelector(name string, attributes ...AttributeSelector) Selector {
	return Selector{
		Name:       name,
		Attributes: attributes,
	}
}

//NewSelectors creates Selectors from xpath
func NewSelectors(xpath string) ([]Selector, error) {
	selectors := make([]Selector, 0)
	cursor := parsly.NewCursor("", []byte(xpath), 0)

	err := parse(cursor, &selectors, false)
	if err != nil {
		return nil, err
	}

	return selectors, nil
}
