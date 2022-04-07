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

	for _, testcase := range testcases[2:3] {
		templatePath := path.Join(testLocation, "testdata", testcase.uri)
		template, err := os.ReadFile(path.Join(templatePath, "index.html"))
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}
		dom, err := NewDOM(template, testcase.attributes)
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		session := dom.Session()
		for _, newAttr := range testcase.newAttributes {
			session.SetAttr(0, []byte(newAttr.newValue), []byte(newAttr.tag), []byte(newAttr.attribute), []byte(newAttr.oldValue))
			attrVal, _, _ := session.Attribute(0, []byte(newAttr.tag), []byte(newAttr.attribute))
			assert.Equal(t, newAttr.newValue, string(attrVal), testcase.uri)
		}

		for _, search := range testcase.innerHTMLGet {
			selectors := make([][]byte, 0)
			if search.tag != "" {
				selectors = append(selectors, []byte(search.tag))
			}

			if search.attribute != "" {
				selectors = append(selectors, []byte(search.attribute))
			}

			innerHTML, _ := session.InnerHTML(0, selectors...)
			assertly.AssertValues(t, search.value, string(innerHTML))
		}

		for _, search := range testcase.innerHTMLSet {
			selectors := make([][]byte, 0)
			if search.tag != "" {
				selectors = append(selectors, []byte(search.tag))
			}

			if search.attribute != "" {
				selectors = append(selectors, []byte(search.attribute))
			}

			_, _ = session.SetInnerHTML(0, []byte(search.value), selectors...)
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
