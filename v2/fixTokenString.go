package aozoraConvert

import (
	"log"
	"strconv"
	"strings"
)

func (t *token) fixLines() {

	for e := t.firstToken(); e != nil; e = e.next {

		e.fixkunoji()
		e.fixaccent()
		e.fixgaiji()

		e.fixruby()

		e.fixBibInfo()

		e.fixAbbreviatedBlockFormat()
		e.fixBlockFormat()
		e.fixInlineAlignBottom()
		e.fixImpliedOpener()
	}

	t.firstToken().fixParagraphs()

}

func (t *token) fixAbbreviatedBlockFormat() {

	if !t.isFirstTokenInLine() {
		return
	}

	if t.tokType != noteToken {
		return
	}

	if strings.HasPrefix(t.innerString(), blockStartStr) {
		return
	}

	m := emptyStr

	for _, e := range impliedCloserMarker {

		if strings.HasSuffix(t.innerString(), e) {

			m = e

			break
		}
	}

	if m == emptyStr {
		return
	}

	//now we know we are dealing with a relevant token

	//first fix t itself.
	sbs := new(strings.Builder)

	addToStringsBuilder(sbs, blockStartStr, t.innerString())

	t.originalContent = t.content

	t.setInnerString(sbs.String())

	t.modified = true

	//need to check if we have to insert a new closing token.

	sbe := new(strings.Builder)

	addToStringsBuilder(sbe, blockEndStr, m, formatEndStr)

	if t.nextLine() != nil {

		if t.nextLine().tokType == noteToken && t.nextLine().innerString() == sbe.String() {

			//need to insert new eol
			t.insertTokenRight(newTokenOfType(endOfLineToken))

			return

		}
	}

	//we need to insert a new closing token.

	t.lastTokenInLine().insertTokenRight(newNote(sbe.String()))

	t.lastTokenInLine().insertTokenLeft(newTokenOfType(endOfLineToken))

	//need to insert new eol
	t.insertTokenRight(newTokenOfType(endOfLineToken))

	return

}

func (t *token) fixInlineAlignBottom() {

	if t.isFirstTokenInLine() {

		return

	}

	if t.tokType != noteToken {

		return

	}

	m := emptyStr

	for _, e := range impliedCloserMarker {

		if strings.HasSuffix(t.innerString(), e) {

			m = e

			break
		}
	}

	if m == emptyStr {
		return
	}

	//now we know we are dealing with a relevant token.

	t.lastTokenInLine().insertTokenRight(newTokenOfType(alignBottomCloserToken))

	return

}

func (t *token) fixkunoji() {

	if t.tokType != kunojiToken {
		return
	}

	if t.content == kunojiStr {

		t.unicodeContent = kunojiStrU

		t.jis0213Content = kunojiStrU

		t.modified = true

		return
	}

	t.unicodeContent = kunojiDakuStrU

	t.jis0213Content = kunojiDakuStrU

	t.modified = true
}

func (t *token) fixruby() {

	if t.tokType != rubyStartToken {

		return
	}

	t.insertTokenLeft(newTokenOfType(rubyParentEndToken))

	e := new(token)

	found := false

	for e = t.next; e != nil; e = e.next {

		if e.tokType == rubyEndToken {

			e.insertTokenRight(newTokenOfType(rubyGroupEndToken))

			found = true

			break
		}

		if e.tokType == endOfLineToken {
			break
		}

		if e.tokType == rubyStartToken {
			break
		}

		if e.tokType == rubyParentStartToken {
			break
		}

	}
	if !found {

		panic("unclosed ruby tag: " + t.info() + "\n surrounding text: " + t.textContext())

	}

	if e2 := t.prev; e2 != nil && e2.tokType == noteToken {

		if e2.isImage() {

			msg := "line " + strconv.Itoa(e2.lineNumber()) + " grphics used as characer? Attempting fix."

			e2.tokType = gaijiImgToken

			log.Println("WARNING: " + msg)

		} else {

			msg := "line " + strconv.Itoa(e2.lineNumber()) + " wrong order of ruby and annotation:" + e2.prev.String() + e2.String() + e2.next.String() + e2.next.next.String()

			if o_strict {
				panic("ERROR: " + msg)
			}

			log.Println("WARNING: " + msg)

			e2.remove()

			e.next.insertTokenRight(e2)
		}
	}

	if !t.rubyParentExplicit() {

		t.insertRubyParentStart()

	}

	return

}

