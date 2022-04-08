package dm

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"os"
	"path"
	"testing"
)

func TestDOM(t *testing.T) {
	testLocation := toolbox.CallerDirectory(3)

	type attrSearch struct {
		attribute string
		oldValue  string
		newValue  string
		fullMatch bool
		tag       string
	}

	type innerHTMLSearch struct {
		attribute string
		tag       string
		value     string
	}

	testcases := []struct {
		description   string
		uri           string
		attributes    Filter
		newAttributes []attrSearch
		innerHTMLGet  []innerHTMLSearch
		innerHTMLSet  []innerHTMLSearch
	}{
		{
			uri: "template001",
		},
		{
			uri: "template002",
			newAttributes: []attrSearch{
				{tag: "img", attribute: "src", oldValue: "[src]", newValue: "abcdef", fullMatch: true},
			},
		},
		{
			uri: "template003",
			newAttributes: []attrSearch{
				{tag: "img", attribute: "src", oldValue: "[src]", newValue: "abcdef", fullMatch: true},
				{tag: "p", attribute: "class", oldValue: "[class]", newValue: "newClasses", fullMatch: true},
				{tag: "div", attribute: "hidden", oldValue: "[hidden]", newValue: "newHidden", fullMatch: true},
			},
			innerHTMLGet: []innerHTMLSearch{
				{tag: "div", attribute: "hidden", value: `This is div inner`},
			},
			innerHTMLSet: []innerHTMLSearch{
				{tag: "div", attribute: "hidden", value: `<iframe hidden><div class="div-class">Hello world</div></iframe>`},
			},
		},
		{
			uri: "template004",
			newAttributes: []attrSearch{
				{tag: "img", attribute: "src", oldValue: "[src]", newValue: "newSrc", fullMatch: true},
				{tag: "img", attribute: "src", oldValue: "newSrc", newValue: "abcdef", fullMatch: true},
			},
			innerHTMLGet: []innerHTMLSearch{
				{tag: "head", value: `
    <title>Index</title>
`},
			},
		},
	}

	for _, testcase := range testcases {
		templatePath := path.Join(testLocation, "testdata", testcase.uri)
		template, err := os.ReadFile(path.Join(templatePath, "index.html"))
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}
		dom, err := New(template, testcase.attributes)
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		session := dom.Template()
		for _, newAttr := range testcase.newAttributes {
			attrIterator := session.SelectAttributes([]byte(newAttr.tag), []byte(newAttr.attribute), []byte(newAttr.oldValue))
			for attrIterator.HasMore() {
				attr, _ := attrIterator.Next()
				attr.Set([]byte(newAttr.newValue))
				assert.Equal(t, newAttr.newValue, string(attr.Value()), testcase.uri)
			}
		}

		for _, search := range testcase.innerHTMLGet {
			selectors := make([][]byte, 0)
			if search.tag != "" {
				selectors = append(selectors, []byte(search.tag))
			}

			if search.attribute != "" {
				selectors = append(selectors, []byte(search.attribute))
			}

			tagIt := session.SelectTags(selectors...)
			for tagIt.HasMore() {
				element, _ := tagIt.Next()
				assertly.AssertValues(t, search.value, string(element.InnerHTML()))
			}
		}

		for _, search := range testcase.innerHTMLSet {
			selectors := make([][]byte, 0)
			if search.tag != "" {
				selectors = append(selectors, []byte(search.tag))
			}

			if search.attribute != "" {
				selectors = append(selectors, []byte(search.attribute))
			}

			selectorIt := session.SelectTags(selectors...)
			for selectorIt.HasMore() {
				element, _ := selectorIt.Next()
				_ = element.SetInnerHTML([]byte(search.value))
			}
		}

		result, err := os.ReadFile(path.Join(templatePath, "expect.html"))
		if !assert.Nil(t, err, testcase.uri) {
			t.Fail()
			continue
		}
		bytes := session.Bytes()
		assertly.AssertValues(t, bytes, result, testcase.uri)
	}
}
