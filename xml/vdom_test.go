package xml_test

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/dm/option"
	"github.com/viant/dm/xml"
	"github.com/viant/toolbox"
	"os"
	"path"
	"strconv"
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
		after     bool
		before    bool
	}

	type attributeChanges struct {
		selectors []xml.Selector
		keys      []string
		values    []string
	}

	type newValue struct {
		selectors []xml.Selector
		values    []string
	}

	var testcases = []struct {
		description       string
		uri               string
		valuesSearch      []valuesSearch
		attributesSearch  []attributeSearch
		attributesChanges []attributeChange
		addedElements     []newElement
		insertedElements  []newElement
		newAttributes     []attributeChanges
		setAttributes     []attributeChanges
		newValues         []newValue
		replaced          []newValue
		filters           *option.Filters
		xpaths            []string
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
					selectors: []xml.Selector{{Name: "foo"}, {Name: "id"}},
					expected:  []string{"1"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "name"}},
					expected:  []string{"foo name"},
				},
			},
			attributesSearch: []attributeSearch{
				{
					selectors: []xml.Selector{{Name: "foo"}},
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
					selectors: []xml.Selector{{Name: "foo"}},
					attribute: "id",
					newValues: []string{"123"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}},
					attribute: "name",
					newValues: []string{"foo name changed"},
				},
			},
		},
		{
			uri:         "xml004",
			description: "add new element",
			addedElements: []newElement{
				{
					selectors: []xml.Selector{{
						Name:       "foo",
						Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}},
					}},
					values: []string{"<price>123.5</price>"},
				},
				{
					selectors: []xml.Selector{{
						Name:       "foo",
						Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}},
					}},
					values: []string{"<quantity>550</quantity>"},
				},
			},
		},
		{
			uri:         "xml005",
			description: "add new attribute",
			newAttributes: []attributeChanges{
				{
					selectors: []xml.Selector{{
						Name:       "foo",
						Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}},
					}},
					keys:   []string{"price"},
					values: []string{"123"},
				},
				{
					selectors: []xml.Selector{{
						Name:       "foo",
						Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}},
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
					selectors: []xml.Selector{{Name: "foo"}, {Name: "id"}},
					values:    []string{"2"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "name"}},
					values:    []string{"foo"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "address"}},
					values:    []string{""},
				},
			},
		},
		{
			uri:         "xml006",
			description: "update element value with filters",
			newValues: []newValue{
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "id"}},
					values:    []string{"2"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "name"}},
					values:    []string{"foo"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "address"}},
					values:    []string{""},
				},
			},
			filters: option.NewFilters(
				option.NewFilter("foo", "test"),
				option.NewFilter("id"),
				option.NewFilter("name"),
				option.NewFilter("address"),
			),
		},
		{
			uri:         "xml007",
			description: "new value, access attributes via map, with filters",
			newValues: []newValue{
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "prop1", Attributes: []xml.AttributeSelector{{Name: "attr1", Value: "abc"}}}},
					values:    []string{"50"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "prop3", Attributes: []xml.AttributeSelector{{Name: "attr4", Value: "jkl"}}}},
					values:    []string{"125"},
				},
			},
			filters: option.NewFilters(
				option.NewFilter("foo"),
				option.NewFilter("prop4", "attr1"),
			),
			xpaths: []string{"foo[test]/prop1[attr1 and attr2]", "prop3[attr4 and attr3]", "prop5[attr1 and attr2]"},
		},
		{
			uri:         "xml007",
			description: "new value, access attributes via map, without filters",
			newValues: []newValue{
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "prop1", Attributes: []xml.AttributeSelector{{Name: "attr1", Value: "abc"}}}},
					values:    []string{"50"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "prop3", Attributes: []xml.AttributeSelector{{Name: "attr4", Value: "jkl"}}}},
					values:    []string{"125"},
				},
			},
		},
		{
			uri:         "xml008",
			description: "update element value without filters",
			insertedElements: []newElement{
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "address"}, {Name: "zip-code"}},
					before:    true,
					values:    []string{"<country-code>FR</country-code>"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "address"}},
					values:    []string{"<type>fType</type>"},
					after:     true,
				},
			},
		},
		{
			uri:         "xml009",
			description: "element -> SetAttribute",
			setAttributes: []attributeChanges{
				{
					selectors: []xml.Selector{{Name: "foo"}},
					keys:      []string{"test"},
					values:    []string{"true"},
				},
				{
					selectors: []xml.Selector{{Name: "foo"}},
					keys:      []string{"newAttr"},
					values:    []string{"125"},
				},
			},
		},
		{
			uri:         "xml010",
			description: "element -> ReplaceWith",
			replaced: []newValue{
				{
					selectors: []xml.Selector{{Name: "foo"}, {Name: "address"}},
					values:    []string{"<type>fType</type>"},
				},
			},
		},
	}
	//for i, testcase := range testcases[len(testcases)-1:] {
	for i, testcase := range testcases {
		fmt.Println("Running testcase: " + strconv.Itoa(i))
		templatePath := path.Join(testLocation, "testdata", testcase.uri)

		filters, err := xml.NewFilters(testcase.xpaths...)
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		testcase.filters.Add(filters...)

		dom, err := readFromFile(path.Join(templatePath, "index.xml"), testcase.filters)
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		document := dom.Document()

		for _, search := range testcase.valuesSearch {
			it := document.Select(search.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				assert.Equal(t, search.expected[counter], element.Value(), testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(search.expected), testcase.description)
		}

		for _, search := range testcase.attributesSearch {
			it := document.Select(search.selectors...)
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
			it := document.Select(attributesMutation.selectors...)
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

		for _, newElem := range testcase.addedElements {
			it := document.Select(newElem.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.AddElement(newElem.values[counter])
				counter++
			}

			assert.Equal(t, counter, len(newElem.values), testcase.description)
		}

		for _, newAttr := range testcase.newAttributes {
			it := document.Select(newAttr.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.AddAttribute(newAttr.keys[counter], newAttr.values[counter])
				counter++
			}

			assert.Equal(t, counter, len(newAttr.values), testcase.description)
		}

		for _, newVal := range testcase.newValues {
			it := document.Select(newVal.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.SetValue(newVal.values[counter])
				assert.Equal(t, newVal.values[counter], element.Value(), testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(newVal.values), testcase.description)
		}

		for _, insertedEl := range testcase.insertedElements {
			it := document.Select(insertedEl.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				if insertedEl.before {
					element.InsertBefore(insertedEl.values[counter])
				}

				if insertedEl.after {
					element.InsertAfter(insertedEl.values[counter])
				}

				counter++
			}

			assert.Equal(t, counter, len(insertedEl.values), testcase.description)
		}

		for _, setAttribute := range testcase.setAttributes {
			it := document.Select(setAttribute.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.SetAttribute(setAttribute.keys[counter], setAttribute.values[counter])
				counter++
			}

			assert.Equal(t, counter, len(setAttribute.values), testcase.description)
		}

		for _, newVal := range testcase.replaced {
			it := document.Select(newVal.selectors...)
			counter := 0
			for it.Has() {
				element, _ := it.Next()
				element.ReplaceWith(newVal.values[counter])
				assert.Equal(t, newVal.values[counter], element.Value(), testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(newVal.values), testcase.description)
		}

		result, err := os.ReadFile(path.Join(templatePath, "expect.txt"))
		if !assert.Nil(t, err) {
			return
		}

		render := document.Render()
		assert.Equal(t, string(result), render, testcase.description)
	}
}

func readFromFile(path string, filters *option.Filters) (*xml.DOM, error) {
	template, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	schema, err := xml.New(string(template), filters)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

//Benchmarks
//go:embed testdata/xml006/index.xml
var benchTemplate string
var benchVDOM *xml.DOM

func init() {
	bufferSize := option.BufferSize(1024)
	filters := option.NewFilters(
		option.NewFilter("foo", "test"),
		option.NewFilter("id"),
		option.NewFilter("name"),
		option.NewFilter("address"),
	)
	benchVDOM, _ = xml.New(benchTemplate, bufferSize, filters)
}

func BenchmarkXml_Render(b *testing.B) {
	b.ReportAllocs()
	var aXml *xml.Document
	for i := 0; i < b.N; i++ {
		aXml = benchVDOM.Document()
		elemIt := aXml.Select(xml.Selector{Name: "foo"}, xml.Selector{Name: "id"})
		for elemIt.Has() {
			elem, _ := elemIt.Next()
			elem.SetValue("10")
		}

		elemIt = aXml.Select(xml.Selector{Name: "foo"}, xml.Selector{Name: "address"})
		for elemIt.Has() {
			elem, _ := elemIt.Next()
			elem.SetValue("")
			elem.AddElement("<new-elem>New element value</new-elem>")
		}

		elemIt = aXml.Select(xml.Selector{Name: "foo", Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}}})
		for elemIt.Has() {
			elem, _ := elemIt.Next()
			elem.AddAttribute("attr1", "value1")
			attribute, ok := elem.Attribute("test")
			if !ok {
				continue
			}
			attribute.Set("new test value")
		}
	}

	assert.Equal(b, `<?xml version="1.0" encoding="UTF-8"?>
<foo test="new test value" attr1="value1">
    <id>10</id>
    <name>foo name</name>
    <address><new-elem>New element value</new-elem></address>
    <quantity>123</quantity>
    <price>50.5</price>
    <type>fType</type>
</foo>`, aXml.Render())
}
