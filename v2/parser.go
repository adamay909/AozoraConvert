package aozoraconvert

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func parse(text string) (*Node, error) {

	tok, err := tokenizeAndFix(text)

	if err != nil {
		return new(Node), err

	}

	return getAST(tok)

}

func getAST(t *token) (nd *Node, err error) {

	if t == nil {

		return new(Node), errors.New("no tokens to process")

	}

	document := newNode("document")

	document.level = 0

	if t.mainTextStart() != t.firstToken() {

		metadataNode := getMetadata(t)

		document.addChild(metadataNode)
	}

	prevNode := new(Node)

	prevNode = document

	nextIsChild := true

	closeNode := false

	secCounter := 0

	figCounter := 0

	imgCounter := 0

	n := new(Node)

	addNewNodeAsChild := nextIsChild

	for e := t.mainTextStart(); e != nil; e = e.next {

		//	fmt.Print(e)

		n = newNode("")

		addNewNodeAsChild = nextIsChild

		switch {

		case e.tokType == endOfLineToken:
			continue

		case e.tokType == textToken:

			n.setType("text")

			nextIsChild = false

			closeNode = false

		case e.tokType == paragraphToken:

			n.setType("paragraph")

			nextIsChild = true

			closeNode = false

		case e.tokType == paragraphEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == rubyGroupStartToken:

			n.setType("ruby group")

			nextIsChild = true

			closeNode = false

		case e.tokType == rubyGroupEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == rubyParentStartToken:

			n.setType("ruby parent")

			nextIsChild = true

			closeNode = false

		case e.tokType == rubyParentEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == rubyStartToken:

			n.setType("ruby")

			nextIsChild = true

			closeNode = false

		case e.tokType == rubyEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == emptyLineToken:

			n.setType("empty line")

			nextIsChild = false

			closeNode = false

		case e.isSectionStart():

			n.setType("section")

			n.SetAttr("level", e.tokType.String())

			nextIsChild = true

			closeNode = false

		case e.isSectionEnd():

			nextIsChild = false

			closeNode = true

		case e.isSectionTitleStart():

			n.setType("section title")

			secCounter++

			n.SetAttr("id", "sec"+strconv.Itoa(secCounter))

			nextIsChild = true

			closeNode = false

		case e.tokType == gaijiCharToken:

			n.setType("gaiji char")

			nextIsChild = false

			closeNode = false

		case e.tokType == kunojiToken:

			n.setType("kunoji")

			nextIsChild = false

			closeNode = false

		case e.tokType == accentStartToken:

			n.setType("accent start")

			nextIsChild = false

			closeNode = false

		case e.tokType == accentEndToken:

			n.setType("accent end")

			nextIsChild = false

			closeNode = false

		case e.isPairClose():

			nextIsChild = false

			closeNode = true

		case e.isFormatOfType(indentationMarker):

			n.setType("indentation")

			n.setIndentationAttr(e)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(bottomalignMarker):

			n.setType("bottom align")

			n.setAlignAttr(e)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(narrowparMarker):

			n.setType("narrow paragraph")

			n.setNarrowParagraphAttr(e)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(decoMarker):

			n.setType("emphasis")

			n.setDecoType(e, decoMarker)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(lineDecoMarker):

			n.setType("line decoration")

			n.setDecoType(e, lineDecoMarker)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(inlineSectionMarker):

			n.setType("inline section")

			secCounter++

			n.SetAttr("id", "sec"+strconv.Itoa(secCounter))

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(windowSectionMarker):

			n.setType("window section")

			secCounter++

			n.SetAttr("id", "sec"+strconv.Itoa(secCounter))

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(inlineNoteMarker):

			n.setType("inline note")

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(fontShapeMarker):

			n.setType("font shape")

			n.setFontShapeAttr(e)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(fontsizeMarker):

			n.setType("font size")

			n.setFontSizeAttr(e)

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(offsetMarker):

			n.setType("offset")

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(directionMarker):

			n.setType("text direction")

			nextIsChild = true

			closeNode = false

		case e.tokType == figureStartToken:

			n.setType("figure")

			figCounter++

			n.SetAttr("id", "fig"+strconv.Itoa(figCounter))

			nextIsChild = true

			closeNode = false

		case e.isImage():

			n.setType("image")

			imgCounter++

			n.setImageData(e)

			n.SetAttr("id", "img"+strconv.Itoa(imgCounter))

			nextIsChild = false

			closeNode = false

		case e.isFormatOfType(captionMarker):

			n.setType("caption")

			if strings.HasPrefix(e.innerString(), blockStartStr) {

				n.SetAttr("scope", "block")

			}

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(paginationMarker):

			n.setType("pagination")

			nextIsChild = false

			closeNode = false

		case e.isFormatOfType(rubylikeNoteMarker):

			n.setType("rubylike note")

			n.setRubylikeAttr(e)

			nextIsChild = true

			closeNode = false

		case e.isKunten():

			n.setType("kunten")

			nextIsChild = false

			closeNode = false

		case e.isOkurigana():

			n.setType("okurigana")

			nextIsChild = false

			closeNode = false

		case e.isFormatOfType(centeringMarker):

			n.setType("centering")

			n.setCenteringAttr(e)

			nextIsChild = true

			closeNode = false

		case e.tokType == bibInfoToken:

			n.setType("bibliographical info")

			nextIsChild = true

			closeNode = false

		case e.tokType == bibInfoEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == gaijiNoteToken:

			n.setType("gaiji note")

			nextIsChild = false

			closeNode = false

		case e.isFormatOfType(warichuLineBreakMarker):

			n.setType("warichu line break")

			nextIsChild = false

			closeNode = false

		case e.tokType == specialCharToken:

			n.setType("special char")

			nextIsChild = false

			closeNode = false

		case e.tokType == gaijiImgToken:

			n.setType("image")

			n.setImageData(e)

			n.SetAttr("style", "inline")

			nextIsChild = false

			closeNode = false

		case e.tokType == mainTextStartToken:

			n.setType("main text")

			nextIsChild = true

			closeNode = false

		case e.tokType == noteEndToken:

			n.setType("special char")

			nextIsChild = false

			closeNode = false

		case e.tokType == endMarkerToken:

			n.setType("end marker")

			nextIsChild = false

			closeNode = false

		case e.tokType == emptyToken:

			continue

		case strings.HasPrefix(e.innerString(), blockStartStr):

			n.setType("unknown block type")

			nextIsChild = true

			closeNode = false

		case strings.HasPrefix(e.innerString(), blockEndStr) && strings.HasSuffix(e.innerString(), formatEndStr):

			nextIsChild = false

			closeNode = true

		default:

			n.setType("unknown")

			msglog.Println("Parser: unknown annotation type: " + e.info())
			if e.next != nil && e.next.tokType == endOfLineToken {
				n.SetAttr("force linebreak", "true")
			}

			nextIsChild = false

			closeNode = false

		}

		switch closeNode {

		case true:

			nextIsChild = false

			prevNode = prevNode.Parent()

			err := isValidStructure(prevNode, e)

			if err != nil {
				if !oTolerant {
					panic(err)
				}
				log.Println(err)
			}

			if prevNode.Attr["raw closer"] == "" {

				prevNode.SetAttr("raw closer", e.innerString())

				prevNode.SetAttr("raw unicode closer", e.innerUnicodeString())

				prevNode.SetAttr("raw jis0213 closer", e.innerJis0213String())

			}

			if prevNode.Attr["maybe kanbun"] != "" {

				prevNode.fixKanbun()

			}

			prevNode.closed = true

		default:

			n.setRaw(e.innerString())

			n.SetAttr("unicode raw", e.innerUnicodeString())

			n.SetAttr("jis0213 raw", e.innerJis0213String())

			n.setBlock(e)

			n.tok = e

			if !nextIsChild {

				n.closed = true

			}

			if e.isSectionStart() {

				prevNode = parentSection(prevNode, e)

				addNewNodeAsChild = true

			}

			if addNewNodeAsChild {

				prevNode.addChild(n)

			} else {

				prevNode.addSibling(n)

			}

			switch {

			case n.Attr["type"] == "image":

				n.fixImage()

			case n.Attr["type"] == "kunten" || n.Attr["type"] == "okurigana":

				if n.Parent().Attr["type"] == "paragraph" {

					n.Parent().SetAttr("maybe kanbun", "true")
				}

			case n.Attr["type"] == "bottom align":

				if n.prev == nil {

					n.SetAttr("scope", "block")

					if !strings.HasPrefix(n.Attr["raw"], blockStartStr) {

						n.SetAttr("raw", blockStartStr+n.Attr["raw"])

						n.SetAttr("unicode raw", n.Attr["raw"])

						n.SetAttr("jis0213 raw", n.Attr["raw"])

					}

					if n.Attr["bottom margin"] == "0" {

						n.SetAttr("raw closer", blockEndStr+"地付き"+formatEndStr)

					} else {

						n.SetAttr("raw closer", blockEndStr+"字上げ"+formatEndStr)

					}

					n.SetAttr("raw unicode closer", n.Attr["raw closer"])

					n.SetAttr("raw jis0213 closer", n.Attr["raw jis0213 closer"])
				}

			case n.Attr["type"] == "special char":

				n.SetAttr("escaped raw", e.specialCharEscape())
			}
			prevNode = n

		}

	}

	if prevNode.Parent() == nil || prevNode.Parent().Attr["type"] != "document" {

		if prevNode.Parent() != nil {
		}

		for _, e := range linearizeNode(document) {

			if !e.closed {
				msglog.Print("unclosed node: ")
				msglog.Print(e.Attr["type"])
				if e.tok != nil {
					msglog.Print(e.tok.info())
				}
				msglog.Print("\n")

			}
		}
	}

	return document, err

}

