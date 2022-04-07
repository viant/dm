package dm

type (
	tags []*tag
	tag  struct {
		InnerHTML *Span
		TagName   *Span
		AttrEnd   int
		Depth     int
	}
)

func (t tags) sliceTo(offsets []int, index int) tags {
	result := make([]*tag, index+1)
	result[0] = t[0]

	for i := 1; i <= index; i++ {
		result[i] = &tag{
			InnerHTML: newSpan(t[i].InnerHTML, t.tagOffset(i, offsets), t.tagOffset(t.lastChild(i), offsets)),
			TagName:   newSpan(t[i].TagName, t.tagOffset(i-1, offsets), t.tagOffset(i-1, offsets)),
			AttrEnd:   t[i].AttrEnd,
		}
	}
	return result
}

//TODO: need to be optimized.
func (t tags) lastChild(tagIndex int) int {
	var i int
	for i = tagIndex + 1; i < len(t); i++ {
		if t[tagIndex].Depth == t[i].Depth {
			i--
			break
		}
	}

	if i >= len(t) {
		return len(t) - 1
	}

	return i
}

func (t tags) tagOffset(index int, offsets []int) int {
	if t[index].AttrEnd-1 < 0 {
		return 0
	}

	return offsets[t[index].AttrEnd-1]
}