func (t *token) rubyParentExplicit() bool {

	for e := t.prev; e != nil; e = e.prev {

		if e.tokType == rubyParentEndToken {
			continue
		}
		if e.tokType == textToken {
			continue
		}

		if e.tokType == rubyParentStartToken {
			e.insertTokenLeft(newTokenOfType(rubyGroupStartToken))
			return true
		}

		if e.tokType == gaijiCharToken {
			continue
		}

		if e.tokType == kunojiToken {
			continue
		}

		if e.tokType == gaijiNoteToken {
			continue
		}

		if e.tokType == accentToken {
			continue
		}

		if e.isKunten() {
			continue
		}

		if e.isOkurigana() {
			continue
		}

		if e.isImage() {
			continue
		}

		if e.tokType == emptyToken {
			continue
		}

		break
	}

	return false
}

func (t *token) insertRubyParentStart() {

	e := new(token)

	e = t.prev

	if e.prev.tokType == accentToken {

		e.prev.insertTokenLeft(newTokenOfType(rubyParentStartToken))

		e.prev.prev.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

		return

	}

	if e.prev.isImage() {

		e.prev.insertTokenLeft(newTokenOfType(rubyParentStartToken))

		e.prev.prev.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

		return

	}
	var ref charTypeID

	var r []rune

	for e = t.prev; e != nil; e = e.prev {

		/*if e.tokType == rubyEndToken {

			for ; e.tokType != rubyStartToken; e = e.prev {
			}

			e = e.prev

		}
		*/
		if e.tokType == textToken {
			break
		}

		if e.tokType == gaijiCharToken {
			break
		}

	}

	r = []rune(e.unicodeString())

	ref = CharType(r[len(r)-1])

	k := 0

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == rubyParentEndToken {
			continue
		}

		if e.tokType == rubyGroupEndToken {
			break
		}

		if e.tokType == endOfLineToken {
			break

		}

		if e.tokType == kunojiToken {
			break
		}

		if e.tokType == noteToken {
			break
		}

		if e.tokType != textToken && e.tokType != gaijiCharToken {
			continue
		}

		r = []rune(e.unicodeString())

		for k = len(r) - 1; ref == CharType(r[k]); k-- {

			if k == 0 {

				break
			}
		}

		if CharType(r[k]) != ref {

			t2 := newTokenOfType(textToken)

			t2.setString(string(r[:k+1]))

			e.originalContent = e.content

			e.setString(string(r[k+1:]))

			e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

			e.prev.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

			e.prev.prev.insertTokenLeft(t2)

			e.modified = true

			return

		}

		if e.prev == nil {

			e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

			e.prev.insertTokenLeft(newTokenOfType(rubyGroupStartToken))

			return

		}

		/*		if e.prev.tokType == rubyParentStartToken {

					return

				}
		*/
	}

	e.insertTokenRight(newTokenOfType(rubyParentStartToken))

	e.insertTokenRight(newTokenOfType(rubyGroupStartToken))

	return
}

func (tok *token) reformgaiji() {

	if tok.tokType != gaijiToken {

		return

	}

	tok.setString(noteStartStr + "※は" + tok.innerString() + noteEndStr)

	tok.tokType = gaijiNoteToken

	tok.modified = true

	n := newTokenOfType(textToken)

	n.setString(referenceMarkStr)

	tok.insertTokenLeft(n)

}

