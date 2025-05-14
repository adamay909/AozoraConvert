package aozoraconvert

//go:generate stringer -type=tokenType
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
	accentProcessed bool
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
	accentStartToken
	accentEndToken
	gaijiCharToken
	kunojiToken
	alignBottomCloserToken
	centeringEndToken
	markupNoteToken
	mainTextStartToken
	mainTextEndToken
	eofToken
	gaijiImgToken
	noteStartToken
	noteEndToken
	gaijiStartToken
	endMarkerToken
)

type tokenSubType int

/*
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

		case accentStartToken:
			return "accentStartToken"

		case accentEndToken:
			return "accentEndToken"

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
*/

func newToken() *token {

	t := new(token)

	t.inserted = false

	return t
}

func (t *token) addTokenRight(t2 *token) {

	if t2 == nil {
		return
	}

	t.insertTokenRight(t2)

	t2.inserted = false

	return

}

func (t *token) insertTokenRight(t2 *token) {

	if t.next != nil {

		t.next.prev = t2

	}

	t2.next = t.next

	t.next = t2

	t2.prev = t

	t2.inserted = true

}

func (t *token) insertTokenLeft(t2 *token) {

	if t.prev != nil {

		t.prev.next = t2

	}

	t2.prev = t.prev

	t.prev = t2

	t2.next = t

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

	default:

		return t.String()

	}
}

func (t *token) innerUnicodeString() string {

	if t == nil {
		return t.unicodeString()
	}

	switch t.tokType {

	case noteToken, gaijiImgToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.unicodeString(), noteStartStr), noteEndStr)

	case gaijiToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.unicodeString(), gaijiMarkerStr), noteEndStr)

	case gaijiNoteToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.unicodeString(), noteStartStr), noteEndStr)

	default:

		return t.unicodeString()

	}
}

func (t *token) innerJis0213String() string {

	if t == nil {
		return t.jis0213String()
	}

	switch t.tokType {

	case noteToken, gaijiImgToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.jis0213String(), noteStartStr), noteEndStr)

	case gaijiToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.jis0213String(), gaijiMarkerStr), noteEndStr)

	case gaijiNoteToken:

		return strings.TrimSuffix(strings.TrimPrefix(t.jis0213String(), noteStartStr), noteEndStr)

	default:

		return t.jis0213String()

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

	return t.getRefStrings()[0]

}

