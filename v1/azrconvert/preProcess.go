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

func convertLeftOverGaijiChuki(src []byte) []byte {

	srcStr := string(src)

	lines := strings.Split(srcStr, "\n")

	var linesNew []string

	for _, l := range lines {

		if !hasLeftOverGaijiChuki(l) {
			linesNew = append(linesNew, l)
			continue
		}

		linesNew = append(linesNew, fixGaijiChuki(l))

	}

	return []byte(strings.Join(linesNew, "\n"))
}

func fixLeftOverGaijiChuki(src []byte) []byte {

	srcStr := string(src)

	lines := strings.Split(srcStr, "\n")

	var linesNew []string

	for _, l := range lines {

		if !hasLeftOverGaijiChuki(l) {
			linesNew = append(linesNew, l)
			continue
		}

		linesNew = append(linesNew, fixGaijiChuki(l))

	}

	return []byte(strings.Join(linesNew, "\n"))
}

func hasLeftOverGaijiChuki(l string) bool {

	r := runes.Runes(l)

	i := runes.Index(r, runes.Runes("※［＃"))

	if i == -1 {
		return false
	}

	j := runes.Index(r[i:], runes.Runes("］"))

	if j <= -1 {
		return false
	}

	return true

}

func fixGaijiChuki(l string) string {

	r := runes.Runes(l)

	i := runes.Index(r, runes.Runes("※［＃"))

	j := runes.Index(r[i:], runes.Runes("］")) + i

	note := r[i : j+1]

	b := runes.NewBuilder()

	b.WriteRunes(r[:i])

	b.WriteRunes(runes.Runes(`<span class="notes">`))

	b.WriteRunes(note)

	b.WriteRunes(runes.Runes(`</span>`))

	if j+1 < len(r) {
		b.WriteRunes(r[j+1:])
	}

	log.Println("fixed left over gaiji chuki", note.String(), "; replaced with ", b.String())

	return b.String()

}
