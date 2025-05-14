package aozoraconvert

import (
	"errors"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

// ToSJIS  converts text to ShiftJIS encoding. Also converts to DOS line endings.
func ToSJIS(text string) string {

	enc := japanese.ShiftJIS.NewEncoder()

	out, err := enc.String(text)

	if err != nil {

		err = errors.Join(errors.New("errors while converting to ShiftJIS:"), err)

		msglog.Println(err)

	}

	return strings.ReplaceAll(out, "\n", "\r\n")

}

// ToUTF8 converts from ShiftJIS to UTF8. Also converts to unix line endings. Note that input is []byte, not string.
func ToUTF8(data []byte) string {

	dec := japanese.ShiftJIS.NewDecoder()

	outdata, err := dec.Bytes(data)

	if err != nil {

		err = errors.Join(errors.New("errors while converting to UTF-8."), err)

		msglog.Println(err)

		return ""

	}

	out := string(outdata)

	if strings.Index(out, "\r\n") == -1 {

		if strings.Index(out, "\n") != -1 {
			return out
		}

		return strings.ReplaceAll(out, "\r", "\n")
	}

	return strings.ReplaceAll(out, "\r\n", "\n")

}
