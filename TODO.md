## Html:

* Optimizations for the building template lazily. Currently, after changing attribute / innerHTML template is being
  updated instantly. It is possible to instead of using the mutations (when attribute at n-th position changes, all
  attributes from the n to the `len(attributes)` will be increased by the difference `len(current) - len(previous)`). It
  is optimal for small changes to the DOM, but if DOM changes a lot, copying template may be expensive. In that case it
  would be better to render template lazily allocating slice for the attributes new value (slice[i] == nil means that
  value was not updated)
* Add InsertAfter and InsertBefore for HTML Elements

### Xml:

* Add support for `or` for attributes selectors