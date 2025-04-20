package aozoraconvert

import (
	"strings"
)

var gaijiNoteExclusionMarker = []string{
	"ruby group",
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

	return strings.TrimSuffix(strings.TrimPrefix(n.rawString(), "（"), "）")

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

	for e = n.Parent(); e != nil; e = e.Parent() {

		if e.isOfNodeType(gaijiNoteExclusionMarker) {
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

	out := []*Node{}

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == t {

			out = append(out, e)

		}
	}

	return out

}

func (n *Node) hasDescendantOfType(t string) bool {

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

func (n *Node) rawString() string {

	switch {

	case oJis0208:
		return regularizeLaTeX(n.Attr["raw"])

	case oJis0213:
		if n.Attr["jis0213 raw"] != "" {
			return regularizeLaTeX(n.Attr["jis0213 raw"])
		}

	case oFull:
		if n.Attr["unicode raw"] != "" {
			return regularizeLaTeX(n.Attr["unicode raw"])
		}
	}

	return regularizeLaTeX(n.Attr["raw"])

}

func (n *Node) rawCloserString() string {

	switch {

	case oJis0208:
		return n.Attr["raw closer"]

	case oJis0213:
		if n.Attr["raw jis0213 closer"] != "" {
			return n.Attr["raw jis0213 closer"]
		}

	case oFull:
		if n.Attr["raw unicode closer"] != "" {
			return n.Attr["raw unicode closer"]
		}
	}

	return n.Attr["raw closer"]

}

func (n *Node) insideSingleLineCommand() bool {

	for e := n.Parent(); e != nil; e = e.Parent() {

		switch e.Attr["type"] {

		case "section":
			return false

		case "main text":
			return false

		case "document":
			return false

		case "indentation":
			return false

		case "bottom align":
			return false

		case "narrow paragraph":
			return false

		case "bibliographical info":
			return false

		}

		return true
	}

	return false
}

func (n *Node) withinScopeOfType(s string) (bool, *Node) {

	for e := n.Parent(); e != nil; e = e.Parent() {

		if e.Attr["type"] == s {
			return true, e
		}
	}

	return false, nil
}

func (n *Node) splitRuby() {

	var rp []rune

	var rt []string

	var gc, gc1, gc2 int

	wt := new(strings.Builder)

	renderInnerTextOnly(n.firstChild, wt)

	rp = []rune(wt.String())

	wt.Reset()

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "ruby" {

			renderInnerTextOnly(e, wt)

			wr := []rune(wt.String())
			for i := 0; i < len(wr); i++ {

				if wr[i] == '〳' {
					rt = append(rt, "{"+string(wr[i:i+2])+"}")
					i++
					continue
				}

				if wr[i] == '〴' {
					rt = append(rt, "{"+string(wr[i:i+2])+"}")
					i++
					continue
				}

				rt = append(rt, string(wr[i:i+1]))

			}

			break
		}
	}

	n.Attr["ruby base"] = string(rp)

	n.Attr["ruby string"] = strings.Join(rt, "")

	if len(rp) < 2 || len(rt) < 2 {

		return
	}

	wt.Reset()

	gc1 = len(rp)/10 + 1

	gc2 = len(rt)/10 + 1

	gc = gc1

	if gc1 == 1 && gc2 == 1 {
		return
	}

	if gc1 > gc2 {

		if gc1 > len(rt) {

			for ; gc1 > len(rt); gc1-- {
			}

		}

		gc = gc1

	}

	if gc1 < gc2 {

		if gc2 > len(rp) {

			for ; gc2 > len(rp); gc2-- {
			}

		}

		gc = gc2

	}

	l1 := len(rp) / gc

	c := 0

	var rps []string

	for c = 1; c < gc; c++ {

		rps = append(rps, string(rp[l1*(c-1):l1*c]))

	}

	if l1*(c-1) < len(rp) {
		rps = append(rps, string(rp[l1*(c-1):]))

	}

	c = 0

	l1 = len(rt) / gc

	var rts []string

	for c = 1; c < gc; c++ {

		rts = append(rts, strings.Join(rt[l1*(c-1):l1*c], ""))

	}

	if l1*(c-1) < len(rt) {
		rts = append(rts, strings.Join(rt[l1*(c-1):], ""))
	}

	n.Attr["ruby base"] = strings.Join(rps, "\t")

	n.Attr["ruby string"] = strings.Join(rts, "\t")

	return
}

func regularizeLaTeX(s string) string {

	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, fwqSpaceStr, fwqRegStr), fwexSpaceStr, fwexRegStr), dashStr, dashStrLatex)

}

func maxSectionDepth(n *Node) int {

	maxdepth := 0

	top := new(Node)

	for top = n; top.Attr["type"] != "document"; top = top.Parent() {
	}

	for _, e := range linearizeNode(top) {

		if e.Attr["type"] == "section" {

			if e.sectionLevel() > maxdepth {
				maxdepth = e.sectionLevel()
			}
		}
	}

	return maxdepth
}

func (n *Node) getTitle() string {

	w := new(strings.Builder)

	for _, e := range linearizeNode(n) {

		if e.Attr["type"] == "meta title" {

			renderInnerTextOnly(e, w)

			return w.String()

		}
	}

	return ""
}

func (n *Node) getTitleNode() *Node {

	for _, e := range linearizeNode(n) {

		if e.Attr["type"] == "meta title" {

			return e

		}
	}

	return nil
}

func (n *Node) getTitleString() string {

	w := new(strings.Builder)

	e := n.getTitleNode()

	if e == nil {
		return ""
	}

	renderInnerTextOnly(e, w)

	return w.String()
}

func (n *Node) sectionStructure() *Node {

	top := newNode("top")

	top.SetAttr("title", n.getTitleString())

	top.SetAttr("id", "main")

	prevNode := top

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] != "section" {
			continue
		}

		sec := newNode("section")

		sec.SetAttr("title", e.getSectionTitle())

		sec.SetAttr("id", e.getSectionID())

		if e.sectionLevel() == prevNode.sectionLevel() {

			prevNode.addSibling(sec)

			prevNode = sec

			continue

		}

		if e.sectionLevel() > prevNode.sectionLevel() {

			prevNode.addChild(sec)

			prevNode = sec

			continue

		}

		for f := prevNode; ; f = f.Parent() {

			if f.sectionLevel() < e.sectionLevel() {

				f.addChild(sec)

				prevNode = sec

				break
			}
		}

	}

	return top

}

func (n *Node) getSectionTitle() string {

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "section title" {

			w := new(strings.Builder)

			renderInnerTextOnly(e, w)

			return w.String()

		}
	}

	return ""
}

func (n *Node) getSectionID() string {

	for _, e := range linearizeDescendants(n) {

		if e.Attr["type"] == "section title" {

			return e.Attr["id"]

		}
	}

	return ""
}
