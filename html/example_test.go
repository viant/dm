package html

import (
	"fmt"
	"github.com/viant/dm/option"
)

func ExampleNew() {
	template := `
<!DOCTYPE html>
<html lang="en">
<head>
	<title>Index</title>
</head>
<body>
	<p class="[class]">Classes</p>
	<img src="[src]" alt="alt"/>
	<div hidden="[hidden]">This is div inner</div>
</body>
</html>`

	dom, err := New(template)
	if err != nil {
		fmt.Println(err)
		return
	}

	filter := option.NewFilters(
		option.NewFilter("div", "hidden"),
		option.NewFilter("img", "src"),
	)

	bufferSize := option.BufferSize(1024)
	document := dom.Document(filter, bufferSize)

	elemIt := document.Select("div", "hidden")
	for elemIt.Has() {
		elem, _ := elemIt.Next()
		fmt.Println(elem.InnerHTML())
		_ = elem.SetInnerHTML("This will be new InnerHTML")
		attribute, ok := elem.MatchAttribute("hidden", "[hidden]")
		if ok {
			attribute.Set("true")
			fmt.Println(attribute.Value())
		}
	}

	attributeIt := document.SelectAttributes("img", "src", "[src]")
	for attributeIt.Has() {
		attribute, _ := attributeIt.Next()
		attribute.Set("abcdef.jpg")
		fmt.Println(attribute.Value())
	}

	fmt.Println(document.Render())

	// Output:
	//This is div inner
	//true
	//abcdef.jpg
	//
	//<!DOCTYPE html>
	//<html lang="en">
	//<head>
	//	<title>Index</title>
	//</head>
	//<body>
	//	<p class="[class]">Classes</p>
	//	<img src="abcdef.jpg" alt="alt"/>
	//	<div hidden="true">This will be new InnerHTML</div>
	//</body>
	//</html>
}
