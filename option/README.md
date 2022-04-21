## Usage
Package options provides generic options shared between `xml` and `html` package. Currently supported options are:
* `BufferSize` - initial buffer size. Buffers allocates `n` bytes upfront, and use offset to returns the currently used space
* `Filters` - can be specified in order to not index all attributes, elements or tags. 


For options specific only  for the `xml` package see [xml options](../xml/README.md#options)

For options specific only for the `html` package see [html options](../html/README.md#options)