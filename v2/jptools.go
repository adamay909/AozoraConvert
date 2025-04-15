package aozoraConvert

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var utf8of map[string]string

// charTypeID represents character types.
type charTypeID int

// Define charater types.
const (
	symbol charTypeID = 1 << iota //Symbol captures everything that isn't captured by the other categories.
	hiragana
	katakana
	kanji
	whitespace
	punctuation
	roman
)

type jisuni struct {
	jis string
	uni string
}

var fwanum = []rune{
	'０',
	'１',
	'２',
	'３',
	'４',
	'５',
	'６',
	'７',
	'８',
	'９',
}

var hwanum = []rune{
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
}

var hexnum = []rune{
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
}

var upperHex = []rune{
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
}

var (
	hwanumOf, fwanumOf map[int]rune

	intOf map[rune]int

	intOfHex map[rune]int

	hexcharOf map[int]rune

	lowercase map[rune]rune
)

var initialized = false

func init() {

	initMap()

}

func initMap() {

	if initialized {
		return
	}

	utf8of = make(map[string]string)

	intOf = make(map[rune]int)

	intOfHex = make(map[rune]int)

	fwanumOf = make(map[int]rune)

	hwanumOf = make(map[int]rune)

	hexcharOf = make(map[int]rune)

	lowercase = make(map[rune]rune)

	for _, d := range data {

		utf8of[d.jis] = d.uni

	}

	for i, numeral := range fwanum {

		intOf[numeral] = i
		fwanumOf[i] = numeral

		intOf[hwanum[i]] = i
		hwanumOf[i] = hwanum[i]

	}

	for i := range hexnum {

		hexcharOf[i] = hexnum[i]
		intOfHex[hexnum[i]] = i

	}

	for i := 10; i < len(hexnum); i++ {

		lowercase[upperHex[i-10]] = hexnum[i]

	}

	initialized = true

}

// String returns the character type of c
// as a string. E.g. "Katakana" if c
// is katakana.
func (c charTypeID) String() string {
	switch c {
	case symbol:
		return "Symbol"
	case kanji:
		return "Kanji"
	case hiragana:
		return "Hiragana"
	case katakana:
		return "Katakana"
	case whitespace:
		return "Whitespace"
	default:
		return "other"
	}
}

// isKatakana checkes if r is katakana.
func isKatakana(r rune) bool {

	return charType(r) == katakana

}

// isHiragana checkes if r is hiragana,
func isHiragana(r rune) bool {

	return charType(r) == hiragana

}

func isKanji(r rune) bool {

	return charType(r)&kanji != 0

}

func sameCharType(r1, r2 rune) bool {

	return charType(r1)&charType(r2) != 0

}

// toKatakana convers r to katakana (iff. r is hiragana).
func toKatakana(r rune) rune {

	if !isHiragana(r) {

		return r
	}

	return r + 0x0060

}

// toHiragana converts r to hiragana (iff. r is katakana)
func toHiragana(r rune) rune {

	if !isKatakana(r) {
		return r
	}

	if r > 0x30f6 {
		return r
	}

	return r - 0x0060

}

func isFWnumeral(r rune) bool {

	if r < 0xFF10 {
		return false
	}

	if r > 0xFF19 {
		return false
	}

	return true

}

/*
func isHWnumeral(r rune) bool {

	return charType(r) == ArabNum

}

func isHexNumeral(r rune) bool {

	_, ok := intOfHex[r]

	return ok

}
*/
func lowerCaseOf(u rune) rune {

	l, ok := lowercase[u]

	if ok {
		return l
	}

	return u

}

var (
	kanjiR = []*unicode.RangeTable{unicode.Han}

	hiraganaR = []*unicode.RangeTable{unicode.Hiragana}

	katakanaR = []*unicode.RangeTable{unicode.Katakana}
)

func charType(r rune) charTypeID {

	switch {
	case ' ' == r:
		return whitespace

	case '　' == r:
		return whitespace

	case '※' == r:
		return kanji

	case '〇' == r:
		return kanji

	case 'ヶ' == r:
		return kanji

	case '〻' == r:
		return kanji

	case '〆' == r:
		return kanji

	case strings.ContainsAny("〳〴〵／″＼", string(r)):
		return symbol

	case unicode.IsOneOf(kanjiR, r):
		return kanji

	case unicode.IsOneOf(hiraganaR, r):
		return hiragana

	case unicode.IsOneOf(katakanaR, r):
		return katakana

	case strings.ContainsAny("、。「」！？・", string(r)):
		return punctuation

	case unicode.IsOneOf([]*unicode.RangeTable{unicode.Latin}, r):
		return roman

	default:
		return symbol
	}
}

func charType_(r rune) charTypeID {

	switch {

	case ' ' == r:
		return whitespace

	case '〇' == r:
		return kanji

	case 'ヶ' == r:
		return kanji

	case '〻' == r:
		return kanji

	case 0x203B == r:
		return kanji

	case 0xFB00 <= r:
		return symbol

	case 0xF900 <= r:
		return kanji

	case 0x9FF0 <= r:
		return symbol

	case 0x4E00 <= r:
		return kanji

	case 0x30FF <= r:
		return symbol

	case 0x30A1 <= r:
		return katakana

	case 0x30A0 <= r:
		return symbol

	case 0x3041 <= r:
		return hiragana

	case 0x3007 <= r:
		return symbol

	case 0x3005 <= r:
		return kanji
	}

	return symbol
}

// convert returns the unicode string corresponding to the
// JIS codepoint in the 面-区-点 (men-ku-ten) format.
// mkt needs to be formatetted as a string of the form "d-dd-dd".
// err is nil if conversion succeeds.
func convert(mkt string) (s string, err error) {

	if !initialized {
		initMap()
	}

	jiscode, err := mktToJis(mkt)
	if err != nil {
		return
	}

	s, ok := utf8of[jiscode]
	if !ok {
		err = errors.New("Unicode undefined for " + mkt)
		return
	}
	return
}

// mktToJis returns the JIS code point corresponding
// to the provided 面-区-点 (men-ku-ten) codepoint.
// mkt needs to be formatetted as a string of the form "d-dd-dd".
func mktToJis(mkt string) (s string, err error) {

	if !initialized {
		initMap()
	}
	fields := strings.Split(mkt, "-")

	if len(fields) != 3 {
		err = errors.New(mkt + ": is not valid 面-区-点 (men-ku-ten) format")

		return
	}

	switch fields[0] {
	case "1":
		s = "3-"
	case "2":
		s = "4-"
	default:
		err = errors.New("invalid 面-区-点 point")
		return

	}

	k, err := strconv.Atoi(fields[1])
	if err != nil {
		err = errors.New("invalid 面-区-点 point")
		return
	}

	s = s + itohex(k+32)

	t, err := strconv.Atoi(fields[2])
	if err != nil {
		err = errors.New("invalid 面-区-点 point")
		return
	}

	s = s + itohex(t+32)

	return

}
