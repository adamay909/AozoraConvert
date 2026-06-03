package aozoraconvert

import (
	"strings"
)

func renderEpubHTML(n *Node, w *strings.Builder) {
	oXHTML = true
	Serialize(n, w, htmlEpubFormatterOpen, htmlEpubFormatterClose)
	oXHTML = false
}

func htmlEpubFormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "empty line":
		h := newHtag("br")
		h.setSelfClose()
		h.AddStringTo(w)
		w.WriteString("\n")

	case "metadata":

		centeringOpenHTML(n, w)
		/*
			case "image":

				imageXHTML(n, w)
		*/
	case "pagination":
		paginationHTML(n, w)

	default:
		azrHTMLFormatterOpen(n, w)

	}

}

func htmlEpubFormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "metadata":

		w.WriteString(newCloseHtag(`div`).String())

	default:
		azrHTMLFormatterClose(n, w)

	}

}

/*
func imageXHTML(n *Node, w *strings.Builder) {

	h := genImageTag(n)

	h.setSelfClose()

	h.AddStringTo(w)

}
*/
/*
func bottomAlignOpenEpubHTML(n *Node, w *strings.Builder) {

	h := newHtag("span")

	if n.isBlockFormat() {
		h.setElement("div")
		h.addClass("alignBottom")
		h.setAfter("\n")

	} else {
		h1 := newHtag("br")
		h1.setSelfClose()
		h1.AddStringTo(w)

		h.addClass("flushBottom")

	}

	h.addClass("bottomMargin" + n.Attr["bottom margin"])

	h.AddStringTo(w)

}
*/
