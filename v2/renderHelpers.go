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

func (n *Node) isBlockFormat() bool {

	return n.Attr["scope"] == "block"

}

func (n *Node) isJisage() bool {

	if n == nil {
		return false
	}

	return n.Attr["type"] == "indentation"

}

func (n *Node) isJiage() bool {

	if n == nil {
		return false
	}

	if n.isBlockFormat() {
		return false
	}

	return n.Attr["type"] == "bottom align"
}

func emphString(s string) string {

	for i, e := range decoMarker {
		if s == e {
			return decoString[i]
		}
	}
	return decoString[0]
}

func (n *Node) decoOnLeft() bool {

	position, ok := n.Attr["position"]

	if !ok {
		return false
	}

	return position == "left"
}

func (n *Node) isChuki() bool {

	for _, m := range rubylikeNoteSimpleMarker {

		if n.Attr["style"] == m {
			return true
		}

	}

	return false
}

func (n *Node) isDeco() bool {

	for _, m := range decoMarker {

		if n.Attr["style"] == m {
			return true
		}

	}

	return false
}

func (n *Node) formatOfType(markerType []string) bool {

	for _, m := range markerType {

		if n.Attr["style"] == m {
			return true
		}

	}

	return false
}

func (n *Node) okuriganaString() string {

	return strings.TrimSuffix(strings.TrimPrefix(n.Attr["raw"], "（"), "）")

}

func (n *Node) headerLevel() int {

	level := 0

	for e := n; e.Parent() != nil; e = e.Parent() {

		if e.Attr["type"] == "section" {
			level++
		}
	}

	return level
}

func (n *Node) withinNoteExclScope() bool {

	e := n

	for e = n; e.Parent().isOfNodeType(gaijiNoteExclusionMarker); e = e.Parent() {
	}

	if e == n {
		return false
	}

	e.hasGaijiWithin = true

	return true

}

func (n *Node) hasChildGaijiNotes() bool {

	nodes := linearizeIsolate(n)

	for _, e := range nodes {

		if e.Attr["type"] == "gaiji note" {

			e.hasGaijiWithin = true

			return true
		}

	}

	return false

}

func (n *Node) isOfNodeType(t []string) bool {

	for _, c := range t {

		if n.Attr["type"] == c {
			return true
		}
	}

	return false
}

func (n *Node) descendantsOfType(t string) []*Node {

	l := linearizeIsolate(n)

	out := []*Node{}

	for _, e := range l {

		if e.Attr["type"] == t {

			out = append(out, e)

		}
	}

	return out

}

func (n *Node) firstDescendantOfType(t string) *Node {

	return n.descendantsOfType(t)[0]

}

func (n *Node) innerText() string {

	return renderSimpleTxt(n.firstChild)

}

func (n *Node) innerTextLength() int {

	return len([]rune(renderSimpleTxt(n)))

}

func (n *Node) innerParagraphCount() int {

	count := 0

	if n.firstChild == nil {
		return count
	}

	for _, e := range linearize(n.firstChild) {

		if e.Attr["type"] == "paragraph" {

			count++
		}

	}

	return count
}
