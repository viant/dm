package html

import "strings"

type (
	tags []*tag
	tag  struct {
		vdom *DOM

		innerHTML *span
		tagName   *span

		attrs     []*attr
		attrIndex map[string]int

		attrEnd int
		depth   int
		index   int
	}
)

func (t tags) tagOffset(index int, offsets []int) int {
	if t[index].attrEnd-1 < 0 {
		return 0
	}

	return offsets[t[index].attrEnd-1]
}

func (t *tag) addAttribute(spans [2]span, index int, offset int) {
	t.attrIndex[strings.ToLower(string(t.vdom.template[spans[0].start:spans[0].end]))] = len(t.attrs)
	t.attrs = append(t.attrs, &attr{
		boundaries: [2]*span{
			{
				start: spans[0].start + offset,
				end:   spans[0].end + offset,
			},
			{
				start: spans[1].start + offset,
				end:   spans[1].end + offset,
			},
		},
		index: index,
		tag:   t,
	})

	t.attrEnd = index
}

func (t *tag) attributeByName(name string) (*attr, bool) {
	if len(t.attrs) > 5 {
		attrIndex, ok := t.attrIndex[name]
		if !ok {
			return nil, false
		}

		return t.attrs[attrIndex], true
	}

	for _, attribute := range t.attrs {
		//TODO: resolve attr name lowercased before
		if strings.EqualFold(string(t.vdom.template[attribute.keyStart():attribute.keyEnd()]), name) {
			return attribute, true
		}
	}

	return nil, false
}
