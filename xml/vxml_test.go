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

	testcases := []struct {
		description      string
		uri              string
		valuesSearch     []valuesSearch
		attributesSearch []attributeSearch
	}{
		{
			uri: "xml001",
		},
		{
			uri: "xml002",
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
	}

	for _, testcase := range testcases[len(testcases)-1:] {
		//for _, testcase := range testcases {
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
				assert.Equal(t, search.expected[counter], attribute, testcase.description)
				counter++
			}

			assert.Equal(t, counter, len(search.expected), testcase.description)
		}

		result, err := os.ReadFile(path.Join(templatePath, "expect.xml"))
		if !assert.Nil(t, err) {
			return
		}

		assert.Equal(t, string(result), xml.Render(), testcase.description)
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
