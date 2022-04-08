package dm

type (
	tags []*tag
	tag  struct {
		innerHTML *span
		tagName   *span
		attrEnd   int
		depth     int
		index     int
	}
)

func (t tags) tagOffset(index int, offsets []int) int {
	if t[index].attrEnd-1 < 0 {
		return 0
	}
	return offsets[t[index].attrEnd-1]
}
