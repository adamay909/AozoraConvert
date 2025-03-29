package aozoratext

import (
	"strings"
)

// first item in slice is metadata, second is main text
func prepText(text string) []string {

	var rt []string

	idx1 := strings.Index(text, markupNoteStartStr)

	if idx1 == -1 {

		rt = append(rt, "")

		rt = append(rt, text)

		return rt
	}

	idx2 := strings.Index(text[idx1+len(markupNoteStartStr):], markupNoteEndStr)

	if idx2 == -1 {

		rt = append(rt, "")

		rt = append(rt, text)

		return rt

	}

	rt = append(rt, text[:idx1])

	rt = append(rt, text[idx1+len(markupNoteStartStr)+idx2+len(markupNoteEndStr):])

	return rt
}

func _isAozoraText(text string) bool {

	idx1 := strings.Index(text, markupNoteStartStr)

	if idx1 == -1 {
		return false
	}

	idx2 := strings.Index(text[idx1+len(markupNoteStartStr):], markupNoteEndStr)

	if idx2 == -1 {
		return false
	}

	return true

}

func aztextMainStart(text string) int {

	start := 0

	lines := strings.Split(text, "\n")

	for i := range lines {

		if lines[i] == "" {

			start = i

			break

		}

	}

	if start == 0 {

		return start

	}

	idx1 := strings.Index(text, markupNoteStartStr)

	if idx1 == -1 {

		return start + 1

	}

	if strings.Index(text[idx1+len(markupNoteStartStr):], markupNoteEndStr) == -1 {

		return len(lines)

	}

	counter := 0

	for i, l := range strings.Split(text, "\n") {

		if l == markupNoteEndStr {
			counter++

			if counter == 2 {
				return i + 1

			}
		}

	}

	return 0
}
