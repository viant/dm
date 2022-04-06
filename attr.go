package dm

import (
	"github.com/viant/xunsafe"
	"golang.org/x/net/html"
	"reflect"
	"unsafe"
)

var attrField *xunsafe.Field
var rawField *xunsafe.Field
var dataField *xunsafe.Field

func init() {
	rType := reflect.TypeOf(&html.Tokenizer{})
	attrField = xunsafe.FieldByName(rType, "attr")
	rawField = xunsafe.FieldByName(rType, "raw")
	dataField = xunsafe.FieldByName(rType, "data")
}

func attributesSpan(tokenizer *html.Tokenizer) [][2]Span {
	return *(*[][2]Span)(attrField.Pointer(unsafe.Pointer(tokenizer)))
}

func rawSpan(tokenizer *html.Tokenizer) Span {
	return *(*Span)(rawField.Pointer(unsafe.Pointer(tokenizer)))
}

func dataSpan(tokenizer *html.Tokenizer) Span {
	return *(*Span)(dataField.Pointer(unsafe.Pointer(tokenizer)))
}
