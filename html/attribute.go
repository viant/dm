package html

type (
	attrs []*attr
	attr  struct {
		boundaries [2]*span
		tag        int32
	}
)

func (a *attr) valueStart() int32 {
	return int32(a.boundaries[1].start)
}

func (a *attr) valueEnd() int32 {
	return int32(a.boundaries[1].end)
}

func (a *attr) keyStart() int32 {
	return int32(a.boundaries[0].start)
}

func (a *attr) keyEnd() int32 {
	return int32(a.boundaries[0].end)
}
