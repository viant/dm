package html

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type (
	//VirtualDOM represents VirtualDOM structure
	VirtualDOM struct {
		template          []byte
		initialBufferSize int
		filter            *Filters
		*builder
	}
)

func (v *VirtualDOM) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			v.initialBufferSize = int(actual)
		case *Filters:
			v.filter = actual
		}
	}
}

//AttributesLen returns number of attributes. Attributes[0] is an empty attribute.
func (v *VirtualDOM) AttributesLen() int {
	return len(v.builder.attributes)
}

//New parses template and creates new VirtualDOM. Filters can be specified to index some tags and attributes.
func New(template string, options ...Option) (*VirtualDOM, error) {
	domBuilder := newBuilder()
	templateBytes := []byte(template)
	d := &VirtualDOM{
		template: templateBytes,
		builder:  domBuilder,
	}
	d.apply(options)

	if err := d.buildTemplate(templateBytes); err != nil {
		return nil, err
	}

	return d, nil
}

func (v *VirtualDOM) buildTemplate(template []byte) error {
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
			nodeSpan := dataSpan(node)
			tagName, _ := node.TagName()
			if v.filter != nil {
				if _, ok := v.filter.tagFilter(string(tagName)); !ok {
					continue outer
				}
			}

			v.builder.newTag(string(tagName), rawSpan(node).end, nodeSpan, html.SelfClosingTagToken == next)
			if v.filter == nil {
				buildAllAttributes(template, node, v.builder)
			} else {
				buildFilteredAttributes(template, tagName, node, v.builder, v.filter)
			}
			v.builder.attributesBuilt()
		case html.EndTagToken:
			tagName, _ := node.TagName()
			if v.filter != nil {
				if _, ok := v.filter.tagFilter(string(tagName)); !ok {
					continue outer
				}
			}
			v.builder.closeTag(rawSpan(node).start)
		}
	}
	return nil
}

func buildAllAttributes(template []byte, z *html.Tokenizer, builder *builder) {
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		builder.attribute(template, attribute)
	}
}

func buildFilteredAttributes(template []byte, tagName []byte, z *html.Tokenizer, builder *builder, tagFilter *Filters) {
	attributeFilter, ok := tagFilter.tagFilter(string(tagName))
	if !ok {
		return
	}
	attributes := attributesSpan(z)
	for _, attribute := range attributes {
		if ok := attributeFilter.matches(string(template[attribute[0].start:attribute[0].end])); !ok {
			continue
		}
		builder.attribute(template, attribute)
	}
}
