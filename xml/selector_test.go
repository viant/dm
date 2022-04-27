package xml_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/dm/xml"
	"github.com/viant/toolbox"
	"testing"
)

func TestXpathSelectors(t *testing.T) {
	testcases := []struct {
		description string
		xpath       string
		output      string
	}{
		{
			xpath:  `foo[test='true']/address`,
			output: `[{"Name":"foo","Attributes":[{"Name":"test","Value":"true","Compare":"="}],"MatchAny":true},{"Name":"address", "MatchAny":false}]`,
		},
		{
			xpath:  `/foo[test='true']/address`,
			output: `[{"Name":"foo","Attributes":[{"Name":"test","Value":"true","Compare":"="}],"MatchAny":false},{"Name":"address","MatchAny":false}]`,
		},
	}

	for _, testcase := range testcases {
		selectors, err := xml.NewSelectors(testcase.xpath)
		if !assert.Nil(t, err, testcase.description) {
			continue
		}

		if !assertly.AssertValues(t, testcase.output, selectors) {
			toolbox.Dump(selectors)
		}
	}
}
