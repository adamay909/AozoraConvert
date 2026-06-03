package aozoraconvert

import (
	"mime"
	"path/filepath"
	"strings"
)

func renderHTMLForAzw3(n *Node, w *strings.Builder) {
	oXHTML = true
	Serialize(n, w, htmlAzw3FormatterOpen, htmlAzw3FormatterClose)
	oXHTML = false

}

func htmlAzw3FormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "image":

		imageAzw3(n, w)

	case "pagination":
		paginationHTML(n, w)

	default:
		azrHTMLFormatterOpen(n, w)

	}

}

func htmlAzw3FormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "image":

		return

	case "metadata":
		metadataCloseHTML(n, w)

	default:
		azrHTMLFormatterClose(n, w)

	}

}

func imageAzw3(n *Node, w *strings.Builder) {

	h := newHtag("img")

	h.addExtraKeyVal("width", n.Attr["width"])

	h.addExtraKeyVal("height", n.Attr["height"])

	h.addExtraKeyVal("alt", n.Attr["alt text"])

	h.setSelfClose()

	h.setAfter("\n")

	srcStr := strings.TrimPrefix(n.Attr["id"], "img")

	for len(srcStr) < 4 {

		srcStr = "0" + srcStr

	}

	srcStr = "kindle:embed:" + srcStr + "?mime=" + mime.TypeByExtension(filepath.Ext(n.Attr["file"]))

	h.addExtraKeyVal("src", srcStr)

	h.AddStringTo(w)

}
