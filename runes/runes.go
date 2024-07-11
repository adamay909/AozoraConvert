package runes

import (
	"errors"
	"io"
	"strings"
)

//Runes is an alias for []rune.
type Runes []rune

//String returns the content of r as a string.
//Implements stringer interface.
func (r Runes) String() string {
	return string(r)
}

//Bytes returns r as slice of bytes.
func (r Runes) Bytes() []byte {
	return []byte(r.String())
}

//Join joins the runes in slice using sep as separator.
//Note that sep is a string.
func Join(slice []Runes, sep Runes) (out Runes) {

	var sliceStr []string

	for _, s := range slice {
		sliceStr = append(sliceStr, s.String())
	}

	return Runes(strings.Join(sliceStr, sep.String()))

}

//Split splits r into a sub-slices using sep as the separator.
func Split(r Runes, sep Runes) (out []Runes) {

	s := strings.Split(r.String(), sep.String())

	for _, i := range s {
		out = append(out, Runes(i))
	}
	return out
}

//HasPrefix returns true if r starts with pre.
func HasPrefix(r Runes, pre Runes) bool {
	return strings.HasPrefix(r.String(), pre.String())
}

//Equal returns true if s1 and s2 are identical.
func Equal(s1, s2 Runes) bool {

	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if string(s1[i]) != string(s2[i]) {
			return false
		}
	}
	return true

}

//Index returns the starting index of s within r.
//Returns -1 if s is not contained in r.
func Index(r Runes, s Runes) int {

	for i := 0; i <= len(r)-len(s); i++ {
		if Equal(r[i:i+len(s)], s) {
			return i
		}
	}
	return -1
}

//IndexFunc returns the index of first rune in r
//that satisfies condition given by f.
func IndexFunc(r Runes, f func(rune) bool) int {

	for i := 0; i < len(r); i++ {
		if f(r[i]) {
			return i
		}
	}
	return -1
}

//IndexAny returns the first index of any of the
//runes appearing in s.
func IndexAny(r Runes, s Runes) int {

	f := func(r rune) bool {
		for _, c := range s {
			if string(r) == string(c) {
				return true
			}
		}
		return false
	}

	return IndexFunc(r, f)
}

//LastIndex returns the idenx of the last instance of
//s within r.
func LastIndex(r, s Runes) int {

	for i := len(r) - 1; i >= 0; i-- {
		if Index(r[i:], s) == 0 {
			return i
		}
	}
	return -1
}

//LastIndexFunc returns the index of first rune in r
//that satisfies condition given by f.
func LastIndexFunc(r Runes, f func(rune) bool) int {

	for i := len(r) - 1; i >= 0; i-- {
		if f(r[i]) {
			return i
		}
	}
	return -1
}

//LastIndexAny returns the last index of any of the
//runes appearing in s.
func LastIndexAny(r Runes, s Runes) int {

	f := func(r rune) bool {
		for _, c := range s {
			if string(r) == string(c) {
				return true
			}
		}
		return false
	}

	return LastIndexFunc(r, f)
}

//TrimSpace removes all leading and trailing white spaces.
func TrimSpace(r Runes) Runes {

	return Runes(strings.TrimSpace(string(r)))

}

//Chunk splits r into a Runes of length at most n and
//returns the result as []Runes
func Chunk(r Runes, n int) (out []Runes) {

	remain := r

	for len(remain) >= n {
		out = append(out, remain[:n])
		if len(remain) > n {
			remain = remain[n:]
		} else {
			remain = Runes{}
		}
	}
	if len(remain) > 0 {
		out = append(out, remain)
	}
	return out
}

//Contains returns true if substr occurs in s.
func Contains(s, substr Runes) bool {
	return strings.Contains(s.String(), substr.String())
}

//ContainsAny returns true if any rune in substr occurs in s.
func ContainsAny(s, substr Runes) bool {
	return strings.ContainsAny(s.String(), substr.String())
}

//Fields splits src into fields where fields are separated
//by white spaces.
func Fields(src Runes) (out []Runes) {

	s := strings.Fields(src.String())

	for _, i := range s {
		out = append(out, Runes(i))
	}

	return out
}

