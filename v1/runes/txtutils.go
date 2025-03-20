package runes

import (
	"strings"
)

//simpleElement is used for holding simple information about
//chunks of text.
type simpleElement struct {
	raw   Runes
	match bool
}

//content returns the raw value of c.
func (c simpleElement) content() Runes {
	return c.raw
}

//isMatch returns the match status of c.
func (c simpleElement) isMatch() bool {
	return c.match
}

func (c simpleElement) String() string {
	return c.raw.String()
}

//SplitIntoBlocks splits src into two types of chunks: those
//that match the block as defined by first and last, and those
//that do not. The result is returned as a slice of Runes. A block will
//always terminate at the first instance of last so
//none of the block related functions in this package will work
//with nested blocks.
func SplitIntoBlocks(src, first, last Runes) (out []Runes) {

	chunks := findBlocks(src, first, last)

	for _, c := range chunks {
		out = append(out, c.content())
	}
	return out
}

//Replace replaces the first n instances of matchString in src
//with replaceString. If n<0, it replaces all. Returns src if n=0.
func Replace(src, matchString, replaceString Runes, n int) (out Runes) {

	return Runes(strings.Replace(src.String(), matchString.String(), replaceString.String(), n))
}

//ReplaceAll replaces all instances of matchString in src with replaceString.
func ReplaceAll(src, matchString, replaceString Runes) (out Runes) {

	return Replace(src, matchString, replaceString, -1)

}

//ReplaceBlocks replaces the first n blocks defined by first-last pair
//by replaceString. If n=-1, all are replaced. If n=0, returns src.
func ReplaceBlocks(src, first, last, replaceString Runes, n int) (out Runes) {

	simpleHandler := func(r Runes) Runes {
		return replaceString
	}

	return ReplaceBlocksFunc(src, first, last, simpleHandler, n)
}

//ReplaceBlocksAll replaces all blocks.
func ReplaceBlocksAll(src, first, last, replaceString Runes) (out Runes) {

	return ReplaceBlocks(src, first, last, replaceString, -1)

}

//ReplaceBlocksFunc replaces the content of blocks defined by first
//and last with the output of handler.
func ReplaceBlocksFunc(src, first, last Runes, handler func(Runes) Runes, n int) (out Runes) {
	var i, j, count int

	if n == 0 {
		return src
	}

	chunks := findBlocks(src, first, last)

	for _, c := range chunks {
		if c.isMatch() {
			count++
		}
	}

	if n < 0 || n > count {
		n = count
	}

	for i = 0; i < len(chunks); i++ {
		if !chunks[i].isMatch() {
			out = append(out, chunks[i].content()...)
		}
		if chunks[i].isMatch() {
			rp := handler(chunks[i].content())
			if len(rp) > 0 {
				out = append(out, rp...)
			}
			j++
			if j == n {
				break
			}
		}
	}
	for i++; i < len(chunks); i++ {
		out = append(out, chunks[i].content()...)
	}

	return out
}

//ReplaceBlocksFuncAll replaces all blocks with output of handler.
func ReplaceBlocksFuncAll(src, first, last Runes, handler func(Runes) Runes) (out Runes) {

	return ReplaceBlocksFunc(src, first, last, handler, -1)
}

//find first position of search string.
func findText(matchString, line Runes) (start, end int) {
	start, end = -1, -1
	start = Index(line, matchString)
	if start == -1 {
		return start, end
	}
	if len(line[start:]) < 2 {
		return start, end
	}
	end = start + len(matchString) - 1

	return start, end

}

//find breaks src into simpleElements that either match or
//don't match the matchString.
func find(src, matchString Runes) (out []simpleElement) {
	var idx, edx int

	for i := 0; i < len(src); {
		var (
			tok simpleElement
		)

		idx, edx = findText(matchString, src[i:])
		if idx == -1 {
			idx = len(src)
			edx = len(src)
		} else {
			idx = idx + i
			edx = edx + i
		}
		if idx != 0 {
			tok.raw = src[i:idx]
			tok.match = false
			out = append(out, tok)
		}
		if edx == len(src) {
			i = edx
			continue
		}

		tok.raw = src[idx : edx+1]
		tok.match = true
		out = append(out, tok)
		i = edx + 1
	}
	return out
}

//findBlocks breaks src into two types of chunks: plain and elements
//defined by first and last pair. Useful for manipulating
//particular type of element.
//The return value is []token, where token.typus is 1 iff.
//it is an element of the defined type.
func findBlocks(src, first, last Runes) (out []simpleElement) {

	for i := 0; i < len(src); {
		var tok simpleElement

		idx, edx := findBlock(first, last, src[i:])
		if idx == -1 {
			tok.raw = src[i:]
			tok.match = false
			out = append(out, tok)
			i = len(src)
			continue
		} else {
			idx = idx + i
			edx = edx + i
		}
		if idx > i {
			tok.raw = src[i:idx]
			tok.match = false
			out = append(out, tok)
		}

		tok.raw = src[idx : edx+1]
		tok.match = true
		out = append(out, tok)
		i = edx + 1
	}
	return out
}

//find first block defined by first and last and return start
//end index.
func findBlock(first, last Runes, line Runes) (start, end int) {
	start, end = -1, -1
	if len(first) == 0 {
		return
	}
	if len(last) == 0 {
		return
	}
	start = Index(line, first)
	if start == -1 {
		return start, end
	}

	end = Index(line[start+len(first):], last)

	if end <= -1 {
		start = -1
		return
	}
	return start, start + end + len(first) + len(last) - 1

}
