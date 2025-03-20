package aozoratext

import (
	"log"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

// ToSJIS  converts text to ShiftJIS encoding. Also converts to DOS line endings.
func ToSJIS(text string) string {

	enc := japanese.ShiftJIS.NewEncoder()

	out, err := enc.String(text)

	if err != nil {

		log.Println("errors while converting to ShiftJIS.")

		log.Println(err)

	}

	return strings.ReplaceAll(out, "\n", "\r\n")

}

// ToUTF8 converts from ShiftJIS to UTF8. Also converts to unix line endings. Note that input is []byte, not string.
func ToUTF8(data []byte) string {

	dec := japanese.ShiftJIS.NewDecoder()

	out, err := dec.Bytes(data)

	if err != nil {

		log.Println("errors while converting to UTF-8.")

		log.Println(err)

		return ""

	}

	return strings.ReplaceAll(string(out), "\r\n", "\n")

}
