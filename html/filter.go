package html

import "strings"

type (
	//Filters represents tag indexing filters
	Filters struct {
		Tags  []*Filter
		index map[string]int
	}

	//Filter represents single tag filter
	Filter struct {
		Name       string
		Attributes []string
		index      map[string]int
	}
)

//Init initializes Filters
func (f *Filters) Init() {
	if len(f.Tags) < 5 {
		return
	}

	f.index = map[string]int{}
	for i, filter := range f.Tags {
		f.index[filter.Name] = i
	}
}

func (f *Filters) tagFilter(tagName string) (*Filter, bool) {
	tagName = strings.ToLower(tagName)

	if f.index != nil {
		tagFilterIndex, ok := f.index[tagName]
		if ok == false {
			return nil, false
		}

		return f.Tags[tagFilterIndex], true
	}

	for _, filter := range f.Tags {
		if filter.Name == tagName {
			return filter, true
		}
	}

	return nil, false
}

//NewFilters creates Filters with given TagFilters
func NewFilters(filters ...*Filter) *Filters {
	filter := &Filters{
		Tags: filters,
	}

	filter.Init()
	return filter
}

//NewFilter creates new Filter against specified attributes and tag name
func NewFilter(tag string, attributes ...string) *Filter {
	newAttributes := make([]string, len(attributes))
	for i, attribute := range attributes {
		newAttributes[i] = strings.ToLower(attribute)
	}

	tagFilter := &Filter{
		Name:       tag,
		Attributes: newAttributes,
		index:      nil,
	}

	tagFilter.Init()
	return tagFilter
}

//Init initializes Filter
func (f *Filter) Init() {
	if len(f.Attributes) < 5 {
		return
	}

	f.index = map[string]int{}

	for i, attribute := range f.Attributes {
		f.index[attribute] = i
	}
}

func (f *Filter) matches(attributeName string) bool {
	attributeName = strings.ToLower(attributeName)

	if f.index != nil {
		_, ok := f.index[attributeName]
		return ok
	}

	for _, filter := range f.Attributes {
		if filter == attributeName {
			return true
		}
	}

	return false
}
