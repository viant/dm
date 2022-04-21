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
		filters           *xml.Filters
	}{
		{
			uri:         "xml001",
			description: "raw template without changes",
		},
		{
			uri:         "xml002",
			description: "check element values",
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
			description: "change attribute value",
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
			description: "add new element",
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
			description: "add new attribute",
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
			description: "update element value without filters",
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
		{
			uri:         "xml006",
			description: "update element value with filters",
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
			filters: xml.NewFilters(
				xml.NewFilter("foo"),
				xml.NewFilter("id"),
				xml.NewFilter("name"),
				xml.NewFilter("address"),
			),
		},
		{
			uri:         "xml007",
			description: "new value, access attributes via map, with filters",
			newValues: []newValue{
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.AttributeSelector{Name: "attr1", Value: "abc"}},
					values:    []string{"50", "100"},
				},
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.AttributeSelector{Name: "attr4", Value: "jkl"}},
					values:    []string{"125"},
				},
			},
			filters: xml.NewFilters(
				xml.NewFilter("foo"),
				xml.NewFilter("prop1", "attr1"),
				xml.NewFilter("prop3", "attr4"),
				xml.NewFilter("prop4", "attr1"),
				xml.NewFilter("prop5", "attr1"),
			),
		},
		{
			uri:         "xml007",
			description: "new value, access attributes via map, without filters",
			newValues: []newValue{
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.AttributeSelector{Name: "attr1", Value: "abc"}},
					values:    []string{"50", "100"},
				},
				{
					selectors: []xml.Selector{xml.ElementSelector("foo"), xml.AttributeSelector{Name: "attr4", Value: "jkl"}},
					values:    []string{"125"},
				},
			},
		},
	}

	//for _, testcase := range testcases[len(testcases)-1:] {
	for _, testcase := range testcases {
		templatePath := path.Join(testLocation, "testdata", testcase.uri)
		schema, err := readFromFile(path.Join(templatePath, "index.xml"), testcase.filters)
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		xml := schema.Xml()

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
				element.SetValue(newVal.values[counter])
				assert.Equal(t, newVal.values[counter], element.Value(), testcase.description)
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

func readFromFile(path string, filters *xml.Filters) (*xml.Schema, error) {
	template, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	dom, err := xml.New(string(template), filters)
	if err != nil {
		return nil, err
	}
	return dom, nil
}
