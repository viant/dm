package xml

//Selector represents generic selectors, can be AttributeSelector or ElementSelector
type Selector interface{}

//AttributeSelector matches Element by Attribute name and value
type AttributeSelector struct {
	Name  string
	Value string
}

//ElementSelector matches Element by name
type ElementSelector string
