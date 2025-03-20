package aozoratext

import (
	"strconv"
	"strings"
)

/*
func newRubyBaseMarker() *azAtom {

	return &azAtom{
		aType: rubyBaseStartTag,
		aData: rubyBaseStartStr,
		aNext: nil,
		aPrev: nil,
	}

}
*/
func (tok *token) noteText() (r string) {

	if tok.tokType != gaijiToken {

		if tok.tokType != noteToken {

			return
		}
	}

	return tok.innerString()

}

func indexNotAnyOf(s string, set string) int {

	for k, e := range s {

		if !strings.ContainsAny(string(e), set) {

			return k

		}
	}

	return len(s)

}

func ucode(txt string) string {

	i := strings.Index(txt, "U+")

	if i == -1 {

		return emptyStr

	}

	_, err := hextoi(txt[i+2 : i+6])

	if err != nil {
		return emptyStr
	}

	return txt[i+2 : i+6]

}

func unicodeOf(txt string) string {

	n, err := hextoi(txt)

	if err != nil {
		return emptyStr
	}

	return string(rune(n))

}

func firstChar(s string) rune {

	for _, e := range s {
		return e
	}

	return 0

}

func jisCodeOf(txt string) string {

	i := strings.Index(txt, "1-")

	if i == -1 {

		i = strings.Index(txt, "2-")

	}

	if i == -1 {
		return emptyStr
	}

	j := indexNotAnyOf(txt[i:], "1234567890-")

	if j == -1 {
		return emptyStr
	}

	c := strings.Split(txt[i:i+j], "-")

	if len(c) != 3 {
		return emptyStr
	}

	if len(c[0]) > 1 {
		return emptyStr
	}

	if len(c[1]) > 2 {
		return emptyStr
	}

	if len(c[2]) > 2 {
		return emptyStr
	}

	return txt[i : i+j]
}

func addToStringsBuilder(b *strings.Builder, s ...string) {

	for _, e := range s {

		b.WriteString(e)

	}

	return
}

func itohex(n int) string {

	return strconv.FormatInt(int64(n), 16)

}

func hextoi(h string) (int, error) {

	n, e := strconv.ParseInt(h, 16, 64)

	return int(n), e
}
