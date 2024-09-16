package azrconvert

import (
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type section struct {
	id, title                                    string
	level                                        int
	parent, prevSibling, nextSibling, firstChild *section
	node                                         *html.Token
	start, end                                   int
}

var counter int

func (b *Book) RenderTOC() string {
	w := new(strings.Builder)

	if b.TopSection.title == "" {
		b.TopSection.title = b.Title
	}

	if b.TopSection.id == "" {
		b.TopSection.id = "azbc_100"
	}

	s := b.TopSection

	addToTOC(s, w)
	return w.String()
}

func addToTOC(s *section, w *strings.Builder) {

	var lead string
	counter++
	//	lead = strings.Repeat("    ", headerLevel(s.node)-3)

	//	if len(s.content) != 0 {
	w.WriteString(lead + `<navPoint id="` + s.id + `" playOrder="` + strconv.Itoa(counter) + `">` + "\n")
	w.WriteString(lead + "\t<navLabel>\n")
	w.WriteString(lead + "\t\t<text>" + s.title + "</text>\n")
	w.WriteString(lead + "\t</navLabel>\n")
	w.WriteString(lead + "\t<content src=" + `"1.html#` + s.id + `" />` + "\n")
	//	}
	if s.firstChild != nil {
		addToTOC(s.firstChild, w)
	}
	w.WriteString(lead + "</navPoint>\n")
	if s.nextSibling != nil {
		addToTOC(s.nextSibling, w)
	}

	return
}

func (b *Book) RenderEP3TOC() string {

	w := new(strings.Builder)

	if b.TopSection.title == "" {
		b.TopSection.title = b.Title
	}

	s := b.TopSection

	w.WriteString("<ol>\n")

	addToEP3TOC(s, w)

	w.WriteString("</ol>\n")

	return w.String()
}

func addToEP3TOC(s *section, w *strings.Builder) {

	var lead string
	counter++
	//lead = strings.Repeat("    ", headerLevel(s.node)-3)

	//	if len(s.content) != 0 {
	w.WriteString(lead + `<li>`)
	w.WriteString(`<a href="1.html#` + s.id + `">` + s.title + "</a>")
	//	}
	if s.firstChild != nil {
		w.WriteString("\n")
		w.WriteString(lead + "<ol>\n")
		addToEP3TOC(s.firstChild, w)
		w.WriteString(lead + "</ol>\n")
	}

	w.WriteString("</li>\n")

	if s.nextSibling != nil {
		addToEP3TOC(s.nextSibling, w)
	}
	return
}

func insertSectionID(tokens []*html.Token) {
	c := 100
	for i, token := range tokens {
		if !isHeader(token) {
			continue
		}
		if hasID(token) {
			continue
		}
		c = c + 10
		tokens[i].Attr = append(tokens[i].Attr, html.Attribute{Namespace: "", Key: "id", Val: "azbc_" + strconv.Itoa(c)})
	}

	if c == 100 {
		tokens[0].Attr = append(tokens[0].Attr, html.Attribute{Namespace: "", Key: "id", Val: "azbc_" + strconv.Itoa(c)})
	}

	return
}

func (b *Book) getStructure() *section {

	tokens := b.Body
	sec := new(section)
	sec.level = 1
	sec.start = 0
	sec.end = len(tokens) - 1
	sec.nextSibling = nil
	sec.firstChild = nil
	sec.parent = nil
	sec.node = new(html.Token)
	sec.node.DataAtom = atom.H3
	first := sec
	i := 0

	for i = 0; i < len(tokens); i++ {
		if isHeader(tokens[i]) {
			break
		}
	}

	if i == len(tokens) {
		return first
	}
	//Guard against the possibility that the first
	//section of document is level 2 rather than
	//level 1.
	if tokens[i].DataAtom != maxLevel(tokens) {
		sec = new(section)
		sec.level = 1
		sec.start = 0
		sec.end = len(tokens) - 1
		sec.nextSibling = nil
		sec.firstChild = nil
		sec.parent = nil
		sec.title = getTextContent(tokens, i)
		sec.node = new(html.Token)
		sec.node.DataAtom = atom.H3
		sec.node.Type = html.StartTagToken
		first.nextSibling = sec
	}

	for i < len(tokens) {
		tok := tokens[i]
		if isHeader(tok) {
			newsec := new(section)
			newsec.node = tok
			newsec.start = i
			newsec.end = len(tokens) - 1
			newsec.id = getID(tok)
			newsec.title = getTextContent(tokens, i)
			switch {
			case isChild(newsec, sec):
				sec.firstChild = newsec
				newsec.parent = sec
				newsec.level = sec.level + 1

			case isParent(newsec, sec):
				sec.parent.nextSibling = newsec
				newsec.parent = sec.parent.parent
				newsec.level = sec.level - 1

			case isGrandParent(newsec, sec):
				sec.parent.parent.nextSibling = newsec
				if sec.parent.parent.parent != nil {
					newsec.parent = sec.parent.parent.parent
					newsec.level = sec.level - 2
				}

			default:
				sec.nextSibling = newsec
				newsec.parent = sec.parent
				newsec.level = sec.level
			}
			sec.end = newsec.start - 1
			sec = newsec
			log.Println("Found section: ", newsec.title)
		}
		i++
	}
	if first.firstChild != nil {
		first.firstChild.start = 0
		return first.firstChild
	}
	if first.nextSibling != nil {
		first.nextSibling.start = 0
		return first.nextSibling
	}

	return first

}

// Find the textual content of t. We remove all rubys.
func getTextContent(t []*html.Token, i int) string {

	node := getNode(t[i:])

	w := new(strings.Builder)

	for k := 0; k < len(node); k++ {
		switch {
		case node[k].DataAtom == atom.Rt:
			k = k + len(getNode(node[k:])) - 1

		case node[k].DataAtom == atom.Rp:
			k = k + len(getNode(node[k:])) - 1

		case isText(node[k]):
			w.WriteString(node[k].Data)

		default:
		}
	}

	return w.String()
}

func isChild(ns, s *section) bool {
	return headerLevel(ns.node) == headerLevel(s.node)+1
}

func isParent(ns, s *section) bool {
	return isChild(s, ns)
}

func isGrandParent(ns, s *section) bool {
	return headerLevel(s.node) == headerLevel(ns.node)+2
}

func maxLevel(tokens []*html.Token) atom.Atom {

	for _, t := range tokens {
		if isHeader(t) {
			if t.DataAtom == atom.H3 {
				return atom.H3
			}
		}
	}
	for _, t := range tokens {
		if isHeader(t) {
			if t.DataAtom == atom.H4 {
				return atom.H4
			}
		}
	}
	return atom.H5
}
