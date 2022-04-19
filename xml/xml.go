package xml

type Xml struct {
	vXml   *VirtualXml
	buffer *Buffer
}

func (v *VirtualXml) Xml() *Xml {
	return &Xml{
		vXml:   v,
		buffer: NewBuffer(v.bufferSize),
	}
}

func (x *Xml) Render() string {
	var prevStart int
	var elem *StartElement
	for _, elem = range x.vXml.allElements() {
		x.buffer.appendBytes(x.vXml.template[prevStart:elem.start])
		prevStart = elem.start
	}

	if elem != nil {
		x.buffer.appendBytes(x.vXml.template[elem.start:elem.end])
		x.buffer.appendBytes(x.vXml.template[elem.end:])
	}

	return x.buffer.String()
}

func (x *Xml) Select(selectors ...Selector) *Iterator {
	return NewIterator(x, selectors)
}
