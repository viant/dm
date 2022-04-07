# DM (DOM manipulator for golang)

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/dm)](https://goreportcard.com/report/github.com/viant/dm)
[![GoDoc](https://godoc.org/github.com/viant/velty?status.svg)](https://godoc.org/github.com/viant/dm)

This library is compatible with Go 1.17+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Usage](#usage)
- [Contribution](#contributing-to-DM)

## Usage

In order to change to update the DOM, you need to create `DOM` representation. `DOM` representation can be shared across the app: 
```go
template := []byte("<html>...</html>")
dom, err := dm.New(template)
// handle error
```

You can specify some options while creating `DOM`:
* `BufferSize` - initial buffer size for each `DOM` session, `int` wrapper
* `Filter` - represents allowed tags and attributes, `map[string]map[string]bool` wrapper. The outer map, specifies the tags, 
the inner - attributes. The value of inner map is not important, if value is specified in the map, the attribute may be indexed.
```go
bufferSize := dm.BufferSize(1024)
filter := dm.Filter(map[string]map[string]bool{
	"div": {"class": true} 
})
dom, err := dm.New(template, bufferSize, filter)
// handle error
```

Then you need to create the `Session`:
```go
newSession := dom.Session()
sessionWithBuffer := dom.Session(dm.NewBuffer(1024))
```
If you don't provide a `Buffer`, will be created one with `BufferSize` specified while creating `DOM`

Now you can get/set Attribute, get/set InnerHTML either by using selectors, or by index. 

Selectors order: `Tag` -> `Attribute` -> `Attribute value`. Selectors are optional, it means if you don't specify `Attribute` only 
tag will be checked. 

`InnerHTML / Attribute` usage:
```go

var offset int
var innerHTML []byte
for {
    innerHTML, offset = session.InnerHTML(offset, []byte("div"))
    if offset == -1 {
    	break
    }
    //...
    session.SetInnerHTMLByIndex(offset, bytes.Trim(innerHTML))
}
```

`SetInnerHTML / SetAttribute` usage:

```go

var offset int
for {
    offset = session.SetInnerHTML(offset, []byte("This is new inner HTML") ,[]byte("div")), []byte("class"), []byte("div-class"))
    if offset == -1 {
        break
    }
    //... 
}
```

## Contributing to DM

DM is an open source project and contributors are welcome!

See [Todo](TODO.md) list.

