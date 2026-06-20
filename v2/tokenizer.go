package aozoraconvert

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type opt int

var oFull, oFragment, oJis0208, oJis0213, oRaw, oStrict, oVerbose, oTolerant, oParsable bool

type tokenizer struct {
	data             string
	position         int
	lineCounter      int
	simpleProcessing bool
	foundMarkupNote  bool
}

var tokenizerOption opt

func init() {

	oFull = true

	oRaw = false

	oJis0208 = false

	oJis0213 = false

	oStrict = false

	oTolerant = true

	oParsable = false

	log.SetFlags(0)

}

func setOutputOption(o string) {

	switch o {
	case "jis0208":
		oFull = false
		oJis0208 = true
		oRaw = false

	case "jis0213":
		oFull = false
		oJis0208 = false
		oJis0213 = true

	case "raw":
		oFull = false
		oJis0208 = true
		oRaw = true

	default:
		oFull = true
		oJis0208 = false
		oJis0213 = true
		oRaw = false
	}
	return
}

func tokenize(s string) *token {

	tknz := newTokenizer(s)

	t0 := tknz.nextToken()

	//for t1 := t0; !tknz.empty(); t1 = t1.next {
	for t1 := t0; t1 != nil; t1 = t1.next {

		t1.addTokenRight(tknz.nextToken())
	}
	return t0

}

