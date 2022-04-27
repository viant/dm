package xml_test

import (
	"fmt"
	"github.com/viant/dm/option"
	"github.com/viant/dm/xml"
	"strings"
)

func ExampleNew() {
	template := `
<?xml version="1.0" encoding="UTF-8"?>
<foo test="true">
    <id>1</id>
    <name>foo name</name>
    <address>
        <street>abc</street>
        <zip-code>123456</zip-code>
        <country>
            <id>1</id>
            <name>def</name>
        </country>
    </address>
    <quantity>123</quantity>
    <price>50.5</price>
    <type>fType</type>
</foo>`

	filters := option.NewFilters(
		option.NewFilter("foo", "test"),
		option.NewFilter("id"),
		option.NewFilter("name"),
		option.NewFilter("address"),
	)

	vdom, err := xml.New(template, filters)
	if err != nil {
		fmt.Println(err)
		return
	}

	dom := vdom.Document()

	elem, ok := dom.SelectFirst(xml.Selector{Name: "foo"}, xml.Selector{Name: "id"})
	if ok {
		elem.SetValue("10")
	}

	elem, ok = dom.SelectFirst(xml.Selector{Name: "foo"}, xml.Selector{Name: "address"})
	if ok {
		elem.SetValue("")
		elem.AddElement("<new-elem>New element value</new-elem>")
	}

	elemIt := dom.Select(xml.Selector{Name: "foo", Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}}})
	for elemIt.Has() {
		elem, _ := elemIt.Next()
		elem.AddAttribute("attr1", "value1")
		attribute, ok := elem.Attribute("test")
		if !ok {
			continue
		}
		attribute.Set(strings.ToUpper(attribute.Value()))
	}

	result := dom.Render()
	fmt.Println(result)

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <foo test="TRUE" attr1="value1">
	//     <id>10</id>
	//     <name>foo name</name>
	//     <address><new-elem>New element value</new-elem></address>
	//     <quantity>123</quantity>
	//     <price>50.5</price>
	//     <type>fType</type>
	// </foo>
}
