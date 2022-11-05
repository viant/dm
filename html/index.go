package html

import "strings"

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
	langAttribute = iota
	srcAttribute
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
	attributes map[string]int32
	tags       map[string]int32
}

func newIndex() *index {
	return &index{
		attributes: map[string]int32{},
		tags:       map[string]int32{},
	}
}

func (i *index) tagIndex(tag string, createIfAbsent bool) int32 {
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

		tagIndex = int32(len(i.tags) + lastTag)
		i.tags[tag] = tagIndex
		return tagIndex
	}
}

func (i *index) attributeIndex(attribute string, createIfAbsent bool) int32 {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case `lang`:
		return langAttribute
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
		attributeIndex, ok := i.attributes[attribute]
		if ok {
			return attributeIndex
		}

		if !createIfAbsent {
			return -1
		}

		attributeIndex = int32(len(i.tags)) + lastAttribute
		i.attributes[attribute] = attributeIndex
		return attributeIndex
	}
}
