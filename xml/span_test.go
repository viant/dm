package xml

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"testing"
)

func TestExtractAttributes(t *testing.T) {
	testcases := []struct {
		offset   int
		template string
		output   [][2]*span
	}{
		{
			template: `<elem attr="value">`,
			output: [][2]*span{
				{
					{
						start: 6,
						end:   10,
					},
					{
						start: 12,
						end:   17,
					},
				},
			},
		},

		{
			template: `<elem attr>`,
			output: [][2]*span{
				{
					{
						start: 6,
						end:   9,
					},
					{
						start: 9,
						end:   9,
					},
				},
			},
		},
		{
			template: `<elem attr="10"    xlns:attr="some attribute">`,
			output: [][2]*span{
				{
					{
						start: 6,
						end:   10,
					},
					{
						start: 12,
						end:   14,
					},
				},
				{
					{
						start: 19,
						end:   28,
					},
					{
						start: 30,
						end:   44,
					},
				},
			},
		},
		{
			template: `<elem attr="10"    xlns:attr="some attribute">`,
			offset:   100,
			output: [][2]*span{
				{
					{
						start: 106,
						end:   110,
					},
					{
						start: 112,
						end:   114,
					},
				},
				{
					{
						start: 119,
						end:   128,
					},
					{
						start: 130,
						end:   144,
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		attributes, err := extractAttributes(testcase.offset, []byte(testcase.template))
		assert.Nil(t, err, testcase.template)
		assertly.AssertValues(t, testcase.output, attributes, testcase.template)
	}
}
