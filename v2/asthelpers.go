package aozoratext

import (
	"log"
	"strconv"
	"strings"
)

func parentSection(pn *Node, t *token) *Node {

	if pn.Parent() == nil {
		return pn
	}

	if pn.Attr["type"] == "main text" {
		return pn
	}

	target := t.tokType.String()

	for e := pn; e.Parent() != nil; e = e.Parent() {

		if e.Attr["type"] == "section" {
			if e.Attr["level"] == target {
				return e.Parent()
			}
		}
	}

	return pn.Parent()
}

func (n *Node) setDecoType(e *token, m []string) {

	ft := strings.TrimPrefix(strings.TrimPrefix(e.innerString(), "ここから"), "左に")
	for _, e := range m {
		if ft == e {
			n.SetAttr("style", e)
			if e == "罫囲み" {
				return
			}
			break
		}
	}

	if e.decorationLeft() {

		n.SetAttr("position", "left")
	} else {
		n.SetAttr("position", "right")
	}

	return

}

func (n *Node) setFontShapeAttr(e *token) {

	n.SetAttr("font shape", strings.TrimPrefix(e.innerString(), "ここから"))

}

func (n *Node) setFontSizeAttr(e *token) {

	str := strings.TrimPrefix(e.innerString(), "ここから")

	ft := ""

	for _, ft = range fontsizeMarker {

		if strings.HasSuffix(str, ft) {

			break
		}
	}

	n.SetAttr("size", ft)

	n.SetAttr("step", getNumberString(strings.TrimSuffix(str, "段階"+ft)))

}

func (n *Node) setIndentationAttr(e *token) {

	str := strings.TrimPrefix(e.innerString(), "ここから")

	switch {

	case strings.HasPrefix(str, "改行天付き"):

		str = strings.TrimPrefix(str, "改行天付き、折り返して")

		n.SetAttr("top margin", getNumberString(str))

		n.SetAttr("indent", "-"+getNumberString(str))

	case strings.Contains(str, "折り返して"):

		part := strings.Split(str, "字下げ、折り返して")

		a, _ := strconv.Atoi(getNumberString(part[0]))

		b, _ := strconv.Atoi(getNumberString(part[1]))

		n.SetAttr("top margin", strconv.Itoa(b))

		n.SetAttr("indent", strconv.Itoa(a-b))

	case strings.HasSuffix(str, "字下げ"):

		str = strings.TrimSuffix(str, "字下げ")

		n.SetAttr("top margin", getNumberString(str))

		n.SetAttr("indent", "0")

	}

	return
}

func (n *Node) setNarrowParagraphAttr(e *token) {

	str := strings.TrimSuffix(strings.TrimPrefix(e.innerString(), "ここから"), "字詰め")

	n.SetAttr("width", getNumberString(str))

	return
}

func (n *Node) setAlignAttr(e *token) {

	if strings.HasPrefix(e.innerString(), "ここから") {
		n.SetAttr("scope", "block")
	}

	str := strings.TrimPrefix(e.innerString(), "ここから")

	switch {

	case strings.HasSuffix(str, "地付き"):

		n.SetAttr("bottom margin", "0")

	case strings.HasSuffix(str, "字上げ"):

		n.SetAttr("bottom margin", getNumberString(strings.TrimPrefix(strings.TrimSuffix(str, "字上げ"), "地から")))

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

func (n *Node) setImageData(e *token) {

	p := strings.Split(e.innerString(), "（")

	n.SetAttr("alt text", p[0])

	p = strings.Split(strings.TrimSuffix(p[1], "）入る"), "、")

	n.SetAttr("file", p[0])

	if len(p) == 1 {

		n.setImageSizeFromFile()

		return

	}

	p = strings.Split(p[1], "×")

	n.SetAttr("width", strings.TrimPrefix(p[0], "横"))

	n.SetAttr("height", strings.TrimPrefix(p[1], "縦"))

	return

}

func (n *Node) setImageSizeFromFile() {

	n.SetAttr("width", "Width")

	n.SetAttr("height", "Height")

	log.Println("setting image from file not implemented")

	return

}

func (n *Node) setPaginationStyle(e *token) {

	n.SetAttr("pagination style", e.innerString())

}

func (n *Node) setRubylikeAttr(e *token) {

	if e.decorationLeft() {

		n.SetAttr("position", "left")
	} else {
		n.SetAttr("position", "right")
	}
}

func (n *Node) setCenteringAttr(e *token) {

	rcounter := 0

	f := new(token)

	if e.next == nil {
		return
	}

	for f = e.next; f.tokType == emptyLineToken || f.tokType == endOfLineToken; f = f.next {

		if f.tokType == emptyLineToken {
			rcounter++
		}
		if f.next == nil {
			return
		}

	}

	for ; !f.isPagination(); f = f.next {

		if f == nil {
			return
		}
	}

	lcounter := 0

	for f := f.prev; f.tokType == emptyLineToken || f.tokType == endOfLineToken; f = f.prev {

		if f.tokType == emptyLineToken {
			lcounter++
		}

		if f.next == nil {
			return
		}
	}

	switch {

	case rcounter == lcounter:
		n.SetAttr("offset", "none")

	case rcounter > lcounter:
		n.SetAttr("offset", "left")

	case rcounter < lcounter:
		n.SetAttr("offset", "right")

	}

	n.SetAttr("scope", "block")

	return
}

func (n *Node) sectionLevel() string {

	level := 0

	for e := n; e != n.topNode(); e = e.Parent() {
		level++
	}

	return strconv.Itoa(level)

}

func (n *Node) setBlock(e *token) {

	if e.tokType != noteToken {
		return
	}

	if strings.HasPrefix(e.innerString(), "ここから") {
		n.SetAttr("scope", "block")
	}

	if e.isFormatOfType(centeringMarker) {
		n.SetAttr("scope", "block")
	}

	return
}

func _matchedPair(n *Node, t *token) bool {

	return true

	return strings.TrimPrefix(n.Attr["raw"], blockStartStr) == strings.TrimPrefix(strings.TrimSuffix(t.innerString(), formatEndStr), blockEndStr)

}
