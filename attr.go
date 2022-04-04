package dm

import (
	"github.com/viant/xunsafe"
	"golang.org/x/net/html"
	"reflect"
	"unsafe"
)

var attrField *xunsafe.Field

func init() {
	rType := reflect.TypeOf(&html.Tokenizer{})
	attrField = xunsafe.FieldByName(rType, "attr")
}

func attributesSpan(tokenizer *html.Tokenizer) [][2]Span {
	return *(*[][2]Span)(attrField.Pointer(unsafe.Pointer(tokenizer)))
}
