package aozoraConvert

import (
	"fmt"
	"strconv"
	"strings"
)

type token struct {
	tokType         tokenType
	next            *token
	prev            *token
	ref             string
	content         string
	jis0213Content  string
	unicodeContent  string
	originalContent string
	id              int
	lineNo          int
	modified        bool
	inserted        bool
	// extra   string
}

type tokenType int8

// These define types of tokens
const (
	emptyToken tokenType = iota
	textToken
	gaijiToken
	rubyGroupStartToken
	rubyGroupEndToken
	rubyParentStartToken
	rubyParentEndToken
	rubyStartToken
	rubyEndToken
	noteToken
	bibInfoToken
	bibInfoEndToken
	emptyLineToken
	endOfLineToken
	paragraphEndToken
	paragraphToken
	sectionTitleStartToken
	sectionTitleEndToken
	sectionToken
	subsectionToken
	subsubsectionToken
	sectionEndToken
	subsectionEndToken
	subsubsectionEndToken
	figureStartToken
	figureEndToken
	gaijiNoteToken
	specialCharToken
	accentToken
	gaijiCharToken
	kunojiToken
	alignBottomCloserToken
	centeringEndToken
	markupNoteToken
	mainTextStartToken
	mainTextEndToken
	eofToken
	gaijiImgToken
)

type tokenSubType int

func (t tokenType) String() string {

	switch t {

	case emptyToken:
		return "emptyToken"

	case textToken:
		return "textToken"

	case gaijiToken:
		return "gaijiToken"

	case rubyGroupStartToken:
		return "rubyGroupStartToken"

	case rubyGroupEndToken:
		return "rubyGroupEndToken"

	case rubyParentStartToken:
		return "rubyParentStartToken"

	case rubyParentEndToken:
		return "rubyParentEndToken"

	case rubyStartToken:
		return "rubyStartToken"

	case rubyEndToken:
		return "rubyEndToken"

	case noteToken:
		return "noteToken"

	case bibInfoToken:
		return "bibInfoToken"

	case bibInfoEndToken:
		return "bibInfoEndToken"

	case emptyLineToken:
		return "emptyLineToken"

	case endOfLineToken:
		return "endOfLineToken"

	case paragraphEndToken:
		return "paragraphEndToken"

	case paragraphToken:
		return "paragraphToken"

	case sectionTitleStartToken:
		return "sectionTitleStartToken"

	case sectionTitleEndToken:
		return "sectionTitleEndToken"

	case sectionToken:
		return "sectionToken"

	case subsectionToken:
		return "subsectionToken"

	case subsubsectionToken:
		return "subsubsectionToken"

	case sectionEndToken:
		return "sectionEndToken"

	case subsectionEndToken:
		return "subsectionEndToken"

	case subsubsectionEndToken:
		return "subsubsectionEndToken"

	case figureStartToken:
		return "figureStartToken"

	case figureEndToken:
		return "figureEndToken"

	case gaijiNoteToken:
		return "gaijiNoteToken"

	case specialCharToken:
		return "specialCharToken"

	case accentToken:
		return "accentToken"

	case gaijiCharToken:
		return "gaijiCharToken"

	case kunojiToken:
		return "kunojiToken"

	case alignBottomCloserToken:
		return "alignBottomCloserToken"

	case centeringEndToken:
		return "centeringEndToken"

	case markupNoteToken:
		return "markupNoteToken"

	case mainTextStartToken:
		return "mainTextStartToken"

	case mainTextEndToken:
		return "mainTextEndToken"

	case eofToken:
		return "eofToken"

	case gaijiImgToken:
		return "gaijiImageToken"

	default:

		return strconv.Itoa(int(t))

	}
}

func newToken() *token {

	t := new(token)

	t.inserted = false

	return t
}

func (t *token) addTokenRight(t2 *token) {

	if t2 == nil {
		return
	}

	t2.next = t.next

	if t2.next != nil {
		t2.next.prev = t2
	}

	t2.prev = t

	t.next = t2

}

func (t *token) insertTokenRight(t2 *token) {

	t2.next = t.next

	if t2.next != nil {
		t2.next.prev = t2
	}

	t2.prev = t

	t.next = t2

	t2.inserted = true

}

func (t *token) insertTokenLeft(t2 *token) {

	e := t.prev

	t.prev = t2

	t2.next = t

	t2.prev = e

	if e != nil {

		e.next = t2
	}

	t2.inserted = true

	return

}

func (t *token) joinTokens(t2 *token) {

	if t == nil {
		t = t2
	}

	tend := t.lastToken()

	t2start := t2.firstToken()

	tend.next = t2start

	t2start.prev = tend

	return

}

func (t *token) remove() {

	e1 := t.prev

	e2 := t.next

	if e1 != nil {
		e1.next = e2
	}

	e2.prev = e1

	t.prev = nil
	t.next = nil

	return

}

func (t *token) firstToken() *token {

	e := new(token)

	for e = t; e.prev != nil; e = e.prev {
	}

	return e

}

