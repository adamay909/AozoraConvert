package aozoratext

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ast(text string, offset int) *node {

	return getAST(tokenizeAll(text, offset))

}

func getAozoraAST(text string) *node {

	offset := aztextMainStart(text)

	fmt.Println("offset", offset)

	document := ast(text, offset)

	if offset == 0 {

		return document

	}

	metadata := getMetadata(text)

	if metadata != nil {

		document.addFirstChild(metadata)
	}

	return document

}

func getAST(t *token) *node {

	document := newNode("document")

	document.level = 0

	prevNode := new(node)

	prevNode = document

	nextIsChild := true

	closeNode := false

	secCounter := 0

	figCounter := 0

	n := new(node)

	addNewNodeAsChild := nextIsChild

	for e := t.firstToken(); e != nil; e = e.next {

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

		case e.tokType == lineBreakToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == rubyParentStartToken:

			n.setType("ruby parent")

			nextIsChild = true

			closeNode = false

		case e.tokType == rubyParentEndToken:

			nextIsChild = false

			closeNode = true

		case e.tokType == dummyCloserToken:

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

			n.setAttr("level", e.tokType.String())

			nextIsChild = true

			closeNode = false

		case e.isSectionEnd():

			nextIsChild = false

			closeNode = true

		case e.isSectionTitleStart():

			n.setType("section title")

			secCounter++

			n.setAttr("id", "sec"+strconv.Itoa(secCounter))

			nextIsChild = true

			closeNode = false

		case e.tokType == gaijiCharToken:

			n.setType("gaiji char")

			n.setAttr("alt raw", e.altContent)

			nextIsChild = false

			closeNode = false

		case e.tokType == kunojiToken:

			n.setType("kunoji")

			n.setAttr("alt raw", e.altContent)

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

			n.setAttr("id", "sec"+strconv.Itoa(secCounter))

			nextIsChild = true

			closeNode = false

		case e.isFormatOfType(windowSectionMarker):

			n.setType("window section")

			secCounter++

			n.setAttr("id", "sec"+strconv.Itoa(secCounter))

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

			n.setAttr("id", "fig"+strconv.Itoa(figCounter))

			nextIsChild = true

			closeNode = false

		case e.tokType == figureEndToken:

			nextIsChild = false

			closeNode = true

		case e.isImage():

			n.setType("image")

			n.setImageData(e)

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

		case e.tokType == bibInfoToken:

			n.setType("bibliographical info")

			nextIsChild = true

			closeNode = false

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

			n.setAttr("alt raw", e.altContent)

			nextIsChild = false

			closeNode = false

		case e.tokType == accentToken:

			n.setType("accent string")

			n.setAttr("alt raw", e.altContent)

			nextIsChild = false

			closeNode = false

		case e.tokType == emptyToken:

			continue

		default:

			n.setType("unknown")

			log.Println("Warning: unknown annotation type: " + e.String())

			nextIsChild = false

			closeNode = false

		}

		switch closeNode {

		case true:

			nextIsChild = false

			if prevNode.parent() == nil {

				log.Fatal("Structure is invalid. Immediate place of  error is around line ", strconv.Itoa(e.lineNumber())+" "+prevNode.String(), n.String())

			}

			prevNode = prevNode.parent()

			prevNode.setAttr("raw closer", e.innerString())

		default:

			n.setRaw(e.innerString())

			n.setBlock(e)

			n.lineNo = e.lineNo

			if e.isSectionStart() {

				prevNode = parentSection(prevNode, e)

				addNewNodeAsChild = true

			}

			if addNewNodeAsChild {

				prevNode.addChild(n)

			} else {

				prevNode.addSibling(n)

			}

			prevNode = n

		}

	}

	return document

}

func getMetadata(text string) (metadataNode *node) {

	metadataNode = newNode("metadata")

	lines := strings.Split(strings.Split(text, "\n\n")[0], "\n")

	titleNode := newNode("meta title")

	titleNode.addChild(ast(lines[0], 0).firstChild.firstChild)

	metadataNode.addChild(titleNode)

	contributorIdx := 1

	minLen := 2

	for k := len(lines) - 1; k > 1; k-- {

		if strings.HasSuffix(lines[k], "訳") {

			minLen++

		}

	}

	if len(lines) > minLen {

		subtitleNode := newNode("meta subtitle")

		subtitleNode.addChild(ast(lines[1], 0).firstChild.firstChild)

		metadataNode.addChild(subtitleNode)

		contributorIdx = 2

	}

	for ; contributorIdx < len(lines); contributorIdx++ {

		contributorNode := newNode("meta contributor")

		contributorNode.addChild(ast(lines[contributorIdx], 0).firstChild.firstChild)

		metadataNode.addChild(contributorNode)

	}

	return metadataNode

}
