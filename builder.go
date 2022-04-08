package dm

type (
	builder struct {
		attributes  attrs
		tags        tags
		tagIndexes  []int
		tagCounter  int
		depth       int
		tagsGrouped [][]int
		*index
	}
)

func newBuilder() *builder {
	builder := &builder{
		tagsGrouped: make([][]int, lastTag),
		index:       newIndex(),
	}

	builder.attributes = append(builder.attributes, &attr{
		boundaries: [2]*span{{}, {}}, //tag 0 is a sentinel
		tag:        0,
	})
	builder.tags = append(builder.tags, &tag{
		depth:     -1,
		innerHTML: &span{},
		tagName:   &span{},
	})

	return builder
}

func (b *builder) attribute(spans [2]span) {
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

	tagGroupPosition := b.index.tag(tagName, true)
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
	b.tags[b.tagIndexes[len(b.tagIndexes)-1]].innerHTML.end = end
	b.tagIndexes = b.tagIndexes[:len(b.tagIndexes)-1]
}

func (b *builder) attributesBuilt() {
	b.tags[b.lastTag()].attrEnd = len(b.attributes)
}

func (b *builder) lastTag() int {
	return len(b.tags) - 1
}
