package xml

import (
	"bytes"
	"encoding/xml"
	"github.com/viant/dm/option"
)

//Schema represents Xml schema
type Schema struct {
	template []byte
	*builder
	bufferSize int

	attributesChangesSize int
	elementsChangesSize   int
}

//New creates new Schema
func New(xml string, options ...option.Option) (*Schema, error) {
	vxml := &Schema{
		template:              []byte(xml),
		bufferSize:            len(xml),
		attributesChangesSize: prealocateSize,
		elementsChangesSize:   prealocateSize,
	}

	vxml.builder = newBuilder(vxml)
	vxml.apply(options)

	if err := vxml.init(); err != nil {
		return nil, err
	}

	return vxml, nil
}

func (s *Schema) init() error {
	decoder := xml.NewDecoder(bytes.NewReader(s.template))
	var prevOffset int
	for {
		token, err := decoder.Token()
		if err != nil {
			if token == nil {
				return nil
			}

			return err
		}

		switch actual := token.(type) {
		case xml.StartElement:
			err := s.builder.addElement(actual, int(decoder.InputOffset()), s.template[prevOffset:decoder.InputOffset()], prevOffset)
			if err != nil {
				return err
			}
		case xml.EndElement:
			s.builder.closeElement()
		case xml.CharData:
			s.builder.addCharData(int(decoder.InputOffset()), actual)
		}

		prevOffset = int(decoder.InputOffset())
	}
}

func (s *Schema) apply(options []option.Option) {
	for _, opt := range options {
		switch actual := opt.(type) {
		case option.BufferSize:
			s.bufferSize = int(actual)
		case *option.Filters:
			s.builder.filters = actual
		case AttributesChangesSize:
			s.attributesChangesSize = int(actual)
		case ElementsChangesSize:
			s.elementsChangesSize = int(actual)
		}
	}
}

func (s *Schema) templateSlice(span *span) string {
	return string(s.template[span.start:span.end])
}
