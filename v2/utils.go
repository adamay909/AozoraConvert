package aozoraconvert

import (
	"strconv"
	"strings"
)

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

	k := 0
	for k = i + 2; k < len(txt) && strings.Contains(`0123456789ABCDEF`, string(txt[k])); k++ {

	}
	_, err := hextoi(txt[i+2 : k])

	if err != nil {
		return emptyStr
	}

	return txt[i+2 : k]

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

	i := strings.Index(txt, "水準1-")

	if i == -1 {

		i = strings.Index(txt, "水準2-")

	}

	if i == -1 {

		i0 := strings.Index(txt, "1-")
		i1 := strings.Index(txt, "2-")

		switch {

		case i0 == -1 && i1 == -1:
			return emptyStr

		case i0 == -1:
			i = i1

		case i1 == -1:
			i = i0

		case i0 < i1:
			i = i0

		case i1 < i0:
			i = i1
		}

		i = i - 6
	}

	i = i + 6 //len("水準")

	j := indexNotAnyOf(txt[i:], "1234567890-")

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
