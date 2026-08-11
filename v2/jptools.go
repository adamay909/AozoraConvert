package aozoraconvert

//go:generate stringer -type=CharTypeID

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var utf8of map[string]string

// CharTypeID represents character types.
type CharTypeID int

// Define charater types.
const (
	Symbol CharTypeID = 1 << iota //Symbol captures everything that isn't captured by the other categories.
	Hiragana
	Katakana
	Kanji
	Whitespace
	Punctuation
	Roman
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
	initialized = false
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

func fwnum(n int) string {

	hstr := strconv.Itoa(n)

	fstr := ""

	for _, d := range hstr {

		k, _ := strconv.Atoi(string(d))

		fstr = fstr + string(fwanumOf[k])

	}

	return fstr
}

// isKatakana checkes if r is katakana.
func isKatakana(r rune) bool {

	return CharType(r) == Katakana

}

// isHiragana checkes if r is hiragana,
func isHiragana(r rune) bool {

	return CharType(r) == Hiragana

}

func isKanji(r rune) bool {

	return CharType(r)&Kanji != 0

}

func sameCharType(r1, r2 rune) bool {

	return CharType(r1)&CharType(r2) != 0

}

// toKatakana convers r to katakana (iff. r is hiragana).
func toKatakana(r rune) rune {

	if !isHiragana(r) {

		return r
	}

	return r + 0x0060

}

// ToHiragana converts r to hiragana (iff. r is katakana)
func ToHiragana(r rune) rune {

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

func isFWNumeralString(s string) bool {

	for _, r := range s {

		if !isFWnumeral(r) {
			return false
		}
	}

	//if s is empty, response is trivially true
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

// CharType returns the character type of r
func CharType(r rune) CharTypeID {

	switch {
	case ' ' == r:
		return Whitespace

	case '　' == r:
		return Whitespace

	case '※' == r:
		return Kanji

	case '〇' == r:
		return Kanji

	case 'ヶ' == r:
		return Kanji

	case '〻' == r:
		return Kanji

	case '〆' == r:
		return Kanji

	case strings.ContainsAny("〳〴〵／″＼", string(r)):
		return Symbol

	case unicode.IsOneOf(kanjiR, r):
		return Kanji

	case unicode.IsOneOf(hiraganaR, r):
		return Hiragana

	case unicode.IsOneOf(katakanaR, r):
		return Katakana

	case strings.ContainsAny("、。「」！？・", string(r)):
		return Punctuation

	case unicode.IsOneOf([]*unicode.RangeTable{unicode.Latin}, r):
		return Roman

	default:
		return Symbol
	}
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
		fmt.Println("fail", err)
		return
	}

	s, ok := utf8of[jiscode]
	if !ok {
		err = errors.New("Unicode undefined for " + mkt)
		fmt.Println(err)
		return
	}
	return
}

// ConvertMKT returns the unicode string corresponding to the JIS code point
// given in the 面区点(men-ku-ten) format.
// mkt needs to be formatetted as a string of the form "d-dd-dd".
func ConvertMKT(mkt string) (s string, err error) {

	return convert(mkt)

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
