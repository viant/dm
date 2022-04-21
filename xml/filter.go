package xml

type (
	Filters struct {
		elementsFilter []*Filter
		filterIndex    map[string]int
	}

	Filter struct {
		elementName string
		attributes  []string

		attributesIndex map[string]int
	}
)

func NewFilters(filters ...*Filter) *Filters {
	f := &Filters{
		elementsFilter: filters,
	}

	f.init()
	return f
}

func (f *Filters) ElementFilter(name string) (*Filter, bool) {
	if f.filterIndex != nil {
		filter, ok := f.filterIndex[name]
		if ok {
			return f.elementsFilter[filter], true
		}

		return nil, false
	}

	for _, filter := range f.elementsFilter {
		if filter.elementName == name {
			return filter, true
		}
	}

	return nil, false
}

func (f *Filters) init() {
	if len(f.elementsFilter) <= mapSize || f.filterIndex == nil {
		return
	}

	f.filterIndex = map[string]int{}
	for i, filter := range f.elementsFilter {
		f.filterIndex[filter.elementName] = i
	}
}

func NewFilter(elementName string, attributes ...string) *Filter {
	f := &Filter{
		elementName: elementName,
		attributes:  attributes,
	}

	f.init()
	return f
}

func (f *Filter) init() {
	if len(f.attributes) <= mapSize || f.attributesIndex == nil {
		return
	}

	f.attributesIndex = map[string]int{}
	for i, attrName := range f.attributes {
		f.attributesIndex[attrName] = i
	}
}

func (f *Filter) Contains(attribute string) bool {
	if f.attributesIndex != nil {
		_, ok := f.attributesIndex[attribute]
		return ok
	}

	for _, attrName := range f.attributes {
		if attrName == attribute {
			return true
		}
	}

	return false
}
