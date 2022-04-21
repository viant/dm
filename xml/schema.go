package xml

import (
	"bytes"
	"encoding/xml"
)

//Schema represents Xml schema
type Schema struct {
	template []byte
	*builder
	bufferSize int
}

//New creates new Schema
func New(xml string, options ...Option) (*Schema, error) {
	vxml := &Schema{
		template:   []byte(xml),
		bufferSize: len(xml),
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

func (s *Schema) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			s.bufferSize = int(actual)
		}
	}
}

func (s *Schema) templateSlice(span *span) string {
	return string(s.template[span.start:span.end])
}
