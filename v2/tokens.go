package aozoratext

import (
	"fmt"
	"strconv"
	"strings"
)

type token struct {
	tokType    tokenType
	next       *token
	prev       *token
	ref        string
	content    string
	altContent string
	id         int
	lineNo     int
	// extra   string
}

type tokenType int8

// These define types of tokens
const (
	emptyToken tokenType = iota
	textToken
	gaijiToken
	rubyParentStartToken
	rubyParentEndToken
	rubyStartToken
	rubyEndToken
	noteToken
	bibInfoToken
	emptyLineToken
	endOfLineToken
	lineBreakToken
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

	case emptyLineToken:
		return "emptyLineToken"

	case endOfLineToken:
		return "endOfLineToken"

	case lineBreakToken:
		return "lineBreakToken"

	case paragraphToken:
		return "paragraphToken"

	case sectionTitleStartToken:
		return "sectionTitleStart"

	case sectionTitleEndToken:
		return "sectionTitleEnd"

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
		return "subsubsectionEndtoken"

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

	default:

		return strconv.Itoa(int(t))

	}
}

func newToken() *token {

	return new(token)

}

func (t *token) insertTokenRight(t2 *token) {

	t2.next = t.next

	if t2.next != nil {
		t2.next.prev = t2
	}

	t2.prev = t

	t.next = t2

}

func (t *token) insertTokenLeft(t2 *token) {

	if t.prev != nil {

		t.prev.insertTokenRight(t2)

		return

	}

	t.prev = t2

	t2.next = t
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

	if t.prev != nil {
		t.prev.next = t.next
	}

	if t.next != nil {
		t.next.prev = t.prev
	}

	t = nil

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

func closingTagOf(o atomType) atomType {

	switch o {

	case rubyStartTag:

		return rubyEndTag

		//	case rubyBaseStartTag:

		//		return rubyStartTag

	case gaijiMarker:

		return noteEndTag

	case noteStartTag:

		return noteEndTag

	case accentStartTag:

		return accentEndTag

	default:

		return emptyAtom

	}

	return emptyAtom
}

func (t *token) String() string {

	if t == nil {

		return "<NIL>"

	}

	if t.content != "" {

		return t.content
	}

	switch t.tokType {

	case emptyToken:
		return emptyStr

	case emptyLineToken:
		return ""

	case endOfLineToken:
		return ""

	case lineBreakToken:
		return ""

	case rubyParentStartToken:
		return rubyBaseStartStr

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

	case noteToken:

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

		i2 := strings.Index(s, refEndStr)

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

	t.tokType = lineBreakToken

	return t
}

func newParagraphToken() *token {

	t := newToken()

	t.tokType = paragraphToken

	return t

}

func (t *token) listAllTokens() string {

	output := new(strings.Builder)

	counter := 0

	for e := t.firstToken(); e != nil; e = e.next {
		counter++
		addToStringsBuilder(output, strconv.Itoa(counter), ": ", e.tokType.String(), ": ", e.String(), " \n")

	}

	return output.String()

}

func (t *token) addTokenBefore(txt string, nt *token) {

	e := t.prev

	c := 0

	cr := 0

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == rubyEndToken {

			for ; e.tokType != rubyStartToken; e = e.prev {
			}

			e = e.prev
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

		tokenizerLog.Println("Can't find place to insert implied opener note. Defaulting to start of line.")

		e.insertTokenLeft(nt)

		return

	}

	if c == len(txt) {

		if e.prev != nil && e.prev.tokType == rubyParentStartToken {

			e = e.prev

		}

		e.insertTokenLeft(nt)

		return

	}

	k := len(txt) - cr

	t2 := newTokenOfType(textToken)

	//	s1 := e.content[:len(e.content)-k]

	t2.setString(e.content[:len(e.content)-k])

	e.setString(e.content[len(e.content)-k:])

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

	return 0

}
