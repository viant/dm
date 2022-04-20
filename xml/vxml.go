package xml

import (
	"bytes"
	"encoding/xml"
)

type VirtualXml struct {
	template []byte
	*builder
	bufferSize int
}

func New(template string, options ...Option) (*VirtualXml, error) {
	vxml := &VirtualXml{
		template:   []byte(template),
		bufferSize: len(template),
	}

	vxml.builder = newBuilder(vxml)
	vxml.apply(options)

	if err := vxml.Init(); err != nil {
		return nil, err
	}

	return vxml, nil
}

func (v *VirtualXml) Init() error {
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

func (v *VirtualXml) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			v.bufferSize = int(actual)
		}
	}
}