func (tok *token) fixgaiji() {

	switch {

	case tok.tokType == gaijiToken:

		tok.replaceGaiji()

	case tok.tokType == noteToken:

		tk2 := tokenize(tok.innerString())

		hasgaiji := false

		for tk3 := tk2; tk3 != nil; tk3 = tk3.next {

			if tk3.tokType == gaijiToken {
				tk3.fixgaiji()
				hasgaiji = true

			}
		}

		if hasgaiji {
			for e := tk2; e != nil; e = e.next {
				tok.unicodeContent = tok.unicodeContent + e.unicodeString()
				tok.jis0213Content = tok.jis0213Content + e.jis0213String()
			}
		}

	default:

		return
	}

	if tok.tokType == gaijiToken {

		tokenizerLog.Println("gaiji conversion failed:", tok.info())

		tok.reformgaiji()

	}
}

func (tok *token) replaceGaiji() {

	var uni string

	originalNote := tok.innerString()

	j := jisCodeOf(originalNote)

	if j != "" {
		uni, _ = convert(j)

		tok.jis0213Content = uni

		tok.unicodeContent = uni

	} else {

		u := ucode(originalNote)

		if u == "" {
			return
		}

		uni = unicodeOf(u)

		tok.unicodeContent = uni

	}

	switch uni {

	case "《":
		tok.tokType = specialCharToken

	case "》":
		tok.tokType = specialCharToken

	case "［":
		tok.tokType = specialCharToken

	case "］":
		tok.tokType = specialCharToken

	case "〔":
		tok.tokType = specialCharToken

	case "〕":
		tok.tokType = specialCharToken

	case "｜":
		tok.tokType = specialCharToken

	case "＃":
		tok.tokType = specialCharToken

	case "※":
		tok.tokType = specialCharToken

	default:
		tok.tokType = gaijiCharToken
	}

	tok.modified = true

	return
}

// convert everything to nested formatting notes
// to ease parsing of notes. The returned token is new opener.
func (note *token) fixImpliedOpener() {

	if note.tokType != noteToken {
		return
	}

	m := ""

	for _, e := range impliedOpenerMarker {

		if strings.HasSuffix(note.innerString(), e) {

			m = e

			break

		}

	}

	if m == "" {
		return
	}

	if note.next.tokType == rubyStartToken {

		msg := "line " + strconv.Itoa(note.lineNumber()) + " wrong order of ruby and annotation:" + note.prev.String() + note.String() + note.next.String() + note.next.next.String()

		log.Println("WARNING: " + msg)

		e := new(token)

		for e := note.next; e.tokType != rubyEndToken; e = e.next {
		}

		n2 := copyOf(note)

		e.insertTokenRight(n2)

		note.tokType = emptyToken

		return
	}

	switch m {

	case "ルビ":
		fixLeftRuby(note)

	case "注記":
		fixChuki(note)

	case "大きな文字":
		fixFontSize(note)

	case "小さな文字":
		fixFontSize(note)

	default:
		fixopener(note, m)

	}

	return
}

func fixopener(note *token, m string) {

	rs := getRefStrings(note.innerString())

	if rs[0] == "" {
		return
	}

	left := false

	if strings.HasPrefix(strings.TrimPrefix(note.innerString(), refStartStr+rs[0]+refEndStr), "の左に") {
		left = true
	}

	b := new(strings.Builder)

	if left {
		addToStringsBuilder(b, "左に")
	}

	addToStringsBuilder(b, m, formatEndStr)

	note.originalContent = note.content

	note.setInnerString(b.String())

	note.modified = true

	nt := new(token)

	if left {

		nt = newNote("左に" + m)

	} else {

		nt = newNote(m)
	}

	note.addTokenBefore(rs[0], nt)

	return
}

func fixFontSize(note *token) {

	rs := getRefStrings(note.innerString())

	if rs[0] == "" {
		return
	}

	sb := strings.Split(note.innerString(), refEndStr+"は")[1]

	note.originalContent = note.content

	if strings.HasSuffix(sb, "大きな文字") {

		note.setInnerString("大きな文字終わり")

	} else {

		note.setInnerString("小さな文字終わり")

	}

	note.modified = true

	note.addTokenBefore(rs[0], newNote(sb))

	return
}

