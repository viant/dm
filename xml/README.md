##Introduction
The subpackage `xml` of `dm` works very similar to the `dm` but for templates in the xml format. 

## Usage
Similarly to the `dm` in order to manipulate xml, first you need to create `Schema`, schema can be shared across the app:

```go
template := []byte("<?xml version=...")
schema, err := xml.New(template)
// handle error
```

```go
bufferSize := xml.BufferSize(1024)
filter := xml.NewFilters(
	xml.NewTagFilter("foo", "attr1", "attr2"), 
	xml.NewTagFilter("name", "attr1", "attr2"),
	)
schema, err := xml.New(template, bufferSize, filter)
// handle error
```

Then you need to create a `Xml`:
```go
xml := schema.Xml()
templateWithBuffer := dom.Template()
```

Now you can get/set Attribute, get/set Value by using selectors. 
Possible selectors:
```go
elementSelector := xml.ElementSelector("foo")
attributeSelector := xml.AttributeSelector{
	Name: "attr1",
	Value: "value1",
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

    filters := xml.NewFilters(
        xml.NewFilter("foo", "test"),
        xml.NewFilter("id"),
        xml.NewFilter("name"),
        xml.NewFilter("address"),
    )

    schema, err := xml.New(template, filters)
    if err != nil {
        fmt.Println(err)
        return
    }

    aXml := schema.Xml()

    elemIt := aXml.Select(xml.ElementSelector("foo"), xml.ElementSelector("id"))
    for elemIt.Has() {
        elem, _ := elemIt.Next()
        elem.SetValue("10")
    }

    elemIt = aXml.Select(xml.ElementSelector("foo"), xml.ElementSelector("address"))
    for elemIt.Has() {
        elem, _ := elemIt.Next()
        elem.SetValue("")
        elem.AddElement("<new-elem>New element value</new-elem>")
    }

    elemIt = aXml.Select(xml.AttributeSelector{Name: "test", Value: "true"})
    for elemIt.Has() {
        elem, _ := elemIt.Next()
        elem.AddAttribute("attr1", "value1")
        attribute, ok := elem.Attribute("test")
        if !ok {
            continue
        }
        attribute.Set(strings.ToUpper(attribute.Value()))
    }

    result := aXml.Render()
    fmt.Println(result)

    // Output:
    // <?xml version="1.0" encoding="UTF-8"?>
    // <foo test="TRUE" attr1="value1">
    //     <id>10</id>
    //     <name>foo name</name>
    //     <address>
    //         <new-elem>New element value</new-elem>
    //     </address>
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
filters := option.NewFilters()
filters, err := xml.FiltersOf("foo/price[currency]", "address[country and city]/street")
// handle error
```
