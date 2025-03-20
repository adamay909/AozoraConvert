package aozoratext

import (
	"strings"
)

func fixfunc() func(t *token) *token {

	var fix func(t *token) *token

	switch {

	case o_raw:
		fix = func(t *token) *token {

			t.fixParagaph()

			return t

		}

	default:

		fix = func(t *token) *token {

			for tok := t.firstToken(); tok != nil; tok = tok.next {

				tok.fixkunoji()
				tok.fixaccent()
				tok.fixgaiji()
				tok.fixruby()

			}

			for tok := t.lastToken(); tok != nil; tok = tok.prev {

				tok.fixImpliedCloser()
				tok.fixImpliedOpener()

			}

			t.fixParagaph()

			return t
		}
	}

	return fix

}

func (t *token) fixline(processor func(*token) *token) *token {

	return processor(t)

}

func (t *token) fixImpliedCloser() {

	if t.tokType != noteToken {
		return
	}

	if strings.HasPrefix(t.innerString(), blockStartStr) {
		return
	}
	m := isMarker(t.innerString(), impliedCloserMarker)

	if m == emptyStr {
		return
	}

	sbs := new(strings.Builder)

	sbe := new(strings.Builder)

	if t.prev == nil {

		addToStringsBuilder(sbs, blockStartStr, t.innerString())

		addToStringsBuilder(sbe, blockEndStr, m, formatEndStr)

		t.setInnerString(sbs.String())

		t.lastToken().insertTokenRight(newNote(sbe.String()))

		return

	}

	addToStringsBuilder(sbs, t.innerString())

	addToStringsBuilder(sbe, m, formatEndStr)

	if t.hasCloser(sbs.String(), sbe.String()) {
		return
	}

	t.setInnerString(sbs.String())

	t.lastToken().insertTokenRight(newNote(sbe.String()))

	return

}

func (t *token) fixkunoji() {

	if t.tokType != kunojiToken {
		return
	}

	if t.content == kunojiStr {

		t.altContent = kunojiStrU

		return
	}

	t.altContent = kunojiDakuStrU

}

func (t *token) fixruby() {

	if t.tokType != rubyStartToken {

		return
	}

	e := new(token)

	for e = t.next; e != nil; e = e.next {

		if e.tokType == rubyEndToken {
			e.insertTokenRight(newTokenOfType(rubyParentEndToken))
			break
		}
	}

	if e == nil {

		t.lastToken().insertTokenRight(newTokenOfType(rubyParentEndToken))

	}

	if !t.rubyParentExplicit() {

		t.insertRubyParentStart()

	}

	return

}

func (t *token) rubyParentExplicit() bool {

	for e := t.prev; e != nil; e = e.prev {

		if e.tokType == textToken {
			continue
		}

		if e.tokType == rubyParentStartToken {
			return true
		}

		if e.tokType == gaijiNoteToken {
			continue
		}

		if e.tokType == accentToken {
			continue
		}

		break
	}

	return false
}

func (t *token) insertRubyParentStart() {

	e := new(token)

	if t.prev.tokType == accentToken {

		t.prev.insertTokenLeft(newTokenOfType(rubyParentStartToken))

		return

	}

	var ref charTypeID

	var r []rune

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == rubyEndToken {

			for ; e.tokType != rubyStartToken; e = e.prev {
			}

			e = e.prev

		}

		if e.tokType != textToken {
			continue
		}

		r = []rune(e.innerString())

		ref = charType(r[len(r)-1])

		break
	}

	k := 0

	for e = t.prev; e != nil; e = e.prev {

		if e.tokType == rubyParentEndToken {
			break
		}

		if e.tokType == noteToken {
			break
		}

		if e.tokType != textToken {
			continue
		}

		r = []rune(e.innerString())

		for k = len(r) - 1; ref == charType(r[k]); k-- {

			if k == 0 {

				break
			}
		}

		if charType(r[k]) != ref {

			t2 := newTokenOfType(textToken)

			t2.setString(string(r[:k+1]))

			e.setString(string(r[k+1:]))

			e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

			e.prev.insertTokenLeft(t2)

			return

		}

		if e.prev == nil {

			e.insertTokenLeft(newTokenOfType(rubyParentStartToken))

			return

		}

		if e.prev.tokType == rubyParentStartToken {

			return

		}

	}

	e.insertTokenRight(newTokenOfType(rubyParentStartToken))

	return
}

func (tok *token) reformgaiji() {

	if tok.tokType != gaijiToken {

		return

	}

	tok.setString(noteStartStr + "※は" + tok.innerString() + noteEndStr)

	tok.tokType = gaijiNoteToken

	n := newTokenOfType(textToken)

	n.setString(referenceMarkStr)

	tok.insertTokenLeft(n)

	n.tokType = textToken

}

func (tok *token) fixgaiji() {

	switch {

	case tok.tokType == gaijiToken:

		tok.replaceGaiji()

	case tok.tokType == noteToken:

		tk2 := tokenize(tok.innerString())

		for tk3 := tk2; tk3 != nil; tk3 = tk3.next {

			tk3.fixgaiji()

		}

		tok.setInnerString(tk2.StringAll())

	default:

		return
	}

	if tok.tokType == gaijiToken {

		tokenizerLog.Println("gaiji conversion failed:", tok.String())

		tok.reformgaiji()

	}
}