func (t *token) getRefStrings() []string {

	var rs []string

	s := strings.TrimSuffix(strings.TrimPrefix(t.unicodeContent, noteStartStr), noteEndStr)

	if s == "" {
		s = t.innerString()
	}

	for i1 := strings.Index(s, refStartStr); i1 != -1; i1 = strings.Index(s, refStartStr) {

		i2 := matchingCloserStringIndex(s, refStartStr, refEndStr)

		if i2 == -1 {
			break
		}

		ostr := s[i1+len(refStartStr) : i2]

		if oTolerant {

			ignore := false
			newstr := ""

			for k, r := range ostr {

				if strings.HasPrefix(ostr[k:], rubyStartStr) {
					ignore = true
					continue
				}

				if strings.HasPrefix(ostr[k:], gaijiMarkerStr) {
					ignore = true
					newstr = newstr + referenceMarkStr
					continue
				}

				if !ignore {
					newstr = newstr + string(r)
				}

				if strings.HasPrefix(ostr[k:], rubyEndStr) {
					ignore = false
					continue
				}

				if strings.HasPrefix(ostr[k:], noteEndStr) {
					ignore = false
					continue
				}

			}
			if newstr != ostr {

				ostr = newstr

				clog.Println("line ", t.lineNumber(), t.String(), "前方参照文字列の指定のルビを削除")
			}
		}

		rs = append(rs, ostr)

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

	ignore := false

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == endOfLineToken {
			break
		}

		if e.tokType == emptyLineToken {
			break
		}

		if e.isPairClose() {
			nesting++
		}

		if e.isPairOpen() {
			nesting--
		}

		if e.tokType == rubyEndToken {
			ignore = true
			continue
		}

		if e.tokType == rubyStartToken {
			ignore = false
			continue
		}

		if ignore {
			continue
		}

		if !e.isText() {
			continue
		}

		cr = c

		c = len(e.unicodeString()) + cr

		if c >= len(txt) {

			break

		}

	}

	if e == nil {

		panic("Can't find place to insert implied opener note. Defaulting to start of line. " + e.info())

		return

	}

	if e.tokType == endOfLineToken || e.tokType == emptyLineToken {

		if !oTolerant {
			panic("Can't find place to insert implied opener note. " + t.info())
		}

		clog.Println(t.lineNumberStr() + "行：前方参照の文字列が見つからず。行頭まで参照とみなす。")
		e.insertTokenRight(nt)
		return

	}

	//	fmt.Println(e, nesting)

	if nesting > 0 {

		for ; nesting > 0; e = e.prev {
			if e.isPairOpen() {
				nesting--
			}
		}

		//e=e.prev is executed once after nesting has gone down to zero
		e.insertTokenRight(nt)
		return
	}

	if c == len(txt) {

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

func (t *token) lineNumberStr() string {

	return strconv.Itoa(t.lineNumber())

}

func (t *token) lastTokenInLine() *token {

	if t == nil {
		return t
	}

	e := new(token)

	for e = t; e != nil; e = e.next {

		if e.tokType == endOfLineToken {
			return e
		}

		if e.tokType == paragraphEndToken {
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

	if t.unicodeContent == "" {
		return t.content
	}

	return t.unicodeContent

}

func (t *token) jis0213String() string {

	if t.jis0213Content == "" {
		return t.content

	}

	return t.jis0213Content

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

	if t.prev.tokType == emptyLineToken {
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

	if oFragment {
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

func (t *token) isPairOpen() bool {

	switch t.tokType {

	case paragraphToken, rubyStartToken, rubyParentStartToken, rubyGroupStartToken, figureStartToken, mainTextStartToken:
		return true

	}

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

	switch t.tokType {

	case paragraphEndToken, rubyEndToken, rubyParentEndToken, rubyGroupEndToken, figureEndToken, centeringEndToken, alignBottomCloserToken, mainTextEndToken:
		return true

	}

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

	if t == nil {
		return false
	}

	if t.tokType != noteToken {

		return false

	}

	str := strings.TrimPrefix(t.innerString(), blockEndStr)

	if !strings.HasSuffix(str, formatEndStr) {
		return false
	}

	str = strings.TrimSuffix(str, formatEndStr)

	for _, m := range sectionMarker {

		if str == m {
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

func (t *token) isIndentationEnd() bool {

	if t.tokType != noteToken {
		return false
	}

	if !strings.HasPrefix(t.innerString(), blockEndStr) {
		return false
	}

	return strings.HasSuffix(strings.TrimSuffix(strings.TrimPrefix(t.innerString(), blockEndStr), formatEndStr), "字下げ")
}

func (t *token) matchingIndentationCloser() *token {

	pos := new(token)

	c := 0

	for pos = t.next; pos != nil; pos = pos.next {

		//need to ensure proper nesting
		if pos.isBlockStartNote() {
			c++
		}

		if pos.isBlockClosingNote() {
			c--
			if c < 0 {
				return pos
			}
		}

		if pos.isFrameStart() {

			for pos2 := pos; pos2 != pos.matchingCloserToken(); pos2 = pos2.next {
				if pos2.isIndentationStart() {
					pos2.fixblockformatting()
				}
			}
			pos = pos.matchingCloserToken()
			continue
		}

		if pos.isNarrowStart() {
			for pos2 := pos; pos2 != pos.matchingCloserToken(); pos2 = pos2.next {
				if pos2.isIndentationStart() {
					pos2.fixblockformatting()
				}
			}
			pos = pos.matchingCloserToken()
			continue
		}

		if pos.isIndentationStart() {
			return pos
		}

		if pos.isIndentationEnd() {
			return pos
		}

		if pos.tokType == mainTextEndToken {
			return pos.prev
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

	if len(s) == 0 {
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

		if pos.isPairOpen() {
			count++
		}

		if pos.isPairClose() {
			count--
		}

		if pos.tokType == paragraphToken {
			count++
		}

		if pos.tokType == paragraphEndToken {
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

	if t.tokType == centeringEndToken {
		return true
	}

	if t.tokType != noteToken {
		return false
	}

	if t.isSectionTitleEnd() {
		return true
	}

	if !strings.HasPrefix(t.innerString(), blockEndStr) {
		return false
	}

	if strings.HasSuffix(t.innerString(), formatEndStr) {
		return true
	}

	return false
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

	nt.jis0213Content = t.jis0213Content

	nt.originalContent = t.originalContent

	return nt
}

func (t *token) isEmptyText() bool {

	if t.tokType != textToken {
		return false
	}
	for _, c := range t.String() {

		if c == ' ' {
			continue
		}

		if c == '　' {
			continue
		}

		return false
	}

	return true
}

func (t *token) isText() bool {

	if t.tokType == textToken {
		return true
	}

	if t.tokType == kunojiToken {
		return true
	}

	if t.tokType == accentStartToken {
		return true
	}

	if t.tokType == accentEndToken {
		return true
	}

	if t.tokType == gaijiCharToken {
		return true
	}

	if t.tokType == specialCharToken {
		return true
	}

	return false
}

func (t *token) isFrameStart() bool {

	if strings.HasSuffix(t.innerString(), "罫囲み") {
		return true
	}

	return false
}

func (t *token) isNarrowStart() bool {

	if strings.HasSuffix(t.innerString(), "字詰め") {
		return true
	}

	return false
}

func (t *token) isLineBreak() bool {

	if t.tokType == endOfLineToken {
		return true
	}

	if t.tokType == emptyLineToken {
		return true
	}

	return false
}

func (t *token) matchingNoteStart() *token {

	if t.tokType != noteEndToken {
		return t
	}

	counter := 0

	for e := t; e != nil; e = e.prev {

		if e.tokType == noteEndToken {
			counter++
			continue
		}

		if e.tokType != noteStartToken {
			continue
		}

		counter--

		if counter == 0 {
			return e
		}
	}

	return nil
}

func (t *token) matchingNoteEnd() *token {

	if t.tokType != noteStartToken {
		return t
	}

	counter := 0

	for e := t; e != nil; e = e.next {

		if e.tokType == noteStartToken {
			counter++
			continue
		}

		if e.tokType != noteEndToken {
			continue
		}

		counter--

		if counter == 0 {
			return e
		}
	}

	return nil
}

func (t *token) specialCharEscape() string {

	switch t.content {

	case "《":
		return "※［＃始め二重山括弧、1-1-52］"

	case "》":
		return "※［＃終わり二重山括弧、1-1-53］"

	case "［":
		return "※［＃始め角括弧、1-1-46］"

	case "］":
		return "※［＃終わり角括弧、1-1-47］"

	case "〔":
		return "※［＃始めきっこう（亀甲）括弧、1-1-44］"

	case "〕":
		return "※［＃終わりきっこう（亀甲）括弧、1-1-45］"

	case "｜":
		return "※［＃縦線、1-1-35］"

	case "＃":
		return "※［＃井げた、1-1-84］"

	case "※":
		return "※［＃米印、1-2-8］"

	default:
		return t.content
	}
}
