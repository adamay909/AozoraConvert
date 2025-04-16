package aozoraConvert

import "strings"

func renderEpubHTML(n *Node, w *strings.Builder) {

	Serialize(n, w, htmlEpubFormatterOpen, htmlEpubFormatterClose)

}

func htmlEpubFormatterOpen(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "metadata":

		centeringOpenHtml(n, w)

	default:
		azrHtmlFormatterOpen(n, w)

	}

}

func htmlEpubFormatterClose(n *Node, w *strings.Builder) {

	switch n.Attr["type"] {

	case "metadata":

		w.WriteString(newCloseHtag(`div`).String())

	default:
		azrHtmlFormatterClose(n, w)

	}

}
