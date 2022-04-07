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
		boundaries: [2]*Span{{}, {}}, //tag 0 is a sentinel
		tag:        0,
	})
	builder.tags = append(builder.tags, &tag{
		Depth:     -1,
		InnerHTML: &Span{},
		TagName:   &Span{},
	})

	return builder
}

func (b *elementsBuilder) attribute(spans [2]Span) {
	b.attributes = append(b.attributes, &attr{
		tag: b.tagCounter,
		boundaries: [2]*Span{
			{
				Start: spans[0].Start + b.offset,
				End:   spans[0].End + b.offset,
			},
			{
				Start: spans[1].Start + b.offset,
				End:   spans[1].End + b.offset,
			},
		},
	})
}

func (b *elementsBuilder) newTag(start int, span Span, selfClosing bool) {
	aTag := &tag{
		Depth: b.depth,
		InnerHTML: &Span{
			Start: start + b.offset,
		},
		TagName: &Span{
			Start: span.Start + b.offset,
			End:   span.End + b.offset,
		},
	}

	b.tagCounter++
	if selfClosing {
		aTag.InnerHTML.End = start + b.offset
	} else {
		b.depth++
		b.tagIndexes = append(b.tagIndexes, b.tagCounter)
	}

	b.tags = append(b.tags, aTag)
}

func (b *elementsBuilder) closeTag(end int) {
	b.depth--
	b.tags[b.tagIndexes[len(b.tagIndexes)-1]].InnerHTML.End = end
	b.tagIndexes = b.tagIndexes[:len(b.tagIndexes)-1]
}

func (b *elementsBuilder) attributesBuilt() {
	b.tags[b.lastTag()].AttrEnd = len(b.attributes)
}

func (b *elementsBuilder) lastTag() int {
	return len(b.tags) - 1
}
