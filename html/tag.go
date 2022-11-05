package html

type (
	tags []*tag
	tag  struct {
		innerHTML *span
		tagName   *span
		attrEnd   int32
		depth     int32
		index     int32
	}
)

func (t tags) tagOffset(index int32, offsets []int32) int32 {
	if t[index].attrEnd-1 < 0 {
		return 0
	}
	return offsets[t[index].attrEnd-1]
}
