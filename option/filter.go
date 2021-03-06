package option

import "strings"

type (
	//Filters represents tag indexing filters
	Filters struct {
		Tags  []*Filter
		index map[string]int
	}

	//Filter represents single tag filter
	Filter struct {
		Name          string
		Attributes    []string
		CaseSensitive bool
		index         map[string]int
	}
)

//Init initializes Filters
func (f *Filters) Init() {
	f.index = map[string]int{}
	for i, filter := range f.Tags {
		f.index[filter.Name] = i
		f.index[strings.ToLower(filter.Name)] = i
		f.index[strings.ToUpper(filter.Name)] = i
	}
}

func (f *Filters) ElementFilter(tagName string, caseSensitive bool) (*Filter, bool) {
	if !caseSensitive {
		tagName = strings.ToLower(tagName)
	}

	if len(f.Tags) >= 5 {
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

func (f *Filters) Add(filters ...*Filter) *Filters {
	for i, filter := range filters {
		elementFilter, ok := f.ElementFilter(filter.Name, filter.CaseSensitive)
		if !ok {
			f.add(filters[i])
			continue
		}

		elementFilter.addAttributes(filters[i].Attributes)
	}

	return f
}

func (f *Filters) add(filter *Filter) {
	f.index[filter.Name] = len(f.Tags)
	f.Tags = append(f.Tags, filter)
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
	return newFilter(tag, attributes, false)
}

func CaseSensitiveFilter(tag string, attributes ...string) *Filter {
	return newFilter(tag, attributes, true)
}

func newFilter(tag string, attributes []string, caseSensitive bool) *Filter {
	newAttributes := make([]string, len(attributes))
	for i, attribute := range attributes {
		newAttributes[i] = strings.ToLower(attribute)
	}

	tagFilter := &Filter{
		Name:          tag,
		Attributes:    newAttributes,
		CaseSensitive: caseSensitive,
		index:         nil,
	}

	tagFilter.Init()
	return tagFilter
}

//Init initializes Filter
func (f *Filter) Init() {
	f.index = map[string]int{}

	for i, attribute := range f.Attributes {
		f.index[attribute] = i
		f.index[strings.ToLower(attribute)] = i
		f.index[strings.ToUpper(attribute)] = i
	}
}

func (f *Filter) Matches(attributeName string) bool {
	if len(f.Attributes) >= 5 {
		_, ok := f.index[attributeName]
		return ok
	}

	for _, filter := range f.Attributes {
		if f.CaseSensitive && filter == attributeName || !f.CaseSensitive && strings.EqualFold(attributeName, filter) {
			return true
		}
	}

	return false
}

func (f *Filter) addAttributes(attributes []string) {
	for i, attribute := range attributes {
		if f.Matches(attribute) {
			continue
		}

		f.index[attributes[i]] = len(f.Attributes)
		f.Attributes = append(f.Attributes, attribute)
	}
}