//HasSuffix returns true if s has suffix.
func HasSuffix(s, suffix Runes) bool {

	return strings.HasSuffix(s.String(), suffix.String())
}

//Map transforms s according to mapping.
func Map(mapping func(rune) rune, s Runes) Runes {

	return Runes(strings.Map(mapping, s.String()))

}

//TrimFunc trims s according to f.
func TrimFunc(s Runes, f func(rune) bool) Runes {

	return Runes(strings.TrimFunc(s.String(), f))

}

//Builder is for building Runes.
type Builder struct {
	buf Runes
}

//NewBuilder returns a pointer to a new builder.
func NewBuilder() *Builder {
	return new(Builder)
}

//Write implements the io.Writer interface
func (b *Builder) Write(p []byte) (n int, err error) {
	b.buf = append(b.buf, Runes(string(p))...)
	return len(p), nil
}

//WriteRunes appends r to b. Returns the length of n and error
//(currently always nil).
func (b *Builder) WriteRunes(r Runes) (n int, err error) {
	b.buf = append(b.buf, r...)
	return len(r), nil
}

//WriteString appends s to b by first converting it to Runes.
//Returns the number of runes written and error (always nil
//at the moment).
func (b *Builder) WriteString(s string) (n int, err error) {
	b.buf = append(b.buf, Runes(s)...)
	return len(Runes(s)), err
}

//Runes reads out the accumulated runes.
func (b *Builder) Runes() Runes {
	return b.buf
}

//String implements the Stringer interface.
func (b *Builder) String() string {
	return b.buf.String()
}

//Bytes returns the content of b in bytes.
func (b *Builder) Bytes() []byte {
	return b.buf.Bytes()
}

//Reset resets b.
func (b *Builder) Reset() {
	b.buf = nil
	return
}

//Reader is an io.Reader reading from s.
type Reader struct {
	s    Runes
	ridx int //current reading index of runes
	bidx int //current reading index of bytes
}

//NewReader returns a new Reader that reads from s.
func NewReader(s Runes) *Reader {
	return &Reader{s, 0, 0}
}

//Size returns the size of the underlying Runes.
func (r *Reader) Size() int64 { return int64(len(r.s)) }

//Read implements io.Reader interface. It reads one rune at
//a time so don't expect to be able to read a sequence of
//single bytes. Most of the time, you probaly want ReadRune.
func (r *Reader) Read(b []byte) (n int, err error) {
	var (
		buf, c  string
		ch      rune
		counter int
	)

	if r.ridx >= len(r.s) {
		return 0, io.EOF
	}
	for {
		counter++
		//c = ""
		ch, _, err = r.ReadRune()
		if err != nil {
			break
		}
		c = string(ch)
		if len(c)+len(buf) > len(b) {
			r.Seek(-1, io.SeekCurrent)
			break
		}
		buf = buf + c
	}
	n = copy(b, []byte(buf))
	if n == 0 && counter == 1 {
		r.Seek(-1, io.SeekCurrent)
		err = io.ErrShortBuffer
	}
	return
}

//ReadRune implements io.RuneReader interface. For compatibility
//only. Most of the time, you probaly want ReadRune.
func (r *Reader) ReadRune() (ch rune, size int, err error) {
	if r.ridx >= len(r.s) {
		return 0, 0, io.EOF
	}
	ch = r.s[r.ridx]
	size = len(string(r.s[r.ridx]))
	r.ridx++
	return
}

//ReadRuneAt returns rune at r[off]. If off is outside the range of r,
//returns appropriate error.
func (r *Reader) ReadRuneAt(off int) (ch rune, size int, err error) {
	if off < 0 {
		err = errors.New("runes.Reader.ReadRunesAt: negative offset")
		return
	}

	if off >= len(r.s) {
		err = io.EOF
		return
	}
	ch = r.s[off]
	size = len(string(r.s[r.ridx]))
	return
}

// Seek implements the io.Seeker interface.
func (r *Reader) Seek(offset int, whence int) (int, error) {
	var abs int
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.ridx + offset
	case io.SeekEnd:
		abs = len(r.s) - 1 - offset
	default:
		return 0, errors.New("strings.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("strings.Reader.Seek: negative position")
	}
	r.ridx = abs
	return abs, nil
}
