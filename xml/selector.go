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
	MatchAny   bool
}

//NewSelectors creates Selectors from xpath
func NewSelectors(xpath string) ([]Selector, error) {
	if len(xpath) == 0 {
		return []Selector{}, nil
	}

	selectors := make([]Selector, 0)

	matchAny := true
	if xpath[0] == '/' {
		matchAny = false
		xpath = xpath[1:]
	}

	cursor := parsly.NewCursor("", []byte(xpath), 0)
	err := parse(cursor, &selectors, false)
	if err != nil {
		return nil, err
	}

	if len(selectors) > 0 {
		selectors[0].MatchAny = matchAny
	}

	return selectors, nil
}
