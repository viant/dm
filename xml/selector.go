package xml

type Selector interface{}

type AttributeSelector struct {
	Name  string
	Value string
}

type ElementSelector string
