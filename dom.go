package dm

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type (
	//DOM represents DOM structure
	DOM struct {
		template          []byte
		initialBufferSize int
		builder           *elementsBuilder
		filter            Filter
	}

	//Filter represents tags and attributes filter
	Filter map[string]map[string]bool
)

func (d *DOM) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			d.initialBufferSize = int(actual)
		case Filter:
			d.filter = actual
		}
	}
}

//AttributesLen returns number of attributes. Attributes[0] is an empty attribute.
func (d *DOM) AttributesLen() int {
	return len(d.builder.attributes)
}

//TagLen returns number of tags. Tag[0] is an empty tag.
func (d *DOM) TagLen() int {
	return len(d.builder.tags)
}

//New parses template and creates new DOM. Filter can be specified to index some tags and attributes.
func New(template []byte, options ...Option) (*DOM, error) {
	domBuilder := newBuilder()
	d := &DOM{
		template: template,
		builder:  domBuilder,
	}
	d.apply(options)

	if err := d.buildTemplate(template); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DOM) buildTemplate(template []byte) error {
	node := html.NewTokenizer(bytes.NewReader(template))
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
				return err
			}

			break outer

		case html.StartTagToken, html.SelfClosingTagToken:
			if d.filter != nil {
				tagName, _ := node.TagName()
				if _, ok := d.filter[string(tagName)]; !ok {
					continue outer
				}
			}

			d.builder.newTag(rawSpan(node).end, dataSpan(node), html.SelfClosingTagToken == next)
			if d.filter == nil {
				buildAllAttributes(node, d.builder)
			} else {
				buildFilteredAttributes(template, node, d.builder, d.filter)
			}
			d.builder.attributesBuilt()
		case html.EndTagToken:
			d.builder.closeTag(rawSpan(node).start)
		}
	}
	return nil
}

func buildAllAttributes(z *html.Tokenizer, builder *elementsBuilder) {
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		builder.attribute(attribute)
	}
}

func buildFilteredAttributes(template []byte, z *html.Tokenizer, builder *elementsBuilder, tagFilter Filter) {
	tagName, _ := z.TagName()
	attributeFilter := tagFilter[string(tagName)]
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		if _, ok := attributeFilter[string(template[attribute[0].start:attribute[0].end])]; !ok {
			continue
		}
		builder.attribute(attribute)
	}
}