func tokenizeAndFix(text string) (tokenString *token, err error) {
	if oRECOVER {
		defer func() {

			if r := recover(); r != nil {
				err = errors.New(r.(error).Error())
				return

			}
		}()
	}

	text = prep(text)
	tokenString = tokenize(text)
	if tokenString.lastToken().tokType != endOfLineToken {
		tokenString.lastToken().addTokenRight(newTokenOfType(endOfLineToken))
	}

	tokenString.lastToken().insertTokenRight(newTokenOfType(endOfLineToken))

	if !oRaw {

		tokenString.fixTextEndNote()

		tokenString.fixLines()

		if !oFragment {
			tokenString.insertAozoraBookMarker()
		}

		tokenString.fixEmptyText()

		tokenString.fixDocument()

		//return

		if oTolerant {
			tokenString.fixIndentationRound2()
		}

		tokenString.cleanUpIndentation()

		tokenString.fixSectionTitles()

		tokenString.insertSectionEnds()

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
	if oRECOVER {
		defer func() {

			if r := recover(); r != nil {

				panic(errors.New("line " + strconv.Itoa(tknz.lineCounter) + ":" + r.(string)))
				return

			}
		}()
	}

	e = newToken()

	end := 0

	e.tokType = typeOf(tknz)

	if e.tokType == eofToken {
		return nil
	}
	/*
		if tknz.simpleProcessing {
			if e.tokType != endOfLineToken && e.tokType != emptyLineToken {
				e.tokType = textToken
			}
		}
	*/
	switch e.tokType {

	case textToken:
		end = findContiguousText(tknz)

	case rubyStartToken:
		end = len(rubyStartStr)

	case rubyEndToken:
		end = len(rubyEndStr)

	case rubyParentStartToken:
		end = len(rubyParentStartStr)

	case emptyLineToken:
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(lineBreakStr)

	case endOfLineToken:
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		end = len(lineBreakStr)

	case gaijiToken:
		end = findMatchingCloser(gaijiToken, tknz)

	case noteStartToken:
		end = len(noteStartStr)

	case noteEndToken:
		end = len(noteEndStr)

	case bibInfoToken:
		tknz.lineCounter++
		e.lineNo = tknz.lineCounter
		tknz.simpleProcessing = true
		end = len(bibInfoStartStr)
		if strings.HasPrefix(tknz.data[tknz.position:], mainTextEndStr) {
			end = len(mainTextEndStr)
		}

	case accentStartToken:
		end = len(accentStartStr)

	case accentEndToken:
		end = len(accentEndStr)

	case kunojiToken:
		end = findMatchingCloser(kunojiToken, tknz)

	case markupNoteToken:
		end = findMatchingCloser(markupNoteToken, tknz)
		tknz.foundMarkupNote = true

	case specialCharToken:
		end = len(noteEndStr)

	}

	e.content = tknz.readNext(end)

	if e.tokType == markupNoteToken {

		tknz.lineCounter = tknz.lineCounter + len(strings.Split(e.content, "\n")) - 1

		//we discard explanation of aozorabunko-style markup!!
		if strings.HasPrefix(tknz.data[tknz.position:], "\n\n") {
			e.tokType = emptyToken
		} else {
			e.tokType = emptyLineToken
		}

	}

	if e.tokType == specialCharToken {

		if !oTolerant {

			if len(e.content) > len(noteEndStr) {

				panic(errors.New(e.content + "は外字中期に置き換えてください"))

			}
		}
	}

	//fmt.Print(e)

	return e
}

func typeOf(s *tokenizer) tokenType {

	if s.simpleProcessing {

		switch {
		case strings.HasPrefix(s.data[s.position:], lineBreakStr):
			if s.position == 0 || s.data[s.position-1] == '\n' {
				return emptyLineToken
			}
			return endOfLineToken

		case len(s.data[s.position:]) == 0:
			return eofToken

		default:
			return textToken
		}
	}

	switch {

	case isMarkupNoteStart(s):
		return markupNoteToken

	case strings.HasPrefix(s.data[s.position:], bibInfoStartStr):
		if s.position != 0 && s.data[s.position-1] == '\n' {
			return bibInfoToken
		}
		return textToken

	case strings.HasPrefix(s.data[s.position:], mainTextEndStr):
		return bibInfoToken

	case strings.HasPrefix(s.data[s.position:], rubyStartStr):
		return rubyStartToken

	case strings.HasPrefix(s.data[s.position:], rubyEndStr):
		return rubyEndToken

	case strings.HasPrefix(s.data[s.position:], rubyParentStartStr):
		return rubyParentStartToken

	case strings.HasPrefix(s.data[s.position:], lineBreakStr):
		if s.position == 0 || s.data[s.position-1] == '\n' {
			return emptyLineToken
		}
		return endOfLineToken

	case strings.HasPrefix(s.data[s.position:], noteStartStr):
		//	return noteToken
		return noteStartToken

	case strings.HasPrefix(s.data[s.position:], gaijiMarkerStr):
		return gaijiToken

	case strings.HasPrefix(s.data[s.position:], noteEndStr):
		return noteEndToken

	case strings.HasPrefix(s.data[s.position:], accentStartStr):
		return accentStartToken

	case strings.HasPrefix(s.data[s.position:], accentEndStr):
		return accentEndToken

	case strings.HasPrefix(s.data[s.position:], kunojiStr):
		return kunojiToken

	case strings.HasPrefix(s.data[s.position:], kunojiDakuStr):
		return kunojiToken

	case strings.HasPrefix(s.data[s.position:], squareBracketOpenStr):
		msg := "始め角括弧 => 外字注記"
		if !oTolerant {
			panic(errors.New(msg))
		}
		clog.Println(strconv.Itoa(s.lineCounter)+"行：", msg)
		return specialCharToken

	case strings.HasPrefix(s.data[s.position:], referenceMarkStr):
		msg := "※ => 外字注記"
		if !oTolerant {
			panic(errors.New(msg))
		}
		clog.Println(strconv.Itoa(s.lineCounter)+"行：", msg)
		return specialCharToken

	case strings.HasPrefix(s.data[s.position:], "＃"):
		msg := "＃ => 外字注記"
		if !oTolerant {
			panic(errors.New(msg))
		}
		clog.Println(strconv.Itoa(s.lineCounter)+"行：", msg)
		return specialCharToken

	case len(s.data[s.position:]) == 0:
		return eofToken

	default:
		return textToken

	}
	return emptyToken //this should never happen
}

func findContiguousText(s *tokenizer) (i int) {

	if s.simpleProcessing {
		i = strings.Index(s.data[s.position:], "\n")
		if i == -1 {
			i = len(s.data[s.position:])
		}
		return i

	}

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

		case strings.HasPrefix(s.data[s.position+i:], squareBracketOpenStr):
			if !oTolerant {
				panic(errors.New("始め角括弧は外字注記か違う字に置き換えてください"))
			}
			return i

		case strings.HasPrefix(s.data[s.position+i:], referenceMarkStr):
			if !oTolerant {
				panic(errors.New(referenceMarkStr + "は外字注記か違う字に置き換えてください"))
			}
			return i

		case strings.HasPrefix(s.data[s.position+i:], "＃"):
			if !oTolerant {
				panic(errors.New("＃は外字注記か違う字に置き換えてください"))
			}
			return i

		case strings.HasPrefix(s.data[s.position+i:], noteEndStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], accentStartStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], accentEndStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], kunojiStr):
			return i

		case strings.HasPrefix(s.data[s.position+i:], kunojiDakuStr):
			return i

			//case strings.HasPrefix(s.data[s.position+i:], bibInfoStartStr):
			//	return i
		}

	}

	return len(s.data[s.position:])

}

