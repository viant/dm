package dm

type (
	attrs []*attr
	attr  struct {
		boundaries [2]*Span
		tag        int
	}
)

func (a attrs) sliceTo(offsets []int, end int) attrs {
	result := make([]*attr, end)
	result[0] = a[0]
	for i := 1; i < end; i++ {
		result[i] = &attr{
			boundaries: [2]*Span{
				newSpan(a[i].boundaries[0], offsets[i-1], offsets[i-1]),
				newSpan(a[i].boundaries[1], offsets[i-1], offsets[i-1]),
			},
			tag: a[i].tag,
		}
	}
	return result
}

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
