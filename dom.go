package dm

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type DOM struct {
	attributes        attrs
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

func NewDOM(template []byte, attributes []string, options ...Option) (*DOM, error) {
	var attrAsMap map[string]bool
	if len(attributes) > 0 {
		attrAsMap = map[string]bool{}
		for _, attribute := range attributes {
			attrAsMap[attribute] = true
		}
	}

	node := html.NewTokenizer(bytes.NewReader(template))
	nodeBuilder := newBuilder()
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
			if attrAsMap == nil {
				buildAllAttributes(node, nodeBuilder)
			} else {
				buildFilteredAttributes(template, node, nodeBuilder, attrAsMap)
			}
		}

	}

	d := &DOM{
		attributes: nodeBuilder.result(),
		template:   template,
	}
	d.apply(options)
	return d, nil
}

func buildAllAttributes(z *html.Tokenizer, builder *attributesBuilder) {
	tagName, _ := z.TagName()
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		builder.attribute(tagName, attribute)
	}
}

func buildFilteredAttributes(template []byte, z *html.Tokenizer, builder *attributesBuilder, allowedAttributes map[string]bool) {
	tagName, _ := z.TagName()
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		if _, ok := allowedAttributes[string(template[attribute[0].Start:attribute[0].End])]; ok {
			continue
		}
		builder.attribute(tagName, attribute)
	}
}
