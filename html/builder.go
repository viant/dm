package html

import "strings"

type (
	builder struct {
		attributes        attrs
		tags              tags
		indexesStack      []int32
		tagCounter        int32
		depth             int32
		tagsGrouped       [][]int32
		attributesGrouped [][]int32
		*index
	}
)

func newBuilder() *builder {
	builder := &builder{
		tagsGrouped:       make([][]int32, lastTag),
		attributesGrouped: make([][]int32, lastAttribute),
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
	if attributeGroup < int32(len(b.attributesGrouped)) {
		b.attributesGrouped[attributeGroup] = append(b.attributesGrouped[attributeGroup], int32(len(b.attributes)))
	} else {
		b.attributesGrouped = append(b.attributesGrouped, []int32{int32(len(b.attributes))})
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

func (b *builder) newTag(tagName string, start int32, tagSpan span, selfClosing bool) {
	aTag := &tag{
		depth: b.depth,
		innerHTML: &span{
			start: int(start),
		},
		tagName: &span{
			start: tagSpan.start,
			end:   tagSpan.end,
		},
		index: b.tagCounter + 1,
	}

	if strings.EqualFold(tagName, "script") {
		aTag.innerHTML.end = int(start)
	}

	tagGroupPosition := b.index.tagIndex(tagName, true)
	if tagGroupPosition >= int32(len(b.tagsGrouped)) {
		b.tagsGrouped = append(b.tagsGrouped, []int32{aTag.index})
	} else {
		b.tagsGrouped[tagGroupPosition] = append(b.tagsGrouped[tagGroupPosition], aTag.index)
	}

	b.tagCounter++
	if selfClosing {
		aTag.innerHTML.end = int(start)
	} else {
		b.depth++
		b.indexesStack = append(b.indexesStack, b.tagCounter)
	}

	b.tags = append(b.tags, aTag)
}

func (b *builder) closeTag(end int32) {
	b.depth--

	lastTag := b.tags[b.indexesStack[len(b.indexesStack)-1]]
	if lastTag.innerHTML.end == 0 {
		lastTag.innerHTML.end = int(end)
	}

	b.indexesStack = b.indexesStack[:len(b.indexesStack)-1]
}

func (b *builder) attributesBuilt() {
	b.tags[b.lastTag()].attrEnd = int32(len(b.attributes))
}

func (b *builder) lastTag() int32 {
	return int32(len(b.tags)) - 1
}
