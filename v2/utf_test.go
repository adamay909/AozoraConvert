package aozoratext

import (
	"fmt"
	"testing"
	"unicode"
)

func TestUni(t *testing.T) {

	s := "「ムツセメルリ」は本文では「ムッセメルリ」"

	rangeT := []*unicode.RangeTable{unicode.Han}

	for _, r := range s {

		fmt.Println(string(r), "is Kanji", unicode.IsOneOf(rangeT, r))

	}

}