func (tok *token) replaceGaiji() {

	var uni string

	originalNote := tok.innerString()

	j := jisCodeOf(originalNote)

	if j != "" {
		uni, _ = convert(j)

	} else {

		u := ucode(originalNote)

		if u == "" {
			return
		}

		uni = unicodeOf(u)

	}

	tok.content = tok.content

	tok.altContent = uni

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

	return
}

// convert everything to nested formatting notes
// to ease parsing of notes. The returned token is new opener.
func (note *token) fixImpliedOpener() {

	if note.tokType != noteToken {
		return
	}

	m := note.getFormattingMarker()

	switch m {

	case "":
		return

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

	note.setInnerString(b.String())

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

	//	sb := strings.TrimPrefix(note.innerString(), refStartStr+rs[0]+refEndStr+"は")

	sb := strings.Split(note.innerString(), refEndStr+"は")[1]
	//	sb := note.innerString()

	if strings.HasSuffix(sb, "大きな文字") {

		note.setInnerString("大きな文字終わり")

	} else {

		note.setInnerString("小さな文字終わり")

	}

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

	note.setString(b.String())

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

	note.setInnerString(b.String())

	nt := new(token)

	if left {

		nt = newNote("左に注記付き")

	} else {

		nt = newNote("注記付き")

	}

	note.addTokenBefore(rs[0], nt)

	return note
}

func (t *token) fixParagaph() {

	pos := t.firstToken()

	if pos.tokType == emptyLineToken {
		return
	}

	if pos.tokType == bibInfoToken {
		pos.lastToken().insertTokenRight(newLineBreakToken())
		return
	}

	if pos.isFormatOfType(captionMarker) {
		/*
			if !strings.HasPrefix(pos.innerString(), blockStartStr) {

				pos.insertTokenRight(newTokenOfType(paragraphToken))

				pos.matchingCloserToken().insertTokenLeft(newLineBreakToken())
		*/
		return
		//		}
	}

	for ; pos.isBlockStartNote(); pos = pos.next {

		if pos.isSectionTitleStart() {
			return
		}

		if pos.next == nil {
			break
		}
	}

	if pos.tokType == noteToken && pos.next == nil {
		return
	}

	if pos.lastToken().isBlockClosingNote() {

		if pos.lastToken().prev != nil {

			pos.insertTokenLeft(newParagraphToken())
			pos.lastToken().prev.insertTokenRight(newLineBreakToken())

		}

		return

	}

	pos.insertTokenLeft(newParagraphToken())
	pos.lastToken().insertTokenRight(newLineBreakToken())

	return
}

func (t *token) length() int {

	c := 0

	for pos := t.firstToken(); ; pos = pos.next {
		c++
		if pos.next == nil {
			break
		}
	}

	return c

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

	pos.next.insertTokenLeft(newNote("ページの左右中央終わり"))

	return
}

func (t *token) fixSectionFormatting() {

	pos := t

	e1 := sectionTitleFormattingStart(pos)

	switch pos.innerString() {
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
			openers = openers[:len(openers)-1]
			open--
		}

		if !pos.isSectionStart() {
			continue
		}

		if len(openers) > 0 {
		}

		if open != 0 {
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

		if s.isBlockFormatterStart() {
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

	return r

}

func (t *token) isFormatted() bool {

	if t.prev == nil {
		return false
	}

	for s := t.prev; ; s = s.prev {

		if s.tokType == endOfLineToken {
			continue
		}

		if s.tokType == emptyLineToken {
			continue
		}

		if s.isBlockFormatterStart() {
			return true
		}

		if s.isFormatOfType(centeringMarker) {
			return true
		}

		break

	}

	return false
}

func (t *token) fixaccent() {

	if t.tokType != accentToken {
		return
	}

	t.altContent = convertAccent(t.innerString())

	if len([]rune(t.altContent)) == len([]rune(t.innerString())) {
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

	/*	for e := t.firstToken(); e.next != nil; e = e.next {

		if !e.isImage() {
			continue
		}
	*/
	e.insertTokenLeft(newTokenOfType(figureStartToken))

	n := newTokenOfType(figureEndToken)

	if !e.next.isCaption() {

		e.insertTokenRight(n)
		if n.next.tokType == lineBreakToken {
			n.next.remove()
		}

		//continue
		return
	}

	e.next.matchingCloserToken().insertTokenRight(n)

	// }
}

func (t *token) fixBibInfo() {

	e := t

	if e.innerString() == "本文終わり" {

		e.tokType = bibInfoToken

		return
	}

	if e.tokType == bibInfoToken {

		e.tokType = textToken

		for pos := e.prev; pos != nil; pos = pos.prev {

			if pos.tokType == bibInfoToken {
				return
			}

			if pos.tokType == emptyLineToken {
				continue
			}

			break

		}

		e.insertTokenLeft(newTokenOfType(bibInfoToken))

		e.insertTokenLeft(newParagraphToken())

	}

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

		case pos.tokType == bibInfoToken:

			pos.fixBibInfo()

		case pos.tokType == noteToken && pos.innerString() == "本文終わり":

			pos.fixBibInfo()
		}

	}

}
