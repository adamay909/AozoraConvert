package mobi

// Chunks produces a list of chunks from one or more strings.
//
// In the resulting list of chunks, each chunk exactly corresponds the
// given string with the same index.  In almost all cases, this is the
// preferred method of converting KF8 HTML into a list of chunks.
func Chunks(list ...string) []Chunk {
	result := make([]Chunk, 0)
	for _, s := range list {
		result = append(result, Chunk{
			Body: s,
		})
	}
	return result
}

func Split(s string) (list []string) {

	var p []rune
	var l int
	for _, c := range s {

		l = l + len(string(c))

		if l > 4096 {
			list = append(list, string(p))
			p = nil
			p = append(p, c)
			l = len(string(p))
		} else {
			p = append(p, c)
		}
	}
	if len(p) > 0 {
		list = append(list, string(p))
	}
	return
}