func fixLeftRuby(note *token) *token {

	rs := getRefStrings(note.innerString())

	if len(rs) != 2 {
		return note
	}

	//	pt := note.prevTextToken()

	b := new(strings.Builder)

	addToStringsBuilder(b, "［＃左に", refStartStr, rs[1], refEndStr, "のルビ付き終わり］")

	note.originalContent = note.content

	note.setString(b.String())

	note.modified = true

	nt := newNote("左にルビ付き")

	note.addTokenBefore(rs[0], nt)

	return note
}

func fixChuki(note *token) *token {

	var left bool

	rs := getRefStrings(note.innerString())

	if len(rs) != 2 {
		return note
	}

	if strings.HasPrefix(strings.TrimPrefix(note.innerString(), refStartStr+rs[0]+refEndStr), "の左に") {
		left = true
	}

	//	pt := note.prevTextToken()

	b := new(strings.Builder)

	if left {
		addToStringsBuilder(b, "左に")
	}

	addToStringsBuilder(b, refStartStr, rs[1], refEndStr, "の注記付き終わり")

	note.originalContent = note.content

	note.setInnerString(b.String())

	note.modified = true

	nt := new(token)

	if left {

		nt = newNote("左に注記付き")

	} else {

		nt = newNote("注記付き")

	}

	note.addTokenBefore(rs[0], nt)

	return note
}

func (t *token) fixParagraphs() {

	start := new(token)
	end := new(token)

	for start = t.firstToken(); start.next != nil; start = end.next {

		switch {

		case start.tokType == endOfLineToken:

			start.tokType = emptyLineToken

			start.modified = true

			end = start

		default:

			end = start.nextTokenOfType(endOfLineToken)

			if end == nil {
				break
			}

			fixParagraph(start, end)

		}

		if end.next == nil {
			break
		}
	}

}

// start should be non-eol and end should be eol
func fixParagraph(start, end *token) {

	switch start.tokType {

	case textToken, rubyGroupStartToken, gaijiCharToken, specialCharToken, accentToken, gaijiImgToken:

		start.insertTokenLeft(newParagraphToken())

		end.tokType = paragraphEndToken

		end.modified = true

		return

	case bibInfoToken:

		if start.String() != "" {
			return
		}

		start.insertTokenRight(newParagraphToken())

		end.tokType = paragraphEndToken

		end.modified = true

		return

	}

	if start.isSectionTitleStart() {
		return
	}

	if start.isFormatOfType(captionMarker) {
		return
	}

	pos := start

	//	if pos.isBlockStartNote() {

	if strings.HasPrefix(pos.innerString(), blockStartStr) {

		for e := pos; e != end; e = e.next {

			if e.isSectionTitleStart() {
				return
			}

			if strings.HasPrefix(pos.innerString(), blockStartStr) {
				//			if e.isBlockStartNote() {
				continue
			}

			if e.tokType != noteToken {

				pos.insertTokenLeft(newParagraphToken())

				end.tokType = paragraphEndToken

				end.modified = true

				return
			}

			if e.next == end {
				return
			}
		}
		return

	}

	if pos.tokType == noteToken {

		for e := pos.next; e != end; e = e.next {

			if e.tokType != noteToken {

				pos.insertTokenLeft(newParagraphToken())

				end.tokType = paragraphEndToken

				end.modified = true

				return

			}
		}
	}
	return
}

func (t *token) fixblockformatting() {

	pos := t

	pos = pos.matchingIndentationCloser()

	if pos.isIndentationStart() {
		pos.insertTokenLeft(newIndentationCloser())
		pos.insertTokenLeft(newEolToken())
	}

	return
}

func (t *token) fixCentering() {

	pos := t

	for ; !pos.next.isPagination(); pos = pos.next {
	}

	pos.next.insertTokenLeft(newTokenOfType(centeringEndToken))

	return
}

