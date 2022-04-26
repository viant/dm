package xml_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/dm/xml"
	"github.com/viant/toolbox"
	"testing"
)

func TestNewFilters(t *testing.T) {
	testcases := []struct {
		description string
		input       []string
		expected    string
	}{
		{
			input:    []string{`foo[test = 'true']/address`},
			expected: `[{"Name":"foo","Attributes":["test"]},{"Name":"address","Attributes":[]}]`,
		},
		{
			input:    []string{`foo[test='true' and lang='eng']/address`},
			expected: `[{"Name":"foo","Attributes":["test","lang"]},{"Name":"address","Attributes":[]}]`,
		},
	}

	//for _, testcase := range testcases[len(testcases)-1:] {
	for _, testcase := range testcases {
		elementsFilters, err := xml.NewFilters(testcase.input...)
		if !assert.Nil(t, err, testcase.description) {
			continue
		}

		if !assertly.AssertValues(t, testcase.expected, elementsFilters, testcase.description) {
			toolbox.Dump(elementsFilters)
		}
	}
}
