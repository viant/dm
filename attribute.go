package vhtml

type (
	Attributes []*Attribute
	Attribute  struct {
		Boundaries [2]*Span
		Tag        []byte
	}

	AttributesBuilder struct {
		attributes Attributes
	}
)

func (a *Attribute) ValueStart() int {
	return a.Boundaries[1].Start
}

func (a *Attribute) ValueEnd() int {
	return a.Boundaries[1].End
}

func (a *Attribute) KeyStart() int {
	return a.Boundaries[0].Start
}

func (a *Attribute) KeyEnd() int {
	return a.Boundaries[0].End
}

func NewBuilder() *AttributesBuilder {
	return &AttributesBuilder{}
}

func (b *AttributesBuilder) Attribute(parent []byte, spans [2]Span) {
	b.attributes = append(b.attributes, &Attribute{
		Tag: parent,
		Boundaries: [2]*Span{
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

func (b *AttributesBuilder) Attributes() Attributes {
	return b.attributes
}
