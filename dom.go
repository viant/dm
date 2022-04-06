package dm

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type (
	DOM struct {
		attributes        attrs
		template          []byte
		initialBufferSize int
		tags              tags
	}

	Filter map[string]map[string]bool
)

func (d *DOM) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			d.initialBufferSize = int(actual)
		}
	}
}

func NewDOM(template []byte, attributes Filter, options ...Option) (*DOM, error) {
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
			nodeBuilder.newTag(rawSpan(node).End, dataSpan(node), html.SelfClosingTagToken == next)
			if attributes == nil {
				buildAllAttributes(node, nodeBuilder)
			} else {
				buildFilteredAttributes(template, node, nodeBuilder, attributes)
			}
			nodeBuilder.attributesBuilt()
		case html.EndTagToken:
			nodeBuilder.closeTag(rawSpan(node).Start)
		}
	}

	d := &DOM{
		attributes: nodeBuilder.attributes,
		template:   template,
		tags:       nodeBuilder.tags,
	}
	d.apply(options)
	return d, nil
}

func buildAllAttributes(z *html.Tokenizer, builder *elementsBuilder) {
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		builder.attribute(attribute)
	}
}

func buildFilteredAttributes(template []byte, z *html.Tokenizer, builder *elementsBuilder, tagFilter Filter) {
	tagName, _ := z.TagName()
	var ok bool
	var attributeFilter map[string]bool
	if attributeFilter, ok = tagFilter[string(tagName)]; !ok {
		return
	}

	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		if _, ok := attributeFilter[string(template[attribute[0].Start:attribute[0].End])]; !ok {
			continue
		}
		builder.attribute(attribute)
	}
}
