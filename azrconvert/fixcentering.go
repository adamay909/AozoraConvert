package azrconvert

import (
	"log"
	"strings"

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

		if !isNote(in[i]) {
			out = append(out, in[i])
			continue
		}

		if !isText(in[i+1]) {
			out = append(out, in[i])
			continue
		}

		if in[i+1].Data != `［＃ページの左右中央］` {
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
			if isEmptyLineBreak(in[i]) {
				brc++
				continue
			}

			out = append(out, in[i])

		}

		nn = new(html.Token)

		nn.Type = html.EndTagToken

		nn.DataAtom = atom.Div

		nn.Data = "div"

		out = append(out, nn)

		log.Println("Fixed centering.")

	}

	return
}
