// Package jptools provides a few simple functions for dealing with
// text in Japanese.
package jptools

import (
	"errors"
	"strconv"
	"strings"
)

// Utf8of maps JIS X 0213:2004 codepoints to Unicode codepoints.
var Utf8of map[string]string

// UnicodeOf returns the Unicode point as an ASCII escaped Go string.
// mkt must be provided in the 面-区-点 (men-ku-ten) format as an
// ASCII encoded string. mkt needs to be formatetted as a string
// of the form "d-dd-dd".
func UnicodeOf(mkt string) (s string, err error) {

	jiscode, err := MktToJis(mkt)
	if err != nil {
		return
	}
	s, ok := Utf8of[jiscode]
	if !ok {
		err = errors.New("Unicode undefined for " + mkt)
	}
	return
}

// Convert returns the unicode string corresponding to the
// JIS codepoint in the 面-区-点 (men-ku-ten) format.
// mkt needs to be formatetted as a string of the form "d-dd-dd".
// err is nil if conversion succeeds.
func Convert(mkt string) (s string, err error) {

	jiscode, err := MktToJis(mkt)
	if err != nil {
		return
	}

	s, ok := Utf8of[jiscode]
	if !ok {
		err = errors.New("Unicode undefined for " + mkt)
		return
	}
	return
}

// MktToJis returns the JIS code point corresponding
// to the provided 面-区-点 (men-ku-ten) codepoint.
// mkt needs to be formatetted as a string of the form "d-dd-dd".
func MktToJis(mkt string) (s string, err error) {

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

	s = s + strconv.FormatInt(int64(k+32), 16)

	t, err := strconv.Atoi(fields[2])
	if err != nil {
		err = errors.New("invalid 面-区-点 point")
		return
	}

	s = s + strconv.FormatInt(int64(t+32), 16)

	return

}
