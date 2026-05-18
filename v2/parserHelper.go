package aozoraconvert

import (
	"errors"
	"strconv"
	"strings"
)

func isValidStructure(prevNode *Node, tok *token) (err error) {

	msgBuilder := new(strings.Builder)

	switch {

	case prevNode == nil:

		addToStringsBuilder(msgBuilder, "\nline ", strconv.Itoa(tok.lineNumber()), ": ", "ERROR: attempting to close node without opener: ", tok.info())

		return errors.New(msgBuilder.String())

	case prevNode.tok == nil:

		addToStringsBuilder(msgBuilder, "\nline ", strconv.Itoa(tok.lineNumber()), ": ", "ERROR: attempting to close node without opener: ", tok.info())

		return errors.New(msgBuilder.String())

	case !matched(tok, prevNode.tok):

		addToStringsBuilder(msgBuilder, "\nERROR: mismatched annotation start and end:\n")

		addToStringsBuilder(msgBuilder, prevNode.tok.info(), prevNode.tok.prev.String(), "\n closed by: \n")

		if prevNode.tok.tokType == rubyParentStartToken {
			addToStringsBuilder(msgBuilder, printContext(prevNode.tok, 10))
		}

		addToStringsBuilder(msgBuilder, tok.info())

		return errors.New(msgBuilder.String())

	default:

		return err

	}
}

func matched(closer, opener *token) bool {

	if closer == nil || opener == nil {
		return false
	}

	switch closer.tokType {

	case paragraphEndToken:
		return opener.tokType == paragraphToken

	case rubyGroupEndToken:
		return opener.tokType == rubyGroupStartToken

	case rubyEndToken:
		return opener.tokType == rubyStartToken

	case rubyParentEndToken:
		return opener.tokType == rubyParentStartToken

	case sectionEndToken:
		return opener.tokType == sectionToken

	case subsectionEndToken:
		return opener.tokType == subsectionToken

	case subsubsectionEndToken:
		return opener.tokType == subsubsectionToken

	case figureEndToken:
		return opener.tokType == figureStartToken

	case centeringEndToken:
		return opener.isFormatOfType(centeringMarker)

	case bibInfoEndToken:
		return opener.tokType == bibInfoToken

	case accentEndToken:
		return opener.tokType == accentStartToken

	case mainTextEndToken:
		return opener.tokType == mainTextStartToken
	}

	if opener.isBlockStartNote() {

		if !closer.isBlockClosingNote() {

			return false

		}
	}

	if closer.isBlockClosingNote() {

		if !opener.isBlockStartNote() {

			return false

		}
	}

	switch {

	case opener.isFormatOfType(bottomalignMarker):
		if !opener.isBlockStartNote() {
			return closer.tokType == alignBottomCloserToken
		}
		fallthrough

	default:
		for _, m := range pairMarker {
			if strings.HasSuffix(opener.innerString(), m) {
				return strings.HasSuffix(strings.TrimPrefix(strings.TrimSuffix(closer.innerString(), formatEndStr), blockEndStr), m)
			}
		}

	}

	if strings.TrimPrefix(opener.innerString(), blockStartStr) == strings.TrimPrefix(strings.TrimSuffix(closer.innerString(), formatEndStr), blockEndStr) {
		return true
	}

	return false
}

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

	defer func() {

		if r := recover(); r != nil {
			msg := `can't parse info from ` + e.info()
			panic(errors.New(msg))
		}
	}()

	n.SetAttr("top margin", strconv.Itoa(e.getTopMargin()))

	ind := e.getIndentation()

	switch {

	case ind != 0:

		n.SetAttr("indent", strconv.Itoa(ind))

	default:

		n.SetAttr("indent", "0")
	}
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

// Convert full width numerals to half width numeral.
// it will work through str until it encounters a
// fw-numeral string
func getFirstNumberString(str string) (num string) {

	rs := []rune(str)

	found := false

	for _, c := range rs {

		if !isFWnumeral(c) {
			if found {
				break
			}
			continue
		}

		found = true

		num = num + strconv.Itoa(intOf[c])

	}

	return num
}
func (n *Node) setImageData(e *token) {

	p := strings.Split(e.innerString(), "（fig")

	n.SetAttr("alt text", p[0])

	p = strings.Split(strings.TrimSuffix(p[1], "）入る"), "、")

	n.SetAttr("file", "fig"+p[0])

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

	msglog.Println("setting image from file not implemented")

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
