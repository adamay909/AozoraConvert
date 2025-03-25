package aozoratext

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type opt int

var o_full, o_raw, o_jis0208 bool

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
func tokenize(s string, offset int) *token {

	tknz := newTokenizer(s, offset)

	t0 := tknz.nextToken()

	for t1 := t0; !tknz.empty(); t1 = t1.next {

		t1.addTokenRight(tknz.nextToken())

	}

	return t0

}

func tokenizeAndFix(text string, offset int) *token {

	tokenString := tokenize(strings.Join(strings.Split(text, "\n")[offset:], "\n"), offset)

	tokenString.lastToken().insertTokenRight(newTokenOfType(endOfLineToken))

	if !o_raw {

		tokenString.fixLines()

		tokenString.fixEmptyText()

		tokenString.fixDocument()

		tokenString.insertSectionEnds()

	}

	return tokenString

}

func newTokenizer(d string, offset int) *tokenizer {

	r := new(tokenizer)

	r.data = d

	r.position = 0

	r.lineCounter = offset

	if !strings.HasPrefix(d, "\n") {

		r.lineCounter++

	}

	return r
}

func (tknz *tokenizer) empty() bool {

	return tknz.position == len(tknz.data)

}

func (tknz *tokenizer) remainder() string {

	return tknz.data[tknz.position:]

}

// read next i bytes
func (tknz *tokenizer) readNext(i int) string {

	tknz.position = tknz.position + i

	return tknz.data[tknz.position-i : tknz.position]

}

func (tknz *tokenizer) nextToken() (e *token) {

	defer func() {

		r := recover()

		if r != nil {

			fmt.Println("The document has errors.")

			fmt.Println("line ", tknz.lineCounter, ": ", r)

		}

		return

	}()

	e = newToken()

	end := 0

	e.tokType = typeOf(tknz.remainder())

	switch e.tokType {

	case textToken:
		end = findContiguousText(tknz.remainder())

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
		end = findMatchingCloser(gaijiToken, tknz.remainder())

	case noteToken:
		end = findMatchingCloser(noteToken, tknz.remainder())

	case bibInfoToken:
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(bibInfoStartStr)

	case accentToken:
		end = findMatchingCloser(accentToken, tknz.remainder())

	case kunojiToken:
		end = findMatchingCloser(kunojiToken, tknz.remainder())

	}
	e.content = tknz.readNext(end)

	if e.tokType == noteToken {

		if e.innerString() == mainTextEndStr {

			e.tokType = bibInfoToken

		}
	}

	return e
}

func typeOf(s string) tokenType {

	switch {

	case len(s) == 0:
		return emptyLineToken

	case strings.HasPrefix(s, bibInfoStartStr):
		//must come before detection of linebreak
		return bibInfoToken

	case strings.HasPrefix(s, rubyStartStr):
		return rubyStartToken

	case strings.HasPrefix(s, rubyEndStr):
		return rubyEndToken

	case strings.HasPrefix(s, rubyParentStartStr):
		return rubyParentStartToken

	case strings.HasPrefix(s, lineBreakStr):
		return endOfLineToken

	case strings.HasPrefix(s, noteStartStr):
		return noteToken

	case strings.HasPrefix(s, gaijiMarkerStr):
		return gaijiToken

	case strings.HasPrefix(s, accentStartStr):
		return accentToken

	case strings.HasPrefix(s, kunojiStr):
		return kunojiToken

	case strings.HasPrefix(s, kunojiDakuStr):
		return kunojiToken

	default:
		return textToken

	}
}

func findContiguousText(s string) (i int) {

	for i = range s {

		switch {

		case strings.HasPrefix(s[i:], rubyStartStr):
			return i

		case strings.HasPrefix(s[i:], rubyEndStr):
			return i

		case strings.HasPrefix(s[i:], rubyParentStartStr):
			return i

		case strings.HasPrefix(s[i:], lineBreakStr):
			return i

		case strings.HasPrefix(s[i:], gaijiMarkerStr):
			return i

		case strings.HasPrefix(s[i:], noteStartStr):
			return i

		case strings.HasPrefix(s[i:], noteEndStr):
			panic("found unexpected: " + noteEndStr)
			return i

		case strings.HasPrefix(s[i:], accentStartStr):
			return i

		case strings.HasPrefix(s[i:], accentEndStr):
			panic("found unexpected: " + accentEndStr)
			return i

		case strings.HasPrefix(s[i:], kunojiStr):
			return i

		case strings.HasPrefix(s[i:], kunojiDakuStr):
			return i

		case strings.HasPrefix(s[i:], bibInfoStartStr):
			return i
		}

	}

	return len(s)

}

func findMatchingCloser(o tokenType, s string) (i int) {

	counter := 0

	i = 0

	for i = range s {

		if typeOf(s[i:]) == o {
			counter++
			continue
		}

		if o == noteToken {

			if typeOf(s[i:]) == gaijiToken {
				counter++
				continue
			}
		}

		if strings.HasPrefix(s[i:], closingStrOf(o)) {

			counter--

			if counter == 0 {
				break
			}
		}

		if typeOf(s[i:]) == endOfLineToken {

			break
		}

	}

	if o == rubyParentStartToken {
		return i
	}

	if counter != 0 {

		panic("Found unclosed opening tag: " + o.String())

	}

	return i + len(closingStrOf(o))
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

	default:

		return emptyStr

	}

}
