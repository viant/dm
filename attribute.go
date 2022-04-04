package dm

type (
	attrs []*attr
	attr  struct {
		boundaries [2]*Span
		tag        []byte
	}

	attributesBuilder struct {
		attributes attrs
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

func newBuilder() *attributesBuilder {
	return &attributesBuilder{}
}

func (b *attributesBuilder) attribute(parent []byte, spans [2]Span) {
	b.attributes = append(b.attributes, &attr{
		tag: parent,
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

func (b *attributesBuilder) result() attrs {
	return b.attributes
}