func getMetadata(t *token) (metadataNode *Node) {

	text := ""

	for e := t.firstToken(); e.tokType != emptyLineToken; e = e.next {

		switch e.tokType {

		case paragraphEndToken:
			text = text + "\n"

		case specialCharToken:
			text = text + e.specialCharEscape()

		default:
			text = text + e.String()
		}
	}

	metadataNode = newNode("metadata")

	lines := strings.Split(text, "\n")

	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	titleNode := newNode("meta title")

	titleText, err := parse(lines[0])

	if err != nil {

		msglog.Fatal(err)

		panic(err.Error())

	}

	titleNode.addChild(titleText.firstChild.firstChild)

	metadataNode.addChild(titleNode)

	titleNode.closed = true

	contributorIdx := 1

	minLen := 2

	for k := len(lines) - 1; k > 1; k-- {

		if strings.HasSuffix(lines[k], "訳") {

			minLen++

		}

	}
	if len(lines) > minLen {

		subtitleNode := newNode("meta subtitle")

		subtitleText, err := parse(lines[1])

		if err != nil {

			msglog.Fatal(err)
			panic(err.Error())

		}

		dw := new(strings.Builder)

		renderAozoraText(subtitleText, dw)

		subtitleNode.addChild(subtitleText.firstChild.firstChild)

		metadataNode.addChild(subtitleNode)

		subtitleNode.closed = true

		contributorIdx = 2

	}

	for ; contributorIdx < len(lines); contributorIdx++ {

		contributorNode := newNode("meta contributor")

		contributorName, err := parse(lines[contributorIdx])

		if err != nil {

			panic(err.Error())

		}

		contributorNode.addChild(contributorName.firstChild.firstChild)

		metadataNode.addChild(contributorNode)

		contributorNode.closed = true

	}

	metadataNode.closed = true

	return metadataNode

}

