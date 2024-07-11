package azrconvert

import (
	"bytes"
	"image"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adamay909/AozoraConvert/jptools"
	"github.com/adamay909/AozoraConvert/mobi/records"
	"github.com/adamay909/AozoraConvert/runes"
	"github.com/google/uuid"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Book represents a book from Aozora Bunko
type Book struct {
	Title, Creator, Publisher string
	Files                     []fileData
	UUID                      string
	Body                      []*html.Token
	Preamble                  []*html.Token
	URI                       string
	TopSection                *section
	CoverImage                image.Image
	Images                    []records.ImageRecord
	CSS                       string
	Log                       string
}

// NewBook returns a new Book.
func NewBook() *Book {
	b := new(Book)
	b.UUID = uuid.NewString()
	return b
}

// NewBookFrom returns a Book based on d. d is assumed to be
// xhtml formatted book from Aozora Bunko.
func NewBookFrom(d []byte) *Book {

	wlog := log.Writer()
	clog := new(strings.Builder)
	nlog := io.MultiWriter(wlog, clog)
	log.SetOutput(nlog)

	bk := NewBook()

	d = fixLineEndings(d)

	d = ToUTF8(d)

	d = fixLineBreaks(d)

	d = fixKunojiten(d)

	d = modifyNotes(d)

	tokens := tokenize(d)
	log.Println("Parsed and tokenized document.")

	bk.Preamble = getPreamble(tokens)

	bk.Body = getBody(tokens)

	bk.TopSection = getStructure(bk.Body)
	if bk.TopSection.firstChild == nil && bk.TopSection.nextSibling == nil {
		bk.TopSection = nil
	}

	bk.Log = clog.String()
	log.SetOutput(nlog)

	return bk
}

// SetURI sets the path of book within
// Aozora Bunko's file structure.
func (b *Book) SetURI(l string) {
	b.URI = l
	return
}

func (b *Book) GetURI() string {
	return b.URI

}

// SetTitle sets the title to t.
func (b *Book) SetTitle(t string) {
	b.Title = t
	return
}

// SetCreator stes the creator to c.
func (b *Book) SetCreator(c string) {
	b.Creator = c
	return
}

// SetPublisher sets the publisher to p.
func (b *Book) SetPublisher(p string) {
	b.Publisher = p
	return
}

func tokenize(d []byte) (out []*html.Token) {

	t := html.NewTokenizer(bytes.NewReader(d))

	for t.Next() != html.ErrorToken {
		tok := t.Token()
		out = append(out, &tok)
	}
	return out
}

func getPreamble(tokens []*html.Token) (preamble []*html.Token) {

	for i, t := range tokens {

		if isBodyStart(t) {
			return tokens[:i]
		}
	}

	return

}

func getBody(tokens []*html.Token) (body []*html.Token) {

	var s, e int

	for i, t := range tokens {
		if isBodyStart(t) {
			s = i
			continue
		}
		if isBodyEnd(t) {
			e = i
			break
		}
	}

	body = tokens[s : e+1]
	log.Println("Got html body of document.")

	body = fixNodes(body)

	body = fixTokens(body)

	insertSectionID(body)

	return body

}

func fixNodes(in []*html.Token) (out []*html.Token) {

	for i := 0; i < len(in); i++ {

		t := in[i]

		switch {

		case isMetadata(t):
			node := getNode(in[i:])
			out = append(out, fixMetadata(node)...)
			i = i + len(node) - 1

		case isHeader(t):
			node := getNode(in[i:])
			if isDiv(in[i-1]) {
				if in[i+len(node)].DataAtom == atom.Div && in[i+len(node)].Type == html.EndTagToken {
					out = out[:len(out)-1]
					out = append(out, node...)
					i = i + len(node)
					log.Println("Removed headers enclosed inside div. Styling should be done via css for h3, h4, etc.")
					continue
				}
			}
			out = append(out, node...)
			i = i + len(node) - 1

		case isRubyStart(t):
			node := getNode(in[i:])
			out = append(out, fixRubyNode(node)...)
			i = i + len(node) - 1

		case isEmStart(t):
			node := getNode(in[i:])
			out = append(out, fixEmNode(node)...)
			i = i + len(node) - 1

		case isNote(t):
			node := getNode(in[i:])
			out = append(out, fixNote(node)...)
			i = i + len(node) - 1

		case isScriptStart(t):
			node := getNode(in[i:])
			i = i + len(node) - 1
			log.Println("Removed script: ", renderTokens(node))

		case isIndex(t):
			node := getNode(in[i:])
			i = i + len(node) - 1
			log.Println("Removed index: ", renderTokens(node))

		default:
			out = append(out, t)
		}
	}
	return
}

func fixTokens(in []*html.Token) (out []*html.Token) {

	for i := 0; i < len(in); i++ {
		t := in[i]

		switch {
		case isJisage(t):
			out = append(out, fixJisage(t))

		case isChitsuki(t):
			out = append(out, fixChitsuki(t))

		case isBurasage(t):
			out = append(out, fixBurasage(t))

		case isGaiji(t):
			out = append(out, fixGaiji(t))

		default:
			out = append(out, in[i])
		}
	}
	return
}

func fixNote(oldNode []*html.Token) (newNode []*html.Token) {

	if strings.Contains(oldNode[1].String(), "U+") {
		r := runes.Runes(oldNode[1].String())
		i := runes.Index(r, runes.Runes("U+"))
		c, _ := strconv.Unquote(`"` + `\u` + string(r[i+2:i+6]) + `"`)
		nn := new(html.Token)
		nn.Data = c
		nn.DataAtom = 0
		nn.Type = html.TextToken
		newNode = append(newNode, nn)
		log.Println("Gaiji. Replaced ", renderTokens(oldNode), "with", renderTokens(newNode))
		return
	}

	if strings.Contains(oldNode[1].String(), "※") {
		delAttr(oldNode[0], "class")
		setAttr(oldNode[0], "class", "charNote")
		log.Println("Detected character note: ", renderTokens(oldNode))
		return oldNode
	}

	if strings.Contains(oldNode[1].String(), "［＃改丁］") {
		newtxt := `<p class="kaiTyou"></p>`
		newNode = tokenize([]byte(newtxt))
		log.Println("Detected page break. Replaced ", renderTokens(oldNode), "with", renderTokens(newNode))
		return
	}

	if strings.Contains(oldNode[1].String(), "［＃改ページ］") {
		newtxt := `<p class="kaiPeiji"></p>`
		newNode = tokenize([]byte(newtxt))
		log.Println("Detected page break. Replaced ", renderTokens(oldNode), "with", renderTokens(newNode))
		return
	}

	if strings.Contains(oldNode[1].String(), "［＃改見開き］") {
		newtxt := `<p class="kaiMihiraki"></p>`
		newNode = tokenize([]byte(newtxt))
		log.Println("Detected page break. Replaced ", renderTokens(oldNode), "with", renderTokens(newNode))
		return
	}

	return oldNode
}

func fixJisage(t *html.Token) *html.Token {

	n := strings.TrimPrefix(classOf(t), "jisage_")

	t.Attr = nil

	setAttr(t, "class", "jisage_"+n)
	setAttr(t, "style", "margin-top: "+n+"em")
	log.Println("Fixed jisage styling.")
	return t
}

func fixChitsuki(t *html.Token) *html.Token {

	n := strings.TrimPrefix(classOf(t), "chitsuki_")

	t.Attr = nil

	setAttr(t, "class", "chitsuki_"+n)
	setAttr(t, "style", "text-align: end; margin-bottom: "+n+"em")
	log.Println("Fixed chitsuki styling.")

	return t
}

func fixBurasage(t *html.Token) *html.Token {

	style := getAttr(t, "style")

	style = strings.ReplaceAll(style, "margin-left", "margin-top")

	delAttr(t, "style")

	setAttr(t, "style", style)

	log.Println("Fixed burasage styling.")

	return t
}

func fixGaiji(t *html.Token) *html.Token {

	jcode := getAttr(t, "jcode")

	if len(jcode) == 0 {
		fn := getAttr(t, "src")
		jcode = strings.TrimSuffix(filepath.Base(fn), ".png")
	}

	r, err := jptools.Convert(jcode)
	if err != nil {
		log.Println("Gaiji. Convert failed for ", t, " Error was: ", err)
		return t
	}

	ot := t.String()

	t.Type = html.TextToken
	t.DataAtom = atom.Atom(int32(0))
	t.Attr = nil

	t.Data = r

	log.Println("Gaiji. Replaced ", ot, "with ", r)
	return t
}

func fixRubyNode(oldNode []*html.Token) (newNode []*html.Token) {

	for i := 0; i < len(oldNode); i++ {

		switch {
		case oldNode[i].DataAtom == atom.Rb:

		case oldNode[i].DataAtom == atom.Rp:
			rpNode := getNode(oldNode[i:])
			i = i + len(rpNode) - 1

		default:
			newNode = append(newNode, oldNode[i])
		}
	}

	return
}

func fixMetadata(oldNode []*html.Token) (newNode []*html.Token) {

	newNode = oldNode

	attr := oldNode[0].Attr

	nn := new(html.Token)

	nn.Type = html.StartTagToken

	nn.DataAtom = atom.Div

	nn.Attr = attr

	nn.Data = "div"

	newNode[0] = nn

	nn = new(html.Token)

	nn.Type = html.EndTagToken

	nn.Data = "div"

	newNode[len(newNode)-1] = nn

	log.Println("Converted metadata tokens from h1, h2 to div.")

	return newNode
}

func fixEmNode(oldNode []*html.Token) (newNode []*html.Token) {

	emphType := classOf(oldNode[0])

	switch {

	case strings.HasPrefix(emphType, "underline"):
		newNode = fixEmphline(oldNode)

	case strings.HasPrefix(emphType, "overline"):
		newNode = fixEmphline(oldNode)

	case strings.HasSuffix(emphType, "after"):
		newNode = fixEmph(oldNode, true)

	default:
		newNode = fixEmph(oldNode, false)
	}
	log.Println("Fixed emph to be kindle friendly. ", renderTokens(oldNode), "to", renderTokens(newNode))

	return newNode
}

func fixEmphline(oldNode []*html.Token) (newNode []*html.Token) {

	newNode = oldNode

	newNode[0].DataAtom = atom.Span

	newNode[len(oldNode)-1].DataAtom = atom.Span

	return newNode
}

func fixEmph(oldNode []*html.Token, left bool) (newNode []*html.Token) {

	text := []rune(renderTokens(oldNode[1 : len(oldNode)-1]))
	style := botenStyle(classOf(oldNode[0]))

	w := new(strings.Builder)

	if !left {
		for _, c := range text {
			w.WriteString(`<ruby class="right-boten">`)
			w.WriteRune(c)
			w.WriteString(`<rt>`)
			w.WriteString(style)
			w.WriteString(`</rt></ruby>`)
		}
	} else {
		for _, c := range text {
			w.WriteString(`<ruby class="left-boten">`)
			w.WriteRune(c)
			w.WriteString(`<rt>`)
			w.WriteString(style)
			w.WriteString(`</rt></ruby>`)
		}
	}

	newNode = tokenize([]byte(w.String()))

	return newNode

}

func botenStyle(class string) string {

	switch strings.TrimSuffix(class, "_after") {

	case "sesame_dot":
		return `﹅`

	case "white_sesame_dot":
		return `﹆`

	case "black_circle":
		return `●`

	case "white_circle":
		return `○`

	case "black_up-pointing_triangle":
		return `▲`

	case "white_up-pointing_triangle":
		return `△`

	case "bullseye":
		return `◎`

	case "fisheye":
		return `⦿`

	case "saltire":
		return `'×'`

	default:
		return ""
	}

}

func getNode(tokens []*html.Token) (node []*html.Token) {

	if tokens[0].Type == html.SelfClosingTagToken {
		return tokens[:1]
	}

	if tokens[0].Type != html.StartTagToken {
		return node
	}

	nt := tokens[0].DataAtom

	for nesting, i := 1, 1; i < len(tokens); i++ {

		if tokens[i].DataAtom == nt {
			if tokens[i].Type == html.StartTagToken {
				nesting++
				continue
			}
			if tokens[i].Type == html.EndTagToken {
				nesting--
			}
			if nesting == 0 {
				node = tokens[0 : i+1]
				break
			}
		}
	}
	return node
}

func mkNewNode(dataAtom atom.Atom) (node []*html.Token) {

	n1 := new(html.Token)
	n1.Type = html.StartTagToken
	n1.DataAtom = dataAtom
	n1.Data = dataAtom.String()

	n2 := new(html.Token)
	n2.Type = html.EndTagToken
	n2.DataAtom = dataAtom
	n2.Data = dataAtom.String()

	node = append(node, n1, n2)

	return
}

func (b *Book) SetMetadataFromPreamble() {

	for _, t := range b.Preamble {

		if t.DataAtom == atom.Meta {

			if t.Attr[0].Key == "name" && t.Attr[0].Val == "DC.Title" {
				b.SetTitle(t.Attr[1].Val)
			}

			if t.Attr[0].Key == "name" && t.Attr[0].Val == "DC.Creator" {
				b.SetCreator(t.Attr[1].Val)
			}

			if t.Attr[0].Key == "name" && t.Attr[0].Val == "DC.Publisher" {
				b.SetPublisher(t.Attr[1].Val)
			}

		}

	}

	return
}
