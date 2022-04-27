package html

import "strings"

type (
	builder struct {
		vdom         *DOM
		tags         tags
		indexesStack []int
		tagCounter   int
		depth        int
		tagsGrouped  [][]int
		*index

		attributeCounter int
	}
)

func newBuilder(vdom *DOM) *builder {
	builder := &builder{
		tagsGrouped: make([][]int, lastTag),
		index:       newIndex(),
		vdom:        vdom,
	}

	builder.tags = append(builder.tags, &tag{
		depth:     -1,
		innerHTML: &span{},
		tagName:   &span{},
		attrEnd:   1,
	})

	return builder
}

func (b *builder) attribute(spans [2]span) {
	b.tags[len(b.tags)-1].addAttribute(spans, b.attributeCounter)
	b.attributeCounter++
}

func (b *builder) newTag(tagName string, start int, tagSpan span, selfClosing bool) {
	tagName = strings.ToLower(tagName)

	aTag := &tag{
		vdom:      b.vdom,
		attrIndex: map[string]int{},
		depth:     b.depth,
		innerHTML: &span{
			start: start,
		},
		tagName: &span{
			start: tagSpan.start,
			end:   tagSpan.end,
		},
		index:   b.tagCounter + 1,
		attrEnd: b.attributeCounter - 1,
	}

	if strings.EqualFold(tagName, "script") {
		aTag.innerHTML.end = start
	}

	tagGroupPosition := b.index.tagIndex(tagName, true)
	if tagGroupPosition >= len(b.tagsGrouped) {
		b.tagsGrouped = append(b.tagsGrouped, []int{aTag.index})
	} else {
		b.tagsGrouped[tagGroupPosition] = append(b.tagsGrouped[tagGroupPosition], aTag.index)
	}

	b.tagCounter++
	if selfClosing {
		aTag.innerHTML.end = start
	} else {
		b.depth++
		b.indexesStack = append(b.indexesStack, b.tagCounter)
	}

	b.tags = append(b.tags, aTag)
}

func (b *builder) closeTag(end int) {
	b.depth--

	lastTag := b.tags[b.indexesStack[len(b.indexesStack)-1]]
	if lastTag.innerHTML.end == 0 {
		lastTag.innerHTML.end = end
	}

	b.indexesStack = b.indexesStack[:len(b.indexesStack)-1]
}

func (b *builder) lastTag() int {
	return len(b.tags) - 1
}
