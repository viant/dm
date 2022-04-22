package html

type (
	attrs []*attr
	attr  struct {
		boundaries [2]*span
		tag        int
	}
)

func (a *attr) valueStart() int {
	return a.boundaries[1].start
}

func (a *attr) valueEnd() int {
	return a.boundaries[1].end
}

func (a *attr) keyStart() int {
	return a.boundaries[0].start
}

func (a *attr) keyEnd() int {
	return a.boundaries[0].end
}