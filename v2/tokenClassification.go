package aozoratext

import (
	"fmt"
	"strings"
)

const (
	normal = iota << 1
	section
	emph
	blockFormat
	pagination
)

// jisage etc.
var indentationMarker = []string{
	"字下げ",
}

var narrowparMarker = []string{
	"字詰め",
}

var bottomalignMarker = []string{
	"地付き",
	"字上げ",
}

// sectioning
var sectionMarker = []string{
	"大見出し",
	"中見出し",
	"小見出し",
}

var inlineSectionMarker = []string{
	"同行中見出し",
	"同行大見出し",
	"同行小見出し",
}

var windowSectionMarker = []string{
	"窓中見出し",
	"窓大見出し",
	"窓中見出し",
}

//emphasis

var lineDecoMarker = []string{
	"二重傍線",
	"傍線",
	"鎖線",
	"破線",
	"波線",
	"罫囲み",
}

var decoMarker = []string{
	"白ゴマ傍点",
	"白丸傍点",
	"黒三角傍点",
	"白三角傍点",
	"二重丸傍点",
	"蛇の目傍点",
	"丸傍点",
	"ばつ傍点",
	"傍点",
}

var decoString = []string{
	"﹆",
	"◦",
	"▲",
	"△",
	"◎ ",
	"◉",
	"•",
	"✕",
	"﹅",
}

// notes
var inlineNoteMarker = []string{
	"割り注",
}

var rubylikeNoteSimpleMarker = []string{
	"ルビ",
	"注記",
}

var rubylikeNoteMarker = []string{
	"ルビ付き",
	"注記付き",
}

var captionMarker = []string{
	"キャプション",
}

// page level formatting
var centeringMarker = []string{
	"ページの左右中央",
}

var paginationMarker = []string{
	"改丁",
	"改ページ",
	"改見開き",
	"改段",
}

//fontspec

var fontShapeMarker = []string{
	"太字",
	"斜体",
}

var fontsizeMarker = []string{

	"大きな文字",
	"小さな文字",
}

var _fontSizeMarker = []string{
	"段階大きな文字",
	"段階小さな文字",
}

//subscript

var offsetMarker = []string{
	"上付き小文字",
	"下付き小文字",
	"行右小書き",
	"行左小書き",
	"上付き小文字",
	"下付き小文字",
}

//text direction

var directionMarker = []string{
	"縦中横",
	"横組み",
}

var warichuLineBreakMarker = []string{
	"改行",
}

var _sectionLevelMarker = []string{
	"大区分",
	"中区分",
	"小区分",
}

var impliedCloserMarker = []string{}

var formattingMarker = []string{}

var rubylikeMarker = []string{}

var impliedOpenerMarker = []string{}

var pairMarker []string

func init() {

	includes := [][]string{
		indentationMarker,
		narrowparMarker,
		bottomalignMarker,
		captionMarker,
		inlineSectionMarker,
		windowSectionMarker,
		sectionMarker,
		lineDecoMarker,
		decoMarker,
		inlineNoteMarker,
		rubylikeNoteSimpleMarker,
		rubylikeNoteMarker,
		fontShapeMarker,
		fontsizeMarker,
		offsetMarker,
		directionMarker,
		centeringMarker,
	}

	for _, s := range includes {

		pairMarker = append(pairMarker, s...)

	}

	includes = [][]string{
		indentationMarker,
		bottomalignMarker,
	}

	for _, s := range includes {

		impliedCloserMarker = append(impliedCloserMarker, s...)

	}

	impliedOpenerMarker = append(impliedOpenerMarker, pairMarker...)

}

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

func (t *token) isPairClose() bool {

	if t.tokType != noteToken {
		return false
	}

	if !strings.HasSuffix(t.innerString(), "終わり") {
		return false
	}

	for _, m := range pairMarker {

		if strings.HasSuffix(strings.TrimSuffix(t.innerString(), "終わり"), m) {
			return true
		}
	}
	return false

}

func isMarker(txt string, markerset []string) string {

	for _, s := range markerset {

		if strings.HasSuffix(txt, s) {
			return s
		}
	}

	return ""

}

func (t *token) getFormattingMarker() string {

	for _, s := range impliedOpenerMarker {

		if strings.HasSuffix(t.innerString(), s) {
			return s
		}
	}

	return ""

}

func (t *token) isSectionTitleStart() bool {

	if t.tokType != noteToken {
		return false
	}

	for _, m := range sectionMarker {

		if t.innerString() == m {
			return true
		}
	}

	return false
}

func (t *token) isSectionTitleEnd() bool {

	if t.tokType != noteToken {

		return false

	}

	if !strings.HasSuffix(t.innerString(), "終わり") {
		return false
	}

	for _, m := range sectionMarker {

		if strings.TrimSuffix(t.innerString(), "終わり") == m {
			return true
		}
	}

	return false

}

func (t *token) isIndentationStart() bool {

	if t.tokType != noteToken {
		return false
	}

	if !strings.HasPrefix(t.innerString(), "ここから") {
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

	panic(t.String() + " has no matching closer")

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

var kuntenchars = []rune("レ一二三四五六七八九十上中下甲乙丙丁天地人元亨利貞乾坤")

func (t *token) isKunten() bool {

	if t.tokType != noteToken {
		return false
	}

	//	fmt.Println("check", t.tokType, printContext(t, 5))

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
		fmt.Println(printContext(t, 10))
		panic("malformed aozora text 1")
	}

	pos := new(token)

	count := 1

	for pos = t.next; ; pos = pos.next {

		if pos.isPairMarkerOpen() {
			count++
		}

		if pos.isPairClose() {
			count--
		}

		if count == 0 {

			return pos

		}

		if pos.next == nil {
			fmt.Println(printContext(t, 10))
			panic("malformed aozora text 2")
		}
	}

	return pos

}

func (t *token) isBlockFormatterStart() bool {

	return strings.HasPrefix(t.innerString(), "ここから")

}

func (t *token) hasCloser(start, end string) bool {

	for e := t.next; e != nil; e = e.next {

		if e.tokType != noteToken {
			continue
		}

		if e.innerString() == end {
			return true
		}

	}

	return false
}

func (t *token) isBlockStartNote() bool {

	if t.tokType != noteToken {
		return false
	}

	if strings.HasPrefix(t.innerString(), "ここから") {
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

	if !strings.HasPrefix(t.innerString(), "ここで") {
		return false
	}

	return strings.HasSuffix(t.innerString(), "終わり")
}

func (t *token) isFormatOfType(m []string) bool {

	if t.tokType != noteToken {
		return false
	}

	s := strings.TrimPrefix(t.innerString(), "ここから")

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

	s := strings.TrimPrefix(t.innerString(), "ここで")

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
