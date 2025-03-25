package aozoratext

import (
	"strings"
)

func renderSimpleTxt(n *Node) string {

	return strings.TrimSpace(Serialize(n, azrSimpleFormatterOpen, azrSimpleFormatterClose))

}

// For very simple output. Assumes only text, ruby, gaiji notes,  and forced
// linebreak within warichu.
func azrSimpleFormatterOpen(n *Node) string {

	switch n.Attr["type"] {

	case "text":
		return n.Attr["raw"]

	case "empty line":
		return lineBreakStr

	case "ruby parent":
		return n.Attr["raw"]

	case "ruby":
		return rubyStartStr + n.Attr["raw"]

	case "gaiji note":
		return gaijiNoteOpenTxt(n)

	case "warichu line break":
		return "\n"

	default:
		return ""

	}

}

func azrSimpleFormatterClose(n *Node) string {

	switch n.Attr["type"] {

	case "paragraph":
		return "\n"

	case "ruby":
		return rubyEndStr

	default:
		return ""
	}

}
