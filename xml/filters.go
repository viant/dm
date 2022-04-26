package xml

import (
	"fmt"
	"github.com/viant/dm/option"
	"github.com/viant/parsly"
	"io"
)

func FiltersOf(xpaths ...string) (*option.Filters, error) {
	result, err := NewFilters(xpaths...)
	if err != nil {
		return nil, err
	}

	return option.NewFilters(result...), nil
}

//NewFilters creates slice option.Filter based on given xpaths
func NewFilters(xpaths ...string) ([]*option.Filter, error) {
	result := make([]Selector, 0)

	for _, xpath := range xpaths {
		cursor := parsly.NewCursor("", []byte(xpath), 0)
		err := parse(cursor, &result, true)
		if err != nil {
			return nil, err
		}
	}

	filters := make([]*option.Filter, len(result))
	for i, selector := range result {
		names := make([]string, len(selector.Attributes))
		for j, attributeSelector := range selector.Attributes {
			names[j] = attributeSelector.Name
		}

		filters[i] = option.NewFilter(selector.Name, names...)
	}

	return filters, nil
}

func parse(cursor *parsly.Cursor, selectors *[]Selector, valueOptional bool) error {
	elemName, err := matchName(cursor, true)
	if err != nil {
		return err
	}

	candidates := []*parsly.Token{attrBlock, newElem}
	matched := cursor.MatchAny(candidates...)

	switch matched.Code {
	case parsly.EOF:
		if elemName != "" {
			*selectors = append(*selectors, ElementSelector(elemName))
		}

	case attrBlockToken:
		attrFragment := matched.Text(cursor)
		attrCursor := parsly.NewCursor("", []byte(attrFragment[1:len(attrFragment)-1]), 0)

		attrs, err := matchAttributes(attrCursor, valueOptional)
		if err != nil {
			return err
		}

		*selectors = append(*selectors, ElementSelector(elemName, attrs...))

		cursor.MatchOne(newElem)
		if matched.Code == newElemToken {
			return parse(cursor, selectors, valueOptional)
		}

	case newElemToken:
		*selectors = append(*selectors, ElementSelector(elemName))
		return parse(cursor, selectors, valueOptional)

	default:
		return cursor.NewError(candidates...)
	}

	return nil
}

func matchAttributes(cursor *parsly.Cursor, valueOptional bool) ([]AttributeSelector, error) {
	var attributes []AttributeSelector
	for {
		cursor.MatchOne(whitespace)
		name, err := matchName(cursor, true)

		if err != nil {
			return nil, err
		}

		candidates := []*parsly.Token{equal, notEqual}
		matched := cursor.MatchAfterOptional(whitespace, candidates...)

		token, err := extractToken(matched, cursor)
		if err != nil && !valueOptional {
			return nil, err
		}

		value, err := matchValue(cursor)
		if err != nil && !valueOptional {
			return nil, err
		}

		candidates = []*parsly.Token{and}
		matched = cursor.MatchAfterOptional(whitespace, candidates...)
		switch matched.Code {
		case parsly.EOF, andToken:
			attributes = append(attributes, AttributeSelector{
				Name:    name,
				Value:   value,
				Compare: token,
			})

			if matched.Code == andToken {
				continue
			}

			return attributes, nil
		}
	}
}

func matchValue(cursor *parsly.Cursor) (string, error) {
	matched := cursor.MatchAfterOptional(whitespace, stringValue)
	switch matched.Code {
	case stringToken:
		value := matched.Text(cursor)
		return value[1 : len(value)-1], nil
	default:
		return "", cursor.NewError(stringValue)
	}
}

func extractToken(matched *parsly.TokenMatch, cursor *parsly.Cursor) (ComparisonToken, error) {
	switch matched.Code {
	case parsly.Invalid, parsly.EOF:
		return "", nil
	default:
		return ComparisonToken(matched.Text(cursor)), nil
	}
}

func matchName(cursor *parsly.Cursor, allowNamespace bool) (string, error) {
	matched := cursor.MatchAfterOptional(whitespace, name)
	prevPos := cursor.Pos

	switch matched.Code {
	case nameToken:
		elemName := matched.Text(cursor)
		matched = cursor.MatchOne(colon)
		if matched.Code == colonToken {
			if !allowNamespace {
				return "", fmt.Errorf("multiple namespaces not allowed")
			}

			actualName, err := matchName(cursor, false)
			if err != nil {
				return "", fmt.Errorf("%w, %v", err, string(cursor.Input[prevPos:cursor.Pos]))
			}

			elemName += ":" + actualName
		}

		return elemName, nil

	case parsly.EOF:
		return "", io.EOF

	default:
		return "", cursor.NewError(name)
	}
}
