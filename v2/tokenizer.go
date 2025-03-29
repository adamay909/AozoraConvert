package aozoratext

import (
	"log"
	"os"
	"strings"
)

type opt int

var o_full, o_fragment, o_jis0208, o_raw bool

type tokenizer struct {
	data        string
	position    int
	lineCounter int
}

var tokenizerOption opt

var tokenizerLog *log.Logger

var dummyLog *strings.Builder

func init() {

	o_full = true

	o_raw = false

	o_jis0208 = false

	tokenizerLog = log.New(os.Stdout, "", 0)

	log.SetFlags(0)

}

func setOutputOption(o string) {

	switch o {
	case "jis0208":
		o_full = false
		o_jis0208 = true
		o_raw = false

	case "raw":
		o_full = false
		o_jis0208 = true
		o_raw = true

	default:
		o_full = true
		o_jis0208 = false
		o_raw = false
	}
	return
}

// offset is setting line number in case s isn't the whole document being processed.
func tokenize(s string) *token {

	tknz := newTokenizer(s)

	t0 := tknz.nextToken()

	//for t1 := t0; !tknz.empty(); t1 = t1.next {
	for t1 := t0; t1 != nil; t1 = t1.next {

		t1.addTokenRight(tknz.nextToken())

	}

	return t0

}

func tokenizeAndFix(text string) (tk *token, err error) {
	/*
		defer func() {

			if r := recover(); r != nil {

				log.Println(r)

				err = errors.New("error")
				return

			}
		}()
	*/
	tokenString := tokenize(text)

	tokenString.lastToken().insertTokenRight(newTokenOfType(endOfLineToken))

	if !o_raw {

		tokenString.fixLines()

		tokenString.fixEmptyText()

		tokenString.fixDocument()

		tokenString.insertSectionEnds()

		if !o_fragment {
			tokenString.insertAozoraBookMarker()
		}

	}

	return tokenString, err

}

func newTokenizer(d string) *tokenizer {

	r := new(tokenizer)

	r.data = d

	r.position = 0

	r.lineCounter = 1

	return r
}

// read next i bytes
func (tknz *tokenizer) readNext(i int) string {

	tknz.position = tknz.position + i

	return tknz.data[tknz.position-i : tknz.position]

}

func (tknz *tokenizer) nextToken() (e *token) {

	e = newToken()

	end := 0

	e.tokType = typeOf(tknz)

	if e.tokType == eofToken {
		return nil
	}

	switch e.tokType {

	case textToken:
		end = findContiguousText(tknz)

	case rubyStartToken:
		end = len(rubyStartStr)

	case rubyEndToken:
		end = len(rubyEndStr)

	case rubyParentStartToken:
		end = len(rubyParentStartStr)

	case endOfLineToken:
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(lineBreakStr)

	case gaijiToken:
		end = findMatchingCloser(gaijiToken, tknz)

	case noteToken:
		end = findMatchingCloser(noteToken, tknz)

	case bibInfoToken:
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(bibInfoStartStr)

	case accentToken:
		end = findMatchingCloser(accentToken, tknz)

	case kunojiToken:
		end = findMatchingCloser(kunojiToken, tknz)

	case markupNoteToken:
		end = findMatchingCloser(markupNoteToken, tknz)

	}
	e.content = tknz.readNext(end)

	if e.tokType == noteToken {

		if e.innerString() == mainTextEndStr {

			e.tokType = bibInfoToken

		}
	}

	if e.tokType == markupNoteToken {

		tknz.lineCounter = tknz.lineCounter + len(strings.Split(e.content, "\n")) - 1

		//we discard explanation of aozorabunko-style markup!!
		return tknz.nextToken()

	}

	return e
}

