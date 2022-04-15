package dm

import "strings"

type (
	builder struct {
		attributes        attrs
		tags              tags
		tagIndexes        []int
		tagCounter        int
		depth             int
		tagsGrouped       [][]int
		attributesGrouped [][]int
		*index
	}
)

func newBuilder() *builder {
	builder := &builder{
		tagsGrouped:       make([][]int, lastTag),
		attributesGrouped: make([][]int, lastAttribute),
		index:             newIndex(),
	}

	builder.attributes = append(builder.attributes, &attr{
		boundaries: [2]*span{{}, {}}, //tag 0 is a sentinel
		tag:        0,
	})
	builder.tags = append(builder.tags, &tag{
		depth:     -1,
		innerHTML: &span{},
		tagName:   &span{},
		attrEnd:   1,
	})

	return builder
}

func (b *builder) attribute(template []byte, spans [2]span) {
	attributeName := string(template[spans[0].start:spans[0].end])
	attributeGroup := b.index.attributeIndex(attributeName, true)
	if attributeGroup < len(b.attributesGrouped) {
		b.attributesGrouped[attributeGroup] = append(b.attributesGrouped[attributeGroup], len(b.attributes))
	} else {
		b.attributesGrouped = append(b.attributesGrouped, []int{len(b.attributes)})
	}

	b.attributes = append(b.attributes, &attr{
		tag: b.tagCounter,
		boundaries: [2]*span{
			{
				start: spans[0].start,
				end:   spans[0].end,
			},
			{
				start: spans[1].start,
				end:   spans[1].end,
			},
		},
	})
}

func (b *builder) newTag(tagName string, start int, tagSpan span, selfClosing bool) {
	aTag := &tag{
		depth: b.depth,
		innerHTML: &span{
			start: start,
		},
		tagName: &span{
			start: tagSpan.start,
			end:   tagSpan.end,
		},
		index: b.tagCounter + 1,
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
		b.tagIndexes = append(b.tagIndexes, b.tagCounter)
	}

	b.tags = append(b.tags, aTag)
}

func (b *builder) closeTag(end int) {
	b.depth--

	lastTag := b.tags[b.tagIndexes[len(b.tagIndexes)-1]]
	if lastTag.innerHTML.end == 0 {
		lastTag.innerHTML.end = end
	}

	b.tagIndexes = b.tagIndexes[:len(b.tagIndexes)-1]
}

func (b *builder) attributesBuilt() {
	b.tags[b.lastTag()].attrEnd = len(b.attributes)
}

func (b *builder) lastTag() int {
	return len(b.tags) - 1
}
