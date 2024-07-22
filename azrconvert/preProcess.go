package azrconvert

import (
	"log"
	"strings"

	"github.com/adamay909/AozoraConvert/runes"
	"golang.org/x/text/encoding/japanese"
)

// ToUTF8 converts src to UTF-8. It assumes src is ShiftJIS.
func ToUTF8(src []byte) []byte {

	r, err := japanese.ShiftJIS.NewDecoder().String(string(src))
	if err != nil {
		log.Println(err)
	}
	log.Println("Converted to UTF-8.")
	return []byte(r)
}

func fixLineEndings(s []byte) (out []byte) {
	p := byte('X')
	for _, c := range s {
		switch c {
		case '\r':
			out = append(out, '\n')
			p = c
		case '\n':
			if p != '\r' {
				out = append(out, c)
			}
			p = c
		default:
			out = append(out, c)
			p = c
		}
	}
	log.Println("Fixed line endings to UNIX style.")

	return
}

// ToSJIS converts src to ShiftJIS. src must be UTF-8.
func ToSJIS(src []byte) []byte {

	r, err := japanese.ShiftJIS.NewEncoder().String(string(src))
	if err != nil {
		log.Println(err)
	}
	log.Println("Converted to ShiftJIS.")
	return []byte(r)
}

func fixLineBreaks(src []byte) []byte {
	var sout string
	b := runes.SplitIntoBlocks(runes.Runes(string(src)), runes.Runes("<body>"), runes.Runes("</body>"))
	for i := range b {
		if runes.HasPrefix(b[i], runes.Runes("<body>")) {
			b[i] = runes.Join(runes.Split(b[i], runes.Runes("\n")), runes.Runes(""))
			b[i] = runes.ReplaceAll(b[i], runes.Runes("<br>"), runes.Runes("<br />"))
			b[i] = runes.ReplaceAll(b[i], runes.Runes("<br />"), runes.Runes("<br />\n\n"))
			b[i] = runes.ReplaceAll(b[i], runes.Runes("</div>"), runes.Runes("</div>\n"))
		}
		sout = sout + b[i].String()
	}
	log.Println("Fixed line breaks to prevent spurious white spaces in some e-readers.")
	return []byte(sout)

}
func fixKunojiten(src []byte) []byte {

	b := string(src)
	b = strings.ReplaceAll(b, "／＼", "〳〵")
	b = strings.ReplaceAll(b, "／″＼", "〴〵")

	log.Println("Fixed kunojiten.")
	return []byte(b)
}

func modifyNotes(src []byte) []byte {

	s := string(src)

	s = strings.ReplaceAll(s, `※<span class="notes">`, `<span class="notes">※`)

	log.Println("Moved all ※ inside notes.")
	return []byte(s)
}
