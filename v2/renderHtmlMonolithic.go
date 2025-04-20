package aozoraconvert

import "strings"

func renderHTMLMonolithic(n *Node, w *strings.Builder) {

	Serialize(n, w, htmlMonoFormatterOpen, htmlMonoFormatterClose)

}

func htmlMonoFormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "image":

		embedImage(n, w)

	default:

		azrHTMLFormatterOpen(n, w)

	}
}

func htmlMonoFormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "image":

		return

	default:
		azrHTMLFormatterClose(n, w)

	}

}

func embedImage(n *Node, w *strings.Builder) {

	h := newHtag("img")

	h.addClass("illustration")

	h.addExtraKeyVal("width", n.Attr["width"])

	h.addExtraKeyVal("height", n.Attr["height"])

	h.addExtraKeyVal("src", n.Attr["data"])

	h.addExtraKeyVal("alt", n.Attr["alt text"])

	h.setSelfClose()

	h.setAfter("\n")

	h.AddStringTo(w)

}
