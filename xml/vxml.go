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

func (v *Schema) init() error {
	decoder := xml.NewDecoder(bytes.NewReader(v.template))
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
			err := v.builder.addElement(actual, int(decoder.InputOffset()), v.template[prevOffset:decoder.InputOffset()], prevOffset)
			if err != nil {
				return err
			}
		case xml.EndElement:
			v.builder.closeElement()
		case xml.CharData:
			v.builder.addCharData(int(decoder.InputOffset()), actual)
		}

		prevOffset = int(decoder.InputOffset())
	}
}

func (v *Schema) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			v.bufferSize = int(actual)
		}
	}
}
