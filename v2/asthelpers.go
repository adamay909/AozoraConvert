package aozoratext

import (
	"log"
	"strconv"
	"strings"
)

func parentSection(pn *node, t *token) *node {

	if pn.parent() == nil {
		return pn
	}

	target := t.tokType.String()

	for e := pn; e.parent() != nil; e = e.parent() {

		if e.attr["type"] == "section" {
			if e.attr["level"] == target {
				return e.parent()
			}
		}
	}

	return pn.parent()
}

func (n *node) setDecoType(e *token, m []string) {

	ft := strings.TrimPrefix(strings.TrimPrefix(e.innerString(), "ここから"), "左に")
	for _, e := range m {
		if ft == e {
			n.setAttr("style", e)
			if e == "罫囲み" {
				return
			}
			break
		}
	}

	if e.decorationLeft() {

		n.setAttr("position", "left")
	} else {
		n.setAttr("position", "right")
	}

	return

}

func (n *node) setFontShapeAttr(e *token) {

	n.setAttr("font shape", strings.TrimPrefix(e.innerString(), "ここから"))

}

func (n *node) setFontSizeAttr(e *token) {

	str := strings.TrimPrefix(e.innerString(), "ここから")

	ft := ""

	for _, ft = range fontsizeMarker {

		if strings.HasSuffix(str, ft) {

			break
		}
	}

	n.setAttr("size", ft)

	n.setAttr("step", getNumberString(strings.TrimSuffix(str, "段階"+ft)))

}

func (n *node) setIndentationAttr(e *token) {

	str := strings.TrimPrefix(e.innerString(), "ここから")

	switch {

	case strings.HasPrefix(str, "改行天付き"):

		str = strings.TrimPrefix(str, "改行天付き、折り返して")

		n.setAttr("top margin", getNumberString(str))

		n.setAttr("indent", "-"+getNumberString(str))

	case strings.Contains(str, "折り返して"):

		part := strings.Split(str, "字下げ、折り返して")

		a, _ := strconv.Atoi(getNumberString(part[0]))

		b, _ := strconv.Atoi(getNumberString(part[1]))

		n.setAttr("top margin", strconv.Itoa(b))

		n.setAttr("indent", strconv.Itoa(a-b))

	case strings.HasSuffix(str, "字下げ"):

		str = strings.TrimSuffix(str, "字下げ")

		n.setAttr("top margin", getNumberString(str))

		n.setAttr("indent", "0")

	}

	return
}

func (n *node) setNarrowParagraphAttr(e *token) {

	str := strings.TrimSuffix(strings.TrimPrefix(e.innerString(), "ここから"), "字詰め")

	n.setAttr("width", getNumberString(str))

	return
}

func (n *node) setAlignAttr(e *token) {

	if strings.HasPrefix(e.innerString(), "ここから") {
		n.setAttr("scope", "block")
	}

	str := strings.TrimPrefix(e.innerString(), "ここから")

	switch {

	case strings.HasSuffix(str, "地付き"):

		n.setAttr("bottom margin", "0")

	case strings.HasSuffix(str, "字上げ"):

		n.setAttr("bottom margin", getNumberString(strings.TrimPrefix(strings.TrimSuffix(str, "字上げ"), "地から")))

	}

	return
}

// Convert full width numerals to half width numeral.
// it will work through str until it encounters a
// character that is not a fw-numeral
func getNumberString(str string) (num string) {

	rs := []rune(str)

	for _, c := range rs {

		if !isFWnumeral(c) {
			break
		}

		num = num + strconv.Itoa(intOf[c])

	}

	return num
}

func (n *node) setImageData(e *token) {

	p := strings.Split(e.innerString(), "（")

	n.setAttr("alt text", p[0])

	p = strings.Split(strings.TrimSuffix(p[1], "）入る"), "、")

	n.setAttr("file", p[0])

	if len(p) == 1 {

		n.setImageSizeFromFile()

		return

	}

	p = strings.Split(p[1], "×")

	n.setAttr("width", strings.TrimPrefix(p[0], "横"))

	n.setAttr("height", strings.TrimPrefix(p[1], "縦"))

	return

}

func (n *node) setImageSizeFromFile() {

	n.setAttr("width", "Width")

	n.setAttr("height", "Height")

	log.Println("setting image from file not implemented")

	return

}

func (n *node) setPaginationStyle(e *token) {

	n.setAttr("pagination style", e.innerString())

}

func (n *node) setRubylikeAttr(e *token) {

	if e.decorationLeft() {

		n.setAttr("position", "left")
	} else {
		n.setAttr("position", "right")
	}
}

func (n *node) setCenteringAttr(e *token) {

	rcounter := 0

	f := new(token)

	for f = e.next; f.tokType == emptyLineToken; f = f.next {

		rcounter++

		if f == nil {
			return
		}

	}

	for ; !f.isPagination(); f = f.next {

		if f == nil {
			return
		}
	}

	lcounter := 0

	for f := f.prev; f.tokType == emptyLineToken; f = f.prev {

		lcounter++

		if f == nil {
			return
		}
	}

	switch {

	case rcounter == lcounter:
		n.setAttr("offset", "none")

	case rcounter > lcounter:
		n.setAttr("offset", "left")

	case rcounter < lcounter:
		n.setAttr("offset", "right")

	}

	n.setAttr("scope", "block")

	return
}

func (n *node) sectionLevel() string {

	level := 0

	for e := n; e != n.topNode(); e = e.parent() {
		level++
	}

	return strconv.Itoa(level)

}

func (n *node) setBlock(e *token) {

	if e.tokType != noteToken {
		return
	}

	if strings.HasPrefix(e.innerString(), "ここから") {
		n.setAttr("scope", "block")
	}

	if e.isFormatOfType(centeringMarker) {
		n.setAttr("scope", "block")
	}

	return
}

func _matchedPair(n *node, t *token) bool {

	return true

	return strings.TrimPrefix(n.attr["raw"], blockStartStr) == strings.TrimPrefix(strings.TrimSuffix(t.innerString(), formatEndStr), blockEndStr)

}
