package dm

const (
	htmlTag = iota
	headTag
	bodyTag
	imgTag
	iframeTag
	titleTag
	h1Tag
	h2Tag
	h3Tag
	pTag
	ulTag
	liTag
	divTag
	aTag
	olTag
	formTag
	tableTag
	theadTag
	tbodyTag
	trTag
	tdTag
	thTag
	lastTag
)

const (
	srcAttribute = iota
	altAttribute
	classAttribute
	idAttribute
	hrefAttribute
	styleAttribute
	widthAttribute
	heightAttribute
	lastAttribute
)

type index struct {
	attributeIndex map[string]int
	tags           map[string]int
}

func newIndex() *index {
	return &index{
		attributeIndex: map[string]int{},
		tags:           map[string]int{},
	}
}

func (i *index) tag(tag string, createIfAbsent bool) int {
	switch tag {
	case `html`:
		return htmlTag
	case `head`:
		return headTag
	case `body`:
		return bodyTag
	case `img`:
		return imgTag
	case `iframe`:
		return iframeTag
	case `title`:
		return titleTag
	case `h1`:
		return h1Tag
	case `h2`:
		return h2Tag
	case `h3`:
		return h3Tag
	case `p`:
		return pTag
	case `ul`:
		return ulTag
	case `li`:
		return liTag
	case `div`:
		return divTag
	case `a`:
		return aTag
	case `ol`:
		return olTag
	case `form`:
		return formTag
	case `table`:
		return tableTag
	case `thead`:
		return theadTag
	case `tbody`:
		return tbodyTag
	case `tr`:
		return trTag
	case `td`:
		return tdTag
	case `th`:
		return thTag
	default:
		tagIndex, ok := i.tags[tag]
		if ok {
			return tagIndex
		}

		if !createIfAbsent {
			return -1
		}

		tagIndex = len(i.tags) + lastTag
		i.tags[tag] = tagIndex
		return tagIndex
	}
}

func (i *index) attribute(attribute string, createIfAbsent bool) int {
	switch attribute {
	case `src`:
		return srcAttribute
	case `alt`:
		return altAttribute
	case `class`:
		return classAttribute
	case `id`:
		return idAttribute
	case `href`:
		return hrefAttribute
	case `style`:
		return styleAttribute
	case `width`:
		return widthAttribute
	case `height`:
		return heightAttribute
	default:
		attributeIndex, ok := i.attributeIndex[attribute]
		if ok {
			return attributeIndex
		}

		if !createIfAbsent {
			return -1
		}

		attributeIndex = len(i.tags) + lastAttribute
		i.attributeIndex[attribute] = attributeIndex
		return attributeIndex
	}
}
