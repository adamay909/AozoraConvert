package aozoraConvert

import (
	"errors"
	"fmt"
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

		return new(Node), errors.New("No tokens to process.")

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

		case e.isPairMarkerClose():

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

		case e.tokType == figureEndToken:

			nextIsChild = false

			closeNode = true

		case e.isImage():

			n.setType("image")

			imgCounter++

			n.setImageData(e)

			n.SetAttr("id", "img"+strconv.Itoa(imgCounter))

			nextIsChild = false

			closeNode = false

		case e.isFormatOfType(captionMarker):

			n.setType("caption")

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

		case e.tokType == centeringEndToken:

			nextIsChild = false

			closeNode = true

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

		case e.tokType == alignBottomCloserToken:

			nextIsChild = false

			closeNode = true

		case e.isFormatOfType(warichuLineBreakMarker):

			n.setType("warichu line break")

			nextIsChild = false

			closeNode = false

		case e.tokType == specialCharToken:

			n.setType("special char")

			nextIsChild = false

			closeNode = false

		case e.tokType == accentToken:

			n.setType("accent string")

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

		case e.tokType == mainTextEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == emptyToken:

			continue

		default:

			n.setType("unknown")

			log.Println("Parser: unknown annotation type: " + e.info())

			nextIsChild = false

			closeNode = false

		}

		switch closeNode {

		case true:

			nextIsChild = false

			prevNode = prevNode.Parent()

			err := isValidStructure(prevNode, e)

			if err != nil {
				panic(err.Error())
			}

			prevNode.SetAttr("raw closer", e.innerString())

			prevNode.SetAttr("raw unicode closer", e.unicodeContent)

			prevNode.SetAttr("raw jis0213 closer", e.jis0213Content)

			if prevNode.Attr["maybe kanbun"] != "" {

				prevNode.fixKanbun()

			}

			prevNode.closed = true

		default:

			n.setRaw(e.innerString())

			n.SetAttr("unicode raw", e.unicodeContent)

			n.SetAttr("jis0213 raw", e.jis0213Content)

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

			if n.Attr["type"] == "image" {
				n.fixImage()
			}

			if n.Attr["type"] == "kunten" || n.Attr["type"] == "okurigana" {
				if n.Parent().Attr["type"] == "paragraph" {
					n.Parent().SetAttr("maybe kanbun", "true")
				}
			}

			prevNode = n

		}

	}

	if prevNode.Parent().Attr["type"] != "document" {

		fmt.Println("last node is ", prevNode.Attr["type"])
		fmt.Println("Parent is ", prevNode.Parent().Attr["type"])

		for _, e := range linearizeNode(document) {

			if !e.closed {
				log.Print("unclosed node: ")
				log.Print(e.Attr["type"])
				if e.tok != nil {
					log.Print(e.tok.info())
				}
				log.Print("\n")

			}
		}
	}

	return document, err

}

func getMetadata(t *token) (metadataNode *Node) {

	text := ""

	for e := t.firstToken(); e.tokType != emptyLineToken; e = e.next {

		if e.tokType == paragraphEndToken {
			text = text + "\n"
		} else {
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

		log.Fatal(err)

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

			log.Fatal(err)
			panic(err.Error())

		}
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

	log.Println("WARNING: line ", strconv.Itoa(n.tok.lineNumber()), " kanbun detected.")

	n.SetAttr("type", "kanbun")
}
