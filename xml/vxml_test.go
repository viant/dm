package xml_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/dm/xml"
	"github.com/viant/toolbox"
	"os"
	"path"
	"testing"
)

func TestNew(t *testing.T) {
	testLocation := toolbox.CallerDirectory(3)

	type valuesSearch struct {
		selectors []xml.Selector
		expected  []string
	}

	type attributeSearch struct {
		selectors []xml.Selector
		attribute string
		expected  []string
	}

	type attributeChange struct {
		selectors []xml.Selector
		attribute string
		newValues []string
	}

	type newElement struct {
		selectors []xml.Selector
		values    []string
	}

	type newAttribute struct {
		selectors []xml.Selector
		keys      []string
		values    []string
	}

	type newValue struct {
		selectors []xml.Selector
		values    []string
	}

	testcases := []struct {
		description       string
		uri               string
		valuesSearch      []valuesSearch
		attributesSearch  []attributeSearch
		attributesChanges []attributeChange
		newElements       []newElement
		newAttributes     []newAttribute
		newValues         []newValue
	}{
		{
			uri:         "xml001",
			description: "xml001",
		},
		{
			uri:         "xml002",
			description: "xml002",
			valuesSearch: []valuesSearch{
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.ElementSelector("id")},
					expected:  []string{"1"},
				},
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.ElementSelector("name")},
					expected:  []string{"foo name"},
				},
			},
			attributesSearch: []attributeSearch{
				{
					selectors: []xml.Selector{xml.ElementSelector("foo")},
					attribute: "test",
					expected:  []string{"true"},
				},
			},
		},
		{
			uri:         "xml003",
			description: "xml003",
			attributesChanges: []attributeChange{
				{
					selectors: []xml.Selector{xml.ElementSelector("foo")},
					attribute: "id",
					newValues: []string{"123"},
				},
				{
					selectors: []xml.Selector{xml.ElementSelector("foo")},
					attribute: "name",
					newValues: []string{"foo name changed"},
				},
			},
		},
		{
			uri:         "xml004",
			description: "xml004",
			newElements: []newElement{
				{
					selectors: []xml.Selector{xml.AttributeSelector{
						Name:  "test",
						Value: "true",
					}},
					values: []string{"<price>123.5</price>"},
				},
				{
					selectors: []xml.Selector{xml.AttributeSelector{
						Name:  "test",
						Value: "true",
					}},
					values: []string{"<quantity>550</quantity>"},
				},
			},
		},
		{
			uri:         "xml005",
			description: "xml005",
			newAttributes: []newAttribute{
				{
					selectors: []xml.Selector{xml.AttributeSelector{
						Name:  "test",
						Value: "true",
					}},
					keys:   []string{"price"},
					values: []string{"123"},
				},
				{
					selectors: []xml.Selector{xml.AttributeSelector{
						Name:  "test",
						Value: "true",
					}},
					keys:   []string{"quantity"},
					values: []string{"50.5"},
				},
			},
		},
		{
			uri:         "xml006",
			description: "xml006",
			newValues: []newValue{
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.ElementSelector("id")},
					values:    []string{"2"},
				},
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.ElementSelector("name")},
					values:    []string{"foo"},
				},
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.ElementSelector("address")},
					values:    []string{""},
				},
			},
		},
	}

	//for _, testcase := range testcases[len(testcases)-1:] {
	for _, testcase := range testcases {
		templatePath := path.Join(testLocation, "testdata", testcase.uri)
		vxml, err := readFromFile(path.Join(templatePath, "index.xml"))
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		xml := vxml.Xml()

		for _, search := range testcase.valuesSearch {
			it := xml.Select(search.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				assert.Equal(t, search.expected[counter], element.Value(), testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(search.expected), testcase.description)
		}

		for _, search := range testcase.attributesSearch {
			it := xml.Select(search.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				attribute, ok := element.Attribute(search.attribute)
				assert.True(t, ok, testcase.description)
				assert.Equal(t, search.expected[counter], attribute.Value(), testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(search.expected), testcase.description)
		}

		for _, attributesMutation := range testcase.attributesChanges {
			it := xml.Select(attributesMutation.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				attribute, ok := element.Attribute(attributesMutation.attribute)
				assert.True(t, ok, testcase.description)
				attribute.Set(attributesMutation.newValues[counter])
				assert.Equal(t, attribute.Value(), attributesMutation.newValues[counter], testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(attributesMutation.newValues), testcase.description)
		}

		for _, newElem := range testcase.newElements {
			it := xml.Select(newElem.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.AddElement(newElem.values[counter])
				counter++
			}

			assert.Equal(t, counter, len(newElem.values), testcase.description)
		}

		for _, newAttr := range testcase.newAttributes {
			it := xml.Select(newAttr.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.AddAttribute(newAttr.keys[counter], newAttr.values[counter])
				counter++
			}

			assert.Equal(t, counter, len(newAttr.values), testcase.description)
		}

		for _, newVal := range testcase.newValues {
			it := xml.Select(newVal.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.Set(newVal.values[counter])
				counter++
			}

			assert.Equal(t, counter, len(newVal.values), testcase.description)
		}

		result, err := os.ReadFile(path.Join(templatePath, "expect.xml"))
		if !assert.Nil(t, err) {
			return
		}

		render := xml.Render()
		assert.Equal(t, string(result), render, testcase.description)
	}
}

func readFromFile(path string) (*xml.VirtualXml, error) {
	template, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	dom, err := xml.New(string(template))
	if err != nil {
		return nil, err
	}
	return dom, nil
}