func (n *Node) fixImage() {

	if n.Parent() == nil {
		return
	}

	if n.Parent().Attr["type"] == "figure" {
		return
	}

	if ok, _ := n.withinScopeOfType("paragraph"); ok {

		n.SetAttr("style", "inline")

	}

}

func (n *Node) fixKanbun() {

	defer delete(n.Attr, "maybe kanbun")

	for _, e := range linearizeDescendants(n) {
		switch e.Attr["type"] {
		case "text", "gaiji char", "kunten", "okurigana":
			continue

		default:
			return
		}
	}

	wt := new(strings.Builder)

	renderInnerTextOnly(n, wt)

	for _, c := range wt.String() {

		switch {
		case CharType(c) == Kanji, c == '、', c == '。':

		default:
			return
		}
	}

	msglog.Println("WARNING: line ", strconv.Itoa(n.tok.lineNumber()), " kanbun detected.")

	n.SetAttr("type", "kanbun")
}

// ClearMetadata clears the metadata node of n.
func (n *Node) ClearMetadata() {

	for _, e := range n.Children() {

		if e.Attr["type"] != "metadata" {
			continue
		}

		e.ClearChildren()

		return
	}
}

// SetDodID sets an id for document.
// n must have a child node with Attr["type"]=="metadata"
func (n *Node) SetDocID(id string) {

	if n.Attr["type"] != "document" {
		return
	}

	for _, e := range n.Children() {

		if e.Attr["type"] != "metadata" {
			continue
		}

		e.SetAttr("data-docid", id)
		return
	}
}

// SetTitle sets the title of n to title.
// n must have a child node with Attr["type"]=="metadata"
func (n *Node) SetTitle(title string) {

	if n.Attr["type"] != "document" {
		return
	}

	for _, e := range n.Children() {

		if e.Attr["type"] != "metadata" {
			continue
		}

		t1 := newNode("meta title")

		e.addChild(t1)

		t2 := newNode("text")

		t2.SetAttr("raw", title)

		t1.addChild(t2)

		return
	}
}

// SetSubtitle sets the subtitle.
func (n *Node) SetSubtitle(subtitle string) {

	if n.Attr["type"] != "document" {
		return
	}

	for _, e := range n.Children() {

		if e.Attr["type"] != "metadata" {
			continue
		}

		t1 := newNode("meta subtitle")

		e.addChild(t1)

		t2 := newNode("text")

		t2.SetAttr("raw", subtitle)

		t1.addChild(t2)

		return
	}
}

// AddContributor adds a contributor
func (n *Node) AddContributor(name string) {

	if n.Attr["type"] != "document" {
		return
	}

	for _, e := range n.Children() {

		if e.Attr["type"] != "metadata" {
			continue
		}

		t1 := newNode("meta contributor")

		e.addChild(t1)

		t2 := newNode("text")

		t2.SetAttr("raw", name)

		t1.addChild(t2)

		return
	}
}