func (t *token) fixSectionFormatting() {

	pos := t

	e1 := sectionTitleFormattingStart(pos)

	switch strings.TrimPrefix(pos.innerString(), blockStartStr) {
	case "大見出し":
		e1.insertTokenLeft(newTokenOfType(sectionToken))

	case "中見出し":
		e1.insertTokenLeft(newTokenOfType(subsectionToken))

	case "小見出し":
		e1.insertTokenLeft(newTokenOfType(subsubsectionToken))
	}

	return

}

func (t *token) insertSectionEnds() {

	var open int

	open = 0

	//var opener, end *token

	var end *token
	var openers []*token

	for pos := t.firstToken(); pos.next != nil; pos = pos.next {

		if pos.isBlockStartNote() {
			openers = append(openers, pos)
			open++
		}

		if pos.isBlockClosingNote() {
			open--
			if open < 0 {
				panic(pos.info() + " closing tag without matching opener! ")
			}

			openers = openers[:len(openers)-1]
		}

		if !pos.isSectionStart() {
			continue
		}

		if open != 0 {
			//ensure proper nesting
			end = openers[len(openers)-1].matchingCloserToken()
		} else {
			end = nil
		}

		switch pos.tokType {

		case sectionToken:
			for pos2 := pos.next; pos2 != nil; pos2 = pos2.next {

				if pos2.tokType == sectionToken {
					pos2.insertTokenLeft(newTokenOfType(sectionEndToken))
					break
				}
				if pos2.tokType == bibInfoToken {
					pos2.insertTokenLeft(newTokenOfType(sectionEndToken))
					break
				}
				if pos2 == end {
					pos2.insertTokenLeft(newTokenOfType(sectionEndToken))
					break
				}
				if pos2.next == nil {
					pos2.insertTokenRight(newTokenOfType(sectionEndToken))
					break
				}
			}

		case subsectionToken:
			for pos2 := pos.next; pos2 != nil; pos2 = pos2.next {

				if pos2.tokType == subsectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2.tokType == sectionEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2.tokType == sectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2.tokType == bibInfoToken {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}
				if pos2 == end {
					pos2.insertTokenLeft(newTokenOfType(subsectionEndToken))
					break
				}

				if pos2.next == nil {
					pos2.insertTokenRight(newTokenOfType(subsectionEndToken))
					break
				}

			}

		case subsubsectionToken:
			for pos2 := pos.next; pos2 != nil; pos2 = pos2.next {

				if pos2.tokType == subsubsectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == subsectionEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == subsectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == sectionEndToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == sectionToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.tokType == bibInfoToken {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2 == end {
					pos2.insertTokenLeft(newTokenOfType(subsubsectionEndToken))
					break
				}
				if pos2.next == nil {
					pos2.insertTokenRight(newTokenOfType(subsubsectionEndToken))
					break
				}

			}

		}
	}
}

func newIndentationCloser() *token {

	t := newNote("ここで字下げ終わり")

	return t

}

func sectionTitleFormattingStart(t *token) (s *token) {

	r := t

	for s = t.prev; ; s = s.prev {

		if s == nil {
			break
		}

		if s.tokType == endOfLineToken {
			continue
		}

		if s.tokType == emptyLineToken {
			continue
		}

		if strings.HasPrefix(s.innerString(), blockStartStr) {
			r = s
			continue
		}

		if s.isFormatOfType(centeringMarker) {
			r = s
			continue
		}

		r = s.next
		break

	}

	for s = r; s.tokType == emptyLineToken; s = s.next {
	}

	if s == nil {

		s = t.firstToken()

	}

	return s

}

func (t *token) fixaccent() {

	if t.tokType != accentToken {
		return
	}

	t.unicodeContent = convertAccent(t.innerString())

	t.jis0213Content = t.unicodeContent

	t.modified = true

	if len([]rune(t.unicodeContent)) == len([]rune(t.innerString())) {
		t.tokType = textToken
		return

	}

}

func (t *token) decorationLeft() bool {

	if t.tokType != noteToken {
		return false
	}

	return strings.HasPrefix(t.innerString(), "左に")

}

func (t *token) fixFigures() {

	e := t

	if e.nextSignificantToken().isCaption() {

		e.insertTokenLeft(newTokenOfType(figureStartToken))

		n := newTokenOfType(figureEndToken)

		e.nextSignificantToken().matchingCloserToken().insertTokenRight(n)

		return

	}

}

func (t *token) fixBibInfo() {

	if t.tokType != bibInfoToken {
		return
	}

	pos := new(token)

	for pos = t.next; pos != nil; pos = pos.next {

		if pos.tokType == bibInfoToken {

			pos.tokType = textToken

			pos.originalContent = pos.content

			pos.setString(strings.TrimPrefix(pos.String(), "\n"))

			pos.modified = true

			pos.insertTokenLeft(newTokenOfType(endOfLineToken))

		}
	}

	for pos = t.lastToken(); pos.tokType != endOfLineToken; pos = pos.prev {
	}

	pos.insertTokenLeft(newTokenOfType(bibInfoEndToken))

	if t.innerString() == noteStartStr+mainTextEndStr+noteEndStr {
		return
	}

	t.tokType = textToken

	t.originalContent = t.content

	t.setString(strings.TrimPrefix(t.String(), "\n"))

	t.modified = true

	t.insertTokenLeft(newTokenOfType(bibInfoToken))

	t.prev.insertTokenLeft(newTokenOfType(endOfLineToken))

	return
}

func (t *token) fixDocument() {

	for pos := t.firstToken(); pos != nil; pos = pos.next {

		switch {

		case pos.tokType == textToken:

			continue

		case pos.tokType == rubyStartToken:

			continue

		case pos.tokType == rubyEndToken:

			continue

		case pos.tokType == rubyParentStartToken:

			continue

		case pos.tokType == rubyParentEndToken:

			continue

		case pos.innerString() == "ページの左右中央":

			pos.fixCentering()

		case pos.isSectionTitleStart():

			pos.fixSectionFormatting()

		case pos.isIndentationStart():

			pos.fixblockformatting()

		case pos.isImage():

			pos.fixFigures()
		}

	}

}

func (t *token) fixEmptyText() {

	for pos := t.firstToken(); pos != nil; pos = pos.next {

		if pos.tokType != textToken {
			continue
		}

		if pos.String() == "" {
			pos = pos.next
			pos.prev.remove()
		}
	}

	return
}

func (t *token) insertAozoraBookMarker() {

	e := new(token)

	for e = t.firstToken(); e.tokType != emptyLineToken; e = e.nextLine() {
		if e.nextLine() == nil {
			return
		}
	}

	if e.lineNumber() < 3 {
		return
	}

	e.insertTokenRight(newTokenOfType(mainTextStartToken))

	for ; e.tokType != bibInfoToken; e = e.next {

		if e.next == nil {
			t.lastToken().insertTokenRight(newTokenOfType(mainTextEndToken))
			return
		}
	}

	e.insertTokenLeft(newTokenOfType(mainTextEndToken))

	return

}

func (t *token) fixBlockFormat() {

	if t.tokType != noteToken {
		return
	}

	if strings.HasPrefix(t.innerString(), blockStartStr) {

		if t.next.tokType == endOfLineToken {
			return
		}

		if o_strict {
			panic("ERROR: line " + strconv.Itoa(t.lineNumber()) + " block start annotation should be on own line.")
			return
		}

		log.Println("WARNING: line", t.lineNumber(), "block start annotation should be on own line. Fixed.")

		t.insertTokenRight(newTokenOfType(endOfLineToken))
		return
	}

	if strings.HasPrefix(t.innerString(), blockEndStr) {

		if t.prev.tokType == endOfLineToken {
			return
		}

		if o_strict {
			panic("ERROR: line " + strconv.Itoa(t.lineNumber()) + " block end annotation should be on own line.")
			return
		}

		log.Println("WARNING: line", t.lineNumber(), "block end annotation should be on own line. Fixed.")

		t.insertTokenLeft(newTokenOfType(endOfLineToken))
	}
	return
}
