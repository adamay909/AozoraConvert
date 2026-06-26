package aozoraconvert

import "strings"

func renderTocNcx(n *Node, w *strings.Builder) {
	SerializeDescendants(n, w, ncxOpen, ncxClose)
}

func ncxOpen(n *Node, w *strings.Builder) {
	addToStringsBuilder(w, `<navPoint id="`, n.Attr["id"], `">`, "\n", `<navLabel>`, `<text>`, n.Attr["title"], `</text>`, `</navLabel>`, "\n", `<content src="1.html#`, n.Attr["id"], `"/>`, "\n")
}

func ncxClose(n *Node, w *strings.Builder) {
	addToStringsBuilder(w, `</navPoint>`, "\n")
}