func typeOf(s *tokenizer) tokenType {

	switch {

	case strings.HasPrefix(s.data[s.position:], markupNoteStartStr):
		return markupNoteToken

	case strings.HasPrefix(s.data[s.position:], bibInfoStartStr):
		//must come before detection of linebreak
		return bibInfoToken

	case strings.HasPrefix(s.data[s.position:], rubyStartStr):
		return rubyStartToken

	case strings.HasPrefix(s.data[s.position:], rubyEndStr):
		return rubyEndToken

	case strings.HasPrefix(s.data[s.position:], rubyParentStartStr):
		return rubyParentStartToken

	case strings.HasPrefix(s.data[s.position:], lineBreakStr):
		return endOfLineToken

	case strings.HasPrefix(s.data[s.position:], noteStartStr):
		return noteToken

	case strings.HasPrefix(s.data[s.position:], gaijiMarkerStr):
		return gaijiToken

	case strings.HasPrefix(s.data[s.position:], accentStartStr):
		return accentToken

	case strings.HasPrefix(s.data[s.position:], kunojiStr):
		return kunojiToken

	case strings.HasPrefix(s.data[s.position:], kunojiDakuStr):
		return kunojiToken

	case len(s.data[s.position:]) == 0:
		return eofToken

	default:
		return textToken

	}
}

func findContiguousText(s *tokenizer) (i int) {

	for i = range s.data[s.position:] {

		switch {

		case strings.HasPrefix(s.data[s.position+i:], rubyStartStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], rubyEndStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], rubyParentStartStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], lineBreakStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], gaijiMarkerStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], noteStartStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], noteEndStr):
			panic("found unexpected: " + noteEndStr)
			return i

		case strings.HasPrefix(s.data[s.position+i:], accentStartStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], accentEndStr):
			panic("found unexpected: " + accentEndStr)
			return i

		case strings.HasPrefix(s.data[s.position+i:], kunojiStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], kunojiDakuStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], bibInfoStartStr):
			return i
		}

	}

	return len(s.data[s.position:])

}

func findMatchingCloser(o tokenType, s *tokenizer) int {

	if o == markupNoteToken {

		i := strings.Index(s.data[s.position+len(openingStrOf(o)):], markupNoteEndStr)

		if i == -1 {

			panic("markup notes do not end")

		}

		return i + len(openingStrOf(o)) + len(closingStrOf(o))

	}

	end := strings.Index(s.data[s.position:], lineBreakStr)

	if o != noteToken && o != gaijiToken {

		i := strings.Index(s.data[s.position:], closingStrOf(o))

		if i == -1 {

			panic("unmatched opening tag: " + o.String() + "\n surrounding text: " + s.textContext())

		}

		if i > end {
			panic("unmatched opening tag: " + o.String() + "\n surrounding text: " + s.textContext())

		}

		return i + len(closingStrOf(o))

	}

	counter := 0

	for i := range s.data[s.position:] {

		if strings.HasPrefix(s.data[s.position+i:], noteStartStr) {
			counter++
			continue
		}

		if strings.HasPrefix(s.data[s.position+i:], noteEndStr) {
			counter--
			if counter == 0 {
				return i + len(noteEndStr)
			}
		}

		if strings.HasPrefix(s.data[s.position+i:], lineBreakStr) {
			break
		}

	}

	panic("note not terminated. " + "\n surrounding text: " + s.textContext())

	return -1

}

func closingStrOf(o tokenType) string {

	switch o {

	case rubyStartToken:

		return rubyEndStr

	case gaijiToken:

		return noteEndStr

	case noteToken:

		return noteEndStr

	case accentToken:

		return accentEndStr

	case kunojiToken:

		return kunojiEndStr

	case markupNoteToken:

		return markupNoteEndStr

	default:

		return emptyStr

	}

}

func openingStrOf(o tokenType) string {

	switch o {

	case rubyStartToken:

		return rubyStartStr

	case gaijiToken:

		return gaijiMarkerStr

	case noteToken:

		return noteStartStr

	case accentToken:

		return accentStartStr

	case kunojiToken:

		return kunojiStr

	case markupNoteToken:

		return markupNoteStartStr

	default:

		return emptyStr

	}

}

func (t *tokenizer) textContext() (tctx string) {

	maxlen := 20

	r1 := []rune(t.data[:t.position])

	r2 := []rune(t.data[t.position:])

	if len(r1) > maxlen/2+1 {

		tctx = string(r1[len(r1)-maxlen/2+1:])

	} else {
		tctx = string(r1)

	}

	if len(r2) > maxlen/2+1 {

		tctx = tctx + string(r2[:maxlen/2+1])

	} else {

		tctx = tctx + string(r2)

	}

	return

}
