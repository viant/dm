package xml

//AttributeSelector matches Element by Attribute name and value
type AttributeSelector struct {
	Name  string
	Value string
}

//Selector matches Element by name
type Selector struct {
	Name       string
	Attributes []AttributeSelector
}

//ElementSelector returns new Selector with only Element Name.
func ElementSelector(name string) Selector {
	return Selector{
		Name: name,
	}
}
