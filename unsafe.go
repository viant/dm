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

func attributesSpan(tokenizer *html.Tokenizer) [][2]span {
	return *(*[][2]span)(attrField.Pointer(unsafe.Pointer(tokenizer)))
}

func rawSpan(tokenizer *html.Tokenizer) span {
	return *(*span)(rawField.Pointer(unsafe.Pointer(tokenizer)))
}

func dataSpan(tokenizer *html.Tokenizer) span {
	return *(*span)(dataField.Pointer(unsafe.Pointer(tokenizer)))
}

func asBytes(value string) []byte {
	return []byte(value)
}
