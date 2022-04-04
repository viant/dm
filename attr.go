package vhtml

import (
	"github.com/viant/xunsafe"
	"golang.org/x/net/html"
	"reflect"
	"unsafe"
)

var spanField *xunsafe.Field

func init() {
	rType := reflect.TypeOf(&html.Tokenizer{})
	spanField = xunsafe.FieldByName(rType, "attr")
}

func AttributesSpan(tokenizer *html.Tokenizer) [][2]Span {
	return *(*[][2]Span)(spanField.Pointer(unsafe.Pointer(tokenizer)))
}
