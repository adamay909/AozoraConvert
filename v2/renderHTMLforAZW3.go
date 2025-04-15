package aozoraConvert

import (
	"mime"
	"path/filepath"
	"strings"
)

func renderHtmlForAzw3(n *Node, w *strings.Builder) {

	Serialize(n, w, htmlAzw3FormatterOpen, htmlAzw3FormatterClose)

}

func htmlAzw3FormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "image":

		imageAzw3(n, w)

	default:
		azrHtmlFormatterOpen(n, w)

	}

}

func htmlAzw3FormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "image":

		return

	default:
		azrHtmlFormatterClose(n, w)

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
