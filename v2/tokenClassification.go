package aozoratext

import (
	"strings"
)

func (t *token) isPairMarkerOpen() bool {

	if t.tokType != noteToken {
		return false
	}

	for _, m := range pairMarker {

		if strings.HasSuffix(t.innerString(), m) {
			return true
		}

	}

	return false
}

func (t *token) isPairMarkerClose() bool {

	if t.tokType != noteToken {
		return false
	}

	if !strings.HasSuffix(t.innerString(), formatEndStr) {
		return false
	}

	for _, m := range pairMarker {

		if strings.HasSuffix(strings.TrimSuffix(t.innerString(), formatEndStr), m) {
			return true
		}
	}
	return false

}

func (t *token) isSectionTitleStart() bool {

	if t.tokType != noteToken {
		return false
	}

	for _, m := range sectionMarker {

		if strings.TrimPrefix(t.innerString(), blockStartStr) == m {
			return true
		}
	}
	return false
}

func (t *token) isSectionTitleEnd() bool {

	if t.tokType != noteToken {

		return false

	}

	if !strings.HasSuffix(strings.TrimPrefix(t.innerString(), blockEndStr), formatEndStr) {
		return false
	}

	for _, m := range sectionMarker {

		if strings.TrimSuffix(t.innerString(), formatEndStr) == m {
			return true
		}
	}

	return false

}

func (t *token) isIndentationStart() bool {

	if t.tokType != noteToken {
		return false
	}

	if !strings.HasPrefix(t.innerString(), blockStartStr) {
		return false
	}

	return strings.HasSuffix(t.innerString(), "字下げ")
}

func (t *token) matchingIndentationCloser() *token {

	pos := new(token)

	for pos = t.next; pos != nil; pos = pos.next {

		if pos.tokType != noteToken {
			continue
		}

		if pos.isIndentationStart() {
			return pos
		}

		if strings.HasSuffix(pos.innerString(), "字下げ終わり") {
			return pos
		}
	}

	panic(t.info() + " No matching closer.")

	return pos
}

func (t *token) isImage() bool {

	if t.tokType != noteToken {
		return false
	}

	if !strings.Contains(t.innerString(), "入る") {
		return false
	}

	return strings.Contains(t.innerString(), ".png")

}

func (t *token) isCaption() bool {

	if t == nil {
		return false
	}

	if t.tokType != noteToken {
		return false
	}

	return strings.HasSuffix(t.innerString(), "キャプション")

}

func (t *token) isPagination() bool {

	c := t.innerString()

	for _, m := range paginationMarker {

		if c == m {
			return true
		}
	}

	return false

}

func (t *token) isKunten() bool {

	if t.tokType != noteToken {
		return false
	}

	s := []rune(t.innerString())

	if len(s) > 2 {
		return false
	}

	ok := false

	//only allow kuntenchars as first character
	for _, c := range kuntenchars {

		if s[0] == c {
			ok = true
			break
		}
	}

	if !ok {
		return false
	}

	if len(s) == 1 {
		return true
	}

	//need to check for combined kunten

	//first char cannot be re-ten
	if s[0] == kuntenchars[0] {
		return false
	}

	//second char must be re-ten
	if s[1] != kuntenchars[0] {
		return false
	}

	return true
}

func (t *token) isOkurigana() bool {

	if t.tokType != noteToken {
		return false
	}

	s := t.innerString()

	if !strings.HasPrefix(s, "（") {
		return false
	}

	if !strings.HasSuffix(s, "）") {
		return false
	}

	return true
}

func (t *token) matchingCloserToken() *token {

	return matchingCloserToken(t)
}

func matchingCloserToken(t *token) *token {

	if t.next == nil {
		panic(t.info() + " missing matching closer")
	}

	pos := new(token)

	count := 1

	for pos = t.next; ; pos = pos.next {

		if pos.isPairMarkerOpen() {
			count++
		}

		if pos.isPairMarkerClose() {
			count--
		}

		if count == 0 {

			return pos

		}

		if pos.next == nil {
			panic(t.info() + " missing matching closer")
		}
	}

	return pos

}

func (t *token) isBlockStartNote() bool {

	if t.tokType != noteToken {
		return false
	}

	if strings.HasPrefix(t.innerString(), blockStartStr) {
		return true
	}

	if t.isSectionTitleStart() {
		return true
	}

	if t.isFormatOfType(centeringMarker) {
		return true
	}

	return false
}

func (t *token) isBlockClosingNote() bool {

	if t.tokType != noteToken {
		return false
	}

	if t.isSectionTitleEnd() {
		return true
	}

	if t.isFormatEndOfType(centeringMarker) {
		return true
	}

	if !strings.HasPrefix(t.innerString(), blockEndStr) {
		return false
	}

	return strings.HasSuffix(t.innerString(), formatEndStr)
}

func (t *token) isFormatOfType(m []string) bool {

	if t.tokType != noteToken {
		return false
	}

	s := strings.TrimPrefix(t.innerString(), blockStartStr)

	for _, e := range m {

		if strings.HasSuffix(s, e) {
			return true
		}
	}

	return false
}

func (t *token) isFormatEndOfType(m []string) bool {

	if t.tokType != noteToken {
		return false
	}

	s := strings.TrimPrefix(t.innerString(), blockEndStr)

	for _, e := range m {

		if strings.HasSuffix(s, e) {
			return true
		}
	}

	return false
}

func (t *token) isSectionStart() bool {

	switch t.tokType {

	case sectionToken:
		return true

	case subsectionToken:
		return true

	case subsubsectionToken:
		return true

	default:
		return false

	}
}

func (t *token) isSectionEnd() bool {

	switch t.tokType {

	case sectionEndToken:
		return true

	case subsectionEndToken:
		return true

	case subsubsectionEndToken:
		return true

	default:
		return false

	}
}

func (t *token) isBlock() bool {

	return strings.HasPrefix(t.innerString(), blockStartStr)

}

func (t *token) isBlockEnd() bool {

	if !strings.HasSuffix(t.innerString(), formatEndStr) {
		return false
	}

	return strings.HasPrefix(t.innerString(), blockEndStr)

}
