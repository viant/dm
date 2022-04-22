package xml

import (
	"bytes"
	"github.com/viant/dm/option"
	"strings"
)

func FiltersOf(xpath ...string) (*option.Filters, error) {
	filters, err := NewFilters(xpath...)
	if err != nil {
		return nil, err
	}

	return option.NewFilters(filters...), nil
}

//NewFilters creates slice option.Filter based on given xpaths
func NewFilters(xpath ...string) ([]*option.Filter, error) {
	filters := make([]*option.Filter, 0)

	for i := 0; i < len(xpath); i++ {
		parsedFilters, err := newFilters(xpath[i])
		if err != nil {
			return nil, err
		}

		filters = append(filters, parsedFilters...)
	}

	return filters, nil
}

func newFilters(xpath string) ([]*option.Filter, error) {
	filters := make([]*option.Filter, 0)

	prev := -1
	var name string
	var attributes []string

	var i = 0
	for {
		switch xpath[i] {
		case '[', '/', ' ', '\n', '\t', '\r', '\v', '\f', ']':
			i++
			continue
		}

		break
	}

	for i = 0; i < len(xpath); {
		b := xpath[i]
		if prev == -1 {
			prev = i
		}

		switch b {
		case '[':
			name = xpath[prev:i]
			prev = -1
		case ' ', '\n', '\t', '\r', '\v', '\f':
			attributes = append(attributes, strings.TrimSpace(xpath[prev:i]))
			prev = -1
		case 'a':
			if isWhitespace(xpath[prev-1]) && bytes.HasPrefix([]byte(xpath[prev:]), []byte("and")) {
				i += 3
				if i < len(xpath) && isWhitespace(xpath[i]) {
					i++
					prev = -1
				}
				continue
			}
		case 'o':
			if isWhitespace(xpath[prev]) && bytes.HasPrefix([]byte(xpath[prev:]), []byte("or")) {
				i += 2
				if i < len(xpath) && isWhitespace(xpath[i]) {
					i++
					prev = -1
				}
				continue
			}
		case ']':
			attributes = append(attributes, strings.TrimSpace(xpath[prev:i]))
			prev = -1
		case '/':
			if name == "" {
				name = xpath[prev : i-1]
			}

			filters = append(filters, option.NewFilter(name, attributes...))
			attributes = nil
			prev = -1
			name = ""
		}

		i++
	}

	if prev != -1 || name != "" {
		if name == "" {
			name = xpath[prev:]
		}

		filters = append(filters, option.NewFilter(name, attributes...))
	}
	return filters, nil
}
