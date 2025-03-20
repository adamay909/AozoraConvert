package aozoratext

import "strings"

var gaijiNoteExclusionMarker = []string{
	"ruby parent",
	"emphasis",
	"line decoration",
	"inline section",
	"window section",
	"rubylike note",
}

func (n *node) isBlockFormat() bool {

	return n.attr["scope"] == "block"

}

func (n *node) isJisage() bool {

	return n.attr["type"] == "indentation"

}

func (n *node) isJiage() bool {

	if n.isBlockFormat() {
		return false
	}

	return n.attr["type"] == "bottom align"
}

func emphString(s string) string {

	for i, e := range decoMarker {
		if s == e {
			return decoString[i]
		}
	}
	return decoString[0]
}

func (n *node) decoOnLeft() bool {

	position, ok := n.attr["position"]

	if !ok {
		return false
	}

	return position == "left"
}

func (n *node) isChuki() bool {

	for _, m := range rubylikeNoteSimpleMarker {

		if n.attr["style"] == m {
			return true
		}

	}

	return false
}

func (n *node) isDeco() bool {

	for _, m := range decoMarker {

		if n.attr["style"] == m {
			return true
		}

	}

	return false
}

func (n *node) formatOfType(markerType []string) bool {

	for _, m := range markerType {

		if n.attr["style"] == m {
			return true
		}

	}

	return false
}

func (n *node) okuriganaString() string {

	return strings.TrimSuffix(strings.TrimPrefix(n.attr["raw"], "（"), "）")

}

func (n *node) headerLevel() int {

	level := 0

	for e := n; e.parent() != nil; e = e.parent() {

		if e.attr["type"] == "section" {
			level++
		}
	}

	return level
}

func (n *node) withinNoteExclScope() bool {

	e := n

	for e = n; e.parent().isOfNodeType(gaijiNoteExclusionMarker); e = e.parent() {
	}

	if e == n {
		return false
	}

	e.hasGaijiWithin = true

	return true

}

func (n *node) hasChildGaijiNotes() bool {

	nodes := linearizeIsolate(n)

	for _, e := range nodes {

		if e.attr["type"] == "gaiji note" {

			e.hasGaijiWithin = true

			return true
		}

	}

	return false

}

func (n *node) isOfNodeType(t []string) bool {

	for _, c := range t {

		if n.attr["type"] == c {
			return true
		}
	}

	return false
}

func (n *node) descendantsOfType(t string) []*node {

	l := linearizeIsolate(n)

	out := []*node{}

	for _, e := range l {

		if e.attr["type"] == t {

			out = append(out, e)

		}
	}

	return out

}

func (n *node) firstDescendantOfType(t string) *node {

	return n.descendantsOfType(t)[0]

}

func (n *node) innerText() string {

	return renderSimpleTxt(n.firstChild)

}

func (n *node) innerTextLength() int {

	return len([]rune(renderSimpleTxt(n)))

}

func (n *node) innerParagraphCount() int {

	count := 0

	if n.firstChild == nil {
		return count
	}

	for _, e := range linearize(n.firstChild) {

		if e.attr["type"] == "paragraph" {

			count++
		}

	}

	return count
}
