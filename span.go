package dm

type Span struct {
	Start int
	End   int
}

func newSpan(span *Span, startOffset, endOffset int) *Span {
	return &Span{
		Start: span.Start + startOffset,
		End:   span.End + endOffset,
	}
}
