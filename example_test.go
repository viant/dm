package dm

import "fmt"

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

	vdom, err := New(template)
	if err != nil {
		fmt.Println(err)
		return
	}

	filter := NewFilter(
		NewTagFilter("div", "hidden"),
		NewTagFilter("img", "src"),
	)

	bufferSize := BufferSize(1024)
	dom := vdom.DOM(filter, bufferSize)

	elemIt := dom.Select("div", "hidden")
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

	attributeIt := dom.SelectAttributes("img", "src", "[src]")
	for attributeIt.Has() {
		attribute, _ := attributeIt.Next()
		attribute.Set("abcdef.jpg")
		fmt.Println(attribute.Value())
	}

	fmt.Println(dom.Render())

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
