##Introduction
The package `xml` of `dm` works very similar to the  `html` package but for templates in the xml format. 

## Usage
Similarly to the `html` in order to manipulate xml, first you need to create `VirtualDOM`, one can be shared across the app:

```go
template := []byte("<?xml version=...")
vdom, err := xml.New(template)
// handle error
```

```go
bufferSize := xml.BufferSize(1024)
filter := option.NewFilters(
	option.NewFilter("foo", "attr1", "attr2"), 
	option.NewFilter("name", "attr1", "attr2"),
	)
vdom, err := xml.New(template, bufferSize, filter)
// handle error
```

Then you need to create a `DOM`:
```go
dom := vdom.DOM()
```

Now you can get/set Attribute, get/set Value by using selectors. 
```go
elementSelector := xml.Selector{
	Name: "foo",
	Attributes: []AttributeSelector{
		Name: "id", 
		Value: "1",
	}
}
```

Usage:
```go
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

dom := vdom.DOM()

elemIt := dom.Select(xml.Selector{Name: "foo"}, xml.Selector{Name: "id"})
for elemIt.Has() {
    elem, _ := elemIt.Next()
    elem.SetValue("10")
}

elemIt = dom.Select(xml.Selector{Name: "foo"}, xml.Selector{Name: "address"})
for elemIt.Has() {
    elem, _ := elemIt.Next()
    elem.SetValue("")
    elem.AddElement("<new-elem>New element value</new-elem>")
}

elemIt = dom.Select(xml.Selector{Name: "foo", Attributes: []xml.AttributeSelector{{Name: "test", Value: "true"}}})
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
```

## Options
Supported options:
* `ElementsChangesSize` - in case of lazy rendering all the changes are buffered. In order to ignore lookup time with the Map,
it is possible to update changes directly in the slice. Default is 30
* `AttributesChangesSize` - same as above, but for the attributes. Default is 30

## Filters
You can create filters parsing xpath:
```go
filters := option.NewFilters()
newFilters, err := xml.NewFilters("foo/price[currency]", "address[country and city]/street")
// handle error
filters.Add(newFilters...)
```

Or
```go
filters, err := xml.FiltersOf("foo/price[currency]", "address[country and city]/street")
// handle error
```
