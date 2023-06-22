package html

import (
	"bytes"
	"github.com/viant/dm/option"
	"golang.org/x/net/html"
	"io"
)

type (
	//DOM represents DOM structure
	DOM struct {
		template          []byte
		initialBufferSize int
		filter            *option.Filters
		*builder
	}
)

func (d *DOM) apply(options []option.Option) {
	for _, opt := range options {
		switch actual := opt.(type) {
		case option.BufferSize:
			d.initialBufferSize = int(actual)
		case *option.Filters:
			d.filter = actual
		}
	}
}

//New parses template and creates new DOM. Filters can be specified to index some tags and attributes.
func New(template string, options ...option.Option) (*DOM, error) {
	templateBytes := []byte(template)
	d := &DOM{
		template: templateBytes,
	}

	domBuilder := newBuilder(d)
	d.builder = domBuilder

	d.apply(options)

	if err := d.buildTemplate(templateBytes); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DOM) buildTemplate(template []byte) error {
	node := html.NewTokenizer(bytes.NewReader(template))
	offset := 0
outer:
	for {
		previousSpan := rawSpan(node)
		next := node.Next()
		nextSpan := rawSpan(node)

		if nextSpan.start < previousSpan.end {
			offset += previousSpan.end
		}

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
			nodeSpan := dataSpan(node)
			tagName, _ := node.TagName()
			if d.filter != nil {
				if _, ok := d.filter.ElementFilter(string(tagName), false); !ok {
					continue outer
				}
			}

			d.builder.newTag(string(tagName), offset+nextSpan.end, nodeSpan, html.SelfClosingTagToken == next)
			if d.filter == nil {
				buildAllAttributes(template, node, d.builder, offset)
			} else {
				buildFilteredAttributes(template, tagName, node, d.builder, d.filter, offset)
			}
		case html.EndTagToken:
			tagName, _ := node.TagName()
			if d.filter != nil {
				if _, ok := d.filter.ElementFilter(string(tagName), false); !ok {
					continue outer
				}
			}
			d.builder.closeTag(nextSpan.start, tagName)
		}
	}
	return nil
}

func buildAllAttributes(template []byte, z *html.Tokenizer, builder *builder, withOffset int) {
	attributes := attributesSpan(z)
	replaceWithSingleQuote(template, attributes, withOffset)

	for _, attribute := range attributes {
		builder.attribute(attribute, withOffset)
	}
}

func replaceWithSingleQuote(template []byte, attributes [][2]span, withOffset int) {
	for _, attribute := range attributes {
		keyEnd := attribute[0].end
		attrValue := attribute[1]

		var foundEqual bool
		for ; keyEnd < attrValue.start; keyEnd++ {
			if template[keyEnd] == '=' {
				foundEqual = true
				break
			}
		}

		if !foundEqual {
			continue
		}

		template[withOffset+attrValue.start-1] = '"'
		template[withOffset+attrValue.end] = '"'
	}
}

func buildFilteredAttributes(template []byte, tagName []byte, z *html.Tokenizer, builder *builder, tagFilter *option.Filters, withOffset int) {
	attributes := attributesSpan(z)
	replaceWithSingleQuote(template, attributes, withOffset)

	attributeFilter, ok := tagFilter.ElementFilter(string(tagName), false)
	if !ok {
		return
	}

	for _, attribute := range attributes {
		if ok := attributeFilter.Matches(string(template[attribute[0].start:attribute[0].end])); !ok {
			continue
		}
		builder.attribute(attribute, withOffset)
	}
}
