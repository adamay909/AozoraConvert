package aozoraconvert

import (
	"strings"
)

func renderEpubHTML(n *Node, w *strings.Builder) {

	Serialize(n, w, htmlEpubFormatterOpen, htmlEpubFormatterClose)

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

	case "image":

		imageXHTML(n, w)

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

func imageXHTML(n *Node, w *strings.Builder) {

	h := genImageTag(n)

	h.setSelfClose()

	h.AddStringTo(w)

}
