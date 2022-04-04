package vhtml

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type DOM struct {
	attributes        Attributes
	template          []byte
	initialBufferSize int
}

func (d *DOM) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			d.initialBufferSize = int(actual)
		}
	}
}

func NewVDom(template []byte, attributes []string, options ...Option) (*DOM, error) {
	node := html.NewTokenizer(bytes.NewReader(template))
	nodeBuilder := NewBuilder()
outer:
	for {
		next := node.Next()
		switch next {
		case html.ErrorToken:
			err := node.Err()
			if err == nil {
				continue
			}

			if err != io.EOF {
				return nil, err
			}

			break outer

		case html.StartTagToken, html.SelfClosingTagToken:
			if err := buildAttributes(node, nodeBuilder); err != nil {
				return nil, err
			}
		}

	}

	nodes := nodeBuilder.Attributes()

	d := &DOM{
		attributes: nodes,
		template:   template,
	}
	d.apply(options)
	return d, nil
}

func buildAttributes(z *html.Tokenizer, builder *AttributesBuilder) error {
	tagName, _ := z.TagName()
	attributes := AttributesSpan(z)

	for _, attribute := range attributes {
		builder.Attribute(tagName, attribute)
	}

	return nil
}
