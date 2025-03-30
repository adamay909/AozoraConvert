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

func (n *Node) okuriganaString() string {

	return strings.TrimSuffix(strings.TrimPrefix(n.RawString(), "（"), "）")

}

func (n *Node) sectionLevel() int {

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

func (n *Node) isOfNodeType(t []string) bool {

	for _, c := range t {

		if n.Attr["type"] == c {
			return true
		}
	}

	return false
}

func (n *Node) descendantsOfType(t string) []*Node {

	out := []*Node{}

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == t {

			out = append(out, e)

		}
	}

	return out

}

func (n *Node) HasDescendantOfType(t string) bool {

	return len(n.descendantsOfType(t)) > 0

}

func (n *Node) firstDescendantOfType(t string) *Node {

	return n.descendantsOfType(t)[0]

}

func (n *Node) innerParagraphCount() int {

	count := 0

	if n.firstChild == nil {
		return count
	}

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "paragraph" {

			count++
		}

	}

	return count
}

func (n *Node) RawString() string {

	if o_jis0208 {
		return n.Attr["raw"]
	}

	if n.Attr["unicode raw"] != "" {
		return n.Attr["unicode raw"]
	}

	return n.Attr["raw"]

}

func (n *Node) RawCloserString() string {

	if o_jis0208 {
		return n.Attr["raw closer"]
	}

	if n.Attr["raw unicode closer"] != "" {
		return n.Attr["raw unicode closer"]
	}

	return n.Attr["raw closer"]

}
