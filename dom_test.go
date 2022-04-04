package vhtml

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

	type newAttribute struct {
		attribute string
		oldValue  string
		newValue  string
		fullMatch bool
		tag       string
	}

	testcases := []struct {
		description   string
		uri           string
		attributes    []string
		newAttributes []newAttribute
	}{
		{
			uri: "template001",
		},
		{
			uri: "template002",
			newAttributes: []newAttribute{
				{tag: "img", attribute: "src", oldValue: "[src]", newValue: "abcdef", fullMatch: true},
			},
		},
		{
			uri: "template003",
			newAttributes: []newAttribute{
				{tag: "img", attribute: "src", oldValue: "[src]", newValue: "abcdef", fullMatch: true},
				{tag: "p", attribute: "class", oldValue: "[class]", newValue: "newClasses", fullMatch: true},
				{tag: "div", attribute: "hidden", oldValue: "[hidden]", newValue: "newHidden", fullMatch: true},
			},
		},
		{
			uri: "template004",
			newAttributes: []newAttribute{
				{tag: "img", attribute: "src", oldValue: "[src]", newValue: "newSrc", fullMatch: true},
				{tag: "img", attribute: "src", oldValue: "newSrc", newValue: "abcdef", fullMatch: true},
			},
		},
	}

	for _, testcase := range testcases[3:4] {
		templatePath := path.Join(testLocation, "testdata", testcase.uri)
		template, err := os.ReadFile(path.Join(templatePath, "index.html"))
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}
		dom, err := NewVDom(template, testcase.attributes)
		if !assert.Nil(t, err, testcase.description) {
			t.Fail()
			continue
		}

		session := dom.Session()
		for _, newAttr := range testcase.newAttributes {
			session.SetAttr([]byte(newAttr.tag), []byte(newAttr.attribute), []byte(newAttr.oldValue), []byte(newAttr.newValue), newAttr.fullMatch)
			attrVal, _ := session.Attribute([]byte(newAttr.tag), []byte(newAttr.attribute))
			assert.Equal(t, newAttr.newValue, string(attrVal), testcase.uri)
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
