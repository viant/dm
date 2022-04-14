package dm

import "strings"

type (
	//Filter represents tag indexing filters
	Filter struct {
		Tags  []*TagFilter
		index map[string]int
	}

	//TagFilter represents single tag filter
	TagFilter struct {
		Name       string
		Attributes []string
		index      map[string]int
	}
)

//Init initializes Filter
func (f *Filter) Init() {
	if len(f.Tags) < 5 {
		return
	}

	f.index = map[string]int{}
	for i, filter := range f.Tags {
		f.index[filter.Name] = i
	}
}

func (f *Filter) tagFilter(tagName string) (*TagFilter, bool) {
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

//NewFilter creates Filter with given TagFilters
func NewFilter(filters ...*TagFilter) *Filter {
	filter := &Filter{
		Tags: filters,
	}

	filter.Init()
	return filter
}

//NewTagFilter creates new TagFilter against specified attributes and tag name
func NewTagFilter(tag string, attributes ...string) *TagFilter {
	newAttributes := make([]string, len(attributes))
	for i, attribute := range attributes {
		newAttributes[i] = strings.ToLower(attribute)
	}

	tagFilter := &TagFilter{
		Name:       tag,
		Attributes: newAttributes,
		index:      nil,
	}

	tagFilter.Init()
	return tagFilter
}

//Init initializes TagFilter
func (f *TagFilter) Init() {
	if len(f.Attributes) < 5 {
		return
	}

	f.index = map[string]int{}

	for i, attribute := range f.Attributes {
		f.index[attribute] = i
	}
}

func (f *TagFilter) matches(attributeName string) bool {
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
