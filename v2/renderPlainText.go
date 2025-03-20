package aozoratext

import (
	"strings"
)

func renderSimpleTxt(n *node) string {

	return strings.TrimSpace(stringify(n, azrSimpleFormatterOpen, azrSimpleFormatterClose))

}

// For very simple output. Assumes only text, ruby, gaiji notes,  and forced
// linebreak within warichu.
func azrSimpleFormatterOpen(n *node) string {

	switch n.attr["type"] {

	case "text":
		return n.attr["raw"]

	case "empty line":
		return lineBreakStr

	case "ruby parent":
		return n.attr["raw"]

	case "ruby":
		return rubyStartStr + n.attr["raw"]

	case "gaiji note":
		return gaijiNoteOpenTxt(n)

	case "warichu line break":
		return "\n"

	default:
		return ""

	}

}

func azrSimpleFormatterClose(n *node) string {

	switch n.attr["type"] {

	case "paragraph":
		return "\n"

	case "ruby":
		return rubyEndStr

	default:
		return ""
	}

}
