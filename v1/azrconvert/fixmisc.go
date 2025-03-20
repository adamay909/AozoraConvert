package azrconvert

import (
	"log"
	"strings"

	"github.com/adamay909/AozoraConvert/runes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func markBlankLines(in []byte) []byte {

	lines := strings.Split(string(in), "\n")

	var newlines []string

	for _, l := range lines {

		if l == `<br />` {

			newlines = append(newlines, `<br class="blankline" />`)
			continue
		}

		newlines = append(newlines, l)

	}

	return []byte(strings.Join(newlines, "\n"))

}

func prettifyEmptyLines(in []byte) []byte {

	//	return in

	lines := strings.Split(string(in), "\n")

	var newlines []string

	ec := 0 //counter for consecutive empty lines

	for _, l := range lines {

		sl := strings.TrimSpace(l)

		if len(sl) > 0 {

			newlines = append(newlines, sl)

			ec = 0

			continue

		}

		ec++

		if ec > 1 {
			continue
		}

		newlines = append(newlines, sl)

	}

	return []byte(strings.Join(newlines, "\n"))

}

func fixCentering(in []*html.Token) (out []*html.Token) {

	for i := 0; i < len(in); i++ {

		if i == len(in)-1 {
			out = append(out, in[i])
			continue
		}

		if !isNote(in[i]) {
			out = append(out, in[i])
			continue
		}

		if !isText(in[i+1]) {
			out = append(out, in[i])
			continue
		}

		if in[i+1].Data != `［＃ページの左右中央］` {
			out = append(out, in[i])
			continue
		}

		i++

		nn := new(html.Token)

		nn.Type = html.StartTagToken

		nn.DataAtom = atom.Div

		nn.Data = "div"

		nn.Attr = append(nn.Attr, html.Attribute{"", "class", "centered"})

		out = append(out, nn)

		i++ //skip closing span tag

		for brc := 0; brc < 4; {

			i++

			if i == len(in) {
				break
			}

			if isEmptyLineBreak(in[i]) {
				brc++
			}

			out = append(out, in[i])

		}

		nm := new(html.Token)

		nm.Type = html.EndTagToken

		nm.DataAtom = atom.Div

		nm.Data = "div"

		out = append(out, nm)

		log.Println("Fixed centering.")

	}

	return
}

func fixKogaki(in []*html.Token) (out []*html.Token) {

	log.Println("fixing Kogaki")

	for i := 0; i < len(in); i++ {

		if i == len(in)-1 {
			out = append(out, in[i])
			continue
		}

		if !isCharNote(in[i]) {
			out = append(out, in[i])
			continue
		}

		if !isText(in[i+1]) {
			out = append(out, in[i])
			continue
		}

		if !strings.HasPrefix(in[i+1].Data, `※［＃小書き`) {
			out = append(out, in[i])
			continue
		}

		oldnode := getNode(in[i:])

		t0 := new(html.Token)
		t0.Type = html.StartTagToken
		t0.Data = "span"
		t0.DataAtom = atom.Span
		t0.Attr = []html.Attribute{{Key: "class", Val: "kogaki"}}

		t2 := new(html.Token)
		t2.Type = html.EndTagToken
		t2.Data = "span"
		t2.DataAtom = atom.Span

		t1 := new(html.Token)

		txt := runes.Runes(in[i+1].Data)[9:10]

		t1.Data = string(txt)
		t1.Type = html.TextToken

		out = append(out, t0)
		out = append(out, t1)
		out = append(out, t2)

		i = i + 2
		log.Println("replaced", renderTokens(oldnode), "with", renderTokens([]*html.Token{t0, t1, t2}))
	}

	return
}
