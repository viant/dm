package xml

import "github.com/viant/parsly"

type identity struct {
}

//Match matches a string
func (n *identity) Match(cursor *parsly.Cursor) (matched int) {
	input := cursor.Input
	pos := cursor.Pos
	if startsWithCharacter := isLetter(input[pos]); startsWithCharacter {
		pos++
		matched++
	} else {
		return 0
	}

	size := len(input)
	for i := pos; i < size; i++ {
		switch input[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
			matched++
			continue
		case '\n', '\r', ' ', ':':
			return matched

		default:
			if isLetter(input[i]) {
				matched++
				continue
			}

			return matched
		}
	}

	return matched
}

func isLetter(b byte) bool {
	if (b < 'a' || b > 'z') && (b < 'A' || b > 'Z') {
		return false
	}
	return true
}

//newIdentity creates a string matcher
func newIdentity() *identity {
	return &identity{}
}
