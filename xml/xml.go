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
		if len(elem.attributes) == 0 {
			x.buffer.appendBytes(x.vXml.template[prevStart:elem.start])
		} else {
			x.buffer.appendBytes(x.vXml.template[prevStart:elem.attributes[0][1].start])
			if len(elem.attributes) > 1 {
				x.buffer.appendBytes(x.vXml.template[elem.attributes[0][1].start:elem.attributes[1][0].start])
			}

			for i := 1; i < len(elem.attributes)-1; i++ {
				x.buffer.appendBytes(x.vXml.template[elem.attributes[i][0].start:elem.attributes[i+1][0].start])
			}

			x.buffer.appendBytes(x.vXml.template[elem.attributes[len(elem.attributes)-1][0].start:elem.start])
		}

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
