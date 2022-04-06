package dm

type (
	attrs []*attr
	attr  struct {
		boundaries [2]*Span
		tag        int
	}

	tags []*tag
	tag  struct {
		InnerHTML *Span
		TagName   *Span
		AttrEnd   int
	}

	elementsBuilder struct {
		attributes attrs
		tags
		tagIndexes []int
		tagCounter int
	}
)

func (a *attr) valueStart() int {
	return a.boundaries[1].Start
}

func (a *attr) valueEnd() int {
	return a.boundaries[1].End
}

func (a *attr) keyStart() int {
	return a.boundaries[0].Start
}

func (a *attr) keyEnd() int {
	return a.boundaries[0].End
}

func newBuilder() *elementsBuilder {
	builder := &elementsBuilder{}
	builder.attributes = append(builder.attributes, &attr{
		boundaries: [2]*Span{{}, {}}, //tag 0 is a sentinel
		tag:        0,
	})
	builder.tags = append(builder.tags, &tag{})

	return builder
}

func (b *elementsBuilder) attribute(spans [2]Span) {
	b.attributes = append(b.attributes, &attr{
		tag: b.tagCounter,
		boundaries: [2]*Span{
			{
				Start: spans[0].Start,
				End:   spans[0].End,
			},
			{
				Start: spans[1].Start,
				End:   spans[1].End,
			},
		},
	})
}

func (b *elementsBuilder) newTag(start int, span Span, selfClosing bool) {
	b.tagCounter++

	aTag := &tag{
		InnerHTML: &Span{
			Start: start,
		},
		TagName: &Span{
			Start: span.Start,
			End:   span.End,
		},
	}
	if selfClosing {
		aTag.InnerHTML.End = start
	} else {
		b.tagIndexes = append(b.tagIndexes, b.tagCounter)
	}

	b.tags = append(b.tags, aTag)
}

func (b *elementsBuilder) closeTag(end int) {
	b.tags[b.tagIndexes[len(b.tagIndexes)-1]].InnerHTML.End = end
	b.tagIndexes = b.tagIndexes[:len(b.tagIndexes)-1]
}

func (b *elementsBuilder) attributesBuilt() {
	b.tags[b.lastTag()].AttrEnd = len(b.attributes)
}

func (b *elementsBuilder) lastTag() int {
	return len(b.tags) - 1
}
