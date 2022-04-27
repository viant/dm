## Usage

In order to change to update the DOM, you need to create `DOM` representation. `DOM` representation can be shared across the app: 
```go
template := []byte("<html>...</html>")
dom, err := html.New(template)
// handle error
```
 

You can specify some options while creating `DOM`:
* `BufferSize` - initial buffer size for each `DOM` session, `int` wrapper
* `*Filters` - represents allowed tags and attributes

```go
bufferSize := option.BufferSize(1024)
filter := option.NewFilters(
	option.NewFilter("div", "class"), 
	option.NewFilter("img", "src"),
	)
dom, err := option.New(template, bufferSize, filter)
// handle error
```

Then you need to create a `Document`:
```go
document := dom.Document()
documentWithBuffer := dom.Document(option.NewBuffer(1024))
```
If you don't provide a `Buffer`, will be created one with `BufferSize` specified while creating `DOM`

Now you can get/set Attribute, get/set InnerHTML by using selectors. 

Selectors order: `Tag` -> `Attribute` -> `Attribute value`. Selectors are optional, it means if you don't specify `Attribute` only 
tag will be checked. 

Usage:

```go
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
```

## Options
Options supported:
* `*Buffer` - you can reuse Buffers by passing one while creating the `DOM` using `VirutalDOM.DOM()` manually, or by using `Pool` created via `NewPool(size int, dom *VirtualDOM)`