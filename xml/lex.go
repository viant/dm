package xml

import (
	"github.com/viant/parsly"
	pmatcher "github.com/viant/parsly/matcher"
)

const (
	whitespaceToken = iota

	nameToken
	colonToken

	attrBlockToken
	newElemToken

	andToken

	equalToken
	notEqualToken

	stringToken
)

var whitespace = parsly.NewToken(whitespaceToken, "Whitespace", pmatcher.NewWhiteSpace())

var name = parsly.NewToken(nameToken, "Name", newIdentity())
var colon = parsly.NewToken(colonToken, "Name after namespace", pmatcher.NewByte(':'))

var attrBlock = parsly.NewToken(attrBlockToken, "Attributes block", pmatcher.NewBlock('[', ']', '\\'))
var newElem = parsly.NewToken(newElemToken, "New element", pmatcher.NewByte('/'))

var and = parsly.NewToken(andToken, "And", pmatcher.NewFragmentsFold([]byte("and")))

var equal = parsly.NewToken(equalToken, "Equal", pmatcher.NewBytes([]byte(EQ)))
var notEqual = parsly.NewToken(notEqualToken, "Not equal", pmatcher.NewBytes([]byte(NEQ)))

var stringValue = parsly.NewToken(stringToken, "String value", newStringMatcher('\''))
