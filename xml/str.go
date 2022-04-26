package xml

import (
	"github.com/viant/parsly"
)

type stringMatcher struct {
	quote byte
}

func (s *stringMatcher) Match(cursor *parsly.Cursor) (matched int) {
	if cursor.Input[cursor.Pos] != s.quote {
		return 0
	}

	matched = 1
	for i := cursor.Pos + matched; i < cursor.InputSize; i++ {
		matched++
		if cursor.Input[i] == s.quote {
			return matched
		}
	}

	return 0
}

func newStringMatcher(quote byte) *stringMatcher {
	return &stringMatcher{
		quote: quote,
	}
}
