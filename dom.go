package dm

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type (
	DOM struct {
		template          []byte
		initialBufferSize int
		builder           *elementsBuilder
		filter            Filter
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

func (d *DOM) rebuildTemplate(newInnerHTML []byte) (int, error) {
	buffer := NewBuffer(len(d.template))
	buffer.appendBytes(d.template)

	diff := buffer.insertBytes(d.builder.tags[d.builder.lastTag()].InnerHTML, 0, 0, newInnerHTML)
	for i := 1; i < len(d.builder.tags); i++ {
		d.builder.tags[i].InnerHTML.End += diff
	}

	if err := d.buildTemplate(newInnerHTML); err != nil {
		return diff, err
	}

	d.template = buffer.bytes()
	return diff, nil
}

func NewDOM(template []byte, attributes Filter, options ...Option) (*DOM, error) {
	domBuilder := newBuilder()
	d := &DOM{
		template: template,
		builder:  domBuilder,
		filter:   attributes,
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
			d.builder.newTag(rawSpan(node).End, dataSpan(node), html.SelfClosingTagToken == next)
			if d.filter == nil {
				buildAllAttributes(node, d.builder)
			} else {
				buildFilteredAttributes(template, node, d.builder, d.filter)
			}
			d.builder.attributesBuilt()
		case html.EndTagToken:
			d.builder.closeTag(rawSpan(node).Start)
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

func newInnerDOM(dom *DOM, e *elementsBuilder, template []byte, d *DOM) *DOM {
	return &DOM{
		template:          template,
		initialBufferSize: dom.initialBufferSize,
		builder:           e,
		filter:            d.filter,
	}
}
