package aozoraconvert

import (
	"log"
	"strings"
)

func renderHTMLMonolithic(n *Node, w *strings.Builder) {

	defer func() {
		r := recover()
		if r != nil {
			log.Println("RENDERING FAILED")
		}
	}()

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

	h := genImageTag(n)

	h.deleteKey("src")

	h.addExtraKeyVal("src", n.Attr["data"])

	h.setSelfClose()

	h.AddStringTo(w)

}
