package xml

import (
	"bytes"
	"encoding/xml"
	"github.com/viant/dm/option"
)

//DOM represents VirtualDOM structure
type DOM struct {
	template []byte
	*builder
	bufferSize int

	attributesChangesSize int
	elementsChangesSize   int
}

//New creates new VirtualDOM
func New(template string, options ...option.Option) (*DOM, error) {
	dom := &DOM{
		template:              []byte(template),
		bufferSize:            len(template),
		attributesChangesSize: prealocateSize,
		elementsChangesSize:   prealocateSize,
	}

	dom.builder = newBuilder(dom)
	dom.apply(options)

	if err := dom.init(); err != nil {
		return nil, err
	}

	return dom, nil
}

func (d *DOM) init() error {
	decoder := xml.NewDecoder(bytes.NewReader(d.template))
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
			err := d.builder.addElement(actual, int(decoder.InputOffset()), d.template[prevOffset:decoder.InputOffset()], prevOffset)
			if err != nil {
				return err
			}
		case xml.EndElement:
			d.builder.closeElement(int(decoder.InputOffset()))
		case xml.CharData:
			d.builder.addCharData(int(decoder.InputOffset()), actual)
		}

		prevOffset = int(decoder.InputOffset())
	}
}

func (d *DOM) apply(options []option.Option) {
	for _, opt := range options {
		switch actual := opt.(type) {
		case option.BufferSize:
			d.bufferSize = int(actual)
		case *option.Filters:
			d.builder.filters = actual
		case AttributesChangesSize:
			d.attributesChangesSize = int(actual)
		case ElementsChangesSize:
			d.elementsChangesSize = int(actual)
		}
	}
}

func (d *DOM) templateSlice(span *span) string {
	return string(d.template[span.start:span.end])
}
