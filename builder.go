package dm

type (
	elementsBuilder struct {
		attributes attrs
		tags
		tagIndexes []int
		tagCounter int
		offset     int
		depth      int
	}
)

func newBuilder() *elementsBuilder {
	builder := &elementsBuilder{}
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

func (b *elementsBuilder) attribute(spans [2]span) {
	b.attributes = append(b.attributes, &attr{
		tag: b.tagCounter,
		boundaries: [2]*span{
			{
				start: spans[0].start + b.offset,
				end:   spans[0].end + b.offset,
			},
			{
				start: spans[1].start + b.offset,
				end:   spans[1].end + b.offset,
			},
		},
	})
}

func (b *elementsBuilder) newTag(start int, tagSpan span, selfClosing bool) {
	aTag := &tag{
		depth: b.depth,
		innerHTML: &span{
			start: start + b.offset,
		},
		tagName: &span{
			start: tagSpan.start + b.offset,
			end:   tagSpan.end + b.offset,
		},
		index: b.tagCounter + 1,
	}

	b.tagCounter++
	if selfClosing {
		aTag.innerHTML.end = start + b.offset
	} else {
		b.depth++
		b.tagIndexes = append(b.tagIndexes, b.tagCounter)
	}

	b.tags = append(b.tags, aTag)
}

func (b *elementsBuilder) closeTag(end int) {
	b.depth--
	b.tags[b.tagIndexes[len(b.tagIndexes)-1]].innerHTML.end = end
	b.tagIndexes = b.tagIndexes[:len(b.tagIndexes)-1]
}

func (b *elementsBuilder) attributesBuilt() {
	b.tags[b.lastTag()].attrEnd = len(b.attributes)
}

func (b *elementsBuilder) lastTag() int {
	return len(b.tags) - 1
}