func findMatchingCloser(o tokenType, s *tokenizer) int {

	if o == markupNoteToken {

		lines := strings.Split(s.data[s.position:], "\n")

		if len(lines) < 3 {

			panic(errors.New("Tokenizer: arkup has no end"))

		}

		i := 2
		found := false

		for i = 2; i < len(lines); i++ {

			if hasOnly(lines[i], []rune(markupNoteDelimiterStr)[0]) {
				found = true
				break
			}
		}

		if !found {
			panic(errors.New("Tokenizer: markup has no end"))
		}

		j := strings.Index(s.data[s.position+len(lines[0])+len(lines[1])+2:], lines[i])

		if j == -1 {

			panic(errors.New("Tokenizer: markup has no end"))

		}

		return j + len(lines[0]) + len(lines[1]) + 2 + len(lines[i])

	}

	//	end := strings.Index(s.data[s.position:], lineBreakStr)

	if o != noteToken && o != gaijiToken {

		i := strings.Index(s.data[s.position:], closingStrOf(o))

		if i == -1 {

			panic(errors.New("Tokenizer: unmatched opening tag: " + strconv.Itoa(s.lineCounter) + " " + o.String() + "\n surrounding text: " + s.textContext()))

		}
		/*
			if i > end {
				panic(errors.New("Tokenizer: matching tag too far away: " + strconv.Itoa(s.lineCounter) + " " + o.String() + "\n surrounding text: " + s.textContext() + strconv.Itoa(end)))

			}
		*/
		j := strings.Index(s.data[s.position+1:], openingStrOf(o))

		if j > -1 && j < i-1 {

			i = findPairMatch(s.data[s.position:], openingStrOf(o), closingStrOf(o))

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

	panic(errors.New("Tokenizer: note not terminated. " + "\n surrounding text: " + s.textContext()))

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

	case accentStartToken:

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

	case kunojiToken:

		return kunojiStr

	case markupNoteToken:

		return markupNoteStartStr

	default:

		return emptyStr

	}

}

func (tknz *tokenizer) textContext() (tctx string) {

	maxlen := 20

	r1 := []rune(tknz.data[:tknz.position])

	r2 := []rune(tknz.data[tknz.position:])

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

func findPairMatch(s string, opener, closer string) int {

	var i, c int

	for i = range s {

		if strings.HasPrefix(s[i:], opener) {
			c++
			continue
		}

		if strings.HasPrefix(s[i:], closer) {
			c--
			if c == 0 {
				return i
			}

		}
	}

	return -1
}

func isMarkupNoteStart(s *tokenizer) bool {

	if s.foundMarkupNote {
		return false
	}

	if s.lineCounter > 20 {
		s.foundMarkupNote = true
		return false
	}

	if s.position != 0 && s.data[s.position-1] != '\n' {
		return false
	}

	lines := strings.Split(s.data[s.position:], "\n")

	if len(lines) < 2 {
		return false
	}

	if !hasOnly(lines[0], '-') {
		return false
	}

	return true
	//	return strings.Contains(lines[1], markupNoteStartStr)

}

func hasOnly(s string, c rune) bool {

	if len(s) == 0 {
		return false
	}

	for _, e := range s {

		if e != c {
			return false
		}
	}

	return true
}

func prep(in string) string {
	var newLines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		tline := strings.Trim(line, "\u0020\u3000")
		if len(tline) == 0 {
			newLines = append(newLines, tline)
		} else {
			newLines = append(newLines, line)
		}
	}
	return strings.Join(newLines, "\n")
}