func (t *token) lastToken() *token {

	e := t

	for e.next != nil {
		e = e.next
	}

	return e

}

func (t *token) String() string {

	if t == nil {

		return "<NIL>"

	}

	switch t.tokType {

	case emptyToken:
		return ""

	case emptyLineToken:
		return ""

	case endOfLineToken:
		return ""

	case paragraphEndToken:
		return ""

		//	case rubyParentStartToken:
		//		return rubyBaseStartStr

	}

	return t.content
}

func (t token) StringAll() string {

	builder := new(strings.Builder)

	for s := t.firstToken(); ; s = s.next {

		builder.WriteString(s.String())

		if s.next == nil {
			break
		}

	}

	return builder.String()

}

func (t *token) innerString() string {

	if t == nil {
		return t.String()
	}

	switch t.tokType {

	case noteToken, gaijiImgToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.String(), noteStartStr), noteEndStr)

	case gaijiToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.String(), gaijiMarkerStr), noteEndStr)

	case gaijiNoteToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.String(), noteStartStr), noteEndStr)

	case accentToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.String(), accentStartStr), accentEndStr)

	default:

		return t.String()

	}
}

func (t *token) setString(str string) {

	t.content = str

	return

}

func (t *token) setInnerString(str string) {

	switch t.tokType {

	case noteToken:

		b := new(strings.Builder)

		addToStringsBuilder(b, noteStartStr, str, noteEndStr)

		t.content = b.String()

	default:

		t.content = str
	}

	return

}

func newNote(str string) *token {

	t := newToken()

	t.tokType = noteToken

	t.setString("［＃" + str + "］")

	return t

}

func (t *token) refString() (ref string) {

	if t.tokType != noteToken {
		return emptyStr
	}

	return getRefStrings(t.innerString())[0]

}

func getRefStrings(s string) []string {

	var rs []string

	for i1 := strings.Index(s, refStartStr); i1 != -1; i1 = strings.Index(s, refStartStr) {

		i2 := matchingCloserStringIndex(s, refStartStr, refEndStr)

		if i2 == -1 {
			break
		}

		rs = append(rs, s[i1+len(refStartStr):i2])

		s = s[i2+len(refEndStr):]

	}

	if len(rs) == 0 {
		rs = append(rs, emptyStr)
	}

	return rs

}

func classof(str string) string {

	return str[len(string(firstChar(str))):]

}

func newEmptyToken() *token {

	t := newToken()

	t.tokType = emptyToken

	return t

}

func newEolToken() *token {

	t := newToken()

	t.tokType = endOfLineToken

	return t

}

func newLineBreakToken() *token {

	t := newToken()

	t.tokType = paragraphEndToken

	return t
}

func newParagraphToken() *token {

	t := newToken()

	t.tokType = paragraphToken

	return t

}

func (t *token) info() string {

	output := new(strings.Builder)

	addToStringsBuilder(output, "line ", strconv.Itoa(t.lineNumber()), ": ", t.tokType.String(), ": ", t.String())

	if t.inserted {
		addToStringsBuilder(output, "(inserted)")
	}

	if t.modified {
		addToStringsBuilder(output, "(modified)")
		if t.originalContent != "" {
			addToStringsBuilder(output, " original string: ", t.originalContent)
		}
	}

	if t.unicodeContent != "" {
		addToStringsBuilder(output, "unicode: ", t.unicodeContent)
	}

	return strings.ReplaceAll(output.String(), "\n", "\\n")

}

func (t *token) addTokenBefore(txt string, nt *token) {

	e := t.prev

	c := 0

	cr := 0

	nesting := 0

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == endOfLineToken {
			break
		}

		if e.tokType == rubyEndToken {

			for ; e.tokType != rubyStartToken; e = e.prev {
			}

			e = e.prev
		}

		if e.isFormatCloseToken() {
			nesting++
			continue
		}

		if e.tokType != textToken {
			continue
		}

		cr = c

		c = len(e.content) + cr

		if c >= len(txt) {

			break

		}

	}

	if e == nil {

		panic("Can't find place to insert implied opener note. Defaulting to start of line. " + e.info())

		return

	}

	if e.tokType == endOfLineToken {

		panic("Can't find place to insert implied opener note. Defaulting to start of line. " + e.info())

		return

	}

	if c == len(txt) {

		if e.prev != nil && e.prev.tokType == rubyParentStartToken {

			e = e.prev.prev

		}

		for range nesting {

			e = e.prev

		}
		e.insertTokenLeft(nt)

		return

	}

	k := len(txt) - cr

	t2 := newTokenOfType(textToken)

	//	s1 := e.content[:len(e.content)-k]

	t2.setString(e.content[:len(e.content)-k])

	e.originalContent = e.content

	e.setString(e.content[len(e.content)-k:])

	e.modified = true

	e.insertTokenLeft(nt)

	e.prev.insertTokenLeft(t2)

	return
}

