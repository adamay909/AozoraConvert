package azrconvert

import (
	"strings"

	"github.com/adamay909/AozoraConvert/runes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func isProperLineBreak(t *html.Token) bool {

	if t.Type != html.SelfClosingTagToken {
		return false
	}

	if t.DataAtom == atom.Br {
		return true
	}

	return false

}

func isEmptyLineBreak(t *html.Token) bool {

	if t.Type != html.SelfClosingTagToken {
		return false
	}

	if t.DataAtom != atom.Br {
		return false
	}

	if len(t.Attr) == 0 {
		return false
	}

	if t.Attr[0].Key != "class" {
		return false
	}

	if t.Attr[0].Val != "blankline" {
		return false
	}

	return true

}

func isBodyStart(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom == atom.Body {
		return true
	}

	return false
}

func isBodyEnd(t *html.Token) bool {

	if t.Type != html.EndTagToken {
		return false
	}

	if t.DataAtom == atom.Body {
		return true
	}

	return false
}

func isBurasage(t *html.Token) bool {

	if !isDiv(t) {
		return false
	}

	return classOf(t) == "burasage"

}

func isChitsuki(t *html.Token) bool {

	if !isDiv(t) {
		return false
	}

	return classNameContains(t, "chitsuki")
}

func isDiv(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	return t.DataAtom == atom.Div
}

func isEmStart(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom == atom.Em {
		return true
	}

	return false
}

func isGaiji(t *html.Token) bool {

	if !isImg(t) {
		return false
	}

	j := jcodeOf(t)
	if len(j) > 0 {
		setAttr(t, "data-jcode", j)
		return true
	}

	s := getAttr(t, "src")
	if strings.HasPrefix(s, "../../../gaiji/") {
		return true
	}

	if classOf(t) == "gaiji" {
		return true
	}

	return false
}

func isHeader(t *html.Token) bool {
	if t.Type != html.StartTagToken {
		return false
	}
	switch t.DataAtom {
	case atom.H1:
		return true
	case atom.H2:
		return true
	case atom.H3:
		return true
	case atom.H4:
		return true
	case atom.H5:
		return true
	case atom.H6:
		return true
	default:
		return false
	}
}

func isImg(t *html.Token) bool {

	return t.DataAtom == atom.Img
}

func isIndex(t *html.Token) bool {

	if !isDiv(t) {
		return false
	}

	return getAttr(t, "id") == "contents"
}

func isJisage(t *html.Token) bool {

	if !isDiv(t) {
		return false
	}

	return classNameContains(t, "jisage")
}

func isMetadata(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom == atom.H1 {
		return true
	}

	if t.DataAtom == atom.H2 {
		return true
	}

	return false
}

func isNote(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom != atom.Span {
		return false
	}

	if classOf(t) != "notes" {
		return false
	}

	return true
}

func isRubyStart(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom == atom.Ruby {
		return true
	}

	return false
}

func isScriptStart(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom == atom.Script {
		return true
	}

	return false
}

func isCharNote(t *html.Token) bool {

	if t.Type != html.StartTagToken {
		return false
	}

	if t.DataAtom != atom.Span {
		return false
	}

	if classOf(t) != "charNote" {
		return false
	}

	return true
}

func isText(t *html.Token) bool {

	return t.Type == html.TextToken

}

func classOf(t *html.Token) string {

	for _, attr := range t.Attr {
		if strings.ToLower(attr.Key) == "class" {
			return attr.Val
		}
	}

	return ""
}

func classNameContains(t *html.Token, s string) bool {

	return strings.Contains(classOf(t), s)
}

func delAttr(t *html.Token, key string) {

	var nattr []html.Attribute

	for _, a := range t.Attr {
		if a.Key == strings.ToLower(key) {
			continue
		}
		nattr = append(nattr, a)
	}

	t.Attr = nattr

	return

}

func setAttr(t *html.Token, key, val string) {

	delAttr(t, key)

	t.Attr = append(t.Attr, html.Attribute{Key: key, Val: val})

	return
}

func getAttr(t *html.Token, key string) (val string) {

	for _, a := range t.Attr {
		if a.Key == strings.ToLower(key) {
			return a.Val
		}
	}

	return ""
}

func hasID(t *html.Token) bool {

	return len(getAttr(t, "id")) > 0

}

func getID(t *html.Token) string {

	return getAttr(t, "id")

}

func headerLevel(token *html.Token) int {
	if token == nil {
		return -1
	}

	if !isHeader(token) {
		return -1
	}
	switch token.DataAtom {
	case atom.H1:
		return 1
	case atom.H2:
		return 2
	case atom.H3:
		return 3
	case atom.H4:
		return 4
	case atom.H5:
		return 5
	case atom.H6:
		return 6
	default:
		return -1
	}
}

func jcodeOf(t *html.Token) (jcode string) {

	a := runes.Runes(getAttr(t, "alt"))
	if !runes.HasPrefix(a, runes.Runes("※")) {
		return
	}

	i := runes.Index(a, runes.Runes("1-"))

	if i == -1 {
		i = runes.Index(a, runes.Runes("2-"))
	}

	if i == -1 {
		return
	}

	j := i

	for j = i; strings.ContainsAny(string(a[j:j+1]), "0123456789-"); j++ {
	}

	s := string(a[i:j])

	c := strings.Split(s, "-")

	if len(c) != 3 {
		return
	}

	if len(c[0]) > 1 {
		return
	}

	if len(c[1]) > 2 {
		return
	}

	if len(c[2]) > 2 {
		return
	}

	jcode = s

	return
}