// Print the strings of up to n tokens prior and after t.
// t.String() is surrounded by "***".
func printContext(t *token, n int) string {

	fmt.Println(t)

	t1 := new(token)

	t2 := new(token)

	c := 0

	for t1, c = t, 0; c < n; t1 = t1.prev {

		if t1.prev == nil {
			break
		}

		if t1.prev.String() == "\n" {
			continue
		}

		if t1.prev.String() == "" {
			continue
		}

		c++

	}

	for t2, c = t, 0; c < n; t2 = t2.next {

		if t2.next == nil {
			break
		}

		if t2.next.String() == "\n" {
			continue
		}

		if t2.next.String() == "" {
			continue
		}

		c++

	}

	w := new(strings.Builder)

	for t0 := t1; ; t0 = t0.next {

		if t0 == t {

			w.WriteString("***" + t0.String() + "***")

		} else {

			w.WriteString(t0.String())

		}

		if t0 == t2 {
			break
		}

	}

	return w.String()

}

func newTokenOfType(s tokenType) *token {

	t := newToken()

	t.tokType = s

	t.inserted = false

	return t
}

func numberTokens(t *token) {

	for e, counter := t.firstToken(), 1; e != nil; e = e.next {

		counter++

		e.id = counter

	}

}

func (t *token) lineNumber() int {

	for e := t; e != nil; e = e.prev {

		if e.lineNo != 0 {

			return e.lineNo

		}

	}

	return 1

}

func (t *token) lastTokenInLine() *token {

	if t == nil {
		return t
	}

	e := new(token)

	for e = t; e != nil; e = e.next {

		if e.next == nil {
			return e
		}

		if e.next.tokType == endOfLineToken {
			return e
		}

		if e.next.tokType == paragraphEndToken {
			return e
		}
	}

	return e
}

func (t *token) nextTokenOfType(c tokenType) *token {

	if t.next == nil {
		return nil
	}

	for pos := t.next; pos != nil; pos = pos.next {

		if pos.tokType == c {
			return pos
		}
	}

	return nil

}

func (t *token) nextSignificantToken() *token {

	e := t

	if e.next == nil {
		return nil
	}

	for e = e.next; e != nil; e = e.next {

		if e.tokType != endOfLineToken {
			break
		}

	}

	return e
}

func (t *token) unicodeString() string {

	switch t.tokType {

	case gaijiCharToken:
		return t.unicodeContent

	case specialCharToken:
		return t.unicodeContent

	case kunojiToken:
		return t.unicodeContent

	default:
		return t.String()

	}
}

func (t *token) jis0213String() string {

	switch t.tokType {

	case gaijiCharToken:
		return t.jis0213Content

	case specialCharToken:
		return t.jis0213Content

	case kunojiToken:
		return t.jis0213Content

	default:
		return t.String()

	}
}
func (t *token) nextLine() *token {

	e := t.lastTokenInLine()

	if e.next == nil {
		return nil
	}
	if e.next.next == nil {
		return nil
	}

	return e.next.next

}

func matchingCloserStringIndex(s string, start, end string) int {

	count := 0

	for i := range s {

		if strings.HasPrefix(s[i:], start) {
			count++
			continue
		}

		if strings.HasPrefix(s[i:], end) {
			count--

			if count == 0 {
				return i
			}
		}
	}
	return -1

}

func (t *token) isFirstTokenInLine() bool {

	if t.prev == nil {
		return true
	}

	if t.prev.tokType == endOfLineToken {
		return true
	}

	if t.prev.tokType == paragraphEndToken {
		return true
	}

	return false

}

func (t *token) textContext() string {

	var s1 []rune

	maxlen := 10 + len(t.content)

	for e := t; len(s1) < maxlen/2+1; e = e.prev {

		if e == nil {
			break
		}

		s1 = append([]rune(e.content), s1...)

	}

	if len(s1) > maxlen/2+1 {

		s1 = s1[len(s1)-maxlen/2+1:]

	}

	for e := t.next; len(s1) < maxlen; e = e.next {

		if e == nil {
			break
		}

		s1 = append(s1, []rune(e.content)...)

	}

	var r []rune

	if len(s1) > maxlen {
		r = append(r, s1[:maxlen]...)
	} else {
		r = append(r, s1...)
	}

	return string(r)

}

func (t *token) mainTextStart() *token {

	if o_fragment {
		return t.firstToken()
	}

	for e := t.firstToken(); e != nil; e = e.next {
		if e.tokType == mainTextStartToken {
			return e
		}
	}

	return t.firstToken()

}

func (t *token) isFormatCloseToken() bool {

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

// t is assumed to be next of t2
func (t *token) switchWith(t2 *token) {

	e1 := t2.prev

	e2 := t.next

	e1.next = t

	t.prev = e1

	t.next = t2

	t2.next = e2

	t2.prev = t
}

func copyOf(t *token) *token {

	nt := new(token)

	nt.tokType = t.tokType

	nt.content = t.content

	nt.unicodeContent = t.unicodeContent

	nt.originalContent = t.originalContent

	nt.lineNo = t.lineNo

	return nt
}
